package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
)

var (
	memSize    string
	growthTime string
)

func init() {
	flag.StringVar(&memSize, "size", "0KB", "")
	flag.StringVar(&growthTime, "time", "0s", "")
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

func main() {
	length, err := humanize.ParseBytes(memSize)
	if err != nil {
		// TODO
		print(err)
	}

	timeLine, err := time.ParseDuration(growthTime)
	if err != nil {
		// TODO
	}

	data, err := syscall.Mmap(-1, 0, int(length), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		// TODO
		print(err)
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
