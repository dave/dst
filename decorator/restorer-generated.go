package decorator

import (
	"fmt"
	"github.com/dave/dst"
	"go/ast"
	"go/token"
)

func (r *fileRestorer) restoreNode(n dst.Node) ast.Node {
	if an, ok := r.Nodes[n]; ok {
		return an
	}
	switch n := n.(type) {
	case *dst.ArrayType:
		out := &ast.ArrayType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Lbrack
		out.Lbrack = r.cursor
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Lbrack
		r.applyDecorations(n.Decs.Lbrack, true)

		// Node: Len
		if n.Len != nil {
			out.Len = r.restoreNode(n.Len).(ast.Expr)
		}

		// Token: Rbrack
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: Len
		r.applyDecorations(n.Decs.Len, true)

		// Node: Elt
		if n.Elt != nil {
			out.Elt = r.restoreNode(n.Elt).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.AssignStmt:
		out := &ast.AssignStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// List: Lhs
		for _, v := range n.Lhs {
			out.Lhs = append(out.Lhs, r.restoreNode(v).(ast.Expr))
		}

		// Decoration: Lhs
		r.applyDecorations(n.Decs.Lhs, true)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: Tok
		r.applyDecorations(n.Decs.Tok, true)

		// List: Rhs
		for _, v := range n.Rhs {
			out.Rhs = append(out.Rhs, r.restoreNode(v).(ast.Expr))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BadDecl:
		out := &ast.BadDecl{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BadExpr:
		out := &ast.BadExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BadStmt:
		out := &ast.BadStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BasicLit:
		out := &ast.BasicLit{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// String: Value
		out.ValuePos = r.cursor
		out.Value = n.Value
		r.cursor += token.Pos(len(n.Value))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Kind
		out.Kind = n.Kind
		r.applySpace(n.Decs.After)

		return out
	case *dst.BinaryExpr:
		out := &ast.BinaryExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Token: Op
		out.Op = n.Op
		out.OpPos = r.cursor
		r.cursor += token.Pos(len(n.Op.String()))

		// Decoration: Op
		r.applyDecorations(n.Decs.Op, true)

		// Node: Y
		if n.Y != nil {
			out.Y = r.restoreNode(n.Y).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BlockStmt:
		out := &ast.BlockStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Lbrace
		out.Lbrace = r.cursor
		r.cursor += token.Pos(len(token.LBRACE.String()))

		// Decoration: Lbrace
		r.applyDecorations(n.Decs.Lbrace, true)

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, r.restoreNode(v).(ast.Stmt))
		}

		// Token: Rbrace
		out.Rbrace = r.cursor
		r.cursor += token.Pos(len(token.RBRACE.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BranchStmt:
		out := &ast.BranchStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: Tok
		r.applyDecorations(n.Decs.Tok, true)

		// Node: Label
		if n.Label != nil {
			out.Label = r.restoreNode(n.Label).(*ast.Ident)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CallExpr:
		out := &ast.CallExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Fun
		if n.Fun != nil {
			out.Fun = r.restoreNode(n.Fun).(ast.Expr)
		}

		// Decoration: Fun
		r.applyDecorations(n.Decs.Fun, true)

		// Token: Lparen
		out.Lparen = r.cursor
		r.cursor += token.Pos(len(token.LPAREN.String()))

		// Decoration: Lparen
		r.applyDecorations(n.Decs.Lparen, true)

		// List: Args
		for _, v := range n.Args {
			out.Args = append(out.Args, r.restoreNode(v).(ast.Expr))
		}

		// Decoration: Args
		r.applyDecorations(n.Decs.Args, true)

		// Token: Ellipsis
		if n.Ellipsis {
			out.Ellipsis = r.cursor
			r.cursor += token.Pos(len(token.ELLIPSIS.String()))
		}

		// Decoration: Ellipsis
		r.applyDecorations(n.Decs.Ellipsis, true)

		// Token: Rparen
		out.Rparen = r.cursor
		r.cursor += token.Pos(len(token.RPAREN.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CaseClause:
		out := &ast.CaseClause{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, false)

		// Token: Case
		out.Case = r.cursor
		r.cursor += token.Pos(len(func() token.Token {
			if n.List == nil {
				return token.DEFAULT
			} else {
				return token.CASE
			}
		}().String()))

		// Decoration: Case
		r.applyDecorations(n.Decs.Case, true)

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, r.restoreNode(v).(ast.Expr))
		}

		// Decoration: List
		r.applyDecorations(n.Decs.List, true)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(n.Decs.Colon, true)

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, r.restoreNode(v).(ast.Stmt))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ChanType:
		out := &ast.ChanType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Begin
		out.Begin = r.cursor
		r.cursor += token.Pos(len(func() token.Token {
			if n.Dir == dst.RECV {
				return token.ARROW
			} else {
				return token.CHAN
			}
		}().String()))

		// Token: Chan
		if n.Dir == dst.RECV {
			r.cursor += token.Pos(len(token.CHAN.String()))
		}

		// Decoration: Begin
		r.applyDecorations(n.Decs.Begin, true)

		// Token: Arrow
		if n.Dir == dst.SEND {
			out.Arrow = r.cursor
			r.cursor += token.Pos(len(token.ARROW.String()))
		}

		// Decoration: Arrow
		r.applyDecorations(n.Decs.Arrow, true)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Dir
		out.Dir = ast.ChanDir(n.Dir)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CommClause:
		out := &ast.CommClause{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, false)

		// Token: Case
		out.Case = r.cursor
		r.cursor += token.Pos(len(func() token.Token {
			if n.Comm == nil {
				return token.DEFAULT
			} else {
				return token.CASE
			}
		}().String()))

		// Decoration: Case
		r.applyDecorations(n.Decs.Case, true)

		// Node: Comm
		if n.Comm != nil {
			out.Comm = r.restoreNode(n.Comm).(ast.Stmt)
		}

		// Decoration: Comm
		r.applyDecorations(n.Decs.Comm, true)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(n.Decs.Colon, true)

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, r.restoreNode(v).(ast.Stmt))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CompositeLit:
		out := &ast.CompositeLit{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type).(ast.Expr)
		}

		// Decoration: Type
		r.applyDecorations(n.Decs.Type, true)

		// Token: Lbrace
		out.Lbrace = r.cursor
		r.cursor += token.Pos(len(token.LBRACE.String()))

		// Decoration: Lbrace
		r.applyDecorations(n.Decs.Lbrace, true)

		// List: Elts
		for _, v := range n.Elts {
			out.Elts = append(out.Elts, r.restoreNode(v).(ast.Expr))
		}

		// Token: Rbrace
		out.Rbrace = r.cursor
		r.cursor += token.Pos(len(token.RBRACE.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Incomplete
		out.Incomplete = n.Incomplete
		r.applySpace(n.Decs.After)

		return out
	case *dst.DeclStmt:
		out := &ast.DeclStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Decl
		if n.Decl != nil {
			out.Decl = r.restoreNode(n.Decl).(ast.Decl)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.DeferStmt:
		out := &ast.DeferStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Defer
		out.Defer = r.cursor
		r.cursor += token.Pos(len(token.DEFER.String()))

		// Decoration: Defer
		r.applyDecorations(n.Decs.Defer, true)

		// Node: Call
		if n.Call != nil {
			out.Call = r.restoreNode(n.Call).(*ast.CallExpr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Ellipsis:
		out := &ast.Ellipsis{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Ellipsis
		out.Ellipsis = r.cursor
		r.cursor += token.Pos(len(token.ELLIPSIS.String()))

		// Decoration: Ellipsis
		r.applyDecorations(n.Decs.Ellipsis, true)

		// Node: Elt
		if n.Elt != nil {
			out.Elt = r.restoreNode(n.Elt).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.EmptyStmt:
		out := &ast.EmptyStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Semicolon
		if !n.Implicit {
			out.Semicolon = r.cursor
			r.cursor += token.Pos(len(token.ARROW.String()))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Implicit
		out.Implicit = n.Implicit
		r.applySpace(n.Decs.After)

		return out
	case *dst.ExprStmt:
		out := &ast.ExprStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Field:
		out := &ast.Field{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, r.restoreNode(v).(*ast.Ident))
		}

		// Decoration: Names
		r.applyDecorations(n.Decs.Names, true)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type).(ast.Expr)
		}

		// Decoration: Type
		r.applyDecorations(n.Decs.Type, true)

		// Node: Tag
		if n.Tag != nil {
			out.Tag = r.restoreNode(n.Tag).(*ast.BasicLit)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FieldList:
		out := &ast.FieldList{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Opening
		if n.Opening {
			out.Opening = r.cursor
			r.cursor += token.Pos(len(token.LPAREN.String()))
		}

		// Decoration: Opening
		r.applyDecorations(n.Decs.Opening, true)

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, r.restoreNode(v).(*ast.Field))
		}

		// Token: Closing
		if n.Closing {
			out.Closing = r.cursor
			r.cursor += token.Pos(len(token.RPAREN.String()))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.File:
		out := &ast.File{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Package
		out.Package = r.cursor
		r.cursor += token.Pos(len(token.PACKAGE.String()))

		// Decoration: Package
		r.applyDecorations(n.Decs.Package, true)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name).(*ast.Ident)
		}

		// Decoration: Name
		r.applyDecorations(n.Decs.Name, true)

		// List: Decls
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, r.restoreNode(v).(ast.Decl))
		}

		// Scope: Scope
		out.Scope = r.restoreScope(n.Scope)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ForStmt:
		out := &ast.ForStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: For
		out.For = r.cursor
		r.cursor += token.Pos(len(token.FOR.String()))

		// Decoration: For
		r.applyDecorations(n.Decs.For, true)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init).(ast.Stmt)
		}

		// Token: InitSemicolon
		if n.Init != nil {
			r.cursor += token.Pos(len(token.SEMICOLON.String()))
		}

		// Decoration: Init
		r.applyDecorations(n.Decs.Init, true)

		// Node: Cond
		if n.Cond != nil {
			out.Cond = r.restoreNode(n.Cond).(ast.Expr)
		}

		// Token: CondSemicolon
		if n.Post != nil {
			r.cursor += token.Pos(len(token.SEMICOLON.String()))
		}

		// Decoration: Cond
		r.applyDecorations(n.Decs.Cond, true)

		// Node: Post
		if n.Post != nil {
			out.Post = r.restoreNode(n.Post).(ast.Stmt)
		}

		// Decoration: Post
		r.applyDecorations(n.Decs.Post, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FuncDecl:
		out := &ast.FuncDecl{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Init: Type
		out.Type = &ast.FuncType{}

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Func
		if true {
			out.Type.Func = r.cursor
			r.cursor += token.Pos(len(token.FUNC.String()))
		}

		// Decoration: Func
		r.applyDecorations(n.Decs.Func, true)

		// Node: Recv
		if n.Recv != nil {
			out.Recv = r.restoreNode(n.Recv).(*ast.FieldList)
		}

		// Decoration: Recv
		r.applyDecorations(n.Decs.Recv, true)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name).(*ast.Ident)
		}

		// Decoration: Name
		r.applyDecorations(n.Decs.Name, true)

		// Node: Params
		if n.Type.Params != nil {
			out.Type.Params = r.restoreNode(n.Type.Params).(*ast.FieldList)
		}

		// Decoration: Params
		r.applyDecorations(n.Decs.Params, true)

		// Node: Results
		if n.Type.Results != nil {
			out.Type.Results = r.restoreNode(n.Type.Results).(*ast.FieldList)
		}

		// Decoration: Results
		r.applyDecorations(n.Decs.Results, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FuncLit:
		out := &ast.FuncLit{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type).(*ast.FuncType)
		}

		// Decoration: Type
		r.applyDecorations(n.Decs.Type, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FuncType:
		out := &ast.FuncType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Func
		if n.Func {
			out.Func = r.cursor
			r.cursor += token.Pos(len(token.FUNC.String()))
		}

		// Decoration: Func
		r.applyDecorations(n.Decs.Func, true)

		// Node: Params
		if n.Params != nil {
			out.Params = r.restoreNode(n.Params).(*ast.FieldList)
		}

		// Decoration: Params
		r.applyDecorations(n.Decs.Params, true)

		// Node: Results
		if n.Results != nil {
			out.Results = r.restoreNode(n.Results).(*ast.FieldList)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.GenDecl:
		out := &ast.GenDecl{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: Tok
		r.applyDecorations(n.Decs.Tok, true)

		// Token: Lparen
		if n.Lparen {
			out.Lparen = r.cursor
			r.cursor += token.Pos(len(token.LPAREN.String()))
		}

		// Decoration: Lparen
		r.applyDecorations(n.Decs.Lparen, true)

		// List: Specs
		for _, v := range n.Specs {
			out.Specs = append(out.Specs, r.restoreNode(v).(ast.Spec))
		}

		// Token: Rparen
		if n.Rparen {
			out.Rparen = r.cursor
			r.cursor += token.Pos(len(token.RPAREN.String()))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.GoStmt:
		out := &ast.GoStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Go
		out.Go = r.cursor
		r.cursor += token.Pos(len(token.GO.String()))

		// Decoration: Go
		r.applyDecorations(n.Decs.Go, true)

		// Node: Call
		if n.Call != nil {
			out.Call = r.restoreNode(n.Call).(*ast.CallExpr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Ident:
		out := &ast.Ident{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// String: Name
		out.NamePos = r.cursor
		out.Name = n.Name
		r.cursor += token.Pos(len(n.Name))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Object: Obj
		out.Obj = r.restoreObject(n.Obj)
		r.applySpace(n.Decs.After)

		return out
	case *dst.IfStmt:
		out := &ast.IfStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: If
		out.If = r.cursor
		r.cursor += token.Pos(len(token.IF.String()))

		// Decoration: If
		r.applyDecorations(n.Decs.If, true)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init).(ast.Stmt)
		}

		// Decoration: Init
		r.applyDecorations(n.Decs.Init, true)

		// Node: Cond
		if n.Cond != nil {
			out.Cond = r.restoreNode(n.Cond).(ast.Expr)
		}

		// Decoration: Cond
		r.applyDecorations(n.Decs.Cond, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Token: ElseTok
		if n.Else != nil {
			r.cursor += token.Pos(len(token.ELSE.String()))
		}

		// Decoration: Else
		r.applyDecorations(n.Decs.Else, true)

		// Node: Else
		if n.Else != nil {
			out.Else = r.restoreNode(n.Else).(ast.Stmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ImportSpec:
		out := &ast.ImportSpec{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name).(*ast.Ident)
		}

		// Decoration: Name
		r.applyDecorations(n.Decs.Name, true)

		// Node: Path
		if n.Path != nil {
			out.Path = r.restoreNode(n.Path).(*ast.BasicLit)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.IncDecStmt:
		out := &ast.IncDecStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.IndexExpr:
		out := &ast.IndexExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Token: Lbrack
		out.Lbrack = r.cursor
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Lbrack
		r.applyDecorations(n.Decs.Lbrack, true)

		// Node: Index
		if n.Index != nil {
			out.Index = r.restoreNode(n.Index).(ast.Expr)
		}

		// Decoration: Index
		r.applyDecorations(n.Decs.Index, true)

		// Token: Rbrack
		out.Rbrack = r.cursor
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.InterfaceType:
		out := &ast.InterfaceType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Interface
		out.Interface = r.cursor
		r.cursor += token.Pos(len(token.INTERFACE.String()))

		// Decoration: Interface
		r.applyDecorations(n.Decs.Interface, true)

		// Node: Methods
		if n.Methods != nil {
			out.Methods = r.restoreNode(n.Methods).(*ast.FieldList)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Incomplete
		out.Incomplete = n.Incomplete
		r.applySpace(n.Decs.After)

		return out
	case *dst.KeyValueExpr:
		out := &ast.KeyValueExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Key
		if n.Key != nil {
			out.Key = r.restoreNode(n.Key).(ast.Expr)
		}

		// Decoration: Key
		r.applyDecorations(n.Decs.Key, true)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(n.Decs.Colon, true)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.LabeledStmt:
		out := &ast.LabeledStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Label
		if n.Label != nil {
			out.Label = r.restoreNode(n.Label).(*ast.Ident)
		}

		// Decoration: Label
		r.applyDecorations(n.Decs.Label, true)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(n.Decs.Colon, true)

		// Node: Stmt
		if n.Stmt != nil {
			out.Stmt = r.restoreNode(n.Stmt).(ast.Stmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.MapType:
		out := &ast.MapType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Map
		out.Map = r.cursor
		r.cursor += token.Pos(len(token.MAP.String()))

		// Token: Lbrack
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Map
		r.applyDecorations(n.Decs.Map, true)

		// Node: Key
		if n.Key != nil {
			out.Key = r.restoreNode(n.Key).(ast.Expr)
		}

		// Token: Rbrack
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: Key
		r.applyDecorations(n.Decs.Key, true)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Package:
		out := &ast.Package{}
		r.Nodes[n] = out

		// Value: Name
		out.Name = n.Name

		// Scope: Scope
		out.Scope = r.restoreScope(n.Scope)

		// Map: Imports
		out.Imports = map[string]*ast.Object{}
		for k, v := range n.Imports {
			out.Imports[k] = r.restoreObject(v)
		}

		// Map: Files
		out.Files = map[string]*ast.File{}
		for k, v := range n.Files {
			out.Files[k] = r.restoreNode(v).(*ast.File)
		}

		return out
	case *dst.ParenExpr:
		out := &ast.ParenExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Lparen
		out.Lparen = r.cursor
		r.cursor += token.Pos(len(token.LPAREN.String()))

		// Decoration: Lparen
		r.applyDecorations(n.Decs.Lparen, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Token: Rparen
		out.Rparen = r.cursor
		r.cursor += token.Pos(len(token.RPAREN.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.RangeStmt:
		out := &ast.RangeStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: For
		out.For = r.cursor
		r.cursor += token.Pos(len(token.FOR.String()))

		// Decoration: For
		r.applyDecorations(n.Decs.For, true)

		// Node: Key
		if n.Key != nil {
			out.Key = r.restoreNode(n.Key).(ast.Expr)
		}

		// Token: Comma
		if n.Value != nil {
			r.cursor += token.Pos(len(token.COMMA.String()))
		}

		// Decoration: Key
		r.applyDecorations(n.Decs.Key, true)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value).(ast.Expr)
		}

		// Decoration: Value
		r.applyDecorations(n.Decs.Value, true)

		// Token: Tok
		if n.Tok != token.ILLEGAL {
			out.Tok = n.Tok
			out.TokPos = r.cursor
			r.cursor += token.Pos(len(n.Tok.String()))
		}

		// Token: Range
		r.cursor += token.Pos(len(token.RANGE.String()))

		// Decoration: Range
		r.applyDecorations(n.Decs.Range, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ReturnStmt:
		out := &ast.ReturnStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Return
		out.Return = r.cursor
		r.cursor += token.Pos(len(token.RETURN.String()))

		// Decoration: Return
		r.applyDecorations(n.Decs.Return, true)

		// List: Results
		for _, v := range n.Results {
			out.Results = append(out.Results, r.restoreNode(v).(ast.Expr))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SelectStmt:
		out := &ast.SelectStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Select
		out.Select = r.cursor
		r.cursor += token.Pos(len(token.SELECT.String()))

		// Decoration: Select
		r.applyDecorations(n.Decs.Select, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SelectorExpr:
		out := &ast.SelectorExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Token: Period
		r.cursor += token.Pos(len(token.PERIOD.String()))

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Node: Sel
		if n.Sel != nil {
			out.Sel = r.restoreNode(n.Sel).(*ast.Ident)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SendStmt:
		out := &ast.SendStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Chan
		if n.Chan != nil {
			out.Chan = r.restoreNode(n.Chan).(ast.Expr)
		}

		// Decoration: Chan
		r.applyDecorations(n.Decs.Chan, true)

		// Token: Arrow
		out.Arrow = r.cursor
		r.cursor += token.Pos(len(token.ARROW.String()))

		// Decoration: Arrow
		r.applyDecorations(n.Decs.Arrow, true)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SliceExpr:
		out := &ast.SliceExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Token: Lbrack
		out.Lbrack = r.cursor
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Lbrack
		r.applyDecorations(n.Decs.Lbrack, true)

		// Node: Low
		if n.Low != nil {
			out.Low = r.restoreNode(n.Low).(ast.Expr)
		}

		// Token: Colon1
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Low
		r.applyDecorations(n.Decs.Low, true)

		// Node: High
		if n.High != nil {
			out.High = r.restoreNode(n.High).(ast.Expr)
		}

		// Token: Colon2
		if n.Slice3 {
			r.cursor += token.Pos(len(token.COLON.String()))
		}

		// Decoration: High
		r.applyDecorations(n.Decs.High, true)

		// Node: Max
		if n.Max != nil {
			out.Max = r.restoreNode(n.Max).(ast.Expr)
		}

		// Decoration: Max
		r.applyDecorations(n.Decs.Max, true)

		// Token: Rbrack
		out.Rbrack = r.cursor
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Slice3
		out.Slice3 = n.Slice3
		r.applySpace(n.Decs.After)

		return out
	case *dst.StarExpr:
		out := &ast.StarExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Star
		out.Star = r.cursor
		r.cursor += token.Pos(len(token.MUL.String()))

		// Decoration: Star
		r.applyDecorations(n.Decs.Star, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.StructType:
		out := &ast.StructType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Struct
		out.Struct = r.cursor
		r.cursor += token.Pos(len(token.STRUCT.String()))

		// Decoration: Struct
		r.applyDecorations(n.Decs.Struct, true)

		// Node: Fields
		if n.Fields != nil {
			out.Fields = r.restoreNode(n.Fields).(*ast.FieldList)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)

		// Value: Incomplete
		out.Incomplete = n.Incomplete
		r.applySpace(n.Decs.After)

		return out
	case *dst.SwitchStmt:
		out := &ast.SwitchStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Switch
		out.Switch = r.cursor
		r.cursor += token.Pos(len(token.SWITCH.String()))

		// Decoration: Switch
		r.applyDecorations(n.Decs.Switch, true)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init).(ast.Stmt)
		}

		// Decoration: Init
		r.applyDecorations(n.Decs.Init, true)

		// Node: Tag
		if n.Tag != nil {
			out.Tag = r.restoreNode(n.Tag).(ast.Expr)
		}

		// Decoration: Tag
		r.applyDecorations(n.Decs.Tag, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.TypeAssertExpr:
		out := &ast.TypeAssertExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Token: Period
		r.cursor += token.Pos(len(token.PERIOD.String()))

		// Decoration: X
		r.applyDecorations(n.Decs.X, true)

		// Token: Lparen
		out.Lparen = r.cursor
		r.cursor += token.Pos(len(token.LPAREN.String()))

		// Decoration: Lparen
		r.applyDecorations(n.Decs.Lparen, true)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type).(ast.Expr)
		}

		// Token: TypeToken
		if n.Type == nil {
			r.cursor += token.Pos(len(token.TYPE.String()))
		}

		// Decoration: Type
		r.applyDecorations(n.Decs.Type, true)

		// Token: Rparen
		out.Rparen = r.cursor
		r.cursor += token.Pos(len(token.RPAREN.String()))

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.TypeSpec:
		out := &ast.TypeSpec{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name).(*ast.Ident)
		}

		// Token: Assign
		if n.Assign {
			out.Assign = r.cursor
			r.cursor += token.Pos(len(token.ASSIGN.String()))
		}

		// Decoration: Name
		r.applyDecorations(n.Decs.Name, true)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.TypeSwitchStmt:
		out := &ast.TypeSwitchStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Switch
		out.Switch = r.cursor
		r.cursor += token.Pos(len(token.SWITCH.String()))

		// Decoration: Switch
		r.applyDecorations(n.Decs.Switch, true)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init).(ast.Stmt)
		}

		// Decoration: Init
		r.applyDecorations(n.Decs.Init, true)

		// Node: Assign
		if n.Assign != nil {
			out.Assign = r.restoreNode(n.Assign).(ast.Stmt)
		}

		// Decoration: Assign
		r.applyDecorations(n.Decs.Assign, true)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.UnaryExpr:
		out := &ast.UnaryExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// Token: Op
		out.Op = n.Op
		out.OpPos = r.cursor
		r.cursor += token.Pos(len(n.Op.String()))

		// Decoration: Op
		r.applyDecorations(n.Decs.Op, true)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ValueSpec:
		out := &ast.ValueSpec{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(n.Decs.Start, true)

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, r.restoreNode(v).(*ast.Ident))
		}

		// Decoration: Names
		r.applyDecorations(n.Decs.Names, true)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type).(ast.Expr)
		}

		// Token: Assign
		if n.Values != nil {
			r.cursor += token.Pos(len(token.ASSIGN.String()))
		}

		// Decoration: Assign
		r.applyDecorations(n.Decs.Assign, true)

		// List: Values
		for _, v := range n.Values {
			out.Values = append(out.Values, r.restoreNode(v).(ast.Expr))
		}

		// Decoration: End
		r.applyDecorations(n.Decs.End, true)
		r.applySpace(n.Decs.After)

		return out
	default:
		panic(fmt.Sprintf("%T", n))
	}
	return nil
}
