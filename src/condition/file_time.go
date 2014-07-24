package condition

import (
	"log"
	"os"
	"reflect"
	"time"
	util "util"
)

type FileTime struct {
	Below    float64
	filename string
	Logger   *log.Logger
}

func NewFileTime(options map[string]interface{}) *FileTime {
	var below float64
	var file string
	var logger *log.Logger = options["logger"].(*log.Logger)
	file = options["filename"].(string)
	type_of_value := reflect.TypeOf(options["below"])

	c := &FileTime{filename: file, Logger: logger}

	switch type_of_value.Kind() {
	case reflect.String:
		v, err := util.TimeParse(options["below"].(string))
		if err != nil {
			logger.Println("error while parsing below options", options["below"])
		}
		below = v.Seconds()
		c.Below = below
	case reflect.Int:
		below = options["below"].(float64)
		c.Below = below
	}

	return c
}

func (c *FileTime) Run(pid int, include_children bool) (float64, error) {
	info, err := os.Stat(c.filename)
	if err != nil {
		// TODO: handle errors (e.g. file not found)
		return 0, err
	} else {
		var d float64
		duration := time.Now().UnixNano() - info.ModTime().UnixNano()
		d = float64(duration)
		return d, nil
	}
}

func (c *FileTime) Check(value float64, include_children bool) (bool, error) {
	//  value.kilobytes < c.Below
	assert := value < c.Below
	return assert, nil
}

/*

   def initialize(options = {})
     @below = options[:below]
     @filename = options[:filename]
   end

   def run(pid, include_children)
     if File.exists?(@filename)
       Time.now()-File::mtime(@filename)
     else
       nil
     end
   rescue
     $!
   end

   def check(value)
     return false if value.nil?
     return value < @below
   end

*/
