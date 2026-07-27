package httpdriver

import (
	"context"
	"errors"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	MeterHandler
}

type MeterHandler interface {
	ListMeters() ListMetersHandler
	GetMeter() GetMeterHandler
	CreateMeter() CreateMeterHandler
	UpdateMeter() UpdateMeterHandler
	DeleteMeter() DeleteMeterHandler
	QueryMeter() QueryMeterHandler
	QueryMeterPost() QueryMeterPostHandler
	QueryMeterPostCSV() QueryMeterPostCSVHandler
	QueryMeterCSV() QueryMeterCSVHandler
	ListSubjects() ListSubjectsHandler
	ListGroupByValues() ListGroupByValuesHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	namespaceDecoder namespacedriver.NamespaceDecoder
	options          []httptransport.HandlerOption
	customerService  customer.Service
	meterService     meter.ManageService
	streaming        streaming.Connector
	subjectService   subject.Service
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
	customerService customer.Service,
	meterService meter.ManageService,
	streaming streaming.Connector,
	subjectService subject.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		namespaceDecoder: namespaceDecoder,
		options:          options,
		customerService:  customerService,
		meterService:     meterService,
		streaming:        streaming,
		subjectService:   subjectService,
	}
}
