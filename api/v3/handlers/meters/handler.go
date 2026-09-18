package meters

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListMeters() ListMetersHandler
	GetMeter() GetMeterHandler
	CreateMeter() CreateMeterHandler
	UpdateMeter() UpdateMeterHandler
	DeleteMeter() DeleteMeterHandler
	QueryMeter() QueryMeterHandler
	QueryMeterCSV() QueryMeterCSVHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	service          meter.ManageService
	streaming        streaming.Connector
	customerService  customer.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	service meter.ManageService,
	streaming streaming.Connector,
	customerService customer.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		service:          service,
		streaming:        streaming,
		customerService:  customerService,
		options:          options,
	}
}
