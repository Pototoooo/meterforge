package meteredentitlement_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/Pototoooo/meterforge/meterforge/credit"
	credit_postgres_adapter "github.com/Pototoooo/meterforge/meterforge/credit/adapter"
	"github.com/Pototoooo/meterforge/meterforge/credit/balance"
	"github.com/Pototoooo/meterforge/meterforge/credit/grant"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	customeradapter "github.com/Pototoooo/meterforge/meterforge/customer/adapter"
	customerservice "github.com/Pototoooo/meterforge/meterforge/customer/service"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	enttx "github.com/Pototoooo/meterforge/meterforge/ent/tx"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	entitlement_postgresadapter "github.com/Pototoooo/meterforge/meterforge/entitlement/adapter"
	entitlementsubscriptionhook "github.com/Pototoooo/meterforge/meterforge/entitlement/hooks/subscription"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	meteradapter "github.com/Pototoooo/meterforge/meterforge/meter/mockadapter"
	productcatalog_postgresadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/adapter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	streamingtestutils "github.com/Pototoooo/meterforge/meterforge/streaming/testutils"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	subjectadapter "github.com/Pototoooo/meterforge/meterforge/subject/adapter"
	subjectservice "github.com/Pototoooo/meterforge/meterforge/subject/service"
	"github.com/Pototoooo/meterforge/meterforge/testutils"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/datetime"
	"github.com/Pototoooo/meterforge/pkg/framework/entutils/entdriver"
	"github.com/Pototoooo/meterforge/pkg/framework/pgdriver"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type dependencies struct {
	dbClient               *db.Client
	pgDriver               *pgdriver.Driver
	entDriver              *entdriver.EntPostgresDriver
	featureRepo            feature.FeatureRepo
	entitlementRepo        entitlement.EntitlementRepo
	usageResetRepo         meteredentitlement.UsageResetRepo
	grantRepo              grant.Repo
	balanceSnapshotService balance.SnapshotService
	balanceConnector       credit.BalanceConnector
	ownerConnector         grant.OwnerConnector
	streamingConnector     *streamingtestutils.MockStreamingConnector
	creditConnector        credit.CreditConnector
	meterAdapter           *meteradapter.TestAdapter
	meterID                string
	subjectService         subject.Service
	customerService        customer.Service
}

// Teardown cleans up the dependencies
func (d *dependencies) Teardown() {
	d.dbClient.Close()
	d.entDriver.Close()
	d.pgDriver.Close()
}

var (
	namespace = "ns1"
	meterSlug = "meter1"
)

// When migrating in parallel with entgo it causes concurrent writes error
var m sync.Mutex

