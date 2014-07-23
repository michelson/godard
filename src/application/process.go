package application

import (
	//"io"
	"io/ioutil"
	"log"
	watcher "watcher"
	//trigger "trigger"
	"errors"
	fsm "github.com/looplab/fsm"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	system "system"
	time "time"
	"util"
	//"sync/atomic"
	proc "process"
)

var wg sync.WaitGroup

type Process struct {
	Name string

	pid     int
	ppid    int
	cpu     float64
	mem     int
	elapsed int
	command string

	Watches    []*watcher.ConditionWatch
	Triggers   []*Trigger
	Children   []*Process
	Statistics *proc.ProcessStatistics

	pid_file          string
	pre_start_command string
	StartCommand      string
	StopCommand       string
	RestartCommand    string

	CacheActualPid      bool
	MonitorChildren     bool
	ChildProcessFactory ProcessFactory

	Stdout string
	Stderr string
	Stdin  string

	Daemonize   bool
	PidFile     string
	WorkingDir  string
	Environment map[string]string

	StartGraceTime   time.Duration
	StopGraceTime    time.Duration
	RestartGraceTime time.Duration

	PreStartCommand string

	Uid string
	Gid string

	actual_pid int64

	PidCommand string
	AutoStart  bool

	SupplementaryGroups []string

	StopSignals []string

	OnStartTimeout string

	Group_start_noblock     bool
	Group_restart_noblock   bool
	Group_stop_noblock      bool
	Group_unmonitor_noblock bool

	skip_ticks_until int64
	process_running  bool
	state            string

	event_mutex *sync.Mutex

	Logger *log.Logger

	state_machine *fsm.FSM

	ListenerChannel chan map[string]string

	Transitioned bool
}

func (c *Process) SetConfigOptions(options map[string]interface{}) {
	//use reflect to dry this
	/*for k, v := range(options){
	  if _,ok := c.Groups[k]; ok {
	     res = true
	  }
	}*/

	if _, ok := options["start_command"]; ok {
		c.StartCommand = options["start_command"].(string)
	}

	if _, ok := options["pre_start_command"]; ok {
		c.PreStartCommand = options["pre_start_command"].(string)
	}

	if _, ok := options["stop_command"]; ok {
		c.StopCommand = options["stop_command"].(string)
	}

	if _, ok := options["restart_command"]; ok {
		c.RestartCommand = options["restart_command"].(string)
	}

	if _, ok := options["stdout"]; ok {
		c.Stdout = options["stdout"].(string)
	}

	if _, ok := options["stderr"]; ok {
		c.Stderr = options["stderr"].(string)
	}

	if _, ok := options["stdin"]; ok {
		c.Stdin = options["stdin"].(string)
	}

	if _, ok := options["pid_file"]; ok {
		c.PidFile = options["pid_file"].(string)
	}

	if _, ok := options["working_dir"]; ok {
		c.WorkingDir = options["working_dir"].(string)
	}

	if _, ok := options["daemonize"]; ok {
		c.Daemonize = options["daemonize"].(bool)
	}

	//if _,ok := options["environment"]; ok {
	//  c.Environment = options["environment"].(bool)
	//}

	if _, ok := options["auto_start"]; ok {
		c.AutoStart = options["auto_start"].(bool)
	}

	if _, ok := options["pid_command"]; ok {
		c.PidCommand = options["pid_command"].(string)
	}

	if _, ok := options["start_grace_time"]; ok {
		c.StartGraceTime, _ = util.TimeParse(options["start_grace_time"].(string))
	}

	if _, ok := options["stop_grace_time"]; ok {
		c.StopGraceTime, _ = util.TimeParse(options["stop_grace_time"].(string))
	}

	if _, ok := options["restart_grace_time"]; ok {
		c.RestartGraceTime, _ = util.TimeParse(options["restart_grace_time"].(string))
	}

	if _, ok := options["on_start_timeout"]; ok {
		c.OnStartTimeout = options["on_start_timeout"].(string)
	}

	if _, ok := options["gid"]; ok {
		c.Gid = options["gid"].(string)
	}

	if _, ok := options["uid"]; ok {
		c.Uid = options["uid"].(string)
	}

	if _, ok := options["cache_actual_pid"]; ok {
		c.CacheActualPid = options["cache_actual_pid"].(bool)
	}

	if _, ok := options["monitor_children"]; ok {
		c.MonitorChildren = options["monitor_children"].(bool)
	}

	if _, ok := options["pid_command"]; ok {
		c.PidCommand = options["pid_command"].(string)
	}

	if _, ok := options["supplementary_groups"]; ok {
		c.SupplementaryGroups = options["supplementary_groups"].([]string)
	}

	if _, ok := options["stop_signals"]; ok {
		c.StopSignals = options["stop_signals"].([]string)
	}

	if _, ok := options["group_start_noblock"]; ok {
		c.Group_start_noblock = options["group_start_noblock"].(bool)
	}

	if _, ok := options["group_restart_noblock"]; ok {
		c.Group_restart_noblock = options["group_restart_noblock"].(bool)
	}

	if _, ok := options["group_stop_noblock"]; ok {
		c.Group_stop_noblock = options["group_stop_noblock"].(bool)
	}

	if _, ok := options["group_unmonitor_noblock"]; ok {
		c.Group_unmonitor_noblock = options["group_unmonitor_noblock"].(bool)
	}

}

