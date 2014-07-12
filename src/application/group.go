package application

import(
  pcs "process"
  "log"
  "strings"
  "reflect"
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
  for _ , process := range(c.Processes){
    if len(process_name) > 0 && process_name != process.Name{
      continue  
    }

    s := []string{c.Name , process.Name}
    affected = append(affected, strings.Join(s, ":") )
    v := reflect.ValueOf(*process)
    noblock := v.FieldByName("Group_"+method+"_noblock")
    
    if noblock {
      //noblock.Interface().(bool) // reflection method value
      log.Println("Command", method ," running in non-blocking mode.")
      //threads << Thread.new { process.handle_user_command("#{event}") }
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

func (c *Group) Status(process_name string) {
  
}
