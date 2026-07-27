package billinginvoices

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListBillingInvoices() ListBillingInvoicesHandler
	GetBillingInvoice() GetBillingInvoiceHandler
	UpdateBillingInvoice() UpdateBillingInvoiceHandler
	DeleteBillingInvoice() DeleteBillingInvoiceHandler
	ProgressInvoice(action ProgressAction) ProgressInvoiceHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	service          billing.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	service billing.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		service:          service,
		options:          options,
	}
}
