package taxcodes

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListTaxCodes() ListTaxCodesHandler
	GetTaxCode() GetTaxCodeHandler
	CreateTaxCode() CreateTaxCodeHandler
	UpdateTaxCode() UpdateTaxCodeHandler
	DeleteTaxCode() DeleteTaxCodeHandler

	GetOrganizationDefaultTaxCodes() GetOrganizationDefaultTaxCodesHandler
	UpsertOrganizationDefaultTaxCodes() UpsertOrganizationDefaultTaxCodesHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	service          taxcode.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	service taxcode.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		service:          service,
		options:          options,
	}
}
