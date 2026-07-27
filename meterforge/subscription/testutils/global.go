package subscriptiontestutils

import "github.com/Pototoooo/meterforge/pkg/datetime"

var ExampleNamespace = "test-namespace"

var ISOMonth, _ = datetime.ISODurationString("P1M").Parse()
