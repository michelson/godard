package process

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	system "system"
	//logger "godard_logger"
	"time"
)

var Logger *log.Logger
var JournalBaseDir string

func SetLogger(new_logger *log.Logger) {
	if Logger == nil {
		Logger = new_logger
	}
}

func SetBaseDir(base_dir string) {

	if len(JournalBaseDir) == 0 {
		JournalBaseDir = path.Join(base_dir, "journals")
	}
	exists, _ := system.FileExists(JournalBaseDir)
	if !exists {
		err := os.MkdirAll(JournalBaseDir, 0777)
		if err != nil {
			log.Println("ERROR CREATING PIDS DIR", err)
		}
	}

	os.Chmod(JournalBaseDir, 0777)
}

func isSkipPid(pid int) bool {
	//!pid.is_a?(Integer) || pid <= 1
	return pid <= 1
}

func isSkipPgid(pid int) bool {
	//!pgid.is_a?(Integer) || pgid <= 1
	return pid <= 1
}

func AcquireAtomicFsLock(name string) {
	times := 0
	name += ".lock"
	err := os.MkdirAll(name, 0700)
	//logger.Debug("Acquired lock #{name}")
	//yield
	if err != nil {
		log.Println("ERROR CREATING PIDS DIR", err)
	}

	times += 1
	//logger.Debug("Waiting for lock #{name}")
	time.Sleep(1 * time.Second)
	if !(times >= 10) {
		// retry
	} else {
		//logger.Info("Timeout waiting for lock #{name}")
		//raise "Timeout waiting for lock #{name}"
	}
	//ensure
	ClearAtomicFsLock(name)
}

func ClearAllAtomicFsLocks() {
	files, _ := filepath.Glob(".*.lock")
	for _, file := range files {
		if system.IsDirectory(file) {
			system.DeleteIfExists(file)
		}
	}
}

func PidJournalFilename(journal_name string) string {
	return path.Join(JournalBaseDir, ".godard_pids_journal."+journal_name)
}

func PgidJournalFilename(journal_name string) string {
	return path.Join(JournalBaseDir, ".godard_pids_journal."+journal_name)
}

