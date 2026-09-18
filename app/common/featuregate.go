package common

import (
	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/pkg/featuregate"
)

var FeatureGateChecker = wire.NewSet(
	featuregate.NewNoop,
	wire.FieldsOf(new(config.Configuration), "FeatureGate"),
	NewFeatureGateChecker,
)

func NewFeatureGateChecker(gate featuregate.Gate, config config.FeatureGateConfiguration, creditsConfig config.CreditsConfiguration) *featuregate.FeatureGateChecker {
	flags := make(featuregate.Flags)
	if config.Flags != nil {
		flags = config.Flags
	}

	flagOverrides := map[featuregate.FeatureFlag]bool{
		featuregate.CtxKeyCredits: creditsConfig.Enabled,
	}

	if !config.Enabled {
		return featuregate.NewFeatureGateChecker(featuregate.NewNoop(), flags, flagOverrides)
	}

	return featuregate.NewFeatureGateChecker(gate, flags, flagOverrides)
}
