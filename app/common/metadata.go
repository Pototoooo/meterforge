package common

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/Pototoooo/meterforge/app/config"
)

// Metadata provides information about the service to components that need it (eg. telemetry).
type Metadata struct {
	ServiceName       string
	Version           string
	Environment       string
	OpenTelemetryName string

	AdditionalAttributes []attribute.KeyValue
}

func NewMetadata(conf config.Configuration, version string, serviceName string, additionalAttributes ...attribute.KeyValue) Metadata {
	return Metadata{
		ServiceName:          fmt.Sprintf("meterforge-%s", serviceName),
		Version:              version,
		Environment:          conf.Environment,
		OpenTelemetryName:    fmt.Sprintf("github.com/Pototoooo/meterforge/%s", serviceName),
		AdditionalAttributes: additionalAttributes,
	}
}
