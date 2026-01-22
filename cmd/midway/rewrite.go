// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

// Rewriter handles the generation of specialized code
type Rewriter struct {
	pkg      *packages.Package
	analyzer *Analyzer
	sizes    []int
}

func NewRewriter(pkg *packages.Package, analyzer *Analyzer, sizes []int) *Rewriter {
	return &Rewriter{
		pkg:      pkg,
		analyzer: analyzer,
		sizes:    sizes,
	}
}

// Rewrite generates the specialized files
func (r *Rewriter) Rewrite() error {
	for _, k := range r.sizes {
		if err := r.generateForSize(k); err != nil {
			return err
		}
	}

	// Generate Dispatchers logic is omitted for now (focus on generation first) or added here.
	if err := r.generateDispatchers(); err != nil {
		return err
	}
	return nil
}

func (r *Rewriter) generateDispatchers() error {
	// Iterate over original files to modify them in place
	for _, fileAST := range r.pkg.Syntax {
		tokenFile := r.pkg.Fset.File(fileAST.Pos())
		if tokenFile == nil {
			continue
		}
		filename := tokenFile.Name()
		if strings.Contains(filename, "_simd") {
			continue
		}

		// We will build a new list of Decls
		var newDecls []ast.Decl
		modified := false

		for _, decl := range fileAST.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				obj := r.pkg.TypesInfo.ObjectOf(d.Name)
				if !r.analyzer.dependentObj[obj] {
					// Not dependent, keep as is
					newDecls = append(newDecls, d)
					continue
				}

				// It IS dependent.
				sig := obj.Type().(*types.Signature)
				if r.analyzer.HasDependentSignature(sig) {
					// Dependent Signature -> Remove (Drop)
					modified = true
					continue
				}

				// Clean Signature -> Dispatcher
				d.Body = r.createDispatcherBody(d.Name.Name, d.Type)
				newDecls = append(newDecls, d)
				modified = true

			case *ast.GenDecl:
				// Filter Specs
				var newSpecs []ast.Spec
				for _, spec := range d.Specs {
					keep := true
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if r.analyzer.dependentObj[r.pkg.TypesInfo.ObjectOf(s.Name)] {
							keep = false
						}
					case *ast.ValueSpec:
						// If ANY name is dependent, is the whole spec dependent?
						// "Variables are rewritten... if their type is one of the simd types"
						// If we have `var x, y simd.Int8s`, both are dependent.
						// If `var x int, y simd.Int8s` (can't happen in one spec unless tuple assign? No, Go syntax restricts type).
						for _, name := range s.Names {
							if r.analyzer.dependentObj[r.pkg.TypesInfo.ObjectOf(name)] {
								keep = false
								break
							}
						}
					}
					if keep {
						newSpecs = append(newSpecs, spec)
					} else {
						modified = true
					}
				}

				if len(newSpecs) > 0 {
					d.Specs = newSpecs
					newDecls = append(newDecls, d)
				} else if len(d.Specs) > 0 {
					// If we removed all specs, we effectively modified the file (by removing the decl)
					// Even if modified=true was set in loop, we confirm it here.
				}

			default:
				newDecls = append(newDecls, decl)
			}
		}

		if modified {
			// Update Decls
			fileAST.Decls = newDecls

			// Filter out existing build tags to prevent duplicates/conflicts
			var newComments []*ast.CommentGroup
			var newBuild = "//go:build !midway"
			for _, cg := range fileAST.Comments {
				keep := true
				for _, c := range cg.List {
					text := strings.TrimSpace(c.Text)
					if strings.HasPrefix(text, "//go:build") || (strings.HasPrefix(text, "// +build") && strings.Contains(text, "midway")) {
						keep = false
						break
					}
					pfx := "//+go:build"
					if strings.HasPrefix(text, pfx) {
						suffix := text[len(pfx):]
						newBuild = newBuild + " &&" + suffix
						keep = false
					}
				}
				if keep {
					newComments = append(newComments, cg)
				}
			}
			fileAST.Comments = newComments

			midwayImport := &ast.ImportSpec{
				Path: &ast.BasicLit{Kind: token.STRING, Value: "\"" + *midwayPackage + "\""},
			}

			// Replace imports (there must be at least one, if the file was modified)
			for _, decl := range fileAST.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
					for _, spec := range genDecl.Specs {
						if imp, ok := spec.(*ast.ImportSpec); ok {
							// Check if path ends with "/simd" or is "simd"
							pathVal := strings.Trim(imp.Path.Value, "\"")
							if pathVal == "simd" || strings.HasSuffix(pathVal, "/simd") {
								// Replace with archsimd
								imp.Name = ast.NewIdent("archsimd")
								imp.Path.Value = fmt.Sprintf("\"%s/archsimd\"", *archsimdPfxFlag)
							}
						}
					}
					if midwayImport != nil {
						genDecl.Specs = append(genDecl.Specs, midwayImport)
						midwayImport = nil
					}
				}
			}

			var buf strings.Builder
			// Prepend build tag
			buf.WriteString(newBuild + "\n\n")

			if err := format.Node(&buf, r.pkg.Fset, fileAST); err != nil {
				return fmt.Errorf("formatting failed: %v", err)
			}

			baseName := strings.TrimSuffix(filepath.Base(filename), ".go")
			outName := filepath.Join(filepath.Dir(filename), baseName+"_simd.go")

			res, err := imports.Process(outName, []byte(buf.String()), nil)
			if err != nil {
				return fmt.Errorf("imports processing failed for %s: %v", outName, err)
			}

			if err := os.WriteFile(outName, res, 0644); err != nil {
				return err
			}
			fmt.Printf("Generated dispatcher (filtered): %s\n", outName)
		}
	}
	return nil
}

