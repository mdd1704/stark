package unleash

import (
	"context"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/palantir/stacktrace"

	"stark/utils/log"
)

type UnleashListener struct {
}

// OnError prints out errors.
func (l UnleashListener) OnError(err error) {
	log.WithContext(context.Background()).Error(stacktrace.Propagate(err, "Unleash received error"))
}

// OnWarning prints out warning.
func (l UnleashListener) OnWarning(warning error) {
	log.WithContext(context.Background()).Warn(stacktrace.Propagate(warning, "Unleash received warning"))
}

// OnReady prints to the console when the repository is ready.
func (l UnleashListener) OnReady() {

}

// OnCount prints to the console when the feature is queried.
func (l UnleashListener) OnCount(name string, enabled bool) {

}

// OnSent prints to the console when the server has uploaded metrics.
func (l UnleashListener) OnSent(payload unleash.MetricsData) {

}

// OnRegistered prints to the console when the client has registered.
func (l UnleashListener) OnRegistered(payload unleash.ClientData) {
	log.WithContext(context.Background()).Infof("Registered: %+v\n", payload)
}
