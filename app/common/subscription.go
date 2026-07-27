package common

import (
	"log/slog"

	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon"
	plansubscription "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription"
	subscriptionchangeservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription/service"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionaddon "github.com/Pototoooo/meterforge/meterforge/subscription/addon"
	subscriptionaddonrepo "github.com/Pototoooo/meterforge/meterforge/subscription/addon/repo"
	subscriptionaddonservice "github.com/Pototoooo/meterforge/meterforge/subscription/addon/service"
	subscriptionentitlement "github.com/Pototoooo/meterforge/meterforge/subscription/entitlement"
	annotationhook "github.com/Pototoooo/meterforge/meterforge/subscription/hooks/annotations"
	subscriptionrepo "github.com/Pototoooo/meterforge/meterforge/subscription/repo"
	subscriptionservice "github.com/Pototoooo/meterforge/meterforge/subscription/service"
	subscriptioncustomer "github.com/Pototoooo/meterforge/meterforge/subscription/validators/customer"
	subscriptionworkflow "github.com/Pototoooo/meterforge/meterforge/subscription/workflow"
	subscriptionworkflowservice "github.com/Pototoooo/meterforge/meterforge/subscription/workflow/service"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/ffx"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

var Subscription = wire.NewSet(
	NewSubscriptionServices,
)

// TODO: break up to multiple initializers
type SubscriptionServiceWithWorkflow struct {
	Service                  subscription.Service
	WorkflowService          subscriptionworkflow.Service
	PlanSubscriptionService  plansubscription.PlanSubscriptionService
	SubscriptionAddonService subscriptionaddon.Service
}

func NewSubscriptionServices(
	logger *slog.Logger,
	db *entdb.Client,
	featureConnector feature.FeatureConnector,
	entitlementRegistry *registry.Entitlement,
	customerService customer.Service,
	planService plan.Service,
	planAddonService planaddon.Service,
	addonService addon.Service,
	eventPublisher eventbus.Publisher,
	lockr *lockr.Locker,
	featureFlags ffx.Service,
	taxCodeService taxcode.Service,
) (SubscriptionServiceWithWorkflow, error) {
	subscriptionRepo := subscriptionrepo.NewSubscriptionRepo(db)
	subscriptionPhaseRepo := subscriptionrepo.NewSubscriptionPhaseRepo(db)
	subscriptionItemRepo := subscriptionrepo.NewSubscriptionItemRepo(db)

	subscriptionEntitlementAdapter := subscriptionentitlement.NewSubscriptionEntitlementAdapter(
		entitlementRegistry.Entitlement,
		subscriptionItemRepo,
		subscriptionItemRepo,
	)

	subscriptionService, err := subscriptionservice.New(subscriptionservice.ServiceConfig{
		SubscriptionRepo:      subscriptionRepo,
		SubscriptionPhaseRepo: subscriptionPhaseRepo,
		SubscriptionItemRepo:  subscriptionItemRepo,
		CustomerService:       customerService,
		EntitlementAdapter:    subscriptionEntitlementAdapter,
		FeatureService:        featureConnector,
		TransactionManager:    subscriptionRepo,
		Publisher:             eventPublisher,
		FeatureFlags:          featureFlags,
		Lockr:                 lockr,
		TaxCode:               taxCodeService,
	})
	if err != nil {
		return SubscriptionServiceWithWorkflow{}, err
	}

	subAddRepo := subscriptionaddonrepo.NewSubscriptionAddonRepo(db)
	subAddQtyRepo := subscriptionaddonrepo.NewSubscriptionAddonQuantityRepo(db)

	subAddSvc, err := subscriptionaddonservice.NewService(subscriptionaddonservice.Config{
		TxManager:        subAddRepo,
		Logger:           logger,
		AddonService:     addonService,
		SubService:       subscriptionService,
		SubAddRepo:       subAddRepo,
		SubAddQtyRepo:    subAddQtyRepo,
		PlanAddonService: planAddonService,
		Publisher:        eventPublisher,
	})
	if err != nil {
		return SubscriptionServiceWithWorkflow{}, err
	}

	subscriptionWorkflowService := subscriptionworkflowservice.NewWorkflowService(subscriptionworkflowservice.WorkflowServiceConfig{
		Service:            subscriptionService,
		CustomerService:    customerService,
		TransactionManager: subscriptionRepo,
		AddonService:       subAddSvc,
		Logger:             logger.With("subsystem", "subscription.workflow.service"),
		Lockr:              lockr,
		FeatureFlags:       featureFlags,
	})

	planSubscriptionService := subscriptionchangeservice.New(subscriptionchangeservice.Config{
		WorkflowService:     subscriptionWorkflowService,
		SubscriptionService: subscriptionService,
		PlanService:         planService,
		CustomerService:     customerService,
		Logger:              logger.With("subsystem", "subscription.change.service"),
	})

	validator, err := subscriptioncustomer.NewValidator(subscriptionService, customerService)
	if err != nil {
		return SubscriptionServiceWithWorkflow{}, err
	}

	annotationCleanupHook, err := annotationhook.NewAnnotationCleanupHook(subscriptionService, subscriptionRepo, logger)
	if err != nil {
		return SubscriptionServiceWithWorkflow{}, err
	}

	if err := subscriptionService.RegisterHook(annotationCleanupHook); err != nil {
		return SubscriptionServiceWithWorkflow{}, err
	}

	customerService.RegisterRequestValidator(validator)

	return SubscriptionServiceWithWorkflow{
		Service:                  subscriptionService,
		WorkflowService:          subscriptionWorkflowService,
		PlanSubscriptionService:  planSubscriptionService,
		SubscriptionAddonService: subAddSvc,
	}, nil
}