func NewProcess(process_name string, checks map[string]interface{}, options map[string]interface{}) *Process {
	c := &Process{}
	c.Name = process_name
	//p.AddWatches(c.Watches)
	c.event_mutex = &sync.Mutex{}
	c.Watches = make([]*watcher.ConditionWatch, 0)
	c.Triggers = make([]*Trigger, 0)
	c.Children = make([]*Process, 0)
	// @threads = []
	c.Statistics = proc.NewProcessStatistics()
	// @actual_pid = options[:actual_pid]
	c.Logger = options["logger"].(*log.Logger)

	c.ListenerChannel = make(chan map[string]string)

	for check_name, value := range checks {
		trigger_exists := false
		for _, v := range c.Triggers {
			if v.Name == check_name {
				trigger_exists = true
				break
			}
		}
		if trigger_exists {
			log.Println("ADD TRIGGER:", check_name, value)
			c.AddTrigger(check_name, value)
		} else {
			//c.Logger.Println("add watch here:", check_name, value)
			c.AddWatch(check_name, value)
		}

	}

	// These defaults are overriden below if it's configured to be something else.
	c.MonitorChildren = false
	c.CacheActualPid = true
	c.StartGraceTime = time.Second * 3
	c.StopGraceTime = time.Second * 3
	c.RestartGraceTime = time.Second * 3
	//@environment = {}
	c.OnStartTimeout = "start"
	c.Group_start_noblock = true
	c.Group_stop_noblock = true
	c.Group_restart_noblock = true
	c.Group_unmonitor_noblock = true

	c.AutoStart = true

	c.SetConfigOptions(options)

	/*
		  CONFIGURABLE_ATTRIBUTES.each do |attribute_name|
			self.send("#{attribute_name}=", options[attribute_name]) if options.has_key?(attribute_name)
		  end

		  # Let state_machine do its initialization stuff
		  super() # no arguments intentional
	*/

	// https://github.com/looplab/fsm/blob/master/fsm_test.go
	c.state_machine = fsm.NewFSM(
		"unmonitored",
		fsm.Events{

			{Name: "tick_up", Src: []string{"starting"}, Dst: "up"},
			{Name: "tick_down", Src: []string{"starting"}, Dst: "down"},

			{Name: "tick_up", Src: []string{"up"}, Dst: "up"},
			{Name: "tick_down", Src: []string{"up"}, Dst: "down"},

			{Name: "tick_up", Src: []string{"stopping"}, Dst: "up"},
			{Name: "tick_down", Src: []string{"stopping"}, Dst: "down"},

			{Name: "tick_up", Src: []string{"down"}, Dst: "up"},
			{Name: "tick_down", Src: []string{"down"}, Dst: "starting"},

			{Name: "tick_up", Src: []string{"restarting"}, Dst: "up"},
			{Name: "tick_down", Src: []string{"restarting"}, Dst: "down"},

			{Name: "start", Src: []string{"unmonitored", "down"}, Dst: "starting"},

			{Name: "restart", Src: []string{"up", "down"}, Dst: "restarting"},
		},
		fsm.Callbacks{
			"before_event": func(e *fsm.Event) {
				c.Logger.Println("EXEC STATE CHANGE FROM", c.state_machine.Current())
				c.NotifyTriggers(c.state_machine.Current())
				if !c.state_machine.Is("stopping") {
					c.CleanThreads()
				}
			},
			"after_event": func(e *fsm.Event) {
				c.Logger.Println("EXEC STATE CHANGE TO", c.state_machine.Current())
				if c.state_machine.Is("starting") {
					c.StartProcess()
				}

				if c.state_machine.Is("stopping") {
					c.StopProcess()
				}

				if c.state_machine.Is("restarting") {
					c.RestartProcess()
				}

				c.RecordTransition(c.state_machine.Current())
			},
		},
	)

	//c.Logger.Println("CREATING PROCESS:", c.Name)

	return c
}

