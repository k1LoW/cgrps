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
	"strings"

	"github.com/gizak/termui"
	"github.com/k1LoW/cgrps/cgroups"
	"github.com/k1LoW/cgrps/util"
)

// NewMemoryStat create new Memory stat vals
func NewMemoryStat(h string) (*termui.Par, *termui.List, *termui.List) {
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

	DrawMemoryStat(h, label, data)

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

// IsEnabledMemoryStat return Memory stat enabled or not
func IsEnabledMemoryStat(h string) bool {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	if c.IsAttachedSubsystem(h, "memory") {
		return true
	}
	return false
}

// DrawMemoryStat gather Memory stat vals and set
func DrawMemoryStat(h string, label *termui.List, data *termui.List) {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	if !IsEnabledMemoryStat(h) {
		return
	}

	d := []string{}
	label.Items = []string{}

	// memory
	memoryLabel, memoryValue := c.Memory(h)
	total := make(map[string]string)
	for k, v := range memoryLabel {
		if strings.Index(v, "memory.stat.total_") == 0 {
			replaced := strings.Replace(v, "memory.stat.total_", "memory.stat.", 1)
			total[replaced] = memoryValue[k]
		}
	}

	for k, v := range memoryValue {
		l := fmt.Sprintf("%s:", memoryLabel[k])
		if strings.Index(l, "memory.stat.total_") == 0 {
			continue
		}
		label.Items = append(label.Items, l)
		if strings.Index(l, "memory.stat") == 0 {
			if t, ok := total[memoryLabel[k]]; ok {
				if strings.Index(l, "memory.stat.pgpg") == 0 {
					d = append(d, fmt.Sprintf("%v / %v", v, t))
				} else {
					d = append(d, fmt.Sprintf("%v / %v", util.Bytes(v), util.Bytes(t)))
				}
			} else {
				if strings.Index(l, "memory.stat.pgpg") == 0 {
					d = append(d, fmt.Sprintf("%v", v))
				} else {
					d = append(d, fmt.Sprintf("%v", util.Bytes(v)))
				}
			}
		} else {
			d = append(d, fmt.Sprintf("%v", util.Bytes(v)))
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
