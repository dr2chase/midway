// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build midway && !amd64

package main

import (
	"github.com/dr2chase/midway/simd"
)

func sum(x simd.Float32s) float32 {
	s := make([]float32, x.Len())
	x.StoreSlice(s)
	var r float32
	for _, e := range s {
		r += e
	}
	return r
}
