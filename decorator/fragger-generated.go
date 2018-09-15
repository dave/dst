package decorator

import "go/ast"

func (f *Fragger) ProcessNode(n ast.Node) {
	f.ProcessToken(n, "", true, 0, n.Pos())
	switch n := n.(type) {
	case *ast.ArrayType:
		// Lbrack
		if n.Lbrack.IsValid() {
			f.ProcessToken(n, "Lbrack", false, 1, n.Lbrack)
		}
		// Len
		if n.Len != nil {
			f.ProcessToken(n, "Len", true, 0, n.Len.Pos())
			f.ProcessNode(n.Len)
			f.ProcessToken(n, "Len", false, 0, n.Len.End())
		}
		// Elt
		if n.Elt != nil {
			f.ProcessToken(n, "Elt", true, 0, n.Elt.Pos())
			f.ProcessNode(n.Elt)
			f.ProcessToken(n, "Elt", false, 0, n.Elt.End())
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
			f.ProcessToken(n, "Tok", false, len(n.Tok.String()), n.TokPos)
		}
		// Rhs
		if n.Rhs != nil {
			for _, v := range n.Rhs {
				f.ProcessNode(v)
			}
		}
	case *ast.BadDecl:
		// From
		if n.From.IsValid() {
			f.ProcessToken(n, "From", false, 0, n.From)
		}
		// To
		if n.To.IsValid() {
			f.ProcessToken(n, "To", false, 0, n.To)
		}
	case *ast.BadExpr:
		// From
		if n.From.IsValid() {
			f.ProcessToken(n, "From", false, 0, n.From)
		}
		// To
		if n.To.IsValid() {
			f.ProcessToken(n, "To", false, 0, n.To)
		}
	case *ast.BadStmt:
		// From
		if n.From.IsValid() {
			f.ProcessToken(n, "From", false, 0, n.From)
		}
		// To
		if n.To.IsValid() {
			f.ProcessToken(n, "To", false, 0, n.To)
		}
	case *ast.BasicLit:
		// Value
		if n.ValuePos.IsValid() {
			f.ProcessToken(n, "Value", false, len(n.Value), n.ValuePos)
		}
	case *ast.BinaryExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Op
		if n.OpPos.IsValid() {
			f.ProcessToken(n, "Op", false, len(n.Op.String()), n.OpPos)
		}
		// Y
		if n.Y != nil {
			f.ProcessToken(n, "Y", true, 0, n.Y.Pos())
			f.ProcessNode(n.Y)
			f.ProcessToken(n, "Y", false, 0, n.Y.End())
		}
	case *ast.BlockStmt:
		// Lbrace
		if n.Lbrace.IsValid() {
			f.ProcessToken(n, "Lbrace", false, 1, n.Lbrace)
		}
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
		// Rbrace
		if n.Rbrace.IsValid() {
			f.ProcessToken(n, "Rbrace", false, 1, n.Rbrace)
		}
	case *ast.BranchStmt:
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", false, len(n.Tok.String()), n.TokPos)
		}
		// Label
		if n.Label != nil {
			f.ProcessToken(n, "Label", true, 0, n.Label.Pos())
			f.ProcessNode(n.Label)
			f.ProcessToken(n, "Label", false, 0, n.Label.End())
		}
	case *ast.CallExpr:
		// Fun
		if n.Fun != nil {
			f.ProcessToken(n, "Fun", true, 0, n.Fun.Pos())
			f.ProcessNode(n.Fun)
			f.ProcessToken(n, "Fun", false, 0, n.Fun.End())
		}
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", false, 1, n.Lparen)
		}
		// Args
		if n.Args != nil {
			for _, v := range n.Args {
				f.ProcessNode(v)
			}
		}
		// Ellipsis
		if n.Ellipsis.IsValid() {
			f.ProcessToken(n, "Ellipsis", false, 3, n.Ellipsis)
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", false, 1, n.Rparen)
		}
	case *ast.CaseClause:
		// Case
		if n.Case.IsValid() {
			f.ProcessToken(n, "Case", false, 4, n.Case)
		}
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", false, 1, n.Colon)
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
			f.ProcessToken(n, "Begin", false, 4, n.Begin)
		}
		// Arrow
		if n.Arrow.IsValid() {
			f.ProcessToken(n, "Arrow", false, 2, n.Arrow)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", true, 0, n.Value.Pos())
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", false, 0, n.Value.End())
		}
	case *ast.CommClause:
		// Case
		if n.Case.IsValid() {
			f.ProcessToken(n, "Case", false, 4, n.Case)
		}
		// Comm
		if n.Comm != nil {
			f.ProcessToken(n, "Comm", true, 0, n.Comm.Pos())
			f.ProcessNode(n.Comm)
			f.ProcessToken(n, "Comm", false, 0, n.Comm.End())
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", false, 1, n.Colon)
		}
		// Body
		if n.Body != nil {
			for _, v := range n.Body {
				f.ProcessNode(v)
			}
		}
	case *ast.Comment:
		// Text
		if n.Slash.IsValid() {
			f.ProcessToken(n, "Text", false, len(n.Text), n.Slash)
		}
	case *ast.CommentGroup:
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
	case *ast.CompositeLit:
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Lbrace
		if n.Lbrace.IsValid() {
			f.ProcessToken(n, "Lbrace", false, 1, n.Lbrace)
		}
		// Elts
		if n.Elts != nil {
			for _, v := range n.Elts {
				f.ProcessNode(v)
			}
		}
		// Rbrace
		if n.Rbrace.IsValid() {
			f.ProcessToken(n, "Rbrace", false, 1, n.Rbrace)
		}
	case *ast.DeclStmt:
		// Decl
		if n.Decl != nil {
			f.ProcessToken(n, "Decl", true, 0, n.Decl.Pos())
			f.ProcessNode(n.Decl)
			f.ProcessToken(n, "Decl", false, 0, n.Decl.End())
		}
	case *ast.DeferStmt:
		// Defer
		if n.Defer.IsValid() {
			f.ProcessToken(n, "Defer", false, 5, n.Defer)
		}
		// Call
		if n.Call != nil {
			f.ProcessToken(n, "Call", true, 0, n.Call.Pos())
			f.ProcessNode(n.Call)
			f.ProcessToken(n, "Call", false, 0, n.Call.End())
		}
	case *ast.Ellipsis:
		// Ellipsis
		if n.Ellipsis.IsValid() {
			f.ProcessToken(n, "Ellipsis", false, 3, n.Ellipsis)
		}
		// Elt
		if n.Elt != nil {
			f.ProcessToken(n, "Elt", true, 0, n.Elt.Pos())
			f.ProcessNode(n.Elt)
			f.ProcessToken(n, "Elt", false, 0, n.Elt.End())
		}
	case *ast.EmptyStmt:
		// Semicolon
		if n.Semicolon.IsValid() {
			f.ProcessToken(n, "Semicolon", false, 0, n.Semicolon)
		}
	case *ast.ExprStmt:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
	case *ast.Field:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Names
		if n.Names != nil {
			for _, v := range n.Names {
				f.ProcessNode(v)
			}
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Tag
		if n.Tag != nil {
			f.ProcessToken(n, "Tag", true, 0, n.Tag.Pos())
			f.ProcessNode(n.Tag)
			f.ProcessToken(n, "Tag", false, 0, n.Tag.End())
		}
		// Comment
		if n.Comment != nil {
			f.ProcessToken(n, "Comment", true, 0, n.Comment.Pos())
			f.ProcessNode(n.Comment)
			f.ProcessToken(n, "Comment", false, 0, n.Comment.End())
		}
	case *ast.FieldList:
		// Opening
		if n.Opening.IsValid() {
			f.ProcessToken(n, "Opening", false, 1, n.Opening)
		}
		// List
		if n.List != nil {
			for _, v := range n.List {
				f.ProcessNode(v)
			}
		}
		// Closing
		if n.Closing.IsValid() {
			f.ProcessToken(n, "Closing", false, 1, n.Closing)
		}
	case *ast.File:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Package
		if n.Package.IsValid() {
			f.ProcessToken(n, "Package", false, 7, n.Package)
		}
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", true, 0, n.Name.Pos())
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", false, 0, n.Name.End())
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
			f.ProcessToken(n, "For", false, 3, n.For)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", true, 0, n.Init.Pos())
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", false, 0, n.Init.End())
		}
		// Cond
		if n.Cond != nil {
			f.ProcessToken(n, "Cond", true, 0, n.Cond.Pos())
			f.ProcessNode(n.Cond)
			f.ProcessToken(n, "Cond", false, 0, n.Cond.End())
		}
		// Post
		if n.Post != nil {
			f.ProcessToken(n, "Post", true, 0, n.Post.Pos())
			f.ProcessNode(n.Post)
			f.ProcessToken(n, "Post", false, 0, n.Post.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.FuncDecl:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Recv
		if n.Recv != nil {
			f.ProcessToken(n, "Recv", true, 0, n.Recv.Pos())
			f.ProcessNode(n.Recv)
			f.ProcessToken(n, "Recv", false, 0, n.Recv.End())
		}
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", true, 0, n.Name.Pos())
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", false, 0, n.Name.End())
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.FuncLit:
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.FuncType:
		// Func
		if n.Func.IsValid() {
			f.ProcessToken(n, "Func", false, 4, n.Func)
		}
		// Params
		if n.Params != nil {
			f.ProcessToken(n, "Params", true, 0, n.Params.Pos())
			f.ProcessNode(n.Params)
			f.ProcessToken(n, "Params", false, 0, n.Params.End())
		}
		// Results
		if n.Results != nil {
			f.ProcessToken(n, "Results", true, 0, n.Results.Pos())
			f.ProcessNode(n.Results)
			f.ProcessToken(n, "Results", false, 0, n.Results.End())
		}
	case *ast.GenDecl:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", false, len(n.Tok.String()), n.TokPos)
		}
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", false, 1, n.Lparen)
		}
		// Specs
		if n.Specs != nil {
			for _, v := range n.Specs {
				f.ProcessNode(v)
			}
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", false, 1, n.Rparen)
		}
	case *ast.GoStmt:
		// Go
		if n.Go.IsValid() {
			f.ProcessToken(n, "Go", false, 2, n.Go)
		}
		// Call
		if n.Call != nil {
			f.ProcessToken(n, "Call", true, 0, n.Call.Pos())
			f.ProcessNode(n.Call)
			f.ProcessToken(n, "Call", false, 0, n.Call.End())
		}
	case *ast.Ident:
		// Name
		if n.NamePos.IsValid() {
			f.ProcessToken(n, "Name", false, len(n.Name), n.NamePos)
		}
	case *ast.IfStmt:
		// If
		if n.If.IsValid() {
			f.ProcessToken(n, "If", false, 2, n.If)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", true, 0, n.Init.Pos())
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", false, 0, n.Init.End())
		}
		// Cond
		if n.Cond != nil {
			f.ProcessToken(n, "Cond", true, 0, n.Cond.Pos())
			f.ProcessNode(n.Cond)
			f.ProcessToken(n, "Cond", false, 0, n.Cond.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
		// Else
		if n.Else != nil {
			f.ProcessToken(n, "Else", true, 0, n.Else.Pos())
			f.ProcessNode(n.Else)
			f.ProcessToken(n, "Else", false, 0, n.Else.End())
		}
	case *ast.ImportSpec:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", true, 0, n.Name.Pos())
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", false, 0, n.Name.End())
		}
		// Path
		if n.Path != nil {
			f.ProcessToken(n, "Path", true, 0, n.Path.Pos())
			f.ProcessNode(n.Path)
			f.ProcessToken(n, "Path", false, 0, n.Path.End())
		}
		// Comment
		if n.Comment != nil {
			f.ProcessToken(n, "Comment", true, 0, n.Comment.Pos())
			f.ProcessNode(n.Comment)
			f.ProcessToken(n, "Comment", false, 0, n.Comment.End())
		}
	case *ast.IncDecStmt:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", false, len(n.Tok.String()), n.TokPos)
		}
	case *ast.IndexExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Lbrack
		if n.Lbrack.IsValid() {
			f.ProcessToken(n, "Lbrack", false, 1, n.Lbrack)
		}
		// Index
		if n.Index != nil {
			f.ProcessToken(n, "Index", true, 0, n.Index.Pos())
			f.ProcessNode(n.Index)
			f.ProcessToken(n, "Index", false, 0, n.Index.End())
		}
		// Rbrack
		if n.Rbrack.IsValid() {
			f.ProcessToken(n, "Rbrack", false, 1, n.Rbrack)
		}
	case *ast.InterfaceType:
		// Interface
		if n.Interface.IsValid() {
			f.ProcessToken(n, "Interface", false, 9, n.Interface)
		}
		// Methods
		if n.Methods != nil {
			f.ProcessToken(n, "Methods", true, 0, n.Methods.Pos())
			f.ProcessNode(n.Methods)
			f.ProcessToken(n, "Methods", false, 0, n.Methods.End())
		}
	case *ast.KeyValueExpr:
		// Key
		if n.Key != nil {
			f.ProcessToken(n, "Key", true, 0, n.Key.Pos())
			f.ProcessNode(n.Key)
			f.ProcessToken(n, "Key", false, 0, n.Key.End())
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", false, 1, n.Colon)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", true, 0, n.Value.Pos())
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", false, 0, n.Value.End())
		}
	case *ast.LabeledStmt:
		// Label
		if n.Label != nil {
			f.ProcessToken(n, "Label", true, 0, n.Label.Pos())
			f.ProcessNode(n.Label)
			f.ProcessToken(n, "Label", false, 0, n.Label.End())
		}
		// Colon
		if n.Colon.IsValid() {
			f.ProcessToken(n, "Colon", false, 1, n.Colon)
		}
		// Stmt
		if n.Stmt != nil {
			f.ProcessToken(n, "Stmt", true, 0, n.Stmt.Pos())
			f.ProcessNode(n.Stmt)
			f.ProcessToken(n, "Stmt", false, 0, n.Stmt.End())
		}
	case *ast.MapType:
		// Map
		if n.Map.IsValid() {
			f.ProcessToken(n, "Map", false, 3, n.Map)
		}
		// Key
		if n.Key != nil {
			f.ProcessToken(n, "Key", true, 0, n.Key.Pos())
			f.ProcessNode(n.Key)
			f.ProcessToken(n, "Key", false, 0, n.Key.End())
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", true, 0, n.Value.Pos())
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", false, 0, n.Value.End())
		}
	case *ast.ParenExpr:
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", false, 1, n.Lparen)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", false, 1, n.Rparen)
		}
	case *ast.RangeStmt:
		// For
		if n.For.IsValid() {
			f.ProcessToken(n, "For", false, 3, n.For)
		}
		// Key
		if n.Key != nil {
			f.ProcessToken(n, "Key", true, 0, n.Key.Pos())
			f.ProcessNode(n.Key)
			f.ProcessToken(n, "Key", false, 0, n.Key.End())
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", true, 0, n.Value.Pos())
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", false, 0, n.Value.End())
		}
		// Tok
		if n.TokPos.IsValid() {
			f.ProcessToken(n, "Tok", false, len(n.Tok.String()), n.TokPos)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.ReturnStmt:
		// Return
		if n.Return.IsValid() {
			f.ProcessToken(n, "Return", false, 6, n.Return)
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
			f.ProcessToken(n, "Select", false, 6, n.Select)
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.SelectorExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Sel
		if n.Sel != nil {
			f.ProcessToken(n, "Sel", true, 0, n.Sel.Pos())
			f.ProcessNode(n.Sel)
			f.ProcessToken(n, "Sel", false, 0, n.Sel.End())
		}
	case *ast.SendStmt:
		// Chan
		if n.Chan != nil {
			f.ProcessToken(n, "Chan", true, 0, n.Chan.Pos())
			f.ProcessNode(n.Chan)
			f.ProcessToken(n, "Chan", false, 0, n.Chan.End())
		}
		// Arrow
		if n.Arrow.IsValid() {
			f.ProcessToken(n, "Arrow", false, 2, n.Arrow)
		}
		// Value
		if n.Value != nil {
			f.ProcessToken(n, "Value", true, 0, n.Value.Pos())
			f.ProcessNode(n.Value)
			f.ProcessToken(n, "Value", false, 0, n.Value.End())
		}
	case *ast.SliceExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Lbrack
		if n.Lbrack.IsValid() {
			f.ProcessToken(n, "Lbrack", false, 1, n.Lbrack)
		}
		// Low
		if n.Low != nil {
			f.ProcessToken(n, "Low", true, 0, n.Low.Pos())
			f.ProcessNode(n.Low)
			f.ProcessToken(n, "Low", false, 0, n.Low.End())
		}
		// High
		if n.High != nil {
			f.ProcessToken(n, "High", true, 0, n.High.Pos())
			f.ProcessNode(n.High)
			f.ProcessToken(n, "High", false, 0, n.High.End())
		}
		// Max
		if n.Max != nil {
			f.ProcessToken(n, "Max", true, 0, n.Max.Pos())
			f.ProcessNode(n.Max)
			f.ProcessToken(n, "Max", false, 0, n.Max.End())
		}
		// Rbrack
		if n.Rbrack.IsValid() {
			f.ProcessToken(n, "Rbrack", false, 1, n.Rbrack)
		}
	case *ast.StarExpr:
		// Star
		if n.Star.IsValid() {
			f.ProcessToken(n, "Star", false, 1, n.Star)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
	case *ast.StructType:
		// Struct
		if n.Struct.IsValid() {
			f.ProcessToken(n, "Struct", false, 6, n.Struct)
		}
		// Fields
		if n.Fields != nil {
			f.ProcessToken(n, "Fields", true, 0, n.Fields.Pos())
			f.ProcessNode(n.Fields)
			f.ProcessToken(n, "Fields", false, 0, n.Fields.End())
		}
	case *ast.SwitchStmt:
		// Switch
		if n.Switch.IsValid() {
			f.ProcessToken(n, "Switch", false, 6, n.Switch)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", true, 0, n.Init.Pos())
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", false, 0, n.Init.End())
		}
		// Tag
		if n.Tag != nil {
			f.ProcessToken(n, "Tag", true, 0, n.Tag.Pos())
			f.ProcessNode(n.Tag)
			f.ProcessToken(n, "Tag", false, 0, n.Tag.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.TypeAssertExpr:
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
		// Lparen
		if n.Lparen.IsValid() {
			f.ProcessToken(n, "Lparen", false, 1, n.Lparen)
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Rparen
		if n.Rparen.IsValid() {
			f.ProcessToken(n, "Rparen", false, 1, n.Rparen)
		}
	case *ast.TypeSpec:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Name
		if n.Name != nil {
			f.ProcessToken(n, "Name", true, 0, n.Name.Pos())
			f.ProcessNode(n.Name)
			f.ProcessToken(n, "Name", false, 0, n.Name.End())
		}
		// Assign
		if n.Assign.IsValid() {
			f.ProcessToken(n, "Assign", false, 1, n.Assign)
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Comment
		if n.Comment != nil {
			f.ProcessToken(n, "Comment", true, 0, n.Comment.Pos())
			f.ProcessNode(n.Comment)
			f.ProcessToken(n, "Comment", false, 0, n.Comment.End())
		}
	case *ast.TypeSwitchStmt:
		// Switch
		if n.Switch.IsValid() {
			f.ProcessToken(n, "Switch", false, 6, n.Switch)
		}
		// Init
		if n.Init != nil {
			f.ProcessToken(n, "Init", true, 0, n.Init.Pos())
			f.ProcessNode(n.Init)
			f.ProcessToken(n, "Init", false, 0, n.Init.End())
		}
		// Assign
		if n.Assign != nil {
			f.ProcessToken(n, "Assign", true, 0, n.Assign.Pos())
			f.ProcessNode(n.Assign)
			f.ProcessToken(n, "Assign", false, 0, n.Assign.End())
		}
		// Body
		if n.Body != nil {
			f.ProcessToken(n, "Body", true, 0, n.Body.Pos())
			f.ProcessNode(n.Body)
			f.ProcessToken(n, "Body", false, 0, n.Body.End())
		}
	case *ast.UnaryExpr:
		// Op
		if n.OpPos.IsValid() {
			f.ProcessToken(n, "Op", false, len(n.Op.String()), n.OpPos)
		}
		// X
		if n.X != nil {
			f.ProcessToken(n, "X", true, 0, n.X.Pos())
			f.ProcessNode(n.X)
			f.ProcessToken(n, "X", false, 0, n.X.End())
		}
	case *ast.ValueSpec:
		// Doc
		if n.Doc != nil {
			f.ProcessToken(n, "Doc", true, 0, n.Doc.Pos())
			f.ProcessNode(n.Doc)
			f.ProcessToken(n, "Doc", false, 0, n.Doc.End())
		}
		// Names
		if n.Names != nil {
			for _, v := range n.Names {
				f.ProcessNode(v)
			}
		}
		// Type
		if n.Type != nil {
			f.ProcessToken(n, "Type", true, 0, n.Type.Pos())
			f.ProcessNode(n.Type)
			f.ProcessToken(n, "Type", false, 0, n.Type.End())
		}
		// Values
		if n.Values != nil {
			for _, v := range n.Values {
				f.ProcessNode(v)
			}
		}
		// Comment
		if n.Comment != nil {
			f.ProcessToken(n, "Comment", true, 0, n.Comment.Pos())
			f.ProcessNode(n.Comment)
			f.ProcessToken(n, "Comment", false, 0, n.Comment.End())
		}
	}
	f.ProcessToken(n, "", false, 0, n.End())
}
