
package godard_config

import (
	"encoding/json"
	//"flag"
	//"github.com/barakmich/glog"
	"os"
	"log"
)

type GodardConfig struct {
	Processes       []map[string]interface{} `json:"processes"`
	ListenHost      string                   `json:"listen_host"`
	ListenPort      string                   `json:"listen_port"`	
	Foreground      bool										 `json:"foreground"`
	LogFile         string									 `json:"log_file"`	
	BaseDir         string									 `json:"base_dir"`	
	KillTimeout     int 										 `json:"kill_timeout"`	
}

func ParseConfigFromFile(filename string) *GodardConfig {
	config := &GodardConfig{}
	
	if filename == "" {
		return config
	}
	f, err := os.Open(filename)
	
	if err != nil {
		log.Fatalln("Couldn't open config file", filename)
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(config)
	
	if err != nil {
		log.Fatalln("Couldn't read config file:", err)
	}
	
	return config
}

func ParseConfigFromFlagsAndFile(fileFlag string) *GodardConfig {
	// Find the file...
	log.Println(fileFlag)
	
	var trueFilename string
	
	if fileFlag != "" {
		if _, err := os.Stat(fileFlag); os.IsNotExist(err) {
			log.Fatalln("Cannot find specified configuration file", fileFlag, ", aborting.")
		} else {
			trueFilename = fileFlag
		}
	} else {
		if _, err := os.Stat(os.Getenv("GODARD_CFG")); err == nil {
			trueFilename = os.Getenv("GODARD_CFG")
		} else {
			if _, err := os.Stat("/etc/godard.cfg"); err == nil {
				trueFilename = "/etc/godard.cfg"
			}
		}
	}
	
	if trueFilename == "" {
		log.Println("Couldn't find a config file in either $GODARD_CFG or /etc/godard.cfg. Going by flag defaults only.")
	}
	
	config := ParseConfigFromFile(trueFilename)

	return config
}
