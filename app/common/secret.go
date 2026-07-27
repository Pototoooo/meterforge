package common

import (
	"log/slog"

	"github.com/google/wire"

	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/secret"
	secretadapter "github.com/Pototoooo/meterforge/meterforge/secret/adapter"
	secretservice "github.com/Pototoooo/meterforge/meterforge/secret/service"
)

var Secret = wire.NewSet(
	wire.Bind(new(secret.Service), new(*secretservice.Service)),

	NewUnsafeSecretService,
)

func NewUnsafeSecretService(logger *slog.Logger, db *entdb.Client) (*secretservice.Service, error) {
	secretAdapter := secretadapter.New()

	return secretservice.New(secretservice.Config{
		Adapter: secretAdapter,
	})
}
