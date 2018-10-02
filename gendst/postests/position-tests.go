// File
/*Start*/ package /*AfterPackage*/ postests /*AfterName*/

// ImportSpec
import (
	/*Start*/ fmt /*AfterName*/ "fmt" /*End*/
)

var a []int
var i int
var b bool
var f interface{}
var p *int
var c chan int

// Field
type A struct {
	/*Start*/ A /*AfterNames*/ int /*AfterType*/ `a:"a"` /*End*/
}

// FieldList
type A1 struct /*Start*/ { /*AfterOpening*/
	a, b int
	c    string
} /*End*/

// Ellipsis
func B(a /*Start*/ ... /*AfterEllipsis*/ int /*End*/) {}

// FuncLit
var C = /*Start*/ func(a int, b ...int) (c int) /*AfterType*/ { return 0 } /*End*/

// CompositeLit
var D = /*Start*/ A /*AfterType*/ { /*AfterLbrace*/ A: 0} /*End*/

// ParenExpr
var E = /*Start*/ ( /*AfterLparen*/ 1 + 1 /*AfterX*/) /*End*/ / 2

// SelectorExpr
var F = /*Start*/ fmt. /*AfterX*/ Sprint /*End*/ (0)

// IndexExpr
var G = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 0 /*AfterIndex*/] /*End*/

// SliceExpr(0)
var H = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/

// SliceExpr(1)
var H1 = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2 /*AfterHigh*/] /*End*/

// SliceExpr(2)
var H2 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/] /*End*/

// SliceExpr(3)
var H3 = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/] /*End*/

// SliceExpr(4)
var H4 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/ 2 /*AfterHigh*/] /*End*/

// SliceExpr(5)
var H5 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/

// TypeAssertExpr
var J = /*Start*/ f. /*AfterX*/ ( /*AfterLparen*/ int /*AfterType*/) /*End*/

// CallExpr
var L = /*Start*/ C /*AfterFun*/ ( /*AfterLparen*/ 0, []int{} /*AfterArgs*/ ... /*AfterEllipsis*/) /*End*/

// StarExpr
var N = /*Start*/ * /*AfterStar*/ p /*End*/

// UnaryExpr
var O = /*Start*/ ^ /*AfterOp*/ 1 /*End*/

// BinaryExpr
var P = /*Start*/ 1 /*AfterX*/ & /*AfterOp*/ 2 /*End*/

// KeyValueExpr
var Q = map[string]string{
	/*Start*/ "a" /*AfterKey*/ : /*AfterColon*/ "a", /*End*/
}

// ArrayType
type R /*Start*/ [ /*AfterLbrack*/ 1] /*AfterLen*/ int /*End*/

// StructType
type S /*Start*/ struct /*AfterStruct*/ {
	A int
} /*End*/

// FuncType
type T /*Start*/ func /*AfterFunc*/ (a int) /*AfterParams*/ (b int) /*End*/

// InterfaceType
type U /*Start*/ interface /*AfterInterface*/ {
	A()
} /*End*/

// MapType
type V /*Start*/ map[ /*AfterMap*/ int] /*AfterKey*/ int /*End*/

// ChanType(0)
type W /*Start*/ chan /*AfterBegin*/ int /*End*/

// ChanType(1)
type X /*Start*/ <-chan /*AfterBegin*/ int /*End*/

// ChanType(2)
type Y /*Start*/ chan /*AfterBegin*/ <- /*AfterArrow*/ int /*End*/

