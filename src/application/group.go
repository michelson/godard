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
	Logger    *log.Logger
}

func NewGroup(name string, log_obj *log.Logger) *Group {
	c := &Group{}
	c.Processes = make([]*Process, 0)
	c.Logger = log_obj
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
	c.Logger.Println("SEND", method, "METHOD TO", process_name, " IS GOING TO BE SO COOL")
	//c.Logger.Println(c.Processes)

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
			c.Logger.Println("Command", method, " running in non-blocking mode.")
			//threads << Thread.new { process.handle_user_command("#{event}") }
			go process.HandleUserCommand(method)

			select {
			case msg := <-process.ListenerChannel:
				c.Logger.Println("PROCESS RECEIVED ACTION:", msg)
				//args := strings.Split(msg, ":")
				//threads = append(threads, msg)
			case <-time.After(time.Second * 2):
				c.Logger.Println("timeout 1")
			default:
			}

		} else {
			c.Logger.Println("Command", method, " running in blocking mode.")
			//thread = Thread.new { process.handle_user_command("#{event}") }
			//thread.join
		}
	}

	// threads.each { |t| t.join } unless threads.nil?
	// affected
	c.Logger.Println("SOME AFFECTED ARE:", affected)
}
