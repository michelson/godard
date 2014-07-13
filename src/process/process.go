package process

import (
	//"io"
	"io/ioutil"
	"log"
	watcher "watcher"
	//trigger "trigger"
	time "time"
	"strings"
	"sync"
	system "system"
	fsm "github.com/looplab/fsm"
  //"sync/atomic"
)

type Process struct {

	Name string

	pid int
	ppid int
	cpu float64
	mem int
	elapsed int
	command string

	Watches  []*watcher.ConditionWatch
	Triggers []*Trigger
	Children []*Process
  Statistics *ProcessStatistics

	pid_file string
	pre_start_command string
	StartCommand string
	StopCommand string
	RestartCommand string

	Stdout string
	Stderr string
	Stdin string

	Daemonize bool
	PidFile string
	WorkingDir string
	Environment map[string]string

	start_grace_time int
	stop_grace_time int
	restart_grace_time int

	Uid string
	Gid string

	actual_pid int64
	cache_actual_pid bool

	monitor_children bool
	child_process_factory string

	pid_command string
	auto_start bool

	supplementary_groups string

	stop_signals string

	on_start_timeout string

	Group_start_noblock bool
	Group_restart_noblock bool
	Group_stop_noblock bool
	Group_unmonitor_noblock bool


	skip_ticks_until time.Time
	process_running bool
	state string

	event_mutex *sync.Mutex

	Logger string

	state_machine *fsm.FSM

  ListenerChannel chan map[string]string

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewProcess(process_name string, checks map[string]interface{}, options map[string]interface{}) *Process {
	c := &Process{}
	c.Name = process_name
  c.StartCommand = options["start_command"].(string)
  c.PidFile = options["pid_file"].(string)
  c.WorkingDir = options["working_dir"].(string)
  //p.AddWatches(c.Watches)
	c.event_mutex = &sync.Mutex{}
	c.Watches = make([]*watcher.ConditionWatch, 0)
	c.Triggers = make([]*Trigger, 0)
	c.Children = make([]*Process, 0)
	// @threads = []
	// @statistics = ProcessStatistics.new
  c.Statistics = NewProcessStatistics()
	// @actual_pid = options[:actual_pid]
	// self.logger = options[:logger]

  c.ListenerChannel = make(chan map[string]string)

	  for check_name , value := range(checks){
	  	trigger_exists := false
	  	for _, v := range(c.Triggers){
	  		if v.Name == check_name{
	  			trigger_exists = true
	  			break
	  		}
	  	}
	  	if trigger_exists {
				log.Println("add trigger here:", check_name, value)
				c.AddTrigger(check_name, value)
  		}else{
  			//log.Println("add watch here:", check_name, value)
  			c.AddWatch(check_name, value)
  		}
	  	
	  }

	  // These defaults are overriden below if it's configured to be something else.
	  c.monitor_children =  false
	  c.cache_actual_pid = true
	  c.start_grace_time = 3
	  c.stop_grace_time = 3
	  c.restart_grace_time = 3
	  //@environment = {}
	  c.on_start_timeout = "start"
	  c.Group_start_noblock = true
	  c.Group_stop_noblock = true
	  c.Group_restart_noblock = true
	  c.Group_unmonitor_noblock = true

	  c.auto_start = true
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

	        {Name: "tick_up", Src: []string{"restarting"}, Dst: "up"},
	        {Name: "tick_down", Src: []string{"restarting"}, Dst: "down"},

	       	{Name: "start", Src: []string{"unmonitored", "down"}, Dst: "starting"},

	       	{Name: "restart", Src: []string{"up", "down"}, Dst: "restarting"},
	    },
	    fsm.Callbacks{
	    		"before_event": func(e *fsm.Event) {
						c.NotifyTriggers( c.state_machine.Current() )
						if c.state_machine.Current() != "stopping" {
							c.CleanThreads()
						}
					},
					"after_event": func(e *fsm.Event) {
          
						if c.state_machine.Current() == "starting" {

							c.StartProcess()
						}

						if c.state_machine.Current() == "stopping" {
							c.StopProcess()
						}

						if c.state_machine.Current() == "stopping" {
							c.StopProcess()
						}

						if c.state_machine.Current() == "restarting" {
							c.RestartProcess()
						}

						c.RecordTransition( c.state_machine.Current() )
					},
	    },
	)

	log.Println("CREATING PROCESS:", c.Name)

	return c
}

