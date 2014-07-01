package application

import(
  pcs "system"
  cfg "godard_config"
  socket "socket"
  "log"
)

type Application struct {
  /*start     string 
  stop      string 
  restart   string 
  unmonitor string
  status    string*/

  foreground bool

  name      string
  logger    string
  base_dir  string
  socket    string
  pid_file  string
  kill_timeout string
  groups    []*Group
  work_queue string
  pids_dir   string
  log_file   string
  Sock     *socket.Socket
}


func NewApplication(name string , options *cfg.GodardConfig) *Application {
  c := &Application{}
  c.name = name
  //c.foreground   = options["foreground"]
  //c.log_file     = options["log_file"]
  //c.base_dir     = options["base_dir"] //|| ENV['BLUEPILL_BASE_DIR'] || (::Process.euid != 0 ? File.join(ENV['HOME'], '.bluepill') : "/var/run/bluepill")

  /*c.pid_file = File.join(self.base_dir, 'pids', self.name + ".pid")
  c.pids_dir = File.join(self.base_dir, 'pids', self.name)
  c.kill_timeout = options[:kill_timeout] || 10*/

  c.groups = make([]*Group, 0)

  /*
  self.logger = ProcessJournal.logger = Bluepill::Logger.new(:log_file => self.log_file, :stdout => foreground?).prefix_with(self.name)

  self.setup_signal_traps
  self.setup_pids_dir

  @mutex = Mutex.new
  */

  return c
}


func (c *Application) isForeground() bool {
  return c.foreground
}

//s := []string{"James", "Jasmine"}
//Greeting("goodbye:", s...)

func (c *Application) Start(names...string)  {
  //group_name string, process_name string
  c.SendToProcessOrGroup("start", names...)
}

func (c *Application) Stop(names...string)  {
  //group_name string, process_name string
  c.SendToProcessOrGroup("stop", names...)
}

func (c *Application) Restart(names...string)  {
  //group_name string, process_name string
  c.SendToProcessOrGroup("restart", names...)
}

func (c *Application) UnMonitor(names...string)  {
  //group_name string, process_name string
  c.SendToProcessOrGroup("unmonitor", names...)
}

func (c *Application) Status(names...string)  {
  //group_name string, process_name string
  c.SendToProcessOrGroup("status", names...)
}

func (c *Application) AddProcess(process *pcs.Process, group_name string ){
  /*group_name = group_name.to_s if group_name

  self.groups[group_name] ||= Group.new(group_name, :logger => self.logger.prefix_with(group_name))
  self.groups[group_name].add_process(process)
  */

}

func (c*Application) Load(){
  c.StartServer()
  /*def load
      begin
        self.start_server
      rescue StandardError => e
        $stderr.puts "Failed to start bluepill:"
        $stderr.puts "%s `%s`" % [e.class.name, e.message]
        $stderr.puts e.backtrace
        exit(5)
      end
  end*/
}

func (c*Application) StartListener(){

  c.Sock.Run()
}    

func (c*Application) StartServer(){
    
    /*def start_server
      self.kill_previous_bluepill
      ProcessJournal.kill_all_from_all_journals
      ProcessJournal.clear_all_atomic_fs_locks

      begin
        ::Process.setpgid(0, 0)
      rescue Errno::EPERM
      end

      Daemonize.daemonize unless foreground?

      self.logger.reopen

      $0 = "bluepilld: #{self.name}"

      self.groups.each {|_, group| group.determine_initial_state }


      self.write_pid_file
      self.socket = Bluepill::Socket.server(self.base_dir, self.name)
      self.start_listener

      self.run
    end*/

    //ss = socket.NewSocket()
    //socket.NewSocket.server(self.base_dir, self.name)
    sock , err := socket.NewSocket()
    
    if err != nil {
      log.Println(err)
    }
    c.Sock = sock
    c.StartListener()

}
/*
    def add_process(process, group_name = nil)
      group_name = group_name.to_s if group_name

      self.groups[group_name] ||= Group.new(group_name, :logger => self.logger.prefix_with(group_name))
      self.groups[group_name].add_process(process)
    end
*/

//Private

func (c *Application) SendToProcessOrGroup(method string , names...string){
  
}

