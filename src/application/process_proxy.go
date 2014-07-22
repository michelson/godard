package application

import (
	"log"
	//app "application"
)

type ProcessProxy struct {
	Attributes map[string]interface{}
	Watches    map[string]interface{}
}

func NewProcessProxy(process_name string, attributes map[string]interface{}) *ProcessProxy {

	c := &ProcessProxy{}
	c.Attributes = attributes
	c.Watches = make(map[string]interface{}, 0)

	if _, ok := attributes["checks"]; ok {
		c.Checks(attributes)
	}

	if _, ok := attributes["monitor_children"]; ok {
		c.MonitorChildren(attributes)
	}
	return c
}

func (c*ProcessProxy) Checks(attributes map[string]interface{}) {
	m := attributes["checks"].(map[string]interface{})
		for k, v := range m {
			switch vv := v.(type) {
			case interface{}:
				log.Println(k, "is interface", vv)
				c.Watches[k] = vv
			default:
				log.Println(k, "is of a type I don't know how to handle")
			}
		}
}

func (c*ProcessProxy) MonitorChildren(child_process_block map[string]interface{}) {
  c.Attributes["monitor_children"] = true
  c.Attributes["child_process_block"] = child_process_block
}

func (c *ProcessProxy) ToProcess() *Process {
	/*p := &app.Process{}
	  p.Name = c.Attributes["name"].(string)
	  p.StartCommand = c.Attributes["start_command"].(string)
	  p.PidFile = c.Attributes["pid_file"].(string)
	  p.AddWatches(c.Watches)*/
	p := NewProcess(c.Attributes["name"].(string), c.Watches, c.Attributes)
	return p
}
