package hook

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/credit/grant"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	EntitlementHook = models.ServiceHook[entitlement.Entitlement]
	NoopHook        = models.NoopServiceHook[entitlement.Entitlement]
)

type entitlementHook struct {
	NoopHook

	grantRepo grant.Repo
}

func NewEntitlementHook(
	grantRepo grant.Repo,
) EntitlementHook {
	return &entitlementHook{
		grantRepo: grantRepo,
	}
}

func (h *entitlementHook) PreDelete(ctx context.Context, ent *entitlement.Entitlement) error {
	if ent == nil {
		return nil
	}

	meteredEnt, err := meteredentitlement.ParseFromGenericEntitlement(ent)
	if err != nil {
		return nil
	}

	return h.grantRepo.DeleteOwnerGrants(ctx, models.NamespacedID{Namespace: meteredEnt.Namespace, ID: meteredEnt.ID})
}
