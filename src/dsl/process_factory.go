package dsl

import (
 //cfg "godard_config"
  "log"
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
  log.Println("CREATING PROCESS:", name)
  p := &proc.Process{}
  p.Name = c.attributes["name"].(string)
  p.StartCommand = c.attributes["start_command"].(string)
  p.PidFile = c.attributes["pid_file"].(string)

  //create child process here
  c.ValidateProcess(p)
  return p
}

func (c *ProcessFactory) AssingDefaultPid() {}

func (c *ProcessFactory) ValidateProcess(process *proc.Process) {}

func (c *ProcessFactory) ValidateChildProcess(process *proc.Process) {}