func Z() {
	// LabeledStmt
			/*Start*/
A /*AfterLabel*/ : /*AfterColon*/
	print("Stmt") /*End*/

	// BranchStmt
	/*Start*/
	goto /*AfterTok*/ A /*End*/

	// SendStmt
	/*Start*/
	c /*AfterChan*/ <- /*AfterArrow*/ 0 /*End*/

	// IncDecStmt
	/*Start*/
	i /*AfterX*/ ++ /*End*/

	// AssignStmt
	/*Start*/
	i /*AfterLhs*/ = /*AfterTok*/ 1 /*End*/

	// GoStmt
	/*Start*/
	go /*AfterGo*/ func() {}() /*End*/

	// DeferStmt
	/*Start*/
	defer /*AfterDefer*/ func() {}() /*End*/

	// ReturnStmt
	func() int {
		/*Start*/ return /*AfterReturn*/ 1 /*End*/
	}()

	// BlockStmt(0)
	if true /*Start*/ { /*AfterLbrace*/
		i++
	} /*End*/

	// BlockStmt(1)
	func() /*Start*/ { /*AfterLbrace*/ i++ } /*End*/ ()

	// IfStmt
	/*Start*/
	if /*AfterIf*/ a := b; /*AfterInit*/ a /*AfterCond*/ {
		i++
	} else /*AfterElse*/ {
		i++
	} /*End*/

	// CaseClause
	switch i {
	/*Start*/ case /*AfterCase*/ 1 /*AfterList*/ : /*AfterColon*/
		i++
	}

	// SwitchStmt(0)
	/*Start*/
	switch /*AfterSwitch*/ i /*AfterTag*/ {
	} /*End*/

	// SwitchStmt(1)
	/*Start*/
	switch /*AfterSwitch*/ a := i; /*AfterInit*/ a /*AfterTag*/ {
	} /*End*/

	// TypeSwitchStmt(0)
	/*Start*/
	switch /*AfterSwitch*/ f.(type) /*AfterAssign*/ {
	} /*End*/

	// TypeSwitchStmt(1)
	/*Start*/
	switch /*AfterSwitch*/ g := f.(type) /*AfterAssign*/ {
	case int:
		print(g)
	} /*End*/

	// TypeSwitchStmt(2)
	/*Start*/
	switch /*AfterSwitch*/ g := f; /*AfterInit*/ g := g.(type) /*AfterAssign*/ {
	case int:
		print(g)
	} /*End*/

	// CommClause
	select {
	/*Start*/ case /*AfterCase*/ a := <-c /*AfterComm*/ : /*AfterColon*/
		print(a)
	}

	// SelectStmt
	/*Start*/
	select /*AfterSelect*/ {
	} /*End*/

	// ForStmt(0)
	/*Start*/
	for /*AfterFor*/ {
		i++
	} /*End*/

	// ForStmt(1)
	/*Start*/
	for /*AfterFor*/ i < 1 /*AfterCond*/ {
		i++
	} /*End*/

	// ForStmt(2)
	/*Start*/
	for /*AfterFor*/ i = 0; /*AfterInit*/ i < 10; /*AfterCond*/ i++ /*AfterPost*/ {
		i++
	} /*End*/

	// RangeStmt(0)
	/*Start*/
	for range /*AfterRange*/ a /*AfterX*/ {
	} /*End*/

	// RangeStmt(1)
	/*Start*/
	for /*AfterFor*/ k /*AfterKey*/ := range /*AfterRange*/ a /*AfterX*/ {
		print(k)
	} /*End*/

	// RangeStmt(2)
	/*Start*/
	for /*AfterFor*/ k /*AfterKey*/, v /*AfterValue*/ := range /*AfterRange*/ a /*AfterX*/ {
		print(k, v)
	} /*End*/

	// ValueSpec(0)
	var (
		/*Start*/ j = /*AfterAssign*/ 1 /*End*/
	)

	// ValueSpec(1)
	var (
		/*Start*/ k, l = /*AfterAssign*/ 1, 2 /*End*/
	)

	// ValueSpec(2)
	var (
		/*Start*/ m, n /*AfterNames*/ int = /*AfterAssign*/ 1, 2 /*End*/
	)

	print(j, k, l, m, n)

	// TypeSpec(0)
	type (
		/*Start*/ T1 /*AfterName*/ []int /*End*/
	)

	// TypeSpec(1)
	type (
		/*Start*/ T2 = /*AfterName*/ T1 /*End*/
	)

	// GenDecl(0)
	/*Start*/
	const /*AfterTok*/ ( /*AfterLparen*/
		a, b = 1, 2
		c    = 3
	) /*End*/

	// GenDecl(1)
	/*Start*/
	const /*AfterTok*/ d = 1 /*End*/

}

// FuncDecl(0)
/*Start*/
func /*AfterFunc*/ d /*AfterName*/ (d, e int) /*AfterParams*/ {
	return
} /*End*/

// FuncDecl(1)
/*Start*/
func /*AfterFunc*/ (a *A) /*AfterRecv*/ e /*AfterName*/ (d, e int) /*AfterParams*/ {
	return
} /*End*/

// FuncDecl(2)
/*Start*/
func /*AfterFunc*/ (a *A) /*AfterRecv*/ f /*AfterName*/ (d, e int) /*AfterParams*/ (f, g int) /*AfterResults*/ {
	return
} /*End*/
