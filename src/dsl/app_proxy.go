package dsl

import (
 cfg "godard_config"
 "log"
 app "application"
 //process "system"
)
type AppProxy struct {
  WorkingDir  string
  Uid         string
  Gid         string
  Environment string
  AutoStart   string
  App         *app.Application
}

func NewAppProxy(app_name string , options *cfg.GodardConfig) *AppProxy {
   app := app.NewApplication(app_name , options)
   c := &AppProxy{}
   c.App = app
   for _, t := range(options.Processes) {
      c.AddProcesses(t)
   }
   return c
}

func (c *AppProxy) AddProcesses(t map[string]interface {} ) {
  log.Println("PROCESS CONFIG:", t["pid_file"] )


  //ATTRS := ["working_dir", "uid", "gid", "environment", "auto_start" ]

  //process_factory = ProcessFactory.new(attributes, process_block)

  //process = process_factory.create_process(process_name, @app.pids_dir)
  //group = process_factory.attributes.delete(:group)

  //group := t["group"] || ""

  //http://stackoverflow.com/questions/19021848/how-to-send-a-message-to-an-object-in-golang-send-equivalent-in-go

  process_factory := NewProcessFactory(t)
  process := process_factory.CreateProcess(t["name"].(string))
  //group = process_factory.attributes.delete(:group)

  c.App.AddProcess(process, "group")

}


