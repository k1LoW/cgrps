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
	"os"
)

func NewCgroupStat(cpath string) (*termui.List, *termui.List) {
	label := termui.NewList()
	label.Border = false
	label.ItemFgColor = termui.ColorYellow
	label.Items = []string{
		"  cgroup path:",
		"  cgroup.procs:",
	}
	label.Height = 3

	data := termui.NewList()
	data.Border = false
	data.Items = []string{
		cpath, "-",
	}
	data.Height = 3

	return label, data
}

func DrawCgroupStat(cpath string, control cgroups.Cgroup, label *termui.List, data *termui.List) {
	subsys := control.Subsystems()
	processes, err := control.Processes(subsys[0].Name(), true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// cgroup.procs
	data.Items[1] = fmt.Sprintf("%d", len(processes)) // @todo use /sys/fs/cgroup/$/cgroup.procs directly
}
