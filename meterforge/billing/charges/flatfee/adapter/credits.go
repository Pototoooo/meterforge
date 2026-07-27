package adapter

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/Pototoooo/meterforge/meterforge/billing/charges/flatfee"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/models/creditrealization"
	"github.com/Pototoooo/meterforge/meterforge/ent/db"
	dbchargeflatfeerun "github.com/Pototoooo/meterforge/meterforge/ent/db/chargeflatfeerun"
	"github.com/Pototoooo/meterforge/pkg/framework/entutils"
	"github.com/Pototoooo/meterforge/pkg/slicesx"
)

var _ flatfee.ChargeCreditAllocationAdapter = (*adapter)(nil)

func (a *adapter) CreateCreditAllocations(ctx context.Context, runID flatfee.RealizationRunID, creditAllocations creditrealization.CreateInputs) (creditrealization.Realizations, error) {
	if err := runID.Validate(); err != nil {
		return creditrealization.Realizations{}, err
	}

	if err := creditAllocations.Validate(); err != nil {
		return creditrealization.Realizations{}, err
	}

	return entutils.TransactingRepo(ctx, a, func(ctx context.Context, tx *adapter) (creditrealization.Realizations, error) {
		if _, err := tx.db.ChargeFlatFeeRun.Query().
			Where(
				dbchargeflatfeerun.NamespaceEQ(runID.Namespace),
				dbchargeflatfeerun.IDEQ(runID.ID),
			).
			Only(ctx); err != nil {
			return creditrealization.Realizations{}, fmt.Errorf("querying flat fee run [run_id=%s]: %w", runID.ID, err)
		}

		dbEntities, err := tx.db.ChargeFlatFeeRunCreditAllocations.CreateBulk(
			lo.Map(creditAllocations, func(creditAllocation creditrealization.CreateInput, idx int) *db.ChargeFlatFeeRunCreditAllocationsCreate {
				create := tx.db.ChargeFlatFeeRunCreditAllocations.Create().
					SetRunID(runID.ID)

				create = creditrealization.Create(create, runID.Namespace, idx, creditAllocation)

				return create
			})...,
		).Save(ctx)
		if err != nil {
			return creditrealization.Realizations{}, err
		}

		realizations, err := slicesx.MapWithErr(dbEntities, func(entity *db.ChargeFlatFeeRunCreditAllocations) (creditrealization.Realization, error) {
			return creditrealization.MapFromDB(entity), nil
		})
		if err != nil {
			return creditrealization.Realizations{}, err
		}

		return realizations, nil
	})
}
