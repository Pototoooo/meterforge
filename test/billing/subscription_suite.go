package billing

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/Pototoooo/meterforge/meterforge/credit"
	grantrepo "github.com/Pototoooo/meterforge/meterforge/credit/adapter"
	"github.com/Pototoooo/meterforge/meterforge/credit/balance"
	credithook "github.com/Pototoooo/meterforge/meterforge/credit/hook"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	enttx "github.com/Pototoooo/meterforge/meterforge/ent/tx"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	entitlementrepo "github.com/Pototoooo/meterforge/meterforge/entitlement/adapter"
	booleanentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/boolean"
	entitlementsubscriptionhook "github.com/Pototoooo/meterforge/meterforge/entitlement/hooks/subscription"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	entitlementservice "github.com/Pototoooo/meterforge/meterforge/entitlement/service"
	staticentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/static"
	meteradapter "github.com/Pototoooo/meterforge/meterforge/meter/mockadapter"
	addonrepo "github.com/Pototoooo/meterforge/meterforge/productcatalog/addon/adapter"
	addonservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/addon/service"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/featureresolver"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	planadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/plan/adapter"
	planservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/plan/service"
	planaddonrepo "github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon/adapter"
	planaddonservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon/service"
	subscriptiontestutils "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription/testutils"
	streamingtestutils "github.com/Pototoooo/meterforge/meterforge/streaming/testutils"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionaddon "github.com/Pototoooo/meterforge/meterforge/subscription/addon"
	subscriptionaddonrepo "github.com/Pototoooo/meterforge/meterforge/subscription/addon/repo"
	subscriptionaddonservice "github.com/Pototoooo/meterforge/meterforge/subscription/addon/service"
	subscriptionentitlementadatapter "github.com/Pototoooo/meterforge/meterforge/subscription/entitlement"
	subscriptionrepo "github.com/Pototoooo/meterforge/meterforge/subscription/repo"
	subscriptionservice "github.com/Pototoooo/meterforge/meterforge/subscription/service"
	subscriptionworkflow "github.com/Pototoooo/meterforge/meterforge/subscription/workflow"
	subscriptionworkflowservice "github.com/Pototoooo/meterforge/meterforge/subscription/workflow/service"
	taxcodeadapter "github.com/Pototoooo/meterforge/meterforge/taxcode/adapter"
	taxcodeservice "github.com/Pototoooo/meterforge/meterforge/taxcode/service"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/datetime"
	"github.com/Pototoooo/meterforge/pkg/ffx"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

type SubscriptionMixin struct {
	PlanService                 plan.Service
	SubscriptionService         subscription.Service
	SubscriptionAddonService    subscriptionaddon.Service
	SubscriptionPlanAdapter     subscriptiontestutils.PlanSubscriptionAdapter
	SubscriptionWorkflowService subscriptionworkflow.Service
	EntitlementConnector        entitlement.Service
}

type SubscriptionMixInDependencies struct {
	DBClient *db.Client

	FeatureRepo     feature.FeatureRepo
	FeatureService  feature.FeatureConnector
	CustomerService customer.Service

	MeterAdapter           *meteradapter.TestAdapter
	MockStreamingConnector *streamingtestutils.MockStreamingConnector
}

func (d SubscriptionMixInDependencies) Validate() error {
	if d.DBClient == nil {
		return errors.New("DBClient is required")
	}

	if d.FeatureRepo == nil {
		return errors.New("FeatureRepo is required")
	}

	if d.FeatureService == nil {
		return errors.New("FeatureService is required")
	}

	if d.CustomerService == nil {
		return errors.New("CustomerService is required")
	}

	if d.MeterAdapter == nil {
		return errors.New("MeterAdapter is required")
	}

	if d.MockStreamingConnector == nil {
		return errors.New("MockStreamingConnector is required")
	}

	return nil
}

