// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
)

// DeepCopier clones AST nodes.
type DeepCopier struct {
	VecLen int

	// OnIdent, if provided, handles identifier cloning.
	// If it returns nil, a default clone is performed.
	OnIdent func(*ast.Ident) *ast.Ident

	// OnSelector, if provided, handles selector expression cloning.
	// If it returns nil, a default clone is performed.
	OnSelector func(*ast.SelectorExpr) ast.Expr
}

func (c *DeepCopier) CopyDecl(d ast.Decl) ast.Decl {
	switch d := d.(type) {
	case *ast.GenDecl:
		return c.CopyGenDecl(d)
	case *ast.FuncDecl:
		return c.CopyFuncDecl(d)
	default:
		return d // Other declarations not handled/needed for this task
	}
}

func (c *DeepCopier) CopyGenDecl(d *ast.GenDecl) *ast.GenDecl {
	newD := &ast.GenDecl{
		Doc:    c.CopyCommentGroup(d.Doc),
		TokPos: d.TokPos,
		Tok:    d.Tok,
		Lparen: d.Lparen,
		Rparen: d.Rparen,
	}
	for _, s := range d.Specs {
		newD.Specs = append(newD.Specs, c.CopySpec(s))
	}
	return newD
}

func (c *DeepCopier) CopySpec(s ast.Spec) ast.Spec {
	switch s := s.(type) {
	case *ast.ValueSpec:
		newS := &ast.ValueSpec{
			Doc:     c.CopyCommentGroup(s.Doc),
			Comment: c.CopyCommentGroup(s.Comment),
			Type:    c.CopyExpr(s.Type),
		}
		for _, n := range s.Names {
			newS.Names = append(newS.Names, c.CopyIdent(n))
		}
		for _, v := range s.Values {
			newS.Values = append(newS.Values, c.CopyExpr(v))
		}
		return newS
	case *ast.TypeSpec:
		newS := &ast.TypeSpec{
			Doc:        c.CopyCommentGroup(s.Doc),
			Comment:    c.CopyCommentGroup(s.Comment),
			Name:       c.CopyIdent(s.Name),
			TypeParams: c.CopyFieldList(s.TypeParams),
			Assign:     s.Assign,
			Type:       c.CopyExpr(s.Type),
		}
		return newS
	default:
		return s // ImportSpec etc
	}
}

func (c *DeepCopier) CopyFuncDecl(d *ast.FuncDecl) *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc:  c.CopyCommentGroup(d.Doc),
		Recv: c.CopyFieldList(d.Recv),
		Name: c.CopyIdent(d.Name),
		Type: c.CopyFuncType(d.Type),
		Body: c.CopyBlockStmt(d.Body),
	}
}

