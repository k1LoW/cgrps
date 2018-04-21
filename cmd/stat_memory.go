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
	"strings"
)

func NewMemoryStat() (*termui.Par, *termui.List, *termui.List) {
	title := termui.NewPar("MEMORY")
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

var cgroupMemory = []string{
	"memory.usage_in_bytes",
	"memory.memsw.usage_in_bytes",
	"memory.max_usage_in_bytes",
	"memory.memsw.max_usage_in_bytes",
	"memory.limit_in_bytes",
	"memory.memsw.limit_in_bytes",
}

func DrawMemoryStat(cpath string, label *termui.List, data *termui.List) {
	c := util.Cgroups{FsPath: "/sys/fs/cgroup"}
	if !c.IsEnableSubsystem(cpath, "memory") {
		return
	}

	d := []string{}

	// cgroupMemory
	for _, s := range cgroupMemory {
		splited := strings.SplitN(s, ".", 2)
		val, err := c.ReadSimple(cpath, splited[0], s)
		if err == nil {
			label.Items = append(label.Items, fmt.Sprintf("%s:", s))
			d = append(d, fmt.Sprintf("%v", util.Bytes(val)))
		}
	}

	// memoty.stat
	stat, err := c.ReadSimple(cpath, "memory", "memory.stat")
	if err == nil {
		in := strings.NewReader(stat)
		scanner := bufio.NewScanner(in)

		stats := make(map[string]string)
		lines := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			splited := strings.SplitN(line, " ", 2)
			k := splited[0]
			v := splited[1]
			if strings.Index(k, "total_") == 0 {
				k = strings.Replace(k, "total_", "", 1)
				stats[k] = v
			} else {
				lines = append(lines, line)
			}
		}

		for _, l := range lines {
			splited := strings.SplitN(l, " ", 2)
			k := splited[0]
			v := splited[1]
			if total, ok := stats[k]; ok {
				if strings.Index(k, "pgpg") == 0 {
					d = append(d, fmt.Sprintf("%v / %v", v, total))
				} else {
					d = append(d, fmt.Sprintf("%v / %v", util.Bytes(v), util.Bytes(total)))
				}
			} else {
				if strings.Index(k, "pgpg") == 0 {
					d = append(d, fmt.Sprintf("%v", v))
				} else {
					d = append(d, fmt.Sprintf("%v", util.Bytes(v)))
				}
			}
			label.Items = append(label.Items, fmt.Sprintf("memory.stat.%s:", k))
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
