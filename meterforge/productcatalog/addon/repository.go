package addon

import (
	"context"

	"github.com/Pototoooo/meterforge/pkg/framework/entutils"
	"github.com/Pototoooo/meterforge/pkg/pagination"
)

// TODO: add bulk api

type Repository interface {
	entutils.TxCreator

	ListAddons(ctx context.Context, params ListAddonsInput) (pagination.Result[Addon], error)
	CreateAddon(ctx context.Context, params CreateAddonInput) (*Addon, error)
	DeleteAddon(ctx context.Context, params DeleteAddonInput) error
	GetAddon(ctx context.Context, params GetAddonInput) (*Addon, error)
	UpdateAddon(ctx context.Context, params UpdateAddonInput) (*Addon, error)
}
