package condition

import (
	"fmt"
	"log"
	system "system"
)

//const mb int64          = 1024 ** 2
const format_str string = "%d%s"
const mb_label string = "MB"
const kb_label string = "KB"

type CpuUsage struct {
	Below float64
	Logger *log.Logger
}

func NewCpuUsage(options map[string]interface{}) *CpuUsage {
	var below float64      = float64(options["below"].(float64))
	var logger *log.Logger = options["logger"].(*log.Logger)
	c := &CpuUsage{Below: below, Logger: logger}
	c.Logger.Println("CREATING PROCESS CONDITION BELOW", c.Below)
	return c
}

func (c *CpuUsage) Run(pid int, include_children bool) (float64, error) {
	val, err := system.CpuUsage(pid) //, include_children)
	c.Logger.Println("CPU USAGE:", val)
	return val, err
}

func (c *CpuUsage) Check(value float64, include_children bool) (bool, error) {
	assert := value < c.Below
	return assert, nil
}

func (c *CpuUsage) FormatValue(value float64) string {
	var int_val int
	int_val = int(value)
	out := string("")
	out = fmt.Sprintf("%d%s", int_val, "%")
	return out
}
