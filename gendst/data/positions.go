// File
/*Start*/ package /*Package*/ data /*Name*/

// notest

// ImportSpec
import (
	/*Start*/ fmt /*Name*/ "fmt" /*End*/
)

// --

var a []int
var i = 1
var b bool
var f interface{} = 1
var p = &i
var c chan int

// Field
type A struct {
	/*Start*/ A int /*Type*/ `a:"a"` /*End*/
}

// FieldList
type A1 struct /*Start*/ { /*Opening*/
	a, b int
	c    string
} /*End*/

// Ellipsis
func B(a /*Start*/ ... /*Ellipsis*/ int /*End*/) {}

// FuncLit
var C = /*Start*/ func(a int, b ...int) (c int) /*Type*/ { return 0 } /*End*/

// CompositeLit
var D = /*Start*/ A /*Type*/ { /*Lbrace*/ A: 0} /*End*/

// ParenExpr
var E = /*Start*/ ( /*Lparen*/ 1 + 1 /*X*/) /*End*/ / 2

// SelectorExpr
var F = /*Start*/ tt. /*X*/ F /*End*/ ()

// IndexExpr
var G = /*Start*/ []int{0} /*X*/ [ /*Lbrack*/ 0 /*Index*/] /*End*/

// SliceExpr(0)
var H = /*Start*/ []int{0, 1, 2} /*X*/ [ /*Lbrack*/ 1: /*Low*/ 2: /*High*/ 3 /*Max*/] /*End*/

// SliceExpr(1)
var H1 = /*Start*/ []int{0, 1, 2} /*X*/ [ /*Lbrack*/ 1: /*Low*/ 2 /*High*/] /*End*/

// SliceExpr(2)
var H2 = /*Start*/ []int{0} /*X*/ [: /*Low*/] /*End*/

// SliceExpr(3)
var H3 = /*Start*/ []int{0} /*X*/ [ /*Lbrack*/ 1: /*Low*/] /*End*/

// SliceExpr(4)
var H4 = /*Start*/ []int{0, 1, 2} /*X*/ [: /*Low*/ 2 /*High*/] /*End*/

// SliceExpr(5)
var H5 = /*Start*/ []int{0, 1, 2} /*X*/ [: /*Low*/ 2: /*High*/ 3 /*Max*/] /*End*/

// TypeAssertExpr
var J = /*Start*/ f. /*X*/ ( /*Lparen*/ int /*Type*/) /*End*/

// CallExpr
var L = /*Start*/ C /*Fun*/ ( /*Lparen*/ 0, []int{}... /*Ellipsis*/) /*End*/

// StarExpr
var N = /*Start*/ * /*Star*/ p /*End*/

// UnaryExpr
var O = /*Start*/ ^ /*Op*/ 1 /*End*/

// BinaryExpr
var P = /*Start*/ 1 /*X*/ & /*Op*/ 2 /*End*/

// KeyValueExpr
var Q = map[string]string{
	/*Start*/ "a" /*Key*/ : /*Colon*/ "a", /*End*/
}

// ArrayType
type R /*Start*/ [ /*Lbrack*/ 1] /*Len*/ int /*End*/

// StructType
type S /*Start*/ struct /*Struct*/ {
	A int
} /*End*/

// FuncType
type T /*Start*/ func /*Func*/ (a int) /*Params*/ (b int) /*End*/

// InterfaceType
type U /*Start*/ interface /*Interface*/ {
	A()
} /*End*/

// MapType
type V /*Start*/ map[ /*Map*/ int] /*Key*/ int /*End*/

// ChanType(0)
type W /*Start*/ chan /*Begin*/ int /*End*/

// ChanType(1)
type X /*Start*/ <-chan /*Begin*/ int /*End*/

// ChanType(2)
type Y /*Start*/ chan /*Begin*/ <- /*Arrow*/ int /*End*/

// --