func (r *Rewriter) createDispatcherBody(funcName string, funcType *ast.FuncType) *ast.BlockStmt {
	// switch archsimd.MaxVectorSize() { ... }

	// Build call arguments
	var args []ast.Expr
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			for _, name := range field.Names {
				args = append(args, name)
			}
		}
	}

	// Build type arguments if any
	var typeArgs []ast.Expr
	if funcType.TypeParams != nil {
		for _, field := range funcType.TypeParams.List {
			for _, name := range field.Names {
				typeArgs = append(typeArgs, name)
			}
		}
	}

	// Create Switch Stmt
	switchStmt := &ast.SwitchStmt{
		Tag: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("midway"),
				Sel: ast.NewIdent("MaxVectorSize"),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}

	for _, k := range r.sizes {
		// case K: return funcName_simdK[T](args...)
		fnIdent := ast.NewIdent(fmt.Sprintf("%s_simd%d", funcName, k))
		var fun ast.Expr = fnIdent

		if len(typeArgs) > 0 {
			if len(typeArgs) == 1 {
				fun = &ast.IndexExpr{
					X:     fnIdent,
					Index: typeArgs[0],
				}
			} else {
				fun = &ast.IndexListExpr{
					X:       fnIdent,
					Indices: typeArgs,
				}
			}
		}

		callExpr := &ast.CallExpr{
			Fun:  fun,
			Args: args,
		}

		var branchStmt ast.Stmt
		if funcType.Results != nil && len(funcType.Results.List) > 0 {
			branchStmt = &ast.ReturnStmt{Results: []ast.Expr{callExpr}}
		} else {
			branchStmt = &ast.ExprStmt{X: callExpr}
			// For void functions, we should explicitly return? Or break?
			// Since it's a switch, `case` falls through only if `fallthrough`.
			// But we need to return to exit the function.
			branchStmt = &ast.BlockStmt{
				List: []ast.Stmt{
					branchStmt,
					&ast.ReturnStmt{},
				},
			}
		}

		caseClause := &ast.CaseClause{
			List: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", k)}},
			Body: []ast.Stmt{branchStmt},
		}
		switchStmt.Body.List = append(switchStmt.Body.List, caseClause)
	}

	// Add default panic?
	switchStmt.Body.List = append(switchStmt.Body.List, &ast.CaseClause{
		List: nil, // default
		Body: []ast.Stmt{
			&ast.ExprStmt{X: &ast.CallExpr{
				Fun:  ast.NewIdent("panic"), // simplistic
				Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"unsupported vector size\""}},
			}},
		},
	})

	return &ast.BlockStmt{List: []ast.Stmt{switchStmt}}
}

