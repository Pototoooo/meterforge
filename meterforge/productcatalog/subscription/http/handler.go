package httpdriver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	appconfig "github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	plansubscription "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionworkflow "github.com/Pototoooo/meterforge/meterforge/subscription/workflow"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	CreateSubscription() CreateSubscriptionHandler
	GetSubscription() GetSubscriptionHandler
	EditSubscription() EditSubscriptionHandler
	CancelSubscription() CancelSubscriptionHandler
	ContinueSubscription() ContinueSubscriptionHandler
	RestoreSubscription() RestoreSubscriptionHandler
	MigrateSubscription() MigrateSubscriptionHandler
	ChangeSubscription() ChangeSubscriptionHandler
	DeleteSubscription() DeleteSubscriptionHandler
	ListCustomerSubscriptions() ListCustomerSubscriptionsHandler
}

type HandlerConfig struct {
	SubscriptionWorkflowService subscriptionworkflow.Service
	SubscriptionService         subscription.Service
	CustomerService             customer.Service
	PlanSubscriptionService     plansubscription.PlanSubscriptionService
	NamespaceDecoder            namespacedriver.NamespaceDecoder
	Logger                      *slog.Logger
	Credits                     appconfig.CreditsConfiguration
}

type handler struct {
	HandlerConfig
	Options []httptransport.HandlerOption
}

func (h *handler) resolveNamespace(ctx context.Context) (string, error) {
	ns, ok := h.NamespaceDecoder.GetNamespace(ctx)
	if !ok {
		return "", commonhttp.NewHTTPError(http.StatusInternalServerError, errors.New("internal server error"))
	}

	return ns, nil
}

func NewHandler(config HandlerConfig, options ...httptransport.HandlerOption) Handler {
	return &handler{
		HandlerConfig: config,
		Options:       options,
	}
}
