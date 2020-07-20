package logger

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

const (
	providerField  = "provider"
	endpointField  = "endpoint"
	requestIDField = "request_id"
)

// ProviderLogger implements basic logger for any endpoints provider.
type ProviderLogger struct {
	name string
}

// NewProviderLogger returns new logger for provider.
func NewProviderLogger(name string) *ProviderLogger {
	return &ProviderLogger{name: name}
}

// Logger returns new logger entry for the specified endpoint and request.
func (pl *ProviderLogger) Logger(request *http.Request) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		providerField:  pl.name,
		endpointField:  request.RequestURI,
		requestIDField: middleware.GetReqID(request.Context()),
	})
}