func (c *Process) Tick() {

	if c.isSkippingTicks() {
		//c.Logger.Println("SKIPPING TICKS")
	} else {
		//c.skip_ticks_until = nil
		c.process_running = false

		// Deal with thread cleanup here since the stopping state isn't used
		//clean_threads if self.unmonitored?
		if c.state_machine.Is("unmonitored") {
			c.CleanThreads()
		}

		// run state machine transitions
		if c.isProcessRunning(false) {
			//c.Logger.Println("TICKS UP WITH", c.state_machine.Current())
			c.state_machine.Event("tick_up")
		} else {
			//c.Logger.Println("TICKS DOWN WITH CURRENT", c.state_machine.Current())
			c.state_machine.Event("tick_down")
		}

		//c.Logger.Println("CURRENT STATE:", c.state_machine.Current())

		if c.isUp() {
			c.RunWatches()
			if c.MonitorChildren {
				c.RefreshChildren()
				for _, child := range c.Children {
					child.Tick()
				}
			}
		}
	}
}

/*
   def logger=(logger)
     @logger = logger
     self.watches.each {|w| w.logger = logger }
     self.triggers.each {|t| t.logger = logger }
   end
*/

func (c *Process) isUp() bool {
	return c.state_machine.Current() == "up"
}

func (c *Process) Dispatch(event string, reason string) {
	c.event_mutex.Lock()
	c.Statistics.RecordEvent(event, reason)
	c.state_machine.Event(event)
	c.Logger.Println("STATS: ", c.Statistics.Events)
	c.event_mutex.Unlock()
}

func (c *Process) RecordTransition(transition string) {
	/*
	   unless transition.loopback?
	     @transitioned = true

	     # When a process changes state, we should clear the memory of all the watches
	     self.watches.each { |w| w.clear_history! }

	     # Also, when a process changes state, we should re-populate its child list
	     if self.monitor_children?
	       self.logger.warning "Clearing child list"
	       self.children.clear
	     end
	     logger.info "Going from #{transition.from_name} => #{transition.to_name}"
	   end

	*/

	c.Transitioned = true
	for _, watch := range c.Watches {
		watch.ClearHistory()
	}
	/*
	   # Also, when a process changes state, we should re-populate its child list
	   if self.monitor_children?
	     self.logger.warning "Clearing child list"
	     self.children.clear
	   end
	*/

	c.Logger.Println("TRANSITION TO: ", c.state_machine.Current())

}

func (c *Process) NotifyTriggers(transition string) {
	// self.triggers.each {|trigger| trigger.notify(transition)}
	for _, tgr := range c.Triggers {
		tgr.Notify(transition)
	}
}

func (c *Process) AddTrigger(name string, value interface{}) {
	//   self.triggers << Trigger[name].new(self, options.merge(:logger => self.logger))
	v := value.(map[string]interface{})
	//m["name"] = name
	v["logger"] = c.Logger
	c.Triggers = append(c.Triggers, NewTrigger(c, v))
	c.Logger.Println("TRIGGER ADDED:", c.Triggers)
}

func (c *Process) AddWatches(options map[string]interface{}) {

	if len(options) > 0 {
		c.Logger.Println("ADDING WATCHES TO PROCESS:", len(options))
		for k, v := range options {
			c.AddWatch(k, v)
		}
	}
}

func (c *Process) AddWatch(name string, value interface{}) {
	//c.Logger.Println("CHECKS:", name, value)

	v := value.(map[string]interface{})
	v["logger"] = c.Logger
	//c.Logger.Println(v["every"])
	condition := watcher.NewConditionWatch(name, v)
	c.Watches = append(c.Watches, condition)
}

type WatcherResponder struct {
	Watcher  *watcher.ConditionWatch
	Response []string
}

