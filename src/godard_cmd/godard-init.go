package godard

import(
  //"log"
  cfg "godard_config"
  //system "system"
  app "dsl"
)

func Init(config *cfg.GodardConfig){

  /*cpu, _ := system.CpuUsage(28278)
  log.Println("CPU: ", cpu)*/

  /*for _, t := range(config.Processes) {
    log.Println("PROCESS CONFIG:", t["pid_file"])
  }*/

  app.InitApplication("myApp", config )

/*
  cpu, _ := system.CpuUsage(28278)
  log.Println("CPU: ", cpu)

  mem, _ := system.MemoryUsage(28278)
  log.Println("MEM: ", mem)

  running_time, _ := system.RunningTime(28278)
  log.Println("ELAPSED TIME: ", running_time)

  cmd, _ := system.Command(28278)
  log.Println("COMMAND: ", cmd)

  childs, _ := system.GetChildren(28278)
  log.Println("CHILD PROCESS: ", childs)
*/
  
}