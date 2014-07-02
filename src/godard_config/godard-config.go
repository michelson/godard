
package godard_config

import (
	"encoding/json"
	"flag"
	"github.com/barakmich/glog"
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
}


var host = flag.String("host", "0.0.0.0", "Host to listen on (defaults to all).")
var port = flag.String("port", "64210", "Port to listen on.")

func ParseConfigFromFile(filename string) *GodardConfig {
	config := &GodardConfig{}
	
	if filename == "" {
		return config
	}
	f, err := os.Open(filename)
	
	if err != nil {
		glog.Fatalln("Couldn't open config file", filename)
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(config)
	
	if err != nil {
		glog.Fatalln("Couldn't read config file:", err)
	}
	
	return config
}

func ParseConfigFromFlagsAndFile(fileFlag string) *GodardConfig {
	// Find the file...
	log.Println(fileFlag)
	
	var trueFilename string
	
	if fileFlag != "" {
		if _, err := os.Stat(fileFlag); os.IsNotExist(err) {
			glog.Fatalln("Cannot find specified configuration file", fileFlag, ", aborting.")
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
		glog.Infoln("Couldn't find a config file in either $GODARD_CFG or /etc/godard.cfg. Going by flag defaults only.")
	}
	
	config := ParseConfigFromFile(trueFilename)

	if config.ListenHost == "" {
		config.ListenHost = *host
	}

	if config.ListenPort == "" {
		config.ListenPort = *port
	}

	//config.ReadOnly = config.ReadOnly || *readOnly

	return config
}
