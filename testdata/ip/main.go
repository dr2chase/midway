// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build midway

package main

import (
	"fmt"
	"github.com/dr2chase/midway/simd"
)

func main() {
	var a, b [50]float32
	for i := 0; i < 50; i++ {
		a[i] = float32(i)
		b[i] = float32(i)
	}
	fmt.Println(ip(a[:5], b[:5]))
	fmt.Println(ip(a[:10], b[:10]))
	fmt.Println(ip(a[:20], b[:20]))
	fmt.Println(ip(a[:30], b[:30]))
	fmt.Println(ip(a[:40], b[:40]))
	fmt.Println(ip(a[:50], b[:50]))
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

	s := make([]float32, a.Len())
	a.StoreSlice(s)
	var r float32
	for _, e := range s {
		r += e
	}
	return r

	// Would like to do this but methods are not
	// defined for all the types we need, sigh.
	// l := a.Len()
	// for l > 1 {
	// 	a = a.AddPairs(a)
	// 	l /= 2
	// }

	// return a.GetElem(0)
}
