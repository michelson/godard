package application

import (
	"log"
	//pcs "process"
	"fmt"
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
	//c.Logger.Println(c.Processes)
	actions := []string{"start", "unmonitor", "stop", "restart"}
	var affected []string

	if stringInSlice(method, actions) {

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

	} else if method == "status" {
		c.Status(process_name)
	}

	c.Logger.Println("SOME AFFECTED ARE:", affected)
}

func (c *Group) Status(process_name string) string {
	c.Logger.Println("OLA PROCESS NAME STATUS!!!!")

	var lines []string
	var prefix string
	if process_name == "" {
		prefix = " "
		if c.Name != "" {
			lines = append(lines, c.Name)
		}
		for _, process := range c.Processes {
			str := fmt.Sprintf("%s%s(pid:%s): %s", prefix, process.Name, process.ActualPid, process.state)
			lines = append(lines, str)
			if process.MonitorChildren {
				for _, child := range process.Children {
					child_str := fmt.Sprintf("  %s%s: %s", prefix, child.Name, child.state)
					lines = append(lines, child_str)
				}
			}
		}
	} else {
		for _, process := range c.Processes {
			if process_name != process.Name {
				continue
			}
			out1 := fmt.Sprintf("%s%s(pid:%d): %s", prefix, process.Name, int(process.actual_pid), process.state)
			lines = append(lines, out1)
			lines = append(lines, process.Statistics.ToS())
		}

	}

	lines = append(lines, " ")
	c.Logger.Println(strings.Join(lines, " "))
	return strings.Join(lines, " ")
}

//dry this to utils
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
