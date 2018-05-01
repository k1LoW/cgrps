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

package cmd

import (
	"fmt"
	"github.com/gizak/termui"
	"github.com/k1LoW/cgrps/cgroups"
	"github.com/k1LoW/cgrps/util"
	"strconv"
	"strings"
)

// CPUStat have CPU/CPUSet/CPUAcct stat
type CPUStat struct {
	Items map[string]uint64
}

// IsEnabledCPUStat return CPU stat enabled or not
func IsEnabledCPUStat(h string) bool {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	if c.IsAttachedSubsystem(h, "cpuset") || c.IsAttachedSubsystem(h, "cpuacct") || c.IsAttachedSubsystem(h, "cpu") {
		return true
	}
	return false
}

// NewCPUStat create new CPU stat vals
func NewCPUStat(h string) (*termui.Par, *termui.List, *termui.List, *CPUStat) {
	title := termui.NewPar("CPU/CPUSet/CPUAcct")
	title.Height = 1
	title.Border = false

	label := termui.NewList()
	label.Border = false
	label.ItemFgColor = termui.ColorCyan
	label.Items = []string{}
	label.Height = len(label.Items)

	data := termui.NewList()
	data.Border = false
	data.Items = []string{}
	data.Height = len(label.Items)

	total := CPUStat{}
	total.Items = make(map[string]uint64)

	DrawCPUStat(h, label, data, &total)

	return title, label, data, &total
}

// DrawCPUStat gather CPU stat vals and set
func DrawCPUStat(h string, label *termui.List, data *termui.List, total *CPUStat) {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	if !IsEnabledCPUStat(h) {
		return
	}

	d := []string{}
	label.Items = []string{}
	tick := cgroups.ClkTck()

	// cpu
	cpuLabel, cpuValue := c.CPU(h)
	for k, v := range cpuValue {
		l := fmt.Sprintf("%s:", cpuLabel[k])
		label.Items = append(label.Items, l)
		if strings.Index(l, "_us") > 0 {
			fval, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}
			d = append(d, util.Usec(fval))
		} else if strings.Index(l, "cpu.stat") == 0 {
			// cpu.stat
			current, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				panic(err)
			}
			if prev, ok := total.Items[l]; ok {
				if strings.Index(l, "throttled_time") > 0 {
					d = append(d, util.UsecPerSec(util.Round(float64(current-prev)/1000)))
				} else {
					d = append(d, fmt.Sprintf("%v       ", current-prev))
				}
			} else {
				d = append(d, "-       ")
			}
			total.Items[l] = current
		} else {
			d = append(d, fmt.Sprintf("%v       ", v))
		}
	}

	// cpuacct
	cpuAcctLabel, cpuAcctValue := c.CPUAcct(h)
	for k, v := range cpuAcctValue {
		l := fmt.Sprintf("%s:", cpuAcctLabel[k])
		label.Items = append(label.Items, l)
		current, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		if prev, ok := total.Items[l]; ok {
			if strings.Index(l, "cpuacct.usage") == 0 {
				d = append(d, util.UsecPerSec(util.Round(float64(current-prev)/1000)))
			} else {
				d = append(d, util.UsecPerSec(util.Round(float64(current-prev)/tick*1000000)))
			}
		} else {
			d = append(d, "-       ")
		}
		total.Items[l] = current
	}

	// cpuset
	cpuSetLabel, cpuSetValue := c.CPUSet(h)
	for k, v := range cpuSetValue {
		l := fmt.Sprintf("%s:", cpuSetLabel[k])
		label.Items = append(label.Items, l)
		d = append(d, fmt.Sprintf("%v       ", v))
	}

	maxlen := 1
	for _, v := range d {
		if maxlen < len(v) {
			maxlen = len(v)
		}
	}

	data.Items = nil
	for _, v := range d {
		data.Items = append(data.Items, fmt.Sprintf(fmt.Sprintf("%%%ds", maxlen), v))
	}

	label.Height = len(data.Items)
	data.Height = len(data.Items)
}
