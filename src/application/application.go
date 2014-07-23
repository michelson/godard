package application

import (
	cfg "godard_config"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	//"process"
	logger "godard_logger"
	ProcessJournal "process"
	socket "socket"
	"strconv"
	"strings"
	"syscall"
	"system"
	"time"
)

type Application struct {
	/*start     string
	  stop      string
	  restart   string
	  unmonitor string
	  status    string*/

	Foreground bool

	Name        string
	Logger      *logger.GodardLogger
	BaseDir     string
	PidFile     string
	KillTimeout int
	Groups      map[string]*Group
	WorkQueue   string
	PidsDir     string
	LogFile     string
	Sock        *socket.Socket
	running     bool
}

var Debug *log.Logger

func NewApplication(name string, options *cfg.GodardConfig) *Application {
	c := &Application{}
	c.Name = name

	c.Foreground = options.Foreground

	if len(options.LogFile) > 0 {
		c.LogFile = options.LogFile
	}

	if len(options.BaseDir) > 1 {
		c.BaseDir = options.BaseDir
	} else {
		c.BaseDir = os.Getenv("GODARD_BASE_DIR")
	}

	c.PidFile = path.Join(c.BaseDir, "pids", c.Name, c.Name+".pid") // File.join(self.base_dir, 'pids', self.name + ".pid")
	c.PidsDir = path.Join(c.BaseDir, "pids", c.Name)                //File.join(self.base_dir, 'pids', self.name)

	if options.KillTimeout > 0 {
		c.KillTimeout = options.KillTimeout
	}

	log.Println("PID FILE_:", c.PidFile)
	c.Groups = make(map[string]*Group, 0)

	logger_opts := make(map[string]interface{}, 0)
	logger_opts["log_file"] = c.LogFile
	logger_opts["stdout"] = c.isForeground()

	c.Logger = logger.NewGodardLogger(logger_opts).PrefixWith(c.Name)
	ProcessJournal.Logger = c.Logger.Logger
	Debug = ProcessJournal.Logger

	c.SetupSignalTraps()

	//@mutex = Mutex.new

	c.SetupPidsDir()

	return c
}

func (c *Application) isForeground() bool {
	return c.Foreground
}

func (c *Application) Start(names ...string) {
	c.sendToProcessOrGroup("start", names...)
}

func (c *Application) Stop(names ...string) {
	c.sendToProcessOrGroup("stop", names...)
}

func (c *Application) Restart(names ...string) {
	c.sendToProcessOrGroup("restart", names...)
}

func (c *Application) UnMonitor(names ...string) {
	c.sendToProcessOrGroup("unmonitor", names...)
}

func (c *Application) Status(names ...string) {
	c.sendToProcessOrGroup("status", names...)
}

func (c *Application) AddProcess(process *Process, group_name string) {

	var group *Group

	if len(c.Groups) == 0 {
		group = NewGroup(group_name, c.Logger.PrefixWith(group_name).Logger)
		c.Groups[group_name] = group
	} else {
		group = c.Groups[group_name]
	}

	process.Logger = group.Logger
	group.AddProcess(process)

	/*
	  Debug.Println("GROUPS COUNT: ", len(c.Groups) )
	  Debug.Println("GROUPS PROCESSES: ", c.Groups["group"].Processes )
	*/
}

