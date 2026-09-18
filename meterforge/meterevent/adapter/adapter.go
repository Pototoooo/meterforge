package adapter

import (
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/meterevent"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
)

func New(
	streamingConnector streaming.Connector,
	customerService customer.Service,
	meterService meter.Service,
) meterevent.Service {
	return &adapter{
		streamingConnector: streamingConnector,
		customerService:    customerService,
		meterService:       meterService,
	}
}

var _ meterevent.Service = (*adapter)(nil)

type adapter struct {
	streamingConnector streaming.Connector
	customerService    customer.Service
	meterService       meter.Service
}
