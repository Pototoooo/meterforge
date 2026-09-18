package service

import "github.com/Pototoooo/meterforge/pkg/framework/lockr"

func NewEntitlementUniqueScopeLock(featureKey string, customerID string) (lockr.Key, error) {
	return lockr.NewKey("fk", featureKey, "cid", customerID)
}
