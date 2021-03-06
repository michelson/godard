package application

import (
	cfg "godard_config"
	"log"
)

type AppProxy struct {
	WorkingDir  string
	Uid         string
	Gid         string
	Environment string
	AutoStart   string
	App         *Application
}

func NewAppProxy(app_name string, options *cfg.GodardConfig) *AppProxy {
	app := NewApplication(app_name, options)
	c := &AppProxy{}
	c.App = app
	for _, t := range options.Processes {
		c.AddProcesses(t)
	}
	return c
}

func (c *AppProxy) AddProcesses(t map[string]interface{}) {
	log.Println("PROCESS CONFIG:", t["pid_file"])

	//ATTRS := ["working_dir", "uid", "gid", "environment", "auto_start" ]
	//http://stackoverflow.com/questions/19021848/how-to-send-a-message-to-an-object-in-golang-send-equivalent-in-go
	t["logger"] = c.App.Logger.Logger
	process_factory := NewProcessFactory(t)
	process := process_factory.CreateProcess(t["name"].(string), c.App.PidsDir)

	var group string
	if _, ok := t["group"]; ok {
		group = t["group"].(string)
	}

	c.App.AddProcess(process, group)
}
