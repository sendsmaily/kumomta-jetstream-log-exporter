package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	v1 "github.com/sendsmaily/kumomta-jetstream-log-exporter/http/v1"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tests")
}

var (
	consumer jetstream.Consumer
	handler  echo.HandlerFunc
)

var _ = BeforeEach(func(ctx SpecContext) {
	dir := Must(os.MkdirTemp(os.TempDir(), "nats-*"))
	DeferCleanup(os.RemoveAll, dir)

	srv := server.New(&server.Options{
		ServerName: lo.RandomString(16, lo.AlphanumericCharset),
		Host:       "127.0.0.1",
		Port:       -1,
		Debug:      true,
		JetStream:  true,
		StoreDir:   dir,
	})

	srv.Start()
	DeferCleanup(srv.WaitForShutdown)
	DeferCleanup(srv.Shutdown)

	nc := Must(nats.Connect(srv.ClientURL()))
	DeferCleanup(nc.Close)

	js := Must(jetstream.New(nc))

	stream := Must(js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "LOG_RECORDS",
		Subjects: []string{"log.test"},
	}))

	consumer = Must(stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		Name: "log-test",
	}))

	handler = Must(v1.CreateHandler("log.test", srv.ClientURL()))
})

func ExecuteRequest(ctx SpecContext, path string) *httptest.ResponseRecorder {
	GinkgoHelper()

	payload, err := os.Open(path)
	Expect(err).NotTo(HaveOccurred())

	req := httptest.NewRequest(http.MethodPost, "/v1/record", payload).WithContext(ctx)
	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()

	Expect(handler(echo.New().NewContext(req, rec))).To(Succeed())

	return rec
}

func Must[T any](t T, err error) T {
	GinkgoHelper()
	if err != nil {
		Fail(err.Error())
	}
	return t
}
