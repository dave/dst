package dst

// ArrayTypeDecorations holds decorations for ArrayType:
//
// 	type R /*Start*/ [ /*Lbrack*/ 1] /*Len*/ int /*End*/
//
type ArrayTypeDecorations struct {
	Start  Decorations
	Lbrack Decorations
	Len    Decorations
	End    Decorations
	Space  SpaceType
}

// AssignStmtDecorations holds decorations for AssignStmt:
//
// 	/*Start*/
// 	i /*Lhs*/ = /*Tok*/ 1 /*End*/
//
type AssignStmtDecorations struct {
	Start Decorations
	Lhs   Decorations
	Tok   Decorations
	End   Decorations
	Space SpaceType
}

// BadDeclDecorations holds decorations for BadDecl:
//
type BadDeclDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// BadExprDecorations holds decorations for BadExpr:
//
type BadExprDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// BadStmtDecorations holds decorations for BadStmt:
//
type BadStmtDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// BasicLitDecorations holds decorations for BasicLit:
//
type BasicLitDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// BinaryExprDecorations holds decorations for BinaryExpr:
//
// 	var P = /*Start*/ 1 /*X*/ & /*Op*/ 2 /*End*/
//
type BinaryExprDecorations struct {
	Start Decorations
	X     Decorations
	Op    Decorations
	End   Decorations
	Space SpaceType
}

// BlockStmtDecorations holds decorations for BlockStmt:
//
// 	if true /*Start*/ { /*Lbrace*/
// 		i++
// 	} /*End*/
//
// 	func() /*Start*/ { /*Lbrace*/ i++ } /*End*/ ()
//
type BlockStmtDecorations struct {
	Start  Decorations
	Lbrace Decorations
	End    Decorations
	Space  SpaceType
}

// BranchStmtDecorations holds decorations for BranchStmt:
//
// 	/*Start*/
// 	goto /*Tok*/ A /*End*/
//
type BranchStmtDecorations struct {
	Start Decorations
	Tok   Decorations
	End   Decorations
	Space SpaceType
}

// CallExprDecorations holds decorations for CallExpr:
//
// 	var L = /*Start*/ C /*Fun*/ ( /*Lparen*/ 0, []int{} /*Args*/ ... /*Ellipsis*/) /*End*/
//
type CallExprDecorations struct {
	Start    Decorations
	Fun      Decorations
	Lparen   Decorations
	Args     Decorations
	Ellipsis Decorations
	End      Decorations
	Space    SpaceType
}

// CaseClauseDecorations holds decorations for CaseClause:
//
// 	switch i {
// 	/*Start*/ case /*Case*/ 1 /*List*/ : /*Colon*/
// 		i++
// 	}
//
type CaseClauseDecorations struct {
	Start Decorations
	Case  Decorations
	List  Decorations
	Colon Decorations
	Space SpaceType
}

// ChanTypeDecorations holds decorations for ChanType:
//
// 	type W /*Start*/ chan /*Begin*/ int /*End*/
//
// 	type X /*Start*/ <-chan /*Begin*/ int /*End*/
//
// 	type Y /*Start*/ chan /*Begin*/ <- /*Arrow*/ int /*End*/
//
type ChanTypeDecorations struct {
	Start Decorations
	Begin Decorations
	Arrow Decorations
	End   Decorations
	Space SpaceType
}

// CommClauseDecorations holds decorations for CommClause:
//
// 	select {
// 	/*Start*/ case /*Case*/ a := <-c /*Comm*/ : /*Colon*/
// 		print(a)
// 	}
//
type CommClauseDecorations struct {
	Start Decorations
	Case  Decorations
	Comm  Decorations
	Colon Decorations
	Space SpaceType
}

// CompositeLitDecorations holds decorations for CompositeLit:
//
// 	var D = /*Start*/ A /*Type*/ { /*Lbrace*/ A: 0} /*End*/
//
type CompositeLitDecorations struct {
	Start  Decorations
	Type   Decorations
	Lbrace Decorations
	End    Decorations
	Space  SpaceType
}

// DeclStmtDecorations holds decorations for DeclStmt:
//
type DeclStmtDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// DeferStmtDecorations holds decorations for DeferStmt:
//
// 	/*Start*/
// 	defer /*Defer*/ func() {}() /*End*/
//
type DeferStmtDecorations struct {
	Start Decorations
	Defer Decorations
	End   Decorations
	Space SpaceType
}

// EllipsisDecorations holds decorations for Ellipsis:
//
// 	func B(a /*Start*/ ... /*Ellipsis*/ int /*End*/) {}
//
type EllipsisDecorations struct {
	Start    Decorations
	Ellipsis Decorations
	End      Decorations
	Space    SpaceType
}