func (c *Process) RunWatches() {

	now := float64(time.Now().Unix())
	threads := make([]*WatcherResponder, 0)
	//c.Logger.Println("RUN WATCHES", c.Watches)
	for _, watch := range c.Watches {
		pid := c.ActualPid()
		wr := &WatcherResponder{Watcher: watch, Response: watch.Run(pid, now)}
		//c.Logger.Println("WATCH RES ON PID:", pid ,  wr.Watcher.Name , "VAL:", wr.Response )
		threads = append(threads, wr)
	}

	c.Transitioned = false
	for _, thread := range threads {
		if len(thread.Response) > 0 {
			c.Logger.Println(thread.Watcher.Name, " dispatched: ", thread.Response)
			for _, event := range thread.Response {
				if c.Transitioned {
					break
				}
				c.Dispatch(event, thread.Watcher.ToS())
			}

		}
	}

	//c.Logger.Println("RUN WATCHES NOW:", c.state_machine.Current())
}

func (c *Process) DetermineInitialState() {

	if c.isProcessRunning(true) {
		c.Logger.Println("IS RUNNING. SET UP STATUS")
		c.state_machine.SetCurrent("up")
	} else {
		//(auto_start == false) ? 'unmonitored' : 'down' # we need to check for false value
		c.Logger.Println("ISN'T RUNNING, SET DOWN STATUS")
		if c.AutoStart == false {
			c.state_machine.SetCurrent("unmonitored")
		} else {
			c.state_machine.SetCurrent("down")
		}
	}

	c.Logger.Println("DETERMINE INITAL STATE", c.state_machine.Current())

}

// System Process Methods

func (c *Process) isProcessRunning(force bool) bool {

	if force {
		c.process_running = false
	}

	if c.ActualPid() != 0 {
		if !c.process_running {
			process, err := os.FindProcess(c.ActualPid())
			if err != nil {
				//log.Printf("Failed to find process: %s\n", err)
			} else {
				err := process.Signal(syscall.Signal(0))
				//log.Printf("process.Signal on pid %d returned: %v\n", c.ActualPid(), err)
				if err == nil {
					c.process_running = true
				}
			}
		}
	}

	// the process isn't running, so we should clear the PID
	if !c.process_running {
		c.ClearPid()
	}
	//c.Logger.Println("PROCESS IS RUNNING?", c.process_running)
	return c.process_running
}

func (c *Process) HandleUserCommand(cmd string) {
	switch cmd {
	case "start":
		if c.isProcessRunning(true) {
			c.Logger.Println("Refusing to re-run start command on an already running process.")
		} else {
			c.Dispatch("start", "user initiated")
		}
	case "stop":
		c.StopProcess()
		c.Dispatch("unmonitor", "user unmonitor")
	case "restart":
		c.RestartProcess()
	case "unmonitor":
		// When the user issues an unmonitor cmd, reset any triggers so that
		// scheduled events gets cleared
		for _, trgr := range c.Triggers {
			trgr.Reset()
		}
		c.Dispatch("unmonitor", "user initiated")
	default:
		c.Logger.Println("default")
	}
}

func (c *Process) StartProcess() {
	proc.KillAllFromJournal(c.Name)
	c.PreStartProcess()
	if c.isDaemonized() {
		c.Logger.Println("Executing start cmd DEMONIZED:", c.StartCommand)
		/* daemon_id = System.daemonize(start_command, self.system_command_options)
		   if daemon_id
		     ProcessJournal.append_pid_to_journal(name, daemon_id)
		     children.each {|child|
		       ProcessJournal.append_pid_to_journal(name, child.actual_id)
		     } if self.monitor_children?
		   end
		   daemon_id
		*/
	} else {
		/*
		   # This is a self-daemonizing process
		   with_timeout(start_grace_time, on_start_timeout) do
		     result = System.execute_blocking(start_command, self.system_command_options)

		     unless result[:exit_code].zero?
		       logger.warning "Start command execution returned non-zero exit code:"
		       logger.warning result.inspect
		     end
		   end
		*/
		c.Logger.Println("Executing start cmd SELF-DEMONIZED:", c.StartCommand)

		result := system.ExecuteBlocking(c.StartCommand, c.SystemCommandOptions())
		c.Logger.Println("EXEC RESULT :", result)
		//c.ListenerChannel <- result

		//result = System.execute_blocking(start_command, self.system_command_options)

		//unless result[:exit_code].zero?
		//  logger.warning "Start command execution returned non-zero exit code:"
		//  logger.warning result.inspect
		//end
	}

	c.SkipTicksFor(c.StartGraceTime.Seconds())
}

