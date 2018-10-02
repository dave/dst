package dst

// ArrayTypeDecorations holds decorations for ArrayType:
//
// 	type R /*Start*/ [ /*AfterLbrack*/ 1] /*AfterLen*/ int /*End*/
//
type ArrayTypeDecorations struct {
	Start       Decorations
	AfterLbrack Decorations
	AfterLen    Decorations
	End         Decorations
}

// AssignStmtDecorations holds decorations for AssignStmt:
//
// 	/*Start*/
// 	i /*AfterLhs*/ = /*AfterTok*/ 1 /*End*/
//
type AssignStmtDecorations struct {
	Start    Decorations
	AfterLhs Decorations
	AfterTok Decorations
	End      Decorations
}

// BadDeclDecorations holds decorations for BadDecl:
//
type BadDeclDecorations struct{}

// BadExprDecorations holds decorations for BadExpr:
//
type BadExprDecorations struct{}

// BadStmtDecorations holds decorations for BadStmt:
//
type BadStmtDecorations struct{}

// BasicLitDecorations holds decorations for BasicLit:
//
type BasicLitDecorations struct {
	Start Decorations
	End   Decorations
}

// BinaryExprDecorations holds decorations for BinaryExpr:
//
// 	var P = /*Start*/ 1 /*AfterX*/ & /*AfterOp*/ 2 /*End*/
//
type BinaryExprDecorations struct {
	Start   Decorations
	AfterX  Decorations
	AfterOp Decorations
	End     Decorations
}

// BlockStmtDecorations holds decorations for BlockStmt:
//
// 	if true /*Start*/ { /*AfterLbrace*/
// 		i++
// 	} /*End*/
//
// 	func() /*Start*/ { /*AfterLbrace*/ i++ } /*End*/ ()
//
type BlockStmtDecorations struct {
	Start       Decorations
	AfterLbrace Decorations
	End         Decorations
}

// BranchStmtDecorations holds decorations for BranchStmt:
//
// 	/*Start*/
// 	goto /*AfterTok*/ A /*End*/
//
type BranchStmtDecorations struct {
	Start    Decorations
	AfterTok Decorations
	End      Decorations
}

// CallExprDecorations holds decorations for CallExpr:
//
// 	var L = /*Start*/ C /*AfterFun*/ ( /*AfterLparen*/ 0, []int{} /*AfterArgs*/ ... /*AfterEllipsis*/) /*End*/
//
type CallExprDecorations struct {
	Start         Decorations
	AfterFun      Decorations
	AfterLparen   Decorations
	AfterArgs     Decorations
	AfterEllipsis Decorations
	End           Decorations
}

// CaseClauseDecorations holds decorations for CaseClause:
//
// 	switch i {
// 	/*Start*/ case /*AfterCase*/ 1 /*AfterList*/ : /*AfterColon*/
// 		i++
// 	}
//
type CaseClauseDecorations struct {
	Start      Decorations
	AfterCase  Decorations
	AfterList  Decorations
	AfterColon Decorations
}

// ChanTypeDecorations holds decorations for ChanType:
//
// 	type W /*Start*/ chan /*AfterBegin*/ int /*End*/
//
// 	type X /*Start*/ <-chan /*AfterBegin*/ int /*End*/
//
// 	type Y /*Start*/ chan /*AfterBegin*/ <- /*AfterArrow*/ int /*End*/
//
type ChanTypeDecorations struct {
	Start      Decorations
	AfterBegin Decorations
	AfterArrow Decorations
	End        Decorations
}

// CommClauseDecorations holds decorations for CommClause:
//
// 	select {
// 	/*Start*/ case /*AfterCase*/ a := <-c /*AfterComm*/ : /*AfterColon*/
// 		print(a)
// 	}
//
type CommClauseDecorations struct {
	Start      Decorations
	AfterCase  Decorations
	AfterComm  Decorations
	AfterColon Decorations
}

// CompositeLitDecorations holds decorations for CompositeLit:
//
// 	var D = /*Start*/ A /*AfterType*/ { /*AfterLbrace*/ A: 0} /*End*/
//
type CompositeLitDecorations struct {
	Start       Decorations
	AfterType   Decorations
	AfterLbrace Decorations
	End         Decorations
}

// DeclStmtDecorations holds decorations for DeclStmt:
//
type DeclStmtDecorations struct {
	Start Decorations
	End   Decorations
}

// DeferStmtDecorations holds decorations for DeferStmt:
//
// 	/*Start*/
// 	defer /*AfterDefer*/ func() {}() /*End*/
//
type DeferStmtDecorations struct {
	Start      Decorations
	AfterDefer Decorations
	End        Decorations
}

