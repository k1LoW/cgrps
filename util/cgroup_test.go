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
	expected := 66
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func testFs() string {
	dir, _ := os.Getwd()
	return dir + "/../test/sys/fs/cgroup"
}
