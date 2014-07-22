package godard

import (
	app "application"
	cfg "godard_config"
)

func Init(config *cfg.GodardConfig) {
	app.InitApplication("myApp", config)
}
