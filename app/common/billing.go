package common

import (
	"log/slog"

	"github.com/google/wire"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/app"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	billingadapter "github.com/Pototoooo/meterforge/meterforge/billing/adapter"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/creditpurchase"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/flatfee"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/usagebased"
	chargesworkeradvance "github.com/Pototoooo/meterforge/meterforge/billing/charges/worker/advance"
	"github.com/Pototoooo/meterforge/meterforge/billing/rating"
	billingratingservice "github.com/Pototoooo/meterforge/meterforge/billing/rating/service"
	billingsequence "github.com/Pototoooo/meterforge/meterforge/billing/sequence"
	billingsequenceadapter "github.com/Pototoooo/meterforge/meterforge/billing/sequence/adapter"
	billingsequenceservice "github.com/Pototoooo/meterforge/meterforge/billing/sequence/service"
	billingservice "github.com/Pototoooo/meterforge/meterforge/billing/service"
	billingcustomer "github.com/Pototoooo/meterforge/meterforge/billing/validators/customer"
	billingsubscription "github.com/Pototoooo/meterforge/meterforge/billing/validators/subscription"
	billingworkerautoadvance "github.com/Pototoooo/meterforge/meterforge/billing/worker/advance"
	billingworkercollect "github.com/Pototoooo/meterforge/meterforge/billing/worker/collect"
	"github.com/Pototoooo/meterforge/meterforge/billing/worker/subscriptionsync"
	subscriptionsyncadapter "github.com/Pototoooo/meterforge/meterforge/billing/worker/subscriptionsync/adapter"
	"github.com/Pototoooo/meterforge/meterforge/billing/worker/subscriptionsync/reconciler"
	subscriptionsyncservice "github.com/Pototoooo/meterforge/meterforge/billing/worker/subscriptionsync/service"
	"github.com/Pototoooo/meterforge/meterforge/currencies"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/ledger"
	ledgeraccount "github.com/Pototoooo/meterforge/meterforge/ledger/account"
	ledgerbreakage "github.com/Pototoooo/meterforge/meterforge/ledger/breakage"
	"github.com/Pototoooo/meterforge/meterforge/ledger/recognizer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/featuregate"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

// BillingRegistry bundles the billing and charges services. External callers that need
// billing or charges should depend on BillingRegistry rather than individual services.
type BillingRegistry struct {
	Billing  billing.Service
	Sequence billingsequence.Service
	Charges  *ChargesRegistry
}

func (r BillingRegistry) ChargesServiceOrNil() charges.Service {
	if r.Charges == nil {
		return nil
	}

	return r.Charges.Service
}

// ChargesRegistry groups all charge-type services.
type ChargesRegistry struct {
	Service               charges.Service
	FlatFeeService        flatfee.Service
	UsageBasedService     usagebased.Service
	CreditPurchaseService creditpurchase.Service
	RecognizerService     recognizer.Service
}

// Billing is the Wire provider set for the billing and charges stack.
// Downstream consumers should depend on BillingRegistry.
var Billing = wire.NewSet(
	BillingAdapter,
	NewBillingRatingService,
	NewLedgerBreakageService,
	NewBillingRegistry,
	NewBillingCustomerOverrideService,
)

func BillingAdapter(
	logger *slog.Logger,
	db *entdb.Client,
) (billing.Adapter, error) {
	return billingadapter.New(billingadapter.Config{
		Client: db,
		Logger: logger,
	})
}

func NewBillingSequenceService(
	logger *slog.Logger,
	db *entdb.Client,
	metricMeter otelmetric.Meter,
) (billingsequence.Service, error) {
	adapter, err := billingsequenceadapter.New(billingsequenceadapter.Config{
		Client: db,
		Logger: logger.With("subsystem", "billing.sequence"),
	})
	if err != nil {
		return nil, err
	}

	return billingsequenceservice.New(billingsequenceservice.Config{
		Adapter: adapter,
		Meter:   metricMeter,
	})
}

// newBillingService creates the billing service and registers validators/hooks.
// Downstream consumers should use BillingRegistry.
func newBillingService(
	logger *slog.Logger,
	appService app.Service,
	billingAdapter billing.Adapter,
	billingRatingService rating.Service,
	sequenceService billingsequence.Service,
	customerService customer.Service,
	featureConnector feature.FeatureConnector,
	meterService meter.Service,
	streamingConnector streaming.Connector,
	eventPublisher eventbus.Publisher,
	billingConfig config.BillingConfiguration,
	subscriptionServices SubscriptionServiceWithWorkflow,
	db *entdb.Client,
	fsConfig config.BillingFeatureSwitchesConfiguration,
	tracer trace.Tracer,
	taxCodeService taxcode.Service,
) (billing.Service, error) {
	service, err := billingservice.New(billingservice.Config{
		Adapter:                      billingAdapter,
		SequenceService:              sequenceService,
		RatingService:                billingRatingService,
		AppService:                   appService,
		CustomerService:              customerService,
		TaxCodeService:               taxCodeService,
		FeatureService:               featureConnector,
		Logger:                       logger,
		MeterService:                 meterService,
		StreamingConnector:           streamingConnector,
		Publisher:                    eventPublisher,
		AdvancementStrategy:          billingConfig.AdvancementStrategy,
		FSNamespaceLockdown:          fsConfig.NamespaceLockdown,
		MaxParallelQuantitySnapshots: billingConfig.MaxParallelQuantitySnapshots,
	})
	if err != nil {
		return nil, err
	}

	return service, nil
}

// NewBillingRegistry assembles the billing and optional charges services.
func NewBillingRegistry(
	logger *slog.Logger,
	appService app.Service,
	billingAdapter billing.Adapter,
	billingRatingService rating.Service,
	customerService customer.Service,
	featureConnector feature.FeatureConnector,
	meterService meter.Service,
	metricMeter otelmetric.Meter,
	streamingConnector streaming.Connector,
	eventPublisher eventbus.Publisher,
	billingConfig config.BillingConfiguration,
	subscriptionServices SubscriptionServiceWithWorkflow,
	db *entdb.Client,
	fsConfig config.BillingFeatureSwitchesConfiguration,
	creditsConfig config.CreditsConfiguration,
	tracer trace.Tracer,
	taxCodeService taxcode.Service,
	currencyResolver currencies.CurrencyResolver,
	locker *lockr.Locker,
	ledgerService ledger.Ledger,
	balanceQuerier ledger.BalanceQuerier,
	accountResolver ledger.AccountResolver,
	accountService ledgeraccount.Service,
	breakageService ledgerbreakage.Service,
	featureGate *featuregate.FeatureGateChecker,
) (BillingRegistry, error) {
	sequenceService, err := NewBillingSequenceService(logger, db, metricMeter)
	if err != nil {
		return BillingRegistry{}, err
	}

	billingService, err := newBillingService(
		logger,
		appService,
		billingAdapter,
		billingRatingService,
		sequenceService,
		customerService,
		featureConnector,
		meterService,
		streamingConnector,
		eventPublisher,
		billingConfig,
		subscriptionServices,
		db,
		fsConfig,
		tracer,
		taxCodeService,
	)
	if err != nil {
		return BillingRegistry{}, err
	}

	var chargesRegistry *ChargesRegistry

	if creditsConfig.Enabled {
		chargesRegistry, err = newChargesRegistry(
			db,
			logger,
			locker,
			billingService,
			billingRatingService,
			featureConnector,
			streamingConnector,
			ledgerService,
			balanceQuerier,
			accountResolver,
			accountService,
			breakageService,
			taxCodeService,
			currencyResolver,
			fsConfig.NamespaceLockdown,
			creditsConfig,
			featureGate,
		)
		if err != nil {
			return BillingRegistry{}, err
		}
	}

	billingRegistry := BillingRegistry{
		Billing:  billingService,
		Sequence: sequenceService,
		Charges:  chargesRegistry,
	}

	// Hook registration

	// Customer validate (and sync subscription on delete)
	// To prevent circular dependencies, we create the validator here
	subscriptionSyncAdapter, err := NewBillingSubscriptionSyncAdapter(db)
	if err != nil {
		return BillingRegistry{}, err
	}
	subscriptionSyncService, err := NewBillingSubscriptionSyncService(logger, subscriptionServices, billingRegistry, subscriptionSyncAdapter, tracer, creditsConfig, fsConfig, featureGate)
	if err != nil {
		return BillingRegistry{}, err
	}

	validator, err := billingcustomer.NewValidator(billingRegistry.Billing, subscriptionSyncService, subscriptionServices.Service)
	if err != nil {
		return BillingRegistry{}, err
	}

	customerService.RegisterRequestValidator(validator)

	// Subscription validate

	subscriptionValidator, err := billingsubscription.NewValidator(billingRegistry.Billing)
	if err != nil {
		return BillingRegistry{}, err
	}

	if err = subscriptionServices.Service.RegisterHook(subscriptionValidator); err != nil {
		return BillingRegistry{}, err
	}

	return billingRegistry, nil
}

func NewBillingCustomerOverrideService(billingRegistry BillingRegistry) billing.CustomerOverrideService {
	return billingRegistry.Billing
}

func NewBillingRatingService(unitConfig config.UnitConfigConfiguration) rating.Service {
	return billingratingservice.New(billingratingservice.Config{
		UnitConfigEnabled: unitConfig.Enabled,
	})
}

func NewBillingAutoAdvancer(logger *slog.Logger, billingRegistry BillingRegistry) (*billingworkerautoadvance.AutoAdvancer, error) {
	return billingworkerautoadvance.NewAdvancer(billingworkerautoadvance.Config{
		BillingService: billingRegistry.Billing,
		Logger:         logger,
	})
}

func NewChargesAutoAdvancer(logger *slog.Logger, billingRegistry BillingRegistry) (*chargesworkeradvance.AutoAdvancer, error) {
	chargesService := billingRegistry.ChargesServiceOrNil()
	if chargesService == nil {
		return nil, nil
	}

	return chargesworkeradvance.NewAdvancer(chargesworkeradvance.Config{
		ChargesService: chargesService,
		Logger:         logger,
	})
}

func NewBillingCollector(logger *slog.Logger, billingRegistry BillingRegistry, fs config.BillingFeatureSwitchesConfiguration) (*billingworkercollect.InvoiceCollector, error) {
	return billingworkercollect.NewInvoiceCollector(billingworkercollect.Config{
		GatheringInvoiceService: billingRegistry.Billing,
		BillingService:          billingRegistry.Billing,
		Logger:                  logger,
		LockedNamespaces:        fs.NamespaceLockdown,
		MaxLinesPerInvoice:      fs.MaxLinesPerCollectedInvoice,
	})
}

func NewBillingSubscriptionReconciler(logger *slog.Logger, subsServices SubscriptionServiceWithWorkflow, subscriptionSync subscriptionsync.Service, customerService customer.Service) (*reconciler.Reconciler, error) {
	return reconciler.NewReconciler(reconciler.ReconcilerConfig{
		SubscriptionService: subsServices.Service,
		SubscriptionSync:    subscriptionSync,
		Logger:              logger,
		CustomerService:     customerService,
	})
}

func NewBillingSubscriptionSyncAdapter(db *entdb.Client) (subscriptionsync.Adapter, error) {
	return subscriptionsyncadapter.New(subscriptionsyncadapter.Config{
		Client: db,
	})
}

func NewBillingSubscriptionSyncService(logger *slog.Logger, subsServices SubscriptionServiceWithWorkflow, billingRegistry BillingRegistry, subscriptionSyncAdapter subscriptionsync.Adapter, tracer trace.Tracer, creditsConfig config.CreditsConfiguration, billingFsConfig config.BillingFeatureSwitchesConfiguration, featureGate *featuregate.FeatureGateChecker) (subscriptionsync.Service, error) {
	return subscriptionsyncservice.New(subscriptionsyncservice.Config{
		SubscriptionService:     subsServices.Service,
		BillingService:          billingRegistry.Billing,
		ChargesService:          billingRegistry.ChargesServiceOrNil(),
		SubscriptionSyncAdapter: subscriptionSyncAdapter,
		FeatureFlags: subscriptionsyncservice.FeatureFlags{
			EnableCreditThenInvoice:     creditsConfig.EnableCreditThenInvoice,
			MaxLinesPerCollectedInvoice: billingFsConfig.MaxLinesPerCollectedInvoice,
		},
		ForceAsyncInvoicePendingLines: billingFsConfig.SubscriptionSyncForceAsyncAdvance,
		Logger:                        logger,
		Tracer:                        tracer,
		FeatureGate:                   featureGate,
	})
}
