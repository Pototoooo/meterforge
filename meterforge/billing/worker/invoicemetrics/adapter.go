package invoicemetrics

import (
	"context"

	"github.com/Pototoooo/meterforge/pkg/framework/entutils"
)

type Adapter interface {
	entutils.TxCreator

	CountOverdueInvoices(ctx context.Context, input CountOverdueInvoicesInput) (OverdueInvoiceCounts, error)
}
