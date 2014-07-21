//package process

package system

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	//"os/signal"
	"bufio"
	"syscall"
)

var Store []map[string]interface{}

func PidAlive(pid int) bool {
	err := syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		log.Println("PID ALIVE?", err)
		return false
	} else {
		return true
	}
	/*   def pid_alive?(pid)
	     begin
	       ::Process.kill(0, pid)
	       true
	     rescue Errno::EPERM # no permission, but it is definitely alive
	       true
	     rescue Errno::ESRCH
	       false
	     end
	   end */
}

func PsAxu() ([]map[string]interface{}, error) {
	cmd := exec.Command("ps", "axo", "pid,ppid,pcpu,rss,etime,command")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]map[string]interface{}, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		//log.Println(tokens)

		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}

		pid, err := strconv.Atoi(ft[0])
		if err != nil {
			continue
		}
		ppid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		mem, err := strconv.Atoi(ft[3])
		if err != nil {
			continue
		}

		d := ParseElapsedTime(string(ft[4]))

		m := map[string]interface{}{
			"pid":     pid,
			"ppid":    ppid,
			"cpu":     cpu,
			"mem":     mem,
			"elapsed": d,
			"command": ft[5],
		}

		processes = append(processes, m)
	}
	Store = processes
	return processes, err
}

func FindByPid(pid int, processes []map[string]interface{}) (map[string]interface{}, error) {
	var pp map[string]interface{}

	for _, p := range processes {
		if p["pid"].(int) == pid {
			pp = p
			break
		}
	}
	return pp, nil
}

func CpuUsage(pid int) (float64, error) {
	ps, err := PsAxu()
	process, err := FindByPid(pid, ps)
	//log.Println("PROCESS:", process["pid"].(int))
	return process["cpu"].(float64), err
}

func MemoryUsage(pid int) (int, error) {
	ps, err := PsAxu()
	process, err := FindByPid(pid, ps)
	return process["mem"].(int) / 1024, err
}

func RunningTime(pid int) (int, error) {
	ps, err := PsAxu()
	process, err := FindByPid(pid, ps)
	return process["elapsed"].(int), err
}

func Command(pid int) (string, error) {
	ps, err := PsAxu()
	process, err := FindByPid(pid, ps)
	return process["command"].(string), err
}

func GetChildren(parent_pid int64) ([]map[string]interface{}, error) {
	child_pids := make([]map[string]interface{}, 0)

	processes, err := PsAxu()

	for _, p := range processes {
		if p["ppid"].(int64) == parent_pid {
			child_pids = append(child_pids, p)
		}
	}

	for _, p := range child_pids {
		gcp, _ := GetChildren(p["pid"].(int64))
		for _, ppp := range gcp {
			child_pids = append(child_pids, ppp)
		}
	}

	return child_pids, err
}

func Daemonize(cmd string, opts map[string]string) {
	/*
	   # Returns the pid of the child that executes the cmd
	   def daemonize(cmd, options = {})
	     rd, wr = IO.pipe

	     if child = Daemonize.safefork
	       # we do not wanna create zombies, so detach ourselves from the child exit status
	       ::Process.detach(child)

	       # parent
	       wr.close

	       daemon_id = rd.read.to_i
	       rd.close

	       return daemon_id if daemon_id > 0

	     else
	       # child
	       rd.close

	       drop_privileges(options[:uid], options[:gid], options[:supplementary_groups])

	       # if we cannot write the pid file as the provided user, err out
	       exit unless can_write_pid_file(options[:pid_file], options[:logger])

	       to_daemonize = lambda do
	         # Setting end PWD env emulates bash behavior when dealing with symlinks
	         Dir.chdir(ENV["PWD"] = options[:working_dir].to_s)  if options[:working_dir]
	         options[:environment].each { |key, value| ENV[key.to_s] = value.to_s } if options[:environment]

	         redirect_io(*options.values_at(:stdin, :stdout, :stderr))

	         ::Kernel.exec(*Shellwords.shellwords(cmd))
	         exit
	       end

	       daemon_id = Daemonize.call_as_daemon(to_daemonize, nil, cmd)

	       File.open(options[:pid_file], "w") {|f| f.write(daemon_id)}

	       wr.write daemon_id
	       wr.close

	       ::Process::exit!(true)
	     end
	   end*/
}

