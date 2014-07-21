package condition

import (
	"testing"
)

func TestCpuFormat(t *testing.T) {
	c := &CpuUsage{}
	format := c.FormatValue(100)
	if format != "100%" {
		t.Error("Expected 100%, got ", format)
	}
}