func Z() {
	// LabeledStmt
		/*Start*/
A /*Label*/ : /*Colon*/
	print("Stmt") /*End*/

	// BranchStmt
	/*Start*/
	goto /*Tok*/ A /*End*/

	// Ident(0)
	/*Start*/
	i /*End*/ ++

	// Ident(1)
	/*Start*/
	fmt. /*X*/ Print /*End*/ ()

	// SendStmt
	/*Start*/
	c /*Chan*/ <- /*Arrow*/ 0 /*End*/

	// IncDecStmt
	/*Start*/
	i /*X*/ ++ /*End*/

	// AssignStmt
	/*Start*/
	i = /*Tok*/ 1 /*End*/

	// GoStmt
	/*Start*/
	go /*Go*/ func() {}() /*End*/

	// DeferStmt
	/*Start*/
	defer /*Defer*/ func() {}() /*End*/

	// ReturnStmt
	func() int {
		/*Start*/ return /*Return*/ 1 /*End*/
	}()

	// BlockStmt(0)
	if true /*Start*/ { /*Lbrace*/
		i++
	} /*End*/

	// BlockStmt(1)
	func() /*Start*/ { /*Lbrace*/ i++ } /*End*/ ()

	// IfStmt
	/*Start*/
	if /*If*/ a := b; /*Init*/ a /*Cond*/ {
		i++
	} else /*Else*/ {
		i++
	} /*End*/

	// CaseClause
	switch i {
	/*Start*/ case /*Case*/ 1: /*Colon*/
		i++ /*End*/
	}

	// SwitchStmt(0)
	/*Start*/
	switch /*Switch*/ i /*Tag*/ {
	} /*End*/

	// SwitchStmt(1)
	/*Start*/
	switch /*Switch*/ a := i; /*Init*/ a /*Tag*/ {
	} /*End*/

	// TypeSwitchStmt(0)
	/*Start*/
	switch /*Switch*/ f.(type) /*Assign*/ {
	} /*End*/

	// TypeSwitchStmt(1)
	/*Start*/
	switch /*Switch*/ g := f.(type) /*Assign*/ {
	case int:
		print(g)
	} /*End*/

	// TypeSwitchStmt(2)
	/*Start*/
	switch /*Switch*/ g := f; /*Init*/ g := g.(type) /*Assign*/ {
	case int:
		print(g)
	} /*End*/

	// CommClause
	select {
	/*Start*/ case /*Case*/ a := <-c /*Comm*/ : /*Colon*/
		print(a) /*End*/
	}

	// SelectStmt
	/*Start*/
	select /*Select*/ {
	} /*End*/

	// ForStmt(0)
	/*Start*/
	for /*For*/ {
		i++
	} /*End*/

	// ForStmt(1)
	/*Start*/
	for /*For*/ i < 1 /*Cond*/ {
		i++
	} /*End*/

	// ForStmt(2)
	/*Start*/
	for /*For*/ i = 0; /*Init*/ i < 10; /*Cond*/ i++ /*Post*/ {
		i++
	} /*End*/

	// RangeStmt(0)
	/*Start*/
	for range /*Range*/ a /*X*/ {
	} /*End*/

	// RangeStmt(1)
	/*Start*/
	for /*For*/ k /*Key*/ := range /*Range*/ a /*X*/ {
		print(k)
	} /*End*/

	// RangeStmt(2)
	/*Start*/
	for /*For*/ k /*Key*/, v /*Value*/ := range /*Range*/ a /*X*/ {
		print(k, v)
	} /*End*/

	// ValueSpec(0)
	var (
		/*Start*/ j = /*Assign*/ 1 /*End*/
	)

	// ValueSpec(1)
	var (
		/*Start*/ k, l = /*Assign*/ 1, 2 /*End*/
	)

	// ValueSpec(2)
	var (
		/*Start*/ m, n int = /*Assign*/ 1, 2 /*End*/
	)

	// --

	print(j, k, l, m, n)

	// TypeSpec(0)
	type (
		/*Start*/ T1 /*Name*/ []int /*End*/
	)

	// TypeSpec(1)
	type (
		/*Start*/ T2 = /*Name*/ T1 /*End*/
	)

	// GenDecl(0)
	/*Start*/
	const /*Tok*/ ( /*Lparen*/
		a, b = 1, 2
		c    = 3
	) /*End*/

	// GenDecl(1)
	/*Start*/
	const /*Tok*/ d = 1 /*End*/

	// --

}

// FuncDecl(0)
/*Start*/
func /*Func*/ d /*Name*/ (d, e int) /*Params*/ {
	return
} /*End*/

// FuncDecl(1)
/*Start*/
func /*Func*/ (a *A) /*Recv*/ e /*Name*/ (d, e int) /*Params*/ {
	return
} /*End*/

// FuncDecl(2)
/*Start*/
func /*Func*/ (a *A) /*Recv*/ f /*Name*/ (d, e int) /*Params*/ (f, g int) /*Results*/ {
	return
} /*End*/

// --

type TT int

func (TT) F() int { return 0 }

var tt TT
