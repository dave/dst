package decorator

import (
	"github.com/dave/dst"
	"go/ast"
)

func (d *Decorator) DecorateNode(n ast.Node) dst.Node {
	if dn, ok := d.nodes[n]; ok {
		return dn
	}
	switch n := n.(type) {
	case *ast.ArrayType:
		out := &dst.ArrayType{}

		// Token: Lbrack

		// Node: Len
		if n.Len != nil {
			out.Len = d.DecorateNode(n.Len).(dst.Expr)
		}

		// Token: Rbrack

		// Node: Elt
		if n.Elt != nil {
			out.Elt = d.DecorateNode(n.Elt).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.AssignStmt:
		out := &dst.AssignStmt{}

		// List: Lhs
		for _, v := range n.Lhs {
			out.Lhs = append(out.Lhs, d.DecorateNode(v).(dst.Expr))
		}

		// Token: Tok
		out.Tok = n.Tok

		// List: Rhs
		for _, v := range n.Rhs {
			out.Rhs = append(out.Rhs, d.DecorateNode(v).(dst.Expr))
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BadDecl:
		out := &dst.BadDecl{}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BadExpr:
		out := &dst.BadExpr{}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BadStmt:
		out := &dst.BadStmt{}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BasicLit:
		out := &dst.BasicLit{}

		// String: Value
		out.Value = n.Value

		// Value: Kind
		out.Kind = n.Kind

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BinaryExpr:
		out := &dst.BinaryExpr{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Op
		out.Op = n.Op

		// Node: Y
		if n.Y != nil {
			out.Y = d.DecorateNode(n.Y).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BlockStmt:
		out := &dst.BlockStmt{}

		// Token: Lbrace

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, d.DecorateNode(v).(dst.Stmt))
		}

		// Token: Rbrace

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.BranchStmt:
		out := &dst.BranchStmt{}

		// Token: Tok
		out.Tok = n.Tok

		// Node: Label
		if n.Label != nil {
			out.Label = d.DecorateNode(n.Label).(*dst.Ident)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.CallExpr:
		out := &dst.CallExpr{}

		// Node: Fun
		if n.Fun != nil {
			out.Fun = d.DecorateNode(n.Fun).(dst.Expr)
		}

		// Token: Lparen

		// List: Args
		for _, v := range n.Args {
			out.Args = append(out.Args, d.DecorateNode(v).(dst.Expr))
		}

		// Token: Ellipsis
		out.Ellipsis = n.Ellipsis.IsValid()

		// Token: Rparen

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.CaseClause:
		out := &dst.CaseClause{}

		// Token: Case

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, d.DecorateNode(v).(dst.Expr))
		}

		// Token: Colon

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, d.DecorateNode(v).(dst.Stmt))
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ChanType:
		out := &dst.ChanType{}

		// Token: Begin

		// Token: Chan

		// Token: Arrow

		// Node: Value
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}

		// Value: Dir
		out.Dir = dst.ChanDir(n.Dir)

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.CommClause:
		out := &dst.CommClause{}

		// Token: Case

		// Node: Comm
		if n.Comm != nil {
			out.Comm = d.DecorateNode(n.Comm).(dst.Stmt)
		}

		// Token: Colon

		// List: Body
		for _, v := range n.Body {
			out.Body = append(out.Body, d.DecorateNode(v).(dst.Stmt))
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.CompositeLit:
		out := &dst.CompositeLit{}

		// Node: Type
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}

		// Token: Lbrace

		// List: Elts
		for _, v := range n.Elts {
			out.Elts = append(out.Elts, d.DecorateNode(v).(dst.Expr))
		}

		// Token: Rbrace

		// Value: Incomplete
		out.Incomplete = n.Incomplete

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.DeclStmt:
		out := &dst.DeclStmt{}

		// Node: Decl
		if n.Decl != nil {
			out.Decl = d.DecorateNode(n.Decl).(dst.Decl)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.DeferStmt:
		out := &dst.DeferStmt{}

		// Token: Defer

		// Node: Call
		if n.Call != nil {
			out.Call = d.DecorateNode(n.Call).(*dst.CallExpr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.Ellipsis:
		out := &dst.Ellipsis{}

		// Token: Ellipsis

		// Node: Elt
		if n.Elt != nil {
			out.Elt = d.DecorateNode(n.Elt).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.EmptyStmt:
		out := &dst.EmptyStmt{}

		// Token: Semicolon

		// Value: Implicit
		out.Implicit = n.Implicit

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ExprStmt:
		out := &dst.ExprStmt{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.Field:
		out := &dst.Field{}

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, d.DecorateNode(v).(*dst.Ident))
		}

		// Node: Type
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}

		// Node: Tag
		if n.Tag != nil {
			out.Tag = d.DecorateNode(n.Tag).(*dst.BasicLit)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.FieldList:
		out := &dst.FieldList{}

		// Token: Opening
		out.Opening = n.Opening.IsValid()

		// List: List
		for _, v := range n.List {
			out.List = append(out.List, d.DecorateNode(v).(*dst.Field))
		}

		// Token: Closing
		out.Closing = n.Closing.IsValid()

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.File:
		out := &dst.File{}

		// Token: Package

		// Node: Name
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}

		// List: Decls
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, d.DecorateNode(v).(dst.Decl))
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ForStmt:
		out := &dst.ForStmt{}

		// Token: For

		// Node: Init
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}

		// Token: InitSemicolon

		// Node: Cond
		if n.Cond != nil {
			out.Cond = d.DecorateNode(n.Cond).(dst.Expr)
		}

		// Token: CondSemicolon

		// Node: Post
		if n.Post != nil {
			out.Post = d.DecorateNode(n.Post).(dst.Stmt)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.FuncDecl:
		out := &dst.FuncDecl{}

		// Init: Type
		out.Type = &dst.FuncType{}

		// Token: Func
		out.Type.Func = true

		// Node: Recv
		if n.Recv != nil {
			out.Recv = d.DecorateNode(n.Recv).(*dst.FieldList)
		}

		// Node: Name
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}

		// Node: Params
		if n.Type.Params != nil {
			out.Type.Params = d.DecorateNode(n.Type.Params).(*dst.FieldList)
		}

		// Node: Results
		if n.Type.Results != nil {
			out.Type.Results = d.DecorateNode(n.Type.Results).(*dst.FieldList)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.FuncLit:
		out := &dst.FuncLit{}

		// Node: Type
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(*dst.FuncType)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.FuncType:
		out := &dst.FuncType{}

		// Token: Func
		out.Func = n.Func.IsValid()

		// Node: Params
		if n.Params != nil {
			out.Params = d.DecorateNode(n.Params).(*dst.FieldList)
		}

		// Node: Results
		if n.Results != nil {
			out.Results = d.DecorateNode(n.Results).(*dst.FieldList)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.GenDecl:
		out := &dst.GenDecl{}

		// Token: Tok
		out.Tok = n.Tok

		// Token: Lparen
		out.Lparen = n.Lparen.IsValid()

		// List: Specs
		for _, v := range n.Specs {
			out.Specs = append(out.Specs, d.DecorateNode(v).(dst.Spec))
		}

		// Token: Rparen
		out.Rparen = n.Rparen.IsValid()

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.GoStmt:
		out := &dst.GoStmt{}

		// Token: Go

		// Node: Call
		if n.Call != nil {
			out.Call = d.DecorateNode(n.Call).(*dst.CallExpr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.Ident:
		out := &dst.Ident{}

		// String: Name
		out.Name = n.Name

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.IfStmt:
		out := &dst.IfStmt{}

		// Token: If

		// Node: Init
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}

		// Node: Cond
		if n.Cond != nil {
			out.Cond = d.DecorateNode(n.Cond).(dst.Expr)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		// Token: ElseTok

		// Node: Else
		if n.Else != nil {
			out.Else = d.DecorateNode(n.Else).(dst.Stmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ImportSpec:
		out := &dst.ImportSpec{}

		// Node: Name
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}

		// Node: Path
		if n.Path != nil {
			out.Path = d.DecorateNode(n.Path).(*dst.BasicLit)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.IncDecStmt:
		out := &dst.IncDecStmt{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Tok
		out.Tok = n.Tok

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.IndexExpr:
		out := &dst.IndexExpr{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Lbrack

		// Node: Index
		if n.Index != nil {
			out.Index = d.DecorateNode(n.Index).(dst.Expr)
		}

		// Token: Rbrack

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.InterfaceType:
		out := &dst.InterfaceType{}

		// Token: Interface

		// Node: Methods
		if n.Methods != nil {
			out.Methods = d.DecorateNode(n.Methods).(*dst.FieldList)
		}

		// Value: Incomplete
		out.Incomplete = n.Incomplete

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.KeyValueExpr:
		out := &dst.KeyValueExpr{}

		// Node: Key
		if n.Key != nil {
			out.Key = d.DecorateNode(n.Key).(dst.Expr)
		}

		// Token: Colon

		// Node: Value
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.LabeledStmt:
		out := &dst.LabeledStmt{}

		// Node: Label
		if n.Label != nil {
			out.Label = d.DecorateNode(n.Label).(*dst.Ident)
		}

		// Token: Colon

		// Node: Stmt
		if n.Stmt != nil {
			out.Stmt = d.DecorateNode(n.Stmt).(dst.Stmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.MapType:
		out := &dst.MapType{}

		// Token: Map

		// Token: Lbrack

		// Node: Key
		if n.Key != nil {
			out.Key = d.DecorateNode(n.Key).(dst.Expr)
		}

		// Token: Rbrack

		// Node: Value
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ParenExpr:
		out := &dst.ParenExpr{}

		// Token: Lparen

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Rparen

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.RangeStmt:
		out := &dst.RangeStmt{}

		// Token: For

		// Node: Key
		if n.Key != nil {
			out.Key = d.DecorateNode(n.Key).(dst.Expr)
		}

		// Token: Comma

		// Node: Value
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}

		// Token: Tok
		out.Tok = n.Tok

		// Token: Range

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ReturnStmt:
		out := &dst.ReturnStmt{}

		// Token: Return

		// List: Results
		for _, v := range n.Results {
			out.Results = append(out.Results, d.DecorateNode(v).(dst.Expr))
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.SelectStmt:
		out := &dst.SelectStmt{}

		// Token: Select

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.SelectorExpr:
		out := &dst.SelectorExpr{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Period

		// Node: Sel
		if n.Sel != nil {
			out.Sel = d.DecorateNode(n.Sel).(*dst.Ident)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.SendStmt:
		out := &dst.SendStmt{}

		// Node: Chan
		if n.Chan != nil {
			out.Chan = d.DecorateNode(n.Chan).(dst.Expr)
		}

		// Token: Arrow

		// Node: Value
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.SliceExpr:
		out := &dst.SliceExpr{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Lbrack

		// Node: Low
		if n.Low != nil {
			out.Low = d.DecorateNode(n.Low).(dst.Expr)
		}

		// Token: Colon1

		// Node: High
		if n.High != nil {
			out.High = d.DecorateNode(n.High).(dst.Expr)
		}

		// Token: Colon2

		// Node: Max
		if n.Max != nil {
			out.Max = d.DecorateNode(n.Max).(dst.Expr)
		}

		// Token: Rbrack

		// Value: Slice3
		out.Slice3 = n.Slice3

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.StarExpr:
		out := &dst.StarExpr{}

		// Token: Star

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.StructType:
		out := &dst.StructType{}

		// Token: Struct

		// Node: Fields
		if n.Fields != nil {
			out.Fields = d.DecorateNode(n.Fields).(*dst.FieldList)
		}

		// Value: Incomplete
		out.Incomplete = n.Incomplete

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.SwitchStmt:
		out := &dst.SwitchStmt{}

		// Token: Switch

		// Node: Init
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}

		// Node: Tag
		if n.Tag != nil {
			out.Tag = d.DecorateNode(n.Tag).(dst.Expr)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.TypeAssertExpr:
		out := &dst.TypeAssertExpr{}

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		// Token: Period

		// Token: Lparen

		// Node: Type
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}

		// Token: TypeToken

		// Token: Rparen

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.TypeSpec:
		out := &dst.TypeSpec{}

		// Node: Name
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}

		// Token: Assign
		out.Assign = n.Assign.IsValid()

		// Node: Type
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.TypeSwitchStmt:
		out := &dst.TypeSwitchStmt{}

		// Token: Switch

		// Node: Init
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}

		// Node: Assign
		if n.Assign != nil {
			out.Assign = d.DecorateNode(n.Assign).(dst.Stmt)
		}

		// Node: Body
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.UnaryExpr:
		out := &dst.UnaryExpr{}

		// Token: Op
		out.Op = n.Op

		// Node: X
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	case *ast.ValueSpec:
		out := &dst.ValueSpec{}

		// List: Names
		for _, v := range n.Names {
			out.Names = append(out.Names, d.DecorateNode(v).(*dst.Ident))
		}

		// Node: Type
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}

		// Token: Assign

		// List: Values
		for _, v := range n.Values {
			out.Values = append(out.Values, d.DecorateNode(v).(dst.Expr))
		}

		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}

		d.nodes[n] = out
		return out
	}
	return nil
}
