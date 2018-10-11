package decorator

import (
	"go/ast"
	"go/token"
)

func (f *Fragger) ProcessNode(n ast.Node) {
	if n.Pos().IsValid() {
		f.cursor = int(n.Pos())
	}
	switch n := n.(type) {
	case *ast.ArrayType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Lbrack
		f.AddToken(n, token.LBRACK, n.Lbrack)

		// Decoration: Lbrack
		f.AddDecoration(n, "Lbrack", token.NoPos)

		// Node: Len
		if n.Len != nil {
			f.ProcessNode(n.Len)
		}

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, token.NoPos)

		// Decoration: Len
		f.AddDecoration(n, "Len", token.NoPos)

		// Node: Elt
		if n.Elt != nil {
			f.ProcessNode(n.Elt)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.AssignStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// List: Lhs
		for _, v := range n.Lhs {
			f.ProcessNode(v)
		}

		// Decoration: Lhs
		f.AddDecoration(n, "Lhs", token.NoPos)

		// Token: Tok
		f.AddToken(n, n.Tok, n.TokPos)

		// Decoration: Tok
		f.AddDecoration(n, "Tok", token.NoPos)

		// List: Rhs
		for _, v := range n.Rhs {
			f.ProcessNode(v)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BadDecl:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BadExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BadStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BasicLit:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// String: Value
		f.AddString(n, n.Value, n.ValuePos)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BinaryExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Token: Op
		f.AddToken(n, n.Op, n.OpPos)

		// Decoration: Op
		f.AddDecoration(n, "Op", token.NoPos)

		// Node: Y
		if n.Y != nil {
			f.ProcessNode(n.Y)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BlockStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Lbrace
		f.AddToken(n, token.LBRACE, n.Lbrace)

		// Decoration: Lbrace
		f.AddDecoration(n, "Lbrace", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.ProcessNode(v)
		}

		// Token: Rbrace
		f.AddToken(n, token.RBRACE, n.Rbrace)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BranchStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Tok
		f.AddToken(n, n.Tok, n.TokPos)

		// Decoration: Tok
		if n.Label != nil {
			f.AddDecoration(n, "Tok", token.NoPos)
		}

		// Node: Label
		if n.Label != nil {
			f.ProcessNode(n.Label)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.CallExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Fun
		if n.Fun != nil {
			f.ProcessNode(n.Fun)
		}

		// Decoration: Fun
		f.AddDecoration(n, "Fun", token.NoPos)

		// Token: Lparen
		f.AddToken(n, token.LPAREN, n.Lparen)

		// Decoration: Lparen
		f.AddDecoration(n, "Lparen", token.NoPos)

		// List: Args
		for _, v := range n.Args {
			f.ProcessNode(v)
		}

		// Decoration: Args
		f.AddDecoration(n, "Args", token.NoPos)

		// Token: Ellipsis
		if n.Ellipsis.IsValid() {
			f.AddToken(n, token.ELLIPSIS, n.Ellipsis)
		}

		// Decoration: Ellipsis
		if n.Ellipsis.IsValid() {
			f.AddDecoration(n, "Ellipsis", token.NoPos)
		}

		// Token: Rparen
		f.AddToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.CaseClause:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Case
		f.AddToken(n, func() token.Token {
			if n.List == nil {
				return token.DEFAULT
			} else {
				return token.CASE
			}
		}(), n.Case)

		// Decoration: Case
		f.AddDecoration(n, "Case", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.ProcessNode(v)
		}

		// Decoration: List
		if n.List != nil {
			f.AddDecoration(n, "List", token.NoPos)
		}

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.AddDecoration(n, "Colon", token.NoPos)

		// List: Body
		for _, v := range n.Body {
			f.ProcessNode(v)
		}

		// Decoration: End
		if false {
			f.AddDecoration(n, "End", n.End())
		}

	case *ast.ChanType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Begin
		f.AddToken(n, func() token.Token {
			if n.Dir == ast.RECV {
				return token.ARROW
			} else {
				return token.CHAN
			}
		}(), n.Begin)

		// Token: Chan
		if n.Dir == ast.RECV {
			f.AddToken(n, token.CHAN, token.NoPos)
		}

		// Decoration: Begin
		f.AddDecoration(n, "Begin", token.NoPos)

		// Token: Arrow
		if n.Dir == ast.SEND {
			f.AddToken(n, token.ARROW, n.Arrow)
		}

		// Decoration: Arrow
		if n.Dir == ast.SEND {
			f.AddDecoration(n, "Arrow", token.NoPos)
		}

		// Node: Value
		if n.Value != nil {
			f.ProcessNode(n.Value)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.CommClause:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Case
		f.AddToken(n, func() token.Token {
			if n.Comm == nil {
				return token.DEFAULT
			} else {
				return token.CASE
			}
		}(), n.Case)

		// Decoration: Case
		f.AddDecoration(n, "Case", token.NoPos)

		// Node: Comm
		if n.Comm != nil {
			f.ProcessNode(n.Comm)
		}

		// Decoration: Comm
		if n.Comm != nil {
			f.AddDecoration(n, "Comm", token.NoPos)
		}

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.AddDecoration(n, "Colon", token.NoPos)

		// List: Body
		for _, v := range n.Body {
			f.ProcessNode(v)
		}

		// Decoration: End
		if false {
			f.AddDecoration(n, "End", n.End())
		}

	case *ast.CompositeLit:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Decoration: Type
		if n.Type != nil {
			f.AddDecoration(n, "Type", token.NoPos)
		}

		// Token: Lbrace
		f.AddToken(n, token.LBRACE, n.Lbrace)

		// Decoration: Lbrace
		f.AddDecoration(n, "Lbrace", token.NoPos)

		// List: Elts
		for _, v := range n.Elts {
			f.ProcessNode(v)
		}

		// Token: Rbrace
		f.AddToken(n, token.RBRACE, n.Rbrace)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.DeclStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Decl
		if n.Decl != nil {
			f.ProcessNode(n.Decl)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.DeferStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Defer
		f.AddToken(n, token.DEFER, n.Defer)

		// Decoration: Defer
		f.AddDecoration(n, "Defer", token.NoPos)

		// Node: Call
		if n.Call != nil {
			f.ProcessNode(n.Call)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.Ellipsis:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Ellipsis
		f.AddToken(n, token.ELLIPSIS, n.Ellipsis)

		// Decoration: Ellipsis
		if n.Elt != nil {
			f.AddDecoration(n, "Ellipsis", token.NoPos)
		}

		// Node: Elt
		if n.Elt != nil {
			f.ProcessNode(n.Elt)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.EmptyStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Semicolon
		if !n.Implicit {
			f.AddToken(n, token.ARROW, n.Semicolon)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.ExprStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.Field:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// List: Names
		for _, v := range n.Names {
			f.ProcessNode(v)
		}

		// Decoration: Names
		f.AddDecoration(n, "Names", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Decoration: Type
		if n.Tag != nil {
			f.AddDecoration(n, "Type", token.NoPos)
		}

		// Node: Tag
		if n.Tag != nil {
			f.ProcessNode(n.Tag)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.FieldList:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Opening
		if n.Opening.IsValid() {
			f.AddToken(n, token.LPAREN, n.Opening)
		}

		// Decoration: Opening
		f.AddDecoration(n, "Opening", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.ProcessNode(v)
		}

		// Token: Closing
		if n.Closing.IsValid() {
			f.AddToken(n, token.RPAREN, n.Closing)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.File:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Package
		f.AddToken(n, token.PACKAGE, n.Package)

		// Decoration: Package
		f.AddDecoration(n, "Package", token.NoPos)

		// Node: Name
		if n.Name != nil {
			f.ProcessNode(n.Name)
		}

		// Decoration: Name
		f.AddDecoration(n, "Name", token.NoPos)

		// List: Decls
		for _, v := range n.Decls {
			f.ProcessNode(v)
		}

	case *ast.ForStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: For
		f.AddToken(n, token.FOR, n.For)

		// Decoration: For
		f.AddDecoration(n, "For", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Token: InitSemicolon
		if n.Init != nil {
			f.AddToken(n, token.SEMICOLON, token.NoPos)
		}

		// Decoration: Init
		if n.Init != nil {
			f.AddDecoration(n, "Init", token.NoPos)
		}

		// Node: Cond
		if n.Cond != nil {
			f.ProcessNode(n.Cond)
		}

		// Token: CondSemicolon
		if n.Post != nil {
			f.AddToken(n, token.SEMICOLON, token.NoPos)
		}

		// Decoration: Cond
		if n.Cond != nil {
			f.AddDecoration(n, "Cond", token.NoPos)
		}

		// Node: Post
		if n.Post != nil {
			f.ProcessNode(n.Post)
		}

		// Decoration: Post
		if n.Post != nil {
			f.AddDecoration(n, "Post", token.NoPos)
		}

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.FuncDecl:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Func
		if true {
			f.AddToken(n, token.FUNC, n.Type.Func)
		}

		// Decoration: Func
		f.AddDecoration(n, "Func", token.NoPos)

		// Node: Recv
		if n.Recv != nil {
			f.ProcessNode(n.Recv)
		}

		// Decoration: Recv
		if n.Recv != nil {
			f.AddDecoration(n, "Recv", token.NoPos)
		}

		// Node: Name
		if n.Name != nil {
			f.ProcessNode(n.Name)
		}

		// Decoration: Name
		f.AddDecoration(n, "Name", token.NoPos)

		// Node: Params
		if n.Type.Params != nil {
			f.ProcessNode(n.Type.Params)
		}

		// Decoration: Params
		f.AddDecoration(n, "Params", token.NoPos)

		// Node: Results
		if n.Type.Results != nil {
			f.ProcessNode(n.Type.Results)
		}

		// Decoration: Results
		if n.Type.Results != nil {
			f.AddDecoration(n, "Results", token.NoPos)
		}

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.FuncLit:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Decoration: Type
		f.AddDecoration(n, "Type", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.FuncType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Func
		if n.Func.IsValid() {
			f.AddToken(n, token.FUNC, n.Func)
		}

		// Decoration: Func
		if n.Func.IsValid() {
			f.AddDecoration(n, "Func", token.NoPos)
		}

		// Node: Params
		if n.Params != nil {
			f.ProcessNode(n.Params)
		}

		// Decoration: Params
		if n.Results != nil {
			f.AddDecoration(n, "Params", token.NoPos)
		}

		// Node: Results
		if n.Results != nil {
			f.ProcessNode(n.Results)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.GenDecl:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Tok
		f.AddToken(n, n.Tok, n.TokPos)

		// Decoration: Tok
		f.AddDecoration(n, "Tok", token.NoPos)

		// Token: Lparen
		if n.Lparen.IsValid() {
			f.AddToken(n, token.LPAREN, n.Lparen)
		}

		// Decoration: Lparen
		if n.Lparen.IsValid() {
			f.AddDecoration(n, "Lparen", token.NoPos)
		}

		// List: Specs
		for _, v := range n.Specs {
			f.ProcessNode(v)
		}

		// Token: Rparen
		if n.Rparen.IsValid() {
			f.AddToken(n, token.RPAREN, n.Rparen)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.GoStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Go
		f.AddToken(n, token.GO, n.Go)

		// Decoration: Go
		f.AddDecoration(n, "Go", token.NoPos)

		// Node: Call
		if n.Call != nil {
			f.ProcessNode(n.Call)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.Ident:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// String: Name
		f.AddString(n, n.Name, n.NamePos)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.IfStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: If
		f.AddToken(n, token.IF, n.If)

		// Decoration: If
		f.AddDecoration(n, "If", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Decoration: Init
		if n.Init != nil {
			f.AddDecoration(n, "Init", token.NoPos)
		}

		// Node: Cond
		if n.Cond != nil {
			f.ProcessNode(n.Cond)
		}

		// Decoration: Cond
		f.AddDecoration(n, "Cond", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Token: ElseTok
		if n.Else != nil {
			f.AddToken(n, token.ELSE, token.NoPos)
		}

		// Decoration: Else
		if n.Else != nil {
			f.AddDecoration(n, "Else", token.NoPos)
		}

		// Node: Else
		if n.Else != nil {
			f.ProcessNode(n.Else)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.ImportSpec:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Name
		if n.Name != nil {
			f.ProcessNode(n.Name)
		}

		// Decoration: Name
		if n.Name != nil {
			f.AddDecoration(n, "Name", token.NoPos)
		}

		// Node: Path
		if n.Path != nil {
			f.ProcessNode(n.Path)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.IncDecStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Token: Tok
		f.AddToken(n, n.Tok, n.TokPos)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.IndexExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Token: Lbrack
		f.AddToken(n, token.LBRACK, n.Lbrack)

		// Decoration: Lbrack
		f.AddDecoration(n, "Lbrack", token.NoPos)

		// Node: Index
		if n.Index != nil {
			f.ProcessNode(n.Index)
		}

		// Decoration: Index
		f.AddDecoration(n, "Index", token.NoPos)

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, n.Rbrack)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.InterfaceType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Interface
		f.AddToken(n, token.INTERFACE, n.Interface)

		// Decoration: Interface
		f.AddDecoration(n, "Interface", token.NoPos)

		// Node: Methods
		if n.Methods != nil {
			f.ProcessNode(n.Methods)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.KeyValueExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Key
		if n.Key != nil {
			f.ProcessNode(n.Key)
		}

		// Decoration: Key
		f.AddDecoration(n, "Key", token.NoPos)

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.AddDecoration(n, "Colon", token.NoPos)

		// Node: Value
		if n.Value != nil {
			f.ProcessNode(n.Value)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.LabeledStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Label
		if n.Label != nil {
			f.ProcessNode(n.Label)
		}

		// Decoration: Label
		f.AddDecoration(n, "Label", token.NoPos)

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.AddDecoration(n, "Colon", token.NoPos)

		// Node: Stmt
		if n.Stmt != nil {
			f.ProcessNode(n.Stmt)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.MapType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Map
		f.AddToken(n, token.MAP, n.Map)

		// Token: Lbrack
		f.AddToken(n, token.LBRACK, token.NoPos)

		// Decoration: Map
		f.AddDecoration(n, "Map", token.NoPos)

		// Node: Key
		if n.Key != nil {
			f.ProcessNode(n.Key)
		}

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, token.NoPos)

		// Decoration: Key
		f.AddDecoration(n, "Key", token.NoPos)

		// Node: Value
		if n.Value != nil {
			f.ProcessNode(n.Value)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.Package:

		// Map: Imports

		// Map: Files
		for _, v := range n.Files {
			f.ProcessNode(v)
		}

	case *ast.ParenExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Lparen
		f.AddToken(n, token.LPAREN, n.Lparen)

		// Decoration: Lparen
		f.AddDecoration(n, "Lparen", token.NoPos)

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Token: Rparen
		f.AddToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.RangeStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: For
		f.AddToken(n, token.FOR, n.For)

		// Decoration: For
		if n.Key != nil {
			f.AddDecoration(n, "For", token.NoPos)
		}

		// Node: Key
		if n.Key != nil {
			f.ProcessNode(n.Key)
		}

		// Token: Comma
		if n.Value != nil {
			f.AddToken(n, token.COMMA, token.NoPos)
		}

		// Decoration: Key
		if n.Key != nil {
			f.AddDecoration(n, "Key", token.NoPos)
		}

		// Node: Value
		if n.Value != nil {
			f.ProcessNode(n.Value)
		}

		// Decoration: Value
		if n.Value != nil {
			f.AddDecoration(n, "Value", token.NoPos)
		}

		// Token: Tok
		if n.Tok != token.ILLEGAL {
			f.AddToken(n, n.Tok, n.TokPos)
		}

		// Token: Range
		f.AddToken(n, token.RANGE, token.NoPos)

		// Decoration: Range
		f.AddDecoration(n, "Range", token.NoPos)

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.ReturnStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Return
		f.AddToken(n, token.RETURN, n.Return)

		// Decoration: Return
		f.AddDecoration(n, "Return", token.NoPos)

		// List: Results
		for _, v := range n.Results {
			f.ProcessNode(v)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.SelectStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Select
		f.AddToken(n, token.SELECT, n.Select)

		// Decoration: Select
		f.AddDecoration(n, "Select", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.SelectorExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Token: Period
		f.AddToken(n, token.PERIOD, token.NoPos)

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Node: Sel
		if n.Sel != nil {
			f.ProcessNode(n.Sel)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.SendStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Chan
		if n.Chan != nil {
			f.ProcessNode(n.Chan)
		}

		// Decoration: Chan
		f.AddDecoration(n, "Chan", token.NoPos)

		// Token: Arrow
		f.AddToken(n, token.ARROW, n.Arrow)

		// Decoration: Arrow
		f.AddDecoration(n, "Arrow", token.NoPos)

		// Node: Value
		if n.Value != nil {
			f.ProcessNode(n.Value)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.SliceExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Token: Lbrack
		f.AddToken(n, token.LBRACK, n.Lbrack)

		// Decoration: Lbrack
		if n.Low != nil {
			f.AddDecoration(n, "Lbrack", token.NoPos)
		}

		// Node: Low
		if n.Low != nil {
			f.ProcessNode(n.Low)
		}

		// Token: Colon1
		f.AddToken(n, token.COLON, token.NoPos)

		// Decoration: Low
		f.AddDecoration(n, "Low", token.NoPos)

		// Node: High
		if n.High != nil {
			f.ProcessNode(n.High)
		}

		// Token: Colon2
		if n.Slice3 {
			f.AddToken(n, token.COLON, token.NoPos)
		}

		// Decoration: High
		if n.High != nil {
			f.AddDecoration(n, "High", token.NoPos)
		}

		// Node: Max
		if n.Max != nil {
			f.ProcessNode(n.Max)
		}

		// Decoration: Max
		if n.Max != nil {
			f.AddDecoration(n, "Max", token.NoPos)
		}

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, n.Rbrack)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.StarExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Star
		f.AddToken(n, token.MUL, n.Star)

		// Decoration: Star
		f.AddDecoration(n, "Star", token.NoPos)

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.StructType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Struct
		f.AddToken(n, token.STRUCT, n.Struct)

		// Decoration: Struct
		f.AddDecoration(n, "Struct", token.NoPos)

		// Node: Fields
		if n.Fields != nil {
			f.ProcessNode(n.Fields)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.SwitchStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Switch
		f.AddToken(n, token.SWITCH, n.Switch)

		// Decoration: Switch
		f.AddDecoration(n, "Switch", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Decoration: Init
		if n.Init != nil {
			f.AddDecoration(n, "Init", token.NoPos)
		}

		// Node: Tag
		if n.Tag != nil {
			f.ProcessNode(n.Tag)
		}

		// Decoration: Tag
		if n.Tag != nil {
			f.AddDecoration(n, "Tag", token.NoPos)
		}

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.TypeAssertExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Token: Period
		f.AddToken(n, token.PERIOD, token.NoPos)

		// Decoration: X
		f.AddDecoration(n, "X", token.NoPos)

		// Token: Lparen
		f.AddToken(n, token.LPAREN, n.Lparen)

		// Decoration: Lparen
		f.AddDecoration(n, "Lparen", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Token: TypeToken
		if n.Type == nil {
			f.AddToken(n, token.TYPE, token.NoPos)
		}

		// Decoration: Type
		f.AddDecoration(n, "Type", token.NoPos)

		// Token: Rparen
		f.AddToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.TypeSpec:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Name
		if n.Name != nil {
			f.ProcessNode(n.Name)
		}

		// Token: Assign
		if n.Assign.IsValid() {
			f.AddToken(n, token.ASSIGN, n.Assign)
		}

		// Decoration: Name
		f.AddDecoration(n, "Name", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.TypeSwitchStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Switch
		f.AddToken(n, token.SWITCH, n.Switch)

		// Decoration: Switch
		f.AddDecoration(n, "Switch", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Decoration: Init
		if n.Init != nil {
			f.AddDecoration(n, "Init", token.NoPos)
		}

		// Node: Assign
		if n.Assign != nil {
			f.ProcessNode(n.Assign)
		}

		// Decoration: Assign
		f.AddDecoration(n, "Assign", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.UnaryExpr:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Op
		f.AddToken(n, n.Op, n.OpPos)

		// Decoration: Op
		f.AddDecoration(n, "Op", token.NoPos)

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.ValueSpec:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// List: Names
		for _, v := range n.Names {
			f.ProcessNode(v)
		}

		// Decoration: Names
		if n.Type != nil {
			f.AddDecoration(n, "Names", token.NoPos)
		}

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Token: Assign
		if n.Values != nil {
			f.AddToken(n, token.ASSIGN, token.NoPos)
		}

		// Decoration: Assign
		if n.Values != nil {
			f.AddDecoration(n, "Assign", token.NoPos)
		}

		// List: Values
		for _, v := range n.Values {
			f.ProcessNode(v)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	}
}