// EmptyStmtDecorations holds decorations for EmptyStmt:
//
type EmptyStmtDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// ExprStmtDecorations holds decorations for ExprStmt:
//
type ExprStmtDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// FieldDecorations holds decorations for Field:
//
// 	type A struct {
// 		/*Start*/ A /*Names*/ int /*Type*/ `a:"a"` /*End*/
// 	}
//
type FieldDecorations struct {
	Start Decorations
	Names Decorations
	Type  Decorations
	End   Decorations
	Space SpaceType
}

// FieldListDecorations holds decorations for FieldList:
//
// 	type A1 struct /*Start*/ { /*Opening*/
// 		a, b int
// 		c    string
// 	} /*End*/
//
type FieldListDecorations struct {
	Start   Decorations
	Opening Decorations
	End     Decorations
	Space   SpaceType
}

// FileDecorations holds decorations for File:
//
// 	/*Start*/ package /*Package*/ postests /*Name*/
//
type FileDecorations struct {
	Start   Decorations
	Package Decorations
	Name    Decorations
	Space   SpaceType
}

// ForStmtDecorations holds decorations for ForStmt:
//
// 	/*Start*/
// 	for /*For*/ {
// 		i++
// 	} /*End*/
//
// 	/*Start*/
// 	for /*For*/ i < 1 /*Cond*/ {
// 		i++
// 	} /*End*/
//
// 	/*Start*/
// 	for /*For*/ i = 0; /*Init*/ i < 10; /*Cond*/ i++ /*Post*/ {
// 		i++
// 	} /*End*/
//
type ForStmtDecorations struct {
	Start Decorations
	For   Decorations
	Init  Decorations
	Cond  Decorations
	Post  Decorations
	End   Decorations
	Space SpaceType
}

// FuncDeclDecorations holds decorations for FuncDecl:
//
// 	/*Start*/
// 	func /*Func*/ d /*Name*/ (d, e int) /*Params*/ {
// 		return
// 	} /*End*/
//
// 	/*Start*/
// 	func /*Func*/ (a *A) /*Recv*/ e /*Name*/ (d, e int) /*Params*/ {
// 		return
// 	} /*End*/
//
// 	/*Start*/
// 	func /*Func*/ (a *A) /*Recv*/ f /*Name*/ (d, e int) /*Params*/ (f, g int) /*Results*/ {
// 		return
// 	}
//
type FuncDeclDecorations struct {
	Start   Decorations
	Func    Decorations
	Recv    Decorations
	Name    Decorations
	Params  Decorations
	Results Decorations
	End     Decorations
	Space   SpaceType
}

// FuncLitDecorations holds decorations for FuncLit:
//
// 	var C = /*Start*/ func(a int, b ...int) (c int) /*Type*/ { return 0 } /*End*/
//
type FuncLitDecorations struct {
	Start Decorations
	Type  Decorations
	End   Decorations
	Space SpaceType
}

// FuncTypeDecorations holds decorations for FuncType:
//
// 	type T /*Start*/ func /*Func*/ (a int) /*Params*/ (b int) /*End*/
//
type FuncTypeDecorations struct {
	Start  Decorations
	Func   Decorations
	Params Decorations
	End    Decorations
	Space  SpaceType
}

// GenDeclDecorations holds decorations for GenDecl:
//
// 	/*Start*/
// 	const /*Tok*/ ( /*Lparen*/
// 		a, b = 1, 2
// 		c    = 3
// 	) /*End*/
//
// 	/*Start*/
// 	const /*Tok*/ d = 1 /*End*/
//
// }
//
type GenDeclDecorations struct {
	Start  Decorations
	Tok    Decorations
	Lparen Decorations
	End    Decorations
	Space  SpaceType
}

// GoStmtDecorations holds decorations for GoStmt:
//
// 	/*Start*/
// 	go /*Go*/ func() {}() /*End*/
//
type GoStmtDecorations struct {
	Start Decorations
	Go    Decorations
	End   Decorations
	Space SpaceType
}

// IdentDecorations holds decorations for Ident:
//
type IdentDecorations struct {
	Start Decorations
	End   Decorations
	Space SpaceType
}

// IfStmtDecorations holds decorations for IfStmt:
//
// 	/*Start*/
// 	if /*If*/ a := b; /*Init*/ a /*Cond*/ {
// 		i++
// 	} else /*Else*/ {
// 		i++
// 	} /*End*/
//
type IfStmtDecorations struct {
	Start Decorations
	If    Decorations
	Init  Decorations
	Cond  Decorations
	Else  Decorations
	End   Decorations
	Space SpaceType
}

