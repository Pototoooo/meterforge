package common

import (
	"fmt"

	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/namespace"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
)

var Namespace = wire.NewSet(
	NewNamespaceManager,
)

func NewNamespaceManager(
	conf config.NamespaceConfiguration,
) (*namespace.Manager, error) {
	manager, err := namespace.NewManager(namespace.ManagerConfig{
		DefaultNamespace:  conf.Default,
		DisableManagement: conf.DisableManagement,
	})
	if err != nil {
		return nil, fmt.Errorf("create namespace manager: %v", err)
	}

	return manager, nil
}

var StaticNamespace = wire.NewSet(
	NewStaticNamespaceDecoder,
)

func NewStaticNamespaceDecoder(
	conf config.NamespaceConfiguration,
) namespacedriver.NamespaceDecoder {
	return namespacedriver.StaticNamespaceDecoder(conf.Default)
}
