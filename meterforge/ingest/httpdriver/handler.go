package httpdriver

import (
	"context"
	"errors"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/ingest"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	IngestHandler
}

type IngestHandler interface {
	IngestEvents() IngestEventsHandler
}

type handler struct {
	service          ingest.Service
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
	service ingest.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		namespaceDecoder: namespaceDecoder,
		service:          service,
		options:          options,
	}
}