// ImportSpecDecorations holds decorations for ImportSpec:
//
// 	import (
// 		/*Start*/ fmt /*Name*/ "fmt" /*End*/
// 	)
//
type ImportSpecDecorations struct {
	Start Decorations
	Name  Decorations
	End   Decorations
	Space SpaceType
}

// IncDecStmtDecorations holds decorations for IncDecStmt:
//
// 	/*Start*/
// 	i /*X*/ ++ /*End*/
//
type IncDecStmtDecorations struct {
	Start Decorations
	X     Decorations
	End   Decorations
	Space SpaceType
}

// IndexExprDecorations holds decorations for IndexExpr:
//
// 	var G = /*Start*/ []int{0} /*X*/ [ /*Lbrack*/ 0 /*Index*/] /*End*/
//
type IndexExprDecorations struct {
	Start  Decorations
	X      Decorations
	Lbrack Decorations
	Index  Decorations
	End    Decorations
	Space  SpaceType
}

// InterfaceTypeDecorations holds decorations for InterfaceType:
//
// 	type U /*Start*/ interface /*Interface*/ {
// 		A()
// 	} /*End*/
//
type InterfaceTypeDecorations struct {
	Start     Decorations
	Interface Decorations
	End       Decorations
	Space     SpaceType
}

// KeyValueExprDecorations holds decorations for KeyValueExpr:
//
// 	var Q = map[string]string{
// 		/*Start*/ "a" /*Key*/ : /*Colon*/ "a", /*End*/
// 	}
//
type KeyValueExprDecorations struct {
	Start Decorations
	Key   Decorations
	Colon Decorations
	End   Decorations
	Space SpaceType
}

// LabeledStmtDecorations holds decorations for LabeledStmt:
//
// 	/*Start*/
// 	A /*Label*/ : /*Colon*/
// 		print("Stmt") /*End*/
//
type LabeledStmtDecorations struct {
	Start Decorations
	Label Decorations
	Colon Decorations
	End   Decorations
	Space SpaceType
}

// MapTypeDecorations holds decorations for MapType:
//
// 	type V /*Start*/ map[ /*Map*/ int] /*Key*/ int /*End*/
//
type MapTypeDecorations struct {
	Start Decorations
	Map   Decorations
	Key   Decorations
	End   Decorations
	Space SpaceType
}

// PackageDecorations holds decorations for Package:
//
type PackageDecorations struct {
	Space SpaceType
}

// ParenExprDecorations holds decorations for ParenExpr:
//
// 	var E = /*Start*/ ( /*Lparen*/ 1 + 1 /*X*/) /*End*/ / 2
//
type ParenExprDecorations struct {
	Start  Decorations
	Lparen Decorations
	X      Decorations
	End    Decorations
	Space  SpaceType
}

// RangeStmtDecorations holds decorations for RangeStmt:
//
// 	/*Start*/
// 	for range /*Range*/ a /*X*/ {
// 	} /*End*/
//
// 	/*Start*/
// 	for /*For*/ k /*Key*/ := range /*Range*/ a /*X*/ {
// 		print(k)
// 	} /*End*/
//
// 	/*Start*/
// 	for /*For*/ k /*Key*/, v /*Value*/ := range /*Range*/ a /*X*/ {
// 		print(k, v)
// 	} /*End*/
//
type RangeStmtDecorations struct {
	Start Decorations
	For   Decorations
	Key   Decorations
	Value Decorations
	Range Decorations
	X     Decorations
	End   Decorations
	Space SpaceType
}

// ReturnStmtDecorations holds decorations for ReturnStmt:
//
// 	func() int {
// 		/*Start*/ return /*Return*/ 1 /*End*/
// 	}()
//
type ReturnStmtDecorations struct {
	Start  Decorations
	Return Decorations
	End    Decorations
	Space  SpaceType
}

// SelectStmtDecorations holds decorations for SelectStmt:
//
// 	/*Start*/
// 	select /*Select*/ {
// 	} /*End*/
//
type SelectStmtDecorations struct {
	Start  Decorations
	Select Decorations
	End    Decorations
	Space  SpaceType
}

// SelectorExprDecorations holds decorations for SelectorExpr:
//
// 	var F = /*Start*/ fmt. /*X*/ Sprint /*End*/ (0)
//
type SelectorExprDecorations struct {
	Start Decorations
	X     Decorations
	End   Decorations
	Space SpaceType
}

// SendStmtDecorations holds decorations for SendStmt:
//
// 	/*Start*/
// 	c /*Chan*/ <- /*Arrow*/ 0 /*End*/
//
type SendStmtDecorations struct {
	Start Decorations
	Chan  Decorations
	Arrow Decorations
	End   Decorations
	Space SpaceType
}

