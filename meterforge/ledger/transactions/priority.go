package transactions

import (
	"github.com/Pototoooo/meterforge/meterforge/ledger"
)

func resolveCustomerFBOCreditPriority(configured *int) int {
	if configured != nil {
		return *configured
	}
	return ledger.DefaultCustomerFBOPriority
}
