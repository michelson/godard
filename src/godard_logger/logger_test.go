package godard_logger

import (
	"testing"
	//"reflect"
	//"log"
)

func TestNewGodardLogger(t *testing.T) {
	var options = make(map[string]interface{})
	L := NewGodardLogger(options)

	if L.Logger != nil {
		t.Error("Expected logger, got ", L.Logger)
	}

}
