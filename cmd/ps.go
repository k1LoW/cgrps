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
	"github.com/k1LoW/cgrps/cgroups"
	"github.com/k1LoW/go-ps"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var OutputJSON bool

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps [CGROUP...]",
	Short: "report a snapshot of the current cgroups processes",
	Long:  `report a snapshot of the current cgroups processes.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if terminal.IsTerminal(0) {
			if len(args) < 1 {
				return errors.New("requires [CGROUP...] or STDIN")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var hs []string

		if terminal.IsTerminal(0) {
			hs = args
		} else {
			b, _ := ioutil.ReadAll(os.Stdin)
			hs = strings.Split(string(b), "\n")
		}

		c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
		processes, err := c.Processes(hs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if OutputJSON {
			printPsAsJSON(processes)
		} else {
			printPs(processes)
		}
	},
}

func printPs(processes []ps.Process) {
	fmt.Println(fmt.Sprintf("%5s", "PID"), fmt.Sprintf("%5s", "PPID"), fmt.Sprintf("%15s", "CMD"), "PATH")
	for _, pr := range processes {
		path, err := filepath.EvalSymlinks(fmt.Sprintf("/proc/%d/exe", pr.Pid()))
		if err != nil {
			path = "-"
		}
		fmt.Println(fmt.Sprintf("%5d", pr.Pid()), fmt.Sprintf("%5d", pr.PPid()), fmt.Sprintf("%15s", pr.Executable()), path)
	}
}

type psJSON struct {
	PID  int    `json:"pid"`
	PPID int    `json:"ppid"`
	CMD  string `json:"cmd"`
	PATH string `json:"path"`
}

func printPsAsJSON(processes []ps.Process) {
	list := make([]psJSON, 0, len(processes))
	for _, pr := range processes {
		path, err := filepath.EvalSymlinks(fmt.Sprintf("/proc/%d/exe", pr.Pid()))
		if err != nil {
			path = "-"
		}
		list = append(list, psJSON{PID: pr.Pid(), PPID: pr.PPid(), CMD: pr.Executable(), PATH: path})
	}
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVarP(&OutputJSON, "json", "", false, "print result as JSON format")
}
