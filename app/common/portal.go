package common

import (
	"fmt"

	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/portal"
	portaladapter "github.com/Pototoooo/meterforge/meterforge/portal/adapter"
)

var Portal = wire.NewSet(
	NewPortalService,
)

func NewPortalService(conf config.PortalConfiguration) (portal.Service, error) {
	if !conf.Enabled {
		return portaladapter.NewNoop(), nil
	}

	p, err := portaladapter.New(portaladapter.Config{
		Secret: conf.TokenSecret,
		Expire: conf.TokenExpiration,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create portal adapter: %w", err)
	}

	return p, nil
}
