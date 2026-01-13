//go:build !midway && !ignore

package testdata

import (
	"fmt"

	// A SIMD-dependent type alias
	archsimd "simd/archsimd"
)

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
	return a
}

// A function using an alias of simd types
func AddAlias_simd512(a, b MyInt8s_simd512) MyInt8s_simd512 {
	return a
}

// A function type that mentions simd
type Adder_simd512 func(a, b archsimd.Int32x16) archsimd.Int32x16

// A struct of a pair of simd pointers
type PtrPair_simd512 struct {
	a *archsimd.Float64x8
	b *archsimd.Uint32x16
}

// A dependent function with a dependent signature
func Process_simd512(v VectorC_simd512) {
	fmt.Println(v)
}

// A dependent function with a standard signature
func ComputeSum_simd512(n int) int {
	// Uses SIMD internally
	var v archsimd.Int32x16
	_ = v

	return n * 2
}

// A dependent function with a standard signature
func MentionsPPair_simd512(a any) any {
	x, _ := a.(PtrPair_simd512)
	return x.b
}

// A dependent function with a standard signature
func MentionsAdder_simd512(a any) any {
	x, _ := a.(Adder_simd512)
	return x
}

// A dependent function with a standard signature
func InAClosure_simd512(x any) any {
	f := func(y any) any {
		sx, _ := x.(archsimd.Int32x16)
		sy, _ := y.(archsimd.Int32x16)
		return Add_simd512(sx, sy)
	}
	return f
}

// A dependent function with a dependent signature
func (v *VectorC_simd512) MethodOfSimd_simd512() bool {
	return false
}

var vc_simd512 VectorC_simd512

// = &vc // this looks like a problem.

// init depends on vc
func init_simd512() {
	anyVc = &vc_simd512
}

// A dependent function with a standard signature
// that calls two dependent functions, one with a dependent
// signature, one without.
// Also assigns a dependent function value to a pointer.
func DepCallsDep_simd512() (x, y any, b bool) {
	x = Add_simd512(ZeroVec_simd512, SomeVec_simd512)
	y = MentionsAdder_simd512(Add_simd512)
	Fvar = InAClosure_simd512
	_, b = anyVc.(Vint)
	return
}

func depGeneric_simd512[T fmt.Stringer](x T) int {
	s := x.String()
	var bs []uint8 = []byte(s)
	v := archsimd.LoadUint8x64SlicePart(bs)
	v = v.Add(v)
	v.StoreSlice(bs)
	return int(bs[0])
}

// signature is not generic, but implementation is
func instGeneric_simd512(x int) int {
	return generic[archsimd.Int8x64](x)
}

// Caller that is NOT dependent calls one that IS dependent.
