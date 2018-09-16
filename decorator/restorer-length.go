package decorator

import (
	"go/token"

	"github.com/dave/dst"
)

/*
	Modifications after generation:

	ArrayType.Len
	SliceExpr.Max
	ChanType.Arrow
	EmptyStmt.Semicolon
	TypeAssertExpr.Type
	CaseClause.Case
	CommClause.Comm
*/

func getLength(n dst.Node, fragment string) (suffix, length, prefix int) {
	switch n := n.(type) {
	case *dst.ArrayType:
		switch fragment {
		case "Lbrack":
			return 0, 1, 0
		case "Len":
			// ************
			// SPECIAL CASE
			// Len has "]" suffix even when Len == nil
			// ************
			return 0, 0, 1
		case "Elt":
			return 0, 0, 0
		}
	case *dst.AssignStmt:
		switch fragment {
		case "Lhs":
			return 0, 0, 0
		case "Tok":
			if n.Tok != token.ILLEGAL {
				return 0, len(n.Tok.String()), 0
			}
			return 0, 0, 0
		case "Rhs":
			return 0, 0, 0
		}
	case *dst.BadDecl:
		switch fragment {
		}
	case *dst.BadExpr:
		switch fragment {
		}
	case *dst.BadStmt:
		switch fragment {
		}
	case *dst.BasicLit:
		switch fragment {
		case "Value":
			if n.Value != "" {
				return 0, len(n.Value), 0
			}
			return 0, 0, 0
		}
	case *dst.BinaryExpr:
		switch fragment {
		case "X":
			return 0, 0, 0
		case "Op":
			if n.Op != token.ILLEGAL {
				return 0, len(n.Op.String()), 0
			}
			return 0, 0, 0
		case "Y":
			return 0, 0, 0
		}
	case *dst.BlockStmt:
		switch fragment {
		case "Lbrace":
			return 0, 1, 0
		case "List":
			return 0, 0, 0
		case "Rbrace":
			return 0, 1, 0
		}
	case *dst.BranchStmt:
		switch fragment {
		case "Tok":
			if n.Tok != token.ILLEGAL {
				return 0, len(n.Tok.String()), 0
			}
			return 0, 0, 0
		case "Label":
			return 0, 0, 0
		}
	case *dst.CallExpr:
		switch fragment {
		case "Fun":
			return 0, 0, 0
		case "Lparen":
			return 0, 1, 0
		case "Args":
			return 0, 0, 0
		case "Ellipsis":
			if n.Ellipsis {
				return 0, 3, 0
			}
			return 0, 0, 0
		case "Rparen":
			return 0, 1, 0
		}
	case *dst.CaseClause:
		switch fragment {
		case "Case":
			// ************
			// SPECIAL CASE
			// When List == nil, "default" is rendered. Otherwise, "case" is rendered before the node.
			// ************
			if n.List != nil {
				return 4, 0, 0 // "case"
			}
			return 0, 7, 0 // "default"
		case "List":
			return 0, 0, 0
		case "Colon":
			return 0, 1, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.ChanType:
		switch fragment {
		case "Begin":
			return 0, 4, 0
		case "Arrow":
			// ************
			// SPECIAL CASE
			// Arrow is not rendered when Dir == 0
			// ************
			if n.Dir == dst.SEND || n.Dir == dst.RECV {
				return 0, 2, 0
			}
			return 0, 0, 0
		case "Value":
			return 0, 0, 0
		}
	case *dst.CommClause:
		switch fragment {
		case "Case":
			return 0, 4, 0
		case "Comm":
			// ************
			// SPECIAL CASE
			// When Comm == nil, "default" is rendered. Otherwise, "case" is rendered before the node.
			// ************
			if n.Comm != nil {
				return 4, 0, 0 // "case"
			}
			return 0, 7, 0 // "default"
		case "Colon":
			return 0, 1, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.Comment:
		switch fragment {
		case "Text":
			if n.Text != "" {
				return 0, len(n.Text), 0
			}
			return 0, 0, 0
		}
	case *dst.CommentGroup:
		switch fragment {
		case "List":
			return 0, 0, 0
		}
	case *dst.CompositeLit:
		switch fragment {
		case "Type":
			return 0, 0, 0
		case "Lbrace":
			return 0, 1, 0
		case "Elts":
			return 0, 0, 0
		case "Rbrace":
			return 0, 1, 0
		}
	case *dst.DeclStmt:
		switch fragment {
		case "Decl":
			return 0, 0, 0
		}
	case *dst.DeferStmt:
		switch fragment {
		case "Defer":
			return 0, 5, 0
		case "Call":
			return 0, 0, 0
		}
	case *dst.Ellipsis:
		switch fragment {
		case "Ellipsis":
			return 0, 3, 0
		case "Elt":
			return 0, 0, 0
		}
	case *dst.EmptyStmt:
		switch fragment {
		case "Semicolon":
			// ************
			// SPECIAL CASE
			// Semicolon is not rendered if Implicit == true
			// ************
			if !n.Implicit {
				return 0, 1, 0
			}
			return 0, 0, 0
		}
	case *dst.ExprStmt:
		switch fragment {
		case "X":
			return 0, 0, 0
		}
	case *dst.Field:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Names":
			return 0, 0, 0
		case "Type":
			return 0, 0, 0
		case "Tag":
			return 0, 0, 0
		case "Comment":
			return 0, 0, 0
		}
	case *dst.FieldList:
		switch fragment {
		case "Opening":
			return 0, 1, 0
		case "List":
			return 0, 0, 0
		case "Closing":
			return 0, 1, 0
		}
	case *dst.File:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Package":
			return 0, 7, 0
		case "Name":
			return 0, 0, 0
		case "Decls":
			return 0, 0, 0
		}
	case *dst.ForStmt:
		switch fragment {
		case "For":
			return 0, 3, 0
		case "Init":
			if n.Init != nil {
				return 0, 0, 1
			}
			return 0, 0, 0
		case "Cond":
			return 0, 0, 0
		case "Post":
			if n.Post != nil {
				return 1, 0, 0
			}
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.FuncDecl:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Recv":
			return 0, 0, 0
		case "Name":
			return 0, 0, 0
		case "Type":
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.FuncLit:
		switch fragment {
		case "Type":
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.FuncType:
		switch fragment {
		case "Func":
			if n.Func {
				return 0, 4, 0
			}
			return 0, 0, 0
		case "Params":
			return 0, 0, 0
		case "Results":
			return 0, 0, 0
		}
	case *dst.GenDecl:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Tok":
			if n.Tok != token.ILLEGAL {
				return 0, len(n.Tok.String()), 0
			}
			return 0, 0, 0
		case "Lparen":
			if n.Lparen {
				return 0, 1, 0
			}
			return 0, 0, 0
		case "Specs":
			return 0, 0, 0
		case "Rparen":
			if n.Rparen {
				return 0, 1, 0
			}
			return 0, 0, 0
		}
	case *dst.GoStmt:
		switch fragment {
		case "Go":
			return 0, 2, 0
		case "Call":
			return 0, 0, 0
		}
	case *dst.Ident:
		switch fragment {
		case "Name":
			if n.Name != "" {
				return 0, len(n.Name), 0
			}
			return 0, 0, 0
		}
	case *dst.IfStmt:
		switch fragment {
		case "If":
			return 0, 2, 0
		case "Init":
			if n.Init != nil {
				return 0, 0, 1
			}
			return 0, 0, 0
		case "Cond":
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		case "Else":
			if n.Else != nil {
				return 0, 0, 4
			}
			return 0, 0, 0
		}
	case *dst.ImportSpec:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Name":
			return 0, 0, 0
		case "Path":
			return 0, 0, 0
		case "Comment":
			return 0, 0, 0
		}
	case *dst.IncDecStmt:
		switch fragment {
		case "X":
			return 0, 0, 0
		case "Tok":
			if n.Tok != token.ILLEGAL {
				return 0, len(n.Tok.String()), 0
			}
			return 0, 0, 0
		}
	case *dst.IndexExpr:
		switch fragment {
		case "X":
			return 0, 0, 0
		case "Lbrack":
			return 0, 1, 0
		case "Index":
			return 0, 0, 0
		case "Rbrack":
			return 0, 1, 0
		}
	case *dst.InterfaceType:
		switch fragment {
		case "Interface":
			return 0, 9, 0
		case "Methods":
			return 0, 0, 0
		}
	case *dst.KeyValueExpr:
		switch fragment {
		case "Key":
			return 0, 0, 0
		case "Colon":
			return 0, 1, 0
		case "Value":
			return 0, 0, 0
		}
	case *dst.LabeledStmt:
		switch fragment {
		case "Label":
			return 0, 0, 0
		case "Colon":
			return 0, 1, 0
		case "Stmt":
			return 0, 0, 0
		}
	case *dst.MapType:
		switch fragment {
		case "Map":
			return 0, 3, 0
		case "Key":
			if n.Key != nil {
				return 1, 0, 1
			}
			return 0, 0, 0
		case "Value":
			return 0, 0, 0
		}
	case *dst.ParenExpr:
		switch fragment {
		case "Lparen":
			return 0, 1, 0
		case "X":
			return 0, 0, 0
		case "Rparen":
			return 0, 1, 0
		}
	case *dst.RangeStmt:
		switch fragment {
		case "For":
			return 0, 3, 0
		case "Key":
			return 0, 0, 0
		case "Value":
			return 0, 0, 0
		case "Tok":
			if n.Tok != token.ILLEGAL {
				return 0, len(n.Tok.String()), 0
			}
			return 0, 0, 0
		case "X":
			if n.X != nil {
				return 5, 0, 0
			}
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.ReturnStmt:
		switch fragment {
		case "Return":
			return 0, 6, 0
		case "Results":
			return 0, 0, 0
		}
	case *dst.SelectStmt:
		switch fragment {
		case "Select":
			return 0, 6, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.SelectorExpr:
		switch fragment {
		case "X":
			if n.X != nil {
				return 0, 0, 1
			}
			return 0, 0, 0
		case "Sel":
			return 0, 0, 0
		}
	case *dst.SendStmt:
		switch fragment {
		case "Chan":
			return 0, 0, 0
		case "Arrow":
			return 0, 2, 0
		case "Value":
			return 0, 0, 0
		}
	case *dst.SliceExpr:
		switch fragment {
		case "X":
			return 0, 0, 0
		case "Lbrack":
			return 0, 1, 0
		case "Low":
			return 0, 0, 0
		case "High":
			if n.High != nil {
				return 1, 0, 0
			}
			return 0, 0, 0
		case "Max":
			// ************
			// SPECIAL CASE
			// If Slice3, we have two colons even with Max == nil
			// ************
			if n.Max != nil || n.Slice3 {
				return 1, 0, 0
			}
			return 0, 0, 0
		case "Rbrack":
			return 0, 1, 0
		}
	case *dst.StarExpr:
		switch fragment {
		case "Star":
			return 0, 1, 0
		case "X":
			return 0, 0, 0
		}
	case *dst.StructType:
		switch fragment {
		case "Struct":
			return 0, 6, 0
		case "Fields":
			return 0, 0, 0
		}
	case *dst.SwitchStmt:
		switch fragment {
		case "Switch":
			return 0, 6, 0
		case "Init":
			if n.Init != nil {
				return 0, 0, 1
			}
			return 0, 0, 0
		case "Tag":
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.TypeAssertExpr:
		switch fragment {
		case "X":
			return 0, 0, 0
		case "Lparen":
			return 0, 1, 0
		case "Type":
			// ************
			// SPECIAL CASE
			// If Type == nil, ".(type)" is rendered. If not, type node is rendered in parens.
			// ************
			if n.Type != nil {
				return 2, 0, 1 // ".(" and ")"
			}
			return 0, 7, 0 // ".(type)"
		case "Rparen":
			return 0, 1, 0
		}
	case *dst.TypeSpec:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Name":
			return 0, 0, 0
		case "Assign":
			if n.Assign {
				return 0, 1, 0
			}
			return 0, 0, 0
		case "Type":
			return 0, 0, 0
		case "Comment":
			return 0, 0, 0
		}
	case *dst.TypeSwitchStmt:
		switch fragment {
		case "Switch":
			return 0, 6, 0
		case "Init":
			if n.Init != nil {
				return 0, 0, 1
			}
			return 0, 0, 0
		case "Assign":
			if n.Assign != nil {
				return 0, 1, 0
			}
			return 0, 0, 0
		case "Body":
			return 0, 0, 0
		}
	case *dst.UnaryExpr:
		switch fragment {
		case "Op":
			if n.Op != token.ILLEGAL {
				return 0, len(n.Op.String()), 0
			}
			return 0, 0, 0
		case "X":
			return 0, 0, 0
		}
	case *dst.ValueSpec:
		switch fragment {
		case "Doc":
			return 0, 0, 0
		case "Names":
			return 0, 0, 0
		case "Type":
			return 0, 0, 0
		case "Values":
			if len(n.Values) > 0 {
				return 1, 0, 0
			}
			return 0, 0, 0
		case "Comment":
			return 0, 0, 0
		}
	}
	return 0, 0, 0
}
