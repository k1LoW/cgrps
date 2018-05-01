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

package cgroups

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

// ClkTck return clocks per sec (CLK_TCK)
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

// Cgroups struct
type Cgroups struct {
	FsPath string
}

// List show all cgroup hierarchies
func (c *Cgroups) List() ([]string, error) {
	subsys := c.AttachedSubsystems("/")

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

// AttachedSubsystems return active subsystem in specific cgroup hierarchy
func (c *Cgroups) AttachedSubsystems(h string) []string {
	enabled := []string{}
	for _, s := range Subsystems {
		path := fmt.Sprintf("%s/%s%s", c.FsPath, s, h)
		if _, err := os.Lstat(path); err != nil {
			continue
		}
		enabled = append(enabled, s)
	}
	return enabled
}

// IsAttachedSubsystem return subsystem is active or not in specific cgroup hierarchy
func (c *Cgroups) IsAttachedSubsystem(h string, sname string) bool {
	path := fmt.Sprintf("%s/%s%s", c.FsPath, sname, h)
	if _, err := os.Lstat(path); err != nil {
		return false
	}
	return true
}

// ListPids return pids in specific cgroup hierarchy
func (c *Cgroups) ListPids(hs []string) []int {
	all := []int{}
	encountered := make(map[int]bool)
	for _, h := range hs {
		all = append(all, c.listPids(h)...)
	}

	merged := []int{}
	for _, pid := range all {
		if !encountered[pid] {
			encountered[pid] = true
			merged = append(merged, pid)
		}
	}

	sort.Ints(merged)

	return merged
}

func (c *Cgroups) listPids(h string) []int {
	subsys := c.AttachedSubsystems(h)

	pids := []int{}

	for _, s := range subsys {
		path := fmt.Sprintf("%s/%s%s", c.FsPath, s, h)
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
				_ = procs.Close()
				return err
			}

			scanner := bufio.NewScanner(procs)
			for scanner.Scan() {
				if t := scanner.Text(); t != "" {
					pid, err := strconv.Atoi(t)
					if err != nil {
						_ = procs.Close()
						return err
					}
					pids = append(pids, pid)
				}
			}
			err = procs.Close()
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return pids
		}
	}

	return pids
}

// ReadSimple read file and return value as string
func (c *Cgroups) ReadSimple(h string, sname string, stat string) (string, error) {
	path := fmt.Sprintf("%s/%s%s/%s", c.FsPath, sname, h, stat)
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
