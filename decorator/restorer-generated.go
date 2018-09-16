package decorator

import (
	"fmt"
	"github.com/dave/dst"
	"go/ast"
	"go/token"
)

func (r *FileRestorer) RestoreNode(n dst.Node) ast.Node {
	if an, ok := r.nodes[n]; ok {
		return an
	}
	switch n := n.(type) {
	case *dst.ArrayType:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ArrayType{}
		{
			r.applyDecorations(n.Decs, "Lbrack", false)
			prefix, length, suffix := getLength(n, "Lbrack")
			out.Lbrack = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lbrack", true)
		}
		{
			r.applyDecorations(n.Decs, "Len", false)
			prefix, length, suffix := getLength(n, "Len")
			r.cursor += token.Pos(prefix)
			if n.Len != nil {
				out.Len = r.RestoreNode(n.Len).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Len", true)
		}
		{
			r.applyDecorations(n.Decs, "Elt", false)
			prefix, length, suffix := getLength(n, "Elt")
			r.cursor += token.Pos(prefix)
			if n.Elt != nil {
				out.Elt = r.RestoreNode(n.Elt).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Elt", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.AssignStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.AssignStmt{}
		{
			r.applyDecorations(n.Decs, "Lhs", false)
			prefix, length, suffix := getLength(n, "Lhs")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Lhs {
				out.Lhs = append(out.Lhs, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lhs", true)
		}
		{
			r.applyDecorations(n.Decs, "Tok", false)
			prefix, length, suffix := getLength(n, "Tok")
			out.TokPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Tok = n.Tok
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tok", true)
		}
		{
			r.applyDecorations(n.Decs, "Rhs", false)
			prefix, length, suffix := getLength(n, "Rhs")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Rhs {
				out.Rhs = append(out.Rhs, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rhs", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BadDecl:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BadDecl{}
		out.From = r.cursor
		r.cursor += token.Pos(n.Length)
		out.To = r.cursor
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BadExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BadExpr{}
		out.From = r.cursor
		r.cursor += token.Pos(n.Length)
		out.To = r.cursor
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BadStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BadStmt{}
		out.From = r.cursor
		r.cursor += token.Pos(n.Length)
		out.To = r.cursor
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BasicLit:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BasicLit{}
		{
			r.applyDecorations(n.Decs, "Value", false)
			prefix, length, suffix := getLength(n, "Value")
			out.ValuePos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Value = n.Value
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Value", true)
		}
		out.Kind = n.Kind
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BinaryExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BinaryExpr{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Op", false)
			prefix, length, suffix := getLength(n, "Op")
			out.OpPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Op = n.Op
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Op", true)
		}
		{
			r.applyDecorations(n.Decs, "Y", false)
			prefix, length, suffix := getLength(n, "Y")
			r.cursor += token.Pos(prefix)
			if n.Y != nil {
				out.Y = r.RestoreNode(n.Y).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Y", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BlockStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BlockStmt{}
		{
			r.applyDecorations(n.Decs, "Lbrace", false)
			prefix, length, suffix := getLength(n, "Lbrace")
			out.Lbrace = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lbrace", true)
		}
		{
			r.applyDecorations(n.Decs, "List", false)
			prefix, length, suffix := getLength(n, "List")
			r.cursor += token.Pos(prefix)
			for _, v := range n.List {
				out.List = append(out.List, r.RestoreNode(v).(ast.Stmt))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "List", true)
		}
		{
			r.applyDecorations(n.Decs, "Rbrace", false)
			prefix, length, suffix := getLength(n, "Rbrace")
			out.Rbrace = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rbrace", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.BranchStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.BranchStmt{}
		{
			r.applyDecorations(n.Decs, "Tok", false)
			prefix, length, suffix := getLength(n, "Tok")
			out.TokPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Tok = n.Tok
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tok", true)
		}
		{
			r.applyDecorations(n.Decs, "Label", false)
			prefix, length, suffix := getLength(n, "Label")
			r.cursor += token.Pos(prefix)
			if n.Label != nil {
				out.Label = r.RestoreNode(n.Label).(*ast.Ident)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Label", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.CallExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.CallExpr{}
		{
			r.applyDecorations(n.Decs, "Fun", false)
			prefix, length, suffix := getLength(n, "Fun")
			r.cursor += token.Pos(prefix)
			if n.Fun != nil {
				out.Fun = r.RestoreNode(n.Fun).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Fun", true)
		}
		{
			r.applyDecorations(n.Decs, "Lparen", false)
			prefix, length, suffix := getLength(n, "Lparen")
			out.Lparen = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lparen", true)
		}
		{
			r.applyDecorations(n.Decs, "Args", false)
			prefix, length, suffix := getLength(n, "Args")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Args {
				out.Args = append(out.Args, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Args", true)
		}
		{
			r.applyDecorations(n.Decs, "Ellipsis", false)
			prefix, length, suffix := getLength(n, "Ellipsis")
			if n.Ellipsis {
				out.Ellipsis = r.cursor
			}
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Ellipsis", true)
		}
		{
			r.applyDecorations(n.Decs, "Rparen", false)
			prefix, length, suffix := getLength(n, "Rparen")
			out.Rparen = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rparen", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.CaseClause:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.CaseClause{}
		{
			r.applyDecorations(n.Decs, "Case", false)
			prefix, length, suffix := getLength(n, "Case")
			out.Case = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Case", true)
		}
		{
			r.applyDecorations(n.Decs, "List", false)
			prefix, length, suffix := getLength(n, "List")
			r.cursor += token.Pos(prefix)
			for _, v := range n.List {
				out.List = append(out.List, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "List", true)
		}
		{
			r.applyDecorations(n.Decs, "Colon", false)
			prefix, length, suffix := getLength(n, "Colon")
			out.Colon = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Colon", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Body {
				out.Body = append(out.Body, r.RestoreNode(v).(ast.Stmt))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ChanType:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ChanType{}
		{
			r.applyDecorations(n.Decs, "Begin", false)
			prefix, length, suffix := getLength(n, "Begin")
			out.Begin = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Begin", true)
		}
		{
			r.applyDecorations(n.Decs, "Arrow", false)
			prefix, length, suffix := getLength(n, "Arrow")
			out.Arrow = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Arrow", true)
		}
		{
			r.applyDecorations(n.Decs, "Value", false)
			prefix, length, suffix := getLength(n, "Value")
			r.cursor += token.Pos(prefix)
			if n.Value != nil {
				out.Value = r.RestoreNode(n.Value).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Value", true)
		}
		out.Dir = ast.ChanDir(n.Dir)
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.CommClause:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.CommClause{}
		{
			r.applyDecorations(n.Decs, "Case", false)
			prefix, length, suffix := getLength(n, "Case")
			out.Case = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Case", true)
		}
		{
			r.applyDecorations(n.Decs, "Comm", false)
			prefix, length, suffix := getLength(n, "Comm")
			r.cursor += token.Pos(prefix)
			if n.Comm != nil {
				out.Comm = r.RestoreNode(n.Comm).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Comm", true)
		}
		{
			r.applyDecorations(n.Decs, "Colon", false)
			prefix, length, suffix := getLength(n, "Colon")
			out.Colon = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Colon", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Body {
				out.Body = append(out.Body, r.RestoreNode(v).(ast.Stmt))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.Comment:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.Comment{}
		{
			r.applyDecorations(n.Decs, "Text", false)
			prefix, length, suffix := getLength(n, "Text")
			out.Slash = r.cursor
			r.cursor += token.Pos(prefix)
			out.Text = n.Text
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Text", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.CommentGroup:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.CommentGroup{}
		{
			r.applyDecorations(n.Decs, "List", false)
			prefix, length, suffix := getLength(n, "List")
			r.cursor += token.Pos(prefix)
			for _, v := range n.List {
				out.List = append(out.List, r.RestoreNode(v).(*ast.Comment))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "List", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.Comments = append(r.Comments, out)
		r.nodes[n] = out
		return out
	case *dst.CompositeLit:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.CompositeLit{}
		{
			r.applyDecorations(n.Decs, "Type", false)
			prefix, length, suffix := getLength(n, "Type")
			r.cursor += token.Pos(prefix)
			if n.Type != nil {
				out.Type = r.RestoreNode(n.Type).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Type", true)
		}
		{
			r.applyDecorations(n.Decs, "Lbrace", false)
			prefix, length, suffix := getLength(n, "Lbrace")
			out.Lbrace = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lbrace", true)
		}
		{
			r.applyDecorations(n.Decs, "Elts", false)
			prefix, length, suffix := getLength(n, "Elts")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Elts {
				out.Elts = append(out.Elts, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Elts", true)
		}
		{
			r.applyDecorations(n.Decs, "Rbrace", false)
			prefix, length, suffix := getLength(n, "Rbrace")
			out.Rbrace = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rbrace", true)
		}
		out.Incomplete = n.Incomplete
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.DeclStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.DeclStmt{}
		{
			r.applyDecorations(n.Decs, "Decl", false)
			prefix, length, suffix := getLength(n, "Decl")
			r.cursor += token.Pos(prefix)
			if n.Decl != nil {
				out.Decl = r.RestoreNode(n.Decl).(ast.Decl)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Decl", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.DeferStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.DeferStmt{}
		{
			r.applyDecorations(n.Decs, "Defer", false)
			prefix, length, suffix := getLength(n, "Defer")
			out.Defer = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Defer", true)
		}
		{
			r.applyDecorations(n.Decs, "Call", false)
			prefix, length, suffix := getLength(n, "Call")
			r.cursor += token.Pos(prefix)
			if n.Call != nil {
				out.Call = r.RestoreNode(n.Call).(*ast.CallExpr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Call", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.Ellipsis:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.Ellipsis{}
		{
			r.applyDecorations(n.Decs, "Ellipsis", false)
			prefix, length, suffix := getLength(n, "Ellipsis")
			out.Ellipsis = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Ellipsis", true)
		}
		{
			r.applyDecorations(n.Decs, "Elt", false)
			prefix, length, suffix := getLength(n, "Elt")
			r.cursor += token.Pos(prefix)
			if n.Elt != nil {
				out.Elt = r.RestoreNode(n.Elt).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Elt", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.EmptyStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.EmptyStmt{}
		{
			r.applyDecorations(n.Decs, "Semicolon", false)
			prefix, length, suffix := getLength(n, "Semicolon")
			out.Semicolon = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Semicolon", true)
		}
		out.Implicit = n.Implicit
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ExprStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ExprStmt{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.Field:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.Field{}
		{
			r.applyDecorations(n.Decs, "Doc", false)
			prefix, length, suffix := getLength(n, "Doc")
			r.cursor += token.Pos(prefix)
			if n.Doc != nil {
				out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Doc", true)
		}
		{
			r.applyDecorations(n.Decs, "Names", false)
			prefix, length, suffix := getLength(n, "Names")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Names {
				out.Names = append(out.Names, r.RestoreNode(v).(*ast.Ident))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Names", true)
		}
		{
			r.applyDecorations(n.Decs, "Type", false)
			prefix, length, suffix := getLength(n, "Type")
			r.cursor += token.Pos(prefix)
			if n.Type != nil {
				out.Type = r.RestoreNode(n.Type).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Type", true)
		}
		{
			r.applyDecorations(n.Decs, "Tag", false)
			prefix, length, suffix := getLength(n, "Tag")
			r.cursor += token.Pos(prefix)
			if n.Tag != nil {
				out.Tag = r.RestoreNode(n.Tag).(*ast.BasicLit)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tag", true)
		}
		{
			r.applyDecorations(n.Decs, "Comment", false)
			prefix, length, suffix := getLength(n, "Comment")
			r.cursor += token.Pos(prefix)
			if n.Comment != nil {
				out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Comment", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.FieldList:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.FieldList{}
		{
			r.applyDecorations(n.Decs, "Opening", false)
			prefix, length, suffix := getLength(n, "Opening")
			out.Opening = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Opening", true)
		}
		{
			r.applyDecorations(n.Decs, "List", false)
			prefix, length, suffix := getLength(n, "List")
			r.cursor += token.Pos(prefix)
			for _, v := range n.List {
				out.List = append(out.List, r.RestoreNode(v).(*ast.Field))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "List", true)
		}
		{
			r.applyDecorations(n.Decs, "Closing", false)
			prefix, length, suffix := getLength(n, "Closing")
			out.Closing = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Closing", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.File:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.File{}
		{
			r.applyDecorations(n.Decs, "Doc", false)
			prefix, length, suffix := getLength(n, "Doc")
			r.cursor += token.Pos(prefix)
			if n.Doc != nil {
				out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Doc", true)
		}
		{
			r.applyDecorations(n.Decs, "Package", false)
			prefix, length, suffix := getLength(n, "Package")
			out.Package = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Package", true)
		}
		{
			r.applyDecorations(n.Decs, "Name", false)
			prefix, length, suffix := getLength(n, "Name")
			r.cursor += token.Pos(prefix)
			if n.Name != nil {
				out.Name = r.RestoreNode(n.Name).(*ast.Ident)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Name", true)
		}
		{
			r.applyDecorations(n.Decs, "Decls", false)
			prefix, length, suffix := getLength(n, "Decls")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Decls {
				out.Decls = append(out.Decls, r.RestoreNode(v).(ast.Decl))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Decls", true)
		}
		for _, v := range n.Imports {
			out.Imports = append(out.Imports, r.RestoreNode(v).(*ast.ImportSpec))
		}
		for _, v := range n.Unresolved {
			out.Unresolved = append(out.Unresolved, r.RestoreNode(v).(*ast.Ident))
		}
		// TODO: Scope (Scope)
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ForStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ForStmt{}
		{
			r.applyDecorations(n.Decs, "For", false)
			prefix, length, suffix := getLength(n, "For")
			out.For = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "For", true)
		}
		{
			r.applyDecorations(n.Decs, "Init", false)
			prefix, length, suffix := getLength(n, "Init")
			r.cursor += token.Pos(prefix)
			if n.Init != nil {
				out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Init", true)
		}
		{
			r.applyDecorations(n.Decs, "Cond", false)
			prefix, length, suffix := getLength(n, "Cond")
			r.cursor += token.Pos(prefix)
			if n.Cond != nil {
				out.Cond = r.RestoreNode(n.Cond).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Cond", true)
		}
		{
			r.applyDecorations(n.Decs, "Post", false)
			prefix, length, suffix := getLength(n, "Post")
			r.cursor += token.Pos(prefix)
			if n.Post != nil {
				out.Post = r.RestoreNode(n.Post).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Post", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.FuncDecl:
		return r.funcDeclOverride(n)
	case *dst.FuncLit:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.FuncLit{}
		{
			r.applyDecorations(n.Decs, "Type", false)
			prefix, length, suffix := getLength(n, "Type")
			r.cursor += token.Pos(prefix)
			if n.Type != nil {
				out.Type = r.RestoreNode(n.Type).(*ast.FuncType)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Type", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.FuncType:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.FuncType{}
		{
			r.applyDecorations(n.Decs, "Func", false)
			prefix, length, suffix := getLength(n, "Func")
			if n.Func {
				out.Func = r.cursor
			}
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Func", true)
		}
		{
			r.applyDecorations(n.Decs, "Params", false)
			prefix, length, suffix := getLength(n, "Params")
			r.cursor += token.Pos(prefix)
			if n.Params != nil {
				out.Params = r.RestoreNode(n.Params).(*ast.FieldList)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Params", true)
		}
		{
			r.applyDecorations(n.Decs, "Results", false)
			prefix, length, suffix := getLength(n, "Results")
			r.cursor += token.Pos(prefix)
			if n.Results != nil {
				out.Results = r.RestoreNode(n.Results).(*ast.FieldList)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Results", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.GenDecl:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.GenDecl{}
		{
			r.applyDecorations(n.Decs, "Doc", false)
			prefix, length, suffix := getLength(n, "Doc")
			r.cursor += token.Pos(prefix)
			if n.Doc != nil {
				out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Doc", true)
		}
		{
			r.applyDecorations(n.Decs, "Tok", false)
			prefix, length, suffix := getLength(n, "Tok")
			out.TokPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Tok = n.Tok
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tok", true)
		}
		{
			r.applyDecorations(n.Decs, "Lparen", false)
			prefix, length, suffix := getLength(n, "Lparen")
			if n.Lparen {
				out.Lparen = r.cursor
			}
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lparen", true)
		}
		{
			r.applyDecorations(n.Decs, "Specs", false)
			prefix, length, suffix := getLength(n, "Specs")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Specs {
				out.Specs = append(out.Specs, r.RestoreNode(v).(ast.Spec))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Specs", true)
		}
		{
			r.applyDecorations(n.Decs, "Rparen", false)
			prefix, length, suffix := getLength(n, "Rparen")
			if n.Rparen {
				out.Rparen = r.cursor
			}
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rparen", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.GoStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.GoStmt{}
		{
			r.applyDecorations(n.Decs, "Go", false)
			prefix, length, suffix := getLength(n, "Go")
			out.Go = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Go", true)
		}
		{
			r.applyDecorations(n.Decs, "Call", false)
			prefix, length, suffix := getLength(n, "Call")
			r.cursor += token.Pos(prefix)
			if n.Call != nil {
				out.Call = r.RestoreNode(n.Call).(*ast.CallExpr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Call", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.Ident:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.Ident{}
		{
			r.applyDecorations(n.Decs, "Name", false)
			prefix, length, suffix := getLength(n, "Name")
			out.NamePos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Name = n.Name
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Name", true)
		}
		// TODO: Obj (Object)
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.IfStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.IfStmt{}
		{
			r.applyDecorations(n.Decs, "If", false)
			prefix, length, suffix := getLength(n, "If")
			out.If = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "If", true)
		}
		{
			r.applyDecorations(n.Decs, "Init", false)
			prefix, length, suffix := getLength(n, "Init")
			r.cursor += token.Pos(prefix)
			if n.Init != nil {
				out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Init", true)
		}
		{
			r.applyDecorations(n.Decs, "Cond", false)
			prefix, length, suffix := getLength(n, "Cond")
			r.cursor += token.Pos(prefix)
			if n.Cond != nil {
				out.Cond = r.RestoreNode(n.Cond).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Cond", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		{
			r.applyDecorations(n.Decs, "Else", false)
			prefix, length, suffix := getLength(n, "Else")
			r.cursor += token.Pos(prefix)
			if n.Else != nil {
				out.Else = r.RestoreNode(n.Else).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Else", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ImportSpec:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ImportSpec{}
		{
			r.applyDecorations(n.Decs, "Doc", false)
			prefix, length, suffix := getLength(n, "Doc")
			r.cursor += token.Pos(prefix)
			if n.Doc != nil {
				out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Doc", true)
		}
		{
			r.applyDecorations(n.Decs, "Name", false)
			prefix, length, suffix := getLength(n, "Name")
			r.cursor += token.Pos(prefix)
			if n.Name != nil {
				out.Name = r.RestoreNode(n.Name).(*ast.Ident)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Name", true)
		}
		{
			r.applyDecorations(n.Decs, "Path", false)
			prefix, length, suffix := getLength(n, "Path")
			r.cursor += token.Pos(prefix)
			if n.Path != nil {
				out.Path = r.RestoreNode(n.Path).(*ast.BasicLit)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Path", true)
		}
		{
			r.applyDecorations(n.Decs, "Comment", false)
			prefix, length, suffix := getLength(n, "Comment")
			r.cursor += token.Pos(prefix)
			if n.Comment != nil {
				out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Comment", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.IncDecStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.IncDecStmt{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Tok", false)
			prefix, length, suffix := getLength(n, "Tok")
			out.TokPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Tok = n.Tok
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tok", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.IndexExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.IndexExpr{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Lbrack", false)
			prefix, length, suffix := getLength(n, "Lbrack")
			out.Lbrack = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lbrack", true)
		}
		{
			r.applyDecorations(n.Decs, "Index", false)
			prefix, length, suffix := getLength(n, "Index")
			r.cursor += token.Pos(prefix)
			if n.Index != nil {
				out.Index = r.RestoreNode(n.Index).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Index", true)
		}
		{
			r.applyDecorations(n.Decs, "Rbrack", false)
			prefix, length, suffix := getLength(n, "Rbrack")
			out.Rbrack = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rbrack", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.InterfaceType:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.InterfaceType{}
		{
			r.applyDecorations(n.Decs, "Interface", false)
			prefix, length, suffix := getLength(n, "Interface")
			out.Interface = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Interface", true)
		}
		{
			r.applyDecorations(n.Decs, "Methods", false)
			prefix, length, suffix := getLength(n, "Methods")
			r.cursor += token.Pos(prefix)
			if n.Methods != nil {
				out.Methods = r.RestoreNode(n.Methods).(*ast.FieldList)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Methods", true)
		}
		out.Incomplete = n.Incomplete
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.KeyValueExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.KeyValueExpr{}
		{
			r.applyDecorations(n.Decs, "Key", false)
			prefix, length, suffix := getLength(n, "Key")
			r.cursor += token.Pos(prefix)
			if n.Key != nil {
				out.Key = r.RestoreNode(n.Key).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Key", true)
		}
		{
			r.applyDecorations(n.Decs, "Colon", false)
			prefix, length, suffix := getLength(n, "Colon")
			out.Colon = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Colon", true)
		}
		{
			r.applyDecorations(n.Decs, "Value", false)
			prefix, length, suffix := getLength(n, "Value")
			r.cursor += token.Pos(prefix)
			if n.Value != nil {
				out.Value = r.RestoreNode(n.Value).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Value", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.LabeledStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.LabeledStmt{}
		{
			r.applyDecorations(n.Decs, "Label", false)
			prefix, length, suffix := getLength(n, "Label")
			r.cursor += token.Pos(prefix)
			if n.Label != nil {
				out.Label = r.RestoreNode(n.Label).(*ast.Ident)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Label", true)
		}
		{
			r.applyDecorations(n.Decs, "Colon", false)
			prefix, length, suffix := getLength(n, "Colon")
			out.Colon = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Colon", true)
		}
		{
			r.applyDecorations(n.Decs, "Stmt", false)
			prefix, length, suffix := getLength(n, "Stmt")
			r.cursor += token.Pos(prefix)
			if n.Stmt != nil {
				out.Stmt = r.RestoreNode(n.Stmt).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Stmt", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.MapType:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.MapType{}
		{
			r.applyDecorations(n.Decs, "Map", false)
			prefix, length, suffix := getLength(n, "Map")
			out.Map = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Map", true)
		}
		{
			r.applyDecorations(n.Decs, "Key", false)
			prefix, length, suffix := getLength(n, "Key")
			r.cursor += token.Pos(prefix)
			if n.Key != nil {
				out.Key = r.RestoreNode(n.Key).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Key", true)
		}
		{
			r.applyDecorations(n.Decs, "Value", false)
			prefix, length, suffix := getLength(n, "Value")
			r.cursor += token.Pos(prefix)
			if n.Value != nil {
				out.Value = r.RestoreNode(n.Value).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Value", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ParenExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ParenExpr{}
		{
			r.applyDecorations(n.Decs, "Lparen", false)
			prefix, length, suffix := getLength(n, "Lparen")
			out.Lparen = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lparen", true)
		}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Rparen", false)
			prefix, length, suffix := getLength(n, "Rparen")
			out.Rparen = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rparen", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.RangeStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.RangeStmt{}
		{
			r.applyDecorations(n.Decs, "For", false)
			prefix, length, suffix := getLength(n, "For")
			out.For = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "For", true)
		}
		{
			r.applyDecorations(n.Decs, "Key", false)
			prefix, length, suffix := getLength(n, "Key")
			r.cursor += token.Pos(prefix)
			if n.Key != nil {
				out.Key = r.RestoreNode(n.Key).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Key", true)
		}
		{
			r.applyDecorations(n.Decs, "Value", false)
			prefix, length, suffix := getLength(n, "Value")
			r.cursor += token.Pos(prefix)
			if n.Value != nil {
				out.Value = r.RestoreNode(n.Value).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Value", true)
		}
		{
			r.applyDecorations(n.Decs, "Tok", false)
			prefix, length, suffix := getLength(n, "Tok")
			out.TokPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Tok = n.Tok
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tok", true)
		}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ReturnStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ReturnStmt{}
		{
			r.applyDecorations(n.Decs, "Return", false)
			prefix, length, suffix := getLength(n, "Return")
			out.Return = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Return", true)
		}
		{
			r.applyDecorations(n.Decs, "Results", false)
			prefix, length, suffix := getLength(n, "Results")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Results {
				out.Results = append(out.Results, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Results", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.SelectStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.SelectStmt{}
		{
			r.applyDecorations(n.Decs, "Select", false)
			prefix, length, suffix := getLength(n, "Select")
			out.Select = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Select", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.SelectorExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.SelectorExpr{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Sel", false)
			prefix, length, suffix := getLength(n, "Sel")
			r.cursor += token.Pos(prefix)
			if n.Sel != nil {
				out.Sel = r.RestoreNode(n.Sel).(*ast.Ident)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Sel", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.SendStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.SendStmt{}
		{
			r.applyDecorations(n.Decs, "Chan", false)
			prefix, length, suffix := getLength(n, "Chan")
			r.cursor += token.Pos(prefix)
			if n.Chan != nil {
				out.Chan = r.RestoreNode(n.Chan).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Chan", true)
		}
		{
			r.applyDecorations(n.Decs, "Arrow", false)
			prefix, length, suffix := getLength(n, "Arrow")
			out.Arrow = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Arrow", true)
		}
		{
			r.applyDecorations(n.Decs, "Value", false)
			prefix, length, suffix := getLength(n, "Value")
			r.cursor += token.Pos(prefix)
			if n.Value != nil {
				out.Value = r.RestoreNode(n.Value).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Value", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.SliceExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.SliceExpr{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Lbrack", false)
			prefix, length, suffix := getLength(n, "Lbrack")
			out.Lbrack = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lbrack", true)
		}
		{
			r.applyDecorations(n.Decs, "Low", false)
			prefix, length, suffix := getLength(n, "Low")
			r.cursor += token.Pos(prefix)
			if n.Low != nil {
				out.Low = r.RestoreNode(n.Low).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Low", true)
		}
		{
			r.applyDecorations(n.Decs, "High", false)
			prefix, length, suffix := getLength(n, "High")
			r.cursor += token.Pos(prefix)
			if n.High != nil {
				out.High = r.RestoreNode(n.High).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "High", true)
		}
		{
			r.applyDecorations(n.Decs, "Max", false)
			prefix, length, suffix := getLength(n, "Max")
			r.cursor += token.Pos(prefix)
			if n.Max != nil {
				out.Max = r.RestoreNode(n.Max).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Max", true)
		}
		{
			r.applyDecorations(n.Decs, "Rbrack", false)
			prefix, length, suffix := getLength(n, "Rbrack")
			out.Rbrack = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rbrack", true)
		}
		out.Slice3 = n.Slice3
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.StarExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.StarExpr{}
		{
			r.applyDecorations(n.Decs, "Star", false)
			prefix, length, suffix := getLength(n, "Star")
			out.Star = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Star", true)
		}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.StructType:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.StructType{}
		{
			r.applyDecorations(n.Decs, "Struct", false)
			prefix, length, suffix := getLength(n, "Struct")
			out.Struct = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Struct", true)
		}
		{
			r.applyDecorations(n.Decs, "Fields", false)
			prefix, length, suffix := getLength(n, "Fields")
			r.cursor += token.Pos(prefix)
			if n.Fields != nil {
				out.Fields = r.RestoreNode(n.Fields).(*ast.FieldList)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Fields", true)
		}
		out.Incomplete = n.Incomplete
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.SwitchStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.SwitchStmt{}
		{
			r.applyDecorations(n.Decs, "Switch", false)
			prefix, length, suffix := getLength(n, "Switch")
			out.Switch = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Switch", true)
		}
		{
			r.applyDecorations(n.Decs, "Init", false)
			prefix, length, suffix := getLength(n, "Init")
			r.cursor += token.Pos(prefix)
			if n.Init != nil {
				out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Init", true)
		}
		{
			r.applyDecorations(n.Decs, "Tag", false)
			prefix, length, suffix := getLength(n, "Tag")
			r.cursor += token.Pos(prefix)
			if n.Tag != nil {
				out.Tag = r.RestoreNode(n.Tag).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Tag", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.TypeAssertExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.TypeAssertExpr{}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		{
			r.applyDecorations(n.Decs, "Lparen", false)
			prefix, length, suffix := getLength(n, "Lparen")
			out.Lparen = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Lparen", true)
		}
		{
			r.applyDecorations(n.Decs, "Type", false)
			prefix, length, suffix := getLength(n, "Type")
			r.cursor += token.Pos(prefix)
			if n.Type != nil {
				out.Type = r.RestoreNode(n.Type).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Type", true)
		}
		{
			r.applyDecorations(n.Decs, "Rparen", false)
			prefix, length, suffix := getLength(n, "Rparen")
			out.Rparen = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Rparen", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.TypeSpec:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.TypeSpec{}
		{
			r.applyDecorations(n.Decs, "Doc", false)
			prefix, length, suffix := getLength(n, "Doc")
			r.cursor += token.Pos(prefix)
			if n.Doc != nil {
				out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Doc", true)
		}
		{
			r.applyDecorations(n.Decs, "Name", false)
			prefix, length, suffix := getLength(n, "Name")
			r.cursor += token.Pos(prefix)
			if n.Name != nil {
				out.Name = r.RestoreNode(n.Name).(*ast.Ident)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Name", true)
		}
		{
			r.applyDecorations(n.Decs, "Assign", false)
			prefix, length, suffix := getLength(n, "Assign")
			if n.Assign {
				out.Assign = r.cursor
			}
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Assign", true)
		}
		{
			r.applyDecorations(n.Decs, "Type", false)
			prefix, length, suffix := getLength(n, "Type")
			r.cursor += token.Pos(prefix)
			if n.Type != nil {
				out.Type = r.RestoreNode(n.Type).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Type", true)
		}
		{
			r.applyDecorations(n.Decs, "Comment", false)
			prefix, length, suffix := getLength(n, "Comment")
			r.cursor += token.Pos(prefix)
			if n.Comment != nil {
				out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Comment", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.TypeSwitchStmt:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.TypeSwitchStmt{}
		{
			r.applyDecorations(n.Decs, "Switch", false)
			prefix, length, suffix := getLength(n, "Switch")
			out.Switch = r.cursor
			r.cursor += token.Pos(prefix)
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Switch", true)
		}
		{
			r.applyDecorations(n.Decs, "Init", false)
			prefix, length, suffix := getLength(n, "Init")
			r.cursor += token.Pos(prefix)
			if n.Init != nil {
				out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Init", true)
		}
		{
			r.applyDecorations(n.Decs, "Assign", false)
			prefix, length, suffix := getLength(n, "Assign")
			r.cursor += token.Pos(prefix)
			if n.Assign != nil {
				out.Assign = r.RestoreNode(n.Assign).(ast.Stmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Assign", true)
		}
		{
			r.applyDecorations(n.Decs, "Body", false)
			prefix, length, suffix := getLength(n, "Body")
			r.cursor += token.Pos(prefix)
			if n.Body != nil {
				out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Body", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.UnaryExpr:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.UnaryExpr{}
		{
			r.applyDecorations(n.Decs, "Op", false)
			prefix, length, suffix := getLength(n, "Op")
			out.OpPos = r.cursor
			r.cursor += token.Pos(prefix)
			out.Op = n.Op
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Op", true)
		}
		{
			r.applyDecorations(n.Decs, "X", false)
			prefix, length, suffix := getLength(n, "X")
			r.cursor += token.Pos(prefix)
			if n.X != nil {
				out.X = r.RestoreNode(n.X).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "X", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	case *dst.ValueSpec:
		r.applyDecorations(n.Decs, "", false)
		out := &ast.ValueSpec{}
		{
			r.applyDecorations(n.Decs, "Doc", false)
			prefix, length, suffix := getLength(n, "Doc")
			r.cursor += token.Pos(prefix)
			if n.Doc != nil {
				out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Doc", true)
		}
		{
			r.applyDecorations(n.Decs, "Names", false)
			prefix, length, suffix := getLength(n, "Names")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Names {
				out.Names = append(out.Names, r.RestoreNode(v).(*ast.Ident))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Names", true)
		}
		{
			r.applyDecorations(n.Decs, "Type", false)
			prefix, length, suffix := getLength(n, "Type")
			r.cursor += token.Pos(prefix)
			if n.Type != nil {
				out.Type = r.RestoreNode(n.Type).(ast.Expr)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Type", true)
		}
		{
			r.applyDecorations(n.Decs, "Values", false)
			prefix, length, suffix := getLength(n, "Values")
			r.cursor += token.Pos(prefix)
			for _, v := range n.Values {
				out.Values = append(out.Values, r.RestoreNode(v).(ast.Expr))
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Values", true)
		}
		{
			r.applyDecorations(n.Decs, "Comment", false)
			prefix, length, suffix := getLength(n, "Comment")
			r.cursor += token.Pos(prefix)
			if n.Comment != nil {
				out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
			}
			r.cursor += token.Pos(length)
			r.cursor += token.Pos(suffix)
			r.applyDecorations(n.Decs, "Comment", true)
		}
		r.applyDecorations(n.Decs, "", true)
		r.nodes[n] = out
		return out
	default:
		panic(fmt.Sprintf("%T", n))
	}
	return nil
}
