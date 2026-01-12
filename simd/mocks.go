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

type Mask8s []int8
type Mask16s []int16
type Mask32s []int32
type Mask64s []int64

func LoadInt8Slice(s []int8) Int8s
func LoadInt16Slice(s []int16) Int16s
func LoadInt32Slice(s []int32) Int32s
func LoadInt64Slice(s []int64) Int64s
func LoadUint8Slice(s []uint8) Uint8s
func LoadUint16Slice(s []uint16) Uint16s
func LoadUint32Slice(s []uint32) Uint32s
func LoadUint64Slice(s []uint64) Uint64s
func LoadFloat32Slice(s []float32) Float32s
func LoadFloat64Slice(s []float64) Float64s
func LoadMask8Slice(s []int8) Mask8s
func LoadMask16Slice(s []int16) Mask16s
func LoadMask32Slice(s []int32) Mask32s
func LoadMask64Slice(s []int64) Mask64s

func LoadInt8SlicePart(s []int8) Int8s
func LoadInt16SlicePart(s []int16) Int16s
func LoadInt32SlicePart(s []int32) Int32s
func LoadInt64SlicePart(s []int64) Int64s
func LoadUint8SlicePart(s []uint8) Uint8s
func LoadUint16SlicePart(s []uint16) Uint16s
func LoadUint32SlicePart(s []uint32) Uint32s
func LoadUint64SlicePart(s []uint64) Uint64s
func LoadFloat32SlicePart(s []float32) Float32s
func LoadFloat64SlicePart(s []float64) Float64s
func LoadMask8SlicePart(s []int8) Mask8s
func LoadMask16SlicePart(s []int16) Mask16s
func LoadMask32SlicePart(s []int32) Mask32s
func LoadMask64SlicePart(s []int64) Mask64s

func (x Int8s) StoreSlice(s []int8)
func (x Int16s) StoreSlice(s []int16)
func (x Int32s) StoreSlice(s []int32)
func (x Int64s) StoreSlice(s []int64)
func (x Uint8s) StoreSlice(s []uint8)
func (x Uint16s) StoreSlice(s []uint16)
func (x Uint32s) StoreSlice(s []uint32)
func (x Uint64s) StoreSlice(s []uint64)
func (x Float32s) StoreSlice(s []float32)
func (x Float64s) StoreSlice(s []float64)
func (x Mask8s) StoreSlice(s []int8)
func (x Mask16s) StoreSlice(s []int16)
func (x Mask32s) StoreSlice(s []int32)
func (x Mask64s) StoreSlice(s []int64)

func (x Int8s) StoreSlicePart(s []int8)
func (x Int16s) StoreSlicePart(s []int16)
func (x Int32s) StoreSlicePart(s []int32)
func (x Int64s) StoreSlicePart(s []int64)
func (x Uint8s) StoreSlicePart(s []uint8)
func (x Uint16s) StoreSlicePart(s []uint16)
func (x Uint32s) StoreSlicePart(s []uint32)
func (x Uint64s) StoreSlicePart(s []uint64)
func (x Float32s) StoreSlicePart(s []float32)
func (x Float64s) StoreSlicePart(s []float64)
func (x Mask8s) StoreSlicePart(s []int8)
func (x Mask16s) StoreSlicePart(s []int16)
func (x Mask32s) StoreSlicePart(s []int32)
func (x Mask64s) StoreSlicePart(s []int64)

func (x Int8s) Add(y Int8s) Int8s
func (x Int16s) Add(y Int16s) Int16s
func (x Int32s) Add(y Int32s) Int32s
func (x Int64s) Add(y Int64s) Int64s
func (x Uint8s) Add(y Uint8s) Uint8s
func (x Uint16s) Add(y Uint16s) Uint16s
func (x Uint32s) Add(y Uint32s) Uint32s
func (x Uint64s) Add(y Uint64s) Uint64s
func (x Float32s) Add(y Float32s) Float32s
func (x Float64s) Add(y Float64s) Float64s

