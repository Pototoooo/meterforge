//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/common"
	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/billing/creditgrant"
	"github.com/Pototoooo/meterforge/meterforge/cost"
	"github.com/Pototoooo/meterforge/meterforge/currencies"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/governance"
	"github.com/Pototoooo/meterforge/meterforge/ingest"
	"github.com/Pototoooo/meterforge/meterforge/ingest/kafkaingest"
	"github.com/Pototoooo/meterforge/meterforge/ledger"
	"github.com/Pototoooo/meterforge/meterforge/ledger/customerbalance"
	"github.com/Pototoooo/meterforge/meterforge/llmcost"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/meterevent"
	"github.com/Pototoooo/meterforge/meterforge/namespace"
	"github.com/Pototoooo/meterforge/meterforge/notification"
	"github.com/Pototoooo/meterforge/meterforge/portal"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon"
	"github.com/Pototoooo/meterforge/meterforge/progressmanager"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	"github.com/Pototoooo/meterforge/meterforge/secret"
	"github.com/Pototoooo/meterforge/meterforge/server"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	subjecthooks "github.com/Pototoooo/meterforge/meterforge/subject/service/hooks"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/featuregate"
	"github.com/Pototoooo/meterforge/pkg/ffx"
	kafkametrics "github.com/Pototoooo/meterforge/pkg/kafka/metrics"
)

type Application struct {
	common.GlobalInitializer
	common.Migrator

	Addon                            addon.Service
	AppRegistry                      common.AppRegistry
	Customer                         customer.Service
	CustomerLedgerHook               common.CustomerLedgerHook
	CustomerSubjectHook              common.CustomerSubjectHook
	CustomerEntitlementValidatorHook common.CustomerEntitlementValidatorHook
	BillingRegistry                  common.BillingRegistry
	CurrencyService                  currencies.Service
	CostService                      cost.Service
	CreditGrantService               creditgrant.Service
	Ledger                           ledger.Ledger
	AccountResolver                  ledger.AccountResolver
	CustomerBalanceFacade            *customerbalance.Facade
	EntClient                        *db.Client
	EventPublisher                   eventbus.Publisher
	EntitlementRegistry              *registry.Entitlement
	FeatureConnector                 feature.FeatureConnector
	FeatureFlags                     ffx.Service
	GovernanceService                governance.Service
	IngestCollector                  ingest.Collector
	IngestService                    ingest.Service
	KafkaProducer                    *kafka.Producer
	KafkaMetrics                     *kafkametrics.Metrics
	KafkaIngestNamespaceHandler      *kafkaingest.NamespaceHandler
	LedgerNamespaceHandler           namespace.Handler
	LLMCostService                   llmcost.Service
	Logger                           *slog.Logger
	MetricMeter                      metric.Meter
	MeterConfigInitializer           common.MeterConfigInitializer
	MeterManageService               meter.ManageService
	MeterEventService                meterevent.Service
	NamespaceManager                 *namespace.Manager
	Notification                     notification.Service
	NotificationEventHandler         notification.EventHandler
	Plan                             plan.Service
	PlanAddon                        planaddon.Service
	Portal                           portal.Service
	ProgressManager                  progressmanager.Service
	RouterHooks                      *server.RouterHooks
	PostAuthMiddlewares              server.PostAuthMiddlewares
	Secret                           secret.Service
	SubjectService                   subject.Service
	SubjectCustomerHook              subjecthooks.CustomerSubjectHook
	Subscription                     common.SubscriptionServiceWithWorkflow
	StreamingConnector               streaming.Connector
	TaxCodeNamespaceHandler          *taxcode.NamespaceHandler
	TaxCodeService                   taxcode.Service
	TelemetryServer                  common.TelemetryServer
	TerminationChecker               *common.TerminationChecker
	RuntimeMetricsCollector          common.RuntimeMetricsCollector
	Tracer                           trace.Tracer
	FeatureGate                      *featuregate.FeatureGateChecker
	ClientIPMiddleware               common.ClientIPMiddleware
}

func initializeApplication(ctx context.Context, conf config.Configuration) (Application, func(), error) {
	wire.Build(
		metadata,
		common.App,
		common.Billing,
		common.ClickHouse,
		common.Config,
		common.CreditGrant,
		common.Currency,
		common.Customer,
		common.CustomerBalance,
		common.NewCustomerLedgerServiceHook,
		common.NewCustomerSubjectServiceHook,
		common.NewCustomerEntitlementValidatorServiceHook,
		common.Database,
		common.Entitlement,
		common.Framework,
		common.FFX,
		common.Governance,
		common.Kafka,
		common.KafkaIngest,
		common.LLMCost,
		common.LedgerStack,
		common.KafkaNamespaceResolver,
		common.MeterManageWithConfigMeters,
		common.MeterEvent,
		common.Namespace,
		common.StaticNamespace,
		common.NewDefaultTextMapPropagator,
		common.NewKafkaIngestCollector,
		common.NewIngestCollector,
		common.NewIngestService,
		common.NewServerPublisher,
		common.Notification,
		common.Streaming,
		common.Portal,
		common.ProductCatalog,
		common.ProgressManager,
		common.Server,
		common.TaxCode,
		common.TaxCodeNamespaceHandler,
		common.Subscription,
		common.Lockr,
		common.Secret,
		common.ServerProvisionTopics,
		common.Subject,
		common.NewSvixAPIClient,
		common.NewSubjectCustomerHook,
		common.Telemetry,
		common.TelemetryLoggerNoAdditionalMiddlewares,
		common.NewTerminationChecker,
		common.WatermillNoPublisher,
		wire.Struct(new(Application), "*"),
		common.FeatureGateChecker,
	)

	return Application{}, nil, nil
}

func metadata(conf config.Configuration) common.Metadata {
	return common.NewMetadata(conf, version, "backend")
}
