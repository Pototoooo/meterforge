package account

import (
	"github.com/Pototoooo/meterforge/meterforge/ledger"
)

type Service interface {
	ledger.AccountCatalog
	ledger.AccountLocker
}

type (
	ListAccountsInput     = ledger.ListAccountsInput
	ListSubAccountsInput  = ledger.ListSubAccountsInput
	CreateAccountInput    = ledger.CreateAccountInput
	CreateSubAccountInput = ledger.CreateSubAccountInput
)
