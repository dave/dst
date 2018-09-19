package decorator

import "go/ast"

func (f *Fragger) ProcessNode(n ast.Node) {
	f.ProcessToken(n, "", n.Pos(), false)
	switch n := n.(type) {
	case *ast.ArrayType:
		// Lbrack
		if n.Lbrack.IsValid() {
			f.ProcessToken(n, "Lbrack", n.Lbrack, false)
			f.ProcessToken(n, "Lbrack", n.Lbrack, true)
		}
		// Len
		if n.Len != nil {
			f.ProcessToken(n, "Len", n.Len.Pos(), false)
			f.ProcessNode(n.Len)
			f.ProcessToken(n, "Len", n.Len.End(), true)
		}
		// Elt
		if n.Elt != nil {
			f.ProcessToken(n, "Elt", n.Elt.Pos(), false)
			f.ProcessNode(n.Elt)
			f.ProcessToken(n, "Elt", n.Elt.End(), true)
		}
	case *ast.AssignStmt:
		// Lhs
		if n.Lhs != nil {
			for _, v := range n.Lhs {
				f.ProcessNode(v)
			}
		}
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", n.TokPos, false)
			f.ProcessToken(n, "Tok", n.TokPos, true)
		}
		// Rhs
		if n.Rhs != nil {
			for _, v := range n.Rhs {
				f.ProcessNode(v)
			}
		}
	case *ast.BadDecl:
	case *ast.BadExpr:
	case *ast.BadStmt:
	case *ast.BasicLit:
		// Value
		if n.ValuePos.IsValid() {
			f.ProcessToken(n, "Value", n.ValuePos, false)
			f.ProcessToken(n, "Value", n.ValuePos, true)
		}
	case *ast.BinaryExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Op
		if n.OpPos.IsValid() {
			f.ProcessToken(n, "Op", n.OpPos, false)
			f.ProcessToken(n, "Op", n.OpPos, true)
		}
		// Y
		if n.Y != nil {
			f.ProcessToken(n, "Y", n.Y.Pos(), false)
			f.ProcessNode(n.Y)
			f.ProcessToken(n, "Y", n.Y.End(), true)
		}
	case *ast.BlockStmt:
		// Lbrace
		if n.Lbrace.IsValid() {
			f.ProcessToken(n, "Lbrace", n.Lbrace, false)
			f.ProcessToken(n, "Lbrace", n.Lbrace, true)
		}
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
		// Rbrace
		if n.Rbrace.IsValid() {
			f.ProcessToken(n, "Rbrace", n.Rbrace, false)
			f.ProcessToken(n, "Rbrace", n.Rbrace, true)
		}
	case *ast.BranchStmt:
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", n.TokPos, false)
			f.ProcessToken(n, "Tok", n.TokPos, true)
		}
		// Label
		if n.Label != nil {
			f.ProcessToken(n, "Label", n.Label.Pos(), false)
			f.ProcessNode(n.Label)
			f.ProcessToken(n, "Label", n.Label.End(), true)
		}
	case *ast.CallExpr:
		// Fun
		if n.Fun != nil {
			f.ProcessToken(n, "Fun", n.Fun.Pos(), false)
			f.ProcessNode(n.Fun)
			f.ProcessToken(n, "Fun", n.Fun.End(), true)
		}
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", n.Lparen, false)
			f.ProcessToken(n, "Lparen", n.Lparen, true)
		}
		// Args
		if n.Args != nil {
			for _, v := range n.Args {
				f.ProcessNode(v)
			}
		}
		// Ellipsis
		if n.Ellipsis.IsValid() {
			f.ProcessToken(n, "Ellipsis", n.Ellipsis, false)
			f.ProcessToken(n, "Ellipsis", n.Ellipsis, true)
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", n.Rparen, false)
			f.ProcessToken(n, "Rparen", n.Rparen, true)
		}
	case *ast.CaseClause:
		// Case
		if n.Case.IsValid() {
			f.ProcessToken(n, "Case", n.Case, false)
			f.ProcessToken(n, "Case", n.Case, true)
		}
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", n.Colon, false)
			f.ProcessToken(n, "Colon", n.Colon, true)
		}
		// Body
		if n.Body != nil {
			for _, v := range n.Body {
				f.ProcessNode(v)
			}
		}
	case *ast.ChanType:
		// Begin
		if n.Begin.IsValid() {
			f.ProcessToken(n, "Begin", n.Begin, false)
			f.ProcessToken(n, "Begin", n.Begin, true)
		}
		// Arrow
		if n.Arrow.IsValid() {
			f.ProcessToken(n, "Arrow", n.Arrow, false)
			f.ProcessToken(n, "Arrow", n.Arrow, true)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", n.Value.Pos(), false)
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", n.Value.End(), true)
		}
	case *ast.CommClause:
		// Case
		if n.Case.IsValid() {
			f.ProcessToken(n, "Case", n.Case, false)
			f.ProcessToken(n, "Case", n.Case, true)
		}
		// Comm
		if n.Comm != nil {
			f.ProcessToken(n, "Comm", n.Comm.Pos(), false)
			f.ProcessNode(n.Comm)
			f.ProcessToken(n, "Comm", n.Comm.End(), true)
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", n.Colon, false)
			f.ProcessToken(n, "Colon", n.Colon, true)
		}
		// Body
		if n.Body != nil {
			for _, v := range n.Body {
				f.ProcessNode(v)
			}
		}
	case *ast.CompositeLit:
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", n.Type.Pos(), false)
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", n.Type.End(), true)
		}
		// Lbrace
		if n.Lbrace.IsValid() {
			f.ProcessToken(n, "Lbrace", n.Lbrace, false)
			f.ProcessToken(n, "Lbrace", n.Lbrace, true)
		}
		// Elts
		if n.Elts != nil {
			for _, v := range n.Elts {
				f.ProcessNode(v)
			}
		}
		// Rbrace
		if n.Rbrace.IsValid() {
			f.ProcessToken(n, "Rbrace", n.Rbrace, false)
			f.ProcessToken(n, "Rbrace", n.Rbrace, true)
		}
	case *ast.DeclStmt:
		// Decl
		if n.Decl != nil {
			f.ProcessToken(n, "Decl", n.Decl.Pos(), false)
			f.ProcessNode(n.Decl)
			f.ProcessToken(n, "Decl", n.Decl.End(), true)
		}
	case *ast.DeferStmt:
		// Defer
		if n.Defer.IsValid() {
			f.ProcessToken(n, "Defer", n.Defer, false)
			f.ProcessToken(n, "Defer", n.Defer, true)
		}
		// Call
		if n.Call != nil {
			f.ProcessToken(n, "Call", n.Call.Pos(), false)
			f.ProcessNode(n.Call)
			f.ProcessToken(n, "Call", n.Call.End(), true)
		}
	case *ast.Ellipsis:
		// Ellipsis
		if n.Ellipsis.IsValid() {
			f.ProcessToken(n, "Ellipsis", n.Ellipsis, false)
			f.ProcessToken(n, "Ellipsis", n.Ellipsis, true)
		}
		// Elt
		if n.Elt != nil {
			f.ProcessToken(n, "Elt", n.Elt.Pos(), false)
			f.ProcessNode(n.Elt)
			f.ProcessToken(n, "Elt", n.Elt.End(), true)
		}
	case *ast.EmptyStmt:
		// Semicolon
		if n.Semicolon.IsValid() {
			f.ProcessToken(n, "Semicolon", n.Semicolon, false)
			f.ProcessToken(n, "Semicolon", n.Semicolon, true)
		}
	case *ast.ExprStmt:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
	case *ast.Field:
		// Names
		if n.Names != nil {
			for _, v := range n.Names {
				f.ProcessNode(v)
			}
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", n.Type.Pos(), false)
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", n.Type.End(), true)
		}
		// Tag
		if n.Tag != nil {
			f.ProcessToken(n, "Tag", n.Tag.Pos(), false)
			f.ProcessNode(n.Tag)
			f.ProcessToken(n, "Tag", n.Tag.End(), true)
		}
	case *ast.FieldList:
		// Opening
		if n.Opening.IsValid() {
			f.ProcessToken(n, "Opening", n.Opening, false)
			f.ProcessToken(n, "Opening", n.Opening, true)
		}
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
		// Closing
		if n.Closing.IsValid() {
			f.ProcessToken(n, "Closing", n.Closing, false)
			f.ProcessToken(n, "Closing", n.Closing, true)
		}
	case *ast.File:
		// Package
		if n.Package.IsValid() {
			f.ProcessToken(n, "Package", n.Package, false)
			f.ProcessToken(n, "Package", n.Package, true)
		}
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", n.Name.Pos(), false)
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", n.Name.End(), true)
		}
		// Decls
		if n.Decls != nil {
			for _, v := range n.Decls {
				f.ProcessNode(v)
			}
		}
	case *ast.ForStmt:
		// For
		if n.For.IsValid() {
			f.ProcessToken(n, "For", n.For, false)
			f.ProcessToken(n, "For", n.For, true)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", n.Init.Pos(), false)
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", n.Init.End(), true)
		}
		// Cond
		if n.Cond != nil {
			f.ProcessToken(n, "Cond", n.Cond.Pos(), false)
			f.ProcessNode(n.Cond)
			f.ProcessToken(n, "Cond", n.Cond.End(), true)
		}
		// Post
		if n.Post != nil {
			f.ProcessToken(n, "Post", n.Post.Pos(), false)
			f.ProcessNode(n.Post)
			f.ProcessToken(n, "Post", n.Post.End(), true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
	case *ast.FuncDecl:
		f.funcDeclOverride(n)
	case *ast.FuncLit:
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", n.Type.Pos(), false)
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", n.Type.End(), true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
	case *ast.FuncType:
		// Func
		if n.Func.IsValid() {
			f.ProcessToken(n, "Func", n.Func, false)
			f.ProcessToken(n, "Func", n.Func, true)
		}
		// Params
		if n.Params != nil {
			f.ProcessToken(n, "Params", n.Params.Pos(), false)
			f.ProcessNode(n.Params)
			f.ProcessToken(n, "Params", n.Params.End(), true)
		}
		// Results
		if n.Results != nil {
			f.ProcessToken(n, "Results", n.Results.Pos(), false)
			f.ProcessNode(n.Results)
			f.ProcessToken(n, "Results", n.Results.End(), true)
		}
	case *ast.GenDecl:
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", n.TokPos, false)
			f.ProcessToken(n, "Tok", n.TokPos, true)
		}
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", n.Lparen, false)
			f.ProcessToken(n, "Lparen", n.Lparen, true)
		}
		// Specs
		if n.Specs != nil {
			for _, v := range n.Specs {
				f.ProcessNode(v)
			}
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", n.Rparen, false)
			f.ProcessToken(n, "Rparen", n.Rparen, true)
		}
	case *ast.GoStmt:
		// Go
		if n.Go.IsValid() {
			f.ProcessToken(n, "Go", n.Go, false)
			f.ProcessToken(n, "Go", n.Go, true)
		}
		// Call
		if n.Call != nil {
			f.ProcessToken(n, "Call", n.Call.Pos(), false)
			f.ProcessNode(n.Call)
			f.ProcessToken(n, "Call", n.Call.End(), true)
		}
	case *ast.Ident:
		// Name
		if n.NamePos.IsValid() {
			f.ProcessToken(n, "Name", n.NamePos, false)
			f.ProcessToken(n, "Name", n.NamePos, true)
		}
	case *ast.IfStmt:
		// If
		if n.If.IsValid() {
			f.ProcessToken(n, "If", n.If, false)
			f.ProcessToken(n, "If", n.If, true)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", n.Init.Pos(), false)
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", n.Init.End(), true)
		}
		// Cond
		if n.Cond != nil {
			f.ProcessToken(n, "Cond", n.Cond.Pos(), false)
			f.ProcessNode(n.Cond)
			f.ProcessToken(n, "Cond", n.Cond.End(), true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
		// Else
		if n.Else != nil {
			f.ProcessToken(n, "Else", n.Else.Pos(), false)
			f.ProcessNode(n.Else)
			f.ProcessToken(n, "Else", n.Else.End(), true)
		}
	case *ast.ImportSpec:
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", n.Name.Pos(), false)
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", n.Name.End(), true)
		}
		// Path
		if n.Path != nil {
			f.ProcessToken(n, "Path", n.Path.Pos(), false)
			f.ProcessNode(n.Path)
			f.ProcessToken(n, "Path", n.Path.End(), true)
		}
	case *ast.IncDecStmt:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", n.TokPos, false)
			f.ProcessToken(n, "Tok", n.TokPos, true)
		}
	case *ast.IndexExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Lbrack
		if n.Lbrack.IsValid() {
			f.ProcessToken(n, "Lbrack", n.Lbrack, false)
			f.ProcessToken(n, "Lbrack", n.Lbrack, true)
		}
		// Index
		if n.Index != nil {
			f.ProcessToken(n, "Index", n.Index.Pos(), false)
			f.ProcessNode(n.Index)
			f.ProcessToken(n, "Index", n.Index.End(), true)
		}
		// Rbrack
		if n.Rbrack.IsValid() {
			f.ProcessToken(n, "Rbrack", n.Rbrack, false)
			f.ProcessToken(n, "Rbrack", n.Rbrack, true)
		}
	case *ast.InterfaceType:
		// Interface
		if n.Interface.IsValid() {
			f.ProcessToken(n, "Interface", n.Interface, false)
			f.ProcessToken(n, "Interface", n.Interface, true)
		}
		// Methods
		if n.Methods != nil {
			f.ProcessToken(n, "Methods", n.Methods.Pos(), false)
			f.ProcessNode(n.Methods)
			f.ProcessToken(n, "Methods", n.Methods.End(), true)
		}
	case *ast.KeyValueExpr:
		// Key
		if n.Key != nil {
			f.ProcessToken(n, "Key", n.Key.Pos(), false)
			f.ProcessNode(n.Key)
			f.ProcessToken(n, "Key", n.Key.End(), true)
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", n.Colon, false)
			f.ProcessToken(n, "Colon", n.Colon, true)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", n.Value.Pos(), false)
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", n.Value.End(), true)
		}
	case *ast.LabeledStmt:
		// Label
		if n.Label != nil {
			f.ProcessToken(n, "Label", n.Label.Pos(), false)
			f.ProcessNode(n.Label)
			f.ProcessToken(n, "Label", n.Label.End(), true)
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", n.Colon, false)
			f.ProcessToken(n, "Colon", n.Colon, true)
		}
		// Stmt
		if n.Stmt != nil {
			f.ProcessToken(n, "Stmt", n.Stmt.Pos(), false)
			f.ProcessNode(n.Stmt)
			f.ProcessToken(n, "Stmt", n.Stmt.End(), true)
		}
	case *ast.MapType:
		// Map
		if n.Map.IsValid() {
			f.ProcessToken(n, "Map", n.Map, false)
			f.ProcessToken(n, "Map", n.Map, true)
		}
		// Key
		if n.Key != nil {
			f.ProcessToken(n, "Key", n.Key.Pos(), false)
			f.ProcessNode(n.Key)
			f.ProcessToken(n, "Key", n.Key.End(), true)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", n.Value.Pos(), false)
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", n.Value.End(), true)
		}
	case *ast.ParenExpr:
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", n.Lparen, false)
			f.ProcessToken(n, "Lparen", n.Lparen, true)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", n.Rparen, false)
			f.ProcessToken(n, "Rparen", n.Rparen, true)
		}
	case *ast.RangeStmt:
		// For
		if n.For.IsValid() {
			f.ProcessToken(n, "For", n.For, false)
			f.ProcessToken(n, "For", n.For, true)
		}
		// Key
		if n.Key != nil {
			f.ProcessToken(n, "Key", n.Key.Pos(), false)
			f.ProcessNode(n.Key)
			f.ProcessToken(n, "Key", n.Key.End(), true)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", n.Value.Pos(), false)
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", n.Value.End(), true)
		}
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", n.TokPos, false)
			f.ProcessToken(n, "Tok", n.TokPos, true)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
	case *ast.ReturnStmt:
		// Return
		if n.Return.IsValid() {
			f.ProcessToken(n, "Return", n.Return, false)
			f.ProcessToken(n, "Return", n.Return, true)
		}
		// Results
		if n.Results != nil {
			for _, v := range n.Results {
				f.ProcessNode(v)
			}
		}
	case *ast.SelectStmt:
		// Select
		if n.Select.IsValid() {
			f.ProcessToken(n, "Select", n.Select, false)
			f.ProcessToken(n, "Select", n.Select, true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
	case *ast.SelectorExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Sel
		if n.Sel != nil {
			f.ProcessToken(n, "Sel", n.Sel.Pos(), false)
			f.ProcessNode(n.Sel)
			f.ProcessToken(n, "Sel", n.Sel.End(), true)
		}
	case *ast.SendStmt:
		// Chan
		if n.Chan != nil {
			f.ProcessToken(n, "Chan", n.Chan.Pos(), false)
			f.ProcessNode(n.Chan)
			f.ProcessToken(n, "Chan", n.Chan.End(), true)
		}
		// Arrow
		if n.Arrow.IsValid() {
			f.ProcessToken(n, "Arrow", n.Arrow, false)
			f.ProcessToken(n, "Arrow", n.Arrow, true)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", n.Value.Pos(), false)
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", n.Value.End(), true)
		}
	case *ast.SliceExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Lbrack
		if n.Lbrack.IsValid() {
			f.ProcessToken(n, "Lbrack", n.Lbrack, false)
			f.ProcessToken(n, "Lbrack", n.Lbrack, true)
		}
		// Low
		if n.Low != nil {
			f.ProcessToken(n, "Low", n.Low.Pos(), false)
			f.ProcessNode(n.Low)
			f.ProcessToken(n, "Low", n.Low.End(), true)
		}
		// High
		if n.High != nil {
			f.ProcessToken(n, "High", n.High.Pos(), false)
			f.ProcessNode(n.High)
			f.ProcessToken(n, "High", n.High.End(), true)
		}
		// Max
		if n.Max != nil {
			f.ProcessToken(n, "Max", n.Max.Pos(), false)
			f.ProcessNode(n.Max)
			f.ProcessToken(n, "Max", n.Max.End(), true)
		}
		// Rbrack
		if n.Rbrack.IsValid() {
			f.ProcessToken(n, "Rbrack", n.Rbrack, false)
			f.ProcessToken(n, "Rbrack", n.Rbrack, true)
		}
	case *ast.StarExpr:
		// Star
		if n.Star.IsValid() {
			f.ProcessToken(n, "Star", n.Star, false)
			f.ProcessToken(n, "Star", n.Star, true)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
	case *ast.StructType:
		// Struct
		if n.Struct.IsValid() {
			f.ProcessToken(n, "Struct", n.Struct, false)
			f.ProcessToken(n, "Struct", n.Struct, true)
		}
		// Fields
		if n.Fields != nil {
			f.ProcessToken(n, "Fields", n.Fields.Pos(), false)
			f.ProcessNode(n.Fields)
			f.ProcessToken(n, "Fields", n.Fields.End(), true)
		}
	case *ast.SwitchStmt:
		// Switch
		if n.Switch.IsValid() {
			f.ProcessToken(n, "Switch", n.Switch, false)
			f.ProcessToken(n, "Switch", n.Switch, true)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", n.Init.Pos(), false)
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", n.Init.End(), true)
		}
		// Tag
		if n.Tag != nil {
			f.ProcessToken(n, "Tag", n.Tag.Pos(), false)
			f.ProcessNode(n.Tag)
			f.ProcessToken(n, "Tag", n.Tag.End(), true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
	case *ast.TypeAssertExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", n.Lparen, false)
			f.ProcessToken(n, "Lparen", n.Lparen, true)
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", n.Type.Pos(), false)
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", n.Type.End(), true)
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", n.Rparen, false)
			f.ProcessToken(n, "Rparen", n.Rparen, true)
		}
	case *ast.TypeSpec:
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", n.Name.Pos(), false)
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", n.Name.End(), true)
		}
		// Assign
		if n.Assign.IsValid() {
			f.ProcessToken(n, "Assign", n.Assign, false)
			f.ProcessToken(n, "Assign", n.Assign, true)
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", n.Type.Pos(), false)
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", n.Type.End(), true)
		}
	case *ast.TypeSwitchStmt:
		// Switch
		if n.Switch.IsValid() {
			f.ProcessToken(n, "Switch", n.Switch, false)
			f.ProcessToken(n, "Switch", n.Switch, true)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", n.Init.Pos(), false)
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", n.Init.End(), true)
		}
		// Assign
		if n.Assign != nil {
			f.ProcessToken(n, "Assign", n.Assign.Pos(), false)
			f.ProcessNode(n.Assign)
			f.ProcessToken(n, "Assign", n.Assign.End(), true)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", n.Body.Pos(), false)
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", n.Body.End(), true)
		}
	case *ast.UnaryExpr:
		// Op
		if n.OpPos.IsValid() {
			f.ProcessToken(n, "Op", n.OpPos, false)
			f.ProcessToken(n, "Op", n.OpPos, true)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", n.X.Pos(), false)
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", n.X.End(), true)
		}
	case *ast.ValueSpec:
		// Names
		if n.Names != nil {
			for _, v := range n.Names {
				f.ProcessNode(v)
			}
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", n.Type.Pos(), false)
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", n.Type.End(), true)
		}
		// Values
		if n.Values != nil {
			for _, v := range n.Values {
				f.ProcessNode(v)
			}
		}
	}
	f.ProcessToken(n, "", n.End(), true)
}
