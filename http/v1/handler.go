package v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// CreateHandler creates a KumoMTA webhook handler.
func CreateHandler(subject, natsURL string, opts ...nats.Option) (echo.HandlerFunc, error) {
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("error acquiring JetStream context: %w", err)
	}

	return func(c echo.Context) error {
		payload, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error reading request body: %s", err)
		}

		if _, err := js.Publish(c.Request().Context(), subject, payload); err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, "error forwarding log record: %s", err)
		}

		return c.NoContent(http.StatusAccepted)
	}, nil
}
