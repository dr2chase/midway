//go:build !midway && amd64

// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"fmt"

	"github.com/dr2chase/midway/midway"
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
	midway.MaxVectorSize() {
	case 128:
		return ComputeSum_simd128(n)
	case 256:
		return ComputeSum_simd256(n)
	case 512:
		return ComputeSum_simd512(n)
	default:
		panic("unsupported vector size")
	}
}

// A dependent function with a standard signature
func MentionsPPair(a any) any {
	switch midway.MaxVectorSize() {
	case 128:
		return MentionsPPair_simd128(a)
	case 256:
		return MentionsPPair_simd256(a)
	case 512:
		return MentionsPPair_simd512(a)
	default:
		panic("unsupported vector size")
	}
}

// A dependent function with a standard signature
func MentionsAdder(a any) any {
	switch midway.MaxVectorSize() {
	case 128:

		// A dependent function with a standard signature
		return MentionsAdder_simd128(a)
	case 256:
		return MentionsAdder_simd256(a)
	case 512:
		return MentionsAdder_simd512(a)
	default:
		panic("unsupported vector size")
	}
}

func InAClosure(x any) any {
	switch midway.MaxVectorSize() {
	case 128:
		return InAClosure_simd128(x)
	case 256:
		return InAClosure_simd256(x)
	case 512:
		return InAClosure_simd512(x)
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
	switch midway.MaxVectorSize(

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
	case 512:
		{
			init_simd512()
			return
		}
	default:
		panic("unsupported vector size")
	}
}

func DepCallsDep() (x, y any, b bool) {
	switch midway.MaxVectorSize() {
	case 128:
		return DepCallsDep_simd128()
	case 256:
		return DepCallsDep_simd256()
	case 512:
		return DepCallsDep_simd512()
	default:
		panic("unsupported vector size")
	}
}

type haslen interface {
	Len() int
}

func generic[T haslen](x int) int {
	var v T
	return x + v.Len()
}

func depGeneric[T fmt.Stringer](x T) int {
	switch midway.MaxVectorSize() {
	case 128:
		return depGeneric_simd128[T](x)
	case 256:
		return depGeneric_simd256[T](x)
	case 512:
		return depGeneric_simd512[T](x)
	default:
		panic("unsupported vector size")
	}
}

// signature is not generic, but implementation is
func instGeneric(x int) int {
	switch midway.MaxVectorSize() {
	case 128:

		// Caller that is NOT dependent calls one that IS dependent.
		return instGeneric_simd128(x)
	case 256:
		return instGeneric_simd256(x)
	case 512:
		return instGeneric_simd512(x)
	default:
		panic("unsupported vector size")
	}
}

func MainCaller() {
	res := ComputeSum(10)
	fmt.Println(res + instGeneric(11))
}
