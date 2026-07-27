package customerservice

import "github.com/Pototoooo/meterforge/meterforge/customer"

var _ customer.RequestValidatorService = (*Service)(nil)

func (s *Service) RegisterRequestValidator(v customer.RequestValidator) {
	s.requestValidatorRegistry.Register(v)
}
