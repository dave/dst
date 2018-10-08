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

		// Decoration: AfterLbrack
		f.AddDecoration(n, "AfterLbrack", token.NoPos)

		// Node: Len
		if n.Len != nil {
			f.ProcessNode(n.Len)
		}

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, token.NoPos)

		// Decoration: AfterLen
		f.AddDecoration(n, "AfterLen", token.NoPos)

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

		// Decoration: AfterLhs
		f.AddDecoration(n, "AfterLhs", token.NoPos)

		// Token: Tok
		f.AddToken(n, n.Tok, n.TokPos)

		// Decoration: AfterTok
		f.AddDecoration(n, "AfterTok", token.NoPos)

		// List: Rhs
		for _, v := range n.Rhs {
			f.ProcessNode(v)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.BadDecl:

	case *ast.BadExpr:

	case *ast.BadStmt:

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

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

		// Token: Op
		f.AddToken(n, n.Op, n.OpPos)

		// Decoration: AfterOp
		f.AddDecoration(n, "AfterOp", token.NoPos)

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

		// Decoration: AfterLbrace
		f.AddDecoration(n, "AfterLbrace", token.NoPos)

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

		// Decoration: AfterTok
		if n.Label != nil {
			f.AddDecoration(n, "AfterTok", token.NoPos)
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

		// Decoration: AfterFun
		f.AddDecoration(n, "AfterFun", token.NoPos)

		// Token: Lparen
		f.AddToken(n, token.LPAREN, n.Lparen)

		// Decoration: AfterLparen
		f.AddDecoration(n, "AfterLparen", token.NoPos)

		// List: Args
		for _, v := range n.Args {
			f.ProcessNode(v)
		}

		// Decoration: AfterArgs
		f.AddDecoration(n, "AfterArgs", token.NoPos)

		// Token: Ellipsis
		if n.Ellipsis.IsValid() {
			f.AddToken(n, token.ELLIPSIS, n.Ellipsis)
		}

		// Decoration: AfterEllipsis
		if n.Ellipsis.IsValid() {
			f.AddDecoration(n, "AfterEllipsis", token.NoPos)
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

		// Decoration: AfterCase
		f.AddDecoration(n, "AfterCase", token.NoPos)

		// List: List
		for _, v := range n.List {
			f.ProcessNode(v)
		}

		// Decoration: AfterList
		if n.List != nil {
			f.AddDecoration(n, "AfterList", token.NoPos)
		}

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: AfterColon
		f.AddDecoration(n, "AfterColon", token.NoPos)

		// List: Body
		for _, v := range n.Body {
			f.ProcessNode(v)
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

		// Decoration: AfterBegin
		f.AddDecoration(n, "AfterBegin", token.NoPos)

		// Token: Arrow
		if n.Dir == ast.SEND {
			f.AddToken(n, token.ARROW, n.Arrow)
		}

		// Decoration: AfterArrow
		if n.Dir == ast.SEND {
			f.AddDecoration(n, "AfterArrow", token.NoPos)
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

		// Decoration: AfterCase
		f.AddDecoration(n, "AfterCase", token.NoPos)

		// Node: Comm
		if n.Comm != nil {
			f.ProcessNode(n.Comm)
		}

		// Decoration: AfterComm
		if n.Comm != nil {
			f.AddDecoration(n, "AfterComm", token.NoPos)
		}

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: AfterColon
		f.AddDecoration(n, "AfterColon", token.NoPos)

		// List: Body
		for _, v := range n.Body {
			f.ProcessNode(v)
		}

	case *ast.CompositeLit:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Decoration: AfterType
		if n.Type != nil {
			f.AddDecoration(n, "AfterType", token.NoPos)
		}

		// Token: Lbrace
		f.AddToken(n, token.LBRACE, n.Lbrace)

		// Decoration: AfterLbrace
		f.AddDecoration(n, "AfterLbrace", token.NoPos)

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

		// Decoration: AfterDefer
		f.AddDecoration(n, "AfterDefer", token.NoPos)

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

		// Decoration: AfterEllipsis
		if n.Elt != nil {
			f.AddDecoration(n, "AfterEllipsis", token.NoPos)
		}

		// Node: Elt
		if n.Elt != nil {
			f.ProcessNode(n.Elt)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.EmptyStmt:

		// Token: Semicolon
		if !n.Implicit {
			f.AddToken(n, token.ARROW, n.Semicolon)
		}

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

		// Decoration: AfterNames
		f.AddDecoration(n, "AfterNames", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Decoration: AfterType
		if n.Tag != nil {
			f.AddDecoration(n, "AfterType", token.NoPos)
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

		// Decoration: AfterOpening
		f.AddDecoration(n, "AfterOpening", token.NoPos)

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

		// Decoration: AfterPackage
		f.AddDecoration(n, "AfterPackage", token.NoPos)

		// Node: Name
		if n.Name != nil {
			f.ProcessNode(n.Name)
		}

		// Decoration: AfterName
		f.AddDecoration(n, "AfterName", token.NoPos)

		// List: Decls
		for _, v := range n.Decls {
			f.ProcessNode(v)
		}

	case *ast.ForStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: For
		f.AddToken(n, token.FOR, n.For)

		// Decoration: AfterFor
		f.AddDecoration(n, "AfterFor", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Token: InitSemicolon
		if n.Init != nil {
			f.AddToken(n, token.SEMICOLON, token.NoPos)
		}

		// Decoration: AfterInit
		if n.Init != nil {
			f.AddDecoration(n, "AfterInit", token.NoPos)
		}

		// Node: Cond
		if n.Cond != nil {
			f.ProcessNode(n.Cond)
		}

		// Token: CondSemicolon
		if n.Post != nil {
			f.AddToken(n, token.SEMICOLON, token.NoPos)
		}

		// Decoration: AfterCond
		if n.Cond != nil {
			f.AddDecoration(n, "AfterCond", token.NoPos)
		}

		// Node: Post
		if n.Post != nil {
			f.ProcessNode(n.Post)
		}

		// Decoration: AfterPost
		if n.Post != nil {
			f.AddDecoration(n, "AfterPost", token.NoPos)
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

		// Decoration: AfterFunc
		f.AddDecoration(n, "AfterFunc", token.NoPos)

		// Node: Recv
		if n.Recv != nil {
			f.ProcessNode(n.Recv)
		}

		// Decoration: AfterRecv
		if n.Recv != nil {
			f.AddDecoration(n, "AfterRecv", token.NoPos)
		}

		// Node: Name
		if n.Name != nil {
			f.ProcessNode(n.Name)
		}

		// Decoration: AfterName
		f.AddDecoration(n, "AfterName", token.NoPos)

		// Node: Params
		if n.Type.Params != nil {
			f.ProcessNode(n.Type.Params)
		}

		// Decoration: AfterParams
		f.AddDecoration(n, "AfterParams", token.NoPos)

		// Node: Results
		if n.Type.Results != nil {
			f.ProcessNode(n.Type.Results)
		}

		// Decoration: AfterResults
		if n.Type.Results != nil {
			f.AddDecoration(n, "AfterResults", token.NoPos)
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

		// Decoration: AfterType
		f.AddDecoration(n, "AfterType", token.NoPos)

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

		// Decoration: AfterFunc
		if n.Func.IsValid() {
			f.AddDecoration(n, "AfterFunc", token.NoPos)
		}

		// Node: Params
		if n.Params != nil {
			f.ProcessNode(n.Params)
		}

		// Decoration: AfterParams
		if n.Results != nil {
			f.AddDecoration(n, "AfterParams", token.NoPos)
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

		// Decoration: AfterTok
		f.AddDecoration(n, "AfterTok", token.NoPos)

		// Token: Lparen
		if n.Lparen.IsValid() {
			f.AddToken(n, token.LPAREN, n.Lparen)
		}

		// Decoration: AfterLparen
		if n.Lparen.IsValid() {
			f.AddDecoration(n, "AfterLparen", token.NoPos)
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

		// Decoration: AfterGo
		f.AddDecoration(n, "AfterGo", token.NoPos)

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

		// Decoration: AfterIf
		f.AddDecoration(n, "AfterIf", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Decoration: AfterInit
		if n.Init != nil {
			f.AddDecoration(n, "AfterInit", token.NoPos)
		}

		// Node: Cond
		if n.Cond != nil {
			f.ProcessNode(n.Cond)
		}

		// Decoration: AfterCond
		f.AddDecoration(n, "AfterCond", token.NoPos)

		// Node: Body
		if n.Body != nil {
			f.ProcessNode(n.Body)
		}

		// Token: ElseTok
		if n.Else != nil {
			f.AddToken(n, token.ELSE, token.NoPos)
		}

		// Decoration: AfterElse
		if n.Else != nil {
			f.AddDecoration(n, "AfterElse", token.NoPos)
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

		// Decoration: AfterName
		if n.Name != nil {
			f.AddDecoration(n, "AfterName", token.NoPos)
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

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

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

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

		// Token: Lbrack
		f.AddToken(n, token.LBRACK, n.Lbrack)

		// Decoration: AfterLbrack
		f.AddDecoration(n, "AfterLbrack", token.NoPos)

		// Node: Index
		if n.Index != nil {
			f.ProcessNode(n.Index)
		}

		// Decoration: AfterIndex
		f.AddDecoration(n, "AfterIndex", token.NoPos)

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, n.Rbrack)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.InterfaceType:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: Interface
		f.AddToken(n, token.INTERFACE, n.Interface)

		// Decoration: AfterInterface
		f.AddDecoration(n, "AfterInterface", token.NoPos)

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

		// Decoration: AfterKey
		f.AddDecoration(n, "AfterKey", token.NoPos)

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: AfterColon
		f.AddDecoration(n, "AfterColon", token.NoPos)

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

		// Decoration: AfterLabel
		f.AddDecoration(n, "AfterLabel", token.NoPos)

		// Token: Colon
		f.AddToken(n, token.COLON, n.Colon)

		// Decoration: AfterColon
		f.AddDecoration(n, "AfterColon", token.NoPos)

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

		// Decoration: AfterMap
		f.AddDecoration(n, "AfterMap", token.NoPos)

		// Node: Key
		if n.Key != nil {
			f.ProcessNode(n.Key)
		}

		// Token: Rbrack
		f.AddToken(n, token.RBRACK, token.NoPos)

		// Decoration: AfterKey
		f.AddDecoration(n, "AfterKey", token.NoPos)

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

		// Decoration: AfterLparen
		f.AddDecoration(n, "AfterLparen", token.NoPos)

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

		// Token: Rparen
		f.AddToken(n, token.RPAREN, n.Rparen)

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	case *ast.RangeStmt:

		// Decoration: Start
		f.AddDecoration(n, "Start", n.Pos())

		// Token: For
		f.AddToken(n, token.FOR, n.For)

		// Decoration: AfterFor
		if n.Key != nil {
			f.AddDecoration(n, "AfterFor", token.NoPos)
		}

		// Node: Key
		if n.Key != nil {
			f.ProcessNode(n.Key)
		}

		// Token: Comma
		if n.Value != nil {
			f.AddToken(n, token.COMMA, token.NoPos)
		}

		// Decoration: AfterKey
		if n.Key != nil {
			f.AddDecoration(n, "AfterKey", token.NoPos)
		}

		// Node: Value
		if n.Value != nil {
			f.ProcessNode(n.Value)
		}

		// Decoration: AfterValue
		if n.Value != nil {
			f.AddDecoration(n, "AfterValue", token.NoPos)
		}

		// Token: Tok
		if n.Tok != token.ILLEGAL {
			f.AddToken(n, n.Tok, n.TokPos)
		}

		// Token: Range
		f.AddToken(n, token.RANGE, token.NoPos)

		// Decoration: AfterRange
		f.AddDecoration(n, "AfterRange", token.NoPos)

		// Node: X
		if n.X != nil {
			f.ProcessNode(n.X)
		}

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

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

		// Decoration: AfterReturn
		f.AddDecoration(n, "AfterReturn", token.NoPos)

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

		// Decoration: AfterSelect
		f.AddDecoration(n, "AfterSelect", token.NoPos)

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

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

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

		// Decoration: AfterChan
		f.AddDecoration(n, "AfterChan", token.NoPos)

		// Token: Arrow
		f.AddToken(n, token.ARROW, n.Arrow)

		// Decoration: AfterArrow
		f.AddDecoration(n, "AfterArrow", token.NoPos)

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

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

		// Token: Lbrack
		f.AddToken(n, token.LBRACK, n.Lbrack)

		// Decoration: AfterLbrack
		if n.Low != nil {
			f.AddDecoration(n, "AfterLbrack", token.NoPos)
		}

		// Node: Low
		if n.Low != nil {
			f.ProcessNode(n.Low)
		}

		// Token: Colon1
		f.AddToken(n, token.COLON, token.NoPos)

		// Decoration: AfterLow
		f.AddDecoration(n, "AfterLow", token.NoPos)

		// Node: High
		if n.High != nil {
			f.ProcessNode(n.High)
		}

		// Token: Colon2
		if n.Slice3 {
			f.AddToken(n, token.COLON, token.NoPos)
		}

		// Decoration: AfterHigh
		if n.High != nil {
			f.AddDecoration(n, "AfterHigh", token.NoPos)
		}

		// Node: Max
		if n.Max != nil {
			f.ProcessNode(n.Max)
		}

		// Decoration: AfterMax
		if n.Max != nil {
			f.AddDecoration(n, "AfterMax", token.NoPos)
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

		// Decoration: AfterStar
		f.AddDecoration(n, "AfterStar", token.NoPos)

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

		// Decoration: AfterStruct
		f.AddDecoration(n, "AfterStruct", token.NoPos)

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

		// Decoration: AfterSwitch
		f.AddDecoration(n, "AfterSwitch", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Decoration: AfterInit
		if n.Init != nil {
			f.AddDecoration(n, "AfterInit", token.NoPos)
		}

		// Node: Tag
		if n.Tag != nil {
			f.ProcessNode(n.Tag)
		}

		// Decoration: AfterTag
		if n.Tag != nil {
			f.AddDecoration(n, "AfterTag", token.NoPos)
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

		// Decoration: AfterX
		f.AddDecoration(n, "AfterX", token.NoPos)

		// Token: Lparen
		f.AddToken(n, token.LPAREN, n.Lparen)

		// Decoration: AfterLparen
		f.AddDecoration(n, "AfterLparen", token.NoPos)

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Token: TypeToken
		if n.Type == nil {
			f.AddToken(n, token.TYPE, token.NoPos)
		}

		// Decoration: AfterType
		f.AddDecoration(n, "AfterType", token.NoPos)

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

		// Decoration: AfterName
		f.AddDecoration(n, "AfterName", token.NoPos)

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

		// Decoration: AfterSwitch
		f.AddDecoration(n, "AfterSwitch", token.NoPos)

		// Node: Init
		if n.Init != nil {
			f.ProcessNode(n.Init)
		}

		// Decoration: AfterInit
		if n.Init != nil {
			f.AddDecoration(n, "AfterInit", token.NoPos)
		}

		// Node: Assign
		if n.Assign != nil {
			f.ProcessNode(n.Assign)
		}

		// Decoration: AfterAssign
		f.AddDecoration(n, "AfterAssign", token.NoPos)

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

		// Decoration: AfterOp
		f.AddDecoration(n, "AfterOp", token.NoPos)

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

		// Decoration: AfterNames
		if n.Type != nil {
			f.AddDecoration(n, "AfterNames", token.NoPos)
		}

		// Node: Type
		if n.Type != nil {
			f.ProcessNode(n.Type)
		}

		// Token: Assign
		if n.Values != nil {
			f.AddToken(n, token.ASSIGN, token.NoPos)
		}

		// Decoration: AfterAssign
		if n.Values != nil {
			f.AddDecoration(n, "AfterAssign", token.NoPos)
		}

		// List: Values
		for _, v := range n.Values {
			f.ProcessNode(v)
		}

		// Decoration: End
		f.AddDecoration(n, "End", n.End())

	}
}
