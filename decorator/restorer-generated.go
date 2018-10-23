package decorator

import (
	"fmt"
	"github.com/dave/dst"
	"go/ast"
	"go/token"
)

func (r *fileRestorer) restoreNode(n dst.Node, allowDuplicate bool) ast.Node {
	if an, ok := r.Nodes[n]; ok {
		if allowDuplicate {
			return an
		} else {
			panic(fmt.Sprintf("duplicate node: %#v", n))
		}
	}
	switch n := n.(type) {
	case *dst.ArrayType:
		out := &ast.ArrayType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Lbrack
		out.Lbrack = r.cursor
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Lbrack
		r.applyDecorations(out, "Lbrack", n.Decs.Lbrack)

		// Node: Len
		if n.Len != nil {
			out.Len = r.restoreNode(n.Len, allowDuplicate).(ast.Expr)
		}

		// Token: Rbrack
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: Len
		r.applyDecorations(out, "Len", n.Decs.Len)

		// Node: Elt
		if n.Elt != nil {
			out.Elt = r.restoreNode(n.Elt, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.AssignStmt:
		out := &ast.AssignStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// List: Lhs
		for _, v := range n.Lhs {
			out.Lhs = append(out.Lhs, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Decoration: Lhs
		r.applyDecorations(out, "Lhs", n.Decs.Lhs)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: Tok
		r.applyDecorations(out, "Tok", n.Decs.Tok)

		// List: Rhs
		for _, v := range n.Rhs {
			out.Rhs = append(out.Rhs, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BadDecl:
		out := &ast.BadDecl{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BadExpr:
		out := &ast.BadExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BadStmt:
		out := &ast.BadStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BasicLit:
		out := &ast.BasicLit{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// String: Value
		r.applyLiteral(n.Value)
		out.ValuePos = r.cursor
		out.Value = n.Value
		r.cursor += token.Pos(len(n.Value))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Kind
		out.Kind = n.Kind
		r.applySpace(n.Decs.After)

		return out
	case *dst.BinaryExpr:
		out := &ast.BinaryExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Token: Op
		out.Op = n.Op
		out.OpPos = r.cursor
		r.cursor += token.Pos(len(n.Op.String()))

		// Decoration: Op
		r.applyDecorations(out, "Op", n.Decs.Op)

		// Node: Y
		if n.Y != nil {
			out.Y = r.restoreNode(n.Y, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BlockStmt:
		out := &ast.BlockStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Lbrace
		out.Lbrace = r.cursor
		r.cursor += token.Pos(len(token.LBRACE.String()))

		// Decoration: Lbrace
		r.applyDecorations(out, "Lbrace", n.Decs.Lbrace)

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, r.restoreNode(v, allowDuplicate).(ast.Stmt))
		}

		// Token: Rbrace
		out.Rbrace = r.cursor
		r.cursor += token.Pos(len(token.RBRACE.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.BranchStmt:
		out := &ast.BranchStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: Tok
		r.applyDecorations(out, "Tok", n.Decs.Tok)

		// Node: Label
		if n.Label != nil {
			out.Label = r.restoreNode(n.Label, allowDuplicate).(*ast.Ident)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CallExpr:
		out := &ast.CallExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Fun
		if n.Fun != nil {
			out.Fun = r.restoreNode(n.Fun, allowDuplicate).(ast.Expr)
		}

		// Decoration: Fun
		r.applyDecorations(out, "Fun", n.Decs.Fun)

		// Token: Lparen
		out.Lparen = r.cursor
		r.cursor += token.Pos(len(token.LPAREN.String()))

		// Decoration: Lparen
		r.applyDecorations(out, "Lparen", n.Decs.Lparen)

		// List: Args
		for _, v := range n.Args {
			out.Args = append(out.Args, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Decoration: Args
		r.applyDecorations(out, "Args", n.Decs.Args)

		// Token: Ellipsis
		if n.Ellipsis {
			out.Ellipsis = r.cursor
			r.cursor += token.Pos(len(token.ELLIPSIS.String()))
		}

		// Decoration: Ellipsis
		r.applyDecorations(out, "Ellipsis", n.Decs.Ellipsis)

		// Token: Rparen
		out.Rparen = r.cursor
		r.cursor += token.Pos(len(token.RPAREN.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CaseClause:
		out := &ast.CaseClause{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

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
		r.applyDecorations(out, "Case", n.Decs.Case)

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Decoration: List
		r.applyDecorations(out, "List", n.Decs.List)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(out, "Colon", n.Decs.Colon)

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, r.restoreNode(v, allowDuplicate).(ast.Stmt))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ChanType:
		out := &ast.ChanType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

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
		r.applyDecorations(out, "Begin", n.Decs.Begin)

		// Token: Arrow
		if n.Dir == dst.SEND {
			out.Arrow = r.cursor
			r.cursor += token.Pos(len(token.ARROW.String()))
		}

		// Decoration: Arrow
		r.applyDecorations(out, "Arrow", n.Decs.Arrow)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Dir
		out.Dir = ast.ChanDir(n.Dir)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CommClause:
		out := &ast.CommClause{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

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
		r.applyDecorations(out, "Case", n.Decs.Case)

		// Node: Comm
		if n.Comm != nil {
			out.Comm = r.restoreNode(n.Comm, allowDuplicate).(ast.Stmt)
		}

		// Decoration: Comm
		r.applyDecorations(out, "Comm", n.Decs.Comm)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(out, "Colon", n.Decs.Colon)

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, r.restoreNode(v, allowDuplicate).(ast.Stmt))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.CompositeLit:
		out := &ast.CompositeLit{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type, allowDuplicate).(ast.Expr)
		}

		// Decoration: Type
		r.applyDecorations(out, "Type", n.Decs.Type)

		// Token: Lbrace
		out.Lbrace = r.cursor
		r.cursor += token.Pos(len(token.LBRACE.String()))

		// Decoration: Lbrace
		r.applyDecorations(out, "Lbrace", n.Decs.Lbrace)

		// List: Elts
		for _, v := range n.Elts {
			out.Elts = append(out.Elts, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Token: Rbrace
		out.Rbrace = r.cursor
		r.cursor += token.Pos(len(token.RBRACE.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Incomplete
		out.Incomplete = n.Incomplete
		r.applySpace(n.Decs.After)

		return out
	case *dst.DeclStmt:
		out := &ast.DeclStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Decl
		if n.Decl != nil {
			out.Decl = r.restoreNode(n.Decl, allowDuplicate).(ast.Decl)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.DeferStmt:
		out := &ast.DeferStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Defer
		out.Defer = r.cursor
		r.cursor += token.Pos(len(token.DEFER.String()))

		// Decoration: Defer
		r.applyDecorations(out, "Defer", n.Decs.Defer)

		// Node: Call
		if n.Call != nil {
			out.Call = r.restoreNode(n.Call, allowDuplicate).(*ast.CallExpr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Ellipsis:
		out := &ast.Ellipsis{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Ellipsis
		out.Ellipsis = r.cursor
		r.cursor += token.Pos(len(token.ELLIPSIS.String()))

		// Decoration: Ellipsis
		r.applyDecorations(out, "Ellipsis", n.Decs.Ellipsis)

		// Node: Elt
		if n.Elt != nil {
			out.Elt = r.restoreNode(n.Elt, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.EmptyStmt:
		out := &ast.EmptyStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Semicolon
		if !n.Implicit {
			out.Semicolon = r.cursor
			r.cursor += token.Pos(len(token.ARROW.String()))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Implicit
		out.Implicit = n.Implicit
		r.applySpace(n.Decs.After)

		return out
	case *dst.ExprStmt:
		out := &ast.ExprStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Field:
		out := &ast.Field{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, r.restoreNode(v, allowDuplicate).(*ast.Ident))
		}

		// Decoration: Names
		r.applyDecorations(out, "Names", n.Decs.Names)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type, allowDuplicate).(ast.Expr)
		}

		// Decoration: Type
		r.applyDecorations(out, "Type", n.Decs.Type)

		// Node: Tag
		if n.Tag != nil {
			out.Tag = r.restoreNode(n.Tag, allowDuplicate).(*ast.BasicLit)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FieldList:
		out := &ast.FieldList{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Opening
		if n.Opening {
			out.Opening = r.cursor
			r.cursor += token.Pos(len(token.LPAREN.String()))
		}

		// Decoration: Opening
		r.applyDecorations(out, "Opening", n.Decs.Opening)

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, r.restoreNode(v, allowDuplicate).(*ast.Field))
		}

		// Token: Closing
		if n.Closing {
			out.Closing = r.cursor
			r.cursor += token.Pos(len(token.RPAREN.String()))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.File:
		out := &ast.File{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Package
		out.Package = r.cursor
		r.cursor += token.Pos(len(token.PACKAGE.String()))

		// Decoration: Package
		r.applyDecorations(out, "Package", n.Decs.Package)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name, allowDuplicate).(*ast.Ident)
		}

		// Decoration: Name
		r.applyDecorations(out, "Name", n.Decs.Name)

		// List: Decls
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, r.restoreNode(v, allowDuplicate).(ast.Decl))
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
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: For
		out.For = r.cursor
		r.cursor += token.Pos(len(token.FOR.String()))

		// Decoration: For
		r.applyDecorations(out, "For", n.Decs.For)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init, allowDuplicate).(ast.Stmt)
		}

		// Token: InitSemicolon
		if n.Init != nil {
			r.cursor += token.Pos(len(token.SEMICOLON.String()))
		}

		// Decoration: Init
		r.applyDecorations(out, "Init", n.Decs.Init)

		// Node: Cond
		if n.Cond != nil {
			out.Cond = r.restoreNode(n.Cond, allowDuplicate).(ast.Expr)
		}

		// Token: CondSemicolon
		if n.Post != nil {
			r.cursor += token.Pos(len(token.SEMICOLON.String()))
		}

		// Decoration: Cond
		r.applyDecorations(out, "Cond", n.Decs.Cond)

		// Node: Post
		if n.Post != nil {
			out.Post = r.restoreNode(n.Post, allowDuplicate).(ast.Stmt)
		}

		// Decoration: Post
		r.applyDecorations(out, "Post", n.Decs.Post)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FuncDecl:
		out := &ast.FuncDecl{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Init: Type
		out.Type = &ast.FuncType{}

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Func
		if true {
			out.Type.Func = r.cursor
			r.cursor += token.Pos(len(token.FUNC.String()))
		}

		// Decoration: Func
		r.applyDecorations(out, "Func", n.Decs.Func)

		// Node: Recv
		if n.Recv != nil {
			out.Recv = r.restoreNode(n.Recv, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: Recv
		r.applyDecorations(out, "Recv", n.Decs.Recv)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name, allowDuplicate).(*ast.Ident)
		}

		// Decoration: Name
		r.applyDecorations(out, "Name", n.Decs.Name)

		// Node: Params
		if n.Type.Params != nil {
			out.Type.Params = r.restoreNode(n.Type.Params, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: Params
		r.applyDecorations(out, "Params", n.Decs.Params)

		// Node: Results
		if n.Type.Results != nil {
			out.Type.Results = r.restoreNode(n.Type.Results, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: Results
		r.applyDecorations(out, "Results", n.Decs.Results)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FuncLit:
		out := &ast.FuncLit{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type, allowDuplicate).(*ast.FuncType)
		}

		// Decoration: Type
		r.applyDecorations(out, "Type", n.Decs.Type)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.FuncType:
		out := &ast.FuncType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Func
		if n.Func {
			out.Func = r.cursor
			r.cursor += token.Pos(len(token.FUNC.String()))
		}

		// Decoration: Func
		r.applyDecorations(out, "Func", n.Decs.Func)

		// Node: Params
		if n.Params != nil {
			out.Params = r.restoreNode(n.Params, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: Params
		r.applyDecorations(out, "Params", n.Decs.Params)

		// Node: Results
		if n.Results != nil {
			out.Results = r.restoreNode(n.Results, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.GenDecl:
		out := &ast.GenDecl{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: Tok
		r.applyDecorations(out, "Tok", n.Decs.Tok)

		// Token: Lparen
		if n.Lparen {
			out.Lparen = r.cursor
			r.cursor += token.Pos(len(token.LPAREN.String()))
		}

		// Decoration: Lparen
		r.applyDecorations(out, "Lparen", n.Decs.Lparen)

		// List: Specs
		for _, v := range n.Specs {
			out.Specs = append(out.Specs, r.restoreNode(v, allowDuplicate).(ast.Spec))
		}

		// Token: Rparen
		if n.Rparen {
			out.Rparen = r.cursor
			r.cursor += token.Pos(len(token.RPAREN.String()))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.GoStmt:
		out := &ast.GoStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Go
		out.Go = r.cursor
		r.cursor += token.Pos(len(token.GO.String()))

		// Decoration: Go
		r.applyDecorations(out, "Go", n.Decs.Go)

		// Node: Call
		if n.Call != nil {
			out.Call = r.restoreNode(n.Call, allowDuplicate).(*ast.CallExpr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.Ident:
		out := &ast.Ident{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// String: Name
		out.NamePos = r.cursor
		out.Name = n.Name
		r.cursor += token.Pos(len(n.Name))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Object: Obj
		out.Obj = r.restoreObject(n.Obj)
		r.applySpace(n.Decs.After)

		return out
	case *dst.IfStmt:
		out := &ast.IfStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: If
		out.If = r.cursor
		r.cursor += token.Pos(len(token.IF.String()))

		// Decoration: If
		r.applyDecorations(out, "If", n.Decs.If)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init, allowDuplicate).(ast.Stmt)
		}

		// Decoration: Init
		r.applyDecorations(out, "Init", n.Decs.Init)

		// Node: Cond
		if n.Cond != nil {
			out.Cond = r.restoreNode(n.Cond, allowDuplicate).(ast.Expr)
		}

		// Decoration: Cond
		r.applyDecorations(out, "Cond", n.Decs.Cond)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Token: ElseTok
		if n.Else != nil {
			r.cursor += token.Pos(len(token.ELSE.String()))
		}

		// Decoration: Else
		r.applyDecorations(out, "Else", n.Decs.Else)

		// Node: Else
		if n.Else != nil {
			out.Else = r.restoreNode(n.Else, allowDuplicate).(ast.Stmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ImportSpec:
		out := &ast.ImportSpec{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name, allowDuplicate).(*ast.Ident)
		}

		// Decoration: Name
		r.applyDecorations(out, "Name", n.Decs.Name)

		// Node: Path
		if n.Path != nil {
			out.Path = r.restoreNode(n.Path, allowDuplicate).(*ast.BasicLit)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.IncDecStmt:
		out := &ast.IncDecStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Token: Tok
		out.Tok = n.Tok
		out.TokPos = r.cursor
		r.cursor += token.Pos(len(n.Tok.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.IndexExpr:
		out := &ast.IndexExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Token: Lbrack
		out.Lbrack = r.cursor
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Lbrack
		r.applyDecorations(out, "Lbrack", n.Decs.Lbrack)

		// Node: Index
		if n.Index != nil {
			out.Index = r.restoreNode(n.Index, allowDuplicate).(ast.Expr)
		}

		// Decoration: Index
		r.applyDecorations(out, "Index", n.Decs.Index)

		// Token: Rbrack
		out.Rbrack = r.cursor
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.InterfaceType:
		out := &ast.InterfaceType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Interface
		out.Interface = r.cursor
		r.cursor += token.Pos(len(token.INTERFACE.String()))

		// Decoration: Interface
		r.applyDecorations(out, "Interface", n.Decs.Interface)

		// Node: Methods
		if n.Methods != nil {
			out.Methods = r.restoreNode(n.Methods, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Incomplete
		out.Incomplete = n.Incomplete
		r.applySpace(n.Decs.After)

		return out
	case *dst.KeyValueExpr:
		out := &ast.KeyValueExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Key
		if n.Key != nil {
			out.Key = r.restoreNode(n.Key, allowDuplicate).(ast.Expr)
		}

		// Decoration: Key
		r.applyDecorations(out, "Key", n.Decs.Key)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(out, "Colon", n.Decs.Colon)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.LabeledStmt:
		out := &ast.LabeledStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Label
		if n.Label != nil {
			out.Label = r.restoreNode(n.Label, allowDuplicate).(*ast.Ident)
		}

		// Decoration: Label
		r.applyDecorations(out, "Label", n.Decs.Label)

		// Token: Colon
		out.Colon = r.cursor
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Colon
		r.applyDecorations(out, "Colon", n.Decs.Colon)

		// Node: Stmt
		if n.Stmt != nil {
			out.Stmt = r.restoreNode(n.Stmt, allowDuplicate).(ast.Stmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.MapType:
		out := &ast.MapType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Map
		out.Map = r.cursor
		r.cursor += token.Pos(len(token.MAP.String()))

		// Token: Lbrack
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Map
		r.applyDecorations(out, "Map", n.Decs.Map)

		// Node: Key
		if n.Key != nil {
			out.Key = r.restoreNode(n.Key, allowDuplicate).(ast.Expr)
		}

		// Token: Rbrack
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: Key
		r.applyDecorations(out, "Key", n.Decs.Key)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
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
			out.Files[k] = r.restoreNode(v, allowDuplicate).(*ast.File)
		}

		return out
	case *dst.ParenExpr:
		out := &ast.ParenExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Lparen
		out.Lparen = r.cursor
		r.cursor += token.Pos(len(token.LPAREN.String()))

		// Decoration: Lparen
		r.applyDecorations(out, "Lparen", n.Decs.Lparen)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Token: Rparen
		out.Rparen = r.cursor
		r.cursor += token.Pos(len(token.RPAREN.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.RangeStmt:
		out := &ast.RangeStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: For
		out.For = r.cursor
		r.cursor += token.Pos(len(token.FOR.String()))

		// Decoration: For
		r.applyDecorations(out, "For", n.Decs.For)

		// Node: Key
		if n.Key != nil {
			out.Key = r.restoreNode(n.Key, allowDuplicate).(ast.Expr)
		}

		// Token: Comma
		if n.Value != nil {
			r.cursor += token.Pos(len(token.COMMA.String()))
		}

		// Decoration: Key
		r.applyDecorations(out, "Key", n.Decs.Key)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value, allowDuplicate).(ast.Expr)
		}

		// Decoration: Value
		r.applyDecorations(out, "Value", n.Decs.Value)

		// Token: Tok
		if n.Tok != token.ILLEGAL {
			out.Tok = n.Tok
			out.TokPos = r.cursor
			r.cursor += token.Pos(len(n.Tok.String()))
		}

		// Token: Range
		r.cursor += token.Pos(len(token.RANGE.String()))

		// Decoration: Range
		r.applyDecorations(out, "Range", n.Decs.Range)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ReturnStmt:
		out := &ast.ReturnStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Return
		out.Return = r.cursor
		r.cursor += token.Pos(len(token.RETURN.String()))

		// Decoration: Return
		r.applyDecorations(out, "Return", n.Decs.Return)

		// List: Results
		for _, v := range n.Results {
			out.Results = append(out.Results, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SelectStmt:
		out := &ast.SelectStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Select
		out.Select = r.cursor
		r.cursor += token.Pos(len(token.SELECT.String()))

		// Decoration: Select
		r.applyDecorations(out, "Select", n.Decs.Select)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SelectorExpr:
		out := &ast.SelectorExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Token: Period
		r.cursor += token.Pos(len(token.PERIOD.String()))

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Node: Sel
		if n.Sel != nil {
			out.Sel = r.restoreNode(n.Sel, allowDuplicate).(*ast.Ident)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SendStmt:
		out := &ast.SendStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Chan
		if n.Chan != nil {
			out.Chan = r.restoreNode(n.Chan, allowDuplicate).(ast.Expr)
		}

		// Decoration: Chan
		r.applyDecorations(out, "Chan", n.Decs.Chan)

		// Token: Arrow
		out.Arrow = r.cursor
		r.cursor += token.Pos(len(token.ARROW.String()))

		// Decoration: Arrow
		r.applyDecorations(out, "Arrow", n.Decs.Arrow)

		// Node: Value
		if n.Value != nil {
			out.Value = r.restoreNode(n.Value, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.SliceExpr:
		out := &ast.SliceExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Token: Lbrack
		out.Lbrack = r.cursor
		r.cursor += token.Pos(len(token.LBRACK.String()))

		// Decoration: Lbrack
		r.applyDecorations(out, "Lbrack", n.Decs.Lbrack)

		// Node: Low
		if n.Low != nil {
			out.Low = r.restoreNode(n.Low, allowDuplicate).(ast.Expr)
		}

		// Token: Colon1
		r.cursor += token.Pos(len(token.COLON.String()))

		// Decoration: Low
		r.applyDecorations(out, "Low", n.Decs.Low)

		// Node: High
		if n.High != nil {
			out.High = r.restoreNode(n.High, allowDuplicate).(ast.Expr)
		}

		// Token: Colon2
		if n.Slice3 {
			r.cursor += token.Pos(len(token.COLON.String()))
		}

		// Decoration: High
		r.applyDecorations(out, "High", n.Decs.High)

		// Node: Max
		if n.Max != nil {
			out.Max = r.restoreNode(n.Max, allowDuplicate).(ast.Expr)
		}

		// Decoration: Max
		r.applyDecorations(out, "Max", n.Decs.Max)

		// Token: Rbrack
		out.Rbrack = r.cursor
		r.cursor += token.Pos(len(token.RBRACK.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Slice3
		out.Slice3 = n.Slice3
		r.applySpace(n.Decs.After)

		return out
	case *dst.StarExpr:
		out := &ast.StarExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Star
		out.Star = r.cursor
		r.cursor += token.Pos(len(token.MUL.String()))

		// Decoration: Star
		r.applyDecorations(out, "Star", n.Decs.Star)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.StructType:
		out := &ast.StructType{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Struct
		out.Struct = r.cursor
		r.cursor += token.Pos(len(token.STRUCT.String()))

		// Decoration: Struct
		r.applyDecorations(out, "Struct", n.Decs.Struct)

		// Node: Fields
		if n.Fields != nil {
			out.Fields = r.restoreNode(n.Fields, allowDuplicate).(*ast.FieldList)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)

		// Value: Incomplete
		out.Incomplete = n.Incomplete
		r.applySpace(n.Decs.After)

		return out
	case *dst.SwitchStmt:
		out := &ast.SwitchStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Switch
		out.Switch = r.cursor
		r.cursor += token.Pos(len(token.SWITCH.String()))

		// Decoration: Switch
		r.applyDecorations(out, "Switch", n.Decs.Switch)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init, allowDuplicate).(ast.Stmt)
		}

		// Decoration: Init
		r.applyDecorations(out, "Init", n.Decs.Init)

		// Node: Tag
		if n.Tag != nil {
			out.Tag = r.restoreNode(n.Tag, allowDuplicate).(ast.Expr)
		}

		// Decoration: Tag
		r.applyDecorations(out, "Tag", n.Decs.Tag)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.TypeAssertExpr:
		out := &ast.TypeAssertExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Token: Period
		r.cursor += token.Pos(len(token.PERIOD.String()))

		// Decoration: X
		r.applyDecorations(out, "X", n.Decs.X)

		// Token: Lparen
		out.Lparen = r.cursor
		r.cursor += token.Pos(len(token.LPAREN.String()))

		// Decoration: Lparen
		r.applyDecorations(out, "Lparen", n.Decs.Lparen)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type, allowDuplicate).(ast.Expr)
		}

		// Token: TypeToken
		if n.Type == nil {
			r.cursor += token.Pos(len(token.TYPE.String()))
		}

		// Decoration: Type
		r.applyDecorations(out, "Type", n.Decs.Type)

		// Token: Rparen
		out.Rparen = r.cursor
		r.cursor += token.Pos(len(token.RPAREN.String()))

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.TypeSpec:
		out := &ast.TypeSpec{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Node: Name
		if n.Name != nil {
			out.Name = r.restoreNode(n.Name, allowDuplicate).(*ast.Ident)
		}

		// Token: Assign
		if n.Assign {
			out.Assign = r.cursor
			r.cursor += token.Pos(len(token.ASSIGN.String()))
		}

		// Decoration: Name
		r.applyDecorations(out, "Name", n.Decs.Name)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.TypeSwitchStmt:
		out := &ast.TypeSwitchStmt{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Switch
		out.Switch = r.cursor
		r.cursor += token.Pos(len(token.SWITCH.String()))

		// Decoration: Switch
		r.applyDecorations(out, "Switch", n.Decs.Switch)

		// Node: Init
		if n.Init != nil {
			out.Init = r.restoreNode(n.Init, allowDuplicate).(ast.Stmt)
		}

		// Decoration: Init
		r.applyDecorations(out, "Init", n.Decs.Init)

		// Node: Assign
		if n.Assign != nil {
			out.Assign = r.restoreNode(n.Assign, allowDuplicate).(ast.Stmt)
		}

		// Decoration: Assign
		r.applyDecorations(out, "Assign", n.Decs.Assign)

		// Node: Body
		if n.Body != nil {
			out.Body = r.restoreNode(n.Body, allowDuplicate).(*ast.BlockStmt)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.UnaryExpr:
		out := &ast.UnaryExpr{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// Token: Op
		out.Op = n.Op
		out.OpPos = r.cursor
		r.cursor += token.Pos(len(n.Op.String()))

		// Decoration: Op
		r.applyDecorations(out, "Op", n.Decs.Op)

		// Node: X
		if n.X != nil {
			out.X = r.restoreNode(n.X, allowDuplicate).(ast.Expr)
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	case *dst.ValueSpec:
		out := &ast.ValueSpec{}
		r.Nodes[n] = out
		r.applySpace(n.Decs.Space)

		// Decoration: Start
		r.applyDecorations(out, "Start", n.Decs.Start)

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, r.restoreNode(v, allowDuplicate).(*ast.Ident))
		}

		// Decoration: Names
		r.applyDecorations(out, "Names", n.Decs.Names)

		// Node: Type
		if n.Type != nil {
			out.Type = r.restoreNode(n.Type, allowDuplicate).(ast.Expr)
		}

		// Token: Assign
		if n.Values != nil {
			r.cursor += token.Pos(len(token.ASSIGN.String()))
		}

		// Decoration: Assign
		r.applyDecorations(out, "Assign", n.Decs.Assign)

		// List: Values
		for _, v := range n.Values {
			out.Values = append(out.Values, r.restoreNode(v, allowDuplicate).(ast.Expr))
		}

		// Decoration: End
		r.applyDecorations(out, "End", n.Decs.End)
		r.applySpace(n.Decs.After)

		return out
	default:
		panic(fmt.Sprintf("%T", n))
	}
	return nil
}
