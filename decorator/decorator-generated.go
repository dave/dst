package decorator

import (
	"github.com/dave/dst"
	"go/ast"
)

func (d *Decorator) decorateNode(n ast.Node) dst.Node {
	if dn, ok := d.Dst.Nodes[n]; ok {
		return dn
	}
	switch n := n.(type) {
	case *ast.ArrayType:
		out := &dst.ArrayType{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Lbrack

		// Node: Len
		if n.Len != nil {
			out.Len = d.decorateNode(n.Len).(dst.Expr)
		}

		// Token: Rbrack

		// Node: Elt
		if n.Elt != nil {
			out.Elt = d.decorateNode(n.Elt).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Lbrack"]; ok {
				out.Decs.Lbrack = decs
			}
			if decs, ok := nd["Len"]; ok {
				out.Decs.Len = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.AssignStmt:
		out := &dst.AssignStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// List: Lhs
		for _, v := range n.Lhs {
			out.Lhs = append(out.Lhs, d.decorateNode(v).(dst.Expr))
		}

		// Token: Tok
		out.Tok = n.Tok

		// List: Rhs
		for _, v := range n.Rhs {
			out.Rhs = append(out.Rhs, d.decorateNode(v).(dst.Expr))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Tok"]; ok {
				out.Decs.Tok = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BadDecl:
		out := &dst.BadDecl{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BadExpr:
		out := &dst.BadExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BadStmt:
		out := &dst.BadStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BasicLit:
		out := &dst.BasicLit{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// String: Value
		out.Value = n.Value

		// Value: Kind
		out.Kind = n.Kind

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BinaryExpr:
		out := &dst.BinaryExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Op
		out.Op = n.Op

		// Node: Y
		if n.Y != nil {
			out.Y = d.decorateNode(n.Y).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["Op"]; ok {
				out.Decs.Op = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BlockStmt:
		out := &dst.BlockStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Lbrace

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, d.decorateNode(v).(dst.Stmt))
		}

		// Token: Rbrace

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Lbrace"]; ok {
				out.Decs.Lbrace = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.BranchStmt:
		out := &dst.BranchStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Tok
		out.Tok = n.Tok

		// Node: Label
		if n.Label != nil {
			out.Label = d.decorateNode(n.Label).(*dst.Ident)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Tok"]; ok {
				out.Decs.Tok = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.CallExpr:
		out := &dst.CallExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Fun
		if n.Fun != nil {
			out.Fun = d.decorateNode(n.Fun).(dst.Expr)
		}

		// Token: Lparen

		// List: Args
		for _, v := range n.Args {
			out.Args = append(out.Args, d.decorateNode(v).(dst.Expr))
		}

		// Token: Ellipsis
		out.Ellipsis = n.Ellipsis.IsValid()

		// Token: Rparen

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Fun"]; ok {
				out.Decs.Fun = decs
			}
			if decs, ok := nd["Lparen"]; ok {
				out.Decs.Lparen = decs
			}
			if decs, ok := nd["Ellipsis"]; ok {
				out.Decs.Ellipsis = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.CaseClause:
		out := &dst.CaseClause{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Case

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, d.decorateNode(v).(dst.Expr))
		}

		// Token: Colon

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, d.decorateNode(v).(dst.Stmt))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Case"]; ok {
				out.Decs.Case = decs
			}
			if decs, ok := nd["Colon"]; ok {
				out.Decs.Colon = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.ChanType:
		out := &dst.ChanType{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Begin

		// Token: Chan

		// Token: Arrow

		// Node: Value
		if n.Value != nil {
			out.Value = d.decorateNode(n.Value).(dst.Expr)
		}

		// Value: Dir
		out.Dir = dst.ChanDir(n.Dir)

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Begin"]; ok {
				out.Decs.Begin = decs
			}
			if decs, ok := nd["Arrow"]; ok {
				out.Decs.Arrow = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.CommClause:
		out := &dst.CommClause{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Case

		// Node: Comm
		if n.Comm != nil {
			out.Comm = d.decorateNode(n.Comm).(dst.Stmt)
		}

		// Token: Colon

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, d.decorateNode(v).(dst.Stmt))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Case"]; ok {
				out.Decs.Case = decs
			}
			if decs, ok := nd["Comm"]; ok {
				out.Decs.Comm = decs
			}
			if decs, ok := nd["Colon"]; ok {
				out.Decs.Colon = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.CompositeLit:
		out := &dst.CompositeLit{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Type
		if n.Type != nil {
			out.Type = d.decorateNode(n.Type).(dst.Expr)
		}

		// Token: Lbrace

		// List: Elts
		for _, v := range n.Elts {
			out.Elts = append(out.Elts, d.decorateNode(v).(dst.Expr))
		}

		// Token: Rbrace

		// Value: Incomplete
		out.Incomplete = n.Incomplete

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Type"]; ok {
				out.Decs.Type = decs
			}
			if decs, ok := nd["Lbrace"]; ok {
				out.Decs.Lbrace = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.DeclStmt:
		out := &dst.DeclStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Decl
		if n.Decl != nil {
			out.Decl = d.decorateNode(n.Decl).(dst.Decl)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.DeferStmt:
		out := &dst.DeferStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Defer

		// Node: Call
		if n.Call != nil {
			out.Call = d.decorateNode(n.Call).(*dst.CallExpr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Defer"]; ok {
				out.Decs.Defer = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.Ellipsis:
		out := &dst.Ellipsis{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Ellipsis

		// Node: Elt
		if n.Elt != nil {
			out.Elt = d.decorateNode(n.Elt).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Ellipsis"]; ok {
				out.Decs.Ellipsis = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.EmptyStmt:
		out := &dst.EmptyStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Semicolon

		// Value: Implicit
		out.Implicit = n.Implicit

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.ExprStmt:
		out := &dst.ExprStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.Field:
		out := &dst.Field{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, d.decorateNode(v).(*dst.Ident))
		}

		// Node: Type
		if n.Type != nil {
			out.Type = d.decorateNode(n.Type).(dst.Expr)
		}

		// Node: Tag
		if n.Tag != nil {
			out.Tag = d.decorateNode(n.Tag).(*dst.BasicLit)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Type"]; ok {
				out.Decs.Type = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.FieldList:
		out := &dst.FieldList{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Opening
		out.Opening = n.Opening.IsValid()

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, d.decorateNode(v).(*dst.Field))
		}

		// Token: Closing
		out.Closing = n.Closing.IsValid()

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Opening"]; ok {
				out.Decs.Opening = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.File:
		out := &dst.File{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Package

		// Node: Name
		if n.Name != nil {
			out.Name = d.decorateNode(n.Name).(*dst.Ident)
		}

		// List: Decls
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, d.decorateNode(v).(dst.Decl))
		}

		// Scope: Scope
		out.Scope = d.decorateScope(n.Scope)

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Package"]; ok {
				out.Decs.Package = decs
			}
			if decs, ok := nd["Name"]; ok {
				out.Decs.Name = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.ForStmt:
		out := &dst.ForStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: For

		// Node: Init
		if n.Init != nil {
			out.Init = d.decorateNode(n.Init).(dst.Stmt)
		}

		// Token: InitSemicolon

		// Node: Cond
		if n.Cond != nil {
			out.Cond = d.decorateNode(n.Cond).(dst.Expr)
		}

		// Token: CondSemicolon

		// Node: Post
		if n.Post != nil {
			out.Post = d.decorateNode(n.Post).(dst.Stmt)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["For"]; ok {
				out.Decs.For = decs
			}
			if decs, ok := nd["Init"]; ok {
				out.Decs.Init = decs
			}
			if decs, ok := nd["Cond"]; ok {
				out.Decs.Cond = decs
			}
			if decs, ok := nd["Post"]; ok {
				out.Decs.Post = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.FuncDecl:
		out := &dst.FuncDecl{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Init: Type
		out.Type = &dst.FuncType{}

		// Token: Func
		out.Type.Func = true

		// Node: Recv
		if n.Recv != nil {
			out.Recv = d.decorateNode(n.Recv).(*dst.FieldList)
		}

		// Node: Name
		if n.Name != nil {
			out.Name = d.decorateNode(n.Name).(*dst.Ident)
		}

		// Node: Params
		if n.Type.Params != nil {
			out.Type.Params = d.decorateNode(n.Type.Params).(*dst.FieldList)
		}

		// Node: Results
		if n.Type.Results != nil {
			out.Type.Results = d.decorateNode(n.Type.Results).(*dst.FieldList)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Func"]; ok {
				out.Decs.Func = decs
			}
			if decs, ok := nd["Recv"]; ok {
				out.Decs.Recv = decs
			}
			if decs, ok := nd["Name"]; ok {
				out.Decs.Name = decs
			}
			if decs, ok := nd["Params"]; ok {
				out.Decs.Params = decs
			}
			if decs, ok := nd["Results"]; ok {
				out.Decs.Results = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.FuncLit:
		out := &dst.FuncLit{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Type
		if n.Type != nil {
			out.Type = d.decorateNode(n.Type).(*dst.FuncType)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Type"]; ok {
				out.Decs.Type = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.FuncType:
		out := &dst.FuncType{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Func
		out.Func = n.Func.IsValid()

		// Node: Params
		if n.Params != nil {
			out.Params = d.decorateNode(n.Params).(*dst.FieldList)
		}

		// Node: Results
		if n.Results != nil {
			out.Results = d.decorateNode(n.Results).(*dst.FieldList)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Func"]; ok {
				out.Decs.Func = decs
			}
			if decs, ok := nd["Params"]; ok {
				out.Decs.Params = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.GenDecl:
		out := &dst.GenDecl{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Tok
		out.Tok = n.Tok

		// Token: Lparen
		out.Lparen = n.Lparen.IsValid()

		// List: Specs
		for _, v := range n.Specs {
			out.Specs = append(out.Specs, d.decorateNode(v).(dst.Spec))
		}

		// Token: Rparen
		out.Rparen = n.Rparen.IsValid()

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Tok"]; ok {
				out.Decs.Tok = decs
			}
			if decs, ok := nd["Lparen"]; ok {
				out.Decs.Lparen = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.GoStmt:
		out := &dst.GoStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Go

		// Node: Call
		if n.Call != nil {
			out.Call = d.decorateNode(n.Call).(*dst.CallExpr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Go"]; ok {
				out.Decs.Go = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.Ident:
		out := &dst.Ident{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// String: Name
		out.Name = n.Name

		// Object: Obj
		out.Obj = d.decorateObject(n.Obj)

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.IfStmt:
		out := &dst.IfStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: If

		// Node: Init
		if n.Init != nil {
			out.Init = d.decorateNode(n.Init).(dst.Stmt)
		}

		// Node: Cond
		if n.Cond != nil {
			out.Cond = d.decorateNode(n.Cond).(dst.Expr)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		// Token: ElseTok

		// Node: Else
		if n.Else != nil {
			out.Else = d.decorateNode(n.Else).(dst.Stmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["If"]; ok {
				out.Decs.If = decs
			}
			if decs, ok := nd["Init"]; ok {
				out.Decs.Init = decs
			}
			if decs, ok := nd["Cond"]; ok {
				out.Decs.Cond = decs
			}
			if decs, ok := nd["Else"]; ok {
				out.Decs.Else = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.ImportSpec:
		out := &dst.ImportSpec{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Name
		if n.Name != nil {
			out.Name = d.decorateNode(n.Name).(*dst.Ident)
		}

		// Node: Path
		if n.Path != nil {
			out.Path = d.decorateNode(n.Path).(*dst.BasicLit)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Name"]; ok {
				out.Decs.Name = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.IncDecStmt:
		out := &dst.IncDecStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Tok
		out.Tok = n.Tok

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.IndexExpr:
		out := &dst.IndexExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Lbrack

		// Node: Index
		if n.Index != nil {
			out.Index = d.decorateNode(n.Index).(dst.Expr)
		}

		// Token: Rbrack

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["Lbrack"]; ok {
				out.Decs.Lbrack = decs
			}
			if decs, ok := nd["Index"]; ok {
				out.Decs.Index = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.InterfaceType:
		out := &dst.InterfaceType{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Interface

		// Node: Methods
		if n.Methods != nil {
			out.Methods = d.decorateNode(n.Methods).(*dst.FieldList)
		}

		// Value: Incomplete
		out.Incomplete = n.Incomplete

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Interface"]; ok {
				out.Decs.Interface = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.KeyValueExpr:
		out := &dst.KeyValueExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Key
		if n.Key != nil {
			out.Key = d.decorateNode(n.Key).(dst.Expr)
		}

		// Token: Colon

		// Node: Value
		if n.Value != nil {
			out.Value = d.decorateNode(n.Value).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Key"]; ok {
				out.Decs.Key = decs
			}
			if decs, ok := nd["Colon"]; ok {
				out.Decs.Colon = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.LabeledStmt:
		out := &dst.LabeledStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Label
		if n.Label != nil {
			out.Label = d.decorateNode(n.Label).(*dst.Ident)
		}

		// Token: Colon

		// Node: Stmt
		if n.Stmt != nil {
			out.Stmt = d.decorateNode(n.Stmt).(dst.Stmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Label"]; ok {
				out.Decs.Label = decs
			}
			if decs, ok := nd["Colon"]; ok {
				out.Decs.Colon = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.MapType:
		out := &dst.MapType{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Map

		// Token: Lbrack

		// Node: Key
		if n.Key != nil {
			out.Key = d.decorateNode(n.Key).(dst.Expr)
		}

		// Token: Rbrack

		// Node: Value
		if n.Value != nil {
			out.Value = d.decorateNode(n.Value).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Map"]; ok {
				out.Decs.Map = decs
			}
			if decs, ok := nd["Key"]; ok {
				out.Decs.Key = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.Package:
		out := &dst.Package{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		// Value: Name
		out.Name = n.Name

		// Scope: Scope
		out.Scope = d.decorateScope(n.Scope)

		// Map: Imports
		out.Imports = map[string]*dst.Object{}
		for k, v := range n.Imports {
			out.Imports[k] = d.decorateObject(v)
		}

		// Map: Files
		out.Files = map[string]*dst.File{}
		for k, v := range n.Files {
			out.Files[k] = d.decorateNode(v).(*dst.File)
		}

		return out
	case *ast.ParenExpr:
		out := &dst.ParenExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Lparen

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Rparen

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Lparen"]; ok {
				out.Decs.Lparen = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.RangeStmt:
		out := &dst.RangeStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: For

		// Node: Key
		if n.Key != nil {
			out.Key = d.decorateNode(n.Key).(dst.Expr)
		}

		// Token: Comma

		// Node: Value
		if n.Value != nil {
			out.Value = d.decorateNode(n.Value).(dst.Expr)
		}

		// Token: Tok
		out.Tok = n.Tok

		// Token: Range

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["For"]; ok {
				out.Decs.For = decs
			}
			if decs, ok := nd["Key"]; ok {
				out.Decs.Key = decs
			}
			if decs, ok := nd["Value"]; ok {
				out.Decs.Value = decs
			}
			if decs, ok := nd["Range"]; ok {
				out.Decs.Range = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.ReturnStmt:
		out := &dst.ReturnStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Return

		// List: Results
		for _, v := range n.Results {
			out.Results = append(out.Results, d.decorateNode(v).(dst.Expr))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Return"]; ok {
				out.Decs.Return = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.SelectStmt:
		out := &dst.SelectStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Select

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Select"]; ok {
				out.Decs.Select = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.SelectorExpr:
		out := &dst.SelectorExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Period

		// Node: Sel
		if n.Sel != nil {
			out.Sel = d.decorateNode(n.Sel).(*dst.Ident)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.SendStmt:
		out := &dst.SendStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Chan
		if n.Chan != nil {
			out.Chan = d.decorateNode(n.Chan).(dst.Expr)
		}

		// Token: Arrow

		// Node: Value
		if n.Value != nil {
			out.Value = d.decorateNode(n.Value).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Chan"]; ok {
				out.Decs.Chan = decs
			}
			if decs, ok := nd["Arrow"]; ok {
				out.Decs.Arrow = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.SliceExpr:
		out := &dst.SliceExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Lbrack

		// Node: Low
		if n.Low != nil {
			out.Low = d.decorateNode(n.Low).(dst.Expr)
		}

		// Token: Colon1

		// Node: High
		if n.High != nil {
			out.High = d.decorateNode(n.High).(dst.Expr)
		}

		// Token: Colon2

		// Node: Max
		if n.Max != nil {
			out.Max = d.decorateNode(n.Max).(dst.Expr)
		}

		// Token: Rbrack

		// Value: Slice3
		out.Slice3 = n.Slice3

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["Lbrack"]; ok {
				out.Decs.Lbrack = decs
			}
			if decs, ok := nd["Low"]; ok {
				out.Decs.Low = decs
			}
			if decs, ok := nd["High"]; ok {
				out.Decs.High = decs
			}
			if decs, ok := nd["Max"]; ok {
				out.Decs.Max = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.StarExpr:
		out := &dst.StarExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Star

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Star"]; ok {
				out.Decs.Star = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.StructType:
		out := &dst.StructType{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Struct

		// Node: Fields
		if n.Fields != nil {
			out.Fields = d.decorateNode(n.Fields).(*dst.FieldList)
		}

		// Value: Incomplete
		out.Incomplete = n.Incomplete

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Struct"]; ok {
				out.Decs.Struct = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.SwitchStmt:
		out := &dst.SwitchStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Switch

		// Node: Init
		if n.Init != nil {
			out.Init = d.decorateNode(n.Init).(dst.Stmt)
		}

		// Node: Tag
		if n.Tag != nil {
			out.Tag = d.decorateNode(n.Tag).(dst.Expr)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Switch"]; ok {
				out.Decs.Switch = decs
			}
			if decs, ok := nd["Init"]; ok {
				out.Decs.Init = decs
			}
			if decs, ok := nd["Tag"]; ok {
				out.Decs.Tag = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.TypeAssertExpr:
		out := &dst.TypeAssertExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		// Token: Period

		// Token: Lparen

		// Node: Type
		if n.Type != nil {
			out.Type = d.decorateNode(n.Type).(dst.Expr)
		}

		// Token: TypeToken

		// Token: Rparen

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["X"]; ok {
				out.Decs.X = decs
			}
			if decs, ok := nd["Lparen"]; ok {
				out.Decs.Lparen = decs
			}
			if decs, ok := nd["Type"]; ok {
				out.Decs.Type = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.TypeSpec:
		out := &dst.TypeSpec{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Node: Name
		if n.Name != nil {
			out.Name = d.decorateNode(n.Name).(*dst.Ident)
		}

		// Token: Assign
		out.Assign = n.Assign.IsValid()

		// Node: Type
		if n.Type != nil {
			out.Type = d.decorateNode(n.Type).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Name"]; ok {
				out.Decs.Name = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.TypeSwitchStmt:
		out := &dst.TypeSwitchStmt{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Switch

		// Node: Init
		if n.Init != nil {
			out.Init = d.decorateNode(n.Init).(dst.Stmt)
		}

		// Node: Assign
		if n.Assign != nil {
			out.Assign = d.decorateNode(n.Assign).(dst.Stmt)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Switch"]; ok {
				out.Decs.Switch = decs
			}
			if decs, ok := nd["Init"]; ok {
				out.Decs.Init = decs
			}
			if decs, ok := nd["Assign"]; ok {
				out.Decs.Assign = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.UnaryExpr:
		out := &dst.UnaryExpr{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// Token: Op
		out.Op = n.Op

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Op"]; ok {
				out.Decs.Op = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	case *ast.ValueSpec:
		out := &dst.ValueSpec{}
		d.Dst.Nodes[n] = out
		d.Ast.Nodes[out] = n

		out.Decs.Space = d.space[n]
		out.Decs.After = d.after[n]

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, d.decorateNode(v).(*dst.Ident))
		}

		// Node: Type
		if n.Type != nil {
			out.Type = d.decorateNode(n.Type).(dst.Expr)
		}

		// Token: Assign

		// List: Values
		for _, v := range n.Values {
			out.Values = append(out.Values, d.decorateNode(v).(dst.Expr))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["Assign"]; ok {
				out.Decs.Assign = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		return out
	}
	return nil
}
