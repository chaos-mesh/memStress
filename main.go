package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var (
	memSize string
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func init() {
	flag.StringVar(&memSize, "size", "0KB", "")
	flag.Parse()
}

func main() {
	length, err := strconv.Atoi(memSize)
	if err != nil {
		// TODO
		print(err)
	}

	data, err := syscall.Mmap(-1, 0, length*1024*1024*1024, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		// TODO
		print(err)
	}
	for i := 0; i < len(data); i += os.Getpagesize() {
		data[i] = 1
	}	
	for {
		time.Sleep(time.Second * 2)
	}
}
