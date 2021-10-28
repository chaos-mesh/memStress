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
	// timeUnit := growthTime[len(growthTime)]
	// if timeUnit = 
	timeLine = int(time.Second) * timeLine

	data, err := syscall.Mmap(-1, 0, length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		// TODO
		print(err)
	}

	sysPageSize := os.Getpagesize()
	duration := time.Duration(timeLine) / time.Duration(length / sysPageSize)
	fmt.Println(duration, length/sysPageSize)
	fmt.Println(time.Now())
	for i := 0; i < length; i += sysPageSize {
		// fmt.Println(time.Now())
		data[i] = 1
		time.Sleep(duration)
	}
	fmt.Println(time.Now())
	time.Sleep(time.Microsecond * 1)
	fmt.Println(time.Now())
	for {
		time.Sleep(time.Second * 2)
	}
}
