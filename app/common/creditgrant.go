package common

import (
	"fmt"

	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/billing/creditgrant"
	creditgrantservice "github.com/Pototoooo/meterforge/meterforge/billing/creditgrant/service"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	enttx "github.com/Pototoooo/meterforge/meterforge/ent/tx"
	"github.com/Pototoooo/meterforge/meterforge/ledger"
	ledgeraccount "github.com/Pototoooo/meterforge/meterforge/ledger/account"
	ledgerbreakage "github.com/Pototoooo/meterforge/meterforge/ledger/breakage"
	"github.com/Pototoooo/meterforge/meterforge/ledger/creditvoid"
	creditvoidadapter "github.com/Pototoooo/meterforge/meterforge/ledger/creditvoid/adapter"
	"github.com/Pototoooo/meterforge/meterforge/ledger/transactions"
)

var CreditGrant = wire.NewSet(
	NewCreditVoidService,
	NewCreditGrantService,
)

func NewCreditVoidService(
	creditsConfig config.CreditsConfiguration,
	db *entdb.Client,
	ledgerService ledger.Ledger,
	balanceQuerier ledger.BalanceQuerier,
	accountResolver ledger.AccountResolver,
	accountService ledgeraccount.Service,
	breakageService ledgerbreakage.Service,
) (creditvoid.Service, error) {
	if !creditsConfig.Enabled {
		return creditvoid.NewNoopService(), nil
	}

	adapter, err := creditvoidadapter.New(creditvoidadapter.Config{
		Client: db,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credit void adapter: %w", err)
	}

	svc, err := creditvoid.NewService(creditvoid.Config{
		Adapter: adapter,
		Ledger:  ledgerService,
		Dependencies: transactions.ResolverDependencies{
			AccountService: accountResolver,
			AccountCatalog: accountService,
			BalanceQuerier: balanceQuerier,
		},
		Breakage:           breakageService,
		AccountLocker:      accountService,
		TransactionManager: enttx.NewCreator(db),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credit void service: %w", err)
	}

	return svc, nil
}

func NewCreditGrantService(
	db *entdb.Client,
	billingRegistry BillingRegistry,
	customerService customer.Service,
	creditVoidService creditvoid.Service,
) (creditgrant.Service, error) {
	if billingRegistry.Charges == nil {
		return nil, nil
	}

	svc, err := creditgrantservice.New(creditgrantservice.Config{
		CreditPurchaseService: billingRegistry.Charges.CreditPurchaseService,
		ChargesService:        billingRegistry.Charges.Service,
		BillingService:        billingRegistry.Billing,
		CustomerService:       customerService,
		CreditVoidService:     creditVoidService,
		TransactionManager:    enttx.NewCreator(db),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credit grant service: %w", err)
	}

	return svc, nil
}
