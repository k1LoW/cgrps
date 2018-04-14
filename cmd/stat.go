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
	"errors"
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/gizak/termui"
	"github.com/k1LoW/cgrps/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
)

// statCmd represents the stat command
var statCmd = &cobra.Command{
	Use:   "stat [CGROUP]",
	Short: "cgroup stat",
	Long:  `cgroup stat.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if terminal.IsTerminal(0) {
			if len(args) < 1 {
				return errors.New("requires [CGROUP] or STDIN")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var cpath string

		if terminal.IsTerminal(0) {
			cpath = args[0]
		} else {
			b, _ := ioutil.ReadAll(os.Stdin)
			cpath = strings.TrimRight(string(b), "\n")
		}

		h := util.Hierarchy(cpath)

		// debug
		control, err := cgroups.Load(h, cgroups.StaticPath(cpath))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stats, err := control.Stat(cgroups.IgnoreNotExist)
		subsys := control.Subsystems()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", stats)
		fmt.Printf("%s\n", subsys)
		// os.Exit(1) // /debug

		err = termui.Init()
		if err != nil {
			panic(err)
		}
		defer termui.Close()

		cgroupLabel, cgroupData := NewCgroupStat(cpath)

		cpuTitle, cpuLabel, cpuData := NewCpuStat()
		memoryTitle, memoryLabel, memoryData := NewMemoryStat()

		termui.Body.AddRows(
			termui.NewRow(
				termui.NewCol(2, 0, cgroupLabel),
				termui.NewCol(4, 0, cgroupData),
			),
			termui.NewRow(
				termui.NewCol(2, 0, cpuTitle),
				termui.NewCol(2, 2, memoryTitle),
			),
			termui.NewRow(
				termui.NewCol(2, 0, cpuLabel),
				termui.NewCol(2, 0, cpuData),
				termui.NewCol(2, 0, memoryLabel),
				termui.NewCol(2, 0, memoryData),
			),
		)
		termui.Body.Align()

		termui.Render(termui.Body)

		termui.Handle("/sys/kbd/<escape>", func(termui.Event) {
			termui.StopLoop()
		})
		termui.Handle("/sys/kbd/q", func(termui.Event) {
			termui.StopLoop()
		})

		termui.Handle("/timer/1s", func(e termui.Event) {
			DrawCgroupStat(cpath, control, cgroupLabel, cgroupData)
			DrawCpuStat(cpath, control, cpuLabel, cpuData)
			DrawMemoryStat(cpath, control, memoryLabel, memoryData)

			termui.Render(termui.Body)
		})

		termui.Handle("/sys/wnd/resize", func(e termui.Event) {
			termui.Body.Width = termui.TermWidth()
			termui.Body.Align()
			termui.Clear()
			termui.Render(termui.Body)
		})

		termui.Loop()

	},
}

func init() {
	rootCmd.AddCommand(statCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
