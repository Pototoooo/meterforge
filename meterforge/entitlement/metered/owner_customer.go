package meteredentitlement

import "github.com/Pototoooo/meterforge/meterforge/streaming"

// ownerCustomer is a lightweight adapter that implements streaming.Customer
// without introducing customer package dependencies into credit.
type ownerCustomer struct {
	id          string
	key         *string
	subjectKeys []string
}

var _ streaming.Customer = ownerCustomer{}

func (c ownerCustomer) GetUsageAttribution() streaming.CustomerUsageAttribution {
	return streaming.NewCustomerUsageAttribution(c.id, c.key, c.subjectKeys)
}
