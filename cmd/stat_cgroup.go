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
)

// NewCgroupStat create new Cgroup stat vals
func NewCgroupStat(h string) (*termui.List, *termui.List) {
	label := termui.NewList()
	label.Border = false
	label.ItemFgColor = termui.ColorGreen
	label.Items = []string{
		"  cgroup hierarchy:",
		"  cgroup.procs:",
		"  subsystems:",
	}
	label.Height = 4

	data := termui.NewList()
	data.Border = false
	data.Items = []string{
		h, "-", "-",
	}
	data.Height = 4

	return label, data
}

// DrawCgroupStat gather cgroup stat vals and set
func DrawCgroupStat(h string, label *termui.List, data *termui.List) {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	pids := c.Pids([]string{h})
	// cgroup.procs
	data.Items[1] = fmt.Sprintf("%d", len(pids))
	// subsystems
	data.Items[2] = fmt.Sprintf("%v", c.AttachedSubsystems(h))
}
