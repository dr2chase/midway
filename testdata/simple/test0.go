// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build midway && !ignore

package testdata

import (
	"fmt"
	"github.com/dr2chase/midway/simd"
)

// A SIMD-dependent type alias
type MyInt8s = simd.Int8s

// A struct dependent on SIMD
type VectorC struct {
	Data simd.Float64s
}

// Two SIMD-dependent global variables
var ZeroVec simd.Int32s
var SomeVec simd.Int32s

// A function using SIMD types directly
func Add(a, b simd.Int32s) simd.Int32s {
	return a
}

// A function using an alias of simd types
func AddAlias(a, b MyInt8s) MyInt8s {
	return a
}

// A function type that mentions simd
type Adder func(a, b simd.Int32s) simd.Int32s

// A struct of a pair of simd pointers
type PtrPair struct {
	a *simd.Float64s
	b *simd.Uint32s
}

// A dependent function with a dependent signature
func Process(v VectorC) {
	fmt.Println(v)
}

// A dependent function with a standard signature
func ComputeSum(n int) int {
	// Uses SIMD internally
	var v simd.Int32s
	_ = v

	return n * 2
}

// A dependent function with a standard signature
func MentionsPPair(a any) any {
	x, _ := a.(PtrPair)
	return x.b
}

// A dependent function with a standard signature
func MentionsAdder(a any) any {
	x, _ := a.(Adder)
	return x
}

// A dependent function with a standard signature
func InAClosure(x any) any {
	f := func(y any) any {
		sx, _ := x.(simd.Int32s)
		sy, _ := y.(simd.Int32s)
		return Add(sx, sy)
	}
	return f
}

type Ftype func(x any) any

var Fvar Ftype

// A dependent function with a dependent signature
func (v *VectorC) MethodOfSimd() bool {
	return false
}

type Vint interface {
	MethodOfSimd() bool
}

var vc VectorC
var anyVc any // = &vc // this looks like a problem.

// init depends on vc
func init() {
	anyVc = &vc
}

// A dependent function with a standard signature
// that calls two dependent functions, one with a dependent
// signature, one without.
// Also assigns a dependent function value to a pointer.
func DepCallsDep() (x, y any, b bool) {
	x = Add(ZeroVec, SomeVec)
	y = MentionsAdder(Add)
	Fvar = InAClosure
	_, b = anyVc.(Vint)
	return
}

type haslen interface {
	Len() int
}

func generic[T haslen](x int) int {
	var v T
	return x + v.Len()
}

func depGeneric[T fmt.Stringer](x T) int {
	s := x.String()
	var bs []uint8 = []byte(s)
	v := simd.LoadUint8SlicePart(bs)
	v = v.Add(v)
	v.StoreSlice(bs)
	return int(bs[0])
}

// signature is not generic, but implementation is
func instGeneric(x int) int {
	return generic[simd.Int8s](x)
}

// Caller that is NOT dependent calls one that IS dependent.
func MainCaller() {
	res := ComputeSum(10)
	fmt.Println(res + instGeneric(11))
}
