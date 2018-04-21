package util

import (
	"os"
	"testing"
)

func TestList(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	cs, err := c.List()
	if err != nil {
		t.Error(err)
	}
	actual := len(cs)
	expected := 73
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func TestEnabledSubsystems(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	subsystems := c.EnabledSubsystems("/my-cgroup")
	actual := len(subsystems)
	expected := 4
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func TestIsEnabledSubsystem(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	actual := c.IsEnableSubsystem("/my-cgroup", "cpu,cpuacct")
	expected := true
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
	actual = c.IsEnableSubsystem("/my-cgroup", "devices")
	expected = false
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func TestProcesses(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	ps, err := c.Processes("/my-cgroup")
	if err != nil {
		t.Error(err)
	}
	actual := len(ps)
	expected := 3
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func testFs() string {
	dir, _ := os.Getwd()
	return dir + "/../test/sys/fs/cgroup"
}
