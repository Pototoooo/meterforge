package httpdriver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionaddon "github.com/Pototoooo/meterforge/meterforge/subscription/addon"
	subscriptionworkflow "github.com/Pototoooo/meterforge/meterforge/subscription/workflow"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	CreateSubscriptionAddon() CreateSubscriptionAddonHandler
	ListSubscriptionAddons() ListSubscriptionAddonsHandler
	GetSubscriptionAddon() GetSubscriptionAddonHandler
	UpdateSubscriptionAddon() UpdateSubscriptionAddonHandler
}

type HandlerConfig struct {
	SubscriptionAddonService    subscriptionaddon.Service
	SubscriptionWorkflowService subscriptionworkflow.Service
	SubscriptionService         subscription.Service
	AddonService                addon.Service
	NamespaceDecoder            namespacedriver.NamespaceDecoder
	Logger                      *slog.Logger
}

func NewHandler(config HandlerConfig, options ...httptransport.HandlerOption) Handler {
	return &handler{
		HandlerConfig: config,
		Options:       options,
	}
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