func (x Int8s) Sub(y Int8s) Int8s
func (x Int16s) Sub(y Int16s) Int16s
func (x Int32s) Sub(y Int32s) Int32s
func (x Int64s) Sub(y Int64s) Int64s
func (x Uint8s) Sub(y Uint8s) Uint8s
func (x Uint16s) Sub(y Uint16s) Uint16s
func (x Uint32s) Sub(y Uint32s) Uint32s
func (x Uint64s) Sub(y Uint64s) Uint64s
func (x Float32s) Sub(y Float32s) Float32s
func (x Float64s) Sub(y Float64s) Float64s

func (x Int8s) Mul(y Int8s) Int8s
func (x Int16s) Mul(y Int16s) Int16s
func (x Int32s) Mul(y Int32s) Int32s
func (x Int64s) Mul(y Int64s) Int64s
func (x Uint8s) Mul(y Uint8s) Uint8s
func (x Uint16s) Mul(y Uint16s) Uint16s
func (x Uint32s) Mul(y Uint32s) Uint32s
func (x Uint64s) Mul(y Uint64s) Uint64s
func (x Float32s) Mul(y Float32s) Float32s
func (x Float64s) Mul(y Float64s) Float64s

func (x Int8s) Div(y Int8s) Int8s
func (x Int16s) Div(y Int16s) Int16s
func (x Int32s) Div(y Int32s) Int32s
func (x Int64s) Div(y Int64s) Int64s
func (x Uint8s) Div(y Uint8s) Uint8s
func (x Uint16s) Div(y Uint16s) Uint16s
func (x Uint32s) Div(y Uint32s) Uint32s
func (x Uint64s) Div(y Uint64s) Uint64s
func (x Float32s) Div(y Float32s) Float32s
func (x Float64s) Div(y Float64s) Float64s

func (x Int8s) And(y Int8s) Int8s
func (x Int16s) And(y Int16s) Int16s
func (x Int32s) And(y Int32s) Int32s
func (x Int64s) And(y Int64s) Int64s
func (x Uint8s) And(y Uint8s) Uint8s
func (x Uint16s) And(y Uint16s) Uint16s
func (x Uint32s) And(y Uint32s) Uint32s
func (x Uint64s) And(y Uint64s) Uint64s

func (x Int8s) Or(y Int8s) Int8s
func (x Int16s) Or(y Int16s) Int16s
func (x Int32s) Or(y Int32s) Int32s
func (x Int64s) Or(y Int64s) Int64s
func (x Uint8s) Or(y Uint8s) Uint8s
func (x Uint16s) Or(y Uint16s) Uint16s
func (x Uint32s) Or(y Uint32s) Uint32s
func (x Uint64s) Or(y Uint64s) Uint64s

func (x Int8s) Xor(y Int8s) Int8s
func (x Int16s) Xor(y Int16s) Int16s
func (x Int32s) Xor(y Int32s) Int32s
func (x Int64s) Xor(y Int64s) Int64s
func (x Uint8s) Xor(y Uint8s) Uint8s
func (x Uint16s) Xor(y Uint16s) Uint16s
func (x Uint32s) Xor(y Uint32s) Uint32s
func (x Uint64s) Xor(y Uint64s) Uint64s

func (x Int8s) AndNot(y Int8s) Int8s
func (x Int16s) AndNot(y Int16s) Int16s
func (x Int32s) AndNot(y Int32s) Int32s
func (x Int64s) AndNot(y Int64s) Int64s
func (x Uint8s) AndNot(y Uint8s) Uint8s
func (x Uint16s) AndNot(y Uint16s) Uint16s
func (x Uint32s) AndNot(y Uint32s) Uint32s
func (x Uint64s) AndNot(y Uint64s) Uint64s

