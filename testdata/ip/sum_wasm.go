// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build midway && wasm

package main

import (
	"github.com/dr2chase/midway/simd"
	"simd/archsimd"
)

func sum(x simd.Float32s) float32 {
	switch a := (any(x)).(type) {
	case archsimd.Float32x4:
		sumWidth = 128
		s := make([]float32, a.Len())
		a.StoreSlice(s)
		var r float32
		for _, e := range s {
			r += e
		}
		return r
	}
	panic("not a known type")
}
