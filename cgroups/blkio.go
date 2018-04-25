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

package cgroups

import (
	"strings"
)

// Blkio return cgroups Blkio values
func (c *Cgroups) Blkio(h string) ([]string, []string) {
	label := []string{}
	value := []string{}

	// params https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt
	var params = []string{
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

	for _, p := range params {
		splited := strings.SplitN(p, ".", 2)
		v, err := c.ReadSimple(h, splited[0], p)
		if v != "" && err == nil {
			label = append(label, p)
			value = append(value, v)
		}
	}

	return label, value
}