func (s *SubscriptionMixin) SetupSuite(t *testing.T, deps SubscriptionMixInDependencies) {
	require.NoError(t, deps.Validate())

	publisher := eventbus.NewMock(t)

	lockr, err := lockr.NewLocker(&lockr.LockerConfig{Logger: slog.Default()})
	require.NoError(t, err)

	planAdapter, err := planadapter.New(planadapter.Config{
		Client: deps.DBClient,
		Logger: slog.Default(),
	})
	require.NoError(t, err)

	taxCodeAdapter, err := taxcodeadapter.New(taxcodeadapter.Config{
		Client: deps.DBClient,
		Logger: slog.Default(),
	})
	require.NoError(t, err)

	taxCodeService, err := taxcodeservice.New(taxcodeservice.Config{
		Adapter: taxCodeAdapter,
		Logger:  slog.Default(),
	})
	require.NoError(t, err)

	featureResolver, err := featureresolver.New(deps.FeatureService)
	require.NoErrorf(t, err, "failed to create feature resolver: %v", err)

	planService, err := planservice.New(planservice.Config{
		FeatureResolver: featureResolver,
		Adapter:         planAdapter,
		TaxCode:         taxCodeService,
		Logger:          slog.Default(),
		Publisher:       publisher,
	})
	require.NoError(t, err)

	s.PlanService = planService

	subsRepo := subscriptionrepo.NewSubscriptionRepo(deps.DBClient)
	subsItemRepo := subscriptionrepo.NewSubscriptionItemRepo(deps.DBClient)

	s.EntitlementConnector = s.SetupEntitlements(t, deps)

	ffService := ffx.NewStaticService(ffx.AccessConfig{
		subscription.MultiSubscriptionEnabledFF: true,
	})
	s.SubscriptionService, err = subscriptionservice.New(subscriptionservice.ServiceConfig{
		SubscriptionRepo:      subsRepo,
		SubscriptionPhaseRepo: subscriptionrepo.NewSubscriptionPhaseRepo(deps.DBClient),
		SubscriptionItemRepo:  subsItemRepo,
		// connectors
		CustomerService: deps.CustomerService,
		FeatureService:  deps.FeatureService,
		// adapters
		EntitlementAdapter: subscriptionentitlementadatapter.NewSubscriptionEntitlementAdapter(
			s.EntitlementConnector,
			subsItemRepo,
			subsRepo,
		),
		// framework
		TransactionManager: subsRepo,
		Lockr:              lockr,
		FeatureFlags:       ffService,
		// events
		Publisher: publisher,
	})
	require.NoError(t, err)

	addonRepo, err := addonrepo.New(addonrepo.Config{
		Client: deps.DBClient,
		Logger: slog.Default(),
	})
	require.NoError(t, err)

	addonService, err := addonservice.New(addonservice.Config{
		Adapter:         addonRepo,
		Logger:          slog.Default(),
		Publisher:       publisher,
		FeatureResolver: featureResolver,
		TaxCode:         taxCodeService,
	})
	require.NoError(t, err)

	planAddonRepo, err := planaddonrepo.New(planaddonrepo.Config{
		Client: deps.DBClient,
		Logger: slog.Default(),
	})
	require.NoError(t, err)

	planAddonService, err := planaddonservice.New(planaddonservice.Config{
		Adapter:   planAddonRepo,
		Logger:    slog.Default(),
		Plan:      planService,
		Addon:     addonService,
		Publisher: publisher,
	})
	require.NoError(t, err)

	subAddRepo := subscriptionaddonrepo.NewSubscriptionAddonRepo(deps.DBClient)
	subAddQtyRepo := subscriptionaddonrepo.NewSubscriptionAddonQuantityRepo(deps.DBClient)

	s.SubscriptionAddonService, err = subscriptionaddonservice.NewService(subscriptionaddonservice.Config{
		TxManager:        subsItemRepo,
		Logger:           slog.Default(),
		AddonService:     addonService,
		SubService:       s.SubscriptionService,
		SubAddRepo:       subAddRepo,
		SubAddQtyRepo:    subAddQtyRepo,
		PlanAddonService: planAddonService,
		Publisher:        publisher,
	})
	require.NoError(t, err)

	s.SubscriptionPlanAdapter = subscriptiontestutils.NewPlanSubscriptionAdapter(subscriptiontestutils.PlanSubscriptionAdapterConfig{
		PlanService: planService,
		Logger:      slog.Default(),
	})

	s.SubscriptionWorkflowService = subscriptionworkflowservice.NewWorkflowService(subscriptionworkflowservice.WorkflowServiceConfig{
		Service:            s.SubscriptionService,
		AddonService:       s.SubscriptionAddonService,
		CustomerService:    deps.CustomerService,
		TransactionManager: subsRepo,
		Logger:             slog.Default(),
		Lockr:              lockr,
		FeatureFlags:       ffService,
	})
}

