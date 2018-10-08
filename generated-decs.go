package dst

// ArrayTypeDecorations holds decorations for ArrayType:
//
// 	type R /*Start*/ [ /*AfterLbrack*/ 1] /*AfterLen*/ int /*End*/
//
type ArrayTypeDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterLbrack Decorations
	AfterLen    Decorations
	End         Decorations
	After       SpaceType
}

// AssignStmtDecorations holds decorations for AssignStmt:
//
// 	/*Start*/
// 	i /*AfterLhs*/ = /*AfterTok*/ 1 /*End*/
//
type AssignStmtDecorations struct {
	Before   SpaceType
	Start    Decorations
	AfterLhs Decorations
	AfterTok Decorations
	End      Decorations
	After    SpaceType
}

// BadDeclDecorations holds decorations for BadDecl:
//
type BadDeclDecorations struct {
	Before SpaceType
	After  SpaceType
}

// BadExprDecorations holds decorations for BadExpr:
//
type BadExprDecorations struct {
	Before SpaceType
	After  SpaceType
}

// BadStmtDecorations holds decorations for BadStmt:
//
type BadStmtDecorations struct {
	Before SpaceType
	After  SpaceType
}

// BasicLitDecorations holds decorations for BasicLit:
//
type BasicLitDecorations struct {
	Before SpaceType
	Start  Decorations
	End    Decorations
	After  SpaceType
}

// BinaryExprDecorations holds decorations for BinaryExpr:
//
// 	var P = /*Start*/ 1 /*AfterX*/ & /*AfterOp*/ 2 /*End*/
//
type BinaryExprDecorations struct {
	Before  SpaceType
	Start   Decorations
	AfterX  Decorations
	AfterOp Decorations
	End     Decorations
	After   SpaceType
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
	Before      SpaceType
	Start       Decorations
	AfterLbrace Decorations
	End         Decorations
	After       SpaceType
}

// BranchStmtDecorations holds decorations for BranchStmt:
//
// 	/*Start*/
// 	goto /*AfterTok*/ A /*End*/
//
type BranchStmtDecorations struct {
	Before   SpaceType
	Start    Decorations
	AfterTok Decorations
	End      Decorations
	After    SpaceType
}

// CallExprDecorations holds decorations for CallExpr:
//
// 	var L = /*Start*/ C /*AfterFun*/ ( /*AfterLparen*/ 0, []int{} /*AfterArgs*/ ... /*AfterEllipsis*/) /*End*/
//
type CallExprDecorations struct {
	Before        SpaceType
	Start         Decorations
	AfterFun      Decorations
	AfterLparen   Decorations
	AfterArgs     Decorations
	AfterEllipsis Decorations
	End           Decorations
	After         SpaceType
}

// CaseClauseDecorations holds decorations for CaseClause:
//
// 	switch i {
// 	/*Start*/ case /*AfterCase*/ 1 /*AfterList*/ : /*AfterColon*/
// 		i++
// 	}
//
type CaseClauseDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterCase  Decorations
	AfterList  Decorations
	AfterColon Decorations
	After      SpaceType
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
	Before     SpaceType
	Start      Decorations
	AfterBegin Decorations
	AfterArrow Decorations
	End        Decorations
	After      SpaceType
}

// CommClauseDecorations holds decorations for CommClause:
//
// 	select {
// 	/*Start*/ case /*AfterCase*/ a := <-c /*AfterComm*/ : /*AfterColon*/
// 		print(a)
// 	}
//
type CommClauseDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterCase  Decorations
	AfterComm  Decorations
	AfterColon Decorations
	After      SpaceType
}

// CompositeLitDecorations holds decorations for CompositeLit:
//
// 	var D = /*Start*/ A /*AfterType*/ { /*AfterLbrace*/ A: 0} /*End*/
//
type CompositeLitDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterType   Decorations
	AfterLbrace Decorations
	End         Decorations
	After       SpaceType
}

// DeclStmtDecorations holds decorations for DeclStmt:
//
type DeclStmtDecorations struct {
	Before SpaceType
	Start  Decorations
	End    Decorations
	After  SpaceType
}

// DeferStmtDecorations holds decorations for DeferStmt:
//
// 	/*Start*/
// 	defer /*AfterDefer*/ func() {}() /*End*/
//
type DeferStmtDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterDefer Decorations
	End        Decorations
	After      SpaceType
}

// EllipsisDecorations holds decorations for Ellipsis:
//
// 	func B(a /*Start*/ ... /*AfterEllipsis*/ int /*End*/) {}
//
type EllipsisDecorations struct {
	Before        SpaceType
	Start         Decorations
	AfterEllipsis Decorations
	End           Decorations
	After         SpaceType
}

// EmptyStmtDecorations holds decorations for EmptyStmt:
//
type EmptyStmtDecorations struct {
	Before SpaceType
	After  SpaceType
}

