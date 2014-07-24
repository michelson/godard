package condition

/*
  MB := 1024 ** 2
  FORMAT_STR := "%d%s"
  MB_LABEL := "MB"
  KB_LABEL := "KB"
*/

import (
	"fmt"
	"log"
	"reflect"
	system "system"
	util "util"
)

type MemoryUsage struct {
	Below  float64
	Logger *log.Logger
}

func NewMemoryUsage(options map[string]interface{}) *MemoryUsage {
	var below float64
	var logger *log.Logger = options["logger"].(*log.Logger)
	type_of_value := reflect.TypeOf(options["below"])

	c := &MemoryUsage{Logger: logger}

	switch type_of_value.Kind() {
	case reflect.String:
		v, err := util.ParseNumber(options["below"].(string))
		if err != nil {
			logger.Println("error while parsing below options", options["below"])
		}
		below = v
		c.Below = below
	case reflect.Int:
		below = options["below"].(float64) * 1024 * 1024
		c.Below = below
	}

	c.Logger.Println("CREATING PROCESS CONDITION BELOW", options["below"])
	return c
}

func (c *MemoryUsage) Run(pid int, include_children bool) (float64, error) { // , include_children bool) {

	val, err := system.MemoryUsage(pid) //, include_children)
	c.Logger.Println("MEM USAGE:", val)
	var usage float64
	usage = float64(val)
	return usage, err
}

func (c *MemoryUsage) Check(value float64, include_children bool) (bool, error) {
	//  value.kilobytes < c.Below
	assert := (value * 1024 * 1024) < c.Below
	return assert, nil
}

func (c *MemoryUsage) FormatValue(value float64) string {
	var int_val int
	int_val = int(value)
	out := string("")
	//MB = 1024 ** 2
	if int_val*1024 >= 1048576 {
		out = fmt.Sprintf(format_str, (int_val / 1024), mb_label)
	} else {
		out = fmt.Sprintf(format_str, int_val, kb_label)
	}
	return out
}