func (c *Application) Load() {
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

func (c *Application) StartServer() {

	c.KillPreviousGodard()
	ProcessJournal.KillAllFromAllJournals()
	ProcessJournal.ClearAllAtomicFsLocks()

	// err := syscall.Setpgid(0, 0)

	//if err != nil {
	//  Debug.Println("Errno::EPERM", err)
	//}

	//Daemonize.daemonize unless foreground?
	//self.logger.reopen
	// $0 = "bluepilld: #{self.name}"

	for _, g := range c.Groups {
		g.DetermineInitialState()
	}

	for _, g := range c.Groups {
		//Debug.Println("GROUP: ", g,  k )
		g.Tick()
	}

	sock, err := socket.NewSocket(c.BaseDir, c.Name)

	if err != nil {
		Debug.Println(err)
	}
	c.WritePidFile()
	c.Sock = sock
	//c.Sock.ListenerChannel = make(chan string)

	c.SetupSignalTraps()

	go c.Sock.Run()

	go c.StartListener()

	c.Run()

}

func (c *Application) StartListener() {

	for {
		select {
		case msg := <-c.Sock.ListenerChannel:
			//Debug.Println("received message:", msg)
			args := strings.Split(msg, ":")
			c.sendToProcessOrGroup(args[0], args[1:]...)
		//case <-time.After(time.Second * 30):
		//    Debug.Println("timeout 1")
		default:

		}
	}
}

func (c *Application) Run() {
	c.running = true // set to false by signal trap

	for {
		//Debug.Println("APP RUNNING FOR:", c.running)
		if c.running {
			system.ResetData()
			for _, group := range c.Groups {
				group.Tick()
				time.Sleep(1 * time.Second)
			}
		}
	}
}

//Private

func (c *Application) SetupSignalTraps() {

	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(cc chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-cc
		Debug.Printf("Caught signal %s: shutting down.", sig)
		c.running = false
		c.CleanUp()
		// Stop listening (and unlink the socket if unix type):
		c.Sock.Listener.Close()
		//os.Remove("/tmp/godard.sock")
		/*
		   puts "Terminating..."
		   cleanup
		   @running = false
		*/
		// And we're done:
		os.Exit(0)
	}(sigc)
}

func (c *Application) SetupPidsDir() {
	/*FileUtils.mkdir_p(self.pids_dir) unless File.exists?(self.pids_dir)
	  # we need everybody to be able to write to the pids_dir as processes managed by
	  # bluepill will be writing to this dir after they've dropped privileges
	  FileUtils.chmod(0777, self.pids_dir)*/
	err := os.MkdirAll(c.PidsDir, 0777)
	if err != nil {
		Debug.Println("ERROR CREATING PIDS DIR", err)
	}
}

func (c *Application) KillPreviousGodard() {
	/*
	     if File.exists?(self.pid_file)
	       previous_pid = File.read(self.pid_file).to_i
	       if System.pid_alive?(previous_pid)
	         begin
	           ::Process.kill(0, previous_pid)
	           puts "Killing previous bluepilld[#{previous_pid}]"
	           ::Process.kill(2, previous_pid)
	         rescue Exception => e
	           $stderr.puts "Encountered error trying to kill previous bluepill:"
	           $stderr.puts "#{e.class}: #{e.message}"
	           exit(4) unless e.is_a?(Errno::ESRCH)
	         else
	           kill_timeout.times do |i|
	             sleep 0.5
	             break unless System.pid_alive?(previous_pid)
	           end

	           if System.pid_alive?(previous_pid)
	             $stderr.puts "Previous bluepilld[#{previous_pid}] didn't die"
	             exit(4)
	           end
	         end
	       end
	     end
	   end
	*/
}
func (c *Application) CleanUp() {
	/*
	   def cleanup
	     ProcessJournal.kill_all_from_all_journals
	     ProcessJournal.clear_all_atomic_fs_locks
	     begin
	       System.delete_if_exists(self.socket.path) if self.socket
	     rescue IOError
	     end
	     System.delete_if_exists(self.pid_file)
	   end
	*/
}

func (c *Application) sendToProcessOrGroup(method string, names ...string) {
	Debug.Println("PREPARE TO SEND", method, "TO PROC OR GROUP", names)
	var group_name string
	var process_name string
	group_name = names[0]
	if len(names) > 1 {
		process_name = names[1]
	}

	if len(group_name) == 0 && len(process_name) == 0 {

		for _, group := range c.Groups {
			Debug.Println("THIS GROUP IS READY TO ,", group)
			group.SendMethod(method, "")
		}

	} else if c.GroupInString(group_name) {
		c.Groups[group_name].SendMethod(method, process_name)

	} else if len(process_name) == 0 {
		// they must be targeting just by process name
		process_name = group_name
		for _, group := range c.Groups {
			//Debug.Println("THIS GROUP IS TARGETING JUST BY PROC ,", group)
			group.SendMethod(method, process_name)
		}
		/*
		   process_name = group_name
		   self.groups.values.collect do |group|
		     group.send(method, process_name)
		   end.flatten */
	} else {
		//[]
	}

	//Debug.Println(group_name , process_name)
}

func (c *Application) GroupInString(name string) bool {
	res := false
	if _, ok := c.Groups[name]; ok {
		res = true
	}
	return res
}

func (c *Application) WritePidFile() {
	//File.open(self.pid_file, 'w') { |x| x.write(::Process.pid) }
	str := []byte(strconv.Itoa(syscall.Getpid()))
	//Debug.Println("WRITTING APP PID:", string(str), c.PidFile)
	err := ioutil.WriteFile(c.PidFile, str, 0644)
	if err != nil {
		Debug.Println("Err creating pid:", err)
	}
}