// builds connector with mock streaming and real PG
func setupConnector(t *testing.T) (meteredentitlement.Connector, *dependencies) {
	testLogger := testutils.NewLogger(t)
	tracer := noop.NewTracerProvider().Tracer("test")
	streamingConnector := streamingtestutils.NewMockStreamingConnector(t)
	testMeterID := ulid.Make().String()
	testMeters := []meter.Meter{{
		ManagedResource: models.ManagedResource{
			ID: testMeterID,
			NamespacedModel: models.NamespacedModel{
				Namespace: namespace,
			},
			ManagedModel: models.ManagedModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "Meter 1",
		},
		Key:         meterSlug,
		Aggregation: meter.MeterAggregationSum,
		// These will be ignored in tests
		EventType:     "test",
		ValueProperty: lo.ToPtr("$.value"),
	}}
	meterAdapter, err := meteradapter.New(testMeters)
	if err != nil {
		t.Fatalf("failed to create meter adapter: %v", err)
	}

	// create isolated pg db for tests
	testdb := testutils.InitPostgresDB(t, testutils.PostgresDBStateEntMigrated)
	dbClient := testdb.EntDriver.Client()
	pgDriver := testdb.PGDriver
	entDriver := testdb.EntDriver

	featureRepo := productcatalog_postgresadapter.NewPostgresFeatureRepo(dbClient, testLogger)
	entitlementRepo := entitlement_postgresadapter.NewPostgresEntitlementRepo(dbClient)
	usageResetRepo := entitlement_postgresadapter.NewPostgresUsageResetRepo(dbClient)
	grantRepo := credit_postgres_adapter.NewPostgresGrantRepo(dbClient)
	balanceSnapshotRepo := credit_postgres_adapter.NewPostgresBalanceSnapshotRepo(dbClient)

	m.Lock()
	defer m.Unlock()
	// migrate db
	require.NoError(t, meterAdapter.SetDBClient(dbClient))
	require.NoError(t, meterAdapter.ReplaceMeters(context.Background(), testMeters))

	mockPublisher := eventbus.NewMock(t)

	subjectRepo, err := subjectadapter.New(dbClient)
	require.NoError(t, err)

	subjectService, err := subjectservice.New(subjectRepo)
	require.NoError(t, err)

	customerAdapter, err := customeradapter.New(customeradapter.Config{
		Client: dbClient,
		Logger: testLogger,
	})
	require.NoError(t, err)

	customerService, err := customerservice.New(customerservice.Config{
		Adapter:   customerAdapter,
		Publisher: mockPublisher,
	})
	require.NoError(t, err)

	// build adapters
	ownerConnector := meteredentitlement.NewEntitlementGrantOwnerAdapter(
		featureRepo,
		entitlementRepo,
		usageResetRepo,
		meterAdapter,
		customerService,
		testLogger,
		tracer,
	)

	transactionManager := enttx.NewCreator(dbClient)

	snapshotService := balance.NewSnapshotService(balance.SnapshotServiceConfig{
		OwnerConnector:     ownerConnector,
		StreamingConnector: streamingConnector,
		Repo:               balanceSnapshotRepo,
	})

	creditConnector := credit.NewCreditConnector(
		credit.CreditConnectorConfig{
			GrantRepo:              grantRepo,
			BalanceSnapshotService: snapshotService,
			OwnerConnector:         ownerConnector,
			StreamingConnector:     streamingConnector,
			Logger:                 testLogger,
			Tracer:                 tracer,
			Granularity:            time.Minute,
			Publisher:              mockPublisher,
			SnapshotGracePeriod:    datetime.MustParseDuration(t, "P1W"),
			TransactionManager:     transactionManager,
		},
	)

	connector := meteredentitlement.NewMeteredEntitlementConnector(
		streamingConnector,
		ownerConnector,
		creditConnector,
		creditConnector,
		grantRepo,
		entitlementRepo,
		mockPublisher,
		testLogger,
		tracer,
	)

	connector.RegisterHooks(
		meteredentitlement.ConvertHook(entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{})),
	)

	return connector, &dependencies{
		dbClient:               dbClient,
		pgDriver:               pgDriver,
		entDriver:              entDriver,
		featureRepo:            featureRepo,
		entitlementRepo:        entitlementRepo,
		usageResetRepo:         usageResetRepo,
		grantRepo:              grantRepo,
		balanceSnapshotService: snapshotService,
		balanceConnector:       creditConnector,
		ownerConnector:         ownerConnector,
		streamingConnector:     streamingConnector,
		creditConnector:        creditConnector,
		meterAdapter:           meterAdapter,
		meterID:                testMeterID,
		subjectService:         subjectService,
		customerService:        customerService,
	}
}

func assertUsagePeriodInputsEquals(t *testing.T, expected, actual *entitlement.UsagePeriodInput) {
	t.Helper()
	assert.NotNil(t, expected, "expected is nil")
	assert.NotNil(t, actual, "actual is nil")
	assert.Equal(t, expected.GetValue().Interval, actual.GetValue().Interval, "periods do not match")
	assert.Equal(t, expected.GetValue().Anchor.UTC().Format(time.RFC3339), actual.GetValue().Anchor.UTC().Format(time.RFC3339), "anchors do not match")
}

func createCustomerAndSubject(t *testing.T, subjectService subject.Service, customerService customer.Service, ns, key, name string) *customer.Customer {
	t.Helper()

	_, err := subjectService.Create(t.Context(), subject.CreateInput{
		Namespace: ns,
		Key:       key,
	})
	require.NoError(t, err)

	cust, err := customerService.CreateCustomer(t.Context(), customer.CreateCustomerInput{
		Namespace: ns,
		CustomerMutate: customer.CustomerMutate{
			Key:  lo.ToPtr(key),
			Name: name,
			UsageAttribution: &customer.CustomerUsageAttribution{
				SubjectKeys: []string{key},
			},
		},
	})
	require.NoError(t, err)

	return cust
}
