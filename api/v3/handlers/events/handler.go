package events

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/ingest"
	"github.com/Pototoooo/meterforge/meterforge/meterevent"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	IngestEvents() IngestEventsHandler
	ListMeteringEvents() ListMeteringEventsHandler
}

type handler struct {
	resolveNamespace  func(ctx context.Context) (string, error)
	service           ingest.Service
	metereventService meterevent.Service
	options           []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	service ingest.Service,
	metereventService meterevent.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace:  resolveNamespace,
		service:           service,
		metereventService: metereventService,
		options:           options,
	}
}
