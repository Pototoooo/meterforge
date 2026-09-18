package common

import (
	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/meterevent"
	"github.com/Pototoooo/meterforge/meterforge/meterevent/adapter"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
)

var MeterEvent = wire.NewSet(
	NewMeterEventService,
)

func NewMeterEventService(
	streamingConnector streaming.Connector,
	customerService customer.Service,
	meterService meter.Service,
) meterevent.Service {
	return adapter.New(streamingConnector, customerService, meterService)
}