func (x Int8s) Not() Int8s
func (x Int16s) Not() Int16s
func (x Int32s) Not() Int32s
func (x Int64s) Not() Int64s
func (x Uint8s) Not() Uint8s
func (x Uint16s) Not() Uint16s
func (x Uint32s) Not() Uint32s
func (x Uint64s) Not() Uint64s
func (x Float32s) Not() Float32s
func (x Float64s) Not() Float64s

func (x Int8s) Neg() Int8s
func (x Int16s) Neg() Int16s
func (x Int32s) Neg() Int32s
func (x Int64s) Neg() Int64s
func (x Float32s) Neg() Float32s
func (x Float64s) Neg() Float64s

func (x Int8s) Abs() Int8s
func (x Int16s) Abs() Int16s
func (x Int32s) Abs() Int32s
func (x Int64s) Abs() Int64s
func (x Float32s) Abs() Float32s
func (x Float64s) Abs() Float64s

func (x Int8s) Equal(y Int8s) Mask8s
func (x Int16s) Equal(y Int16s) Mask16s
func (x Int32s) Equal(y Int32s) Mask32s
func (x Int64s) Equal(y Int64s) Mask64s
func (x Uint8s) Equal(y Uint8s) Mask8s
func (x Uint16s) Equal(y Uint16s) Mask16s
func (x Uint32s) Equal(y Uint32s) Mask32s
func (x Uint64s) Equal(y Uint64s) Mask64s
func (x Float32s) Equal(y Float32s) Mask32s
func (x Float64s) Equal(y Float64s) Mask64s

func (x Int8s) NotEqual(y Int8s) Mask8s
func (x Int16s) NotEqual(y Int16s) Mask16s
func (x Int32s) NotEqual(y Int32s) Mask32s
func (x Int64s) NotEqual(y Int64s) Mask64s
func (x Uint8s) NotEqual(y Uint8s) Mask8s
func (x Uint16s) NotEqual(y Uint16s) Mask16s
func (x Uint32s) NotEqual(y Uint32s) Mask32s
func (x Uint64s) NotEqual(y Uint64s) Mask64s
func (x Float32s) NotEqual(y Float32s) Mask32s
func (x Float64s) NotEqual(y Float64s) Mask64s

func (x Int8s) Less(y Int8s) Mask8s
func (x Int16s) Less(y Int16s) Mask16s
func (x Int32s) Less(y Int32s) Mask32s
func (x Int64s) Less(y Int64s) Mask64s
func (x Uint8s) Less(y Uint8s) Mask8s
func (x Uint16s) Less(y Uint16s) Mask16s
func (x Uint32s) Less(y Uint32s) Mask32s
func (x Uint64s) Less(y Uint64s) Mask64s
func (x Float32s) Less(y Float32s) Mask32s
func (x Float64s) Less(y Float64s) Mask64s

func (x Int8s) LessEqual(y Int8s) Mask8s
func (x Int16s) LessEqual(y Int16s) Mask16s
func (x Int32s) LessEqual(y Int32s) Mask32s
func (x Int64s) LessEqual(y Int64s) Mask64s
func (x Uint8s) LessEqual(y Uint8s) Mask8s
func (x Uint16s) LessEqual(y Uint16s) Mask16s
func (x Uint32s) LessEqual(y Uint32s) Mask32s
func (x Uint64s) LessEqual(y Uint64s) Mask64s
func (x Float32s) LessEqual(y Float32s) Mask32s
func (x Float64s) LessEqual(y Float64s) Mask64s

func (x Int8s) Greater(y Int8s) Mask8s
func (x Int16s) Greater(y Int16s) Mask16s
func (x Int32s) Greater(y Int32s) Mask32s
func (x Int64s) Greater(y Int64s) Mask64s
func (x Uint8s) Greater(y Uint8s) Mask8s
func (x Uint16s) Greater(y Uint16s) Mask16s
func (x Uint32s) Greater(y Uint32s) Mask32s
func (x Uint64s) Greater(y Uint64s) Mask64s
func (x Float32s) Greater(y Float32s) Mask32s
func (x Float64s) Greater(y Float64s) Mask64s

