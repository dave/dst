package decorator

import (
	"github.com/dave/dst"
	"go/ast"
	"go/token"
)

func (d *Decorator) DecorateNode(n ast.Node) dst.Node {
	if dn, ok := d.nodes[n]; ok {
		return dn
	}
	switch n := n.(type) {
	case *ast.ArrayType:
		out := &dst.ArrayType{}
		if n.Len != nil {
			out.Len = d.DecorateNode(n.Len).(dst.Expr)
		}
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
		for _, v := range n.Lhs {
			out.Lhs = append(out.Lhs, d.DecorateNode(v).(dst.Expr))
		}
		out.Tok = n.Tok
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
		out.Length = int(n.End() - n.Pos())
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.BadExpr:
		out := &dst.BadExpr{}
		out.Length = int(n.End() - n.Pos())
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.BadStmt:
		out := &dst.BadStmt{}
		out.Length = int(n.End() - n.Pos())
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.BasicLit:
		out := &dst.BasicLit{}
		out.Value = n.Value
		out.Kind = n.Kind
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.BinaryExpr:
		out := &dst.BinaryExpr{}
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
		out.Op = n.Op
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
		for _, v := range n.List {
			out.List = append(out.List, d.DecorateNode(v).(dst.Stmt))
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.BranchStmt:
		out := &dst.BranchStmt{}
		out.Tok = n.Tok
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
		if n.Fun != nil {
			out.Fun = d.DecorateNode(n.Fun).(dst.Expr)
		}
		for _, v := range n.Args {
			out.Args = append(out.Args, d.DecorateNode(v).(dst.Expr))
		}
		if n.Ellipsis != token.NoPos {
			out.Ellipsis = true
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.CaseClause:
		out := &dst.CaseClause{}
		for _, v := range n.List {
			out.List = append(out.List, d.DecorateNode(v).(dst.Expr))
		}
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
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}
		out.Dir = dst.ChanDir(n.Dir)
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.CommClause:
		out := &dst.CommClause{}
		if n.Comm != nil {
			out.Comm = d.DecorateNode(n.Comm).(dst.Stmt)
		}
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
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}
		for _, v := range n.Elts {
			out.Elts = append(out.Elts, d.DecorateNode(v).(dst.Expr))
		}
		out.Incomplete = n.Incomplete
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.DeclStmt:
		out := &dst.DeclStmt{}
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
		out.Implicit = n.Implicit
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.ExprStmt:
		out := &dst.ExprStmt{}
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
		for _, v := range n.Names {
			out.Names = append(out.Names, d.DecorateNode(v).(*dst.Ident))
		}
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}
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
		if n.Opening != token.NoPos {
			out.Opening = true
		}
		for _, v := range n.List {
			out.List = append(out.List, d.DecorateNode(v).(*dst.Field))
		}
		if n.Closing != token.NoPos {
			out.Closing = true
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.File:
		out := &dst.File{}
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, d.DecorateNode(v).(dst.Decl))
		}
		for _, v := range n.Imports {
			out.Imports = append(out.Imports, d.DecorateNode(v).(*dst.ImportSpec))
		}
		for _, v := range n.Unresolved {
			out.Unresolved = append(out.Unresolved, d.DecorateNode(v).(*dst.Ident))
		}
		// TODO: Scope (Scope)
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.ForStmt:
		out := &dst.ForStmt{}
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}
		if n.Cond != nil {
			out.Cond = d.DecorateNode(n.Cond).(dst.Expr)
		}
		if n.Post != nil {
			out.Post = d.DecorateNode(n.Post).(dst.Stmt)
		}
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
		if n.Recv != nil {
			out.Recv = d.DecorateNode(n.Recv).(*dst.FieldList)
		}
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(*dst.FuncType)
		}
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
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(*dst.FuncType)
		}
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
		if n.Func != token.NoPos {
			out.Func = true
		}
		if n.Params != nil {
			out.Params = d.DecorateNode(n.Params).(*dst.FieldList)
		}
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
		out.Tok = n.Tok
		if n.Lparen != token.NoPos {
			out.Lparen = true
		}
		for _, v := range n.Specs {
			out.Specs = append(out.Specs, d.DecorateNode(v).(dst.Spec))
		}
		if n.Rparen != token.NoPos {
			out.Rparen = true
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.GoStmt:
		out := &dst.GoStmt{}
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
		out.Name = n.Name
		// TODO: Obj (Object)
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.IfStmt:
		out := &dst.IfStmt{}
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}
		if n.Cond != nil {
			out.Cond = d.DecorateNode(n.Cond).(dst.Expr)
		}
		if n.Body != nil {
			out.Body = d.DecorateNode(n.Body).(*dst.BlockStmt)
		}
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
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}
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
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
		out.Tok = n.Tok
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.IndexExpr:
		out := &dst.IndexExpr{}
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
		if n.Index != nil {
			out.Index = d.DecorateNode(n.Index).(dst.Expr)
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.InterfaceType:
		out := &dst.InterfaceType{}
		if n.Methods != nil {
			out.Methods = d.DecorateNode(n.Methods).(*dst.FieldList)
		}
		out.Incomplete = n.Incomplete
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.KeyValueExpr:
		out := &dst.KeyValueExpr{}
		if n.Key != nil {
			out.Key = d.DecorateNode(n.Key).(dst.Expr)
		}
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
		if n.Label != nil {
			out.Label = d.DecorateNode(n.Label).(*dst.Ident)
		}
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
		if n.Key != nil {
			out.Key = d.DecorateNode(n.Key).(dst.Expr)
		}
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
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.RangeStmt:
		out := &dst.RangeStmt{}
		if n.Key != nil {
			out.Key = d.DecorateNode(n.Key).(dst.Expr)
		}
		if n.Value != nil {
			out.Value = d.DecorateNode(n.Value).(dst.Expr)
		}
		out.Tok = n.Tok
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
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
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
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
		if n.Chan != nil {
			out.Chan = d.DecorateNode(n.Chan).(dst.Expr)
		}
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
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
		if n.Low != nil {
			out.Low = d.DecorateNode(n.Low).(dst.Expr)
		}
		if n.High != nil {
			out.High = d.DecorateNode(n.High).(dst.Expr)
		}
		if n.Max != nil {
			out.Max = d.DecorateNode(n.Max).(dst.Expr)
		}
		out.Slice3 = n.Slice3
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.StarExpr:
		out := &dst.StarExpr{}
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
		if n.Fields != nil {
			out.Fields = d.DecorateNode(n.Fields).(*dst.FieldList)
		}
		out.Incomplete = n.Incomplete
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.SwitchStmt:
		out := &dst.SwitchStmt{}
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}
		if n.Tag != nil {
			out.Tag = d.DecorateNode(n.Tag).(dst.Expr)
		}
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
		if n.X != nil {
			out.X = d.DecorateNode(n.X).(dst.Expr)
		}
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}
		if decs, ok := d.decorations[n]; ok {
			out.Decs = decs
		}
		d.nodes[n] = out
		return out
	case *ast.TypeSpec:
		out := &dst.TypeSpec{}
		if n.Name != nil {
			out.Name = d.DecorateNode(n.Name).(*dst.Ident)
		}
		if n.Assign != token.NoPos {
			out.Assign = true
		}
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
		if n.Init != nil {
			out.Init = d.DecorateNode(n.Init).(dst.Stmt)
		}
		if n.Assign != nil {
			out.Assign = d.DecorateNode(n.Assign).(dst.Stmt)
		}
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
		out.Op = n.Op
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
		for _, v := range n.Names {
			out.Names = append(out.Names, d.DecorateNode(v).(*dst.Ident))
		}
		if n.Type != nil {
			out.Type = d.DecorateNode(n.Type).(dst.Expr)
		}
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
