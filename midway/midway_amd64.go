// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package midway

import (
	"fmt"
	"os"
	"simd/archsimd"
	"strconv"
)

var maxVectorSize int

func init() {
	actualMax := archMaxVectorSize()
	if gosimd := os.Getenv("GOSIMD"); gosimd != "" {
		val, err := strconv.Atoi(gosimd)
		if err != nil {
			panic(fmt.Errorf("Could not parse GOSIMD(='%s') as a decimal number, %v", gosimd, err))
		}
		if val > actualMax {
			panic(fmt.Errorf("Requested GOSIMD(='%d') is larger than the simd length (%d) supported on this cpu ", val, actualMax))
		}
		if val < 0 {
			panic(fmt.Errorf("Requested GOSIMD(='%d') is negative", val))
		}
		maxVectorSize = val
		return
	}
	maxVectorSize = actualMax
}

func MaxVectorSize() int {
	return maxVectorSize
}

// MaxVectorSize returns the bit length of the longest vector available
// on the current hardware.  For AVX and AVX2, this is 256, for AVX512, 512.
func archMaxVectorSize() int {
	if archsimd.X86.AVX512() {
		return 512
	}
	if archsimd.X86.AVX2() {
		return 256
	}
	// AVX has 256 bit float ops but only 128-bit integer ops
	// therefore it is 128.
	if archsimd.X86.AVX() {
		return 128
	}
	return 0
}

func Assert512() {
	if !archsimd.X86.AVX512() {
		panic("vector length is not 512")
	}
}

func Assert256() {
	if !archsimd.X86.AVX2() {
		panic("vector length is not 256")
	}
}

func Assert128() {
	if !archsimd.X86.AVX() {
		panic("vector length is not 128")
	}
}

func Assert2048() {
	// Arm64/SVE vectors can be this large, in theory.
	panic("vector length is not 2048")
}

func Assert65536() {
	// RISCV vectors can be this large, in theory.
	panic("vector length is not 65536")
}
