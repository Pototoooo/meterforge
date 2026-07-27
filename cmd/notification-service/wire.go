//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/common"
	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/notification"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	watermillkafka "github.com/Pototoooo/meterforge/meterforge/watermill/driver/kafka"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
)

type Application struct {
	common.GlobalInitializer
	common.Migrator

	BrokerOptions           watermillkafka.BrokerOptions
	EventPublisher          eventbus.Publisher
	EntClient               *db.Client
	FeatureConnector        feature.FeatureConnector
	Logger                  *slog.Logger
	MessagePublisher        message.Publisher
	Meter                   metric.Meter
	Tracer                  trace.Tracer
	Metadata                common.Metadata
	MeterService            meter.Service
	Notification            notification.Service
	RuntimeMetricsCollector common.RuntimeMetricsCollector
	StreamingConnector      streaming.Connector
	TelemetryServer         common.TelemetryServer
}

func initializeApplication(ctx context.Context, conf config.Configuration) (Application, func(), error) {
	wire.Build(
		metadata,
		common.ClickHouse,
		common.Config,
		common.Database,
		common.Feature,
		common.Framework,
		common.Meter,
		common.Namespace,
		common.NewDefaultTextMapPropagator,
		common.NewKafkaTopicProvisioner,
		common.Notification,
		common.NotificationServiceProvisionTopics,
		common.ProgressManager,
		common.Streaming,
		common.NewSvixAPIClient,
		common.Telemetry,
		common.TelemetryLoggerNoAdditionalMiddlewares,
		common.Watermill,
		wire.Struct(new(Application), "*"),
	)
	return Application{}, nil, nil
}

func metadata(conf config.Configuration) common.Metadata {
	return common.NewMetadata(conf, version, "notification-worker")
}
