// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"fmt"
	"simd_flex/simd"
)

// A SIMD-dependent type alias

// A struct dependent on SIMD

// Two SIMD-dependent global variables

// A function using SIMD types directly

// A function using an alias of simd types

// A function type that mentions simd

// A struct of a pair of simd pointers

// A dependent function with a dependent signature

// A dependent function with a standard signature
func ComputeSum(n int) int {
	switch
	// Uses SIMD internally
	simd.VectorSize() {
	case 128:
		return ComputeSum_simd128(n)
	case 256:
		return ComputeSum_simd256(n)
	default:
		panic("unsupported vector size")
	}
}

// A dependent function with a standard signature
func MentionsPPair(a any) any {
	switch simd.VectorSize() {
	case 128:
		return MentionsPPair_simd128(a)
	case 256:
		return MentionsPPair_simd256(a)
	default:
		panic("unsupported vector size")
	}
}

// A dependent function with a standard signature
func MentionsAdder(a any) any {
	switch simd.VectorSize() {
	case 128:
		return MentionsAdder_simd128(a)
	case 256:
		return MentionsAdder_simd256(a)
	default:
		panic("unsupported vector size")
	}
}

// A dependent function with a standard signature
func InAClosure(x any) any {
	switch simd.VectorSize() {
	case 128:
		return InAClosure_simd128(x)
	case 256:
		return InAClosure_simd256(x)
	default:
		panic("unsupported vector size")
	}
}

type Ftype func(x any) any

var Fvar Ftype

// A dependent function with a dependent signature

type Vint interface {
	MethodOfSimd() bool
}

var anyVc any // = &vc // this looks like a problem.

// init depends on vc
func init() {
	switch simd.VectorSize(

	// A dependent function with a standard signature
	// that calls two dependent functions, one with a dependent
	// signature, one without.
	// Also assigns a dependent function value to a pointer.
	) {
	case 128:
		{
			init_simd128()
			return
		}
	case 256:
		{
			init_simd256()
			return
		}
	default:
		panic("unsupported vector size")
	}
}

func DepCallsDep() (x, y any, b bool) {
	switch simd.VectorSize() {
	case 128:
		return DepCallsDep_simd128()
	case 256:
		return DepCallsDep_simd256()
	default:
		panic(

			// Caller that is NOT dependent calls one that IS dependent.
			"unsupported vector size")
	}
}

func MainCaller() {
	res := ComputeSum(10)
	fmt.Println(res)
}
