package flushhandler

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/sink/models"
)

type FlushEventHandler interface {
	OnFlushSuccess(ctx context.Context, events []models.SinkMessage) error
	Start(context.Context) error
	WaitForDrain(context.Context) error
	Close() error
}

type (
	FlushCallback func(context.Context, []models.SinkMessage) error
)
