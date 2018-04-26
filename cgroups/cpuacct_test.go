package cgroups

import (
	"reflect"
	"testing"
)

func TestCPUAcct(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	h := "/"

	actualLabel, actualValue := c.CPUAcct(h)
	expectedLabel := []string{
		"cpuacct.usage",
		"cpuacct.stat.user",
		"cpuacct.stat.system",
	}
	expectedValue := []string{
		"11063909873548",
		"912048",
		"99830",
	}

	if !reflect.DeepEqual(actualLabel, expectedLabel) {
		t.Errorf("actual %v\nwant %v", actualLabel, expectedLabel)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Errorf("actual %v\nwant %v", actualValue, expectedValue)
	}
}