func PidJournal(filename string) []int {
	//logger.Debug("pid journal file: #{filename}")
	dat := system.ReadLines(filename)
	var arr []int
	for _, d := range dat {
		i, _ := strconv.Atoi(d)
		arr = append(arr, i)
	}
	for i, pid := range arr {
		if isSkipPid(pid) {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	//logger.Debug("pid journal = #{result.join(' ')}")
	//rescue Errno::ENOENT
	//[]
	return arr
}

func PgidJournal(filename string) []int {
	//logger.Debug("pgid journal file: #{filename}")
	dat := system.ReadLines(filename)
	var arr []int
	for _, d := range dat {
		i, _ := strconv.Atoi(d)
		arr = append(arr, i)
	}
	for i, pgid := range arr {
		if isSkipPgid(pgid) {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	//logger.Debug("pgid journal = #{result.join(' ')}")
	//rescue Errno::ENOENT
	//[]
	return arr
}

func ClearAtomicFsLock(name string) {
	if system.IsDirectory(name) {
		os.Remove(name)
		//logger.Debug("Cleared lock #{name}")
	}
}

func KillAllFromAllJournals() {

	files, _ := filepath.Glob(".godard_pids_journal.*")
	var xx []string
	var yy []string
	for _, x := range files {
		xx = append(xx, strings.Replace(x, ".bluepill_pids_journal.", "", 1))
	}
	for _, y := range xx {
		if !strings.Contains(y, ".lock") {
			yy = append(yy, y)
		}
	}
	for _, journal_name := range yy {
		KillAllFromJournal(journal_name)
	}
}

func KillAllFromJournal(journal_name string) {
	KillAllPidsFromJournal(journal_name)
	KillAllPgidsFromJournal(journal_name)
}

func KillAllPgidsFromJournal(journal_name string) {

	filename := PgidJournalFilename(journal_name)
	j := PgidJournal(filename)
	if len(j) > 0 {
		//AcquireAtomicFsLock(filename) do ??
		for _, pgid := range j {
			err := syscall.Kill(-pgid, syscall.SIGTERM)
			//logger.Info("Termed old process group #{pgid}")
			if err != nil {
				//logger.Debug("Unable to term missing process group #{pgid}")
			}
		}
		arr := make([]int, 0)
		for _, pgid := range j {
			if system.IsPidAlive(pgid) {
				arr = append(arr, pgid)
			}
		}
		if len(arr) > 1 {
			time.Sleep(1 * time.Second)
			for _, pgid := range j {
				err := syscall.Kill(-pgid, syscall.SIGTERM)
				//logger.Info("Killed old process group #{pgid}")
				if err != nil {
					//logger.Debug("Unable to kill missing process group #{pgid}")
				}
			}
			system.DeleteIfExists(filename) // reset journal
			//logger.Debug('Journal cleanup completed')
		}

		//end
	} else {
		//logger.Debug('No previous process journal - Skipping cleanup')
	}
}

func KillAllPidsFromJournal(journal_name string) {
	filename := PgidJournalFilename(journal_name)
	j := PgidJournal(filename)
	if len(j) > 0 {
		//acquire_atomic_fs_lock(filename) do
		for _, pid := range j {
			err := syscall.Kill(pid, syscall.SIGTERM)
			//logger.Info("Termed old process group #{pid}")
			if err != nil {
				//logger.Debug("Unable to term missing process group #{pid}")
			}
		}

		arr := make([]int, 0)
		for _, pid := range j {
			if system.IsPidAlive(pid) {
				arr = append(arr, pid)
			}
		}
		if len(arr) > 1 {
			time.Sleep(1 * time.Second)
			for _, pid := range j {
				err := syscall.Kill(pid, syscall.SIGTERM)
				//logger.Info("Killed old process group #{pid}")
				if err != nil {
					//logger.Debug("Unable to kill missing process group #{pid}")
				}
			}
			system.DeleteIfExists(filename) // reset journal
			//logger.Debug('Journal cleanup completed')
		}
		//end
	} else {
		//logger.Debug('No previous process journal - Skipping cleanup')
	}
}

func AppendPgidToJournal(journal_name string, pgid int) {

	if isSkipPgid(pgid) {
		//logger.Debug("Skipping invalid pgid #{pgid} (class #{pgid.class})")
		//return
	}
	filename := PgidJournalFilename(journal_name)
	//acquire_atomic_fs_lock(filename) do
	count := 0
	for _, p := range PgidJournal(filename) {
		if p == pgid {
			count += 1
		}
	}
	if count > 0 {
		//logger.Debug("Saving pgid #{pgid} to process journal #{journal_name}")
		d1 := []byte(strconv.Itoa(pgid))
		err := ioutil.WriteFile(filename, d1, 0600)
		if err == nil {
			//logger.Info("Saved pgid #{pgid} to journal #{journal_name}")
			//logger.Debug("Journal now = #{File.open(filename, 'r').read}")
		}
	} else {
		//logger.Debug("Skipping duplicate pgid #{pgid} already in journal #{journal_name}")
	}
}

func AppendPidToJournal(journal_name string, pid int) {

	pgid, err := syscall.Getpgid(pid)
	if err != nil {

	} else {
		AppendPgidToJournal(journal_name, pgid)
	}

	if isSkipPid(pid) {
		//logger.Debug("Skipping invalid pid #{pid} (class #{pid.class})")
		return
	}
	filename := PidJournalFilename(journal_name)
	//acquire_atomic_fs_lock(filename) do
	count := 0
	for _, p := range PidJournal(filename) {
		if p == pgid {
			count += 1
		}
	}
	if count > 0 {
		//logger.Debug("Saving pid #{pid} to process journal #{journal_name}")
		d1 := []byte(strconv.Itoa(pid))
		err := ioutil.WriteFile(filename, d1, 0600)
		if err == nil {
			//logger.Info("Saved pid #{pid} to journal #{journal_name}")
			//logger.Debug("Journal now = #{File.open(filename, 'r').read}")
		}
	} else {
		//logger.Debug("Skipping duplicate pid #{pid} already in journal #{journal_name}")
	}
}