func (c *Process) PreStartProcess() {
	if c.pre_start_command != "" {
		c.Logger.Println("Executing pre start command:", c.pre_start_command)
		result := system.ExecuteBlocking(c.pre_start_command, c.SystemCommandOptions())
		//c.Logger.Println("PRE START COMMAND RESULT :", result)
		if result["exit_code"] != "0" {
			c.Logger.Println("Pre start command execution returned non-zero exit code:")
			c.Logger.Println(result)
		}
	}
}

//NOK
func (c *Process) StopProcess() {
	if c.MonitorChildren {
		childs, _ := system.GetChildren(c.actual_pid)
		for _, child_pid := range childs {
			proc.AppendPidToJournal(c.Name, child_pid["pid"].(int))
			c.Logger.Println("Stop process : ", child_pid)
		}
	}
	if len(c.StopCommand) > 0 {
		cmd := c.PrepareCommand(c.StopCommand)
		c.Logger.Println("Executing stop command:", cmd)

		result := system.ExecuteBlocking(cmd, c.SystemCommandOptions())
		//c.Logger.Println("EXEC RESULT:", result)
		c.ListenerChannel <- result
		/*
		   with_timeout(stop_grace_time, "stop") do
		     result = System.execute_blocking(cmd, self.system_command_options)

		     unless result[:exit_code].zero?
		       logger.warning "Stop command execution returned non-zero exit code:"
		       logger.warning result.inspect
		     end
		   end
		*/
	} else if len(c.StopSignals) > 0 {
		//issue stop signals with configurable delay between each
		c.Logger.Println("Sending stop signals to", c.actual_pid)
		/*
		   @threads << Thread.new(self, stop_signals.clone) do |process, stop_signals|
		     signal = stop_signals.shift
		     logger.info "Sending signal #{signal} to #{process.actual_pid}"
		     process.signal_process(signal) # send first signal

		     until stop_signals.empty?
		       # we already checked to make sure stop_signals had an odd number of items
		       delay = stop_signals.shift
		       signal = stop_signals.shift

		       logger.debug "Sleeping for #{delay} seconds"
		       sleep delay
		       #break unless signal_process(0) #break unless the process can be reached
		       unless process.signal_process(0)
		         logger.debug "Process has terminated."
		         break
		       end
		       logger.info "Sending signal #{signal} to #{process.actual_pid}"
		       process.signal_process(signal)
		     end
		   end
		*/
	} else {
		c.Logger.Println("Executing default stop command. Sending TERM signal to", c.ActualPid)
		c.SignalProcess(syscall.SIGTERM)
	}
	proc.KillAllFromJournal(c.Name) // finish cleanup
	c.UnlinkPid()                   // TODO: we only write the pid file if we daemonize, should we only unlink it if we daemonize?

	c.SkipTicksFor(c.StopGraceTime.Seconds())

}

func (c *Process) RestartProcess() {

	if c.RestartCommand != "" {
		cmd := c.PrepareCommand(c.RestartCommand)
		c.Logger.Println("Executing restart command:", cmd)
		//MAKE FUNCTIONAL HERE!!!
		/*
		   with_timeout(restart_grace_time, "restart") do
		     result = System.execute_blocking(cmd, self.system_command_options)

		     unless result[:exit_code].zero?
		       logger.warning "Restart command execution returned non-zero exit code:"
		       logger.warning result.inspect
		     end
		   end
		*/

		result := system.ExecuteBlocking(cmd, c.SystemCommandOptions())
		c.Logger.Println("EXEC RESULT:", result)
		c.ListenerChannel <- result

		c.SkipTicksFor(c.RestartGraceTime.Seconds())

	} else {
		c.Logger.Println("No RestartCommand specified. Must stop and start to restart")

		wg.Add(1)

		go func() {
			c.StopProcess()
			c.StartProcess()
			defer wg.Done()
		}()

		wg.Wait()
	}

}

func (c *Process) CleanThreads() {
	//@threads.each { |t| t.kill }
	//@threads.clear
}

func (c *Process) isDaemonized() bool {
	return !!c.Daemonize
}

func (c *Process) isMonitorChildren() bool {
	return !!c.MonitorChildren
}

