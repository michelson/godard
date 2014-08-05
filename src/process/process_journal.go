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
			Logger.Println("ERROR CREATING PIDS DIR", err)
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
	//Logger.Println("Acquired lock #{name}")
	//yield
	if err != nil {
		Logger.Println("ERROR CREATING PIDS DIR", err)
	}

	times += 1
	Logger.Println("Waiting for lock", name)
	time.Sleep(1 * time.Second)
	if !(times >= 10) {
		// retry
	} else {
		Logger.Println("Timeout waiting for lock", name)
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
	return path.Join(JournalBaseDir, ".godard_pgids_journal."+journal_name)
}

func PidJournal(filename string) []int {
	Logger.Println("pid journal file:", filename)

	var arr []int
	exists, _ := system.FileExists(filename)
	if exists {
		os.Open(filename)
	} else {
		return arr
	}

	dat := system.ReadLines(filename)

	for _, d := range dat {
		i, _ := strconv.Atoi(d)
		arr = append(arr, i)
	}
	for i, pid := range arr {
		if isSkipPid(pid) {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	Logger.Println("pid journal =", arr)
	//rescue Errno::ENOENT
	//[]
	return arr
}

func PgidJournal(filename string) []int {
	//Logger.Println("pgid journal file: #{filename}")
	var arr []int
	exists, _ := system.FileExists(filename)
	if exists {
		os.Open(filename)
	} else {
		return arr
	}

	dat := system.ReadLines(filename)
	for _, d := range dat {
		i, _ := strconv.Atoi(d)
		arr = append(arr, i)
	}
	for i, pgid := range arr {
		if isSkipPgid(pgid) {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	Logger.Println("pgid journal =", arr)
	//rescue Errno::ENOENT
	//[]
	return arr
}

func ClearAtomicFsLock(name string) {
	if system.IsDirectory(name) {
		os.Remove(name)
		Logger.Println("Cleared lock", name)
	}
}

func KillAllFromAllJournals() {
	filelocation := path.Join(JournalBaseDir, ".godard_pids_journal.")
	files, _ := filepath.Glob(filelocation + "*")
	var filenames []string
	var names []string
	for _, file := range files {
		filenames = append(filenames, strings.Replace(file, filelocation, "", 1))
	}
	for _, name := range filenames {
		if !strings.Contains(name, ".lock") {
			names = append(names, name)
		}
	}
	//log.Println("FILES TO KILL", names)
	for _, journal_name := range names {
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
			Logger.Println("Termed old process group", pgid)
			if err != nil {
				Logger.Println("Unable to term missing process group", pgid)
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
				Logger.Println("Killed old process group", pgid)
				if err != nil {
					Logger.Println("Unable to kill missing process group", pgid)
				}
			}

		}
		system.DeleteIfExists(filename) // reset journal
		Logger.Println("Journal cleanup completed")
		//end
	} else {
		Logger.Println("No previous process journal - Skipping cleanup")
	}
}

func KillAllPidsFromJournal(journal_name string) {
	filename := PidJournalFilename(journal_name)
	j := PgidJournal(filename)
	if len(j) > 0 {
		//acquire_atomic_fs_lock(filename) do
		for _, pid := range j {
			err := syscall.Kill(pid, syscall.SIGTERM)
			Logger.Println("Termed old process group", pid)
			if err != nil {
				Logger.Println("Unable to term missing process group", pid)
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
				Logger.Println("Killed old process group", pid)
				if err != nil {
					Logger.Println("Unable to kill missing process group", pid)
				}
			}

		}
		system.DeleteIfExists(filename) // reset journal
		Logger.Println("Journal cleanup completed")
		//end
	} else {
		Logger.Println("No previous process journal - Skipping cleanup")
	}
}

func AppendPgidToJournal(journal_name string, pgid int) {

	if isSkipPgid(pgid) {
		Logger.Println("Skipping invalid pgid", pgid)
		//return
	} else {
		filename := PgidJournalFilename(journal_name)
		//acquire_atomic_fs_lock(filename) do
		count := 0
		for _, p := range PgidJournal(filename) {
			if p == pgid {
				count += 1
			}
		}
		if count == 0 {
			Logger.Println("Saving pgid", pgid, " to process journal", journal_name, filename)
			d1 := strconv.Itoa(pgid)

			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			f.WriteString(d1 + "\n")
			f.Close()

			//err := ioutil.WriteFile(filename, d1, 0600)
			if err == nil {
				Logger.Println("Saved pgid", pgid, " to journal", journal_name)
				dat, _ := ioutil.ReadFile(filename)
				Logger.Println("Journal now =", string(dat))
			}

		} else {
			//Logger.Println("Skipping duplicate pgid" , pgid , " already in journal", journal_name )
		}
	}

}

func AppendPidToJournal(journal_name string, pid int) {

	if isSkipPid(pid) {
		//Logger.Println("Skipping invalid pid #{pid} (class #{pid.class})")
		return
	}
	pgid, err := syscall.Getpgid(pid)
	if err != nil {

	} else {
		//Logger.Println("PGID IS:", pgid, "FROM", pid)
		AppendPgidToJournal(journal_name, pgid)
	}

	filename := PidJournalFilename(journal_name)
	//Logger.Println("FILENAME JOURNAL", filename)

	//acquire_atomic_fs_lock(filename) do
	count := 0
	for _, p := range PidJournal(filename) {
		if p == pgid {
			count += 1
		}
	}
	if count == 0 {
		Logger.Println("Saving pid", pid, " to process journal", journal_name)
		d1 := strconv.Itoa(pid)
		//err := ioutil.WriteFile(filename, d1, 0600)

		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		f.WriteString(d1 + "\n")
		f.Close()

		if err == nil {
			Logger.Println("Saved pid", pid, " to journal ", journal_name)
			dat, _ := ioutil.ReadFile(filename)
			Logger.Println("Journal now =", string(dat))
		}

	} else {
		//Logger.Println("Skipping duplicate pid", pid, " already in journal", journal_name )
	}
}
