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
	// We assume we can just overwrite them.

	for _, fileAST := range r.pkg.Syntax {
		tokenFile := r.pkg.Fset.File(fileAST.Pos())
		if tokenFile == nil {
			continue
		}
		filename := tokenFile.Name()
		if strings.Contains(filename, "_simd") {
			continue
		}

		modified := false

		for _, decl := range fileAST.Decls {
			fnDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Check if Dependent
			obj := r.pkg.TypesInfo.ObjectOf(fnDecl.Name)
			if !r.analyzer.dependentObj[obj] {
				continue
			}

			// Check if Signature is Clean (no Dependent types in Params/Results)
			// We use the already computed dependency info for the signature components.
			// Wait, `analyzer.dependentObj` tracks objects.
			// We need to check if the `Type` of the function involves SIMD types.
			// `obj.Type()` is a `*types.Signature`.
			sig := obj.Type().(*types.Signature)

			// Re-use `analyzer.isDependentType` logic?
			// `Analyzer` needs to expose `isDependentType` or we check manually.
			// `r.analyzer.isDependentType` is private. I should make it public or duplicate check.
			// Let's assume we can access it or I'll export it.
			// For now, I'll export `IsDependentType` in `analysis.go` in next step if needed, or stick to this plan.
			// Better: `markIfDependent` already did `dependentObj[obj] = true` based on signature.
			// Accessing `IsDependentType` is cleaner.

			if r.analyzer.HasDependentSignature(sig) {
				// HAS dependent signature -> DO NOT Dispatch (as per instructions)
				continue
			}

			// If we are here, it is Dependent (calls internal SIMD stuff) but has Clean Signature.
			// Rewrite Body!
			fnDecl.Body = r.createDispatcherBody(fnDecl.Name.Name, fnDecl.Type)
			modified = true
		}

		if modified {
			// Overwrite file
			// Imports? We might need to ensure `simd` is imported if we reference `simd.VectorSize()`.
			// Check if `simd` is imported. If not, add it.
			// `imports.Process` will handle it!

			var buf strings.Builder
			if err := format.Node(&buf, r.pkg.Fset, fileAST); err != nil {
				return fmt.Errorf("formatting failed: %v", err)
			}

			res, err := imports.Process(filename, []byte(buf.String()), nil)
			if err != nil {
				return fmt.Errorf("imports processing failed for %s: %v", filename, err)
			}

			if err := os.WriteFile(filename, res, 0644); err != nil {
				return err
			}
			fmt.Printf("Updated dispatcher: %s\n", filename)
		}
	}
	return nil
}

func (r *Rewriter) createDispatcherBody(funcName string, funcType *ast.FuncType) *ast.BlockStmt {
	// switch simd.VectorSize() { ... }

	// Build call arguments
	var args []ast.Expr
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			for _, name := range field.Names {
				args = append(args, name)
			}
		}
	}

	// Create Switch Stmt
	switchStmt := &ast.SwitchStmt{
		Tag: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("simd"),
				Sel: ast.NewIdent("VectorSize"),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}

	for _, k := range r.sizes {
		// case K: return funcName_simdK(args...)
		callExpr := &ast.CallExpr{
			Fun:  ast.NewIdent(fmt.Sprintf("%s_simd%d", funcName, k)),
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
		// Use TypesInfo from original package to check original object
		// But wait, `id` here is the node we are visiting in the ORIGINAL AST (before copy)?
		// NO, DeepCopier visits the ORIGINAL AST nodes in the input arguments.
		// `DeepCopier.CopyIdent(id *ast.Ident)` -> `id` is the original node.
		// So we can look it up!

		obj := r.pkg.TypesInfo.ObjectOf(id)
		if obj == nil {
			return nil
		}

		// If the object is dependent (AND is defined in this package or is a SIMD type), we rename it.
		// Special case: SIMD types from `simd` package -> `Int8s` -> `Int8s_simd128`
		// Local variables marked dependent -> `x` -> `x_simd128`
		// Top level funcs marked dependent -> `F` -> `F_simd128`

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

	copier := &DeepCopier{OnIdent: onIdent}

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

		newFileAST := &ast.File{
			Name:  ast.NewIdent(r.pkg.Name),
			Decls: newDecls,
		}

		// Add imports
		for _, decl := range fileAST.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				newFileAST.Decls = append([]ast.Decl{copier.CopyDecl(genDecl)}, newFileAST.Decls...)
			}
		}

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