func (s *SubscriptionMixin) SetupEntitlements(t *testing.T, deps SubscriptionMixInDependencies) entitlement.Service {
	tracer := noop.NewTracerProvider().Tracer("test")

	// Init grants/credit
	grantRepo := grantrepo.NewPostgresGrantRepo(deps.DBClient)
	balanceSnapshotRepo := grantrepo.NewPostgresBalanceSnapshotRepo(deps.DBClient)

	// Init entitlements
	entitlementRepo := entitlementrepo.NewPostgresEntitlementRepo(deps.DBClient)
	usageResetRepo := entitlementrepo.NewPostgresUsageResetRepo(deps.DBClient)

	mockPublisher := eventbus.NewMock(t)

	owner := meteredentitlement.NewEntitlementGrantOwnerAdapter(
		deps.FeatureRepo,
		entitlementRepo,
		usageResetRepo,
		deps.MeterAdapter,
		deps.CustomerService,
		slog.Default(),
		tracer,
	)

	balanceSnapshotService := balance.NewSnapshotService(balance.SnapshotServiceConfig{
		OwnerConnector:     owner,
		StreamingConnector: deps.MockStreamingConnector,
		Repo:               balanceSnapshotRepo,
	})

	transactionManager := enttx.NewCreator(deps.DBClient)

	creditConnector := credit.NewCreditConnector(
		credit.CreditConnectorConfig{
			GrantRepo:              grantRepo,
			BalanceSnapshotService: balanceSnapshotService,
			OwnerConnector:         owner,
			StreamingConnector:     deps.MockStreamingConnector,
			Logger:                 slog.Default(),
			Tracer:                 tracer,
			Granularity:            time.Minute,
			Publisher:              mockPublisher,
			TransactionManager:     transactionManager,
			SnapshotGracePeriod:    datetime.MustParseDuration(t, "P1W"),
		},
	)

	meteredEntitlementConnector := meteredentitlement.NewMeteredEntitlementConnector(
		deps.MockStreamingConnector,
		owner,
		creditConnector,
		creditConnector,
		grantRepo,
		entitlementRepo,
		mockPublisher,
		slog.Default(),
		tracer,
	)

	meteredEntitlementConnector.RegisterHooks(
		meteredentitlement.ConvertHook(entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{})),
	)

	staticEntitlementConnector := staticentitlement.NewStaticEntitlementConnector()
	booleanEntitlementConnector := booleanentitlement.NewBooleanEntitlementConnector()

	locker, err := lockr.NewLocker(&lockr.LockerConfig{
		Logger: slog.Default(),
	})
	require.NoError(t, err)

	service := entitlementservice.NewEntitlementService(
		entitlementservice.ServiceConfig{
			EntitlementRepo:             entitlementRepo,
			FeatureConnector:            deps.FeatureService,
			CustomerService:             deps.CustomerService,
			MeterService:                deps.MeterAdapter,
			MeteredEntitlementConnector: meteredEntitlementConnector,
			StaticEntitlementConnector:  staticEntitlementConnector,
			BooleanEntitlementConnector: booleanEntitlementConnector,
			Publisher:                   mockPublisher,
			Locker:                      locker,
		},
	)

	service.RegisterHooks(
		entitlementsubscriptionhook.NewEntitlementSubscriptionHook(entitlementsubscriptionhook.EntitlementSubscriptionHookConfig{}),
		credithook.NewEntitlementHook(grantRepo),
	)

	return service
}
