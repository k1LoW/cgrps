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
	"bufio"
	"fmt"
	"strings"
)

// Memory return cgroups memory values
func (c *Cgroups) Memory(h string) ([]string, []string) {
	label := []string{}
	value := []string{}

	// params
	var params = []string{
		"memory.usage_in_bytes",
		"memory.memsw.usage_in_bytes",
		"memory.max_usage_in_bytes",
		"memory.memsw.max_usage_in_bytes",
		"memory.limit_in_bytes",
		"memory.memsw.limit_in_bytes",
	}
	for _, p := range params {
		splited := strings.SplitN(p, ".", 2)
		v, err := c.ReadSimple(h, splited[0], p)
		if v != "" && err == nil {
			label = append(label, p)
			value = append(value, v)
		}
	}

	// stat
	stat, err := c.ReadSimple(h, "memory", "memory.stat")
	if err == nil {
		in := strings.NewReader(stat)
		scanner := bufio.NewScanner(in)

		for scanner.Scan() {
			line := scanner.Text()
			splited := strings.SplitN(line, " ", 2)
			s := fmt.Sprintf("memory.stat.%s", splited[0])
			v := splited[1]
			label = append(label, s)
			value = append(value, v)
		}
	}

	return label, value
}
