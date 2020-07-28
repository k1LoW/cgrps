package cgroups

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

func TestListPids(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	pids := c.ListPids([]string{""})
	actual := len(pids)
	expected := 0
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func TestAttachedSubsystems(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	subsystems := c.AttachedSubsystems("/my-cgroup")
	actual := len(subsystems)
	expected := 4
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func TestIsAttacheddSubsystem(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	actual := c.IsAttachedSubsystem("/my-cgroup", "cpuacct")
	expected := true
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
	actual = c.IsAttachedSubsystem("/my-cgroup", "devices")
	expected = false
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func TestReadSimple(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	actual, err := c.ReadSimple("/my-cgroup", "memory", "memory.usage_in_bytes")
	if err != nil {
		t.Error(err)
	}
	expected := "1234567890"
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func testFs() string {
	dir, _ := os.Getwd()
	return dir + "/../test/sys/fs/cgroup"
}
