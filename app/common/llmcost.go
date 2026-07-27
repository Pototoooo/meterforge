package common

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/wire"

	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/llmcost"
	llmcostadapter "github.com/Pototoooo/meterforge/meterforge/llmcost/adapter"
	llmcostservice "github.com/Pototoooo/meterforge/meterforge/llmcost/service"
	llmcostsync "github.com/Pototoooo/meterforge/meterforge/llmcost/sync"
)

var LLMCost = wire.NewSet(
	NewLLMCostService,
	NewLLMCostSyncJob,
)

func NewLLMCostService(logger *slog.Logger, db *entdb.Client) (llmcost.Service, error) {
	adapter, err := llmcostadapter.New(llmcostadapter.Config{
		Client: db,
		Logger: logger.With("subsystem", "llmcost"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize llmcost adapter: %w", err)
	}

	return llmcostservice.New(adapter, logger.With("subsystem", "llmcost")), nil
}

func NewLLMCostSyncJob(logger *slog.Logger, db *entdb.Client) (*llmcostsync.SyncJob, error) {
	adapter, err := llmcostadapter.New(llmcostadapter.Config{
		Client: db,
		Logger: logger.With("subsystem", "llmcost.sync"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize llmcost sync adapter: %w", err)
	}

	return llmcostsync.NewSyncJob(llmcostsync.SyncJobConfig{
		HTTPClient: &http.Client{},
		Repo:       adapter,
		Logger:     logger.With("subsystem", "llmcost.sync"),
	}), nil
}