func (c *Process) SignalProcess(code syscall.Signal) bool {
	/*
	  HUP (hang up)
	  INT (interrupt)
	  QUIT (quit)
	  ABRT (abort)
	  KILL (non-catchable, non-ignorable kill)
	  ALRM (alarm clock)
	  TERM (software termination signal)
	*/
	c.Logger.Println("WE ARE GOING TO KILL PROCESS PID:", c.ActualPid())
	if c.actual_pid == 0 {
		c.Logger.Println("No pid to kill")
		return false
	}

	err := syscall.Kill(c.ActualPid(), code)
	var res bool
	if err == nil {
		res = true
	} else {
		c.Logger.Println("Failed to signal process", c.actual_pid, " with code", code, ":", err)
		res = false
	}
	return res
}

func (c *Process) isActualPidCached() bool {
	return !!c.CacheActualPid
}

func (c *Process) ActualPid() int {
	value := ""
	if c.PidCommand != "" {
		value, _ = c.PidFromCommand()
	} else {
		value, _ = c.PidFromFile()
	}
	//c.Logger.Println("PID ACTUAL:", value)
	var int_str int
	int_str, _ = strconv.Atoi(value)
	return int_str
}

func (c *Process) PidFromFile() (string, error) {
	//ap := strconv.Atoi(c.actual_pid)

	var int_pid int
	int_pid = int(c.actual_pid)
	stringed_pid := strconv.Itoa(int_pid)

	if c.CacheActualPid && c.actual_pid > 0 {
		return stringed_pid, nil
	} else {
		if len(c.PidFile) > 0 {
			dat, err := ioutil.ReadFile(c.PidFile)
			if err != nil {
				err := errors.New("pid_file " + c.PidFile + " does not exist or cannot be read")
				return "", err
			}
			var num_pid string
			num_pid = string(dat)
			int_str, _ := strconv.Atoi(num_pid)
			c.SetActualPid(int64(int_str))
			return string(dat), err
		} else {
			c.Logger.Println("pid_file", c.PidFile, " does not exist or cannot be read")
			err := errors.New("pid_file " + c.PidFile + " does not exist or cannot be read")
			return "", err
		}
	}

}
func (c *Process) PidFromCommand() (string, error) {
	//ps -ef | awk '/memcached$/{ print $2 }'
	// pid = %x{#{pid_command}}.strip
	// (pid =~ /\A\d+\z/) ? pid.to_i : nil
	opts := strings.Split(c.PidCommand, " ")
	out, err := exec.Command(opts[0], opts[1:]...).Output()
	str := string(out)
	return str, err
}

func (c *Process) SetActualPid(pid int64) {
	var p int = int(pid)
	proc.AppendPidToJournal(c.Name, p) // be sure to always log the pid
	c.actual_pid = pid
}

func (c *Process) ClearPid() {
	c.SetActualPid(0)
}

func (c *Process) UnlinkPid() {
	system.DeleteIfExists(c.pid_file)
}

func (c *Process) SkipTicksFor(seconds float64) {
	/*
	   TODO: should this be addative or longest wins?
	   i.e. if two calls for skip_ticks_for come in for 5 and 10, should it skip for 10 or 15?
	   self.skip_ticks_until = (self.skip_ticks_until || Time.now.to_i) + seconds.to_i
	*/

	var secs int64
	secs = int64(seconds)
	if c.skip_ticks_until > 0 {
		c.skip_ticks_until = time.Now().Unix() + secs
		c.Logger.Println("SKIP TICKS UNTIL", c.skip_ticks_until)
	} else {
		c.skip_ticks_until = c.skip_ticks_until + secs
		c.Logger.Println("SKIP TICKS UNTIL", c.skip_ticks_until)
	}

}

func (c *Process) isSkippingTicks() bool {
	t := time.Now()
	//c.skip_ticks_until = time.Now()
	value := false
	//if c.skip_ticks_until != nil { //&& c.skip_ticks_until > t { //time.Since(t).Seconds()
	if c.skip_ticks_until > t.Unix() {
		value = true
	}
	return value
}

