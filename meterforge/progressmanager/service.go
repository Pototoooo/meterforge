package progressmanager

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/progressmanager/entity"
)

type Service interface {
	ProgressManagerService
}

type ProgressManagerService interface {
	GetProgress(ctx context.Context, input entity.GetProgressInput) (*entity.Progress, error)
	UpsertProgress(ctx context.Context, input entity.UpsertProgressInput) error
}
