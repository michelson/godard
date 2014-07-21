package util

import (
	"testing"
	"time"
)

func TestTimeParse(t *testing.T) {
	v, _ := TimeParse("1.seconds")
	if v != time.Second*1 {
		t.Error("Expected 1 sec, got ", v)
	}

	v2, _ := TimeParse("1.hour")
	if v2 != time.Hour*1 {
		t.Error("Expected 1 hour, got ", v2)
	}

	v3, _ := TimeParse("1.minutes")
	if v3 != time.Minute*1 {
		t.Error("Expected 1 min, got ", v2)
	}

	v4, _ := TimeParse("1.5.minutes")
	f := float64(1.5)
	if v4 != time.Minute*time.Duration(f) {
		t.Error("Expected 1 min, got ", v4)
	}

}

func TestFailingTimeParse(t *testing.T) {
	_, err := TimeParse("ccc.minutes")
	if err == nil {
		t.Error("Expected err, got ", err)
	}

	_, err2 := TimeParse("12.sparks")
	if err2 == nil {
		t.Error("Expected err, got ", err2)
	}

	_, err3 := TimeParse("1,minutes")
	if err3 == nil {
		t.Error("Expected err, got ", err3)
	}
}