func (c *Process) RefreshChildren() {

	// First prune the list of dead children
	for _, child := range c.Children {
		if !child.isProcessRunning(true) {
			//delete HERE!!!!!
		}
	}

	// Add new found children to the list
	new_children_pids := make([]map[string]interface{}, 0)
	childs_arr, _ := system.GetChildren(c.actual_pid)
	for _, pid := range childs_arr {
		if c.actual_pid != pid["Pid"].(int64) {
			new_children_pids = append(new_children_pids, pid)
		}
	}

	if len(new_children_pids) == 0 {
		//logger.info "Existing children: #{@children.collect{|c| c.actual_pid}.join(",")}. Got new children: #{new_children_pids.inspect} for #{actual_pid}"
		c.Logger.Println("Existing children: ")
		for _, ch := range c.Children {
			c.Logger.Println(ch.ActualPid())
		}
	}

	//Construct a new process wrapper for each new found children
	for _, child_pid := range new_children_pids {
		//ProcessJournal.append_pid_to_journal(name, child_pid)
		child_name := "<child(pid:" + child_pid["pid"].(string) + ")>"
		//logger = self.logger.prefix_with(child_name)
		child := c.ChildProcessFactory.CreateChildProcess(child_name, child_pid["pid"].(string), "logger")
		c.Children = append(c.Children, child)
	}
}

func (c *Process) SystemCommandOptions() map[string]interface{} {
	m := make(map[string]interface{}, 0)

	m["uid"] = c.Uid
	m["gid"] = c.Gid
	m["working_dir"] = c.WorkingDir
	m["environment"] = c.Environment
	m["pid_file"] = c.PidFile
	m["logger"] = c.Logger
	m["stdin"] = c.Stdin
	m["stdout"] = c.Stdout
	m["stderr"] = c.Stderr
	m["supplementary_groups"] = c.SupplementaryGroups

	return m
}

func (c *Process) PrepareCommand(command string) string {

	cmd := strings.Replace(command, "{{PID}}", strconv.Itoa(c.ActualPid()), 1)
	return cmd
}

/*
	func with_timeout(secs int , block func(int) int ) int{
	  secs += 100
	  secs = block(secs)
	  return secs
	}


	func callbackable(uno int) func(int) int  {

	    return func(uno int) int {
	        uno = uno + 20
	        return uno
	    }
	}

	func some_func(num int ) int{
	  num2 := with_timeout(num , callbackable(1) )
	  return num2
	}
*/

func (c *Process) WithTimeout(secs int, next_state string) { //secs int, next_state = nil, &blk) {
	/*
	   def with_timeout(secs, next_state = nil, &blk)
	     # Attempt to execute the passed block. If the block takes
	     # too long, transition to the indicated next state.
	     begin
	       Timeout.timeout(secs.to_f, &blk)
	     rescue Timeout::Error
	       logger.err "Execution is taking longer than expected."
	       logger.err "Did you forget to tell bluepill to daemonize this process?"
	       dispatch!(next_state)
	     end
	   end
	*/
	//Timeout.timeout(secs.to_f, &blk)

	c.Dispatch(next_state, "")
}

//PROCESS

type Trigger struct {
	Process *Process
	Logger  *log.Logger
	//mutex
	Name            string
	ScheduledEvents []string
}

func NewTrigger(process *Process, options map[string]interface{}) *Trigger {
	c := &Trigger{}
	c.Name = options["name"].(string)
	c.Process = process
	//c.Logger = options["logger"]
	c.ScheduledEvents = make([]string, 0)
	return c
}

func (c *Trigger) Reset() {
	//self.cancel_all_events

}

func (c *Trigger) Notify(transition string) {
	//raise "Implement in subclass"
}

func (c *Trigger) Dispatch() {
	//self.process.dispatch!(event, self.class.name.split("::").last)

}

func (c *Trigger) ScheduleEvent() {
	/*

	   # TODO: maybe wrap this in a ScheduledEvent class with methods like cancel
	   thread = Thread.new(self) do |trigger|
	     begin
	       sleep delay.to_f
	       trigger.dispatch!(event)
	       trigger.mutex.synchronize do
	         trigger.scheduled_events.delete_if { |_, thread| thread == Thread.current }
	       end
	     rescue StandardError => e
	       trigger.logger.err(e)
	       trigger.logger.err(e.backtrace.join("\n"))
	     end
	   end

	   self.scheduled_events.push([event, thread])

	*/
}

func (c *Trigger) CancellAllEvents() {
	/*
	   self.logger.info "Canceling all scheduled events"
	    self.mutex.synchronize do
	      self.scheduled_events.each {|_, thread| thread.kill}
	    end
	*/
}
