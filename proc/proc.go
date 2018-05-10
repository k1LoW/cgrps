package proc

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Proc struct
type Proc struct {
	FsPath string
}

// Cgroup List cgroup
func (p *Proc) Cgroup(pids []string) ([]string, error) {
	cs := []string{}
	encountered := make(map[string]bool)

	for _, pid := range pids {
		path := fmt.Sprintf("%s/%s/cgroup", p.FsPath, pid)
		cgroup, err := os.Open(path)
		if err != nil {
			_ = cgroup.Close()
			return cs, err
		}
		scanner := bufio.NewScanner(cgroup)
		for scanner.Scan() {
			if t := scanner.Text(); t != "" {
				splited := strings.SplitN(t, ":", 3)
				c := splited[2]
				if c != "" && !encountered[c] {
					encountered[c] = true
					cs = append(cs, c)
				}
			}
		}
		err = cgroup.Close()
		if err != nil {
			return cs, err
		}
	}

	sort.Strings(cs)

	return cs, nil
}
