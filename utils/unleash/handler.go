package unleash

import (
	"os"
	"time"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/Unleash/unleash-client-go/v3/context"
)

var unleashHandler Handler

type Handler interface {
	IsEnabled(feature string) bool
	IsEnabledForID(feature, identifier string) bool
	Close() error
	WaitForReady()
}

type handler struct {
	client *unleash.Client
}

func GetHandler() (Handler, error) {
	if unleashHandler != nil {
		return unleashHandler, nil
	}

	client, err := unleash.NewClient(
		unleash.WithListener(&UnleashListener{}),
		unleash.WithAppName("unleash-server"),
		unleash.WithUrl(os.Getenv("UNLEASH_URL")),
		unleash.WithRefreshInterval(time.Second),
	)

	if err != nil {
		return nil, err
	}

	unleashHandler = &handler{client: client}
	unleashHandler.WaitForReady()

	return unleashHandler, nil
}

func (h *handler) IsEnabled(feature string) bool {
	return h.client.IsEnabled(feature)
}

func (h *handler) IsEnabledForID(feature, identifier string) bool {
	ctx := context.Context{UserId: identifier}

	return h.client.IsEnabled(feature, unleash.WithContext(ctx))
}

func (h *handler) WaitForReady() {
	h.client.WaitForReady()
}

func (h *handler) Close() error {
	return h.client.Close()
}
