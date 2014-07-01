package application

import (
 cfg "godard_config"
)
type AppProxy struct {
  WorkingDir  string
  Uid         string
  Gid         string
  Environment string
  AutoStart   string
  App         *Application
}

func NewAppProxy(app_name string , options *cfg.GodardConfig) *AppProxy {
   app := NewApplication(app_name , options)
   c := &AppProxy{}
   c.App = app
   return c
}

