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
	"github.com/k1LoW/cgrps/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list cgroups.",
	Long:  `list cgroups.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		subsys, err := util.Subsystems()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cs := []string{}
		encountered := make(map[string]bool)

		for _, s := range subsys {
			searchDir := fmt.Sprintf("/sys/fs/cgroup/%s", s)

			err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
				if f.IsDir() {
					c := strings.Replace(path, searchDir, "", 1)
					if c != "" && !encountered[c] {
						encountered[c] = true
						cs = append(cs, c)
					}
				}
				return nil
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		for _, c := range cs {
			fmt.Println(c)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
