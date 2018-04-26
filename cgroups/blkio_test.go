package cgroups

import (
	"reflect"
	"testing"
)

func TestBlkio(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	h := "/"

	actualLabel, actualValue := c.Blkio(h)
	expectedLabel := []string{
		"blkio.weight",
	}
	expectedValue := []string{
		"1000",
	}

	if !reflect.DeepEqual(actualLabel, expectedLabel) {
		t.Errorf("actual %v\nwant %v", actualLabel, expectedLabel)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Errorf("actual %v\nwant %v", actualValue, expectedValue)
	}
}
