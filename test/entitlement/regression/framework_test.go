package framework_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/Pototoooo/meterforge/meterforge/credit"
	grantrepo "github.com/Pototoooo/meterforge/meterforge/credit/adapter"
	"github.com/Pototoooo/meterforge/meterforge/credit/balance"
	"github.com/Pototoooo/meterforge/meterforge/credit/grant"
	credithook "github.com/Pototoooo/meterforge/meterforge/credit/hook"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	customeradapter "github.com/Pototoooo/meterforge/meterforge/customer/adapter"
	customerservice "github.com/Pototoooo/meterforge/meterforge/customer/service"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	enttx "github.com/Pototoooo/meterforge/meterforge/ent/tx"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	entitlementrepo "github.com/Pototoooo/meterforge/meterforge/entitlement/adapter"
	booleanentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/boolean"
	entitlementsubscriptionhook "github.com/Pototoooo/meterforge/meterforge/entitlement/hooks/subscription"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	entitlementservice "github.com/Pototoooo/meterforge/meterforge/entitlement/service"
	staticentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/static"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	meteradapter "github.com/Pototoooo/meterforge/meterforge/meter/mockadapter"
	productcatalogrepo "github.com/Pototoooo/meterforge/meterforge/productcatalog/adapter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	streamingtestutils "github.com/Pototoooo/meterforge/meterforge/streaming/testutils"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	subjectadapter "github.com/Pototoooo/meterforge/meterforge/subject/adapter"
	subjectservice "github.com/Pototoooo/meterforge/meterforge/subject/service"
	"github.com/Pototoooo/meterforge/meterforge/testutils"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/datetime"
	"github.com/Pototoooo/meterforge/pkg/framework/entutils/entdriver"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
	"github.com/Pototoooo/meterforge/pkg/framework/pgdriver"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type Dependencies struct {
	DBClient  *db.Client
	PGDriver  *pgdriver.Driver
	EntDriver *entdriver.EntPostgresDriver

	GrantRepo              grant.Repo
	BalanceSnapshotService balance.SnapshotService
	GrantConnector         credit.GrantConnector

	EntitlementRepo entitlement.EntitlementRepo

	EntitlementConnector        entitlement.Service
	StaticEntitlementConnector  staticentitlement.Connector
	BooleanEntitlementConnector booleanentitlement.Connector
	MeteredEntitlementConnector meteredentitlement.Connector

	Streaming *streamingtestutils.MockStreamingConnector

	FeatureRepo      feature.FeatureRepo
	FeatureConnector feature.FeatureConnector

	CustomerService customer.Service
	SubjectService  subject.Service

	Meter1ID string

	Log *slog.Logger
}

func (d *Dependencies) Close() {
	d.DBClient.Close()
	d.EntDriver.Close()
	d.PGDriver.Close()
}

