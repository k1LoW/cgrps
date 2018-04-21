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
	"github.com/k1LoW/cgrps/util"
	"strings"
)

func NewBlkioStat() (*termui.Par, *termui.List, *termui.List) {
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

	return title, label, data
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt
var cgroupBlkio = []string{
	"blkio.throttle.read_bps_device",
	"blkio.throttle.write_bps_device",
	"blkio.throttle.read_iops_device",
	"blkio.throttle.write_iops_device",

	// @todo
	// "blkio.throttle.io_serviced",
	// "blkio.throttle.io_service_bytes",

	"blkio.weight",
	"blkio.weight_device",

	// @todo
	// "blkio.leaf_weight[_device]",
	// "blkio.time",
	// "blkio.sectors",
	// "blkio.io_service_bytes",
	// "blkio.io_serviced",
	// "blkio.io_service_time",
	// "blkio.io_wait_time",
	// "blkio.io_merged",
	// "blkio.io_queued",
	// "blkio.avg_queue_size",
	// "blkio.group_wait_time",
	// "blkio.empty_time",
	// "blkio.idle_time",
	// "blkio.dequeue",
}

func DrawBlkioStat(cpath string, label *termui.List, data *termui.List) {
	c := util.Cgroups{FsPath: "/sys/fs/cgroup"}
	if !c.IsEnableSubsystem(cpath, "blkio") {
		return
	}

	d := []string{}

	// cgroupBlkio
	for _, s := range cgroupBlkio {
		splited := strings.SplitN(s, ".", 2)
		val, err := c.ReadSimple(cpath, splited[0], s)
		if val != "" && err == nil {
			label.Items = append(label.Items, fmt.Sprintf("%s:", s))
			d = append(d, fmt.Sprintf("%v", val))
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