func (c *DeepCopier) CopyExpr(e ast.Expr) ast.Expr {
	if e == nil {
		return nil
	}
	switch e := e.(type) {
	case *ast.Ident:
		return c.CopyIdent(e)
	case *ast.StarExpr:
		return &ast.StarExpr{Star: e.Star, X: c.CopyExpr(e.X)}
	case *ast.ArrayType:
		return &ast.ArrayType{Lbrack: e.Lbrack, Len: c.CopyExpr(e.Len), Elt: c.CopyExpr(e.Elt)}
	case *ast.SelectorExpr:
		if c.OnSelector != nil {
			if sub := c.OnSelector(e); sub != nil {
				return sub
			}
		}
		return &ast.SelectorExpr{X: c.CopyExpr(e.X), Sel: c.CopyIdent(e.Sel)}
	case *ast.CallExpr:
		newE := &ast.CallExpr{
			Fun:      c.CopyExpr(e.Fun),
			Lparen:   e.Lparen,
			Ellipsis: e.Ellipsis,
			Rparen:   e.Rparen,
		}
		for _, a := range e.Args {
			newE.Args = append(newE.Args, c.CopyExpr(a))
		}
		return newE
	case *ast.ParenExpr:
		return &ast.ParenExpr{Lparen: e.Lparen, X: c.CopyExpr(e.X), Rparen: e.Rparen}
	case *ast.TypeAssertExpr:
		return &ast.TypeAssertExpr{X: c.CopyExpr(e.X), Lparen: e.Lparen, Type: c.CopyExpr(e.Type), Rparen: e.Rparen}
	case *ast.IndexExpr:
		return &ast.IndexExpr{X: c.CopyExpr(e.X), Lbrack: e.Lbrack, Index: c.CopyExpr(e.Index), Rbrack: e.Rbrack}
	case *ast.SliceExpr:
		return &ast.SliceExpr{X: c.CopyExpr(e.X), Lbrack: e.Lbrack, Low: c.CopyExpr(e.Low), High: c.CopyExpr(e.High), Max: c.CopyExpr(e.Max), Slice3: e.Slice3, Rbrack: e.Rbrack}
	case *ast.BinaryExpr:
		return &ast.BinaryExpr{X: c.CopyExpr(e.X), OpPos: e.OpPos, Op: e.Op, Y: c.CopyExpr(e.Y)}
	case *ast.UnaryExpr:
		return &ast.UnaryExpr{OpPos: e.OpPos, Op: e.Op, X: c.CopyExpr(e.X)}
	case *ast.CompositeLit:
		newE := &ast.CompositeLit{
			Type:   c.CopyExpr(e.Type),
			Lbrace: e.Lbrace,
			Rbrace: e.Rbrace,
		}
		for _, el := range e.Elts {
			newE.Elts = append(newE.Elts, c.CopyExpr(el))
		}
		return newE
	case *ast.FuncLit:
		return &ast.FuncLit{Type: c.CopyFuncType(e.Type), Body: c.CopyBlockStmt(e.Body)}
	case *ast.StructType:
		newE := &ast.StructType{
			Struct:     e.Struct,
			Fields:     c.CopyFieldList(e.Fields),
			Incomplete: e.Incomplete,
		}
		return newE
	case *ast.InterfaceType:
		newE := &ast.InterfaceType{
			Interface:  e.Interface,
			Methods:    c.CopyFieldList(e.Methods),
			Incomplete: e.Incomplete,
		}
		return newE
	case *ast.MapType:
		return &ast.MapType{Map: e.Map, Key: c.CopyExpr(e.Key), Value: c.CopyExpr(e.Value)}
	case *ast.ChanType:
		return &ast.ChanType{Begin: e.Begin, Arrow: e.Arrow, Dir: e.Dir, Value: c.CopyExpr(e.Value)}
	case *ast.FuncType:
		return c.CopyFuncType(e)
	default:
		// TODO: Handle other expressions (InterfaceType, StructType, etc if they appear in bodies/signatures we care about)
		// For now, return as is (risky if modified in place later) or implement more.
		return e
	}
}

