package cgroups

import (
	"reflect"
	"testing"
)

func TestPids(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	h := "/user.slice"

	actualLabel, actualValue := c.Pids(h)
	expectedLabel := []string{
		"pids.max",
		"pids.current",
	}
	expectedValue := []string{
		"max",
		"24",
	}

	if !reflect.DeepEqual(actualLabel, expectedLabel) {
		t.Errorf("actual %v\nwant %v", actualLabel, expectedLabel)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Errorf("actual %v\nwant %v", actualValue, expectedValue)
	}
}
