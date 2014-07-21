package process


import (
  system "system"
  "path"
  "path/filepath"
  "os"
  "syscall"
  "log"
  "io/ioutil"
  "strconv"
)

type ProcessJournal struct {
  Logger string
  JournalBaseDir string
}


func (c*ProcessJournal) SetLogger(new_logger string) {
  if len(c.Logger) == 0 {
    c.Logger = new_logger
  } 
}

func (c*ProcessJournal) SetBaseDir(base_dir string){

  if len(c.JournalBaseDir) == 0 {
    c.JournalBaseDir = path.Join(base_dir , "journals")
  } 
  exists , _ := system.FileExists(c.JournalBaseDir)
  if !exists{
    err := os.MkdirAll(c.JournalBaseDir, 0777)
    if err != nil {
      log.Println("ERROR CREATING PIDS DIR" , err)
    }
  }

  os.Chmod(c.JournalBaseDir, 0777)
}

func (c*ProcessJournal) isSkipPid(pid int) bool{
  //!pid.is_a?(Integer) || pid <= 1
  return pid <= 1
}

func (c*ProcessJournal) isSkipPgid(pid int) bool{
  //!pgid.is_a?(Integer) || pgid <= 1
  return pid <= 1
}

func (c*ProcessJournal) AcquireAtomicFsLock(name string){
  times := 0
  name += ".lock"
  err := os.MkdirAll(name, 0700)
  //logger.debug("Acquired lock #{name}")
  //yield
  if err != nil {
    log.Println("ERROR CREATING PIDS DIR" , err)
  }

  times += 1
  //logger.debug("Waiting for lock #{name}")
  //sleep time.Second * 1
  if !(times >= 10) {
    // retry
  } else {
    //logger.info("Timeout waiting for lock #{name}")
    //raise "Timeout waiting for lock #{name}"
  }
  //ensure
  c.ClearAtomicFsLock(name)
}

func (c*ProcessJournal) ClearAllAtomicFsLocks(){
  files, _ := filepath.Glob(".*.lock")
  for _ , file := range(files){
    if system.IsDirectory(file){
      system.DeleteIfExists(file)
    }
  }
}

func (c*ProcessJournal) PidJournalFilename(journal_name string) string{
  return path.Join(c.JournalBaseDir , ".godard_pids_journal" + journal_name)
}

func (c*ProcessJournal) PgidJournalFilename(journal_name string) string{
  return path.Join(c.JournalBaseDir , ".godard_pids_journal" + journal_name)
}

func (c*ProcessJournal) PidJournal(filename string) []int {
   //logger.debug("pid journal file: #{filename}")
   dat := system.ReadLines(filename)
   var arr []int
   for _ , d := range(dat) {
      i , _ := strconv.Atoi(d)
      arr = append(arr , i)
   }
   for i , pid := range(arr){
      if c.isSkipPid(pid){
        arr = append(arr[:i], arr[i+1:]...)
      }
   }
   //logger.debug("pid journal = #{result.join(' ')}")
   //rescue Errno::ENOENT
   //[]
  return arr

}

func (c*ProcessJournal) PgidJournal(filename string) []int {
   //logger.debug("pgid journal file: #{filename}")
   dat := system.ReadLines(filename)
   var arr []int
   for _ , d := range(dat) {
      i , _ := strconv.Atoi(d)
      arr = append(arr , i)
   }
   for i , pgid := range(arr){
      if c.isSkipPgid(pgid){
        arr = append(arr[:i], arr[i+1:]...)
      }
   }
   //logger.debug("pgid journal = #{result.join(' ')}")
   //rescue Errno::ENOENT
   //[]
  return arr
}

func (c*ProcessJournal) ClearAtomicFsLock(name string){
  if system.IsDirectory(name) {
    os.Remove(name)
    //logger.debug("Cleared lock #{name}")
  }
}

func (c*ProcessJournal) KillAllFromAllJournals(){
/*
      Dir[".bluepill_pids_journal.*"].map { |x|
        x.sub(/^\.bluepill_pids_journal\./,"")
      }.reject { |y|
        y =~ /\.lock$/
      }.each do |journal_name|
        kill_all_from_journal(journal_name)
      end
*/
}

func (c*ProcessJournal) KillAllFromJournal(journal_name string){
  c.KillAllPidsFromJournal(journal_name)
  c.KillAllPgidsFromJournal(journal_name)
}