// SliceExprDecorations holds decorations for SliceExpr:
//
// 	var H = /*Start*/ []int{0} /*X*/ [ /*Lbrack*/ 1: /*Low*/ 2: /*High*/ 3 /*Max*/] /*End*/
//
// 	var H1 = /*Start*/ []int{0} /*X*/ [ /*Lbrack*/ 1: /*Low*/ 2 /*High*/] /*End*/
//
// 	var H2 = /*Start*/ []int{0} /*X*/ [: /*Low*/] /*End*/
//
// 	var H3 = /*Start*/ []int{0} /*X*/ [ /*Lbrack*/ 1: /*Low*/] /*End*/
//
// 	var H4 = /*Start*/ []int{0} /*X*/ [: /*Low*/ 2 /*High*/] /*End*/
//
// 	var H5 = /*Start*/ []int{0} /*X*/ [: /*Low*/ 2: /*High*/ 3 /*Max*/] /*End*/
//
type SliceExprDecorations struct {
	Start  Decorations
	X      Decorations
	Lbrack Decorations
	Low    Decorations
	High   Decorations
	Max    Decorations
	End    Decorations
	Space  SpaceType
}

// StarExprDecorations holds decorations for StarExpr:
//
// 	var N = /*Start*/ * /*Star*/ p /*End*/
//
type StarExprDecorations struct {
	Start Decorations
	Star  Decorations
	End   Decorations
	Space SpaceType
}

// StructTypeDecorations holds decorations for StructType:
//
// 	type S /*Start*/ struct /*Struct*/ {
// 		A int
// 	} /*End*/
//
type StructTypeDecorations struct {
	Start  Decorations
	Struct Decorations
	End    Decorations
	Space  SpaceType
}

// SwitchStmtDecorations holds decorations for SwitchStmt:
//
// 	/*Start*/
// 	switch /*Switch*/ i /*Tag*/ {
// 	} /*End*/
//
// 	/*Start*/
// 	switch /*Switch*/ a := i; /*Init*/ a /*Tag*/ {
// 	} /*End*/
//
type SwitchStmtDecorations struct {
	Start  Decorations
	Switch Decorations
	Init   Decorations
	Tag    Decorations
	End    Decorations
	Space  SpaceType
}

// TypeAssertExprDecorations holds decorations for TypeAssertExpr:
//
// 	var J = /*Start*/ f. /*X*/ ( /*Lparen*/ int /*Type*/) /*End*/
//
type TypeAssertExprDecorations struct {
	Start  Decorations
	X      Decorations
	Lparen Decorations
	Type   Decorations
	End    Decorations
	Space  SpaceType
}

// TypeSpecDecorations holds decorations for TypeSpec:
//
// 	type (
// 		/*Start*/ T1 /*Name*/ []int /*End*/
// 	)
//
// 	type (
// 		/*Start*/ T2 = /*Name*/ T1 /*End*/
// 	)
//
type TypeSpecDecorations struct {
	Start Decorations
	Name  Decorations
	End   Decorations
	Space SpaceType
}

// TypeSwitchStmtDecorations holds decorations for TypeSwitchStmt:
//
// 	/*Start*/
// 	switch /*Switch*/ f.(type) /*Assign*/ {
// 	} /*End*/
//
// 	/*Start*/
// 	switch /*Switch*/ g := f.(type) /*Assign*/ {
// 	case int:
// 		print(g)
// 	} /*End*/
//
// 	/*Start*/
// 	switch /*Switch*/ g := f; /*Init*/ g := g.(type) /*Assign*/ {
// 	case int:
// 		print(g)
// 	} /*End*/
//
type TypeSwitchStmtDecorations struct {
	Start  Decorations
	Switch Decorations
	Init   Decorations
	Assign Decorations
	End    Decorations
	Space  SpaceType
}

// UnaryExprDecorations holds decorations for UnaryExpr:
//
// 	var O = /*Start*/ ^ /*Op*/ 1 /*End*/
//
type UnaryExprDecorations struct {
	Start Decorations
	Op    Decorations
	End   Decorations
	Space SpaceType
}

// ValueSpecDecorations holds decorations for ValueSpec:
//
// 	var (
// 		/*Start*/ j = /*Assign*/ 1 /*End*/
// 	)
//
// 	var (
// 		/*Start*/ k, l = /*Assign*/ 1, 2 /*End*/
// 	)
//
// 	var (
// 		/*Start*/ m, n /*Names*/ int = /*Assign*/ 1, 2 /*End*/
// 	)
//
type ValueSpecDecorations struct {
	Start  Decorations
	Names  Decorations
	Assign Decorations
	End    Decorations
	Space  SpaceType
}
