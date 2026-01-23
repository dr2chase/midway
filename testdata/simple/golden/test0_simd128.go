// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !midway && !ignore

package testdata

import (
	"fmt"
	"simd/archsimd"

	"github.com/dr2chase/midway/midway"
)

// A SIMD-dependent type alias
type MyInt8s_simd128 = archsimd.

	// A struct dependent on SIMD
	Int8x16

type VectorC_simd128 struct {
	Data archsimd.Float64x2
}

// Two SIMD-dependent global variables
var ZeroVec_simd128 archsimd.Int32x4
var SomeVec_simd128 archsimd.

	// A function using SIMD types directly
	Int32x4

func Add_simd128(a, b archsimd.Int32x4) archsimd.Int32x4 {
	midway.Assert128(

	// A function using an alias of simd types
	)
	return a
}

func AddAlias_simd128(a, b MyInt8s_simd128) MyInt8s_simd128 {
	midway.Assert128(

	// A function type that mentions simd
	)
	return a
}

type Adder_simd128 func(a, b archsimd.Int32x4) archsimd.Int32x4

// A struct of a pair of simd pointers
type PtrPair_simd128 struct {
	a *archsimd.Float64x2
	b *archsimd.Uint32x4
}

// A dependent function with a dependent signature
func Process_simd128(v VectorC_simd128) {
	midway.Assert128()
	fmt.Println(v)
}

// A dependent function with a standard signature
func ComputeSum_simd128(n int) int {
	midway.
		// Uses SIMD internally
		Assert128()

	var v archsimd.Int32x4
	_ = v

	return n * 2
}

// A dependent function with a standard signature
func MentionsPPair_simd128(a any) any {
	midway.Assert128()
	x, _ := a.(PtrPair_simd128)
	return x.b
}

// A dependent function with a standard signature
func MentionsAdder_simd128(a any) any {
	midway.Assert128()
	x, _ := a.(Adder_simd128)
	return x
}

// A dependent function with a standard signature
func InAClosure_simd128(x any) any {
	midway.Assert128()
	f := func(y any) any {
		midway.Assert128()
		sx, _ := x.(archsimd.Int32x4)
		sy, _ := y.(archsimd.Int32x4)
		return Add_simd128(sx, sy)
	}
	return f
}

// A dependent function with a dependent signature
func (v *VectorC_simd128) MethodOfSimd_simd128() bool {
	midway.Assert128()
	return false
}

var vc_simd128 VectorC_simd128

// = &vc // this looks like a problem.

// init depends on vc
func init_simd128() {
	midway.Assert128()
	anyVc = &vc_simd128
}

// A dependent function with a standard signature
// that calls two dependent functions, one with a dependent
// signature, one without.
// Also assigns a dependent function value to a pointer.
func DepCallsDep_simd128() (x, y any, b bool) {
	midway.Assert128()
	x = Add_simd128(ZeroVec_simd128, SomeVec_simd128)
	y = MentionsAdder_simd128(Add_simd128)
	Fvar = InAClosure_simd128
	_, b = anyVc.(Vint)
	return
}

func depGeneric_simd128[T fmt.Stringer](x T) int {
	midway.Assert128()
	s := x.String()
	var bs []uint8 = []byte(s)
	v := archsimd.LoadUint8x16SlicePart(bs)
	v = v.Add(v)
	v.StoreSlice(bs)
	return int(bs[0])
}

// signature is not generic, but implementation is
func instGeneric_simd128(x int) int {
	midway.Assert128()
	return generic[archsimd.Int8x16](x)
}

// Caller that is NOT dependent calls one that IS dependent.
