package httpdriver

import (
	"context"
	"errors"
	"net/http"

	appcustominvoicing "github.com/Pototoooo/meterforge/meterforge/app/custominvoicing"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	AppHandler
}

type AppHandler interface {
	DraftSyncronized() DraftSyncronizedHandler
	IssuingSyncronized() IssuingSyncronizedHandler
	UpdatePaymentStatus() UpdatePaymentStatusHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	service appcustominvoicing.SyncService

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
	service appcustominvoicing.SyncService,
	namespaceDecoder namespacedriver.NamespaceDecoder,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		service:          service,
		namespaceDecoder: namespaceDecoder,
		options:          options,
	}
}
