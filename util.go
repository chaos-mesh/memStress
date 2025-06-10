package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	psutil "github.com/shirou/gopsutil/mem"
)

const (
	cgroupV2Path    = "/sys/fs/cgroup/memory.max"
	cgroupV1Path    = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	cgroupNoLimitV1 = 0x7FFFFFFFFFFFF000
)

// Read file content and parse as uint64
func readUintFromFile(path string) (uint64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	content := strings.TrimSpace(string(data))
	return strconv.ParseUint(content, 10, 64)
}

// Check cgroup v2 memory limit
func getCgroupV2Limit() (uint64, error) {
	data, err := ioutil.ReadFile(cgroupV2Path)
	if err != nil {
		return 0, err
	}
	content := strings.TrimSpace(string(data))
	if content == "max" {
		return 0, fmt.Errorf("cgroup v2: no memory limit set")
	}
	limit, err := strconv.ParseUint(content, 10, 64)
	if err != nil || limit == 0 {
		return 0, fmt.Errorf("cgroup v2: invalid memory limit")
	}
	return limit, nil
}

// Check cgroup v1 memory limit
func getCgroupV1Limit() (uint64, error) {
	limit, err := readUintFromFile(cgroupV1Path)
	if err != nil {
		return 0, err
	}
	// 0 or cgroup's "infinity" value means no limit
	if limit == 0 || limit >= cgroupNoLimitV1 {
		return 0, fmt.Errorf("cgroup v1: no memory limit set")
	}
	return limit, nil
}

// Get total memory, prefer cgroup v2 -> cgroup v1 -> host
func getTotalMemory() (uint64, error) {
	if limit, err := getCgroupV2Limit(); err == nil {
		return limit, nil
	}
	if limit, err := getCgroupV1Limit(); err == nil {
		return limit, nil
	}
	mem, err := psutil.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to get system memory: %v", err)
	}
	return mem.Total, nil
}
