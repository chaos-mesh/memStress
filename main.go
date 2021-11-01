package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
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

func linearGrow(data []byte, interval time.Duration, length int, startTime time.Time, timeLine int) {
	fmt.Println(interval)
	sysPageSize := os.Getpagesize()
	minPageQuantity := int(time.Millisecond * 100 / interval)
	pageCount := 0
	resLength := length
	for i := 0; i < length; i += sysPageSize {
		data[i] = 1
		if minPageQuantity > 0 {
			pageCount += 1
			acculatedPage := pageCount % minPageQuantity
			if acculatedPage == 0 {
				time.Sleep(time.Duration(minPageQuantity) * interval)
				resLength = length - i
				interval = updateInterval(timeLine - int(time.Since(startTime)), resLength)
			}
		} else {
			time.Sleep(interval)
		}
	}

	resTime := time.Duration(resLength / sysPageSize * int(interval))
	if resTime > 100 * time.Millisecond {
		time.Sleep(resTime)
	}
}

func updateInterval(timeLine int, length int) time.Duration {
	sysPageSize := os.Getpagesize()
	interval := time.Duration(timeLine) / time.Duration(length / sysPageSize)
	return interval
}

func main() {
	memSize = strings.ToUpper(memSize)
	length, err := strconv.Atoi(memSize[:len(memSize)-2])
	if err != nil {
		// TODO
		print(err)
	}
	sizeUnit := memSize[len(memSize)-2:]
	if sizeUnit == "KB" {
		length *= 1024
	} else if sizeUnit == "MB" {
		length *= 1024 * 1024
	} else if sizeUnit == "GB" {
		length *= 1024 * 1024 * 1024
	} else {
		// TODO
	}

	growthTime = strings.ToLower(growthTime)
	timeLine, err := strconv.Atoi(growthTime[:len(growthTime)-1])
	if err != nil {
		// TODO
		print(err)
	}
	timeUnit := growthTime[len(growthTime)-1:]
	if timeUnit == "s" {
		timeLine = int(time.Second) * timeLine
	} else if timeUnit == "m" {
		timeLine = int(time.Minute) * timeLine
	} else if timeUnit == "h" {
		timeLine = int(time.Hour) * timeLine
	} else {
		// TODO
	}
	

	data, err := syscall.Mmap(-1, 0, length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		// TODO
		print(err)
	}

	sysPageSize := os.Getpagesize()
	interval := time.Duration(timeLine) / time.Duration(length / sysPageSize)

	if interval > time.Nanosecond {
		linearGrow(data, interval, length, time.Now(), timeLine)
	} else {
		for i := 0; i < length; i += sysPageSize {
			data[i] = 1
		}
	}
	for {
		time.Sleep(time.Second * 2)
	}
}
