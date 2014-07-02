package dsl

import(
    proc "process"
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

  if _,ok := attributes["checks"]; ok {
    m := attributes["checks"].(map[string]interface{})
    for k, v := range m {
      switch vv := v.(type) {
      case interface{}:
        //log.Println(k, "is interface", vv)
        c.Checks(k , vv)
      default:
        log.Println(k, "is of a type I don't know how to handle")
      }
    }
  }

  //log.Println("CHECKS:", c.Watches)

  return c
}


func (c*ProcessProxy) Checks(name string, options interface{}) {
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
  p.AddWatches(c.Watches)
  log.Println("CREATING PROCESS:", p.Name)
  return p
}