func (c *DeepCopier) CopyStmt(s ast.Stmt) ast.Stmt {
	if s == nil {
		return nil
	}
	switch s := s.(type) {
	case *ast.DeclStmt:
		return &ast.DeclStmt{Decl: c.CopyDecl(s.Decl)}
	case *ast.ExprStmt:
		return &ast.ExprStmt{X: c.CopyExpr(s.X)}
	case *ast.AssignStmt:
		newS := &ast.AssignStmt{TokPos: s.TokPos, Tok: s.Tok}
		for _, lhs := range s.Lhs {
			newS.Lhs = append(newS.Lhs, c.CopyExpr(lhs))
		}
		for _, rhs := range s.Rhs {
			newS.Rhs = append(newS.Rhs, c.CopyExpr(rhs))
		}
		return newS
	case *ast.ReturnStmt:
		newS := &ast.ReturnStmt{Return: s.Return}
		for _, r := range s.Results {
			newS.Results = append(newS.Results, c.CopyExpr(r))
		}
		return newS
	case *ast.BlockStmt:
		return c.CopyBlockStmt(s)
	case *ast.IfStmt:
		return &ast.IfStmt{
			If:   s.If,
			Init: c.CopyStmt(s.Init),
			Cond: c.CopyExpr(s.Cond),
			Body: c.CopyBlockStmt(s.Body),
			Else: c.CopyStmt(s.Else),
		}
	case *ast.ForStmt:
		return &ast.ForStmt{
			For:  s.For,
			Init: c.CopyStmt(s.Init),
			Cond: c.CopyExpr(s.Cond),
			Post: c.CopyStmt(s.Post),
			Body: c.CopyBlockStmt(s.Body),
		}
	case *ast.RangeStmt:
		return &ast.RangeStmt{
			For:    s.For,
			Key:    c.CopyExpr(s.Key),
			Value:  c.CopyExpr(s.Value),
			TokPos: s.TokPos,
			Tok:    s.Tok,
			X:      c.CopyExpr(s.X),
			Body:   c.CopyBlockStmt(s.Body),
		}
	case *ast.SwitchStmt:
		newS := &ast.SwitchStmt{
			Switch: s.Switch,
			Init:   c.CopyStmt(s.Init),
			Tag:    c.CopyExpr(s.Tag),
			Body:   c.CopyBlockStmt(s.Body),
		}
		return newS
	case *ast.CaseClause:
		newS := &ast.CaseClause{
			Case:  s.Case,
			Colon: s.Colon,
		}
		for _, l := range s.List {
			newS.List = append(newS.List, c.CopyExpr(l))
		}
		for _, b := range s.Body {
			newS.Body = append(newS.Body, c.CopyStmt(b))
		}
		return newS
	default:
		// GoStmt, DeferStmt, etc.
		return s
	}
}

func (c *DeepCopier) CopyBlockStmt(b *ast.BlockStmt) *ast.BlockStmt {
	if b == nil {
		return nil
	}
	newB := &ast.BlockStmt{Lbrace: b.Lbrace, Rbrace: b.Rbrace}

		assertName := fmt.Sprintf("Assert%d", c.VecLen)
		assertCall := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "midway"},
					Sel: &ast.Ident{Name: assertName},
				},
			},
		}
		newB.List = append(newB.List, assertCall)

	for _, s := range b.List {
		newB.List = append(newB.List, c.CopyStmt(s))
	}
	return newB
}

func (c *DeepCopier) CopyFieldList(f *ast.FieldList) *ast.FieldList {
	if f == nil {
		return nil
	}
	newF := &ast.FieldList{Opening: f.Opening, Closing: f.Closing}
	for _, field := range f.List {
		newF.List = append(newF.List, c.CopyField(field))
	}
	return newF
}

func (c *DeepCopier) CopyField(f *ast.Field) *ast.Field {
	newF := &ast.Field{
		Doc:     c.CopyCommentGroup(f.Doc),
		Comment: c.CopyCommentGroup(f.Comment),
		Type:    c.CopyExpr(f.Type),
		Tag:     f.Tag,
	}
	for _, n := range f.Names {
		newF.Names = append(newF.Names, c.CopyIdent(n))
	}
	return newF
}

func (c *DeepCopier) CopyFuncType(t *ast.FuncType) *ast.FuncType {
	if t == nil {
		return nil
	}
	return &ast.FuncType{
		Func:       t.Func,
		Params:     c.CopyFieldList(t.Params),
		Results:    c.CopyFieldList(t.Results),
		TypeParams: c.CopyFieldList(t.TypeParams),
	}
}

func (c *DeepCopier) CopyIdent(id *ast.Ident) *ast.Ident {
	if id == nil {
		return nil
	}
	if c.OnIdent != nil {
		if match := c.OnIdent(id); match != nil {
			return match
		}
	}
	newId := &ast.Ident{
		NamePos: id.NamePos,
		Name:    id.Name,
		Obj:     id.Obj, // Note: We keep the Obj reference, but later logic must depend on the original object
	}
	return newId
}

func (c *DeepCopier) CopyCommentGroup(cg *ast.CommentGroup) *ast.CommentGroup {
	if cg == nil {
		return nil
	}
	// Shallow copy comments is usually fine
	newCg := *cg
	return &newCg
}
