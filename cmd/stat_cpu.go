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
	"bufio"
	"fmt"
	"github.com/gizak/termui"
	"github.com/k1LoW/cgrps/util"
	"strconv"
	"strings"
)

type CPUStat struct {
	Items map[string]uint64
}

// NewCPUStat create new CPU stat vals
func NewCPUStat() (*termui.Par, *termui.List, *termui.List, *CPUStat) {
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

// DrawCPUStat gather CPU stat vals and set
func DrawCPUStat(cpath string, label *termui.List, data *termui.List, total *CPUStat) {
	if !util.IsEnableSubsystem(cpath, "cpu,cpuacct") && !util.IsEnableSubsystem(cpath, "cpuset") {
		return
	}

	d := []string{}
	var l string
	var t uint64
	tick := util.ClkTck()

	// cgroupCPU
	for _, s := range cgroupCPU {
		splited := strings.SplitN(s, ".", 2)
		val, err := util.ReadSimple(cpath, splited[0], s)
		if val != "" && err == nil {
			l = fmt.Sprintf("%s:", s)
			label.Items = append(label.Items, l)
			if strings.Index(l, "_us") > 0 {
				fval, err := strconv.ParseFloat(val, 64)
				if err != nil {
					panic(err)
				}
				d = append(d, util.Usec(fval))
			} else {
				d = append(d, fmt.Sprintf("%v       ", val))
			}
		}
	}

	// cpuacct.stat
	stat, err := util.ReadSimple(cpath, "cpu", "cpuacct.stat")
	if err == nil {
		in := strings.NewReader(stat)
		scanner := bufio.NewScanner(in)

		for scanner.Scan() {
			line := scanner.Text()
			splited := strings.SplitN(line, " ", 2)
			k := splited[0]
			v := splited[1]
			l = fmt.Sprintf("cpuacct.stat.%s:", k)
			t, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				panic(err)
			}
			if prev, ok := total.Items[l]; ok {
				d = append(d, util.UsecPerSec(util.Round(float64(t-prev)/tick*1000000)))
			} else {
				d = append(d, "-       ")
			}
			total.Items[l] = t
			label.Items = append(label.Items, l)
		}
	}

	// cpuacct.usage
	v, err := util.ReadSimple(cpath, "cpu", "cpuacct.usage")
	l = "cpuacct.usage:"
	t, err = strconv.ParseUint(v, 10, 64)
	if err != nil {
		panic(err)
	}
	if prev, ok := total.Items[l]; ok {
		d = append(d, util.UsecPerSec(util.Round(float64(t-prev)/1000)))
	} else {
		d = append(d, "-       ")
	}
	total.Items[l] = t
	label.Items = append(label.Items, l)

	// cpuacct.usage_percpu

	// cpu.stat
	stat, err = util.ReadSimple(cpath, "cpu", "cpu.stat")
	if err == nil {
		in := strings.NewReader(stat)
		scanner := bufio.NewScanner(in)

		for scanner.Scan() {
			line := scanner.Text()
			splited := strings.SplitN(line, " ", 2)
			k := splited[0]
			v := splited[1]
			l = fmt.Sprintf("cpu.stat.%s:", k)
			t, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				panic(err)
			}
			if prev, ok := total.Items[l]; ok {
				if strings.Index(l, "throttled_time") > 0 {
					d = append(d, util.UsecPerSec(util.Round(float64(t-prev)/1000)))
				} else {
					d = append(d, fmt.Sprintf("%v       ", t-prev))
				}
			} else {
				d = append(d, "-       ")
			}
			total.Items[l] = t
			label.Items = append(label.Items, l)
		}
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
