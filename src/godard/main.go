package main

import (
  godard "godard_cmd"
  cfg "godard_config"
  "flag"
  "fmt"
  //"github.com/barakmich/glog"
  //"graph"
  "os"
  "runtime"
)


var configFile = flag.String("config", "", "Path to an explicit configuration file.")

func Usage() {
  fmt.Println("Godard is monitoring tool.\n")
  fmt.Println("Usage:")
  fmt.Println("  godard COMMAND [flags]\n")
  fmt.Println("Commands:")
  fmt.Println("  init\tStart server , load config.")
  fmt.Println("\nFlags:")
  flag.Parse()
  flag.PrintDefaults()
}

func main() {
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

  config := cfg.ParseConfigFromFlagsAndFile(*configFile)
  
  if os.Getenv("GOMAXPROCS") == "" {
    runtime.GOMAXPROCS(runtime.NumCPU())
  }

  switch cmd {
    case "init":
      godard.Init(config)
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