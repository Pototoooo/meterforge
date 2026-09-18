package registrybuilder

import (
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/credit"
	creditpgadapter "github.com/Pototoooo/meterforge/meterforge/credit/adapter"
	"github.com/Pototoooo/meterforge/meterforge/credit/balance"
	credithook "github.com/Pototoooo/meterforge/meterforge/credit/hook"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	enttx "github.com/Pototoooo/meterforge/meterforge/ent/tx"
	entitlementpgadapter "github.com/Pototoooo/meterforge/meterforge/entitlement/adapter"
	booleanentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/boolean"
	entitlementsubscriptionhook "github.com/Pototoooo/meterforge/meterforge/entitlement/hooks/subscription"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	entitlementservice "github.com/Pototoooo/meterforge/meterforge/entitlement/service"
	staticentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/static"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	productcatalogpgadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/adapter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

type EntitlementOptions struct {
	DatabaseClient            *db.Client
	EntitlementsConfiguration config.EntitlementsConfiguration
	StreamingConnector        streaming.Connector
	Logger                    *slog.Logger
	MeterService              meter.Service
	CustomerService           customer.Service
	Publisher                 eventbus.Publisher
	Tracer                    trace.Tracer
	Locker                    *lockr.Locker
}

func GetEntitlementRegistry(opts EntitlementOptions) *registry.Entitlement {
	// Initialize database adapters
	featureDBAdapter := productcatalogpgadapter.NewPostgresFeatureRepo(opts.DatabaseClient, opts.Logger)
	entitlementDBAdapter := entitlementpgadapter.NewPostgresEntitlementRepo(opts.DatabaseClient)
	usageResetDBAdapter := entitlementpgadapter.NewPostgresUsageResetRepo(opts.DatabaseClient)
	grantDBAdapter := creditpgadapter.NewPostgresGrantRepo(opts.DatabaseClient)
	balanceSnashotDBAdapter := creditpgadapter.NewPostgresBalanceSnapshotRepo(opts.DatabaseClient)

	// Initialize connectors
	featureConnector := feature.NewFeatureConnector(featureDBAdapter, opts.MeterService, opts.Publisher)
	entitlementOwnerConnector := meteredentitlement.NewEntitlementGrantOwnerAdapter(
		featureDBAdapter,
		entitlementDBAdapter,
		usageResetDBAdapter,
		opts.MeterService,
		opts.CustomerService,
		opts.Logger,
		opts.Tracer,
	)
	transactionManager := enttx.NewCreator(opts.DatabaseClient)

	balanceSnapshotService := balance.NewSnapshotService(balance.SnapshotServiceConfig{
		OwnerConnector:     entitlementOwnerConnector,
		StreamingConnector: opts.StreamingConnector,
		Repo:               balanceSnashotDBAdapter,
	})

	creditConnector := credit.NewCreditConnector(
		credit.CreditConnectorConfig{
			GrantRepo:              grantDBAdapter,
			BalanceSnapshotService: balanceSnapshotService,
			OwnerConnector:         entitlementOwnerConnector,
			StreamingConnector:     opts.StreamingConnector,
			Logger:                 opts.Logger,
			Tracer:                 opts.Tracer,
			Granularity:            time.Minute,
			SnapshotGracePeriod:    opts.EntitlementsConfiguration.GetGracePeriod(),
			TransactionManager:     transactionManager,
			Publisher:              opts.Publisher,
		},
	)
	creditBalanceConnector := creditConnector
	grantConnector := creditConnector
	meteredEntitlementConnector := meteredentitlement.NewMeteredEntitlementConnector(
		opts.StreamingConnector,
		entitlementOwnerConnector,
		creditBalanceConnector,
		grantConnector,
		grantDBAdapter,
		entitlementDBAdapter,
		opts.Publisher,
		opts.Logger,
		opts.Tracer,
	)

	meteredEntitlementConnector.RegisterHooks(
		meteredentitlement.ConvertHook(entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{})),
	)

	entitlementConnector := entitlementservice.NewEntitlementService(
		entitlementservice.ServiceConfig{
			EntitlementRepo:             entitlementDBAdapter,
			FeatureConnector:            featureConnector,
			CustomerService:             opts.CustomerService,
			MeterService:                opts.MeterService,
			MeteredEntitlementConnector: meteredEntitlementConnector,
			StaticEntitlementConnector:  staticentitlement.NewStaticEntitlementConnector(),
			BooleanEntitlementConnector: booleanentitlement.NewBooleanEntitlementConnector(),
			Publisher:                   opts.Publisher,
			Locker:                      opts.Locker,
		},
	)

	entitlementConnector.RegisterHooks(
		entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{}),
		credithook.NewEntitlementHook(grantDBAdapter),
	)

	return &registry.Entitlement{
		Feature:            featureConnector,
		FeatureRepo:        featureDBAdapter,
		EntitlementOwner:   entitlementOwnerConnector,
		CreditBalance:      creditBalanceConnector,
		Grant:              grantConnector,
		GrantRepo:          grantDBAdapter,
		MeteredEntitlement: meteredEntitlementConnector,
		Entitlement:        entitlementConnector,
		EntitlementRepo:    entitlementDBAdapter,
	}
}
