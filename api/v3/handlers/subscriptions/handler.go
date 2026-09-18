package subscriptions

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	plansubscription "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListSubscriptions() ListSubscriptionsHandler
	GetSubscription() GetSubscriptionHandler
	CreateSubscription() CreateSubscriptionHandler
	CancelSubscription() CancelSubscriptionHandler
	UnscheduleCancelation() UnscheduleCancelationHandler
	ChangeSubscription() ChangeSubscriptionHandler
}

type handler struct {
	resolveNamespace        func(ctx context.Context) (string, error)
	customerService         customer.Service
	planService             plan.Service
	planSubscriptionService plansubscription.PlanSubscriptionService
	subscriptionService     subscription.Service
	options                 []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	customerService customer.Service,
	planService plan.Service,
	planSubscriptionService plansubscription.PlanSubscriptionService,
	subscriptionService subscription.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace:        resolveNamespace,
		customerService:         customerService,
		planService:             planService,
		planSubscriptionService: planSubscriptionService,
		subscriptionService:     subscriptionService,
		options:                 options,
	}
}
