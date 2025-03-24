package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	v1 "github.com/sendsmaily/kumomta-jetstream-log-exporter/http/v1"
	"github.com/spf13/cobra"
)

var (
	version = "test"
	commit  = "unknown"

	listenAddr      string
	natsURL         string
	natsTLSKeyFile  string
	natsTLSCertFile string
	natsTLSCAFile   string
)

func main() {
	cmd := &cobra.Command{
		Use:   "exporter",
		Short: "Export KumoMTA Log Records to JetStream",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var opts []nats.Option

			if !lo.IsEmpty(natsTLSKeyFile) || !lo.IsEmpty(natsTLSCertFile) {
				opts = append(opts, nats.ClientCert(natsTLSCertFile, natsTLSKeyFile))
			}

			if !lo.IsEmpty(natsTLSCAFile) {
				opts = append(opts, nats.RootCAs(natsTLSCAFile))
			}

			handler, err := v1.CreateHandler(args[0], natsURL, opts...)
			if err != nil {
				return err
			}

			e := echo.New()

			e.Use(middleware.Recover())

			e.POST("/v1/record", handler)

			return e.Start(listenAddr)
		},
	}

	cmd.Flags().StringVar(&listenAddr, "listen", "127.0.0.1:8080", "Listen address")
	cmd.Flags().StringVar(&natsURL, "nats-url", "nats://localhost:4222", "NATS cluster's URL")
	cmd.Flags().StringVar(&natsTLSKeyFile, "nats-tls-key", "", "Path to the NATS' TLS key file")
	cmd.Flags().StringVar(&natsTLSCertFile, "nats-tls-cert", "", "Path to the NATS' TLS certificate file")
	cmd.Flags().StringVar(&natsTLSCAFile, "nats-tls-ca-cert", "", "Path to the NATS' CA certificate file")

	cmd.Printf("Starting the exporter (version=%s, commit=%s)\n", version, commit)

	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}