// ExprStmtDecorations holds decorations for ExprStmt:
//
type ExprStmtDecorations struct {
	Before SpaceType
	Start  Decorations
	End    Decorations
	After  SpaceType
}

// FieldDecorations holds decorations for Field:
//
// 	type A struct {
// 		/*Start*/ A /*AfterNames*/ int /*AfterType*/ `a:"a"` /*End*/
// 	}
//
type FieldDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterNames Decorations
	AfterType  Decorations
	End        Decorations
	After      SpaceType
}

// FieldListDecorations holds decorations for FieldList:
//
// 	type A1 struct /*Start*/ { /*AfterOpening*/
// 		a, b int
// 		c    string
// 	} /*End*/
//
type FieldListDecorations struct {
	Before       SpaceType
	Start        Decorations
	AfterOpening Decorations
	End          Decorations
	After        SpaceType
}

// FileDecorations holds decorations for File:
//
// 	/*Start*/ package /*AfterPackage*/ postests /*AfterName*/
//
type FileDecorations struct {
	Before       SpaceType
	Start        Decorations
	AfterPackage Decorations
	AfterName    Decorations
	After        SpaceType
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
	Before    SpaceType
	Start     Decorations
	AfterFor  Decorations
	AfterInit Decorations
	AfterCond Decorations
	AfterPost Decorations
	End       Decorations
	After     SpaceType
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
	Before       SpaceType
	Start        Decorations
	AfterFunc    Decorations
	AfterRecv    Decorations
	AfterName    Decorations
	AfterParams  Decorations
	AfterResults Decorations
	End          Decorations
	After        SpaceType
}

// FuncLitDecorations holds decorations for FuncLit:
//
// 	var C = /*Start*/ func(a int, b ...int) (c int) /*AfterType*/ { return 0 } /*End*/
//
type FuncLitDecorations struct {
	Before    SpaceType
	Start     Decorations
	AfterType Decorations
	End       Decorations
	After     SpaceType
}

// FuncTypeDecorations holds decorations for FuncType:
//
// 	type T /*Start*/ func /*AfterFunc*/ (a int) /*AfterParams*/ (b int) /*End*/
//
type FuncTypeDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterFunc   Decorations
	AfterParams Decorations
	End         Decorations
	After       SpaceType
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
	Before      SpaceType
	Start       Decorations
	AfterTok    Decorations
	AfterLparen Decorations
	End         Decorations
	After       SpaceType
}

// GoStmtDecorations holds decorations for GoStmt:
//
// 	/*Start*/
// 	go /*AfterGo*/ func() {}() /*End*/
//
type GoStmtDecorations struct {
	Before  SpaceType
	Start   Decorations
	AfterGo Decorations
	End     Decorations
	After   SpaceType
}

// IdentDecorations holds decorations for Ident:
//
type IdentDecorations struct {
	Before SpaceType
	Start  Decorations
	End    Decorations
	After  SpaceType
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
	Before    SpaceType
	Start     Decorations
	AfterIf   Decorations
	AfterInit Decorations
	AfterCond Decorations
	AfterElse Decorations
	End       Decorations
	After     SpaceType
}

// ImportSpecDecorations holds decorations for ImportSpec:
//
// 	import (
// 		/*Start*/ fmt /*AfterName*/ "fmt" /*End*/
// 	)
//
type ImportSpecDecorations struct {
	Before    SpaceType
	Start     Decorations
	AfterName Decorations
	End       Decorations
	After     SpaceType
}

// IncDecStmtDecorations holds decorations for IncDecStmt:
//
// 	/*Start*/
// 	i /*AfterX*/ ++ /*End*/
//
type IncDecStmtDecorations struct {
	Before SpaceType
	Start  Decorations
	AfterX Decorations
	End    Decorations
	After  SpaceType
}

// IndexExprDecorations holds decorations for IndexExpr:
//
// 	var G = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 0 /*AfterIndex*/] /*End*/
//
type IndexExprDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterX      Decorations
	AfterLbrack Decorations
	AfterIndex  Decorations
	End         Decorations
	After       SpaceType
}

// InterfaceTypeDecorations holds decorations for InterfaceType:
//
// 	type U /*Start*/ interface /*AfterInterface*/ {
// 		A()
// 	} /*End*/
//
type InterfaceTypeDecorations struct {
	Before         SpaceType
	Start          Decorations
	AfterInterface Decorations
	End            Decorations
	After          SpaceType
}

// KeyValueExprDecorations holds decorations for KeyValueExpr:
//
// 	var Q = map[string]string{
// 		/*Start*/ "a" /*AfterKey*/ : /*AfterColon*/ "a", /*End*/
// 	}
//
type KeyValueExprDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterKey   Decorations
	AfterColon Decorations
	End        Decorations
	After      SpaceType
}

