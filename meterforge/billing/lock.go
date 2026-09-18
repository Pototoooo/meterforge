package billing

import "github.com/Pototoooo/meterforge/meterforge/customer"

type (
	LockCustomerForUpdateAdapterInput = customer.CustomerID
	UpsertCustomerLockAdapterInput    = LockCustomerForUpdateAdapterInput
)