func (r *Rewriter) generateForSize(k int) error {
	suffix := fmt.Sprintf("_simd%d", k)

	// We handle identifiers by checking if they resolve to Dependent objects
	onIdent := func(id *ast.Ident) *ast.Ident {
		obj := r.pkg.TypesInfo.ObjectOf(id)
		if obj == nil {
			return nil
		}

		shouldRename := false
		if r.analyzer.dependentObj[obj] {
			shouldRename = true
		} else if isBaseSimdTypeObj(obj) {
			// It might not be in dependentObj set if we didn't track it explicitly there,
			// but we want to rewrite `simd.Int8s` to `simd.Int8s_simd128`.
			shouldRename = true
		}

		if shouldRename {
			return &ast.Ident{
				NamePos: id.NamePos,
				Name:    id.Name + suffix, // The rewriting happens here
				Obj:     nil,              // New object is not resolved yet
			}
		}

		return nil // Use default copy behavior
	}

	onSelector := func(se *ast.SelectorExpr) ast.Expr {
		if x, ok := se.X.(*ast.Ident); ok {
			if obj, ok := r.pkg.TypesInfo.ObjectOf(x).(*types.PkgName); ok {
				if obj.Imported().Name() == "simd" {
					isLoad := false
					suffix := "Slice"
					name := se.Sel.Name
					// Looking for simd.Load<Type><Size>Slice[Part].
					// If so, extract <Type><Size> and append "s" to get the type translation.
					if p := strings.Index(name, "Slice"); p > 0 && strings.HasPrefix(name, "Load") {
						isLoad = true
						if strings.HasSuffix(name, "SlicePart") {
							suffix = "SlicePart"
						}
						name = name[len("Load"):p] + "s"
					}
					width := nameToWidth(name)
					if width > 0 {
						count := k / width
						base := name[:len(name)-1]
						newName := fmt.Sprintf("%sx%d", base, count)
						if isLoad {
							newName = "Load" + newName + suffix
						}
						return &ast.SelectorExpr{
							X:   ast.NewIdent("archsimd"),
							Sel: ast.NewIdent(newName),
						}
					}
				}
			}
		}
		return nil
	}

	copier := &DeepCopier{OnIdent: onIdent, OnSelector: onSelector, VecLen: k}

	for _, fileAST := range r.pkg.Syntax {
		tokenFile := r.pkg.Fset.File(fileAST.Pos())
		if tokenFile == nil {
			continue
		}
		filename := tokenFile.Name()
		if strings.Contains(filename, "_simd") {
			continue
		}

		var newDecls []ast.Decl

		for _, decl := range fileAST.Decls {
			if r.shouldIncludeDecl(decl) {
				newDecl := copier.CopyDecl(decl)
				newDecls = append(newDecls, newDecl)
			}
		}

		if len(newDecls) == 0 {
			continue
		}

		var importDecls []ast.Decl

		// Add imports
		archSimdName := fmt.Sprintf("\"%s/archsimd\"", *archsimdPfxFlag)

		// Inject archsimd import and midway import (for assert<Size>)
		archSimdImport := &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Name: ast.NewIdent("archsimd"),
					Path: &ast.BasicLit{Kind: token.STRING, Value: archSimdName},
				},
				&ast.ImportSpec{
					Path: &ast.BasicLit{Kind: token.STRING, Value: "\"" + *midwayPackage + "\""},
				},
			},
		}
		importDecls = append(importDecls, archSimdImport)

	declLoop:
		for _, decl := range fileAST.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				for _, spec := range genDecl.Specs {
					// don't copy archsimd if it was imported in the input file
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						if importSpec.Path.Value == archSimdName {
							continue declLoop
						}
					}
				}
				importDecls = append(importDecls, copier.CopyDecl(genDecl))
			}
		}

		newFileAST := &ast.File{
			Name:    ast.NewIdent(r.pkg.Name),
			Package: fileAST.Package, // Preserve Package Pos to avoid comment interleaving
			Decls:   append(importDecls, newDecls...),
		}

		// Replace "midway" with "!midway"
		var newComments []*ast.CommentGroup

		for _, cg := range fileAST.Comments {
			newcg := &ast.CommentGroup{}
			for _, c := range cg.List {
				text := strings.TrimSpace(c.Text)
				if i := strings.Index(text, "midway"); i > 0 && strings.HasPrefix(text, "//go:build") && text[i-1] != '!' {
					newC := *c
					newC.Text = strings.ReplaceAll(c.Text, "midway", "!midway")
					newcg.List = append(newcg.List, &newC)
				} else {
					newcg.List = append(newcg.List, c)
				}
			}
			newComments = append(newComments, newcg)
		}
		newFileAST.Comments = newComments

		baseName := strings.TrimSuffix(filepath.Base(filename), ".go")
		outName := filepath.Join(filepath.Dir(filename), baseName+suffix+".go")

		var buf strings.Builder

		if err := format.Node(&buf, r.pkg.Fset, newFileAST); err != nil {
			return fmt.Errorf("formatting failed: %v", err)
		}

		res, err := imports.Process(outName, []byte(buf.String()), nil)
		if err != nil {
			if writeErr := os.WriteFile(outName, []byte(buf.String()), 0644); writeErr != nil {
				return writeErr
			}
			return fmt.Errorf("imports processing failed for %s: %v", outName, err)
		}

		if err := os.WriteFile(outName, res, 0644); err != nil {
			return err
		}
		fmt.Printf("Generated %s\n", outName)
	}
	return nil
}

func nameToWidth(name string) int {
	var width int
	switch name {
	case "Int8s", "Uint8s", "Mask8s":
		width = 8
	case "Int16s", "Uint16s", "Mask16s":
		width = 16
	case "Int32s", "Uint32s", "Float32s", "Mask32s":
		width = 32
	case "Int64s", "Uint64s", "Float64s", "Mask64s":
		width = 64
	}
	return width
}

func (r *Rewriter) shouldIncludeDecl(decl ast.Decl) bool {
	// Check if decl contains dependent objects.
	switch d := decl.(type) {
	case *ast.FuncDecl:
		obj := r.pkg.TypesInfo.ObjectOf(d.Name)
		return r.analyzer.dependentObj[obj]
	case *ast.GenDecl:
		if d.Tok == token.TYPE || d.Tok == token.VAR {
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if r.analyzer.dependentObj[r.pkg.TypesInfo.ObjectOf(s.Name)] {
						return true
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if r.analyzer.dependentObj[r.pkg.TypesInfo.ObjectOf(name)] {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// Helper duplication from analysis.go since we are avoiding circular deps or just keeping it simple
// isBaseSimdTypeObj maintained in analysis.go

func parserParseFileWrapper(fset *token.FileSet, filename string, src interface{}) (*ast.File, error) {
	return parser.ParseFile(fset, filename, src, parser.ParseComments)
}
