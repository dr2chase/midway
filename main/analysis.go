// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// Analyzer holds the state for SIMD dependency analysis
type Analyzer struct {
	pkg          *packages.Package
	simdTypes    map[types.Type]bool
	dependentObj map[types.Object]bool
}

func NewAnalyzer(pkg *packages.Package) *Analyzer {
	return &Analyzer{
		pkg:          pkg,
		simdTypes:    make(map[types.Type]bool),
		dependentObj: make(map[types.Object]bool),
	}
}

// Analyze builds the set of SIMD-dependent objects
func (a *Analyzer) Analyze() error {
	// Phase 1: Seed dependence from types and signatures

	// Scan all Defs
	for _, obj := range a.pkg.TypesInfo.Defs {
		if obj != nil {
			a.markIfDependent(obj)
		}
	}
	// Scan all Uses
	for _, obj := range a.pkg.TypesInfo.Uses {
		if obj != nil {
			a.markIfDependent(obj)
		}
	}

	// Phase 2: Transitive closure via function bodies
	// We need to iterate until no new objects are marked dependent.
	// Objects that can become dependent:
	// - Functions (if they use dependent objects)
	// - Variables (if initialized with dependent values? Types handle most vars, but maybe type inference?)
	// - Types (structs containing dependent fields - already handled by recursive type check)

	changed := true
	for changed {
		changed = false

		// Scan AST for function body dependencies
		for _, fileAST := range a.pkg.Syntax {
			for _, decl := range fileAST.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					obj := a.pkg.TypesInfo.ObjectOf(fn.Name)
					if obj == nil || a.dependentObj[obj] {
						continue // Already marked or unknown
					}

					if a.hasBodyDependency(fn) {
						a.dependentObj[obj] = true
						changed = true
					}
				}
			}
		}
	}

	return nil
}

func (a *Analyzer) hasBodyDependency(fn *ast.FuncDecl) bool {
	// Walk the body and check identifiers
	found := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if found {
			return false
		} // Early exit

		if id, ok := n.(*ast.Ident); ok {
			obj := a.pkg.TypesInfo.ObjectOf(id)
			if obj != nil {
				// If it's a variable or type that is dependent, we are dependent.
				if _, ok := obj.(*types.Func); !ok {
					if a.dependentObj[obj] {
						found = true
						return false
					}
				} else {
					// It is a function. Only propagate if it has a dependent signature.
					// If it has a clean signature, it will be dispatched, so we don't need to be specialized just to call it.
					// BUT wait, if we are NOT specialized, we can only call the dispatcher.
					// The dispatcher is the original function.
					// So if we call "ComputeSum", and we are "MainCaller" (Clean).
					// we remain "MainCaller" and call "ComputeSum".
					// Perfect.
					
					sig := obj.Type().(*types.Signature)
					if a.HasDependentSignature(sig) {
						found = true
						return false
					}
				}
				// Also check if it's a SIMD type directly (e.g. simd.Int32s)
				// (markIfDependent handles checking dependentObj and isDependentType, but obj might be from another pkg)
				if a.isDependentType(obj.Type()) {
					found = true
					return false
				}
				if isBaseSimdTypeObj(obj) {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

func (a *Analyzer) markIfDependent(obj types.Object) bool {
	if a.dependentObj[obj] {
		return true
	}

	isDep := false
	switch obj := obj.(type) {
	case *types.Var:
		isDep = a.isDependentType(obj.Type())
	case *types.TypeName:
		isDep = a.isDependentType(obj.Type())
	case *types.Func:
		// Check signature
		sig := obj.Type().(*types.Signature)
		if a.isDependentType(sig.Params()) || a.isDependentType(sig.Results()) || (sig.Recv() != nil && a.isDependentType(sig.Recv().Type())) {
			isDep = true
		}
	}

	// Also check if obj name is "simd.Type" (base case)
	if isBaseSimdTypeObj(obj) {
		isDep = true
	}

	if isDep {
		a.dependentObj[obj] = true
	}
	return isDep
}

func (a *Analyzer) isDependentType(t types.Type) bool {
	return a.checkTypeRecursive(t, make(map[types.Type]bool))
}

func (a *Analyzer) checkTypeRecursive(t types.Type, visited map[types.Type]bool) bool {
	if visited[t] {
		return false // Break cycles
	}
	visited[t] = true
    if t == nil { return false }
    // fmt.Printf("DEBUG: checkTypeRecursive %s (%T)\n", t.String(), t)

	// Unwrap aliases
	if named, ok := t.(*types.Named); ok {
		if isBaseSimdType(named) {
			return true
		}
		if a.checkTypeRecursive(named.Underlying(), visited) {
			return true
		}
	}

	switch t := t.(type) {
	case *types.Basic:
		return false
	case *types.Pointer:
		return a.checkTypeRecursive(t.Elem(), visited)
	case *types.Slice:
		return a.checkTypeRecursive(t.Elem(), visited)
	case *types.Array:
		return a.checkTypeRecursive(t.Elem(), visited)
	case *types.Map:
		return a.checkTypeRecursive(t.Key(), visited) || a.checkTypeRecursive(t.Elem(), visited)
	case *types.Chan:
		return a.checkTypeRecursive(t.Elem(), visited)
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			if a.checkTypeRecursive(t.Field(i).Type(), visited) {
				return true
			}
		}
	case *types.Signature:
		return a.checkTypeRecursive(t.Params(), visited) || a.checkTypeRecursive(t.Results(), visited) || (t.Recv() != nil && a.checkTypeRecursive(t.Recv().Type(), visited))
	case *types.Tuple:
		for i := 0; i < t.Len(); i++ {
			if a.checkTypeRecursive(t.At(i).Type(), visited) {
				return true
			}
		}
    case *types.Alias:
        return a.checkTypeRecursive(t.Rhs(), visited)
	}
	return false
}

func isBaseSimdType(t *types.Named) bool {
	obj := t.Obj()
	return isBaseSimdTypeObj(obj)
}

func isBaseSimdTypeObj(obj types.Object) bool {
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	if obj.Pkg().Name() != "simd" {
		return false
	}
	switch obj.Name() {
	case "Int8s", "Int16s", "Int32s", "Int64s",
		"Uint8s", "Uint16s", "Uint32s", "Uint64s",
		"Float32s", "Float64s":
		return true
	}
	return false
}

// HasDependentSignature checks if the signature involves SIMD types directly (params/results/receiver)
func (a *Analyzer) HasDependentSignature(sig *types.Signature) bool {
	return a.isDependentType(sig.Params()) || a.isDependentType(sig.Results()) || (sig.Recv() != nil && a.isDependentType(sig.Recv().Type()))
}
