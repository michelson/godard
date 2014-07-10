package main

import (
  godard "godard_cmd"
  cfg "godard_config"
  "bitbucket.org/kardianos/osext"
  "flag"
  "fmt"
  //"github.com/barakmich/glog"
  //"graph"
  "os"
  "runtime"
  "syscall"
  "path/filepath"
)


var Application string
var configFile = flag.String("config", "", "Path to an explicit configuration file.")
var LogFile = flag.String("l", "", "Path to logfile, defaults to #{options[:log_file]}")
var BaseDir = flag.String("c", "", "Directory to store bluepill socket and pid files, defaults to #{options[:base_dir]}")
var Privileged = flag.Bool("no-privileged", false, "Allow/disallow to run #{$0} as non-privileged process. disallowed by default")
var Timeout = flag.Int("t", 10, "Timeout for commands sent to the daemon, in seconds. Defaults to 10.")
var Attempts = flag.Int("attempts Count", 1 , "Attempts for commands sent to the daemon, in seconds. Defaults to 1.") 

var ApplicationCommands = []string{"status", "start", "stop", "restart", "unmonitor", "quit", "log"}


func Usage() {
  /*fmt.Println("Godard is monitoring tool.\n")
  fmt.Println("Usage:")
  fmt.Println("  godard COMMAND [flags]\n")
  fmt.Println("Commands:")
  fmt.Println("  init\tStart server , load config.")
  fmt.Println("\nFlags:")*/

  fmt.Println("Usage:")
  fmt.Println("  godard COMMAND [flags]\n")
  fmt.Println( "Commands:")
  fmt.Println( "    load CONFIG_FILE\t\tLoads new instance of godard using the specified config file")
  fmt.Println( "    status\t\t\tLists the status of the proceses for the specified app")
  fmt.Println( "    start [TARGET]\t\tIssues the start command for the target process or group, defaults to all processes")
  fmt.Println( "    stop [TARGET]\t\tIssues the stop command for the target process or group, defaults to all processes")
  fmt.Println( "    restart [TARGET]\t\tIssues the restart command for the target process or group, defaults to all processes")
  fmt.Println( "    unmonitor [TARGET]\t\tStop monitoring target process or group, defaults to all processes")
  fmt.Println( "    log [TARGET]\t\tShow the log for the specified process or group, defaults to all for app")
  fmt.Println( "    quit\t\t\tStop godard")
  fmt.Println( "See http://github.com/godard/godard#readme")

  fmt.Println("\nFlags:")
  flag.Parse()
  flag.PrintDefaults()
}

func main() {

  if os.Getenv("GOMAXPROCS") == "" {
    runtime.GOMAXPROCS(runtime.NumCPU())
  }

  // No command? It's time for usage.
  if len(os.Args) == 1 {
    Usage()
    os.Exit(1)
  }
  cmd := os.Args[1]
  newargs := make([]string, 0)
  newargs = append(newargs, os.Args[0])
  newargs = append(newargs, os.Args[2:]...)
  os.Args = newargs
  flag.Parse()

  // Check for root
  //fmt.Println(syscall.Geteuid())
  //fmt.Println(*Privileged)
  if *Privileged && syscall.Getgid() != 0 {
    os.Stderr.Write([]byte("You must run godard as root or use --no-privileged option."))
    os.Exit(3)
  }
  
  var controller_opts map[string]interface{} = make(map[string]interface{}, 0)

  controller_opts["base_dir"] = *BaseDir
  controller_opts["log_file"] = *LogFile
  controller := NewController( controller_opts )

  fmt.Println(controller)
  basefile, _ := osext.Executable()
  //if controller.running_applications.include?(File.basename($0)) && File.symlink?($0)
  //  options[:application] = File.basename($0)
  if stringInSlice(BaseName(basefile), controller.RunningApplications()) && isSymlink(basefile) {
  } else if stringInSlice(os.Args[0], controller.RunningApplications()){
    //options[:application] = File.basename($0)
  } else if stringInSlice(os.Args[0], ApplicationCommands){
    /*

    if controller.running_applications.length == 1
      # there is only one, let's just use that
      options[:application] = controller.running_applications.first
    elsif controller.running_applications.length > 1
      # There is more than one, tell them the list and exit
      $stderr.puts "You must specify an application name to run that command. Here's the list of running applications:"
      controller.running_applications.each_with_index do |app, index|
        $stderr.puts "  #{index + 1}. #{app}"
      end
      $stderr.puts "Usage: bluepill [app] cmd [options]"
      exit(1)
    else
      # There are none running AND they aren't trying to start one
      $stderr.puts "Error: There are no running bluepill daemons.\nTo start a bluepill daemon, use: bluepill load <config file>"
      exit(2)
    end

    */
  }

  
  fmt.Println("ARGS:", os.Args[0])

  switch cmd {
    case "init":
      config := cfg.ParseConfigFromFlagsAndFile(*configFile)
      godard.Init(config )
    case "status":
      fmt.Println("CURRENT STATUS", cmd)
    case "start":
      fmt.Println("STARTING CMD", cmd)
    case "stop":
      fmt.Println("STOP PROCESS", cmd)
    case "quit":
      fmt.Println("TERMINATING QUIT", cmd)
    case "log":
      fmt.Println("INIT LOGGING", cmd)
    default:
      fmt.Println("No command", cmd)
      flag.Usage()
  }
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func isSymlink(filename string) bool {
  fi, err := os.Lstat(filename)
  res := false
  if fi.Mode() & os.ModeSymlink == os.ModeSymlink {
    res = true 
  } else {
    res = false
  }
  fmt.Println("error detecting symlink", err)
  return res
}

func BaseName(file string) string {
  fName := filepath.Base(file)
  return string(fName)
}
