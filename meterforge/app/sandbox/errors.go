package appsandbox

import "github.com/Pototoooo/meterforge/meterforge/billing"

var ErrSimulatedPaymentFailure = billing.NewValidationError("simulated_payment_failure", "simulated payment failure")