func setupDependencies(t *testing.T) Dependencies {
	log := slog.Default()
	driver := testutils.InitPostgresDB(t, testutils.PostgresDBStateEntMigrated)
	// init db
	dbClient := db.NewClient(db.Driver(driver.EntDriver.Driver()))
	tracer := noop.NewTracerProvider().Tracer("test")

	// Init product catalog
	featureRepo := productcatalogrepo.NewPostgresFeatureRepo(dbClient, log)

	meter1ID := ulid.Make().String()
	meters := []meter.Meter{
		{
			ManagedResource: models.ManagedResource{
				ID: meter1ID,
				NamespacedModel: models.NamespacedModel{
					Namespace: "namespace-1",
				},
				ManagedModel: models.ManagedModel{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Name: "Meter 1",
			},
			Key:         "meter-1",
			Aggregation: meter.MeterAggregationCount,
			EventType:   "test",
		},
	}

	streaming := streamingtestutils.NewMockStreamingConnector(t)

	meterAdapter, err := meteradapter.New(meters)
	if err != nil {
		t.Fatalf("failed to create meter adapter: %v", err)
	}
	require.NoError(t, meterAdapter.SetDBClient(dbClient))

	featureConnector := feature.NewFeatureConnector(featureRepo, meterAdapter, eventbus.NewMock(t)) // TODO: meter repo is needed

	// Init grants/credit
	grantRepo := grantrepo.NewPostgresGrantRepo(dbClient)
	balanceSnapshotRepo := grantrepo.NewPostgresBalanceSnapshotRepo(dbClient)

	// Init entitlements
	entitlementRepo := entitlementrepo.NewPostgresEntitlementRepo(dbClient)
	usageResetRepo := entitlementrepo.NewPostgresUsageResetRepo(dbClient)

	mockPublisher := eventbus.NewMock(t)

	subjectRepo, err := subjectadapter.New(dbClient)
	require.NoError(t, err)

	subjectService, err := subjectservice.New(subjectRepo)
	require.NoError(t, err)

	customerAdapter, err := customeradapter.New(customeradapter.Config{
		Client: dbClient,
		Logger: log,
	})
	require.NoError(t, err)

	customerService, err := customerservice.New(customerservice.Config{
		Adapter:   customerAdapter,
		Publisher: mockPublisher,
	})
	require.NoError(t, err)

	owner := meteredentitlement.NewEntitlementGrantOwnerAdapter(
		featureRepo,
		entitlementRepo,
		usageResetRepo,
		meterAdapter,
		customerService,
		log,
		tracer,
	)

	balanceSnapshotService := balance.NewSnapshotService(balance.SnapshotServiceConfig{
		OwnerConnector:     owner,
		StreamingConnector: streaming,
		Repo:               balanceSnapshotRepo,
	})

	transactionManager := enttx.NewCreator(dbClient)

	creditConnector := credit.NewCreditConnector(
		credit.CreditConnectorConfig{
			GrantRepo:              grantRepo,
			BalanceSnapshotService: balanceSnapshotService,
			OwnerConnector:         owner,
			StreamingConnector:     streaming,
			Logger:                 log,
			Tracer:                 tracer,
			Granularity:            time.Minute,
			Publisher:              mockPublisher,
			TransactionManager:     transactionManager,
			SnapshotGracePeriod:    datetime.NewISODuration(0, 0, 0, 1, 0, 0, 0),
		},
	)

	meteredEntitlementConnector := meteredentitlement.NewMeteredEntitlementConnector(
		streaming,
		owner,
		creditConnector,
		creditConnector,
		grantRepo,
		entitlementRepo,
		mockPublisher,
		log,
		tracer,
	)

	meteredEntitlementConnector.RegisterHooks(
		meteredentitlement.ConvertHook(entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{})),
	)

	staticEntitlementConnector := staticentitlement.NewStaticEntitlementConnector()
	booleanEntitlementConnector := booleanentitlement.NewBooleanEntitlementConnector()

	locker, err := lockr.NewLocker(&lockr.LockerConfig{
		Logger: log,
	})
	require.NoError(t, err)

	entitlementConnector := entitlementservice.NewEntitlementService(
		entitlementservice.ServiceConfig{
			EntitlementRepo:             entitlementRepo,
			FeatureConnector:            featureConnector,
			CustomerService:             customerService,
			MeterService:                meterAdapter,
			MeteredEntitlementConnector: meteredEntitlementConnector,
			StaticEntitlementConnector:  staticEntitlementConnector,
			BooleanEntitlementConnector: booleanEntitlementConnector,
			Publisher:                   mockPublisher,
			Locker:                      locker,
		},
	)

	entitlementConnector.RegisterHooks(
		entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{}),
		credithook.NewEntitlementHook(grantRepo),
	)

	return Dependencies{
		DBClient:  dbClient,
		PGDriver:  driver.PGDriver,
		EntDriver: driver.EntDriver,

		GrantRepo:      grantRepo,
		GrantConnector: creditConnector,

		EntitlementRepo: entitlementRepo,

		EntitlementConnector:        entitlementConnector,
		StaticEntitlementConnector:  staticEntitlementConnector,
		BooleanEntitlementConnector: booleanEntitlementConnector,
		MeteredEntitlementConnector: meteredEntitlementConnector,

		BalanceSnapshotService: balanceSnapshotService,

		Streaming: streaming,

		FeatureRepo:      featureRepo,
		FeatureConnector: featureConnector,

		CustomerService: customerService,
		SubjectService:  subjectService,

		Meter1ID: meter1ID,

		Log: log,
	}
}

func createCustomerAndSubject(t *testing.T, subjectService subject.Service, customerService customer.Service, ns, key, name string) *customer.Customer {
	t.Helper()
	_, err := subjectService.Create(context.Background(), subject.CreateInput{
		Namespace: ns,
		Key:       key,
	})
	require.NoError(t, err)

	cust, err := customerService.CreateCustomer(context.Background(), customer.CreateCustomerInput{
		Namespace: ns,
		CustomerMutate: customer.CustomerMutate{
			Key: lo.ToPtr(key),
			UsageAttribution: &customer.CustomerUsageAttribution{
				SubjectKeys: []string{key},
			},
			Name: name,
		},
	})
	require.NoError(t, err)

	return cust
}
