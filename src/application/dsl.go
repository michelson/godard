package application

import (
	cfg "godard_config"
)

func InitApplication(app_name string, options *cfg.GodardConfig) {
	app_proxy := NewAppProxy(app_name, options)
	app_proxy.App.Load()
}
