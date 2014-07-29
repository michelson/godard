package application

import (
	//app "application"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	socket "socket"
	"strconv"
	"syscall"
	system "system"
)

type Controller struct {
	BaseDir    string
	LogFile    string
	SocketsDir string
	PidsFile   string
	pid        string
}

func NewController(options map[string]interface{}) *Controller {
	c := &Controller{}
	c.LogFile = options["log_file"].(string)
	c.BaseDir = options["base_dir"].(string)
	c.SocketsDir = path.Join(c.BaseDir, "sock")
	c.PidsFile = path.Join(c.BaseDir, "pids")

	c.setup_dir_structure()
	c.cleanup_godard_directory()
	//log.Println("CREATING NEW CONTROLLER")
	return c
}

func (c *Controller) RunningApplications() []string {
	var arr []string
	pp := path.Join(c.SocketsDir)
	files, _ := ioutil.ReadDir(pp)
	for _, f := range files {
		fName := filepath.Base(f.Name())
		extName := filepath.Ext(f.Name())
		bName := fName[:len(fName)-len(extName)]
		arr = append(arr, bName)
	}
	return arr
}

func (c *Controller) HandleCommand(application string, command string, args ...string) {
	switch command {
	case "status":
		log.Println("status")
		c.send_to_daemon(application, command, args)
	case "start", "stop", "restart", "unmonitor":
		affected := c.send_to_daemon(application, command, args)
		log.Println("AFFECTED:", affected)
		/*if len(affected) == 0 {
			//log.Println("No processes affected")
		} else {
			//log.Println("SOME EXTRA ARGS:", args)
			//log.Println("SENT", command, "CMD TO:", affected)
		}*/
	case "quit":
		log.Println("HANDLE", command, " COMMAND FOR:", application, args)

		pid, err := c.PidFor(application)
		log.Println("quit", pid)

		if err != nil {
		} else {
			if system.IsPidAlive(pid) {
				err := syscall.Kill(pid, syscall.SIGTERM)
				//process, _ := os.FindProcess(int(pid))
				//err := process.Kill()
				if err != nil {
					log.Println("error Killing Godard: ", err)
				} else {
					log.Println("Killing Godard", pid)
				}
			} else {
				log.Println("godard", pid, " not running")
			}
		}

	default:
		log.Println("Unknown command", command, "or application", command, "has not been loaded yer")
		os.Exit(1)
	}
}

func (c *Controller) send_to_daemon(application string, command string, args []string) string {
	var res string
	if c.verify_version(application) {
		cmd := command
		for _, arg := range args {
			cmd = cmd + ":" + arg
		}
		response, err := socket.ClientCommand(c.BaseDir, application, cmd)
		if err != nil {
			log.Println("Received error from server:")
			log.Println(response)
			res = response
			os.Exit(8)
		} else {
			log.Println("sucess: ", response)
			res = response
		}

	}
	return res
}

func (c *Controller) grep_pattern(application *Application, query string) {
	/*
	   pattern = [application, query].compact.join(':')
	   ['\[.*', Regexp.escape(pattern), '.*'].compact.join
	*/
}

func (c *Controller) cleanup_godard_directory() {

	for _, app := range c.RunningApplications() {
		pid, err := c.PidFor(app)
		log.Println("CLEANUP:", pid, system.IsPidAlive(pid))
		if err != nil || !system.IsPidAlive(pid) {
			pid_file := path.Join(c.PidsFile, app, app+".pid")
			sock_file := path.Join(c.SocketsDir, app, app+".sock")
			system.DeleteIfExists(pid_file)
			system.DeleteIfExists(sock_file)
		}
	}
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//def PidFor(app)
func (c *Controller) PidFor(app string) (int, error) {

	pid_file := path.Join(c.PidsFile, app, app+".pid")
	dat, err := ioutil.ReadFile(pid_file)

	var num_pid string
	num_pid = string(dat)
	int_str, _ := strconv.Atoi(num_pid)
	log.Println("PID FOR ", pid_file, app, "IS", int_str)
	return int_str, err
}

func (c *Controller) setup_dir_structure() {
	var arr []string
	arr = append(arr, c.SocketsDir)
	arr = append(arr, c.PidsFile)
	for _, a := range arr {

		err := os.MkdirAll(a, 0777)
		if err != nil {
			log.Println("error creating dir", a, err)
		}
	}
}

func (c *Controller) verify_version(application string) bool {
	/*  begin
	      version = Socket.client_command(base_dir, application, "version")
	      if version != Bluepill::VERSION
	        abort("The running version of your daemon seems to be out of date.\nDaemon Version: #{version}, CLI Version: #{Bluepill::VERSION}")
	      end
	    rescue ArgumentError
	      abort("The running version of your daemon seems to be out of date.")
	    end
	*/
	return true
}