func (x Int8s) GreaterEqual(y Int8s) Mask8s
func (x Int16s) GreaterEqual(y Int16s) Mask16s
func (x Int32s) GreaterEqual(y Int32s) Mask32s
func (x Int64s) GreaterEqual(y Int64s) Mask64s
func (x Uint8s) GreaterEqual(y Uint8s) Mask8s
func (x Uint16s) GreaterEqual(y Uint16s) Mask16s
func (x Uint32s) GreaterEqual(y Uint32s) Mask32s
func (x Uint64s) GreaterEqual(y Uint64s) Mask64s
func (x Float32s) GreaterEqual(y Float32s) Mask32s
func (x Float64s) GreaterEqual(y Float64s) Mask64s

func (x Int8s) Min(y Int8s) Int8s
func (x Int16s) Min(y Int16s) Int16s
func (x Int32s) Min(y Int32s) Int32s
func (x Int64s) Min(y Int64s) Int64s
func (x Uint8s) Min(y Uint8s) Uint8s
func (x Uint16s) Min(y Uint16s) Uint16s
func (x Uint32s) Min(y Uint32s) Uint32s
func (x Uint64s) Min(y Uint64s) Uint64s
func (x Float32s) Min(y Float32s) Float32s
func (x Float64s) Min(y Float64s) Float64s

func (x Int8s) Max(y Int8s) Int8s
func (x Int16s) Max(y Int16s) Int16s
func (x Int32s) Max(y Int32s) Int32s
func (x Int64s) Max(y Int64s) Int64s
func (x Uint8s) Max(y Uint8s) Uint8s
func (x Uint16s) Max(y Uint16s) Uint16s
func (x Uint32s) Max(y Uint32s) Uint32s
func (x Uint64s) Max(y Uint64s) Uint64s
func (x Float32s) Max(y Float32s) Float32s
func (x Float64s) Max(y Float64s) Float64s

func (x Int8s) GetElem(index uint8) int8
func (x Int16s) GetElem(index uint8) int16
func (x Int32s) GetElem(index uint8) int32
func (x Int64s) GetElem(index uint8) int64
func (x Uint8s) GetElem(index uint8) uint8
func (x Uint16s) GetElem(index uint8) uint16
func (x Uint32s) GetElem(index uint8) uint32
func (x Uint64s) GetElem(index uint8) uint64
func (x Float32s) GetElem(index uint8) float32
func (x Float64s) GetElem(index uint8) float64

func (x Int8s) SetElem(index uint8, value int8) Int8s
func (x Int16s) SetElem(index uint8, value int16) Int16s
func (x Int32s) SetElem(index uint8, value int32) Int32s
func (x Int64s) SetElem(index uint8, value int64) Int64s
func (x Uint8s) SetElem(index uint8, value uint8) Uint8s
func (x Uint16s) SetElem(index uint8, value uint16) Uint16s
func (x Uint32s) SetElem(index uint8, value uint32) Uint32s
func (x Uint64s) SetElem(index uint8, value uint64) Uint64s
func (x Float32s) SetElem(index uint8, value float32) Float32s
func (x Float64s) SetElem(index uint8, value float64) Float64s

func (x Int16s) AddPairs(y Int16s) Int16s
func (x Int32s) AddPairs(y Int32s) Int32s
func (x Uint16s) AddPairs(y Uint16s) Uint16s
func (x Uint32s) AddPairs(y Uint32s) Uint32s
func (x Float32s) AddPairs(y Float32s) Float32s
func (x Float64s) AddPairs(y Float64s) Float64s

func (x Int8s) Len() int
func (x Int16s) Len() int
func (x Int32s) Len() int
func (x Int64s) Len() int
func (x Uint8s) Len() int
func (x Uint16s) Len() int
func (x Uint32s) Len() int
func (x Uint64s) Len() int
func (x Float32s) Len() int
func (x Float64s) Len() int
