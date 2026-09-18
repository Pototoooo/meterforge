package subscriptiontestutils

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	customerservicehooks "github.com/Pototoooo/meterforge/meterforge/customer/service/hooks"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	meteradapter "github.com/Pototoooo/meterforge/meterforge/meter/mockadapter"
	addonrepo "github.com/Pototoooo/meterforge/meterforge/productcatalog/addon/adapter"
	addonservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/addon/service"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/featureresolver"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	planrepo "github.com/Pototoooo/meterforge/meterforge/productcatalog/plan/adapter"
	planservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/plan/service"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon"
	planaddonrepo "github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon/adapter"
	planaddonservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon/service"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	registrybuilder "github.com/Pototoooo/meterforge/meterforge/registry/builder"
	streamingtestutils "github.com/Pototoooo/meterforge/meterforge/streaming/testutils"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	subjecthooks "github.com/Pototoooo/meterforge/meterforge/subject/service/hooks"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionaddon "github.com/Pototoooo/meterforge/meterforge/subscription/addon"
	subscriptionaddonrepo "github.com/Pototoooo/meterforge/meterforge/subscription/addon/repo"
	subscriptionaddonservice "github.com/Pototoooo/meterforge/meterforge/subscription/addon/service"
	subscriptionentitlement "github.com/Pototoooo/meterforge/meterforge/subscription/entitlement"
	annotationhook "github.com/Pototoooo/meterforge/meterforge/subscription/hooks/annotations"
	"github.com/Pototoooo/meterforge/meterforge/subscription/service"
	subscriptionworkflow "github.com/Pototoooo/meterforge/meterforge/subscription/workflow"
	subscriptionworkflowservice "github.com/Pototoooo/meterforge/meterforge/subscription/workflow/service"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	taxcodeadapter "github.com/Pototoooo/meterforge/meterforge/taxcode/adapter"
	taxcodeservice "github.com/Pototoooo/meterforge/meterforge/taxcode/service"
	"github.com/Pototoooo/meterforge/meterforge/testutils"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/datetime"
	"github.com/Pototoooo/meterforge/pkg/ffx"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type SubscriptionDependencies struct {
	ItemRepo                 subscription.SubscriptionItemRepository
	CustomerAdapter          *testCustomerRepo
	CustomerService          customer.Service
	SubjectService           subject.Service
	FeatureConnector         *testFeatureConnector
	ExampleMeterID           string
	MeterService             meter.Service
	MockStreamingConnector   *streamingtestutils.MockStreamingConnector
	EntitlementAdapter       subscription.EntitlementAdapter
	PlanHelper               *planHelper
	PlanService              plan.Service
	DBDeps                   *DBDeps
	EntitlementRegistry      *registry.Entitlement
	SubscriptionService      subscription.Service
	WorkflowService          subscriptionworkflow.Service
	SubscriptionAddonService subscriptionaddon.Service
	AddonService             *testAddonService
	PlanAddonService         planaddon.Service
	TaxCodeService           taxcode.Service
}

