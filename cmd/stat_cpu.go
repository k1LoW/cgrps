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

func NewCpuStat() (*termui.Par, *termui.List, *termui.List) {
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

	return title, label, data
}

var cgroupParam = []string{
	"cpu.cfs_period_us",
	"cpu.cfs_quota_us",
	"cpu.rt_period_us",
	"cpu.rt_runtime_us",
	"cpu.shares",
	"cpuset.cpus",
	"cpuset.mems",
}

func DrawCpuStat(cpath string, control cgroups.Cgroup, label *termui.List, data *termui.List) {
	if !util.IsEnableSubsystem("cpu", control) {
		return
	}
	stats, err := control.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	d := []string{}

	// cgroupParam
	for _, s := range cgroupParam {
		splited := strings.SplitN(s, ".", 2)
		val, err := util.ReadSimple(cpath, splited[0], s)
		if err == nil {
			label.Items = append(label.Items, fmt.Sprintf("%s:", s))
			d = append(d, fmt.Sprintf("%v", val))
		}
	}
	// cpuacct.stat.user
	d = append(d, fmt.Sprintf("%v", stats.CPU.Usage.User))
	label.Items = append(label.Items, "cpuacct.stat.user:")
	// cpuacct.stat.system
	d = append(d, fmt.Sprintf("%v", stats.CPU.Usage.Kernel))
	label.Items = append(label.Items, "cpuacct.stat.system:")
	// cpuacct.usage
	d = append(d, fmt.Sprintf("%v", stats.CPU.Usage.Total))
	label.Items = append(label.Items, "cpuacct.usage:")
	// cpuacct.usage_percpu
	// d = append(d, fmt.Sprintf("%v", stats.CPU.Usage.PerCPU))
	// cpu.stat.nr_periods
	d = append(d, fmt.Sprintf("%v", stats.CPU.Throttling.Periods))
	label.Items = append(label.Items, "cpu.stat.nr_periods:")
	// cpu.stat.nr_throttled
	d = append(d, fmt.Sprintf("%v", stats.CPU.Throttling.ThrottledPeriods))
	label.Items = append(label.Items, "cpu.stat.nr_throttled:")
	// cpu.stat.throttled_time
	d = append(d, fmt.Sprintf("%v", stats.CPU.Throttling.ThrottledTime))
	label.Items = append(label.Items, "cpu.stat.throttled_time:")

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
