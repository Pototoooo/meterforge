package httpdriver

import (
	"context"
	"errors"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	AddonHandler
}

type AddonHandler interface {
	ListAddons() ListAddonsHandler
	CreateAddon() CreateAddonHandler
	DeleteAddon() DeleteAddonHandler
	GetAddon() GetAddonHandler
	UpdateAddon() UpdateAddonHandler
	PublishAddon() PublishAddonHandler
	ArchiveAddon() ArchiveAddonHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	service          addon.Service
	namespaceDecoder namespacedriver.NamespaceDecoder
	options          []httptransport.HandlerOption
}

func (h *handler) resolveNamespace(ctx context.Context) (string, error) {
	ns, ok := h.namespaceDecoder.GetNamespace(ctx)
	if !ok {
		return "", commonhttp.NewHTTPError(http.StatusInternalServerError, errors.New("internal server error"))
	}

	return ns, nil
}

func New(
	namespaceDecoder namespacedriver.NamespaceDecoder,
	service addon.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		service:          service,
		namespaceDecoder: namespaceDecoder,
		options:          options,
	}
}
