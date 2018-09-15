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
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ArrayType{}

		r.applyDecorations(n.Decs, "Lbrack", true)
		out.Lbrack = r.cursor
		r.cursor += 1 // Lbrack has fixed length 1
		r.applyDecorations(n.Decs, "Lbrack", false)

		r.applyDecorations(n.Decs, "Len", true)
		if n.Len != nil {
			out.Len = r.RestoreNode(n.Len).(ast.Expr)
			r.cursor += 1 // Len has suffix length 1
		}
		r.applyDecorations(n.Decs, "Len", false)

		r.applyDecorations(n.Decs, "Elt", true)
		if n.Elt != nil {
			out.Elt = r.RestoreNode(n.Elt).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Elt", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.AssignStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.AssignStmt{}

		r.applyDecorations(n.Decs, "Lhs", true)
		for _, v := range n.Lhs {
			out.Lhs = append(out.Lhs, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "Lhs", false)

		r.applyDecorations(n.Decs, "Tok", true)
		out.TokPos = r.cursor
		out.Tok = n.Tok
		r.cursor += token.Pos(len(n.Tok.String()))
		r.applyDecorations(n.Decs, "Tok", false)

		r.applyDecorations(n.Decs, "Rhs", true)
		for _, v := range n.Rhs {
			out.Rhs = append(out.Rhs, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "Rhs", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BadDecl:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BadDecl{}

		r.applyDecorations(n.Decs, "From", true)
		out.From = r.cursor
		r.cursor += token.Pos(n.Length)
		r.applyDecorations(n.Decs, "From", false)

		r.applyDecorations(n.Decs, "To", true)
		out.To = r.cursor
		r.applyDecorations(n.Decs, "To", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BadExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BadExpr{}

		r.applyDecorations(n.Decs, "From", true)
		out.From = r.cursor
		r.cursor += token.Pos(n.Length)
		r.applyDecorations(n.Decs, "From", false)

		r.applyDecorations(n.Decs, "To", true)
		out.To = r.cursor
		r.applyDecorations(n.Decs, "To", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BadStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BadStmt{}

		r.applyDecorations(n.Decs, "From", true)
		out.From = r.cursor
		r.cursor += token.Pos(n.Length)
		r.applyDecorations(n.Decs, "From", false)

		r.applyDecorations(n.Decs, "To", true)
		out.To = r.cursor
		r.applyDecorations(n.Decs, "To", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BasicLit:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BasicLit{}

		out.Kind = n.Kind

		r.applyDecorations(n.Decs, "Value", true)
		out.ValuePos = r.cursor
		out.Value = n.Value
		r.cursor += token.Pos(len(n.Value))
		r.applyDecorations(n.Decs, "Value", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BinaryExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BinaryExpr{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Op", true)
		out.OpPos = r.cursor
		out.Op = n.Op
		r.cursor += token.Pos(len(n.Op.String()))
		r.applyDecorations(n.Decs, "Op", false)

		r.applyDecorations(n.Decs, "Y", true)
		if n.Y != nil {
			out.Y = r.RestoreNode(n.Y).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Y", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BlockStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BlockStmt{}

		r.applyDecorations(n.Decs, "Lbrace", true)
		out.Lbrace = r.cursor
		r.cursor += 1 // Lbrace has fixed length 1
		r.applyDecorations(n.Decs, "Lbrace", false)

		r.applyDecorations(n.Decs, "List", true)
		for _, v := range n.List {
			out.List = append(out.List, r.RestoreNode(v).(ast.Stmt))
		}
		r.applyDecorations(n.Decs, "List", false)

		r.applyDecorations(n.Decs, "Rbrace", true)
		out.Rbrace = r.cursor
		r.cursor += 1 // Rbrace has fixed length 1
		r.applyDecorations(n.Decs, "Rbrace", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.BranchStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.BranchStmt{}

		r.applyDecorations(n.Decs, "Tok", true)
		out.TokPos = r.cursor
		out.Tok = n.Tok
		r.cursor += token.Pos(len(n.Tok.String()))
		r.applyDecorations(n.Decs, "Tok", false)

		r.applyDecorations(n.Decs, "Label", true)
		if n.Label != nil {
			out.Label = r.RestoreNode(n.Label).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Label", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.CallExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.CallExpr{}

		r.applyDecorations(n.Decs, "Fun", true)
		if n.Fun != nil {
			out.Fun = r.RestoreNode(n.Fun).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Fun", false)

		r.applyDecorations(n.Decs, "Lparen", true)
		out.Lparen = r.cursor
		r.cursor += 1 // Lparen has fixed length 1
		r.applyDecorations(n.Decs, "Lparen", false)

		r.applyDecorations(n.Decs, "Args", true)
		for _, v := range n.Args {
			out.Args = append(out.Args, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "Args", false)

		r.applyDecorations(n.Decs, "Ellipsis", true)
		if n.Ellipsis {
			out.Ellipsis = r.cursor
		}
		r.cursor += 3 // Ellipsis has fixed length 3
		r.applyDecorations(n.Decs, "Ellipsis", false)

		r.applyDecorations(n.Decs, "Rparen", true)
		out.Rparen = r.cursor
		r.cursor += 1 // Rparen has fixed length 1
		r.applyDecorations(n.Decs, "Rparen", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.CaseClause:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.CaseClause{}

		r.applyDecorations(n.Decs, "Case", true)
		out.Case = r.cursor
		r.cursor += 4 // Case has fixed length 4
		r.applyDecorations(n.Decs, "Case", false)

		r.applyDecorations(n.Decs, "List", true)
		for _, v := range n.List {
			out.List = append(out.List, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "List", false)

		r.applyDecorations(n.Decs, "Colon", true)
		out.Colon = r.cursor
		r.cursor += 1 // Colon has fixed length 1
		r.applyDecorations(n.Decs, "Colon", false)

		r.applyDecorations(n.Decs, "Body", true)
		for _, v := range n.Body {
			out.Body = append(out.Body, r.RestoreNode(v).(ast.Stmt))
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ChanType:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ChanType{}

		r.applyDecorations(n.Decs, "Begin", true)
		out.Begin = r.cursor
		r.cursor += 4 // Begin has fixed length 4
		r.applyDecorations(n.Decs, "Begin", false)

		r.applyDecorations(n.Decs, "Arrow", true)
		out.Arrow = r.cursor
		r.cursor += 2 // Arrow has fixed length 2
		r.applyDecorations(n.Decs, "Arrow", false)

		out.Dir = ast.ChanDir(n.Dir)

		r.applyDecorations(n.Decs, "Value", true)
		if n.Value != nil {
			out.Value = r.RestoreNode(n.Value).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Value", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.CommClause:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.CommClause{}

		r.applyDecorations(n.Decs, "Case", true)
		out.Case = r.cursor
		r.cursor += 4 // Case has fixed length 4
		r.applyDecorations(n.Decs, "Case", false)

		r.applyDecorations(n.Decs, "Comm", true)
		if n.Comm != nil {
			out.Comm = r.RestoreNode(n.Comm).(ast.Stmt)
		}
		r.applyDecorations(n.Decs, "Comm", false)

		r.applyDecorations(n.Decs, "Colon", true)
		out.Colon = r.cursor
		r.cursor += 1 // Colon has fixed length 1
		r.applyDecorations(n.Decs, "Colon", false)

		r.applyDecorations(n.Decs, "Body", true)
		for _, v := range n.Body {
			out.Body = append(out.Body, r.RestoreNode(v).(ast.Stmt))
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.Comment:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.Comment{}

		r.applyDecorations(n.Decs, "Text", true)
		out.Slash = r.cursor
		out.Text = n.Text
		r.cursor += token.Pos(len(n.Text))
		r.applyDecorations(n.Decs, "Text", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.CommentGroup:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.CommentGroup{}

		r.applyDecorations(n.Decs, "List", true)
		for _, v := range n.List {
			out.List = append(out.List, r.RestoreNode(v).(*ast.Comment))
		}
		r.applyDecorations(n.Decs, "List", false)
		r.applyDecorations(n.Decs, "", false)
		r.Comments = append(r.Comments, out)
		r.nodes[n] = out
		return out
	case *dst.CompositeLit:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.CompositeLit{}

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Lbrace", true)
		out.Lbrace = r.cursor
		r.cursor += 1 // Lbrace has fixed length 1
		r.applyDecorations(n.Decs, "Lbrace", false)

		r.applyDecorations(n.Decs, "Elts", true)
		for _, v := range n.Elts {
			out.Elts = append(out.Elts, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "Elts", false)

		r.applyDecorations(n.Decs, "Rbrace", true)
		out.Rbrace = r.cursor
		r.cursor += 1 // Rbrace has fixed length 1
		r.applyDecorations(n.Decs, "Rbrace", false)

		out.Incomplete = n.Incomplete
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.DeclStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.DeclStmt{}

		r.applyDecorations(n.Decs, "Decl", true)
		if n.Decl != nil {
			out.Decl = r.RestoreNode(n.Decl).(ast.Decl)
		}
		r.applyDecorations(n.Decs, "Decl", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.DeferStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.DeferStmt{}

		r.applyDecorations(n.Decs, "Defer", true)
		out.Defer = r.cursor
		r.cursor += 5 // Defer has fixed length 5
		r.applyDecorations(n.Decs, "Defer", false)

		r.applyDecorations(n.Decs, "Call", true)
		if n.Call != nil {
			out.Call = r.RestoreNode(n.Call).(*ast.CallExpr)
		}
		r.applyDecorations(n.Decs, "Call", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.Ellipsis:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.Ellipsis{}

		r.applyDecorations(n.Decs, "Ellipsis", true)
		out.Ellipsis = r.cursor
		r.cursor += 3 // Ellipsis has fixed length 3
		r.applyDecorations(n.Decs, "Ellipsis", false)

		r.applyDecorations(n.Decs, "Elt", true)
		if n.Elt != nil {
			out.Elt = r.RestoreNode(n.Elt).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Elt", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.EmptyStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.EmptyStmt{}

		r.applyDecorations(n.Decs, "Semicolon", true)
		out.Semicolon = r.cursor
		r.applyDecorations(n.Decs, "Semicolon", false)

		out.Implicit = n.Implicit
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ExprStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ExprStmt{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.Field:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.Field{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Names", true)
		for _, v := range n.Names {
			out.Names = append(out.Names, r.RestoreNode(v).(*ast.Ident))
		}
		r.applyDecorations(n.Decs, "Names", false)

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Tag", true)
		if n.Tag != nil {
			out.Tag = r.RestoreNode(n.Tag).(*ast.BasicLit)
		}
		r.applyDecorations(n.Decs, "Tag", false)

		r.applyDecorations(n.Decs, "Comment", true)
		if n.Comment != nil {
			out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Comment", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.FieldList:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.FieldList{}

		r.applyDecorations(n.Decs, "Opening", true)
		out.Opening = r.cursor
		r.cursor += 1 // Opening has fixed length 1
		r.applyDecorations(n.Decs, "Opening", false)

		r.applyDecorations(n.Decs, "List", true)
		for _, v := range n.List {
			out.List = append(out.List, r.RestoreNode(v).(*ast.Field))
		}
		r.applyDecorations(n.Decs, "List", false)

		r.applyDecorations(n.Decs, "Closing", true)
		out.Closing = r.cursor
		r.cursor += 1 // Closing has fixed length 1
		r.applyDecorations(n.Decs, "Closing", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.File:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.File{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Package", true)
		out.Package = r.cursor
		r.cursor += 7 // Package has fixed length 7
		r.applyDecorations(n.Decs, "Package", false)

		r.applyDecorations(n.Decs, "Name", true)
		if n.Name != nil {
			out.Name = r.RestoreNode(n.Name).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Name", false)

		r.applyDecorations(n.Decs, "Decls", true)
		for _, v := range n.Decls {
			out.Decls = append(out.Decls, r.RestoreNode(v).(ast.Decl))
		}
		r.applyDecorations(n.Decs, "Decls", false)

		// TODO: Scope (Scope)

		for _, v := range n.Imports {
			out.Imports = append(out.Imports, r.RestoreNode(v).(*ast.ImportSpec))
		}

		for _, v := range n.Unresolved {
			out.Unresolved = append(out.Unresolved, r.RestoreNode(v).(*ast.Ident))
		}
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ForStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ForStmt{}

		r.applyDecorations(n.Decs, "For", true)
		out.For = r.cursor
		r.cursor += 3 // For has fixed length 3
		r.applyDecorations(n.Decs, "For", false)

		r.applyDecorations(n.Decs, "Init", true)
		if n.Init != nil {
			out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			r.cursor += 1 // Init has suffix length 1
		}
		r.applyDecorations(n.Decs, "Init", false)

		r.applyDecorations(n.Decs, "Cond", true)
		if n.Cond != nil {
			out.Cond = r.RestoreNode(n.Cond).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Cond", false)

		r.applyDecorations(n.Decs, "Post", true)
		if n.Post != nil {
			r.cursor += 1 // Post has prefix length 1
			out.Post = r.RestoreNode(n.Post).(ast.Stmt)
		}
		r.applyDecorations(n.Decs, "Post", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.FuncDecl:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.FuncDecl{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Recv", true)
		if n.Recv != nil {
			out.Recv = r.RestoreNode(n.Recv).(*ast.FieldList)
		}
		r.applyDecorations(n.Decs, "Recv", false)

		r.applyDecorations(n.Decs, "Name", true)
		if n.Name != nil {
			out.Name = r.RestoreNode(n.Name).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Name", false)

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(*ast.FuncType)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.FuncLit:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.FuncLit{}

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(*ast.FuncType)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.FuncType:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.FuncType{}

		r.applyDecorations(n.Decs, "Func", true)
		if n.Func {
			out.Func = r.cursor
		}
		r.cursor += 4 // Func has fixed length 4
		r.applyDecorations(n.Decs, "Func", false)

		r.applyDecorations(n.Decs, "Params", true)
		if n.Params != nil {
			out.Params = r.RestoreNode(n.Params).(*ast.FieldList)
		}
		r.applyDecorations(n.Decs, "Params", false)

		r.applyDecorations(n.Decs, "Results", true)
		if n.Results != nil {
			out.Results = r.RestoreNode(n.Results).(*ast.FieldList)
		}
		r.applyDecorations(n.Decs, "Results", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.GenDecl:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.GenDecl{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Tok", true)
		out.TokPos = r.cursor
		out.Tok = n.Tok
		r.cursor += token.Pos(len(n.Tok.String()))
		r.applyDecorations(n.Decs, "Tok", false)

		r.applyDecorations(n.Decs, "Lparen", true)
		if n.Lparen {
			out.Lparen = r.cursor
		}
		r.cursor += 1 // Lparen has fixed length 1
		r.applyDecorations(n.Decs, "Lparen", false)

		r.applyDecorations(n.Decs, "Specs", true)
		for _, v := range n.Specs {
			out.Specs = append(out.Specs, r.RestoreNode(v).(ast.Spec))
		}
		r.applyDecorations(n.Decs, "Specs", false)

		r.applyDecorations(n.Decs, "Rparen", true)
		if n.Rparen {
			out.Rparen = r.cursor
		}
		r.cursor += 1 // Rparen has fixed length 1
		r.applyDecorations(n.Decs, "Rparen", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.GoStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.GoStmt{}

		r.applyDecorations(n.Decs, "Go", true)
		out.Go = r.cursor
		r.cursor += 2 // Go has fixed length 2
		r.applyDecorations(n.Decs, "Go", false)

		r.applyDecorations(n.Decs, "Call", true)
		if n.Call != nil {
			out.Call = r.RestoreNode(n.Call).(*ast.CallExpr)
		}
		r.applyDecorations(n.Decs, "Call", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.Ident:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.Ident{}

		r.applyDecorations(n.Decs, "Name", true)
		out.NamePos = r.cursor
		out.Name = n.Name
		r.cursor += token.Pos(len(n.Name))
		r.applyDecorations(n.Decs, "Name", false)

		// TODO: Obj (Object)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.IfStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.IfStmt{}

		r.applyDecorations(n.Decs, "If", true)
		out.If = r.cursor
		r.cursor += 2 // If has fixed length 2
		r.applyDecorations(n.Decs, "If", false)

		r.applyDecorations(n.Decs, "Init", true)
		if n.Init != nil {
			out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			r.cursor += 1 // Init has suffix length 1
		}
		r.applyDecorations(n.Decs, "Init", false)

		r.applyDecorations(n.Decs, "Cond", true)
		if n.Cond != nil {
			out.Cond = r.RestoreNode(n.Cond).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Cond", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)

		r.applyDecorations(n.Decs, "Else", true)
		if n.Else != nil {
			out.Else = r.RestoreNode(n.Else).(ast.Stmt)
			r.cursor += 4 // Else has suffix length 4
		}
		r.applyDecorations(n.Decs, "Else", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ImportSpec:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ImportSpec{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Name", true)
		if n.Name != nil {
			out.Name = r.RestoreNode(n.Name).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Name", false)

		r.applyDecorations(n.Decs, "Path", true)
		if n.Path != nil {
			out.Path = r.RestoreNode(n.Path).(*ast.BasicLit)
		}
		r.applyDecorations(n.Decs, "Path", false)

		r.applyDecorations(n.Decs, "Comment", true)
		if n.Comment != nil {
			out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Comment", false)

		out.EndPos = r.cursor
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.IncDecStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.IncDecStmt{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Tok", true)
		out.TokPos = r.cursor
		out.Tok = n.Tok
		r.cursor += token.Pos(len(n.Tok.String()))
		r.applyDecorations(n.Decs, "Tok", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.IndexExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.IndexExpr{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Lbrack", true)
		out.Lbrack = r.cursor
		r.cursor += 1 // Lbrack has fixed length 1
		r.applyDecorations(n.Decs, "Lbrack", false)

		r.applyDecorations(n.Decs, "Index", true)
		if n.Index != nil {
			out.Index = r.RestoreNode(n.Index).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Index", false)

		r.applyDecorations(n.Decs, "Rbrack", true)
		out.Rbrack = r.cursor
		r.cursor += 1 // Rbrack has fixed length 1
		r.applyDecorations(n.Decs, "Rbrack", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.InterfaceType:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.InterfaceType{}

		r.applyDecorations(n.Decs, "Interface", true)
		out.Interface = r.cursor
		r.cursor += 9 // Interface has fixed length 9
		r.applyDecorations(n.Decs, "Interface", false)

		r.applyDecorations(n.Decs, "Methods", true)
		if n.Methods != nil {
			out.Methods = r.RestoreNode(n.Methods).(*ast.FieldList)
		}
		r.applyDecorations(n.Decs, "Methods", false)

		out.Incomplete = n.Incomplete
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.KeyValueExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.KeyValueExpr{}

		r.applyDecorations(n.Decs, "Key", true)
		if n.Key != nil {
			out.Key = r.RestoreNode(n.Key).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Key", false)

		r.applyDecorations(n.Decs, "Colon", true)
		out.Colon = r.cursor
		r.cursor += 1 // Colon has fixed length 1
		r.applyDecorations(n.Decs, "Colon", false)

		r.applyDecorations(n.Decs, "Value", true)
		if n.Value != nil {
			out.Value = r.RestoreNode(n.Value).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Value", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.LabeledStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.LabeledStmt{}

		r.applyDecorations(n.Decs, "Label", true)
		if n.Label != nil {
			out.Label = r.RestoreNode(n.Label).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Label", false)

		r.applyDecorations(n.Decs, "Colon", true)
		out.Colon = r.cursor
		r.cursor += 1 // Colon has fixed length 1
		r.applyDecorations(n.Decs, "Colon", false)

		r.applyDecorations(n.Decs, "Stmt", true)
		if n.Stmt != nil {
			out.Stmt = r.RestoreNode(n.Stmt).(ast.Stmt)
		}
		r.applyDecorations(n.Decs, "Stmt", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.MapType:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.MapType{}

		r.applyDecorations(n.Decs, "Map", true)
		out.Map = r.cursor
		r.cursor += 3 // Map has fixed length 3
		r.applyDecorations(n.Decs, "Map", false)

		r.applyDecorations(n.Decs, "Key", true)
		if n.Key != nil {
			r.cursor += 1 // Key has prefix length 1
			out.Key = r.RestoreNode(n.Key).(ast.Expr)
			r.cursor += 1 // Key has suffix length 1
		}
		r.applyDecorations(n.Decs, "Key", false)

		r.applyDecorations(n.Decs, "Value", true)
		if n.Value != nil {
			out.Value = r.RestoreNode(n.Value).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Value", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ParenExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ParenExpr{}

		r.applyDecorations(n.Decs, "Lparen", true)
		out.Lparen = r.cursor
		r.cursor += 1 // Lparen has fixed length 1
		r.applyDecorations(n.Decs, "Lparen", false)

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Rparen", true)
		out.Rparen = r.cursor
		r.cursor += 1 // Rparen has fixed length 1
		r.applyDecorations(n.Decs, "Rparen", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.RangeStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.RangeStmt{}

		r.applyDecorations(n.Decs, "For", true)
		out.For = r.cursor
		r.cursor += 3 // For has fixed length 3
		r.applyDecorations(n.Decs, "For", false)

		r.applyDecorations(n.Decs, "Key", true)
		if n.Key != nil {
			out.Key = r.RestoreNode(n.Key).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Key", false)

		r.applyDecorations(n.Decs, "Value", true)
		if n.Value != nil {
			out.Value = r.RestoreNode(n.Value).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Value", false)

		r.applyDecorations(n.Decs, "Tok", true)
		out.TokPos = r.cursor
		out.Tok = n.Tok
		r.cursor += token.Pos(len(n.Tok.String()))
		r.applyDecorations(n.Decs, "Tok", false)

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			r.cursor += 5 // X has prefix length 5
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ReturnStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ReturnStmt{}

		r.applyDecorations(n.Decs, "Return", true)
		out.Return = r.cursor
		r.cursor += 6 // Return has fixed length 6
		r.applyDecorations(n.Decs, "Return", false)

		r.applyDecorations(n.Decs, "Results", true)
		for _, v := range n.Results {
			out.Results = append(out.Results, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "Results", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.SelectStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.SelectStmt{}

		r.applyDecorations(n.Decs, "Select", true)
		out.Select = r.cursor
		r.cursor += 6 // Select has fixed length 6
		r.applyDecorations(n.Decs, "Select", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.SelectorExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.SelectorExpr{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
			r.cursor += 1 // X has suffix length 1
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Sel", true)
		if n.Sel != nil {
			out.Sel = r.RestoreNode(n.Sel).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Sel", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.SendStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.SendStmt{}

		r.applyDecorations(n.Decs, "Chan", true)
		if n.Chan != nil {
			out.Chan = r.RestoreNode(n.Chan).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Chan", false)

		r.applyDecorations(n.Decs, "Arrow", true)
		out.Arrow = r.cursor
		r.cursor += 2 // Arrow has fixed length 2
		r.applyDecorations(n.Decs, "Arrow", false)

		r.applyDecorations(n.Decs, "Value", true)
		if n.Value != nil {
			out.Value = r.RestoreNode(n.Value).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Value", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.SliceExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.SliceExpr{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Lbrack", true)
		out.Lbrack = r.cursor
		r.cursor += 1 // Lbrack has fixed length 1
		r.applyDecorations(n.Decs, "Lbrack", false)

		r.applyDecorations(n.Decs, "Low", true)
		if n.Low != nil {
			out.Low = r.RestoreNode(n.Low).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Low", false)

		r.applyDecorations(n.Decs, "High", true)
		if n.High != nil {
			r.cursor += 1 // High has prefix length 1
			out.High = r.RestoreNode(n.High).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "High", false)

		r.applyDecorations(n.Decs, "Max", true)
		if n.Max != nil {
			r.cursor += 1 // Max has prefix length 1
			out.Max = r.RestoreNode(n.Max).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Max", false)

		out.Slice3 = n.Slice3

		r.applyDecorations(n.Decs, "Rbrack", true)
		out.Rbrack = r.cursor
		r.cursor += 1 // Rbrack has fixed length 1
		r.applyDecorations(n.Decs, "Rbrack", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.StarExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.StarExpr{}

		r.applyDecorations(n.Decs, "Star", true)
		out.Star = r.cursor
		r.cursor += 1 // Star has fixed length 1
		r.applyDecorations(n.Decs, "Star", false)

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.StructType:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.StructType{}

		r.applyDecorations(n.Decs, "Struct", true)
		out.Struct = r.cursor
		r.cursor += 6 // Struct has fixed length 6
		r.applyDecorations(n.Decs, "Struct", false)

		r.applyDecorations(n.Decs, "Fields", true)
		if n.Fields != nil {
			out.Fields = r.RestoreNode(n.Fields).(*ast.FieldList)
		}
		r.applyDecorations(n.Decs, "Fields", false)

		out.Incomplete = n.Incomplete
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.SwitchStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.SwitchStmt{}

		r.applyDecorations(n.Decs, "Switch", true)
		out.Switch = r.cursor
		r.cursor += 6 // Switch has fixed length 6
		r.applyDecorations(n.Decs, "Switch", false)

		r.applyDecorations(n.Decs, "Init", true)
		if n.Init != nil {
			out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			r.cursor += 1 // Init has suffix length 1
		}
		r.applyDecorations(n.Decs, "Init", false)

		r.applyDecorations(n.Decs, "Tag", true)
		if n.Tag != nil {
			out.Tag = r.RestoreNode(n.Tag).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Tag", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.TypeAssertExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.TypeAssertExpr{}

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)

		r.applyDecorations(n.Decs, "Lparen", true)
		out.Lparen = r.cursor
		r.cursor += 1 // Lparen has fixed length 1
		r.applyDecorations(n.Decs, "Lparen", false)

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Rparen", true)
		out.Rparen = r.cursor
		r.cursor += 1 // Rparen has fixed length 1
		r.applyDecorations(n.Decs, "Rparen", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.TypeSpec:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.TypeSpec{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Name", true)
		if n.Name != nil {
			out.Name = r.RestoreNode(n.Name).(*ast.Ident)
		}
		r.applyDecorations(n.Decs, "Name", false)

		r.applyDecorations(n.Decs, "Assign", true)
		if n.Assign {
			out.Assign = r.cursor
		}
		r.cursor += 1 // Assign has fixed length 1
		r.applyDecorations(n.Decs, "Assign", false)

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Comment", true)
		if n.Comment != nil {
			out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Comment", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.TypeSwitchStmt:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.TypeSwitchStmt{}

		r.applyDecorations(n.Decs, "Switch", true)
		out.Switch = r.cursor
		r.cursor += 6 // Switch has fixed length 6
		r.applyDecorations(n.Decs, "Switch", false)

		r.applyDecorations(n.Decs, "Init", true)
		if n.Init != nil {
			out.Init = r.RestoreNode(n.Init).(ast.Stmt)
			r.cursor += 1 // Init has suffix length 1
		}
		r.applyDecorations(n.Decs, "Init", false)

		r.applyDecorations(n.Decs, "Assign", true)
		if n.Assign != nil {
			out.Assign = r.RestoreNode(n.Assign).(ast.Stmt)
		}
		r.applyDecorations(n.Decs, "Assign", false)

		r.applyDecorations(n.Decs, "Body", true)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.applyDecorations(n.Decs, "Body", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.UnaryExpr:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.UnaryExpr{}

		r.applyDecorations(n.Decs, "Op", true)
		out.OpPos = r.cursor
		out.Op = n.Op
		r.cursor += token.Pos(len(n.Op.String()))
		r.applyDecorations(n.Decs, "Op", false)

		r.applyDecorations(n.Decs, "X", true)
		if n.X != nil {
			out.X = r.RestoreNode(n.X).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "X", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	case *dst.ValueSpec:
		r.applyDecorations(n.Decs, "", true)
		out := &ast.ValueSpec{}

		r.applyDecorations(n.Decs, "Doc", true)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Doc", false)

		r.applyDecorations(n.Decs, "Names", true)
		for _, v := range n.Names {
			out.Names = append(out.Names, r.RestoreNode(v).(*ast.Ident))
		}
		r.applyDecorations(n.Decs, "Names", false)

		r.applyDecorations(n.Decs, "Type", true)
		if n.Type != nil {
			out.Type = r.RestoreNode(n.Type).(ast.Expr)
		}
		r.applyDecorations(n.Decs, "Type", false)

		r.applyDecorations(n.Decs, "Values", true)
		for _, v := range n.Values {
			out.Values = append(out.Values, r.RestoreNode(v).(ast.Expr))
		}
		r.applyDecorations(n.Decs, "Values", false)

		r.applyDecorations(n.Decs, "Comment", true)
		if n.Comment != nil {
			out.Comment = r.RestoreNode(n.Comment).(*ast.CommentGroup)
		}
		r.applyDecorations(n.Decs, "Comment", false)
		r.applyDecorations(n.Decs, "", false)
		r.nodes[n] = out
		return out
	default:
		panic(fmt.Sprintf("%T", n))
	}
	return nil
}
