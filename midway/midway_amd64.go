// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package midway

import "simd/archsimd"

// MaxVectorSize returns the bit length of the longest vector available
// on the current hardware.  For AVX and AVX2, this is 256, for AVX512, 512.
func MaxVectorSize() int {
	if archsimd.X86.AVX512() {
		return 512
	}
	if archsimd.X86.AVX() {
		return 256
	}
	return 0
}

func Assert512() {
	if !archsimd.X86.AVX512() {
		panic("vector length is not 512")
	}
}

func Assert256() {
	if !archsimd.X86.AVX() {
		panic("vector length is not 256")
	}
}

func Assert128() {
	// 128 is Arm64/Neon, WASM, PPC64, Loong64/LSX, maybe others.
	panic("vector length is not 128")
}

func Assert2048() {
	// Arm64/SVE vectors can be this large, in theory.
	panic("vector length is not 2048")
}

func Assert65536() {
	// RISCV vectors can be this large, in theory.
	panic("vector length is not 65536")
}
