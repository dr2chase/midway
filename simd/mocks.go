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

type Int8x16 = Int8s_simd128
type Int16x8 = Int16s_simd128
type Int32x4 = Int32s_simd128
type Int64x2 = Int64s_simd128
type Uint8x16 = Uint8s_simd128
type Uint16x8 = Uint16s_simd128
type Uint32x4 = Uint32s_simd128
type Uint64x2 = Uint64s_simd128
type Float32x4 = Float32s_simd128
type Float64x2 = Float64s_simd128

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

type Int8x32 = Int8s_simd256
type Int16x16 = Int16s_simd256
type Int32x8 = Int32s_simd256
type Int64x4 = Int64s_simd256
type Uint8x32 = Uint8s_simd256
type Uint16x16 = Uint16s_simd256
type Uint32x8 = Uint32s_simd256
type Uint64x4 = Uint64s_simd256
type Float32x8 = Float32s_simd256
type Float64x4 = Float64s_simd256

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

type Int8x64 = Int8s_simd512
type Int16x32 = Int16s_simd512
type Int32x16 = Int32s_simd512
type Int64x8 = Int64s_simd512
type Uint8x64 = Uint8s_simd512
type Uint16x32 = Uint16s_simd512
type Uint32x16 = Uint32s_simd512
type Uint64x8 = Uint64s_simd512
type Float32x16 = Float32s_simd512
type Float64x8 = Float64s_simd512

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

type Int8x256 = Int8s_simd2048
type Int16x128 = Int16s_simd2048
type Int32x64 = Int32s_simd2048
type Int64x32 = Int64s_simd2048
type Uint8x256 = Uint8s_simd2048
type Uint16x128 = Uint16s_simd2048
type Uint32x64 = Uint32s_simd2048
type Uint64x32 = Uint64s_simd2048
type Float32x64 = Float32s_simd2048
type Float64x32 = Float64s_simd2048
