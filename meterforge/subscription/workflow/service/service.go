package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionaddon "github.com/Pototoooo/meterforge/meterforge/subscription/addon"
	subscriptionworkflow "github.com/Pototoooo/meterforge/meterforge/subscription/workflow"
	"github.com/Pototoooo/meterforge/pkg/ffx"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
	"github.com/Pototoooo/meterforge/pkg/framework/transaction"
)

type WorkflowServiceConfig struct {
	Service      subscription.Service
	AddonService subscriptionaddon.Service
	// connectors
	CustomerService customer.Service
	// framework
	TransactionManager transaction.Creator
	Logger             *slog.Logger
	Lockr              *lockr.Locker
	FeatureFlags       ffx.Service
}

type service struct {
	WorkflowServiceConfig
}

func NewWorkflowService(cfg WorkflowServiceConfig) subscriptionworkflow.Service {
	return &service{
		WorkflowServiceConfig: cfg,
	}
}

var _ subscriptionworkflow.Service = &service{}

func (s *service) lockCustomer(ctx context.Context, customerId string) error {
	key, err := subscription.GetCustomerLock(customerId)
	if err != nil {
		return fmt.Errorf("failed to get customer lock: %w", err)
	}

	if err := s.Lockr.LockForTX(ctx, key); err != nil {
		return fmt.Errorf("failed to lock customer: %w", err)
	}

	return nil
}
