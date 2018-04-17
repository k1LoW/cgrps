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
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Subsystems() ([]string, error) {
	ss := []string{}
	subsystems, err := cgroups.V1()
	if err != nil {
		return nil, err
	}
	for _, s := range subsystems {
		ss = append(ss, string(s.Name()))
	}
	return ss, nil
}

func Hierarchy(c string) cgroups.Hierarchy {
	f := func() ([]cgroups.Subsystem, error) {
		enabled := []cgroups.Subsystem{}
		subsystems, err := cgroups.V1()
		if err != nil {
			return nil, err
		}
		for _, s := range subsystems {
			path := fmt.Sprintf("/sys/fs/cgroup/%s%s", s.Name(), c)
			if _, err := os.Lstat(path); err != nil {
				continue
			}
			enabled = append(enabled, s)
		}
		return enabled, nil
	}
	return f
}

func IsEnableSubsystem(sname string, control cgroups.Cgroup) bool {
	subsys := control.Subsystems()
	for _, s := range subsys {
		if string(s.Name()) == sname {
			return true
		}
	}
	return false
}

func ReadSimple(cpath string, sname string, stat string) (string, error) {
	path := fmt.Sprintf("/sys/fs/cgroup/%s%s/%s", sname, cpath, stat)
	val, err := ioutil.ReadFile(path)
	if err != nil {
		return "-", err
	}
	str := strings.TrimRight(string(val), "\n")
	if str == "" {
		str = "-"
	}
	return str, nil
}

func Bytes(v string) string {
	parsed, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return v
	}
	return humanize.Bytes(parsed)
}

func Usec(v float64) string {
	vint := int64(v)
	return fmt.Sprintf("%v us    ", humanize.Comma(vint))
}

func UsecPerSec(v float64) string {
	vint := int64(v)
	return fmt.Sprintf("%v us/sec", humanize.Comma(vint))
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

func Round(f float64) float64 {
	return math.Floor(f + 0.5)
}
