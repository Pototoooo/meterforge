package customerbalance

import (
	"context"
	"time"

	"github.com/samber/mo"

	"github.com/Pototoooo/meterforge/meterforge/billing/charges/creditpurchase"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/ledger"
	"github.com/Pototoooo/meterforge/pkg/currencyx"
)

type creditTransactionLoaderInput struct {
	Limit         int
	After         *ledger.TransactionCursor
	Before        *ledger.TransactionCursor
	CustomerID    customer.CustomerID
	AccountID     string
	Currency      *currencyx.Code
	AsOf          time.Time
	FeatureFilter mo.Option[creditpurchase.FeatureFilters]
}

type creditTransactionLoaderResult struct {
	Items   []CreditTransaction
	HasMore bool
}

type creditTransactionLoader interface {
	Load(ctx context.Context, input creditTransactionLoaderInput) (creditTransactionLoaderResult, error)
}

type creditTransactionLoaderFactory func(*service) creditTransactionLoader

var creditTransactionLoaderOrder = []CreditTransactionType{
	CreditTransactionTypeFunded,
	CreditTransactionTypeConsumed,
	CreditTransactionTypeExpired,
	CreditTransactionTypeVoided,
}

var creditTransactionLoaderFactories = map[CreditTransactionType]creditTransactionLoaderFactory{
	CreditTransactionTypeFunded: newFundedCreditTransactionLoader,
	CreditTransactionTypeConsumed: func(s *service) creditTransactionLoader {
		return newLedgerCreditTransactionLoader(s, ledger.ListTransactionsCreditMovementNegative)
	},
	CreditTransactionTypeExpired: newExpiredCreditTransactionLoader,
	CreditTransactionTypeVoided:  newVoidedCreditTransactionLoader,
}

func (s *service) creditTransactionLoaders(txType *CreditTransactionType) ([]creditTransactionLoader, error) {
	if txType == nil {
		loaders := make([]creditTransactionLoader, 0, len(creditTransactionLoaderOrder))
		for _, transactionType := range creditTransactionLoaderOrder {
			loaders = append(loaders, creditTransactionLoaderFactories[transactionType](s))
		}

		return loaders, nil
	}

	if err := txType.Validate(); err != nil {
		return nil, err
	}

	factory, ok := creditTransactionLoaderFactories[*txType]
	if !ok {
		return nil, txType.Validate()
	}

	return []creditTransactionLoader{factory(s)}, nil
}
