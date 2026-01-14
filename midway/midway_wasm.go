// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package midway

// MaxVectorSize returns the bit length of the longest vector available
// on the current hardware.  For wasm, this is 128.
func MaxVectorSize() int {
	return 128
}

func Assert512() {
	panic("vector length is not 512")
}

func Assert256() {
	panic("vector length is not 256")
}

func Assert128() {

}

func Assert2048() {
	// Arm64/SVE vectors can be this large, in theory.
	panic("vector length is not 2048")
}

func Assert65536() {
	// RISCV vectors can be this large, in theory.
	panic("vector length is not 65536")
}
