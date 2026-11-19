// Copyright 2026 Chaos Mesh Authors.
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
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRunReturnsMappingError(t *testing.T) {
	if err := run(0, 0); !errors.Is(err, syscall.EINVAL) {
		t.Fatalf("mapping zero bytes: got %v, want EINVAL", err)
	}
}

func TestLinearGrowSubMillisecond(t *testing.T) {
	pageSize := os.Getpagesize()
	data := make([]byte, 2*pageSize)
	for i := 0; i < len(data); i += pageSize {
		data[i] = 1
	}

	linearGrow(data, uint64(len(data)), 500*time.Microsecond)

	for i := 0; i < len(data); i += pageSize {
		if data[i] != 0 {
			t.Errorf("page %d was not touched", i/pageSize)
		}
	}
}
