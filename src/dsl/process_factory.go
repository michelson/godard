package dsl

import (
 //cfg "godard_config"
  //"log"
  proc "system"
)


type ProcessFactory struct {
  process_keys map[string]int
  pid_files    map[string]int
  attributes    map[string]interface{}
}

func NewProcessFactory(attributes map[string]interface{}) *ProcessFactory {
  c := &ProcessFactory{}
  c.attributes = attributes
  //c.process_block = process_block
  return c
}

func (c *ProcessFactory) CreateProcess(name string) *proc.Process {
  
  //log.Println("PROXY CREATING PROCESS:", name)

  process := NewProcessProxy(name, c.attributes)
  //child_process_block = @attributes.delete(:child_process_block)
  //self.validate_process! process
  c.ValidateProcess(process)
  p := process.ToProcess()
  
  return p

}

func (c *ProcessFactory) AssingDefaultPid() {}

func (c *ProcessFactory) ValidateProcess(process *ProcessProxy) {}

func (c *ProcessFactory) ValidateChildProcess(process *proc.Process) {}
