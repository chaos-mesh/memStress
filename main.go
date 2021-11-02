package main

import (
	"flag"
	"fmt"
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

func linearGrow(data []byte, interval time.Duration, length uint64, startTime time.Time, timeLine time.Duration) {
	fmt.Println(interval)
	sysPageSize := os.Getpagesize()
	minPageQuantity := int(time.Millisecond * 100 / interval)
	pageCount := 0
	resLength := length
	for i := 0; uint64(i) < length; i += sysPageSize {
		data[i] = 1
		if minPageQuantity > 0 {
			pageCount += 1
			acculatedPage := pageCount % minPageQuantity
			if acculatedPage == 0 {
				time.Sleep(time.Duration(minPageQuantity) * interval)
				resLength = length - uint64(i)
				interval = updateInterval(timeLine - time.Since(startTime), resLength)
			}
		} else {
			time.Sleep(interval)
		}
	}

	resTime := time.Duration(resLength / uint64(sysPageSize) * uint64(interval))
	if resTime > 100 * time.Millisecond {
		time.Sleep(resTime)
	}
}

func updateInterval(timeLine time.Duration, length uint64) time.Duration {
	sysPageSize := uint64(os.Getpagesize())
	interval := time.Duration(timeLine) / time.Duration(length / sysPageSize)
	return interval
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

	sysPageSize := os.Getpagesize()
	interval := time.Duration(timeLine) / time.Duration(length / uint64(sysPageSize))

	if interval > time.Nanosecond {
		linearGrow(data, interval, length, time.Now(), timeLine)
	} else {
		for i := 0; uint64(i) < length; i += sysPageSize {
			data[i] = 1
		}
	}
	
	for {
		time.Sleep(time.Second * 2)
	}
}
