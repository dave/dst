package decorator

import (
	"github.com/dave/dst"
	"go/ast"
)

func (d *Decorator) decorateNode(n ast.Node) dst.Node {
	if dn, ok := d.Nodes[n]; ok {
		return dn
	}
	switch n := n.(type) {
	case *ast.ArrayType:
		out := &dst.ArrayType{}

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
			if decs, ok := nd["AfterLbrack"]; ok {
				out.Decs.AfterLbrack = decs
			}
			if decs, ok := nd["AfterLen"]; ok {
				out.Decs.AfterLen = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.AssignStmt:
		out := &dst.AssignStmt{}

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
			if decs, ok := nd["AfterLhs"]; ok {
				out.Decs.AfterLhs = decs
			}
			if decs, ok := nd["AfterTok"]; ok {
				out.Decs.AfterTok = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.BadDecl:
		out := &dst.BadDecl{}

		d.Nodes[n] = out
		return out
	case *ast.BadExpr:
		out := &dst.BadExpr{}

		d.Nodes[n] = out
		return out
	case *ast.BadStmt:
		out := &dst.BadStmt{}

		d.Nodes[n] = out
		return out
	case *ast.BasicLit:
		out := &dst.BasicLit{}

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

		d.Nodes[n] = out
		return out
	case *ast.BinaryExpr:
		out := &dst.BinaryExpr{}

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
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["AfterOp"]; ok {
				out.Decs.AfterOp = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.BlockStmt:
		out := &dst.BlockStmt{}

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
			if decs, ok := nd["AfterLbrace"]; ok {
				out.Decs.AfterLbrace = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.BranchStmt:
		out := &dst.BranchStmt{}

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
			if decs, ok := nd["AfterTok"]; ok {
				out.Decs.AfterTok = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.CallExpr:
		out := &dst.CallExpr{}

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
			if decs, ok := nd["AfterFun"]; ok {
				out.Decs.AfterFun = decs
			}
			if decs, ok := nd["AfterLparen"]; ok {
				out.Decs.AfterLparen = decs
			}
			if decs, ok := nd["AfterArgs"]; ok {
				out.Decs.AfterArgs = decs
			}
			if decs, ok := nd["AfterEllipsis"]; ok {
				out.Decs.AfterEllipsis = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.CaseClause:
		out := &dst.CaseClause{}

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
			if decs, ok := nd["AfterCase"]; ok {
				out.Decs.AfterCase = decs
			}
			if decs, ok := nd["AfterList"]; ok {
				out.Decs.AfterList = decs
			}
			if decs, ok := nd["AfterColon"]; ok {
				out.Decs.AfterColon = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.ChanType:
		out := &dst.ChanType{}

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
			if decs, ok := nd["AfterBegin"]; ok {
				out.Decs.AfterBegin = decs
			}
			if decs, ok := nd["AfterArrow"]; ok {
				out.Decs.AfterArrow = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.CommClause:
		out := &dst.CommClause{}

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
			if decs, ok := nd["AfterCase"]; ok {
				out.Decs.AfterCase = decs
			}
			if decs, ok := nd["AfterComm"]; ok {
				out.Decs.AfterComm = decs
			}
			if decs, ok := nd["AfterColon"]; ok {
				out.Decs.AfterColon = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.CompositeLit:
		out := &dst.CompositeLit{}

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
			if decs, ok := nd["AfterType"]; ok {
				out.Decs.AfterType = decs
			}
			if decs, ok := nd["AfterLbrace"]; ok {
				out.Decs.AfterLbrace = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.DeclStmt:
		out := &dst.DeclStmt{}

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

		d.Nodes[n] = out
		return out
	case *ast.DeferStmt:
		out := &dst.DeferStmt{}

		// Token: Defer

		// Node: Call
		if n.Call != nil {
			out.Call = d.decorateNode(n.Call).(*dst.CallExpr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterDefer"]; ok {
				out.Decs.AfterDefer = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.Ellipsis:
		out := &dst.Ellipsis{}

		// Token: Ellipsis

		// Node: Elt
		if n.Elt != nil {
			out.Elt = d.decorateNode(n.Elt).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterEllipsis"]; ok {
				out.Decs.AfterEllipsis = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.EmptyStmt:
		out := &dst.EmptyStmt{}

		// Token: Semicolon

		// Value: Implicit
		out.Implicit = n.Implicit

		d.Nodes[n] = out
		return out
	case *ast.ExprStmt:
		out := &dst.ExprStmt{}

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

		d.Nodes[n] = out
		return out
	case *ast.Field:
		out := &dst.Field{}

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
			if decs, ok := nd["AfterNames"]; ok {
				out.Decs.AfterNames = decs
			}
			if decs, ok := nd["AfterType"]; ok {
				out.Decs.AfterType = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.FieldList:
		out := &dst.FieldList{}

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
			if decs, ok := nd["AfterOpening"]; ok {
				out.Decs.AfterOpening = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.File:
		out := &dst.File{}

		// Token: Package

		// Node: Name
		if n.Name != nil {
			out.Name = d.decorateNode(n.Name).(*dst.Ident)
		}

		// List: Decls
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, d.decorateNode(v).(dst.Decl))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterPackage"]; ok {
				out.Decs.AfterPackage = decs
			}
			if decs, ok := nd["AfterName"]; ok {
				out.Decs.AfterName = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.ForStmt:
		out := &dst.ForStmt{}

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
			if decs, ok := nd["AfterFor"]; ok {
				out.Decs.AfterFor = decs
			}
			if decs, ok := nd["AfterInit"]; ok {
				out.Decs.AfterInit = decs
			}
			if decs, ok := nd["AfterCond"]; ok {
				out.Decs.AfterCond = decs
			}
			if decs, ok := nd["AfterPost"]; ok {
				out.Decs.AfterPost = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.FuncDecl:
		out := &dst.FuncDecl{}

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
			if decs, ok := nd["AfterFunc"]; ok {
				out.Decs.AfterFunc = decs
			}
			if decs, ok := nd["AfterRecv"]; ok {
				out.Decs.AfterRecv = decs
			}
			if decs, ok := nd["AfterName"]; ok {
				out.Decs.AfterName = decs
			}
			if decs, ok := nd["AfterParams"]; ok {
				out.Decs.AfterParams = decs
			}
			if decs, ok := nd["AfterResults"]; ok {
				out.Decs.AfterResults = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.FuncLit:
		out := &dst.FuncLit{}

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
			if decs, ok := nd["AfterType"]; ok {
				out.Decs.AfterType = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.FuncType:
		out := &dst.FuncType{}

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
			if decs, ok := nd["AfterFunc"]; ok {
				out.Decs.AfterFunc = decs
			}
			if decs, ok := nd["AfterParams"]; ok {
				out.Decs.AfterParams = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.GenDecl:
		out := &dst.GenDecl{}

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
			if decs, ok := nd["AfterTok"]; ok {
				out.Decs.AfterTok = decs
			}
			if decs, ok := nd["AfterLparen"]; ok {
				out.Decs.AfterLparen = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.GoStmt:
		out := &dst.GoStmt{}

		// Token: Go

		// Node: Call
		if n.Call != nil {
			out.Call = d.decorateNode(n.Call).(*dst.CallExpr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterGo"]; ok {
				out.Decs.AfterGo = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.Ident:
		out := &dst.Ident{}

		// String: Name
		out.Name = n.Name

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.IfStmt:
		out := &dst.IfStmt{}

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
			if decs, ok := nd["AfterIf"]; ok {
				out.Decs.AfterIf = decs
			}
			if decs, ok := nd["AfterInit"]; ok {
				out.Decs.AfterInit = decs
			}
			if decs, ok := nd["AfterCond"]; ok {
				out.Decs.AfterCond = decs
			}
			if decs, ok := nd["AfterElse"]; ok {
				out.Decs.AfterElse = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.ImportSpec:
		out := &dst.ImportSpec{}

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
			if decs, ok := nd["AfterName"]; ok {
				out.Decs.AfterName = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.IncDecStmt:
		out := &dst.IncDecStmt{}

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
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.IndexExpr:
		out := &dst.IndexExpr{}

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
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["AfterLbrack"]; ok {
				out.Decs.AfterLbrack = decs
			}
			if decs, ok := nd["AfterIndex"]; ok {
				out.Decs.AfterIndex = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.InterfaceType:
		out := &dst.InterfaceType{}

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
			if decs, ok := nd["AfterInterface"]; ok {
				out.Decs.AfterInterface = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.KeyValueExpr:
		out := &dst.KeyValueExpr{}

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
			if decs, ok := nd["AfterKey"]; ok {
				out.Decs.AfterKey = decs
			}
			if decs, ok := nd["AfterColon"]; ok {
				out.Decs.AfterColon = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.LabeledStmt:
		out := &dst.LabeledStmt{}

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
			if decs, ok := nd["AfterLabel"]; ok {
				out.Decs.AfterLabel = decs
			}
			if decs, ok := nd["AfterColon"]; ok {
				out.Decs.AfterColon = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.MapType:
		out := &dst.MapType{}

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
			if decs, ok := nd["AfterMap"]; ok {
				out.Decs.AfterMap = decs
			}
			if decs, ok := nd["AfterKey"]; ok {
				out.Decs.AfterKey = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.Package:
		out := &dst.Package{}

		// Value: Name
		out.Name = n.Name

		// Map: Files
		for k, v := range n.Files {
			out.Files[k] = d.decorateNode(v).(*dst.File)
		}

		d.Nodes[n] = out
		return out
	case *ast.ParenExpr:
		out := &dst.ParenExpr{}

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
			if decs, ok := nd["AfterLparen"]; ok {
				out.Decs.AfterLparen = decs
			}
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.RangeStmt:
		out := &dst.RangeStmt{}

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
			if decs, ok := nd["AfterFor"]; ok {
				out.Decs.AfterFor = decs
			}
			if decs, ok := nd["AfterKey"]; ok {
				out.Decs.AfterKey = decs
			}
			if decs, ok := nd["AfterValue"]; ok {
				out.Decs.AfterValue = decs
			}
			if decs, ok := nd["AfterRange"]; ok {
				out.Decs.AfterRange = decs
			}
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.ReturnStmt:
		out := &dst.ReturnStmt{}

		// Token: Return

		// List: Results
		for _, v := range n.Results {
			out.Results = append(out.Results, d.decorateNode(v).(dst.Expr))
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterReturn"]; ok {
				out.Decs.AfterReturn = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.SelectStmt:
		out := &dst.SelectStmt{}

		// Token: Select

		// Node: Body
		if n.Body != nil {
			out.Body = d.decorateNode(n.Body).(*dst.BlockStmt)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterSelect"]; ok {
				out.Decs.AfterSelect = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.SelectorExpr:
		out := &dst.SelectorExpr{}

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
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.SendStmt:
		out := &dst.SendStmt{}

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
			if decs, ok := nd["AfterChan"]; ok {
				out.Decs.AfterChan = decs
			}
			if decs, ok := nd["AfterArrow"]; ok {
				out.Decs.AfterArrow = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.SliceExpr:
		out := &dst.SliceExpr{}

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
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["AfterLbrack"]; ok {
				out.Decs.AfterLbrack = decs
			}
			if decs, ok := nd["AfterLow"]; ok {
				out.Decs.AfterLow = decs
			}
			if decs, ok := nd["AfterHigh"]; ok {
				out.Decs.AfterHigh = decs
			}
			if decs, ok := nd["AfterMax"]; ok {
				out.Decs.AfterMax = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.StarExpr:
		out := &dst.StarExpr{}

		// Token: Star

		// Node: X
		if n.X != nil {
			out.X = d.decorateNode(n.X).(dst.Expr)
		}

		if nd, ok := d.decorations[n]; ok {
			if decs, ok := nd["Start"]; ok {
				out.Decs.Start = decs
			}
			if decs, ok := nd["AfterStar"]; ok {
				out.Decs.AfterStar = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.StructType:
		out := &dst.StructType{}

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
			if decs, ok := nd["AfterStruct"]; ok {
				out.Decs.AfterStruct = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.SwitchStmt:
		out := &dst.SwitchStmt{}

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
			if decs, ok := nd["AfterSwitch"]; ok {
				out.Decs.AfterSwitch = decs
			}
			if decs, ok := nd["AfterInit"]; ok {
				out.Decs.AfterInit = decs
			}
			if decs, ok := nd["AfterTag"]; ok {
				out.Decs.AfterTag = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.TypeAssertExpr:
		out := &dst.TypeAssertExpr{}

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
			if decs, ok := nd["AfterX"]; ok {
				out.Decs.AfterX = decs
			}
			if decs, ok := nd["AfterLparen"]; ok {
				out.Decs.AfterLparen = decs
			}
			if decs, ok := nd["AfterType"]; ok {
				out.Decs.AfterType = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.TypeSpec:
		out := &dst.TypeSpec{}

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
			if decs, ok := nd["AfterName"]; ok {
				out.Decs.AfterName = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.TypeSwitchStmt:
		out := &dst.TypeSwitchStmt{}

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
			if decs, ok := nd["AfterSwitch"]; ok {
				out.Decs.AfterSwitch = decs
			}
			if decs, ok := nd["AfterInit"]; ok {
				out.Decs.AfterInit = decs
			}
			if decs, ok := nd["AfterAssign"]; ok {
				out.Decs.AfterAssign = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.UnaryExpr:
		out := &dst.UnaryExpr{}

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
			if decs, ok := nd["AfterOp"]; ok {
				out.Decs.AfterOp = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	case *ast.ValueSpec:
		out := &dst.ValueSpec{}

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
			if decs, ok := nd["AfterNames"]; ok {
				out.Decs.AfterNames = decs
			}
			if decs, ok := nd["AfterAssign"]; ok {
				out.Decs.AfterAssign = decs
			}
			if decs, ok := nd["End"]; ok {
				out.Decs.End = decs
			}
		}

		d.Nodes[n] = out
		return out
	}
	return nil
}
