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

// IsEnabledBlkioStat return Memory stat enabled or not
func IsEnabledBlkioStat(h string) bool {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	if c.IsAttachedSubsystem(h, "blkio") {
		return true
	}
	return false
}

// NewBlkioStat create new Blkio stat vals
func NewBlkioStat(h string) (*termui.Par, *termui.List, *termui.List) {
	title := termui.NewPar("BLKIO")
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

	DrawBlkioStat(h, label, data)

	return title, label, data
}

// DrawBlkioStat gather Blkio stat vals and set
func DrawBlkioStat(h string, label *termui.List, data *termui.List) {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	if !IsEnabledBlkioStat(h) {
		return
	}

	d := []string{}

	// blkio
	blkioLabel, blkioValue := c.Blkio(h)
	for k, v := range blkioValue {
		l := fmt.Sprintf("%s:", blkioLabel[k])
		label.Items = append(label.Items, l)
		d = append(d, fmt.Sprintf("%v", v))
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
