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
	/*Start*/ A /*AfterName*/ int /*AfterType*/ `a:"a"` /*End*/
}

// Ellipsis
func B(a /*Start*/ ... /*AfterEllipsis*/ int /*End*/) {}

// FuncLit
var C = /*Start*/ func(a int, b ...int) (c int) /*AfterType*/ { return 0 } /*End*/

// CompositeLit
var D = /*Start*/ A /*AfterType*/ { /*AfterLbrace*/ A: 0 /*AfterElts*/} /*End*/

// ParenExpr
var E = /*Start*/ ( /*AfterLparen*/ 1 + 1 /*AfterX*/) /*End*/ / 2

// SelectorExpr (no comment possible before period).
var F = /*Start*/ fmt. /*AfterX*/ Sprint /*End*/ (0)

// IndexExpr
var G = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 0 /*AfterIndex*/] /*End*/

// SliceExpr (no comment possible before colons)
var H = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/

var H1 = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2 /*AfterHigh*/] /*End*/

var H2 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/] /*End*/

var H3 = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/] /*End*/

var H4 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/ 2 /*AfterHigh*/] /*End*/

var H5 = /*Start*/ []int{0} /*AfterX*/ [: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/

// TypeAssertExpr (no comment possible before period)
var J = /*Start*/ f. /*AfterX*/ ( /*AfterLparen*/ int /*AfterType*/) /*End*/

// CallExpr
var L = /*Start*/ C /*AfterFun*/ ( /*AfterLparen*/ 0, []int{} /*AfterArgs*/ ... /*AfterEllipsis*/) /*End*/

// StarExpr
var N = /*Start*/ * /*AfterStar*/ p /*End*/

// UnaryExpr
var O = /*Start*/ ^ /*AfterOp*/ 1 /*End*/

// BinaryExpr
var P = /*Start*/ 1 /*AfterX*/ & /*AfterOp*/ 2 /*End*/

// KeyValueExpr (no comment possible before comma)
var Q = map[string]string{
	/*Start*/ "a" /*AfterKey*/ : /*AfterColon*/ "a", /*End*/
}

// ArrayType (no comment possible before closing bracket)
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

// MapType (no comment possible before opening/closing bracket)
type V /*Start*/ map[ /*AfterMap*/ int] /*AfterKey*/ int /*End*/

// ChanType (no comment is possible between "<-" and "chan" in SEND variation
type W /*Start*/ chan /*AfterBegin*/ int /*End*/

type X /*Start*/ <-chan /*AfterBegin*/ int /*End*/

type Y /*Start*/ chan /*AfterBegin*/ <- /*AfterArrow*/ int /*End*/

func Z() {
	// LabeledStmt (gofmt has a bug moving the start comment to a weird place)
			/*Start*/
A /*AfterLabel*/ : /*AfterColon*/
	print("Stmt") /*End*/

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

	// BranchStmt
	/*Start*/
	goto /*AfterTok*/ A /*End*/

	// BlockStmt
	if true /*Start*/ { /*AfterLbrace*/
		i++
	} /*End*/

	func() /*Start*/ { /*AfterLbrace*/ i++ } /*End*/ ()

	// IfStmt (no comment possible before semicolon)
	/*Start*/
	if /*AfterIf*/ a := b; /*AfterInit*/ a /*AfterCond*/ {
		i++
	} else /*AfterElse*/ {
		i++
	} /*End*/

	// CaseClause
	switch i {
	/*Start*/ case /*AfterCase*/ 1 /*AfterList*/ : /*AfterColon*/
		i++ /*End(?)*/
	}

	// SwitchStmt (no comment possible before semicolon)
	/*Start*/
	switch /*AfterSwitch*/ i /*AfterTag*/ {
	} /*End*/

	/*Start*/
	switch /*AfterSwitch*/ a := i; /*AfterInit*/ a /*AfterTag*/ {
	} /*End*/

	// TypeSwitchStmt
	/*Start*/
	switch /*AfterSwitch*/ f.(type) /*AfterAssign*/ {
	} /*End*/

	/*Start*/
	switch /*AfterSwitch*/ g := f.(type) /*AfterAssign*/ {
	case int:
		print(g)
	} /*End*/

	/*Start*/
	switch /*AfterSwitch*/ g := f; /*AfterInit*/ g := g.(type) /*AfterAssign*/ {
	case int:
		print(g)
	} /*End*/

	// CommClause
	select {
	/*Start*/ case /*AfterCase*/ a := <-c /*AfterComm*/ : /*AfterColon*/
		print(a) /*End(?)*/
	}

	// SelectStmt
	/*Start*/
	select /*AfterSelect*/ {
	} /*End*/

	// ForStmt
	/*Start*/
	for /*AfterFor*/ {
		i++
	} /*End*/

	/*Start*/
	for /*AfterFor*/ i < 1 /*AfterCond*/ {
		i++
	} /*End*/

	/*Start*/
	for /*AfterFor*/ i = 0; /*AfterInit*/ i < 10; /*AfterCond*/ i++ /*AfterPost*/ {
		i++
	} /*End*/

	// RangeStmt (no comment possible before "range" if Key == nil, or between ":=" and "range")
	/*Start*/
	for range /*AfterRange*/ a /*AfterX*/ {
	} /*End*/

	/*Start*/
	for /*AfterFor*/ k /*AfterKey*/ := range /*AfterRange*/ a /*AfterX*/ {
		print(k)
	} /*End*/

	/*Start*/
	for /*AfterFor*/ k /*AfterKey*/, v /*AfterValue*/ := range /*AfterRange*/ a /*AfterX*/ {
		print(k, v)
	} /*End*/

	// ValueSpec (no comment possible before "=")
	var (
		/*Start*/ j = /*AfterAssign*/ 1 /*End*/
	)
	var (
		/*Start*/ k, l = /*AfterAssign*/ 1, 2 /*End*/
	)
	var (
		/*Start*/ m, n /*AfterNames*/ int = /*AfterAssign*/ 1, 2 /*End*/
	)

	print(j, k, l, m, n)

	// TypeSpec (no comment possible before "=")
	type (
		/*Start*/ T1 /*AfterName*/ []int /*End*/
	)

	type (
		/*Start*/ T2 = /*AfterName*/ T1 /*End*/
	)

	// GenDecl
	/*Start*/
	const /*AfterTok*/ ( /*AfterLparen*/
		a, b = 1, 2
		c    = 3
	) /*End*/

	/*Start*/
	const /*AfterTok*/ d = 1 /*End*/

}

/*Start*/
func /*AfterFunc*/ d /*AfterName*/ (d, e int) /*AfterParams*/ {
	return
} /*End*/

/*Start*/
func /*AfterFunc*/ (a *A) /*AfterRecv*/ e /*AfterName*/ (d, e int) /*AfterParams*/ {
	return
} /*End*/

/*Start*/
func /*AfterFunc*/ (a *A) /*AfterRecv*/ f /*AfterName*/ (d, e int) /*AfterParams*/ (f, g int) /*AfterResults*/ {
	return
} /*End*/
