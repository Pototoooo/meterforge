package subscription

import "github.com/Pototoooo/meterforge/pkg/framework/lockr"

func GetCustomerLock(customerId string) (lockr.Key, error) {
	return lockr.NewKey("customer", customerId, "subscription")
}
