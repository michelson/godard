package application

import(
  pcs "process"
  "log"
  "strings"
  "reflect"
  "time"
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
  log.Println("Start Group")
}

func (c *Group) UnMonitor(process_name string) {
  log.Println("UnMonitor Group")
}

func (c *Group) Stop(process_name string) {
  log.Println("Stop Group")
}

func (c *Group) Restart(process_name string) {
  log.Println("Restart Group")
}

func (c*Group) SendMethod(method string , process_name string){
  log.Println("SEND",method,"METHOD TO" , process_name, " IS GOING TO BE SO COOL")
  log.Println(c.Processes)
  //threads = []

  var affected []string 
  //var threads  []map[string]int64
  //threads = make([]map[string]int64 ,0)
  for _ , process := range(c.Processes){
    if len(process_name) > 0 && process_name != process.Name{
      continue  
    }

    s := []string{c.Name , process.Name}
    affected = append(affected, strings.Join(s, ":") )
    v := reflect.ValueOf(*process)
    noblock_field := v.FieldByName("Group_"+method+"_noblock")
    noblock := noblock_field.Interface().(bool)
    if noblock {
      log.Println("Command", method ," running in non-blocking mode.")
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

      
    }else{
      log.Println("Command", method ," running in blocking mode.")
      //thread = Thread.new { process.handle_user_command("#{event}") }
      //thread.join
    }
  }

  // threads.each { |t| t.join } unless threads.nil?
  // affected 
  log.Println("SOME AFFECTED ARE:" , affected)
}

/*
      lines = []
      if process_name.nil?
        prefix = self.name ? "  " : ""
        lines << "#{self.name}:" if self.name

        self.processes.each do |process|
          lines << "%s%s(pid:%s): %s" % [prefix, process.name, process.actual_pid, process.state]
          if process.monitor_children?
            process.children.each do |child|
              lines << "  %s%s: %s" % [prefix, child.name, child.state]
            end
          end
        end
      else
        self.processes.each do |process|
          next if process_name != process.name
          lines << "%s%s(pid:%s): %s" % [prefix, process.name, process.actual_pid, process.state]
          lines << process.statistics.to_s
        end
      end
      lines << ""
*/  

