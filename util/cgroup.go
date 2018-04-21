// Copyright © 2018 Ken'ichiro Oyama <k1lowxb@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package util

import (
	"bufio"
	"fmt"
	"github.com/k1LoW/go-ps"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Subsystems cgroups subsystems list
var Subsystems = []string{
	"cpuset",
	"cpu",
	"cpuacct",
	"blkio",
	"memory",
	"devices",
	"freezer",
	"net_cls",
	"net_prio",
	"perf_event",
	"hugetlb",
	"pids",
	"rdma",
}

func ClkTck() float64 {
	tck := float64(128)
	out, err := exec.Command("/usr/bin/getconf", "CLK_TCK").Output()
	if err == nil {
		i, err := strconv.ParseFloat(string(out), 64)
		if err == nil {
			tck = float64(i)
		}
	}
	return tck
}

type Cgroups struct {
	FsPath string
}

func (c *Cgroups) List() ([]string, error) {
	subsys := c.EnabledSubsystems("/")

	cs := []string{"/"}
	encountered := make(map[string]bool)

	for _, s := range subsys {
		searchDir := filepath.Clean(fmt.Sprintf("%s/%s", c.FsPath, s))

		err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				c := strings.Replace(path, searchDir, "", 1)
				if c != "" && !encountered[c] {
					encountered[c] = true
					cs = append(cs, c)
				}
			}
			return nil
		})

		if err != nil {
			return cs, err
		}
	}

	return cs, nil
}

func (c *Cgroups) EnabledSubsystems(cpath string) []string {
	enabled := []string{}
	for _, s := range Subsystems {
		path := fmt.Sprintf("%s/%s%s", c.FsPath, s, cpath)
		if _, err := os.Lstat(path); err != nil {
			continue
		}
		enabled = append(enabled, s)
	}
	return enabled
}

func (c *Cgroups) IsEnableSubsystem(cpath string, sname string) bool {
	path := fmt.Sprintf("%s/%s%s", c.FsPath, sname, cpath)
	if _, err := os.Lstat(path); err != nil {
		return false
	}
	return true
}

func (c *Cgroups) Processes(cpath string) ([]ps.Process, error) {
	subsys := c.EnabledSubsystems(cpath)

	pids := []int{}
	processes := []ps.Process{}
	encountered := make(map[int]bool)

	for _, s := range subsys {
		path := fmt.Sprintf("%s/%s%s", c.FsPath, s, cpath)
		err := filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			base := filepath.Base(p)
			if base != "cgroup.procs" {
				return nil
			}
			procs, err := os.Open(filepath.Join(path, "cgroup.procs"))
			if err != nil {
				return err
			}
			defer procs.Close()

			scanner := bufio.NewScanner(procs)
			for scanner.Scan() {
				if t := scanner.Text(); t != "" {
					pid, err := strconv.Atoi(t)
					if err != nil {
						return err
					}
					if !encountered[pid] {
						encountered[pid] = true
						pids = append(pids, pid)
					}
				}
			}
			return nil
		})
		if err != nil {
			return processes, err
		}
	}

	for _, pid := range pids {
		pr, err := ps.FindProcess(pid)
		if err != nil {
			return processes, err
		}
		processes = append(processes, pr)
	}

	return processes, nil
}

func (c *Cgroups) ReadSimple(cpath string, sname string, stat string) (string, error) {
	path := fmt.Sprintf("%s/%s%s/%s", c.FsPath, sname, cpath, stat)
	val, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := strings.TrimRight(string(val), "\n")
	if str == "" {
		str = ""
	}
	return str, nil
}
