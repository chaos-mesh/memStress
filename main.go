package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
	psutil "github.com/shirou/gopsutil/mem"
)

var (
	memSize    string
	growthTime string
	workers    int
)

func init() {
	flag.StringVar(&memSize, "size", "0KB", "")
	flag.StringVar(&growthTime, "time", "0s", "")
	flag.IntVar(&workers, "workers", 1, "")
	flag.Parse()
}

func linearGrow(data []byte, length uint64, timeLine time.Duration) {
	startTime := time.Now()
	endTime := startTime.Add(timeLine)

	var allocated uint64 = 0
	pageSize := uint64(syscall.Getpagesize())
	interval := time.Millisecond * 10

	for {
		now := time.Now()
		if now.After(endTime) {
			now = endTime
		}
		expected := length * uint64(now.Sub(startTime).Milliseconds()) / uint64(endTime.Sub(startTime).Milliseconds()) / pageSize

		for i := allocated; uint64(i) < expected; i++ {
			data[uint64(i)*pageSize] = 0
		}

		allocated = expected
		if now.Equal(endTime) {
			break
		} else {
			time.Sleep(interval)
		}
	}

}

func run(length uint64, timeLine time.Duration) {
	data, err := syscall.Mmap(-1, 0, int(length), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		// TODO
		fmt.Println(err)
		os.Exit(1)
	}

	if timeLine > time.Nanosecond {
		linearGrow(data, length, timeLine)
	} else {
		sysPageSize := os.Getpagesize()
		for i := 0; uint64(i) < length; i += sysPageSize {
			data[i] = 1
		}
	}

	for {
		time.Sleep(time.Second * 2)
	}
}

func main() {
	memInfo, _ := psutil.VirtualMemory()
	var length uint64

	if memSize[len(memSize)-1] != '%' {
		var err error
		length, err = humanize.ParseBytes(memSize)
		if err != nil {
			// TODO
			fmt.Println(err)
		}
	} else {
		percentage, err := strconv.ParseFloat(memSize[0:len(memSize)-1], 64)
		if err != nil {
			fmt.Println(err)
		}
		length = uint64(float64(memInfo.Total) / 100.0 * percentage)
	}

	timeLine, err := time.ParseDuration(growthTime)
	if err != nil {
		// TODO
	}

	for i := 0; i < workers; i++ {
		go run(length, timeLine)
	}

	for {
		time.Sleep(time.Second * 2)
	}
}
