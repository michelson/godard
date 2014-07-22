package application

import (
	"log"
	//pcs "process"
	"reflect"
	"strings"
	"time"
)

type Group struct {
	Name      string
	Processes []*Process
	Options   []map[string]interface{}
}

func NewGroup(name string) *Group {
	c := &Group{}
	c.Processes = make([]*Process, 0)
	return c
}

func (c *Group) AddProcess(process *Process) {
	c.Processes = append(c.Processes, process)
}

func (c *Group) Tick() {
	for _, process := range c.Processes {
		process.Tick()
	}
}

func (c *Group) DetermineInitialState() {
	for _, process := range c.Processes {
		process.DetermineInitialState()
	}
}

func (c *Group) SendMethod(method string, process_name string) {
	log.Println("SEND", method, "METHOD TO", process_name, " IS GOING TO BE SO COOL")
	//log.Println(c.Processes)

	var affected []string
	for _, process := range c.Processes {
		if len(process_name) > 0 && process_name != process.Name {
			continue
		}

		s := []string{c.Name, process.Name}
		affected = append(affected, strings.Join(s, ":"))
		v := reflect.ValueOf(*process)
		noblock_field := v.FieldByName("Group_" + method + "_noblock")
		noblock := noblock_field.Interface().(bool)
		if noblock {
			log.Println("Command", method, " running in non-blocking mode.")
			//threads << Thread.new { process.handle_user_command("#{event}") }
			go process.HandleUserCommand(method)

			select {
			case msg := <-process.ListenerChannel:
				log.Println("PROCESS RECEIVED ACTION:", msg)
				//args := strings.Split(msg, ":")
				//threads = append(threads, msg)
			case <-time.After(time.Second * 2):
				log.Println("timeout 1")
			default:
			}

		} else {
			log.Println("Command", method, " running in blocking mode.")
			//thread = Thread.new { process.handle_user_command("#{event}") }
			//thread.join
		}
	}

	// threads.each { |t| t.join } unless threads.nil?
	// affected
	log.Println("SOME AFFECTED ARE:", affected)
}
