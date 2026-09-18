package adapter

import (
	"context"

	"github.com/samber/lo"

	"github.com/Pototoooo/meterforge/meterforge/billing/charges/models/creditrealization"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/usagebased"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/pkg/framework/entutils"
)

var _ usagebased.RealizationRunCreditAllocationAdapter = (*adapter)(nil)

func (a *adapter) CreateRunCreditRealization(ctx context.Context, runID usagebased.RealizationRunID, creditAllocations creditrealization.CreateInputs) (creditrealization.Realizations, error) {
	if err := runID.Validate(); err != nil {
		return nil, err
	}

	if err := creditAllocations.Validate(); err != nil {
		return nil, err
	}

	return entutils.TransactingRepo(ctx, a, func(ctx context.Context, tx *adapter) (creditrealization.Realizations, error) {
		creates := lo.Map(creditAllocations, func(creditAllocation creditrealization.CreateInput, idx int) *entdb.ChargeUsageBasedRunCreditAllocationsCreate {
			create := tx.db.ChargeUsageBasedRunCreditAllocations.Create().
				SetRunID(runID.ID).
				SetNamespace(runID.Namespace)

			create = creditrealization.Create[*entdb.ChargeUsageBasedRunCreditAllocationsCreate](create, runID.Namespace, idx, creditAllocation)

			return create
		})

		dbEntities, err := tx.db.ChargeUsageBasedRunCreditAllocations.CreateBulk(creates...).Save(ctx)
		if err != nil {
			return nil, err
		}

		realizations := lo.Map(dbEntities, func(entity *entdb.ChargeUsageBasedRunCreditAllocations, _ int) creditrealization.Realization {
			return creditrealization.MapFromDB(entity)
		})

		return realizations, nil
	})
}
