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
  fmt.Println("Godard is a graph store and graph query layer.\n")
  fmt.Println("Usage:")
  fmt.Println("  godard COMMAND [flags]\n")
  fmt.Println("Commands:")
  fmt.Println("  init\tCreate an empty database.")
  //fmt.Println("  http\tServe an HTTP endpoint on the given host and port.")
  //fmt.Println("  repl\tDrop into a REPL of the given query language.")
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
  //var ts graph.TripleStore
  config := cfg.ParseConfigFromFlagsAndFile(*configFile)
  if os.Getenv("GOMAXPROCS") == "" {
    runtime.GOMAXPROCS(runtime.NumCPU())
  } else {
  }
  switch cmd {
  case "init":
    godard.Init(config)
  /*case "load":
    ts = godard.OpenTSFromConfig(config)
    godard.GodardLoad(ts, config, *tripleFile, false)
    ts.Close()
  case "repl":
    ts = godard.OpenTSFromConfig(config)
    godard.GodardRepl(ts, *queryLanguage, config)
    ts.Close()
  case "http":
    ts = godard.OpenTSFromConfig(config)
    godard_http.GodardHTTP(ts, config)
    ts.Close()*/
  default:
    fmt.Println("No command", cmd)
    flag.Usage()
  }
}