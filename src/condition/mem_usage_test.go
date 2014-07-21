package condition

import (
	"testing"
)

func TestMemFormat(t *testing.T) {
	c := &MemoryUsage{}
	format := c.FormatValue(100)
	if format != "100KB" {
		t.Error("Expected 100KB, got ", format)
	}
}
