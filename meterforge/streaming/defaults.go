package streaming

import (
	"time"

	"github.com/Pototoooo/meterforge/meterforge/meter"
)

const (
	// MinimumWindowSizeDuration is the minimum window size the aggregation can represent.
	MinimumWindowSizeDuration = time.Second
	// MinimumWindowSize is the minimum window size the aggregation can represent.
	MinimumWindowSize = meter.WindowSizeSecond
)
