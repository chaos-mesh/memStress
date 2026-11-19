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
	flag.StringVar(&memSize, "size", "0KB", "total memory to allocate across all workers (e.g. 64MiB)")
	flag.StringVar(&growthTime, "time", "0s", "time to reach the size of memory you allocated")
	flag.IntVar(&workers, "workers", 1, "number of workers allocating memory")
	flag.BoolVar(&client, "client", false, "the process runs as a client")
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
		progress := float64(now.Sub(startTime)) / float64(timeLine)
		expected := uint64(float64(length)*progress) / pageSize

		for i := allocated; uint64(i) < expected; i++ {
			data[uint64(i)*pageSize] = 0
		}

		allocated = expected
		if now.Equal(endTime) {
			break
		}
		time.Sleep(min(interval, time.Until(endTime)))
	}

}

func run(length uint64, timeLine time.Duration) error {
	data, err := syscall.Mmap(-1, 0, int(length), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		return fmt.Errorf("allocate %d bytes: %w", length, err)
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
	flag.Parse()
	length, err := humanize.ParseBytes(memSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid --size %q: expected a byte size such as 64MiB (%v)\n", memSize, err)
		os.Exit(2)
	}
	timeLine, err := time.ParseDuration(growthTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid --time %q: %v\n", growthTime, err)
		os.Exit(2)
	}

	if !client {
		workQueue := make(chan struct{}, workers)
		initialWorkers := workers
		for {
			workQueue <- struct{}{}
			if initialWorkers > 0 {
				initialWorkers--
			} else {
				time.Sleep(time.Second)
			}
			go func() {
				cmd := exec.Command("memStress", "--size", memSize, "--workers", fmt.Sprintf("%d", workers),
					"--time", growthTime, "--client", "1")
				cmd.Stderr = os.Stderr
				cmd.SysProcAttr = &syscall.SysProcAttr{
					Pdeathsig: syscall.SIGTERM,
				}
				if err := cmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "memStress worker: %v\n", err)
				}
				<-workQueue
			}()
		}
	} else if err := run(length/uint64(workers), timeLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