// EllipsisDecorations holds decorations for Ellipsis:
//
// 	func B(a /*Start*/ ... /*AfterEllipsis*/ int /*End*/) {}
//
type EllipsisDecorations struct {
	Start         Decorations
	AfterEllipsis Decorations
	End           Decorations
}

// EmptyStmtDecorations holds decorations for EmptyStmt:
//
type EmptyStmtDecorations struct{}

// ExprStmtDecorations holds decorations for ExprStmt:
//
type ExprStmtDecorations struct {
	Start Decorations
	End   Decorations
}

// FieldDecorations holds decorations for Field:
//
// 	type A struct {
// 		/*Start*/ A /*AfterNames*/ int /*AfterType*/ `a:"a"` /*End*/
// 	}
//
type FieldDecorations struct {
	Start      Decorations
	AfterNames Decorations
	AfterType  Decorations
	End        Decorations
}

// FieldListDecorations holds decorations for FieldList:
//
// 	type A1 struct /*Start*/ { /*AfterOpening*/
// 		a, b int
// 		c    string
// 	} /*End*/
//
type FieldListDecorations struct {
	Start        Decorations
	AfterOpening Decorations
	End          Decorations
}

// FileDecorations holds decorations for File:
//
// 	/*Start*/ package /*AfterPackage*/ postests /*AfterName*/
//
type FileDecorations struct {
	Start        Decorations
	AfterPackage Decorations
	AfterName    Decorations
}

// ForStmtDecorations holds decorations for ForStmt:
//
// 	/*Start*/
// 	for /*AfterFor*/ {
// 		i++
// 	} /*End*/
//
// 	/*Start*/
// 	for /*AfterFor*/ i < 1 /*AfterCond*/ {
// 		i++
// 	} /*End*/
//
// 	/*Start*/
// 	for /*AfterFor*/ i = 0; /*AfterInit*/ i < 10; /*AfterCond*/ i++ /*AfterPost*/ {
// 		i++
// 	} /*End*/
//
type ForStmtDecorations struct {
	Start     Decorations
	AfterFor  Decorations
	AfterInit Decorations
	AfterCond Decorations
	AfterPost Decorations
	End       Decorations
}

// FuncDeclDecorations holds decorations for FuncDecl:
//
// 	/*Start*/
// 	func /*AfterFunc*/ d /*AfterName*/ (d, e int) /*AfterParams*/ {
// 		return
// 	} /*End*/
//
// 	/*Start*/
// 	func /*AfterFunc*/ (a *A) /*AfterRecv*/ e /*AfterName*/ (d, e int) /*AfterParams*/ {
// 		return
// 	} /*End*/
//
// 	/*Start*/
// 	func /*AfterFunc*/ (a *A) /*AfterRecv*/ f /*AfterName*/ (d, e int) /*AfterParams*/ (f, g int) /*AfterResults*/ {
// 		return
// 	}
//
type FuncDeclDecorations struct {
	Start        Decorations
	AfterFunc    Decorations
	AfterRecv    Decorations
	AfterName    Decorations
	AfterParams  Decorations
	AfterResults Decorations
	End          Decorations
}

// FuncLitDecorations holds decorations for FuncLit:
//
// 	var C = /*Start*/ func(a int, b ...int) (c int) /*AfterType*/ { return 0 } /*End*/
//
type FuncLitDecorations struct {
	Start     Decorations
	AfterType Decorations
	End       Decorations
}

// FuncTypeDecorations holds decorations for FuncType:
//
// 	type T /*Start*/ func /*AfterFunc*/ (a int) /*AfterParams*/ (b int) /*End*/
//
type FuncTypeDecorations struct {
	Start       Decorations
	AfterFunc   Decorations
	AfterParams Decorations
	End         Decorations
}

// GenDeclDecorations holds decorations for GenDecl:
//
// 	/*Start*/
// 	const /*AfterTok*/ ( /*AfterLparen*/
// 		a, b = 1, 2
// 		c    = 3
// 	) /*End*/
//
// 	/*Start*/
// 	const /*AfterTok*/ d = 1 /*End*/
//
// }
//
type GenDeclDecorations struct {
	Start       Decorations
	AfterTok    Decorations
	AfterLparen Decorations
	End         Decorations
}

// GoStmtDecorations holds decorations for GoStmt:
//
// 	/*Start*/
// 	go /*AfterGo*/ func() {}() /*End*/
//
type GoStmtDecorations struct {
	Start   Decorations
	AfterGo Decorations
	End     Decorations
}

// IdentDecorations holds decorations for Ident:
//
type IdentDecorations struct {
	Start Decorations
	End   Decorations
}

