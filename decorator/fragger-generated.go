package decorator

import (
	"go/ast"
	"go/token"
)

func (f *fragger) processNode(n ast.Node) {
	if n.Pos().IsValid() {
		f.cursor = int(n.Pos())
	}
	switch n := n.(type) {
	case *ast.ArrayType:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Lbrack
		f.addToken(n, token.LBRACK, n.Lbrack)

		// Decoration: Lbrack
		f.addDecoration(n, "Lbrack", token.NoPos)

		// Node: Len
		if n.Len != nil {
			f.processNode(n.Len)
		}

		// Token: Rbrack
		f.addToken(n, token.RBRACK, token.NoPos)

		// Decoration: Len
		f.addDecoration(n, "Len", token.NoPos)

		// Node: Elt
		if n.Elt != nil {
			f.processNode(n.Elt)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.AssignStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// List: Lhs
		for _, v := range n.Lhs {
			f.processNode(v)
		}

		// Decoration: Lhs
		f.addDecoration(n, "Lhs", token.NoPos)

		// Token: Tok
		f.addToken(n, n.Tok, n.TokPos)

		// Decoration: Tok
		f.addDecoration(n, "Tok", token.NoPos)

		// List: Rhs
		for _, v := range n.Rhs {
			f.processNode(v)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BadDecl:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BadExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BadStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BasicLit:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// String: Value
		f.addString(n, n.Value, n.ValuePos)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BinaryExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Token: Op
		f.addToken(n, n.Op, n.OpPos)

		// Decoration: Op
		f.addDecoration(n, "Op", token.NoPos)

		// Node: Y
		if n.Y != nil {
			f.processNode(n.Y)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BlockStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Lbrace
		f.addToken(n, token.LBRACE, n.Lbrace)

		// Decoration: Lbrace
		f.addDecoration(n, "Lbrace", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.processNode(v)
		}

		// Token: Rbrace
		f.addToken(n, token.RBRACE, n.Rbrace)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.BranchStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Tok
		f.addToken(n, n.Tok, n.TokPos)

		// Decoration: Tok
		if n.Label != nil {
			f.addDecoration(n, "Tok", token.NoPos)
		}

		// Node: Label
		if n.Label != nil {
			f.processNode(n.Label)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.CallExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Fun
		if n.Fun != nil {
			f.processNode(n.Fun)
		}

		// Decoration: Fun
		f.addDecoration(n, "Fun", token.NoPos)

		// Token: Lparen
		f.addToken(n, token.LPAREN, n.Lparen)

		// Decoration: Lparen
		f.addDecoration(n, "Lparen", token.NoPos)

		// List: Args
		for _, v := range n.Args {
			f.processNode(v)
		}

		// Decoration: Args
		f.addDecoration(n, "Args", token.NoPos)

		// Token: Ellipsis
		if n.Ellipsis.IsValid() {
			f.addToken(n, token.ELLIPSIS, n.Ellipsis)
		}

		// Decoration: Ellipsis
		if n.Ellipsis.IsValid() {
			f.addDecoration(n, "Ellipsis", token.NoPos)
		}

		// Token: Rparen
		f.addToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.CaseClause:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Case
		f.addToken(n, func() token.Token {
			if n.List == nil {
				return token.DEFAULT
			} else {
				return token.CASE
			}
		}(), n.Case)

		// Decoration: Case
		f.addDecoration(n, "Case", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.processNode(v)
		}

		// Decoration: List
		if n.List != nil {
			f.addDecoration(n, "List", token.NoPos)
		}

		// Token: Colon
		f.addToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.addDecoration(n, "Colon", token.NoPos)

		// List: Body
		for _, v := range n.Body {
			f.processNode(v)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.ChanType:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Begin
		f.addToken(n, func() token.Token {
			if n.Dir == ast.RECV {
				return token.ARROW
			} else {
				return token.CHAN
			}
		}(), n.Begin)

		// Token: Chan
		if n.Dir == ast.RECV {
			f.addToken(n, token.CHAN, token.NoPos)
		}

		// Decoration: Begin
		f.addDecoration(n, "Begin", token.NoPos)

		// Token: Arrow
		if n.Dir == ast.SEND {
			f.addToken(n, token.ARROW, n.Arrow)
		}

		// Decoration: Arrow
		if n.Dir == ast.SEND {
			f.addDecoration(n, "Arrow", token.NoPos)
		}

		// Node: Value
		if n.Value != nil {
			f.processNode(n.Value)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.CommClause:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Case
		f.addToken(n, func() token.Token {
			if n.Comm == nil {
				return token.DEFAULT
			} else {
				return token.CASE
			}
		}(), n.Case)

		// Decoration: Case
		f.addDecoration(n, "Case", token.NoPos)

		// Node: Comm
		if n.Comm != nil {
			f.processNode(n.Comm)
		}

		// Decoration: Comm
		if n.Comm != nil {
			f.addDecoration(n, "Comm", token.NoPos)
		}

		// Token: Colon
		f.addToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.addDecoration(n, "Colon", token.NoPos)

		// List: Body
		for _, v := range n.Body {
			f.processNode(v)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.CompositeLit:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Type
		if n.Type != nil {
			f.processNode(n.Type)
		}

		// Decoration: Type
		if n.Type != nil {
			f.addDecoration(n, "Type", token.NoPos)
		}

		// Token: Lbrace
		f.addToken(n, token.LBRACE, n.Lbrace)

		// Decoration: Lbrace
		f.addDecoration(n, "Lbrace", token.NoPos)

		// List: Elts
		for _, v := range n.Elts {
			f.processNode(v)
		}

		// Token: Rbrace
		f.addToken(n, token.RBRACE, n.Rbrace)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.DeclStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Decl
		if n.Decl != nil {
			f.processNode(n.Decl)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.DeferStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Defer
		f.addToken(n, token.DEFER, n.Defer)

		// Decoration: Defer
		f.addDecoration(n, "Defer", token.NoPos)

		// Node: Call
		if n.Call != nil {
			f.processNode(n.Call)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.Ellipsis:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Ellipsis
		f.addToken(n, token.ELLIPSIS, n.Ellipsis)

		// Decoration: Ellipsis
		if n.Elt != nil {
			f.addDecoration(n, "Ellipsis", token.NoPos)
		}

		// Node: Elt
		if n.Elt != nil {
			f.processNode(n.Elt)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.EmptyStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Semicolon
		if !n.Implicit {
			f.addToken(n, token.ARROW, n.Semicolon)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.ExprStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.Field:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// List: Names
		for _, v := range n.Names {
			f.processNode(v)
		}

		// Decoration: Names
		f.addDecoration(n, "Names", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.processNode(n.Type)
		}

		// Decoration: Type
		if n.Tag != nil {
			f.addDecoration(n, "Type", token.NoPos)
		}

		// Node: Tag
		if n.Tag != nil {
			f.processNode(n.Tag)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.FieldList:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Opening
		if n.Opening.IsValid() {
			f.addToken(n, token.LPAREN, n.Opening)
		}

		// Decoration: Opening
		f.addDecoration(n, "Opening", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.processNode(v)
		}

		// Token: Closing
		if n.Closing.IsValid() {
			f.addToken(n, token.RPAREN, n.Closing)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.File:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Package
		f.addToken(n, token.PACKAGE, n.Package)

		// Decoration: Package
		f.addDecoration(n, "Package", token.NoPos)

		// Node: Name
		if n.Name != nil {
			f.processNode(n.Name)
		}

		// Decoration: Name
		f.addDecoration(n, "Name", token.NoPos)

		// List: Decls
		for _, v := range n.Decls {
			f.processNode(v)
		}

	case *ast.ForStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: For
		f.addToken(n, token.FOR, n.For)

		// Decoration: For
		f.addDecoration(n, "For", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.processNode(n.Init)
		}

		// Token: InitSemicolon
		if n.Init != nil {
			f.addToken(n, token.SEMICOLON, token.NoPos)
		}

		// Decoration: Init
		if n.Init != nil {
			f.addDecoration(n, "Init", token.NoPos)
		}

		// Node: Cond
		if n.Cond != nil {
			f.processNode(n.Cond)
		}

		// Token: CondSemicolon
		if n.Post != nil {
			f.addToken(n, token.SEMICOLON, token.NoPos)
		}

		// Decoration: Cond
		if n.Cond != nil {
			f.addDecoration(n, "Cond", token.NoPos)
		}

		// Node: Post
		if n.Post != nil {
			f.processNode(n.Post)
		}

		// Decoration: Post
		if n.Post != nil {
			f.addDecoration(n, "Post", token.NoPos)
		}

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.FuncDecl:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Func
		if true {
			f.addToken(n, token.FUNC, n.Type.Func)
		}

		// Decoration: Func
		f.addDecoration(n, "Func", token.NoPos)

		// Node: Recv
		if n.Recv != nil {
			f.processNode(n.Recv)
		}

		// Decoration: Recv
		if n.Recv != nil {
			f.addDecoration(n, "Recv", token.NoPos)
		}

		// Node: Name
		if n.Name != nil {
			f.processNode(n.Name)
		}

		// Decoration: Name
		f.addDecoration(n, "Name", token.NoPos)

		// Node: Params
		if n.Type.Params != nil {
			f.processNode(n.Type.Params)
		}

		// Decoration: Params
		f.addDecoration(n, "Params", token.NoPos)

		// Node: Results
		if n.Type.Results != nil {
			f.processNode(n.Type.Results)
		}

		// Decoration: Results
		if n.Type.Results != nil {
			f.addDecoration(n, "Results", token.NoPos)
		}

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.FuncLit:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Type
		if n.Type != nil {
			f.processNode(n.Type)
		}

		// Decoration: Type
		f.addDecoration(n, "Type", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.FuncType:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Func
		if n.Func.IsValid() {
			f.addToken(n, token.FUNC, n.Func)
		}

		// Decoration: Func
		if n.Func.IsValid() {
			f.addDecoration(n, "Func", token.NoPos)
		}

		// Node: Params
		if n.Params != nil {
			f.processNode(n.Params)
		}

		// Decoration: Params
		if n.Results != nil {
			f.addDecoration(n, "Params", token.NoPos)
		}

		// Node: Results
		if n.Results != nil {
			f.processNode(n.Results)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.GenDecl:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Tok
		f.addToken(n, n.Tok, n.TokPos)

		// Decoration: Tok
		f.addDecoration(n, "Tok", token.NoPos)

		// Token: Lparen
		if n.Lparen.IsValid() {
			f.addToken(n, token.LPAREN, n.Lparen)
		}

		// Decoration: Lparen
		if n.Lparen.IsValid() {
			f.addDecoration(n, "Lparen", token.NoPos)
		}

		// List: Specs
		for _, v := range n.Specs {
			f.processNode(v)
		}

		// Token: Rparen
		if n.Rparen.IsValid() {
			f.addToken(n, token.RPAREN, n.Rparen)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.GoStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Go
		f.addToken(n, token.GO, n.Go)

		// Decoration: Go
		f.addDecoration(n, "Go", token.NoPos)

		// Node: Call
		if n.Call != nil {
			f.processNode(n.Call)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.Ident:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// String: Name
		f.addString(n, n.Name, n.NamePos)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.IfStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: If
		f.addToken(n, token.IF, n.If)

		// Decoration: If
		f.addDecoration(n, "If", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.processNode(n.Init)
		}

		// Decoration: Init
		if n.Init != nil {
			f.addDecoration(n, "Init", token.NoPos)
		}

		// Node: Cond
		if n.Cond != nil {
			f.processNode(n.Cond)
		}

		// Decoration: Cond
		f.addDecoration(n, "Cond", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Token: ElseTok
		if n.Else != nil {
			f.addToken(n, token.ELSE, token.NoPos)
		}

		// Decoration: Else
		if n.Else != nil {
			f.addDecoration(n, "Else", token.NoPos)
		}

		// Node: Else
		if n.Else != nil {
			f.processNode(n.Else)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.ImportSpec:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Name
		if n.Name != nil {
			f.processNode(n.Name)
		}

		// Decoration: Name
		if n.Name != nil {
			f.addDecoration(n, "Name", token.NoPos)
		}

		// Node: Path
		if n.Path != nil {
			f.processNode(n.Path)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.IncDecStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Token: Tok
		f.addToken(n, n.Tok, n.TokPos)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.IndexExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Token: Lbrack
		f.addToken(n, token.LBRACK, n.Lbrack)

		// Decoration: Lbrack
		f.addDecoration(n, "Lbrack", token.NoPos)

		// Node: Index
		if n.Index != nil {
			f.processNode(n.Index)
		}

		// Decoration: Index
		f.addDecoration(n, "Index", token.NoPos)

		// Token: Rbrack
		f.addToken(n, token.RBRACK, n.Rbrack)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.InterfaceType:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Interface
		f.addToken(n, token.INTERFACE, n.Interface)

		// Decoration: Interface
		f.addDecoration(n, "Interface", token.NoPos)

		// Node: Methods
		if n.Methods != nil {
			f.processNode(n.Methods)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.KeyValueExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Key
		if n.Key != nil {
			f.processNode(n.Key)
		}

		// Decoration: Key
		f.addDecoration(n, "Key", token.NoPos)

		// Token: Colon
		f.addToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.addDecoration(n, "Colon", token.NoPos)

		// Node: Value
		if n.Value != nil {
			f.processNode(n.Value)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.LabeledStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Label
		if n.Label != nil {
			f.processNode(n.Label)
		}

		// Decoration: Label
		f.addDecoration(n, "Label", token.NoPos)

		// Token: Colon
		f.addToken(n, token.COLON, n.Colon)

		// Decoration: Colon
		f.addDecoration(n, "Colon", token.NoPos)

		// Node: Stmt
		if n.Stmt != nil {
			f.processNode(n.Stmt)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.MapType:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Map
		f.addToken(n, token.MAP, n.Map)

		// Token: Lbrack
		f.addToken(n, token.LBRACK, token.NoPos)

		// Decoration: Map
		f.addDecoration(n, "Map", token.NoPos)

		// Node: Key
		if n.Key != nil {
			f.processNode(n.Key)
		}

		// Token: Rbrack
		f.addToken(n, token.RBRACK, token.NoPos)

		// Decoration: Key
		f.addDecoration(n, "Key", token.NoPos)

		// Node: Value
		if n.Value != nil {
			f.processNode(n.Value)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.Package:

		// Map: Imports

		// Map: Files
		for _, v := range n.Files {
			f.processNode(v)
		}

	case *ast.ParenExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Lparen
		f.addToken(n, token.LPAREN, n.Lparen)

		// Decoration: Lparen
		f.addDecoration(n, "Lparen", token.NoPos)

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Token: Rparen
		f.addToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.RangeStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: For
		f.addToken(n, token.FOR, n.For)

		// Decoration: For
		if n.Key != nil {
			f.addDecoration(n, "For", token.NoPos)
		}

		// Node: Key
		if n.Key != nil {
			f.processNode(n.Key)
		}

		// Token: Comma
		if n.Value != nil {
			f.addToken(n, token.COMMA, token.NoPos)
		}

		// Decoration: Key
		if n.Key != nil {
			f.addDecoration(n, "Key", token.NoPos)
		}

		// Node: Value
		if n.Value != nil {
			f.processNode(n.Value)
		}

		// Decoration: Value
		if n.Value != nil {
			f.addDecoration(n, "Value", token.NoPos)
		}

		// Token: Tok
		if n.Tok != token.ILLEGAL {
			f.addToken(n, n.Tok, n.TokPos)
		}

		// Token: Range
		f.addToken(n, token.RANGE, token.NoPos)

		// Decoration: Range
		f.addDecoration(n, "Range", token.NoPos)

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.ReturnStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Return
		f.addToken(n, token.RETURN, n.Return)

		// Decoration: Return
		f.addDecoration(n, "Return", token.NoPos)

		// List: Results
		for _, v := range n.Results {
			f.processNode(v)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.SelectStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Select
		f.addToken(n, token.SELECT, n.Select)

		// Decoration: Select
		f.addDecoration(n, "Select", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.SelectorExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Token: Period
		f.addToken(n, token.PERIOD, token.NoPos)

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Node: Sel
		if n.Sel != nil {
			f.processNode(n.Sel)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.SendStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Chan
		if n.Chan != nil {
			f.processNode(n.Chan)
		}

		// Decoration: Chan
		f.addDecoration(n, "Chan", token.NoPos)

		// Token: Arrow
		f.addToken(n, token.ARROW, n.Arrow)

		// Decoration: Arrow
		f.addDecoration(n, "Arrow", token.NoPos)

		// Node: Value
		if n.Value != nil {
			f.processNode(n.Value)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.SliceExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Token: Lbrack
		f.addToken(n, token.LBRACK, n.Lbrack)

		// Decoration: Lbrack
		if n.Low != nil {
			f.addDecoration(n, "Lbrack", token.NoPos)
		}

		// Node: Low
		if n.Low != nil {
			f.processNode(n.Low)
		}

		// Token: Colon1
		f.addToken(n, token.COLON, token.NoPos)

		// Decoration: Low
		f.addDecoration(n, "Low", token.NoPos)

		// Node: High
		if n.High != nil {
			f.processNode(n.High)
		}

		// Token: Colon2
		if n.Slice3 {
			f.addToken(n, token.COLON, token.NoPos)
		}

		// Decoration: High
		if n.High != nil {
			f.addDecoration(n, "High", token.NoPos)
		}

		// Node: Max
		if n.Max != nil {
			f.processNode(n.Max)
		}

		// Decoration: Max
		if n.Max != nil {
			f.addDecoration(n, "Max", token.NoPos)
		}

		// Token: Rbrack
		f.addToken(n, token.RBRACK, n.Rbrack)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.StarExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Star
		f.addToken(n, token.MUL, n.Star)

		// Decoration: Star
		f.addDecoration(n, "Star", token.NoPos)

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.StructType:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Struct
		f.addToken(n, token.STRUCT, n.Struct)

		// Decoration: Struct
		f.addDecoration(n, "Struct", token.NoPos)

		// Node: Fields
		if n.Fields != nil {
			f.processNode(n.Fields)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.SwitchStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Switch
		f.addToken(n, token.SWITCH, n.Switch)

		// Decoration: Switch
		f.addDecoration(n, "Switch", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.processNode(n.Init)
		}

		// Decoration: Init
		if n.Init != nil {
			f.addDecoration(n, "Init", token.NoPos)
		}

		// Node: Tag
		if n.Tag != nil {
			f.processNode(n.Tag)
		}

		// Decoration: Tag
		if n.Tag != nil {
			f.addDecoration(n, "Tag", token.NoPos)
		}

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.TypeAssertExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Token: Period
		f.addToken(n, token.PERIOD, token.NoPos)

		// Decoration: X
		f.addDecoration(n, "X", token.NoPos)

		// Token: Lparen
		f.addToken(n, token.LPAREN, n.Lparen)

		// Decoration: Lparen
		f.addDecoration(n, "Lparen", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.processNode(n.Type)
		}

		// Token: TypeToken
		if n.Type == nil {
			f.addToken(n, token.TYPE, token.NoPos)
		}

		// Decoration: Type
		f.addDecoration(n, "Type", token.NoPos)

		// Token: Rparen
		f.addToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.TypeSpec:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Node: Name
		if n.Name != nil {
			f.processNode(n.Name)
		}

		// Token: Assign
		if n.Assign.IsValid() {
			f.addToken(n, token.ASSIGN, n.Assign)
		}

		// Decoration: Name
		f.addDecoration(n, "Name", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.processNode(n.Type)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.TypeSwitchStmt:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Switch
		f.addToken(n, token.SWITCH, n.Switch)

		// Decoration: Switch
		f.addDecoration(n, "Switch", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.processNode(n.Init)
		}

		// Decoration: Init
		if n.Init != nil {
			f.addDecoration(n, "Init", token.NoPos)
		}

		// Node: Assign
		if n.Assign != nil {
			f.processNode(n.Assign)
		}

		// Decoration: Assign
		f.addDecoration(n, "Assign", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.processNode(n.Body)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.UnaryExpr:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// Token: Op
		f.addToken(n, n.Op, n.OpPos)

		// Decoration: Op
		f.addDecoration(n, "Op", token.NoPos)

		// Node: X
		if n.X != nil {
			f.processNode(n.X)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	case *ast.ValueSpec:

		// Decoration: Start
		f.addDecoration(n, "Start", n.Pos())

		// List: Names
		for _, v := range n.Names {
			f.processNode(v)
		}

		// Decoration: Names
		if n.Type != nil {
			f.addDecoration(n, "Names", token.NoPos)
		}

		// Node: Type
		if n.Type != nil {
			f.processNode(n.Type)
		}

		// Token: Assign
		if n.Values != nil {
			f.addToken(n, token.ASSIGN, token.NoPos)
		}

		// Decoration: Assign
		if n.Values != nil {
			f.addDecoration(n, "Assign", token.NoPos)
		}

		// List: Values
		for _, v := range n.Values {
			f.processNode(v)
		}

		// Decoration: End
		f.addDecoration(n, "End", n.End())

	}
}
