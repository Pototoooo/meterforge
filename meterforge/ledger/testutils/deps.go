package testutils

import (
	"log/slog"

	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	ledgeraccount "github.com/Pototoooo/meterforge/meterforge/ledger/account"
	accountadapter "github.com/Pototoooo/meterforge/meterforge/ledger/account/adapter"
	accountservice "github.com/Pototoooo/meterforge/meterforge/ledger/account/service"
	"github.com/Pototoooo/meterforge/meterforge/ledger/historical"
	historicaladapter "github.com/Pototoooo/meterforge/meterforge/ledger/historical/adapter"
	"github.com/Pototoooo/meterforge/meterforge/ledger/resolvers"
	resolversadapter "github.com/Pototoooo/meterforge/meterforge/ledger/resolvers/adapter"
	"github.com/Pototoooo/meterforge/meterforge/ledger/routingrules"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

type Deps struct {
	AccountService   ledgeraccount.Service
	ResolversService *resolvers.AccountResolver
	HistoricalLedger *historical.Ledger
}

func InitDeps(db *entdb.Client, logger *slog.Logger) (Deps, error) {
	repo := accountadapter.NewRepo(db)
	locker, err := lockr.NewLocker(&lockr.LockerConfig{
		Logger: logger,
	})
	if err != nil {
		return Deps{}, err
	}

	historicalRepo := historicaladapter.NewRepo(db)
	accountService := accountservice.New(repo, locker)
	historicalLedger := historical.NewLedger(historicalRepo, accountService, accountService, routingrules.DefaultValidator)
	customerAccountRepo := resolversadapter.NewRepo(db)
	accountResolver := resolvers.NewAccountResolver(resolvers.AccountResolverConfig{
		AccountService: accountService,
		Repo:           customerAccountRepo,
		Locker:         locker,
	})

	return Deps{
		AccountService:   accountService,
		ResolversService: accountResolver,
		HistoricalLedger: historicalLedger,
	}, nil
}
