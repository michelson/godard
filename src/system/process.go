package system

import (
    //"io"
    "io/ioutil"
    "log"
)

type Process struct {
    //PID  PPID  %CPU    RSS     ELAPSED COMMAND
    Name string

    pid int
    ppid int
    cpu float64
    mem int
    elapsed int
    command string


    pre_start_command string
    StartCommand string
    stop_command string
    restart_command string

    stdout string
    stderr string
    stdin string

    daemonize string
    PidFile string
    working_dir string
    environment string

    start_grace_time string
    stop_grace_time string
    restart_grace_time string

    uid string
    gid string

    cache_actual_pid string

    monitor_children string
    child_process_factory string

    pid_command string
    auto_start string

    supplementary_groups string

    stop_signals string

    on_start_timeout string

    group_start_noblock string
    group_restart_noblock string
    group_stop_noblock string
    group_unmonitor_noblock string
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

  //self.watches << ConditionWatch.new(name, options.merge(:logger => self.logger))
}

func (c *Process) AddTrigger(name string) {
   //   self.triggers << Trigger[name].new(self, options.merge(:logger => self.logger))
}


/*
    def start_process

    end

    def pre_start_process

    end

    def stop_process


    end

    def restart_process

    end
*/

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func (c *Process) PidFromFile(name string) (string, error) {
    dat, err := ioutil.ReadFile(name)
    check(err)
    //log.Println(string(dat))
    return string(dat) , err

}


func (c *Process) Tick()  {

}

func (c *Process) DetermineInitialState(){
    
}

/*
    def pid_from_file
      return @actual_pid if cache_actual_pid? && @actual_pid
      @actual_pid = begin
        if pid_file
          if File.exists?(pid_file)
            str = File.read(pid_file)
            str.to_i if str.size > 0
          else
            logger.warning("pid_file #{pid_file} does not exist or cannot be read")
            nil
          end
        end
      end
    end
*/
//run watches