// LabeledStmtDecorations holds decorations for LabeledStmt:
//
// 	/*Start*/
// 	A /*AfterLabel*/ : /*AfterColon*/
// 		print("Stmt") /*End*/
//
type LabeledStmtDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterLabel Decorations
	AfterColon Decorations
	End        Decorations
	After      SpaceType
}

// MapTypeDecorations holds decorations for MapType:
//
// 	type V /*Start*/ map[ /*AfterMap*/ int] /*AfterKey*/ int /*End*/
//
type MapTypeDecorations struct {
	Before   SpaceType
	Start    Decorations
	AfterMap Decorations
	AfterKey Decorations
	End      Decorations
	After    SpaceType
}

// PackageDecorations holds decorations for Package:
//
type PackageDecorations struct {
	Before SpaceType
	After  SpaceType
}

// ParenExprDecorations holds decorations for ParenExpr:
//
// 	var E = /*Start*/ ( /*AfterLparen*/ 1 + 1 /*AfterX*/) /*End*/ / 2
//
type ParenExprDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterLparen Decorations
	AfterX      Decorations
	End         Decorations
	After       SpaceType
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
	Before     SpaceType
	Start      Decorations
	AfterFor   Decorations
	AfterKey   Decorations
	AfterValue Decorations
	AfterRange Decorations
	AfterX     Decorations
	End        Decorations
	After      SpaceType
}

// ReturnStmtDecorations holds decorations for ReturnStmt:
//
// 	func() int {
// 		/*Start*/ return /*AfterReturn*/ 1 /*End*/
// 	}()
//
type ReturnStmtDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterReturn Decorations
	End         Decorations
	After       SpaceType
}

// SelectStmtDecorations holds decorations for SelectStmt:
//
// 	/*Start*/
// 	select /*AfterSelect*/ {
// 	} /*End*/
//
type SelectStmtDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterSelect Decorations
	End         Decorations
	After       SpaceType
}

// SelectorExprDecorations holds decorations for SelectorExpr:
//
// 	var F = /*Start*/ fmt. /*AfterX*/ Sprint /*End*/ (0)
//
type SelectorExprDecorations struct {
	Before SpaceType
	Start  Decorations
	AfterX Decorations
	End    Decorations
	After  SpaceType
}

// SendStmtDecorations holds decorations for SendStmt:
//
// 	/*Start*/
// 	c /*AfterChan*/ <- /*AfterArrow*/ 0 /*End*/
//
type SendStmtDecorations struct {
	Before     SpaceType
	Start      Decorations
	AfterChan  Decorations
	AfterArrow Decorations
	End        Decorations
	After      SpaceType
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
	Before      SpaceType
	Start       Decorations
	AfterX      Decorations
	AfterLbrack Decorations
	AfterLow    Decorations
	AfterHigh   Decorations
	AfterMax    Decorations
	End         Decorations
	After       SpaceType
}

// StarExprDecorations holds decorations for StarExpr:
//
// 	var N = /*Start*/ * /*AfterStar*/ p /*End*/
//
type StarExprDecorations struct {
	Before    SpaceType
	Start     Decorations
	AfterStar Decorations
	End       Decorations
	After     SpaceType
}

// StructTypeDecorations holds decorations for StructType:
//
// 	type S /*Start*/ struct /*AfterStruct*/ {
// 		A int
// 	} /*End*/
//
type StructTypeDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterStruct Decorations
	End         Decorations
	After       SpaceType
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
	Before      SpaceType
	Start       Decorations
	AfterSwitch Decorations
	AfterInit   Decorations
	AfterTag    Decorations
	End         Decorations
	After       SpaceType
}

// TypeAssertExprDecorations holds decorations for TypeAssertExpr:
//
// 	var J = /*Start*/ f. /*AfterX*/ ( /*AfterLparen*/ int /*AfterType*/) /*End*/
//
type TypeAssertExprDecorations struct {
	Before      SpaceType
	Start       Decorations
	AfterX      Decorations
	AfterLparen Decorations
	AfterType   Decorations
	End         Decorations
	After       SpaceType
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
	Before    SpaceType
	Start     Decorations
	AfterName Decorations
	End       Decorations
	After     SpaceType
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
	Before      SpaceType
	Start       Decorations
	AfterSwitch Decorations
	AfterInit   Decorations
	AfterAssign Decorations
	End         Decorations
	After       SpaceType
}

// UnaryExprDecorations holds decorations for UnaryExpr:
//
// 	var O = /*Start*/ ^ /*AfterOp*/ 1 /*End*/
//
type UnaryExprDecorations struct {
	Before  SpaceType
	Start   Decorations
	AfterOp Decorations
	End     Decorations
	After   SpaceType
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
	Before      SpaceType
	Start       Decorations
	AfterNames  Decorations
	AfterAssign Decorations
	End         Decorations
	After       SpaceType
}
