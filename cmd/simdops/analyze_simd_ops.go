// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type MethodSet map[string]*ast.FuncDecl
type TypeMethods map[string]MethodSet

var goRoot = flag.String("goroot", "/Users/drchase/work/go", "Go root")
var verbose = flag.Bool("v", false, "Be much chattier about processing")

func main() {
	flag.Parse()

	p := func(s ...any) { fmt.Print(s...) }
	pf := func(f string, s ...any) { fmt.Printf(f, s...) }
	pv := func(f string, s ...any) {
		if *verbose {
			fmt.Fprintf(os.Stderr, f, s...)
		}
	}
	nl := func() { fmt.Println() }

	// Hardcoded path to archsimd
	archSimdPath := *goRoot + "/src/simd/archsimd"

	// Hardcoded list of files
	files := []string{"ops_amd64.go", "types_amd64.go", "other_gen_amd64.go", "extra_amd64.go", "maskmerge_gen_amd64.go", "shuffles_amd64.go", "slice_gen_amd64.go", "slicepart_amd64.go", "string.go"}

	// Categories based on bit size
	// 128-bit map: ElementType -> TypeName
	map128 := map[string]string{
		"Int8":    "Int8x16",
		"Int16":   "Int16x8",
		"Int32":   "Int32x4",
		"Int64":   "Int64x2",
		"Uint8":   "Uint8x16",
		"Uint16":  "Uint16x8",
		"Uint32":  "Uint32x4",
		"Uint64":  "Uint64x2",
		"Float32": "Float32x4",
		"Float64": "Float64x2",
		"Mask8":   "Mask8x16",
		"Mask16":  "Mask16x8",
		"Mask32":  "Mask32x4",
		"Mask64":  "Mask64x2",
	}

	// 256-bit map: ElementType -> TypeName
	map256 := map[string]string{
		"Int8":    "Int8x32",
		"Int16":   "Int16x16",
		"Int32":   "Int32x8",
		"Int64":   "Int64x4",
		"Uint8":   "Uint8x32",
		"Uint16":  "Uint16x16",
		"Uint32":  "Uint32x8",
		"Uint64":  "Uint64x4",
		"Float32": "Float32x8",
		"Float64": "Float64x4",
		"Mask8":   "Mask8x32",
		"Mask16":  "Mask16x16",
		"Mask32":  "Mask32x8",
		"Mask64":  "Mask64x4",
	}

	map512 := map[string]string{
		"Int8":    "Int8x64",
		"Int16":   "Int16x32",
		"Int32":   "Int32x16",
		"Int64":   "Int64x8",
		"Uint8":   "Uint8x64",
		"Uint16":  "Uint16x32",
		"Uint32":  "Uint32x16",
		"Uint64":  "Uint64x8",
		"Float32": "Float32x16",
		"Float64": "Float64x8",
		"Mask8":   "Mask8x64",
		"Mask16":  "Mask16x32",
		"Mask32":  "Mask32x16",
		"Mask64":  "Mask64x8",
	}

	methodsByType := make(TypeMethods)

	fset := token.NewFileSet()

	knownReceivers := make(map[string]string)
	for k, v := range map128 {
		knownReceivers[v] = k + "s"
	}
	for k, v := range map256 {
		knownReceivers[v] = k + "s"
	}
	for k, v := range map512 {
		knownReceivers[v] = k + "s"
	}

	for _, fname := range files {
		path := filepath.Join(archSimdPath, fname)
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse %s: %v", path, err)
		}

		for _, decl := range f.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Recv == nil {
					continue
				}

				// Identify receiver type
				var recvType string
				for _, field := range funcDecl.Recv.List {
					// We assume single receiver
					if ident, ok := field.Type.(*ast.Ident); ok {
						recvType = ident.Name
					} else if star, ok := field.Type.(*ast.StarExpr); ok {
						if ident, ok := star.X.(*ast.Ident); ok {
							recvType = ident.Name
						}
					}
				}

				if recvType == "" {
					continue
				}

				if knownReceivers[recvType] == "" {
					continue
				}

				eltType := recvType[:strings.Index(recvType, "x")]

				methodName := funcDecl.Name.Name

				// Allow reinterpret vectors.
				if xAt := strings.Index(methodName, "x"); xAt != -1 && (strings.HasPrefix(methodName, "As") || strings.HasPrefix(methodName, "ToInt") && strings.HasPrefix(eltType, "Mask")) {
					// We think this is fine, even if it changes the number of elements in the vector.
					// Tweak the method name so that they will line up properly.
					methodName = methodName[:xAt] + "s"
				} else if strings.HasPrefix(methodName, "Broadcast") {
					// Broadcast is okay
				} else {
					// Exclude "grouped", "Store" (not slice), and vector-size-changing methods.
					if strings.Contains(methodName, "Group") {
						pv("Skipping grouped method %s.%s\n", recvType, methodName)
						continue
					}
					if methodName == "Store" || methodName == "StoreMasked" {
						pv("Skipping fixed-size Store method method %s.%s\n", recvType, methodName)
						continue
					}
					if methodName == "ToBits" {
						pv("Skipping ToBits method (has varying return type) %s.%s\n", recvType, methodName)
						continue
					}
					if lastChar := methodName[len(methodName)-1]; unicode.IsDigit(rune(lastChar)) && lastChar != eltType[len(eltType)-1] {
						pv("Skipping size-changing method %s.%s\n", recvType, methodName)
						continue
					}
				}

				if methodsByType[recvType] == nil {
					methodsByType[recvType] = make(MethodSet)
				}
				methodsByType[recvType][methodName] = funcDecl
			}
		}
	}

	elems := []string{"Int8", "Int16", "Int32", "Int64", "Uint8", "Uint16", "Uint32", "Uint64", "Float32", "Float64", "Mask8", "Mask16", "Mask32", "Mask64"}

	fmt.Println("// Intersection of methods for 128-bit, 256-bit, and 512-bit vectors:")
	fmt.Println(
		`// Code generated by 'go run cmd/simdops/analyze_simd_ops.go > simd/mocks.go'; DO NOT EDIT.

package simd

func VectorSize() int`)

	sigForMethod := make(map[string]*ast.FuncDecl)

	// xlateType translates a type by replacing instances of types with keys in knownReceivers with their values,
	// and generates the string representation of the resulting type.  E.g., []Int8x32 -> []Int8s
	// (because Int8x32 -> Int8s in knownReceivers
	var xlateType func(ast.Expr) string
	xlateType = func(e ast.Expr) string {
		switch t := e.(type) {
		case *ast.Ident:
			if mapped, ok := knownReceivers[t.Name]; ok {
				return mapped
			}
			return t.Name
		case *ast.StarExpr:
			return "*" + xlateType(t.X)
		case *ast.ArrayType:
			lenStr := ""
			if t.Len != nil {
				var buf strings.Builder
				format.Node(&buf, token.NewFileSet(), t.Len)
				lenStr = buf.String()
			}
			return "[" + lenStr + "]" + xlateType(t.Elt)
		case *ast.SelectorExpr:
			return xlateType(t.X) + "." + t.Sel.Name
		case *ast.Ellipsis:
			return "..." + xlateType(t.Elt)
		default:
			var buf strings.Builder
			format.Node(&buf, token.NewFileSet(), t)
			return buf.String()
		}
	}

	toScalar := func(s string) string {
		if strings.HasPrefix(s, "Mask") {
			return "int" + s[4:]
		}
		return strings.ToLower(s)
	}

	for _, elem := range elems {
		type128 := map128[elem]
		type256 := map256[elem]
		type512 := map512[elem]

		methods128 := methodsByType[type128]
		methods256 := methodsByType[type256]
		methods512 := methodsByType[type512]

		var intersection []string
		for m := range methods128 {
			if methods256[m] != nil && methods512[m] != nil {
				intersection = append(intersection, m)
				sigForMethod[m] = methods512[m] // Use 512-bit signature (arbitrary choice, they should match)
			}
		}
		sort.Strings(intersection)

		pf("\n// Element Type: %s\n//\n", elem)
		pf("//\t128-bit Type: %s (Methods: %d)\n", type128, len(methods128))
		pf("//\t256-bit Type: %s (Methods: %d)\n", type256, len(methods256))
		pf("//\t512-bit Type: %s (Methods: %d)\n", type512, len(methods512))

		pf("type %ss []%s\n\n", elem, toScalar(elem))
		pf("func Load%sSlice([]%s) %ss\n", elem, toScalar(elem), elem)
		pf("func Load%sSlicePart([]%s) %ss\n", elem, toScalar(elem), elem)

		// pf("//   Intersection (%d): %v\n", len(intersection), intersection)
		for _, m := range intersection {
			fd := sigForMethod[m]
			pf("func (x %s) %s(", elem+"s", m)

			if fd.Type.Params != nil {
				for i, field := range fd.Type.Params.List {
					if i > 0 {
						p(", ")
					}
					if len(field.Names) > 0 {
						for j, name := range field.Names {
							if j > 0 {
								p(", ")
							}
							p(name.Name)
						}
						p(" ")
					}
					p(xlateType(field.Type))
				}
			}
			p(")")

			if fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
				p(" ")
				needsParens := len(fd.Type.Results.List) > 1 || (len(fd.Type.Results.List) == 1 && len(fd.Type.Results.List[0].Names) > 0)
				if needsParens {
					p("(")
				}
				for i, field := range fd.Type.Results.List {
					if i > 0 {
						p(", ")
					}
					if len(field.Names) > 0 {
						for j, name := range field.Names {
							if j > 0 {
								p(", ")
							}
							p(name.Name)
						}
						p(" ")
					}
					p(xlateType(field.Type))
				}
				if needsParens {
					p(")")
				}
			}
			nl()
		}
	}
}
