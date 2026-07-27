package internal

import "github.com/Pototoooo/meterforge/app/config"

//nolint:gochecknoglobals
var (
	App         Application
	AppShutdown func()
	Config      config.Configuration
	ConfigFile  string
)
