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
	"github.com/k1LoW/cgrps/cgroups"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps [CGROUP]",
	Short: "report a snapshot of the current cgroups processes",
	Long:  `report a snapshot of the current cgroups processes.`,
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

		c := cgroups.Cgroups{FsPath: "/sys/fs/cgroup"}
		processes, err := c.Processes(cpath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(fmt.Sprintf("%5s", "PID"), fmt.Sprintf("%5s", "PPID"), fmt.Sprintf("%15s", "CMD"), "PATH")
		for _, pr := range processes {
			path, err := filepath.EvalSymlinks(fmt.Sprintf("/proc/%d/exe", pr.Pid()))
			if err != nil {
				path = "-"
			}
			fmt.Println(fmt.Sprintf("%5d", pr.Pid()), fmt.Sprintf("%5d", pr.PPid()), fmt.Sprintf("%15s", pr.Executable()), path)
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// psCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
