package dsl

import(
    proc "system"
    "log" 
)

type ProcessProxy struct {
  Attributes map[string]interface{}
  Watches    map[string]interface{}
}

func NewProcessProxy(process_name string, attributes map[string]interface{}) *ProcessProxy {

  c := &ProcessProxy{}
  c.Attributes = attributes
  c.Watches = make(map[string]interface{}, 0)
  return c
}


func (c*ProcessProxy) checks(name string, options interface{}) {
  c.Watches[name] = options
}

/*
func (c*ProcessProxy) monitor_children(&child_process_block) {
      @attributes[:monitor_children] = true
      @attributes[:child_process_block] = child_process_block
}
*/

func (c*ProcessProxy) ToProcess() *proc.Process {
  p := &proc.Process{}
  p.Name = c.Attributes["name"].(string)
  p.StartCommand = c.Attributes["start_command"].(string)
  p.PidFile = c.Attributes["pid_file"].(string)
  log.Println("CREATING PROCESS:", p.Name)
  return p
}

