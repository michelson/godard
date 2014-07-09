package application

import(
  pcs "process"
  "log"
)

type Group struct {
  Name        string 
  Processes   []*pcs.Process 
  Options     []map[string]interface{}
}


func NewGroup(name string) *Group {
  c := &Group{}
  c.Processes = make([]*pcs.Process , 0)
  return c
}


func (c *Group) AddProcess(process *pcs.Process) {
  c.Processes = append(c.Processes , process )
}

func (c *Group) Tick() {
  for _ , process := range(c.Processes){
    log.Println("tick process")
    process.Tick()
  }
}

func (c *Group) DetermineInitialState() {
  for _ , process := range(c.Processes){
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
