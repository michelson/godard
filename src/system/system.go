//package process

package system

import (
    "bytes"
    "log"
    "os/exec"
    "strconv"
    "strings"
    "regexp"
)

func PsAxu() ([]map[string]interface{} , error) {
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
      if err!=nil {
          break;
      }
      tokens := strings.Split(line, " ")
      //log.Println(tokens)

      ft := make([]string, 0)
      for _, t := range(tokens) {
          if t!="" && t!="\t" {
              ft = append(ft, t)
          }
      }

      pid, err := strconv.Atoi(ft[0])
      if err!=nil {
          continue
      }
      ppid, err := strconv.Atoi(ft[1])
      if err!=nil {
          continue
      }
      cpu, err := strconv.ParseFloat(ft[2], 64)
      if err!=nil {
          log.Fatal(err)
      }
      mem, err := strconv.Atoi(ft[3])
      if err!=nil {
          continue
      }


      d := ParseElapsedTime(string(ft[4]))

      m := map[string]interface{}{
          "pid": pid,
          "ppid":  ppid,
          "cpu":  cpu,
          "mem":  mem,
          "elapsed":  d,
          "command":  ft[5],
      }

      processes = append(processes, m )
    }

    return processes , err
}

func FindByPid(pid int , processes []map[string]interface{}) (map[string]interface{}, error) {
  var pp map[string]interface{}

  for _, p := range(processes) {
    if p["pid"].(int) == pid {
      pp = p
      break
    }
  }
  return pp, nil
}


func CpuUsage(pid int) (float64, error) {
  ps , err := PsAxu()
  process , err := FindByPid(pid, ps)
  //log.Println("PROCESS:", process["pid"].(int))
  return process["cpu"].(float64), err
}


func MemoryUsage(pid int) (int, error) {
  ps , err := PsAxu()
  process , err := FindByPid(pid, ps)
  return process["mem"].(int)/1024, err
}

func RunningTime(pid int) (int, error){
  ps , err := PsAxu()
  process , err := FindByPid(pid, ps)
  return process["elapsed"].(int), err
}

func Command(pid int) (string, error){
  ps , err := PsAxu()
  process , err := FindByPid(pid, ps)
  return process["command"].(string), err
}

func GetChildren(parent_pid int) ([]map[string]interface{}, error){
  child_pids := make([]map[string]interface{}, 0)

  processes , err := PsAxu()

  for _, p := range(processes) {
    if p["ppid"].(int) == parent_pid {
      child_pids = append(child_pids, p)
    }
  }

  for _, p := range(child_pids) {
    gcp, _ := GetChildren(p["pid"].(int))
    for _, ppp := range(gcp) {
      child_pids = append(child_pids, ppp)
    }
  }

  return child_pids , err
} 

func Daemonize(cmd string, opts map[string]string){}


func ParseElapsedTime(str string) int {
  // [[dd-]hh:]mm:ss
    //s := "02-09:38:38"
    r, _ := regexp.Compile(`(?:(?:(\d+)-)?(\d\d):)?(\d\d):(\d\d)`)
    res:= r.FindStringSubmatch(str)

    days , _  := strconv.ParseFloat(res[1], 64)
    hours , _ := strconv.ParseFloat(res[2], 64)
    mins , _  := strconv.ParseFloat(res[3], 64)
    secs , _  := strconv.ParseFloat(res[4], 64)
    //var tt time.Duration = (( (days*24 + hours)*60 + mins)*60 + secs )* time.Second 
    //log.Println( days , hours,mins, secs  )
   
    minutes := ( ( (days*24 + hours)*60 + mins)*60 + secs ) / 60

    return int(minutes)
}

