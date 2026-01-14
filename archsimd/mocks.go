// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package archsimd

type Int8x16 struct { v [16]int8 }
type Int16x8 struct { v [8]int16 }
type Int32x4 struct { v [4]int32 }
type Int64x2 struct { v [2]int64 }
type Uint8x16 struct { v [16]uint8 }
type Uint16x8 struct { v [8]uint16 }
type Uint32x4 struct { v [4]uint32 }
type Uint64x2 struct { v [2]uint64 }
type Float32x4 struct { v [4]float32 }
type Float64x2 struct { v [2]float64 }

type Int8x32 struct { v [32]int8 }
type Int16x16 struct { v [16]int16 }
type Int32x8 struct { v [8]int32 }
type Int64x4 struct { v [4]int64 }
type Uint8x32 struct { v [32]uint8 }
type Uint16x16 struct { v [16]uint16 }
type Uint32x8 struct { v [8]uint32 }
type Uint64x4 struct { v [4]uint64 }
type Float32x8 struct { v [8]float32 }
type Float64x4 struct { v [4]float64 }

type Int8x64 struct { v [64]int8 }
type Int16x32 struct { v [32]int16 }
type Int32x16 struct { v [16]int32 }
type Int64x8 struct { v [8]int64 }
type Uint8x64 struct { v [64]uint8 }
type Uint16x32 struct { v [32]uint16 }
type Uint32x16 struct { v [16]uint32 }
type Uint64x8 struct { v [8]uint64 }
type Float32x16 struct { v [16]float32 }
type Float64x8 struct { v [8]float64 }

type Mask8x16 struct { v [16]int8 }
type Mask16x8 struct { v [8]int16 }
type Mask32x4 struct { v [4]int32 }
type Mask64x2 struct { v [2]int64 }
type Mask8x32 struct { v [32]int8 }
type Mask16x16 struct { v [16]int16 }
type Mask32x8 struct { v [8]int32 }
type Mask64x4 struct { v [4]int64 }
type Mask8x64 struct { v [64]int8 }
type Mask16x32 struct { v [32]int16 }
type Mask32x16 struct { v [16]int32 }
type Mask64x8 struct { v [8]int64 }
