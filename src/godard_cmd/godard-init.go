package godard

import (
	app "dsl"
	cfg "godard_config"
)

func Init(config *cfg.GodardConfig) {
	app.InitApplication("myApp", config)
}
