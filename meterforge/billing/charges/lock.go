package charges

import (
	"fmt"

	"github.com/Pototoooo/meterforge/meterforge/billing/charges/meta"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

func NewLockKeyForCharge(chargeID meta.ChargeID) (lockr.Key, error) {
	if err := chargeID.Validate(); err != nil {
		return nil, fmt.Errorf("charge ID: %w", err)
	}

	return lockr.NewKey("namespace", chargeID.Namespace, "charge", chargeID.ID)
}
