package dsl

import (
	"log"
	"path"
	proc "process"
	"regexp"
	"strings"
)

type ProcessFactory struct {
	process_keys map[string]int
	pid_files    map[string]int
	attributes   map[string]interface{}
}

func NewProcessFactory(attributes map[string]interface{}) *ProcessFactory {
	c := &ProcessFactory{}
	c.attributes = attributes
	//c.process_block = process_block
	return c
}

func (c *ProcessFactory) CreateProcess(name string, pids_dir string) *proc.Process {

	c.assignDefaultPidFile(name, pids_dir)

	//log.Println("PROXY CREATING PROCESS:", name)
	process := NewProcessProxy(name, c.attributes)
	//child_process_block = @attributes.delete(:child_process_block)
	//self.validate_process! process
	c.ValidateProcess(process)
	p := process.ToProcess()

	return p

}

func (c *ProcessFactory) assignDefaultPidFile(process_name string, pids_dir string) {
	if _, ok := c.attributes["pid_file"]; ok {

	} else {
		group_name := c.attributes["group"].(string)
		s := []string{group_name, process_name}
		d := strings.Join(s, "_")

		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		default_pid_name := reg.ReplaceAllString(d, "_")

		c.attributes["pid_file"] = path.Join(pids_dir, default_pid_name+".pid")
	}
}

func (c *ProcessFactory) ValidateProcess(process *ProcessProxy) {
	//TODO
}

func (c *ProcessFactory) ValidateChildProcess(process *proc.Process) {
	//TODO
}