func (c *Process) Tick(){

	if c.isSkippingTicks(){
		log.Println("SKIPPING TICKS")
	}else{
		//c.skip_ticks_until = nil
		c.process_running = false
		
		// Deal with thread cleanup here since the stopping state isn't used
    //clean_threads if self.unmonitored?
    if c.state_machine.Current() == "unmonitored"{
			c.CleanThreads()
    }
    // run state machine transitions
    if c.isProcessRunning(false){
    	log.Println("TICKS UP")
    	c.state_machine.Event("tick_up")	
    }else{
    	log.Println("TICKS DOWN")
    	c.state_machine.Event("tick_down")	
    }
    
    log.Println("CURRENT STATE:", c.state_machine.Current())

		if c.isUp() {
			c.RunWatches()
			if c.monitor_children {
         c.RefreshChildren()
         for _ , child := range(c.Children){
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

func (c *Process) isUp() bool{
	return c.state_machine.Current() == "up"
}

func (c *Process) Dispatch(event string, reason string) {
 /*
      c.event_mutex.synchronize do
        @statistics.record_event(event, reason)
        self.send("#{event}")
      end
 */
    c.event_mutex.Lock()
    //@statistics.record_event(event, reason)
    //self.send("#{event}")

    /*aa := make(map[string]int64,0)
    aa[c.state_machine.Current()] = 23
    c.ListenerChannel <- aa*/
    c.state_machine.Event(event)
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
}

func (c *Process) NotifyTriggers(transition string) {
	// self.triggers.each {|trigger| trigger.notify(transition)}
	for _ , tgr := range(c.Triggers){
		tgr.Notify(transition)
	}
}


func (c *Process) AddTrigger(name string, value interface{}) {
   //   self.triggers << Trigger[name].new(self, options.merge(:logger => self.logger))
	v := value.(map[string]interface{})
	//m["name"] = name
	//m["logger"] = c.logger
	c.Triggers = append(c.Triggers, NewTrigger(c , v))
}

func (c *Process) AddWatches(options map[string]interface{}){

	if len(options) > 0 {
		log.Println("ADDING WATCHES TO PROCESS:", len(options))
		for k, v := range options {
			c.AddWatch(k, v)
		}
	}
}

func (c *Process) AddWatch(name string, value interface{}) {
  log.Println("CHECKS:", name, value )
  
  v := value.(map[string]interface{})
  //log.Println(v["every"])

  condition := watcher.NewConditionWatch(name, v)
  c.Watches = append(c.Watches , condition)
}

//NOK
func (c *Process) RunWatches() {

	/*now := time.Now().Unix()
	for _ , watch := range(c.Watches){

	}*/
	/*
      now = Time.now.to_i

      threads = self.watches.collect do |watch|
        [watch, Thread.new { Thread.current[:events] = watch.run(self.actual_pid, now) }]
      end

      @transitioned = false

      threads.inject([]) do |events, (watch, thread)|
        thread.join
        if thread[:events].size > 0
          logger.info "#{watch.name} dispatched: #{thread[:events].join(',')}"
          thread[:events].each do |event|
            events << [event, watch.to_s]
          end
        end
        events
      end.each do |(event, reason)|
        break if @transitioned
        self.dispatch!(event, reason)
      end

	*/

}

func (c *Process) DetermineInitialState(){
/*
      if self.process_running?(true)
        self.state = 'up'
      else
        self.state = (auto_start == false) ? 'unmonitored' : 'down' # we need to check for false value
      end

*/
   	if c.isProcessRunning(true){
   		c.state_machine.SetCurrent("up")
   	}else{
   		//(auto_start == false) ? 'unmonitored' : 'down' # we need to check for false value
   		if c.auto_start == false {
   			c.state_machine.SetCurrent("unmonitored")
   		}else{
   			c.state_machine.SetCurrent("down")
   		}
   	}

   	log.Println("DETERMINE INITAL STATE", c.state_machine.Current())

}


// System Process Methods

func (c*Process) isProcessRunning(force bool) bool{

  if force {
	  c.process_running = false 
  } 
 
	//@process_running ||= signal_process(0)

	// the process isn't running, so we should clear the PID
	if !c.process_running {
		c.ClearPid()	
	}
	
  return c.process_running
}


func (c *Process) HandleUserCommand(cmd string){
    switch cmd {
    case "start": 
    	  if c.isProcessRunning(true){
          log.Println("Refusing to re-run start command on an already running process.")
    	  }else{
          c.Dispatch("start", "user initiated")
    	  }
    case "stop": 
        c.StopProcess()
        c.Dispatch("unmonitor", "user initiated")
    case "restart": c.RestartProcess()
    case "unmonitor": 
        // When the user issues an unmonitor cmd, reset any triggers so that
        // scheduled events gets cleared
    		for _, trgr := range(c.Triggers){
    			trgr.Reset()
    		}
        c.Dispatch("unmonitor", "user initiated")
    default: log.Println("default")
    }
}

func (c *Process) StartProcess(){
    c.PreStartProcess()
    log.Println( "Executing start command:", c.StartCommand)
    if c.isDaemonized() {
    	/* daemon_id = System.daemonize(start_command, self.system_command_options)
        if daemon_id
          ProcessJournal.append_pid_to_journal(name, daemon_id)
          children.each {|child|
            ProcessJournal.append_pid_to_journal(name, child.actual_id)
          } if self.monitor_children?
        end
        daemon_id*/

    }else{
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
    
      	result := system.ExecuteBlocking(c.StartCommand, c.SystemCommandOptions())
				log.Println("EXEC RESULT :", result)
				c.ListenerChannel <- result
        //result = System.execute_blocking(start_command, self.system_command_options)

        //unless result[:exit_code].zero?
        //  logger.warning "Start command execution returned non-zero exit code:"
        //  logger.warning result.inspect
        //end


    }

    c.SkipTicksFor(c.start_grace_time)
}

func (c *Process) PreStartProcess(){
	if c.pre_start_command != ""{
		log.Println("Executing pre start command:", c.pre_start_command )
		result := system.ExecuteBlocking(c.pre_start_command, c.SystemCommandOptions())
		log.Println("PRE START COMMAND RESULT :", result)
		if result["exit_code"] != "0" {
			log.Println("Pre start command execution returned non-zero exit code:")
			log.Println(result)
		}
	}
}
//NOK
func (c *Process) StopProcess(){
	/*
      if monitor_children
        System.get_children(self.actual_pid).each do |child_pid|
          ProcessJournal.append_pid_to_journal(name, child_pid)
        end
      end

      if stop_command
        cmd = self.prepare_command(stop_command)
        logger.warning "Executing stop command: #{cmd}"

        with_timeout(stop_grace_time, "stop") do
          result = System.execute_blocking(cmd, self.system_command_options)

          unless result[:exit_code].zero?
            logger.warning "Stop command execution returned non-zero exit code:"
            logger.warning result.inspect
          end
        end

      elsif stop_signals
        # issue stop signals with configurable delay between each
        logger.warning "Sending stop signals to #{actual_pid}"
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
      else
        logger.warning "Executing default stop command. Sending TERM signal to #{actual_pid}"
        signal_process("TERM")
      end
      ProcessJournal.kill_all_from_journal(name) # finish cleanup
      self.unlink_pid # TODO: we only write the pid file if we daemonize, should we only unlink it if we daemonize?

      self.skip_ticks_for(stop_grace_time)
	*/
}

func (c *Process) RestartProcess(){

	if c.RestartCommand != ""{
		cmd := c.PrepareCommand(c.RestartCommand)
		log.Println("Executing restart command:", cmd)
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

    c.SkipTicksFor(c.restart_grace_time)

	} else {
		log.Println("No RestartCommand specified. Must stop and start to restart")
    c.StopProcess()
    c.StartProcess()
	}

}

func (c *Process) CleanThreads(){
	//@threads.each { |t| t.kill }
  //@threads.clear
}

func (c *Process) isDaemonized() bool{
	return !!c.Daemonize
}

func (c *Process) isMonitorChildren() bool{
	return !!c.monitor_children
}

func (c *Process) SignalProcess(code int){
	/*

      code = code.to_s.upcase if code.is_a?(String) || code.is_a?(Symbol)
      ::Process.kill(code, actual_pid)
      true
    rescue Exception => e
      logger.err "Failed to signal process #{actual_pid} with code #{code}: #{e}"
      false

	*/
}

func (c *Process) isActualPidCached() bool{
	return !!c.cache_actual_pid
}

func (c *Process) ActualPid() string{
	value := ""
  if c.pid_command != "" {
  	value, _ = c.PidFromCommand()
  } else {
  	value, _ = c.PidFromFile()
  } 
  return value
}

func (c *Process) PidFromFile() (string, error) {
	dat, err := ioutil.ReadFile(c.Name)
	check(err)
	//log.Println(string(dat))
	return string(dat) , err

}
func (c *Process) PidFromCommand() (string, error) {
  // pid = %x{#{pid_command}}.strip
  // (pid =~ /\A\d+\z/) ? pid.to_i : nil
  log.Println("PID COMMAND NOT IMPLEMENTED YET")
  return "none", nil
}

func (c *Process) SetActualPid(pid int64) {
	//ProcessJournal.append_pid_to_journal(name, pid) # be sure to always log the pid
  c.actual_pid = pid
}

func (c *Process) ClearPid() {
	c.actual_pid = 0
}

func (c *Process) UnlinkPid() {
  system.DeleteIfExists(c.pid_file)
}

func (c *Process) SkipTicksFor(seconds int) {}

func (c *Process) isSkippingTicks() bool {
	t := time.Now()
	//c.skip_ticks_until = time.Now()
	value := false
	//if c.skip_ticks_until != nil { //&& c.skip_ticks_until > t { //time.Since(t).Seconds() 
	if c.skip_ticks_until.Unix() > t.Unix() {
		value = true
	}
	return value
}

func (c *Process) RefreshChildren() {
	/*

      # First prune the list of dead children
      @children.delete_if {|child| !child.process_running?(true) }

      # Add new found children to the list
      new_children_pids = System.get_children(self.actual_pid) - @children.map {|child| child.actual_pid}

      unless new_children_pids.empty?
        logger.info "Existing children: #{@children.collect{|c| c.actual_pid}.join(",")}. Got new children: #{new_children_pids.inspect} for #{actual_pid}"
      end

      # Construct a new process wrapper for each new found children
      new_children_pids.each do |child_pid|
        ProcessJournal.append_pid_to_journal(name, child_pid)
        child_name = "<child(pid:#{child_pid})>"
        logger = self.logger.prefix_with(child_name)

        child = self.child_process_factory.create_child_process(child_name, child_pid, logger)
        @children << child
      end

	*/
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
  m["supplementary_groups"] = c.supplementary_groups
 
  return m
}

func (c *Process) PrepareCommand(command string) string {
	cmd := strings.Replace( command , "{{PID}}", c.ActualPid(), 1)
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
  Logger string
  //mutex 
  Name string
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

func (c*Trigger) Reset(){
  //self.cancel_all_events

}

func (c*Trigger) Notify(transition string){
  //raise "Implement in subclass"
}

func (c*Trigger) Dispatch(){
  //self.process.dispatch!(event, self.class.name.split("::").last)

}

func (c*Trigger) ScheduleEvent(){
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

func (c*Trigger) CancellAllEvents(){
  /*
     self.logger.info "Canceling all scheduled events"
      self.mutex.synchronize do
        self.scheduled_events.each {|_, thread| thread.kill}
      end
  */
}