func (c*ProcessJournal) KillAllPgidsFromJournal(journal_name string){

  filename := c.PgidJournalFilename(journal_name)
  j := c.PgidJournal(filename)
  if len(j) > 0 {
    //c.AcquireAtomicFsLock(filename) do ??
    for _, pgid := range(j){
      err := syscall.Kill(-pgid, syscall.SIGTERM)
      //logger.info("Termed old process group #{pgid}")
      if err != nil {
        //logger.debug("Unable to term missing process group #{pgid}")
      }
    }
    arr := make([]int, 0)
    for _,pgid := range(j){
      if system.IsPidAlive(pgid){
        arr = append(arr , pgid)
      }
    }
    if len(arr) > 1 {
      //sleep 1
      for _ , pgid := range(j) {
        err := syscall.Kill(-pgid, syscall.SIGTERM)
          //logger.info("Killed old process group #{pgid}")
        if err != nil {
          //logger.debug("Unable to kill missing process group #{pgid}")
        }
      }
      system.DeleteIfExists(filename) // reset journal
      //logger.debug('Journal cleanup completed')
    }

    //end
  }else{
    //logger.debug('No previous process journal - Skipping cleanup')
  }
}

func (c*ProcessJournal) KillAllPidsFromJournal(journal_name string){
  filename := c.PgidJournalFilename(journal_name)
  j := c.PgidJournal(filename)
  if len(j) > 0 {
    //acquire_atomic_fs_lock(filename) do
    for _, pid := range(j){
      err := syscall.Kill(pid, syscall.SIGTERM)
      //logger.info("Termed old process group #{pid}")
      if err != nil {
        //logger.debug("Unable to term missing process group #{pid}")
      }
    }

    arr := make([]int, 0)
    for _,pid := range(j){
      if system.IsPidAlive(pid){
        arr = append(arr , pid)
      }
    }
    if len(arr) > 1 {
      //sleep 1
      for _ , pid := range(j) {
        err := syscall.Kill(pid, syscall.SIGTERM)
        //logger.info("Killed old process group #{pid}")
        if err != nil {
          //logger.debug("Unable to kill missing process group #{pid}")
        }
      }
      system.DeleteIfExists(filename) // reset journal
      //logger.debug('Journal cleanup completed')
    }
    //end
  }else{
    //logger.debug('No previous process journal - Skipping cleanup')
  }
}

func (c*ProcessJournal) AppendPgidToJournal(journal_name string, pgid int){

  if c.isSkipPgid(pgid){
    //logger.debug("Skipping invalid pgid #{pgid} (class #{pgid.class})")
    //return
  }
  filename := c.PgidJournalFilename(journal_name)
  //acquire_atomic_fs_lock(filename) do
  count := 0
  for _ , p := range(c.PgidJournal(filename)){
    if p == pgid {
      count += 1
    }
  }
  if count > 0 {   
    //logger.debug("Saving pgid #{pgid} to process journal #{journal_name}")
    //File.open(filename, 'a+', 0600) { |f| f.puts(pgid) }
    d1 := []byte(strconv.Itoa(pgid))
    err := ioutil.WriteFile(filename, d1, 0600)
    if err == nil{
      //logger.info("Saved pgid #{pgid} to journal #{journal_name}")
      //logger.debug("Journal now = #{File.open(filename, 'r').read}")
    }
  }else{
    //logger.debug("Skipping duplicate pgid #{pgid} already in journal #{journal_name}")
  }

}

func (c*ProcessJournal) AppendPidToJournal(journal_name string, pid int){

  pgid, err := syscall.Getpgid(pid)
  if err != nil {
    
  }else{
    c.AppendPgidToJournal(journal_name , pgid  )    
  }

  if c.isSkipPid(pid){
    //logger.debug("Skipping invalid pid #{pid} (class #{pid.class})")
    return
  }
  filename := c.PidJournalFilename(journal_name)
  //acquire_atomic_fs_lock(filename) do
  count := 0
  for _ , p := range(c.PidJournal(filename)){
    if p == pgid {
      count += 1
    }
  }
  if count > 0 {
    //logger.debug("Saving pid #{pid} to process journal #{journal_name}")
    d1 := []byte(strconv.Itoa(pid))
    err := ioutil.WriteFile(filename, d1, 0600)
    if err == nil{
      //logger.info("Saved pid #{pid} to journal #{journal_name}")
      //logger.debug("Journal now = #{File.open(filename, 'r').read}")      
    }
  }else{
    //logger.debug("Skipping duplicate pid #{pid} already in journal #{journal_name}")
  }
}
