package customer

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/pkg/models"
	"github.com/Pototoooo/meterforge/pkg/pagination"
)

type Service interface {
	CustomerService
	RequestValidatorService

	models.ServiceHooks[Customer]
}

type RequestValidatorService interface {
	RegisterRequestValidator(RequestValidator)
}

type CustomerService interface {
	ListCustomers(ctx context.Context, params ListCustomersInput) (pagination.Result[Customer], error)
	ListCustomerUsageAttributions(ctx context.Context, input ListCustomerUsageAttributionsInput) (pagination.Result[streaming.CustomerUsageAttribution], error)
	CreateCustomer(ctx context.Context, params CreateCustomerInput) (*Customer, error)
	DeleteCustomer(ctx context.Context, customer DeleteCustomerInput) error
	GetCustomer(ctx context.Context, customer GetCustomerInput) (*Customer, error)
	GetCustomerByUsageAttribution(ctx context.Context, input GetCustomerByUsageAttributionInput) (*Customer, error)
	GetCustomersByUsageAttribution(ctx context.Context, input GetCustomersByUsageAttributionInput) ([]Customer, error)
	UpdateCustomer(ctx context.Context, params UpdateCustomerInput) (*Customer, error)
}
