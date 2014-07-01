package application

import(
  pcs "system"
)

type Group struct {
  name        string 
  processes   []*pcs.Process 
}


func NewGroup(name string , options []map[string]interface{}) *Group {
  c := &Group{}
  c.processes = make([]*pcs.Process , 0)
  return c
}


func (c *Group) AddProcess(process *pcs.Process) {
  c.processes = append(c.processes , process )
}

func (c *Group) Tick() {
  for _ , process := range(c.processes){
    process.Tick()
  }
}

func (c *Group) DetermineInitialState() {
  for _ , process := range(c.processes){
    process.DetermineInitialState()
  }
}


//[:start, :unmonitor, :stop, :restart]
func (c *Group) Start(process_name string) {
}

func (c *Group) UnMonitor(process_name string) {
}

func (c *Group) Stop(process_name string) {
}

func (c *Group) Restart(process_name string) {
}

/*
    def #{event}(process_name = nil)
      threads = []
      affected = []
      self.processes.each do |process|
        next if process_name && process_name != process.name
        affected << [self.name, process.name].join(":")
        noblock = process.group_#{event}_noblock
        if noblock
          self.logger.debug("Command #{event} running in non-blocking mode.")
          threads << Thread.new { process.handle_user_command("#{event}") }
        else
          self.logger.debug("Command #{event} running in blocking mode.")
          thread = Thread.new { process.handle_user_command("#{event}") }
          thread.join
        end
      end
      threads.each { |t| t.join } unless threads.nil?
      affected
    end
*/

func (c *Group) Status(process_name string) {
  
}
