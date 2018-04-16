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
	"github.com/containerd/cgroups"
	"github.com/gizak/termui"
	"github.com/k1LoW/cgrps/util"
	"os"
	"strings"
)

type CPUStat struct {
	Items map[string]uint64
}

// NewCPUStat ...
func NewCPUStat() (*termui.Par, *termui.List, *termui.List, *CPUStat) {
	title := termui.NewPar("CPU")
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

	return title, label, data, &total
}

var cgroupCPU = []string{
	"cpu.cfs_period_us",
	"cpu.cfs_quota_us",
	"cpu.rt_period_us",
	"cpu.rt_runtime_us",
	"cpu.shares",
	"cpuset.cpus",
	"cpuset.mems",
}

// DrawCPUStat ...
func DrawCPUStat(cpath string, control cgroups.Cgroup, label *termui.List, data *termui.List, total *CPUStat) {
	if !util.IsEnableSubsystem("cpu", control) {
		return
	}
	stats, err := control.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	d := []string{}
	var l string
	var t uint64
	// tick := util.ClkTck()

	// cgroupCPU
	for _, s := range cgroupCPU {
		splited := strings.SplitN(s, ".", 2)
		val, err := util.ReadSimple(cpath, splited[0], s)
		if err == nil {
			l = fmt.Sprintf("%s:", s)
			label.Items = append(label.Items, l)
			d = append(d, fmt.Sprintf("%v", val))
		}
	}
	// cpuacct.stat.user
	l = "cpuacct.stat.user:"
	t = stats.CPU.Usage.User
	if prev, ok := total.Items[l]; ok {
		d = append(d, fmt.Sprintf("%v", t-prev))
	} else {
		d = append(d, fmt.Sprintf("%v", t))
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)
	// cpuacct.stat.system
	l = "cpuacct.stat.system:"
	t = stats.CPU.Usage.Kernel
	if prev, ok := total.Items[l]; ok {
		d = append(d, fmt.Sprintf("%v", t-prev))
	} else {
		d = append(d, fmt.Sprintf("%v", t))
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)
	// cpuacct.usage
	l = "cpuacct.usage:"
	t = stats.CPU.Usage.Total
	if prev, ok := total.Items[l]; ok {
		d = append(d, fmt.Sprintf("%v", t-prev))
	} else {
		d = append(d, fmt.Sprintf("%v", t))
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)
	// cpuacct.usage_percpu
	// d = append(d, fmt.Sprintf("%v", stats.CPU.Usage.PerCPU))
	// cpu.stat.nr_periods
	l = "cpu.stat.nr_periods:"
	t = stats.CPU.Throttling.Periods
	if prev, ok := total.Items[l]; ok {
		d = append(d, fmt.Sprintf("%v", t-prev))
	} else {
		d = append(d, fmt.Sprintf("%v", t))
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)
	// cpu.stat.nr_throttled
	l = "cpu.stat.nr_throttled:"
	t = stats.CPU.Throttling.ThrottledPeriods
	if prev, ok := total.Items[l]; ok {
		d = append(d, fmt.Sprintf("%v", t-prev))
	} else {
		d = append(d, fmt.Sprintf("%v", t))
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)
	// cpu.stat.throttled_time
	l = "cpu.stat.throttled_time:"
	t = stats.CPU.Throttling.ThrottledTime
	if prev, ok := total.Items[l]; ok {
		d = append(d, fmt.Sprintf("%v", t-prev))
	} else {
		d = append(d, fmt.Sprintf("%v", t))
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)

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
