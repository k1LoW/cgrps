package cgroups

import (
	"reflect"
	"testing"
)

func TestCPU(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	h := "/"

	actualLabel, actualValue := c.CPU(h)
	expectedLabel := []string{
		"cpu.cfs_period_us",
		"cpu.cfs_quota_us",
		"cpu.shares",
		"cpu.stat.nr_periods",
		"cpu.stat.nr_throttled",
		"cpu.stat.throttled_time",
	}
	expectedValue := []string{
		"100000",
		"-1",
		"1024",
		"0",
		"0",
		"0",
	}

	if !reflect.DeepEqual(actualLabel, expectedLabel) {
		t.Errorf("actual %v\nwant %v", actualLabel, expectedLabel)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Errorf("actual %v\nwant %v", actualValue, expectedValue)
	}
}
