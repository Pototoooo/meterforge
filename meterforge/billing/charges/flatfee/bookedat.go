package flatfee

import (
	"time"

	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	"github.com/Pototoooo/meterforge/pkg/timeutil"
)

// UsageBookedAt returns the ledger booking time for a flat-fee service period.
func UsageBookedAt(paymentTerm productcatalog.PaymentTermType, servicePeriod timeutil.ClosedPeriod) time.Time {
	if paymentTerm == productcatalog.InArrearsPaymentTerm {
		return servicePeriod.To
	}

	return servicePeriod.From
}
