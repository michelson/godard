package condition

/* MB := 1024 ** 2
FORMAT_STR := "%d%s"
MB_LABEL := "MB"
KB_LABEL := "KB"

*/

import (
	"log"
	"reflect"
	system "system"
	util "util"
)

type RunningTime struct {
	Below  float64
	Logger *log.Logger
}

func NewRunningTime(options map[string]interface{}) *RunningTime {
	var below float64
	var logger *log.Logger = options["logger"].(*log.Logger)
	type_of_value := reflect.TypeOf(options["below"])

	c := &RunningTime{Logger: logger}

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

func (c *RunningTime) Run(pid int, include_children bool) (int, error) {
	return system.RunningTime(pid)
}

func (c *RunningTime) Check(value float64) (bool, error) {
	assert := value < c.Below
	return assert, nil
}