// IfStmtDecorations holds decorations for IfStmt:
//
// 	/*Start*/
// 	if /*AfterIf*/ a := b; /*AfterInit*/ a /*AfterCond*/ {
// 		i++
// 	} else /*AfterElse*/ {
// 		i++
// 	} /*End*/
//
type IfStmtDecorations struct {
	Start     Decorations
	AfterIf   Decorations
	AfterInit Decorations
	AfterCond Decorations
	AfterElse Decorations
	End       Decorations
}

// ImportSpecDecorations holds decorations for ImportSpec:
//
// 	import (
// 		/*Start*/ fmt /*AfterName*/ "fmt" /*End*/
// 	)
//
type ImportSpecDecorations struct {
	Start     Decorations
	AfterName Decorations
	End       Decorations
}

// IncDecStmtDecorations holds decorations for IncDecStmt:
//
// 	/*Start*/
// 	i /*AfterX*/ ++ /*End*/
//
type IncDecStmtDecorations struct {
	Start  Decorations
	AfterX Decorations
	End    Decorations
}

// IndexExprDecorations holds decorations for IndexExpr:
//
// 	var G = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 0 /*AfterIndex*/] /*End*/
//
type IndexExprDecorations struct {
	Start       Decorations
	AfterX      Decorations
	AfterLbrack Decorations
	AfterIndex  Decorations
	End         Decorations
}

// InterfaceTypeDecorations holds decorations for InterfaceType:
//
// 	type U /*Start*/ interface /*AfterInterface*/ {
// 		A()
// 	} /*End*/
//
type InterfaceTypeDecorations struct {
	Start          Decorations
	AfterInterface Decorations
	End            Decorations
}

// KeyValueExprDecorations holds decorations for KeyValueExpr:
//
// 	var Q = map[string]string{
// 		/*Start*/ "a" /*AfterKey*/ : /*AfterColon*/ "a", /*End*/
// 	}
//
type KeyValueExprDecorations struct {
	Start      Decorations
	AfterKey   Decorations
	AfterColon Decorations
	End        Decorations
}

// LabeledStmtDecorations holds decorations for LabeledStmt:
//
// 	/*Start*/
// 	A /*AfterLabel*/ : /*AfterColon*/
// 		print("Stmt") /*End*/
//
type LabeledStmtDecorations struct {
	Start      Decorations
	AfterLabel Decorations
	AfterColon Decorations
	End        Decorations
}

// MapTypeDecorations holds decorations for MapType:
//
// 	type V /*Start*/ map[ /*AfterMap*/ int] /*AfterKey*/ int /*End*/
//
type MapTypeDecorations struct {
	Start    Decorations
	AfterMap Decorations
	AfterKey Decorations
	End      Decorations
}

// ParenExprDecorations holds decorations for ParenExpr:
//
// 	var E = /*Start*/ ( /*AfterLparen*/ 1 + 1 /*AfterX*/) /*End*/ / 2
//
type ParenExprDecorations struct {
	Start       Decorations
	AfterLparen Decorations
	AfterX      Decorations
	End         Decorations
}

// RangeStmtDecorations holds decorations for RangeStmt:
//
// 	/*Start*/
// 	for range /*AfterRange*/ a /*AfterX*/ {
// 	} /*End*/
//
// 	/*Start*/
// 	for /*AfterFor*/ k /*AfterKey*/ := range /*AfterRange*/ a /*AfterX*/ {
// 		print(k)
// 	} /*End*/
//
// 	/*Start*/
// 	for /*AfterFor*/ k /*AfterKey*/, v /*AfterValue*/ := range /*AfterRange*/ a /*AfterX*/ {
// 		print(k, v)
// 	} /*End*/
//
type RangeStmtDecorations struct {
	Start      Decorations
	AfterFor   Decorations
	AfterKey   Decorations
	AfterValue Decorations
	AfterRange Decorations
	AfterX     Decorations
	End        Decorations
}

// ReturnStmtDecorations holds decorations for ReturnStmt:
//
// 	func() int {
// 		/*Start*/ return /*AfterReturn*/ 1 /*End*/
// 	}()
//
type ReturnStmtDecorations struct {
	Start       Decorations
	AfterReturn Decorations
	End         Decorations
}

// SelectStmtDecorations holds decorations for SelectStmt:
//
// 	/*Start*/
// 	select /*AfterSelect*/ {
// 	} /*End*/
//
type SelectStmtDecorations struct {
	Start       Decorations
	AfterSelect Decorations
	End         Decorations
}

// SelectorExprDecorations holds decorations for SelectorExpr:
//
// 	var F = /*Start*/ fmt. /*AfterX*/ Sprint /*End*/ (0)
//
type SelectorExprDecorations struct {
	Start  Decorations
	AfterX Decorations
	End    Decorations
}