func NewService(t *testing.T, dbDeps *DBDeps) SubscriptionDependencies {
	t.Helper()
	logger := testutils.NewLogger(t)
	subRepo := NewSubscriptionRepo(t, dbDeps)
	subPhaseRepo := NewSubscriptionPhaseRepo(t, dbDeps)
	subItemRepo := NewSubscriptionItemRepo(t, dbDeps)
	publisher := eventbus.NewMock(t)

	lockr, err := lockr.NewLocker(&lockr.LockerConfig{Logger: logger})
	require.NoError(t, err)

	meterID := ulid.Make().String()
	meterAdapter, err := meteradapter.New([]meter.Meter{{
		ManagedResource: models.ManagedResource{
			ID: meterID,
			NamespacedModel: models.NamespacedModel{
				Namespace: ExampleNamespace,
			},
			ManagedModel: models.ManagedModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "Meter 1",
		},
		Key:           ExampleFeatureMeterSlug,
		Aggregation:   meter.MeterAggregationSum,
		ValueProperty: lo.ToPtr("$.value"),
		EventType:     "test",
	}})
	require.NoError(t, err)
	require.NotNil(t, meterAdapter)
	require.NoError(t, meterAdapter.SetDBClient(dbDeps.DBClient))

	// After SetDBClient, the adapter may have resolved to an existing DB meter ID
	// (e.g., from a shared template DB). Read back the resolved ID.
	resolvedMeter, err := meterAdapter.GetMeterByIDOrSlug(context.Background(), meter.GetMeterInput{
		Namespace: ExampleNamespace,
		IDOrSlug:  ExampleFeatureMeterSlug,
	})
	require.NoError(t, err)
	meterID = resolvedMeter.ID

	mockStreaming := streamingtestutils.NewMockStreamingConnector(t)

	customerAdapter := NewCustomerAdapter(t, dbDeps)
	customerService := NewCustomerService(t, dbDeps)
	subjectService := NewSubjectService(t, dbDeps)

	entitlementRegistry := registrybuilder.GetEntitlementRegistry(registrybuilder.EntitlementOptions{
		DatabaseClient:     dbDeps.DBClient,
		StreamingConnector: mockStreaming,
		Logger:             logger,
		Tracer:             noop.NewTracerProvider().Tracer("test"),
		MeterService:       meterAdapter,
		CustomerService:    customerService,
		Publisher:          publisher,
		EntitlementsConfiguration: config.EntitlementsConfiguration{
			GracePeriod: datetime.ISODurationString("P1D"),
		},
		Locker: lockr,
	})

	entitlementAdapter := subscriptionentitlement.NewSubscriptionEntitlementAdapter(
		entitlementRegistry.Entitlement,
		subItemRepo,
		subPhaseRepo,
	)

	// Hooks

	// Subject hooks

	subjectCustomerHook, err := subjecthooks.NewCustomerSubjectHook(subjecthooks.CustomerSubjectHookConfig{
		Subject: subjectService,
		Logger:  logger,
		Tracer:  noop.NewTracerProvider().Tracer("test_env"),
	})
	require.NoError(t, err)
	customerService.RegisterHooks(subjectCustomerHook)

	// customer hooks
	customerSubjectHook, err := customerservicehooks.NewSubjectCustomerHook(customerservicehooks.SubjectCustomerHookConfig{
		Customer:         customerService,
		CustomerOverride: NoopCustomerOverrideService{},
		Logger:           logger,
		Tracer:           noop.NewTracerProvider().Tracer("test_env"),
	})
	require.NoError(t, err)
	subjectService.RegisterHooks(customerSubjectHook)

	entitlementValidatorHook, err := customerservicehooks.NewEntitlementValidatorHook(customerservicehooks.EntitlementValidatorHookConfig{
		EntitlementService: entitlementRegistry.Entitlement,
	})
	require.NoError(t, err)
	customerService.RegisterHooks(entitlementValidatorHook)

	planRepo, err := planrepo.New(planrepo.Config{
		Client: dbDeps.DBClient,
		Logger: logger,
	})
	require.NoError(t, err)

	taxCodeAdapter, err := taxcodeadapter.New(taxcodeadapter.Config{
		Client: dbDeps.DBClient,
		Logger: logger,
	})
	require.NoError(t, err)

	taxCodeService, err := taxcodeservice.New(taxcodeservice.Config{
		Adapter: taxCodeAdapter,
		Logger:  logger,
	})
	require.NoError(t, err)

	featureResolver, err := featureresolver.New(entitlementRegistry.Feature)
	require.NoErrorf(t, err, "failed to create feature resolver: %v", err)

	planService, err := planservice.New(planservice.Config{
		FeatureResolver: featureResolver,
		Logger:          logger,
		Adapter:         planRepo,
		Publisher:       publisher,
		TaxCode:         taxCodeService,
	})
	require.NoError(t, err)

	planHelper := NewPlanHelper(planService)

	ffService := ffx.NewTestContextService(ffx.AccessConfig{
		subscription.MultiSubscriptionEnabledFF: false,
	})

	svc, err := service.New(service.ServiceConfig{
		SubscriptionRepo:      subRepo,
		SubscriptionPhaseRepo: subPhaseRepo,
		SubscriptionItemRepo:  subItemRepo,
		CustomerService:       customerService,
		EntitlementAdapter:    entitlementAdapter,
		FeatureService:        entitlementRegistry.Feature,
		TransactionManager:    subItemRepo,
		Publisher:             publisher,
		Lockr:                 lockr,
		FeatureFlags:          ffService,
		TaxCode:               taxCodeService,
	})
	require.NoError(t, err)

	addonRepo, err := addonrepo.New(addonrepo.Config{
		Client: dbDeps.DBClient,
		Logger: logger,
	})
	require.NoError(t, err)

	addonService, err := addonservice.New(addonservice.Config{
		Adapter:         addonRepo,
		Logger:          logger,
		Publisher:       publisher,
		FeatureResolver: featureResolver,
		TaxCode:         taxCodeService,
	})
	require.NoError(t, err)

	planAddonRepo, err := planaddonrepo.New(planaddonrepo.Config{
		Client: dbDeps.DBClient,
		Logger: logger,
	})
	require.NoError(t, err)

	planAddonService, err := planaddonservice.New(planaddonservice.Config{
		Adapter:   planAddonRepo,
		Logger:    logger,
		Plan:      planService,
		Addon:     addonService,
		Publisher: publisher,
	})
	require.NoError(t, err)
	subAddRepo := subscriptionaddonrepo.NewSubscriptionAddonRepo(dbDeps.DBClient)
	subAddQtyRepo := subscriptionaddonrepo.NewSubscriptionAddonQuantityRepo(dbDeps.DBClient)

	subAddSvc, err := subscriptionaddonservice.NewService(subscriptionaddonservice.Config{
		TxManager:        subItemRepo,
		Logger:           logger,
		AddonService:     addonService,
		SubService:       svc,
		SubAddRepo:       subAddRepo,
		SubAddQtyRepo:    subAddQtyRepo,
		PlanAddonService: planAddonService,
		Publisher:        publisher,
	})
	require.NoError(t, err)

	annotationCleanupHook, err := annotationhook.NewAnnotationCleanupHook(svc, subRepo, logger)
	require.NoError(t, err)
	require.NoError(t, svc.RegisterHook(annotationCleanupHook))

	workflowSvc := subscriptionworkflowservice.NewWorkflowService(subscriptionworkflowservice.WorkflowServiceConfig{
		Service:            svc,
		CustomerService:    customerService,
		TransactionManager: subItemRepo,
		AddonService:       subAddSvc,
		Logger:             logger.With("subsystem", "subscription.workflow.service"),
		Lockr:              lockr,
		FeatureFlags:       ffService,
	})

	return SubscriptionDependencies{
		SubscriptionService:      svc,
		WorkflowService:          workflowSvc,
		CustomerAdapter:          customerAdapter,
		CustomerService:          customerService,
		SubjectService:           subjectService,
		FeatureConnector:         NewTestFeatureConnector(entitlementRegistry.Feature),
		ExampleMeterID:           meterID,
		EntitlementAdapter:       entitlementAdapter,
		DBDeps:                   dbDeps,
		PlanHelper:               planHelper,
		PlanService:              planService,
		ItemRepo:                 subItemRepo,
		EntitlementRegistry:      entitlementRegistry,
		SubscriptionAddonService: subAddSvc,
		AddonService:             NewTestAddonService(addonService),
		PlanAddonService:         planAddonService,
		MeterService:             meterAdapter,
		MockStreamingConnector:   mockStreaming,
		TaxCodeService:           taxCodeService,
	}
}
