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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/gizak/termui"
	"github.com/k1LoW/cgrps/cgroups"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// statCmd represents the stat command
var statCmd = &cobra.Command{
	Use:   "stat [CGROUP]",
	Short: "show current cgroup statistics",
	Long:  `show current cgroup statistics.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if terminal.IsTerminal(0) {
			if len(args) < 1 {
				return errors.New("requires [CGROUP] or STDIN")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var h string

		if terminal.IsTerminal(0) {
			h = args[0]
		} else {
			b, _ := ioutil.ReadAll(os.Stdin)
			h = strings.TrimRight(string(b), "\n")
		}

		if OutputJSON {
			printStatAsJSON(h)
		} else {
			printStat(h)
		}
	},
}

// printStat ...
func printStat(h string) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	title := termui.NewPar("stat")
	title.Height = 1
	title.Border = false
	cgroupLabel, cgroupData := NewCgroupStat(h)

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(2, 0, title),
		),
		termui.NewRow(
			termui.NewCol(2, 0, cgroupLabel),
			termui.NewCol(10, 0, cgroupData),
		),
	)

	cpuTitle, cpuLabel, cpuData, cpuDataTotal := NewCPUStat(h)
	memoryTitle, memoryLabel, memoryData := NewMemoryStat(h)
	blkioTitle, blkioLabel, blkioData := NewBlkioStat(h)
	pidsTitle, pidsLabel, pidsData := NewPidsStat(h)

	titleLists := []*termui.Par{}
	labelLists := []*termui.List{}
	dataLists := []*termui.List{}

	if IsEnabledCPUStat(h) {
		titleLists = append(titleLists, cpuTitle)
		labelLists = append(labelLists, cpuLabel)
		dataLists = append(dataLists, cpuData)
	}
	if IsEnabledMemoryStat(h) {
		titleLists = append(titleLists, memoryTitle)
		labelLists = append(labelLists, memoryLabel)
		dataLists = append(dataLists, memoryData)
	}
	if IsEnabledBlkioStat(h) {
		titleLists = append(titleLists, blkioTitle)
		labelLists = append(labelLists, blkioLabel)
		dataLists = append(dataLists, blkioData)
	}
	if IsEnabledPidsStat(h) {
		titleLists = append(titleLists, pidsTitle)
		labelLists = append(labelLists, pidsLabel)
		dataLists = append(dataLists, pidsData)
	}

	row := int(math.Ceil(float64(len(titleLists)) / 3))
	for i := 0; i < row; i++ {
		t := []*termui.Row{}
		d := []*termui.Row{}

		index := i * 3

		if len(titleLists) > index {
			t = append(t, termui.NewCol(2, 0, titleLists[index]))
			d = append(d, termui.NewCol(2, 0, labelLists[index]))
			d = append(d, termui.NewCol(2, 0, dataLists[index]))
		}
		if len(titleLists) > index+1 {
			t = append(t, termui.NewCol(2, 2, titleLists[index+1]))
			d = append(d, termui.NewCol(2, 0, labelLists[index+1]))
			d = append(d, termui.NewCol(2, 0, dataLists[index+1]))
		}
		if len(titleLists) > index+2 {
			t = append(t, termui.NewCol(2, 2, titleLists[index+2]))
			d = append(d, termui.NewCol(2, 0, labelLists[index+2]))
			d = append(d, termui.NewCol(2, 0, dataLists[index+2]))
		}
		termui.Body.AddRows(
			termui.NewRow(t...),
		)
		termui.Body.AddRows(
			termui.NewRow(d...),
		)
	}

	termui.Body.Align()

	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/<escape>", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		DrawCgroupStat(h, cgroupLabel, cgroupData)
		DrawCPUStat(h, cpuLabel, cpuData, cpuDataTotal)
		DrawMemoryStat(h, memoryLabel, memoryData)
		DrawBlkioStat(h, blkioLabel, blkioData)
		DrawPidsStat(h, pidsLabel, pidsData)

		termui.Render(termui.Body)
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Loop()
}

func printStatAsJSON(h string) {
	c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
	stat := map[string]interface{}{}
	if c.IsAttachedSubsystem(h, "cpu") {
		cpu := map[string]interface{}{}
		cpuLabel, cpuValue := c.CPU(h)
		for k, l := range cpuLabel {
			cpu[l] = cpuValue[k]
		}
		stat["cpu"] = cpu
	}
	if c.IsAttachedSubsystem(h, "cpuacct") {
		cpuacct := map[string]interface{}{}
		cpuAcctLabel, cpuAcctValue := c.CPUAcct(h)
		for k, l := range cpuAcctLabel {
			cpuacct[l] = cpuAcctValue[k]
		}
		stat["cpuacct"] = cpuacct
	}
	if c.IsAttachedSubsystem(h, "cpuset") {
		cpuset := map[string]interface{}{}
		cpuSetLabel, cpuSetValue := c.CPUSet(h)
		for k, l := range cpuSetLabel {
			cpuset[l] = cpuSetValue[k]
		}
		stat["cpuset"] = cpuset
	}
	if c.IsAttachedSubsystem(h, "memory") {
		memory := map[string]interface{}{}
		memoryLabel, memoryValue := c.Memory(h)
		for k, l := range memoryLabel {
			memory[l] = memoryValue[k]
		}
		stat["memory"] = memory
	}
	if c.IsAttachedSubsystem(h, "blkio") {
		blkio := map[string]interface{}{}
		blkioLabel, blkioValue := c.Blkio(h)
		for k, l := range blkioLabel {
			blkio[l] = blkioValue[k]
		}
		stat["blkio"] = blkio
	}
	if c.IsAttachedSubsystem(h, "pids") {
		pids := map[string]interface{}{}
		pidsLabel, pidsValue := c.Pids(h)
		for k, l := range pidsLabel {
			pids[l] = pidsValue[k]
		}
		stat["pids"] = pids
	}

	jsonBytes, err := json.Marshal(stat)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(statCmd)
	statCmd.Flags().BoolVarP(&OutputJSON, "json", "", false, "print result as JSON format")
}
