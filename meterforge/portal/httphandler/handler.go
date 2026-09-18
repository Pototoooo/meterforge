package httpdriver

import (
	"context"
	"errors"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/meterforge/portal"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	TokenHandler
}

type TokenHandler interface {
	CreateToken() CreateTokenHandler
	ListTokens() ListTokensHandler
	InvalidateToken() InvalidateTokenHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	namespaceDecoder namespacedriver.NamespaceDecoder
	options          []httptransport.HandlerOption
	portalService    portal.Service
	meterService     meter.Service
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
	portalService portal.Service,
	meterService meter.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		namespaceDecoder: namespaceDecoder,
		portalService:    portalService,
		meterService:     meterService,
		options:          options,
	}
}
