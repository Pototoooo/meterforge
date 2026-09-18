package adapter

import (
	"github.com/Pototoooo/meterforge/meterforge/secret"
)

func New() secret.Adapter {
	return &adapter{}
}

var _ secret.Adapter = (*adapter)(nil)

type adapter struct{}
