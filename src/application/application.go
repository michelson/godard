package application

import(
  pcs "process"
  cfg "godard_config"
  socket "socket"
  "log"
  "os"
  "path"
)

type Application struct {
  /*start     string 
  stop      string 
  restart   string 
  unmonitor string
  status    string*/

  Foreground bool

  Name        string
  Logger      string
  BaseDir     string
  PidFile     string
  KillTimeout string
  Groups      map[string]*Group
  WorkQueue   string
  PidsDir     string
  LogFile     string
  Sock        *socket.Socket
}


func NewApplication(name string , options *cfg.GodardConfig) *Application {
  c := &Application{}
  c.Name = name

  c.Foreground   = options.Foreground
 
  if len(options.LogFile) > 0 {
    c.LogFile      = options.LogFile
  }

  if len(options.BaseDir) == 0 {
    c.BaseDir      = os.Getenv("GODARD_BASE_DIR")
  }else{
    c.BaseDir      = options.BaseDir
  }
  //c.base_dir     = options["base_dir"] //|| ENV['BLUEPILL_BASE_DIR'] || (::Process.euid != 0 ? File.join(ENV['HOME'], '.bluepill') : "/var/run/bluepill")
  
  c.PidFile = path.Join(c.BaseDir, "pids", c.Name, ".pid") // File.join(self.base_dir, 'pids', self.name + ".pid")
  c.PidsDir = path.Join(c.BaseDir, "pids", c.Name) //File.join(self.base_dir, 'pids', self.name)
  //c.kill_timeout = options.KillTimeout || 10

  c.Groups = make(map[string]*Group, 0)

  /*
  self.logger = ProcessJournal.logger = Bluepill::Logger.new(:log_file => self.log_file, :stdout => foreground?).prefix_with(self.name)

  self.setup_signal_traps

  @mutex = Mutex.new
  */
  c.SetupPidsDir()

  return c
}


func (c *Application) isForeground() bool {
  return c.Foreground
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
  log.Println("ADDING PROCESS TO GROUP" )

  var group *Group 

  if len(c.Groups) == 0 {
    group = NewGroup(group_name) // :logger => self.logger.prefix_with(group_name))
    c.Groups[group_name] = group
  } else {
    group = c.Groups[group_name]
  }

  group.AddProcess(process)

  for k, _ := range(c.Groups) {
     log.Println("GROUP: ", k )
  }

  log.Println("GROUPS COUNT: ", len(c.Groups) )
  log.Println("GROUPS PROCESSES: ", c.Groups["group"].Processes )

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

  
}    

func (c*Application) StartServer(){
    
    //os.Remove("/tmp/godard.sock") // kill previous

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


    for _, g := range(c.Groups) {
      g.DetermineInitialState()
    }

    for k, g := range(c.Groups) {
      log.Println("GROUP: ", g,  k )
      g.Tick()
    }

    sock , err := socket.NewSocket()
    
    if err != nil {
      log.Println(err)
    }
    c.Sock = sock
    c.StartListener()
    c.Sock.Run()



}

//Private

func (c *Application) SetupPidsDir(){
      /*FileUtils.mkdir_p(self.pids_dir) unless File.exists?(self.pids_dir)
      # we need everybody to be able to write to the pids_dir as processes managed by
      # bluepill will be writing to this dir after they've dropped privileges
      FileUtils.chmod(0777, self.pids_dir)*/
      err := os.MkdirAll(c.PidsDir, 0777)
      if err != nil {
        log.Println("ERROR CREATING PIDS DIR" , err)
      }
}

func (c *Application) SendToProcessOrGroup(method string , names...string){
  
}

