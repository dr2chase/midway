// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build midway

package main

import (
	"fmt"
	"github.com/dr2chase/midway/simd"
	"simd/archsimd"
)

var sumWidth int

func main() {
	var a, b [50]float32
	for i := 0; i < 50; i++ {
		a[i] = float32(i)
		b[i] = float32(i)
	}
	println(ip(a[:5], b[:5])) // use println so fmt not imported into main_simd.go
	println(ip(a[:10], b[:10]))
	println(ip(a[:20], b[:20]))
	println(ip(a[:30], b[:30]))
	println(ip(a[:40], b[:40]))
	println(ip(a[:50], b[:50]))

	println("sum was computed in width", sumWidth)
}

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
	panic(fmt.Errorf("Unexpected simd type %T", x))
}

func ip(x, y []float32) float32 {
	var a simd.Float32s
	var i int
	for i = 0; i < len(x)-a.Len()+1; i += a.Len() {
		u := simd.LoadFloat32Slice(x[i : i+a.Len()])
		v := simd.LoadFloat32Slice(y[i : i+a.Len()])
		a = a.Add(u.Mul(v))
	}
	if i < len(x) {
		a = a.Add(simd.LoadFloat32SlicePart(x[i:]).
			Mul(simd.LoadFloat32SlicePart(y[i:])))
	}

	return sum(a)
}
