//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-slog/otelslog"
	"github.com/google/wire"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/common"
	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/ingest/kafkaingest/topicresolver"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/sink"
	"github.com/Pototoooo/meterforge/meterforge/sink/flushhandler"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	pkgkafka "github.com/Pototoooo/meterforge/pkg/kafka"
)

type Application struct {
	common.GlobalInitializer

	FlushHandler            flushhandler.FlushEventHandler
	Logger                  *slog.Logger
	Metadata                common.Metadata
	Meter                   metric.Meter
	Streaming               streaming.Connector
	TelemetryServer         common.TelemetryServer
	TopicProvisioner        pkgkafka.TopicProvisioner
	TopicResolver           *topicresolver.NamespacedTopicResolver
	Tracer                  trace.Tracer
	MeterService            meter.Service
	RuntimeMetricsCollector common.RuntimeMetricsCollector
	Sink                    *sink.Sink
}

func initializeApplication(ctx context.Context, conf config.Configuration) (Application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(config.Configuration), "Sink"),

		metadata,
		common.ClickHouse,
		common.Config,
		common.Database,
		common.Framework,
		common.KafkaNamespaceResolver,
		common.NewKafkaTopicProvisioner,
		common.Meter,
		common.Namespace,
		common.NewDefaultTextMapPropagator,
		common.ProgressManager,
		common.SinkWorkerProvisionTopics,
		common.Streaming,
		common.Sink,
		common.Telemetry,
		common.TelemetryLoggerNoAdditionalMiddlewares,
		common.WatermillNoPublisher,
		wire.Struct(new(Application), "*"),
	)
	return Application{}, nil, nil
}

func metadata(conf config.Configuration) common.Metadata {
	return common.NewMetadata(conf, version, "sink-worker")
}

// TODO: use the primary logger
func NewLogger(conf config.LogTelemetryConfig, res *resource.Resource) *slog.Logger {
	return slog.New(slogmulti.Pipe(
		otelslog.ResourceMiddleware(res),
		otelslog.NewHandler,
	).Handler(conf.NewHandler(os.Stdout)))
}
