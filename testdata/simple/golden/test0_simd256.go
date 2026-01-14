// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !midway && !ignore

package testdata

import (
	"fmt"

	// A SIMD-dependent type alias
	archsimd "simd/archsimd"
)

type MyInt8s_simd256 = archsimd.

	// A struct dependent on SIMD
	Int8x32

type VectorC_simd256 struct {
	Data archsimd.Float64x4
}

// Two SIMD-dependent global variables
var ZeroVec_simd256 archsimd.Int32x8
var SomeVec_simd256 archsimd.

	// A function using SIMD types directly
	Int32x8

func Add_simd256(a, b archsimd.Int32x8) archsimd.Int32x8 {
	return a
}

// A function using an alias of simd types
func AddAlias_simd256(a, b MyInt8s_simd256) MyInt8s_simd256 {
	return a
}

// A function type that mentions simd
type Adder_simd256 func(a, b archsimd.Int32x8) archsimd.Int32x8

// A struct of a pair of simd pointers
type PtrPair_simd256 struct {
	a *archsimd.Float64x4
	b *archsimd.Uint32x8
}

// A dependent function with a dependent signature
func Process_simd256(v VectorC_simd256) {
	fmt.Println(v)
}

// A dependent function with a standard signature
func ComputeSum_simd256(n int) int {
	// Uses SIMD internally
	var v archsimd.Int32x8
	_ = v

	return n * 2
}

// A dependent function with a standard signature
func MentionsPPair_simd256(a any) any {
	x, _ := a.(PtrPair_simd256)
	return x.b
}

// A dependent function with a standard signature
func MentionsAdder_simd256(a any) any {
	x, _ := a.(Adder_simd256)
	return x
}

// A dependent function with a standard signature
func InAClosure_simd256(x any) any {
	f := func(y any) any {
		sx, _ := x.(archsimd.Int32x8)
		sy, _ := y.(archsimd.Int32x8)
		return Add_simd256(sx, sy)
	}
	return f
}

// A dependent function with a dependent signature
func (v *VectorC_simd256) MethodOfSimd_simd256() bool {
	return false
}

var vc_simd256 VectorC_simd256

// = &vc // this looks like a problem.

// init depends on vc
func init_simd256() {
	anyVc = &vc_simd256
}

// A dependent function with a standard signature
// that calls two dependent functions, one with a dependent
// signature, one without.
// Also assigns a dependent function value to a pointer.
func DepCallsDep_simd256() (x, y any, b bool) {
	x = Add_simd256(ZeroVec_simd256, SomeVec_simd256)
	y = MentionsAdder_simd256(Add_simd256)
	Fvar = InAClosure_simd256
	_, b = anyVc.(Vint)
	return
}

func depGeneric_simd256[T fmt.Stringer](x T) int {
	s := x.String()
	var bs []uint8 = []byte(s)
	v := archsimd.LoadUint8x32SlicePart(bs)
	v = v.Add(v)
	v.StoreSlice(bs)
	return int(bs[0])
}

// signature is not generic, but implementation is
func instGeneric_simd256(x int) int {
	return generic[archsimd.Int8x32](x)
}

// Caller that is NOT dependent calls one that IS dependent.