func ParseElapsedTime(str string) int {
	// [[dd-]hh:]mm:ss
	//s := "02-09:38:38"
	r, _ := regexp.Compile(`(?:(?:(\d+)-)?(\d\d):)?(\d\d):(\d\d)`)
	res := r.FindStringSubmatch(str)

	days, _ := strconv.ParseFloat(res[1], 64)
	hours, _ := strconv.ParseFloat(res[2], 64)
	mins, _ := strconv.ParseFloat(res[3], 64)
	secs, _ := strconv.ParseFloat(res[4], 64)
	//var tt time.Duration = (( (days*24 + hours)*60 + mins)*60 + secs )* time.Second
	//log.Println( days , hours,mins, secs  )

	minutes := (((days*24+hours)*60+mins)*60 + secs) / 60

	return int(minutes)
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

func DeleteIfExists(filename string) {

	exists, _ := FileExists(filename)

	if exists {
		err := os.Remove(filename)
		if err != nil {
			log.Fatal("Warning: permission denied trying to delete #{filename}")
		}
	}
	/*
	   tries = 0

	   begin
	     File.unlink(filename) if filename && File.exists?(filename)
	   rescue IOError, Errno::ENOENT
	   rescue Errno::EACCES
	     retry if (tries += 1) < 3
	     $stderr.puts("Warning: permission denied trying to delete #{filename}")
	   end
	*/

}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

//http://stackoverflow.com/questions/10781516/how-to-pipe-several-commands
//http://stackoverflow.com/questions/10385551/get-exit-code-go
func disabledExecuteBlocking(command string, options map[string]interface{}) map[string]string {
	c1 := exec.Command(command)

	r, w := io.Pipe()
	c1.Stdout = w
	c1.Stdin = r

	var b2 bytes.Buffer
	c1.Stdout = &b2

	var b1 bytes.Buffer
	c1.Stdin = &b1

	c1.Start()
	c1.Wait()

	m := make(map[string]string, 0)

	if err := c1.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
				var i64 string
				i64 = string(status.ExitStatus())
				m["exit_status"] = i64
			}
		} else {
			m["exit_status"] = "0" //status.ExitStatus()
			// log.Fatalf("cmd.Wait: %v", err)
		}
	}

	w.Close()
	//so,  := io.Copy(os.Stdout, &b2)
	//m["stdout"], _ = io.Copy(os.Stdout, &b2)
	//m["stderr"], _ = io.Copy(os.Stderr, &b1)
	m["exit_code"] = "0"
	log.Println("EXEC OPTIONS:", options)
	return m
}

func ExecuteBlocking(command string, options map[string]interface{}) map[string]string {

	//if options["working_dir"]
	working_dir := options["working_dir"].(string)
	os.Setenv("PWD", working_dir)
	os.Chdir(working_dir)

	args := strings.Split(command, " ")
	m := make(map[string]string)
	subProcess := exec.Command(args[0], args[1:]...) //Just for testing, replace with your subProcess

	_, err := subProcess.StdinPipe()
	/*
	   if err != nil {
	       log.Println(err) //replace with logger, or anything you want
	   }

	   stdout, err := subProcess.StdoutPipe()
	   if err != nil {
	       log.Println(err) //replace with logger, or anything you want
	   }*/

	defer subProcess.Wait()
	//defer stdin.Close() // the doc says subProcess.Wait will close it, but I'm not sure, so I kept this line

	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stderr

	//log.Println("START") //for debug
	if err = subProcess.Start(); err != nil { //Use start, not run
		log.Println("An error occured: ", err) //replace with logger, or anything you want
	}

	/*buf := new(bytes.Buffer)
	  buf.ReadFrom(stdout)
	  s := buf.String()*/

	//subProcess.Wait()

	m["stdout"] = "blabla"
	m["exit_code"] = "0"
	//m["stdin"] = sin

	//log.Println(m)
	//log.Println("END") //for debug
	//log.Println("EXEC OPTIONS:", options)
	return m
}

