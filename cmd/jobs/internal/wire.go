//go:build wireinject
// +build wireinject

package internal

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/metric"

	"github.com/Pototoooo/meterforge/app/common"
	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/app"
	appstripe "github.com/Pototoooo/meterforge/meterforge/app/stripe"
	chargesworkeradvance "github.com/Pototoooo/meterforge/meterforge/billing/charges/worker/advance"
	billingworkerautoadvance "github.com/Pototoooo/meterforge/meterforge/billing/worker/advance"
	billingworkercollect "github.com/Pototoooo/meterforge/meterforge/billing/worker/collect"
	billingworkersubscriptionreconciler "github.com/Pototoooo/meterforge/meterforge/billing/worker/subscriptionsync/reconciler"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	llmcostsync "github.com/Pototoooo/meterforge/meterforge/llmcost/sync"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/namespace"
	"github.com/Pototoooo/meterforge/meterforge/notification"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	"github.com/Pototoooo/meterforge/meterforge/secret"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	kafkametrics "github.com/Pototoooo/meterforge/pkg/kafka/metrics"
)

type Application struct {
	common.GlobalInitializer
	common.Migrator

	App                           app.Service
	AppStripe                     appstripe.Service
	AppSandboxProvisioner         common.AppSandboxProvisioner
	Customer                      customer.Service
	BillingRegistry               common.BillingRegistry
	BillingAutoAdvancer           *billingworkerautoadvance.AutoAdvancer
	ChargesAutoAdvancer           *chargesworkeradvance.AutoAdvancer
	BillingCollector              *billingworkercollect.InvoiceCollector
	BillingSubscriptionReconciler *billingworkersubscriptionreconciler.Reconciler
	EntClient                     *db.Client
	EventPublisher                eventbus.Publisher
	EntitlementRegistry           *registry.Entitlement
	FeatureConnector              feature.FeatureConnector
	KafkaProducer                 *kafka.Producer
	KafkaMetrics                  *kafkametrics.Metrics
	Logger                        *slog.Logger
	MeterService                  meter.Service
	NamespaceManager              *namespace.Manager
	Meter                         metric.Meter
	NotificationService           notification.Service
	Plan                          plan.Service
	Secret                        secret.Service
	Subscription                  common.SubscriptionServiceWithWorkflow
	Subject                       subject.Service
	StreamingConnector            streaming.Connector
	LLMCostSyncJob                *llmcostsync.SyncJob
}

func initializeApplication(ctx context.Context, conf config.Configuration) (Application, func(), error) {
	wire.Build(
		metadata,
		common.App,
		common.Billing,
		common.ClickHouse,
		common.Config,
		common.Customer,
		common.Currency,
		common.Subject,
		common.Database,
		common.Entitlement,
		common.Framework,
		common.Kafka,
		common.LLMCost,
		common.Meter,
		common.FFX,
		common.Namespace,
		common.NewBillingAutoAdvancer,
		common.NewChargesAutoAdvancer,
		common.NewBillingCollector,
		common.NewBillingSubscriptionSyncAdapter,
		common.NewBillingSubscriptionSyncService,
		common.NewBillingSubscriptionReconciler,
		common.NewDefaultTextMapPropagator,
		common.NewServerPublisher,
		common.Notification,
		common.Streaming,
		common.TaxCode,
		common.ProductCatalog,
		common.ProgressManager,
		common.Subscription,
		common.LedgerStack,
		common.Lockr,
		common.Secret,
		common.ServerProvisionTopics,
		common.NewSvixAPIClient,
		common.TelemetryWithoutServer,
		common.TelemetryLoggerNoAdditionalMiddlewares,
		common.WatermillNoPublisher,
		wire.Struct(new(Application), "*"),
		common.FeatureGateChecker,
	)

	return Application{}, nil, nil
}

func metadata(conf config.Configuration) common.Metadata {
	return common.NewMetadata(conf, version, "jobs")
}
