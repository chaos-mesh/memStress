// Copyright 2022 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
)

var (
	memSize    string
	growthTime string
	workers    int
	client     bool
)

func init() {
	flag.StringVar(&memSize, "size", "0KB", "size of memory you want to allocate")
	flag.StringVar(&growthTime, "time", "0s", "time to reach the size of memory you allocated")
	flag.IntVar(&workers, "workers", 1, "number of workers allocating memory")
	flag.BoolVar(&client, "client", false, "the process runs as a client")
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
			data[i] = 0
		}
	}

	for {
		time.Sleep(time.Second * 2)
	}
}

func main() {
	if !client {
		workQueue := make(chan struct{}, workers)
		for {
			workQueue <- struct{}{}
			go func() {
				cmd := exec.Command("./memStress", "--size", memSize, "--workers", fmt.Sprintf("%d", workers),
					"--time", growthTime, "--client", "1")
				cmd.SysProcAttr = &syscall.SysProcAttr{
					Pdeathsig: syscall.SIGTERM,
				}
				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
				<-workQueue
			}()
			time.Sleep(time.Second)
		}
	} else {
		totalMem, err := getTotalMemory()
		if err != nil {
			fmt.Println("get total memory err:", err)
			return
		}
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
			length = uint64(float64(totalMem) / 100.0 * percentage)
		}

		timeLine, err := time.ParseDuration(growthTime)
		if err != nil {
			// TODO
		}
		run(length/uint64(workers), timeLine)
	}
}
