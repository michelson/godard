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
	system "system"
)

type MemoryUsage struct {
	Below float64
}

func NewMemoryUsage(options map[string]interface{}) *MemoryUsage {
	var below float64
	below = float64(options["below"].(float64))
	c := &MemoryUsage{Below: below}
	log.Println("CREATING PROCESS CONDITION BELOW", c.Below)
	return c
}

func (c *MemoryUsage) Run(pid int, include_children bool) (float64, error) { // , include_children bool) {

	val, err := system.MemoryUsage(pid) //, include_children)
	log.Println("MEM USAGE:", val)
	var usage float64
	usage = float64(val)
	return usage, err
}

func (c *MemoryUsage) Check(value float64, include_children bool) (bool, error) {
	//  value.kilobytes < c.Below
	assert := value < c.Below
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
