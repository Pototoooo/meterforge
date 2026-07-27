package httphandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/meterevent"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	EventHandler
}

type EventHandler interface {
	ListEvents() ListEventsHandler
	ListEventsV2() ListEventsV2Handler
}

var _ Handler = (*handler)(nil)

type handler struct {
	namespaceDecoder  namespacedriver.NamespaceDecoder
	options           []httptransport.HandlerOption
	metereventService meterevent.Service
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
	metereventService meterevent.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		namespaceDecoder:  namespaceDecoder,
		metereventService: metereventService,
		options:           options,
	}
}