func aaExecuteBlocking(command string, options map[string]interface{}) map[string]string {
	m := make(map[string]string)
	m["stdout"] = "ss"
	m["exit_code"] = "0"
	return m
}

/*
   # Returns the stdout, stderr and exit code of the cmd
   def execute_blocking(cmd, options = {})
     rd, wr = IO.pipe

     if child = Daemonize.safefork
       # parent
       wr.close

       cmd_status = rd.read
       rd.close

       ::Process.waitpid(child)

       cmd_status.strip != '' ? Marshal.load(cmd_status) : {:exit_code => 0, :stdout => '', :stderr => ''}
     else
       # child
       rd.close

       # create a child in which we can override the stdin, stdout and stderr
       cmd_out_read, cmd_out_write = IO.pipe
       cmd_err_read, cmd_err_write = IO.pipe

       pid = fork {
         begin
           # grandchild
           drop_privileges(options[:uid], options[:gid], options[:supplementary_groups])

           Dir.chdir(ENV["PWD"] = options[:working_dir].to_s) if options[:working_dir]
           options[:environment].each { |key, value| ENV[key.to_s] = value.to_s } if options[:environment]

           # close unused fds so ancestors wont hang. This line is the only reason we are not
           # using something like popen3. If this fd is not closed, the .read call on the parent
           # will never return because "wr" would still be open in the "exec"-ed cmd
           wr.close

           # we do not care about stdin of cmd
           STDIN.reopen("/dev/null")

           # point stdout of cmd to somewhere we can read
           cmd_out_read.close
           STDOUT.reopen(cmd_out_write)
           cmd_out_write.close

           # same thing for stderr
           cmd_err_read.close
           STDERR.reopen(cmd_err_write)
           cmd_err_write.close

           # finally, replace grandchild with cmd
           ::Kernel.exec(*Shellwords.shellwords(cmd))
         rescue Exception => e
           (cmd_err_write.closed? ? STDERR : cmd_err_write).puts "Exception in grandchild: #{e.to_s}."
           (cmd_err_write.closed? ? STDERR : cmd_err_write).puts e.backtrace
           exit 1
         end
       }

       # we do not use these ends of the pipes in the child
       cmd_out_write.close
       cmd_err_write.close

       # wait for the cmd to finish executing and acknowledge it's death
       ::Process.waitpid(pid)

       # collect stdout, stderr and exitcode
       result = {
         :stdout => cmd_out_read.read,
         :stderr => cmd_err_read.read,
         :exit_code => $?.exitstatus
       }

       # We're done with these ends of the pipes as well
       cmd_out_read.close
       cmd_err_read.close

       # Time to tell the parent about what went down
       wr.write Marshal.dump(result)
       wr.close

       ::Process.exit!
     end
   end
*/

func IsPidAlive(pid int) bool {
	res := false
	process, err := os.FindProcess(int(pid))
	if err != nil {
		log.Printf("Failed to find process: %s\n", err)
	} else {
		err := process.Signal(syscall.Signal(0))
		log.Printf("process.Signal on pid %d returned: %v\n", pid, err)
		if err == nil {
			res = true
		}
	}
	return res
}

func ResetData() {
	if len(Store) > 0 {
		Store = make([]map[string]interface{}, 0)
	}
}

func IsDirectory(name string) bool {
	s, err := os.Stat(name)
	return err == nil && s.IsDir()
}

func ReadLines(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	var response []string
	bf := bufio.NewReader(f)

	for {
		line, isPrefix, err := bf.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if isPrefix {
			log.Fatal("Error: Unexpected long line reading", f.Name())
		}

		//fmt.Println(string(line))
		response = append(response, string(line))
	}
	return response
}
