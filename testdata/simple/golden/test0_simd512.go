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
type MyInt8s_simd512 = archsimd.

	// A struct dependent on SIMD
	Int8x64

type VectorC_simd512 struct {
	Data archsimd.Float64x8
}

// Two SIMD-dependent global variables
var ZeroVec_simd512 archsimd.Int32x16
var SomeVec_simd512 archsimd.

	// A function using SIMD types directly
	Int32x16

func Add_simd512(a, b archsimd.Int32x16) archsimd.Int32x16 {
	midway.Assert512(

	// A function using an alias of simd types
	)
	return a
}

func AddAlias_simd512(a, b MyInt8s_simd512) MyInt8s_simd512 {
	midway.Assert512(

	// A function type that mentions simd
	)
	return a
}

type Adder_simd512 func(a, b archsimd.Int32x16) archsimd.Int32x16

// A struct of a pair of simd pointers
type PtrPair_simd512 struct {
	a *archsimd.Float64x8
	b *archsimd.Uint32x16
}

// A dependent function with a dependent signature
func Process_simd512(v VectorC_simd512) {
	midway.Assert512()
	fmt.Println(v)
}

// A dependent function with a standard signature
func ComputeSum_simd512(n int) int {
	midway.
		// Uses SIMD internally
		Assert512()

	var v archsimd.Int32x16
	_ = v

	return n * 2
}

// A dependent function with a standard signature
func MentionsPPair_simd512(a any) any {
	midway.Assert512()
	x, _ := a.(PtrPair_simd512)
	return x.b
}

// A dependent function with a standard signature
func MentionsAdder_simd512(a any) any {
	midway.Assert512()
	x, _ := a.(Adder_simd512)
	return x
}

// A dependent function with a standard signature
func InAClosure_simd512(x any) any {
	midway.Assert512()
	f := func(y any) any {
		midway.Assert512()
		sx, _ := x.(archsimd.Int32x16)
		sy, _ := y.(archsimd.Int32x16)
		return Add_simd512(sx, sy)
	}
	return f
}

// A dependent function with a dependent signature
func (v *VectorC_simd512) MethodOfSimd_simd512() bool {
	midway.Assert512()
	return false
}

var vc_simd512 VectorC_simd512

// = &vc // this looks like a problem.

// init depends on vc
func init_simd512() {
	midway.Assert512()
	anyVc = &vc_simd512
}

// A dependent function with a standard signature
// that calls two dependent functions, one with a dependent
// signature, one without.
// Also assigns a dependent function value to a pointer.
func DepCallsDep_simd512() (x, y any, b bool) {
	midway.Assert512()
	x = Add_simd512(ZeroVec_simd512, SomeVec_simd512)
	y = MentionsAdder_simd512(Add_simd512)
	Fvar = InAClosure_simd512
	_, b = anyVc.(Vint)
	return
}

func depGeneric_simd512[T fmt.Stringer](x T) int {
	midway.Assert512()
	s := x.String()
	var bs []uint8 = []byte(s)
	v := archsimd.LoadUint8x64SlicePart(bs)
	v = v.Add(v)
	v.StoreSlice(bs)
	return int(bs[0])
}

// signature is not generic, but implementation is
func instGeneric_simd512(x int) int {
	midway.Assert512()
	return generic[archsimd.Int8x64](x)
}

// Caller that is NOT dependent calls one that IS dependent.