// SendStmtDecorations holds decorations for SendStmt:
//
// 	/*Start*/
// 	c /*AfterChan*/ <- /*AfterArrow*/ 0 /*End*/
//
type SendStmtDecorations struct {
	Start      Decorations
	AfterChan  Decorations
	AfterArrow Decorations
	End        Decorations
}

// SliceExprDecorations holds decorations for SliceExpr:
//
// 	var H = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/
//
// 	var H1 = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2 /*AfterHigh*/] /*End*/
//
// 	var H2 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/] /*End*/
//
// 	var H3 = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/] /*End*/
//
// 	var H4 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/ 2 /*AfterHigh*/] /*End*/
//
// 	var H5 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/
//
type SliceExprDecorations struct {
	Start       Decorations
	AfterX      Decorations
	AfterLbrack Decorations
	AfterLow    Decorations
	AfterHigh   Decorations
	AfterMax    Decorations
	End         Decorations
}

// StarExprDecorations holds decorations for StarExpr:
//
// 	var N = /*Start*/ * /*AfterStar*/ p /*End*/
//
type StarExprDecorations struct {
	Start     Decorations
	AfterStar Decorations
	End       Decorations
}

// StructTypeDecorations holds decorations for StructType:
//
// 	type S /*Start*/ struct /*AfterStruct*/ {
// 		A int
// 	} /*End*/
//
type StructTypeDecorations struct {
	Start       Decorations
	AfterStruct Decorations
	End         Decorations
}

// SwitchStmtDecorations holds decorations for SwitchStmt:
//
// 	/*Start*/
// 	switch /*AfterSwitch*/ i /*AfterTag*/ {
// 	} /*End*/
//
// 	/*Start*/
// 	switch /*AfterSwitch*/ a := i; /*AfterInit*/ a /*AfterTag*/ {
// 	} /*End*/
//
type SwitchStmtDecorations struct {
	Start       Decorations
	AfterSwitch Decorations
	AfterInit   Decorations
	AfterTag    Decorations
	End         Decorations
}

// TypeAssertExprDecorations holds decorations for TypeAssertExpr:
//
// 	var J = /*Start*/ f. /*AfterX*/ ( /*AfterLparen*/ int /*AfterType*/) /*End*/
//
type TypeAssertExprDecorations struct {
	Start       Decorations
	AfterX      Decorations
	AfterLparen Decorations
	AfterType   Decorations
	End         Decorations
}

// TypeSpecDecorations holds decorations for TypeSpec:
//
// 	type (
// 		/*Start*/ T1 /*AfterName*/ []int /*End*/
// 	)
//
// 	type (
// 		/*Start*/ T2 = /*AfterName*/ T1 /*End*/
// 	)
//
type TypeSpecDecorations struct {
	Start     Decorations
	AfterName Decorations
	End       Decorations
}

// TypeSwitchStmtDecorations holds decorations for TypeSwitchStmt:
//
// 	/*Start*/
// 	switch /*AfterSwitch*/ f.(type) /*AfterAssign*/ {
// 	} /*End*/
//
// 	/*Start*/
// 	switch /*AfterSwitch*/ g := f.(type) /*AfterAssign*/ {
// 	case int:
// 		print(g)
// 	} /*End*/
//
// 	/*Start*/
// 	switch /*AfterSwitch*/ g := f; /*AfterInit*/ g := g.(type) /*AfterAssign*/ {
// 	case int:
// 		print(g)
// 	} /*End*/
//
type TypeSwitchStmtDecorations struct {
	Start       Decorations
	AfterSwitch Decorations
	AfterInit   Decorations
	AfterAssign Decorations
	End         Decorations
}

// UnaryExprDecorations holds decorations for UnaryExpr:
//
// 	var O = /*Start*/ ^ /*AfterOp*/ 1 /*End*/
//
type UnaryExprDecorations struct {
	Start   Decorations
	AfterOp Decorations
	End     Decorations
}

// ValueSpecDecorations holds decorations for ValueSpec:
//
// 	var (
// 		/*Start*/ j = /*AfterAssign*/ 1 /*End*/
// 	)
//
// 	var (
// 		/*Start*/ k, l = /*AfterAssign*/ 1, 2 /*End*/
// 	)
//
// 	var (
// 		/*Start*/ m, n /*AfterNames*/ int = /*AfterAssign*/ 1, 2 /*End*/
// 	)
//
type ValueSpecDecorations struct {
	Start       Decorations
	AfterNames  Decorations
	AfterAssign Decorations
	End         Decorations
}
