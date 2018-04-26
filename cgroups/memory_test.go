package cgroups

import (
	"reflect"
	"testing"
)

func TestMemory(t *testing.T) {
	c := Cgroups{FsPath: testFs()}
	h := "/"

	actualLabel, actualValue := c.Memory(h)
	expectedLabel := []string{
		"memory.usage_in_bytes", "memory.max_usage_in_bytes", "memory.limit_in_bytes",
		"memory.stat.cache", "memory.stat.rss", "memory.stat.rss_huge",
		"memory.stat.mapped_file", "memory.stat.dirty", "memory.stat.writeback",
		"memory.stat.pgpgin", "memory.stat.pgpgout", "memory.stat.pgfault",
		"memory.stat.pgmajfault", "memory.stat.inactive_anon", "memory.stat.active_anon",
		"memory.stat.inactive_file", "memory.stat.active_file", "memory.stat.unevictable",
		"memory.stat.hierarchical_memory_limit", "memory.stat.total_cache",
		"memory.stat.total_rss", "memory.stat.total_rss_huge", "memory.stat.total_mapped_file",
		"memory.stat.total_dirty", "memory.stat.total_writeback", "memory.stat.total_pgpgin",
		"memory.stat.total_pgpgout", "memory.stat.total_pgfault", "memory.stat.total_pgmajfault",
		"memory.stat.total_inactive_anon", "memory.stat.total_active_anon",
		"memory.stat.total_inactive_file", "memory.stat.total_active_file",
		"memory.stat.total_unevictable",
	}
	expectedValue := []string{
		"804159488", "1014558720", "9223372036854771712", "46403584", "552960", "0",
		"5423104", "49152", "0", "43236", "31772", "45737", "87", "28672", "606208",
		"4177920", "40206336", "1937408", "9223372036854771712", "701714432", "102445056",
		"0", "53059584", "69632", "0", "38107452", "37911124", "69743793", "1803",
		"5382144", "101584896", "163041280", "530374656", "3739648",
	}

	if !reflect.DeepEqual(actualLabel, expectedLabel) {
		t.Errorf("actual %v\nwant %v", actualLabel, expectedLabel)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Errorf("actual %v\nwant %v", actualValue, expectedValue)
	}
}
