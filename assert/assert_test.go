package assert

import (
	"os"
	"testing"
)

func TestAssert(t *testing.T) {
	True(true, "This is true")
	if TotalPanics != 0 {
		t.Errorf("TotalPanics should be 0")
	}

	var before = os.Getenv(TaskNinjaSkipAssert)
	os.Setenv(TaskNinjaSkipAssert, "true")
	True(false, "This is false")
	if TotalPanics != 1 {
		t.Errorf("TotalPanics should be 1")
	}
	os.Setenv(TaskNinjaSkipAssert, before)
}
