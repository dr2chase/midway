// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build midway && amd64

package main

import (
	"github.com/dr2chase/midway/simd"
	"simd/archsimd"
)

func sum(x simd.Float32s) float32 {

	switch a := (any(x)).(type) {
	case archsimd.Float32x8:
		sumWidth = 256
		a = a.AddPairsGrouped(a)
		a = a.AddPairsGrouped(a)
		return a.GetLo().GetElem(0) + a.GetHi().GetElem(0)
	case archsimd.Float32x16:
		sumWidth = 512
		s := make([]float32, a.Len())
		a.StoreSlice(s)
		var r float32
		for _, e := range s {
			r += e
		}
		return r
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
