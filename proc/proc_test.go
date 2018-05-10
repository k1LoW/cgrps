package proc

import (
	"os"
	"testing"
)

func TestCgroup(t *testing.T) {
	p := Proc{FsPath: testFs()}
	pids := []string{"12345"}
	cs, err := p.Cgroup(pids)
	if err != nil {
		t.Error(err)
	}
	actual := len(cs)
	expected := 5
	if actual != expected {
		t.Errorf("actual %v\nwant %v", actual, expected)
	}
}

func testFs() string {
	dir, _ := os.Getwd()
	return dir + "/../test/proc"
}
