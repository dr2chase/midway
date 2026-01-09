// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simd

func VectorSize() int {
	return 128
}

type Int8s []int8
type Int16s []int16
type Int32s []int32
type Int64s []int64
type Uint8s []uint8
type Uint16s []uint16
type Uint32s []uint32
type Uint64s []uint64
type Float32s []float32
type Float64s []float64

// K=128
type Int8s_simd128 []int8
type Int16s_simd128 []int16
type Int32s_simd128 []int32
type Int64s_simd128 []int64
type Uint8s_simd128 []uint8
type Uint16s_simd128 []uint16
type Uint32s_simd128 []uint32
type Uint64s_simd128 []uint64
type Float32s_simd128 []float32
type Float64s_simd128 []float64

// K=256
type Int8s_simd256 []int8
type Int16s_simd256 []int16
type Int32s_simd256 []int32
type Int64s_simd256 []int64
type Uint8s_simd256 []uint8
type Uint16s_simd256 []uint16
type Uint32s_simd256 []uint32
type Uint64s_simd256 []uint64
type Float32s_simd256 []float32
type Float64s_simd256 []float64

// K=512
type Int8s_simd512 []int8
type Int16s_simd512 []int16
type Int32s_simd512 []int32
type Int64s_simd512 []int64
type Uint8s_simd512 []uint8
type Uint16s_simd512 []uint16
type Uint32s_simd512 []uint32
type Uint64s_simd512 []uint64
type Float32s_simd512 []float32
type Float64s_simd512 []float64

// K=2048
type Int8s_simd2048 []int8
type Int16s_simd2048 []int16
type Int32s_simd2048 []int32
type Int64s_simd2048 []int64
type Uint8s_simd2048 []uint8
type Uint16s_simd2048 []uint16
type Uint32s_simd2048 []uint32
type Uint64s_simd2048 []uint64
type Float32s_simd2048 []float32
type Float64s_simd2048 []float64
