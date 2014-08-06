package main

import (
	"bitbucket.org/kardianos/osext"
	"flag"
	"fmt"
	godard "godard_cmd"
	cfg "godard_config"
	//"github.com/barakmich/glog"
	//"graph"
	app "application"
	//"github.com/davecheney/profile"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"syscall"
)

var Application string
var configFile = flag.String("config", "", "Path to an explicit configuration file.")
var LogFile = flag.String("l", "", "Path to logfile, defaults to #{options[:log_file]}")
var BaseDir = flag.String("c", "", "Directory to store godard socket and pid files, defaults to #{options[:base_dir]}")
var Privileged = flag.Bool("no-privileged", false, "Allow/disallow to run #{$0} as non-privileged process. disallowed by default")
var Timeout = flag.Int("t", 10, "Timeout for commands sent to the daemon, in seconds. Defaults to 10.")
var Attempts = flag.Int("attempts Count", 1, "Attempts for commands sent to the daemon, in seconds. Defaults to 1.")

var ApplicationCommands = []string{"status", "start", "stop", "restart", "unmonitor", "quit", "log"}

func Usage() {

	fmt.Println("Usage:")
	fmt.Println("  godard COMMAND [flags]\n")
	fmt.Println("Commands:")
	fmt.Println("    load CONFIG_FILE\t\tLoads new instance of godard using the specified config file")
	fmt.Println("    status\t\t\tLists the status of the proceses for the specified app")
	fmt.Println("    start [TARGET]\t\tIssues the start command for the target process or group, defaults to all processes")
	fmt.Println("    stop [TARGET]\t\tIssues the stop command for the target process or group, defaults to all processes")
	fmt.Println("    restart [TARGET]\t\tIssues the restart command for the target process or group, defaults to all processes")
	fmt.Println("    unmonitor [TARGET]\t\tStop monitoring target process or group, defaults to all processes")
	fmt.Println("    log [TARGET]\t\tShow the log for the specified process or group, defaults to all for app")
	fmt.Println("    quit\t\t\tStop godard")
	fmt.Println("See http://github.com/godard/godard#readme")

	fmt.Println("\nFlags:")
	flag.Parse()
	flag.PrintDefaults()
}

func main() {

	/*pprof_cfg := profile.Config{
		MemProfile:     true,
		NoShutdownHook: true, // do not hook SIGINT
	}
	// ensure profiling information is written to disk.
	p := profile.Start(&pprof_cfg)
	defer p.Stop() //must be called before the program exits to
	*/
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

	if len(*BaseDir) == 0 {
		if len(os.Getenv("GODARD_BASE_DIR")) > 0 {
			controller_opts["base_dir"] = os.Getenv("GODARD_BASE_DIR")
		} else {
			if syscall.Geteuid() != 0 {
				controller_opts["base_dir"] = path.Join(os.Getenv("HOME"), ".godard")
			} else {
				controller_opts["base_dir"] = "/var/run/godard"
			}
		}
	} else {
		controller_opts["base_dir"] = *BaseDir
	}

	controller_opts["log_file"] = *LogFile
	controller := app.NewController(controller_opts)

	//fmt.Println(controller_opts)
	basefile, _ := osext.Executable()
	//fmt.Println("CMD", cmd, "ARGS", os.Args, "CONTROLLER:", controller, "BASE FILE", basefile)
	//fmt.Println("BASE IN APPS", BaseName(basefile), controller.RunningApplications(),  stringInSlice(BaseName(basefile), controller.RunningApplications()))
	//fmt.Println("ARG IN APPS", os.Args[1], controller.RunningApplications(),  stringInSlice(os.Args[1], controller.RunningApplications()) )

	if stringInSlice(BaseName(basefile), controller.RunningApplications()) && isSymlink(basefile) {
		// godard was called as a symlink with the name of the target application
		controller_opts["application"] = basefile
		fmt.Println("godard was called as a symlink with the name of the target applicatio")

	} else if stringInSlice(os.Args[0], controller.RunningApplications()) {
		//the first arg is the application name
		controller_opts["application"] = os.Args[0]
		fmt.Println("ARGV SHIFT NAME:", controller_opts["application"])

	} else if stringInSlice(cmd, ApplicationCommands) {
		//fmt.Println("OPT 3")
		if len(controller.RunningApplications()) == 1 {
			// there is only one, let's just use that
			controller_opts["application"] = controller.RunningApplications()[0]
		} else if len(controller.RunningApplications()[0]) > 1 {
			// There is more than one, tell them the list and exit
			fmt.Println("You must specify an application name to run that command. Here's the list of running applications:")
			for app, index := range controller.RunningApplications() {
				fmt.Println("INDEX APP ", index, app)
			}
			fmt.Println("Usage: godard [app] cmd [options]")
			os.Exit(1)
		} else {
			// There are none running AND they aren't trying to start one
			fmt.Println("Error: There are no running godard daemons.\nTo start a godard daemon, use: godard load <config file>")
			os.Exit(2)
		}

	}
	//else {
	//	fmt.Println("NONE OF THE ABOVE")
	//}
	//fmt.Println("ARGS:", os.Args[0])

	if cmd == "load" {
		os.Setenv("GODARD_BASE_DIR", controller_opts["base_dir"].(string))
		config := cfg.ParseConfigFromFlagsAndFile(*configFile)
		//fmt.Println(config)
		godard.Init(config)
	} else {
		target := os.Args[len(os.Args)-1] //ARGV.shift
		if controller_opts["application"] != nil {
			//fmt.Println("HANDLE COMMAND NOW:", controller_opts["application"], cmd, target)
			app := controller_opts["application"].(string)
			controller.HandleCommand(app, cmd, target)
		} else {
			fmt.Println("Don´t know how to handle", cmd, target)
		}

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
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
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
