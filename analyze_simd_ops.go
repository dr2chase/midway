// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
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

func main() {
	// Hardcoded path to archsimd
	archSimdPath := "/Users/drchase/work/go/src/simd/archsimd"

	files := []string{"ops_amd64.go", "types_amd64.go", "other_gen_amd64.go", "extra_amd64.go", "maskmerge_gen_amd64.go", "shuffles_amd64.go", "slice_gen_amd64.go", "slicepart_amd64.go", "string.go"}

	// Categories based on bit size
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

				// Exclude "grouped" methods
				methodName := funcDecl.Name.Name

				// Allow reinterpret vectors.
				if xAt := strings.Index(methodName, "x"); xAt != -1 && (strings.HasPrefix(methodName, "As") || strings.HasPrefix(methodName, "ToInt") && strings.HasPrefix(eltType, "Mask")) {
					// We think this is fine, even if it changes the number of elements in the vector.
					// Tweak the method name so that they will line up properly.
					methodName = methodName[:xAt] + "s"
				} else if strings.HasPrefix(methodName, "Broadcast") {
					// Broadcast is okay
				} else {
					if strings.Contains(methodName, "Group") {
						fmt.Fprintf(os.Stderr, "Skipping grouped method %s.%s\n", recvType, methodName)
						continue
					}
					if lastChar := methodName[len(methodName)-1]; unicode.IsDigit(rune(lastChar)) && lastChar != eltType[len(eltType)-1] {
						fmt.Fprintf(os.Stderr, "Skipping size-changing method %s.%s\n", recvType, methodName)
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

	fmt.Println("Intersection of methods for 256-bit and 512-bit vectors:")

	sigForMethod := make(map[string]*ast.FuncDecl)

	var xlateType func(*ast.Expr) string
	xlateType = func(e *ast.Expr) string {

	}

	for _, elem := range elems {
		type256 := map256[elem]
		type512 := map512[elem]

		methods256 := methodsByType[type256]
		methods512 := methodsByType[type512]

		var intersection []string
		for m := range methods256 {
			if fd := methods512[m]; fd != nil {
				intersection = append(intersection, m)
				sigForMethod[m] = fd
			}
		}
		sort.Strings(intersection)

		fmt.Printf("\n// Element Type: %s\n", elem)
		fmt.Printf("//   256-bit Type: %s (Methods: %d)\n", type256, len(methods256))
		fmt.Printf("//   512-bit Type: %s (Methods: %d)\n", type512, len(methods512))

		fmt.Printf("//   Intersection (%d): %v\n", len(intersection), intersection)
		for _, m := range intersection {
			fd := sigForMethod[m]
			fmt.Printf("func (%s)%s(", elem+"s")
			fmt.Printf(") %s", xlateType(fd.Type.TypeParams.List[0].Type))

		}
	}
}
