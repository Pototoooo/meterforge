package adapter

import (
	"context"
	"fmt"

	"github.com/Pototoooo/meterforge/meterforge/billing/charges/flatfee"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/models/invoicedusage"
	dbchargeflatfeerun "github.com/Pototoooo/meterforge/meterforge/ent/db/chargeflatfeerun"
	"github.com/Pototoooo/meterforge/pkg/framework/entutils"
)

var _ flatfee.ChargeInvoicedUsageAdapter = (*adapter)(nil)

func (a *adapter) CreateInvoicedUsage(ctx context.Context, input flatfee.CreateInvoicedUsageInput) (invoicedusage.AccruedUsage, error) {
	if err := input.Validate(); err != nil {
		return invoicedusage.AccruedUsage{}, err
	}

	return entutils.TransactingRepo(ctx, a, func(ctx context.Context, tx *adapter) (invoicedusage.AccruedUsage, error) {
		if _, err := tx.db.ChargeFlatFeeRun.UpdateOneID(input.RunID.ID).
			Where(dbchargeflatfeerun.Namespace(input.RunID.Namespace)).
			SetLineID(input.LineID).
			SetInvoiceID(input.InvoiceID).
			Save(ctx); err != nil {
			return invoicedusage.AccruedUsage{}, fmt.Errorf("updating flat fee run invoice refs [run_id=%s]: %w", input.RunID.ID, err)
		}

		create := tx.db.ChargeFlatFeeRunInvoicedUsage.Create().
			SetRunID(input.RunID.ID)

		create = invoicedusage.Create(create, input.RunID.Namespace, input.InvoicedUsage)

		entity, err := create.Save(ctx)
		if err != nil {
			return invoicedusage.AccruedUsage{}, err
		}

		return invoicedusage.MapAccruedUsageFromDB(entity), nil
	})
}
