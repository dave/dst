package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/andreyvit/diff"
)

func TestRestorer(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
	}{
		{
			name: "empty func",
			code: `package a

				func b() {
				} // c
				`,
		},
		{
			name: "inside if block",
			code: `package main

				func main() {
					if true {
						// a
					}
				}`,
		},
		{
			skip: true,
			name: "indented comment in case",
			code: `package main

			func main() {
				switch true {
				case true:
					// a
				case false:
				}
			}`,
		},
		{
			name: "non indented comment in case",
			code: `package main

			func main() {
				switch true {
				case true:
				// a
				case false:
				}
			}`,
		},
		{
			name: "comment in final case",
			code: `package main

			func main() {
				switch true {
				case true:
					// a
				}
			}`,
		},
		{
			name: "Interface param",
			code: `package main

			func a(a interface{}) {}`,
		},
		{
			name: "Field",
			code: `package main

			// Field
			type A struct {
    			A /*FieldAfterName*/ int /*FieldAfterType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
		},
		{
			name: "Ellipsis",
			code: `package main
			
			// Ellipsis
			func B(a ... /*EllipsisAfterEllipsis*/ int) {}`,
		},
		{
			name: "FuncLit",
			code: `package main
			
			// FuncLit
			var C = func(a int, b ...int) (c int) /*FuncLitAfterType*/ { return 0 }`,
		},
		{
			name: "CompositeLit",
			code: `package main
			
			// CompositeLit
			var D = A /*CompositeLitAfterType*/ { /*CompositeLitAfterLbrace*/ A: 0 /*CompositeLitAfterElts*/}`,
		},
		{
			name: "ParenExpr",
			code: `package main
			
			// ParenExpr
			var E = ( /*ParenExprAfterLparen*/ 1 + 1 /*ParenExprAfterX*/) / 2`,
		},
		{
			name: "SelectorExpr",
			code: `package main
			
			// SelectorExpr
			var F = fmt. /*SelectorExprAfterX*/ Sprint(0)`,
		},
		{
			name: "IndexExpr",
			code: `package main
			
			// IndexExpr
			var G = []int{0} /*IndexExprAfterX*/ [ /*IndexExprAfterLbrack*/ 0 /*IndexExprAfterIndex*/]`,
		},
		{
			name: "SliceExpr",
			code: `package main
			
			// SliceExpr
			var H = []int{0} /*SliceExprAfterX*/ [ /*SliceExprAfterLbrack*/ 1: /*SliceExprAfterLow*/ 2: /*SliceExprAfterHigh*/ 3 /*SliceExprAfterMax*/]`,
		},
		{
			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprAfterX*/ ( /*TypeAssertExprAfterLparen*/ int /*TypeAssertExprAfterType*/)`,
		},
		{
			name: "CallExpr",
			code: `package main
			
			// CallExpr
			var L = C /*CallExprAfterFun*/ ( /*CallExprAfterLparen*/ 0, []int{} /*CallExprAfterArgs*/ ... /*CallExprAfterEllipsis*/)`,
		},
		{
			name: "StarExpr",
			code: `package main
			
			// StarExpr
			var M = * /*StarExprAfterStar*/ (&I)`,
		},
		{
			name: "UnaryExpr",
			code: `package main
			
			// UnaryExpr
			var N = ^ /*UnaryExprAfterOp*/ 1`,
		},
		{
			name: "BinaryExpr",
			code: `package main
			
			// BinaryExpr
			var O = 1 /*BinaryExprAfterX*/ & /*BinaryExprAfterOp*/ 2`,
		},
		{
			name: "KeyValueExpr",
			code: `package main
			
			// KeyValueExpr
			var P = map[string]string{"a" /*KeyValueExprAfterKey*/ : /*KeyValueExprAfterColon*/ "a"}`,
		},
		{
			name: "ArrayType",
			code: `package main
			
			// ArrayType
			type Q [ /*ArrayTypeAfterLbrack*/ 1] /*ArrayTypeAfterLen*/ int`,
		},
		{
			name: "StructType",
			code: `package main
			
			// StructType
			type R struct /*StructTypeAfterStruct*/ {
				A int
			}`,
		},
		{
			name: "FuncType",
			code: `package main
			
			// FuncType
			type S func /*FuncTypeAfterFunc*/ (a int) /*FuncTypeAfterParams*/ (b int)`,
		},
		{
			name: "InterfaceType",
			code: `package main
			
			// InterfaceType
			type T interface /*InterfaceTypeAfterInterface*/ {
				A()
			}`,
		},
		{
			name: "MapType",
			code: `package main
			
			// MapType
			type U map[ /*MapTypeAfterMap*/ int] /*MapTypeAfterKey*/ int`,
		},
		{
			name: "ChanType",
			code: `package main
			
			// ChanType
			type V chan /*ChanTypeAfterBegin*/ int

			type W <-chan /*ChanTypeAfterBegin*/ int

			type X chan /*ChanTypeAfterBegin*/ <- /*ChanTypeAfterArrow*/ int`,
		},
		{
			name: "LabeledStmt, BranchStmt",
			code: `package main
			
			func main() {
				// LabeledStmt TODO: create cmd/gofmt issue for wonky BeforeNode comment positioning
				A /*LabeledStmtAfterLabel*/ : /*LabeledStmtAfterColon*/
				1++

				// BranchStmt
				goto /*BranchStmtAfterTok*/ A
			}`,
		},
		{
			name: "SendStmt",
			code: `package main
			
			func main() {
				// SendStmt
				B := make(chan int)
				B /*SendStmtAfterChan*/ <- /*SendStmtAfterArrow*/ 0
			}`,
		},
		{
			name: "IncDecStmt",
			code: `package main
			
			func main() {
				// IncDecStmt
				var C int
				C /*IncDecStmtAfterX*/ ++
			}`,
		},
		{
			name: "AssignStmt",
			code: `package main
			
			func main() {
				// AssignStmt
				D, E, F /*AssignStmtAfterLhs*/ := /*AssignStmtAfterTok*/ 1, 2, 3

				fmt.Println(D, E, F)
			}`,
		},
		{
			name: "GoStmt",
			code: `package main
			
			func main() {
				// GoStmt
				go /*GoStmtAfterGo*/ func() {}()
			}`,
		},
		{
			name: "DeferStmt",
			code: `package main
			
			func main() {
				// DeferStmt
				defer /*DeferStmtAfterDefer*/ func() {}()
			}`,
		},
		{
			name: "ReturnStmt",
			code: `package main
			
			func main() {
				// ReturnStmt
				func() (int, int, int) {
					return /*ReturnStmtAfterReturn*/ 1, 2, 3
				}()
			}`,
		},
		{
			name: "BlockStmt",
			code: `package main
			
			func main() {
				// BlockStmt
				if true { /*BlockStmtAfterLbrace*/
					1++
				}

				func() { /*BlockStmtAfterLbrace*/ 1++ }()
			}`,
		},
		{
			name: "IfStmt",
			code: `package main
			
			func main() {
				// IfStmt
				if /*IfStmtAfterIf*/ a := true; /*IfStmtAfterInit*/ a /*IfStmtAfterCond*/ {
					1++
				} else /*IfStmtAfterElse*/ {
					1++
				}
			}`,
		},
		{
			name: "CaseClause",
			code: `package main
			
			func main() {
				// CaseClause
				switch C {
				case /*CaseClauseAfterCase*/ 1, 2, 3 /*CaseClauseAfterList*/ : /*CaseClauseAfterColon*/
					1++
				default:
					1++
				}
			}`,
		},
		{
			name: "SwitchStmt",
			code: `package main
			
			func main() {
				// SwitchStmt
				switch /*SwitchStmtAfterSwitch*/ C /*SwitchStmtAfterTag*/ {
				}

				switch /*SwitchStmtAfterSwitch*/ a := C; /*SwitchStmtAfterInit*/ a /*SwitchStmtAfterTag*/ {
				}
			}`,
		},
		{
			name: "TypeSwitchStmt",
			code: `package main
			
			func main() {
				// TypeSwitchStmt
				switch /*TypeSwitchStmtAfterSwitch*/ I.(type) /*TypeSwitchStmtAfterAssign*/ {
				}

				switch /*TypeSwitchStmtAfterSwitch*/ j := I.(type) /*TypeSwitchStmtAfterAssign*/ {
				case int:
					fmt.Print(j)
				}

				switch /*TypeSwitchStmtAfterSwitch*/ j := I; /*TypeSwitchStmtAfterInit*/ j := j.(type) /*TypeSwitchStmtAfterAssign*/ {
				case int:
					fmt.Print(j)
				}
			}`,
		},
		{
			name: "CommClause",
			code: `package main
			
			func main() {
				// CommClause
				var a chan int
				select {
				case /*CommClauseAfterCase*/ b := <-a /*CommClauseAfterComm*/ : /*CommClauseAfterColon*/
					b++
				default:
				}		
			}`,
		},
		{
			name: "SelectStmt",
			code: `package main
			
			func main() {
				// SelectStmt
				select /*SelectStmtAfterSelect*/ {
				default:
				}
			}`,
		},
		{
			name: "ForStmt",
			code: `package main
			
			func main() {
				// ForStmt
				var i int
				for /*ForStmtAfterFor*/ {
					i++
				}

				for /*ForStmtAfterFor*/ i < 1 /*ForStmtAfterCond*/ {
					i++
				}

				for /*ForStmtAfterFor*/ i = 0; /*ForStmtAfterInit*/ i < 10; /*ForStmtAfterCond*/ i++ /*ForStmtAfterPost*/ {
					i++
				}
			}`,
		},
		{
			name: "RangeStmt",
			code: `package main
			
			func main() {	
				// RangeStmt
				var a []int

				for range /*RangeStmtAfterFor*/ a /*RangeStmtAfterX*/ {
				}

				for /*RangeStmtAfterFor*/ k /*RangeStmtAfterKey*/ := range /*RangeStmtAfterTok*/ a /*RangeStmtAfterX*/ {
				}

				for /*RangeStmtAfterFor*/ k /*RangeStmtAfterKey*/, v /*RangeStmtAfterValue*/ := range /*RangeStmtAfterTok*/ a /*RangeStmtAfterX*/ {
				}
			}`,
		},
		{
			name: "ImportSpec",
			code: `package main

			// ImportSpec
			import (
				// Doc
				/*ImportSpecAfterDoc*/ fmt /*ImportSpecAfterName*/ "fmt" /*ImportSpecAfterPath*/ // Comment
			)
	
			func main() {
				fmt.Println()
			}`,
		},
		{
			name: "ValueSpec",
			code: `package main
			
			func main() {
				// ValueSpec
				var /*ValueSpecAfterDoc*/ a = /*ValueSpecAfterNames*/ 1 /*ValueSpecAfterValues*/

				var /*ValueSpecAfterDoc*/ a, b = /*ValueSpecAfterNames*/ 1, 2 /*ValueSpecAfterValues*/

				var /*ValueSpecAfterDoc*/ a, b /*ValueSpecAfterNames*/ int = /*ValueSpecAfterType*/ 1, 2 /*ValueSpecAfterValues*/

			}`,
		},
		{
			name: "TypeSpec",
			code: `package main
			
			func main() {
				// TypeSpec
				type /*TypeSpecAfterDoc*/ T1 /*TypeSpecAfterName*/ []int /*TypeSpecAfterType*/

				type /*TypeSpecAfterDoc*/ T2 = /*TypeSpecAfterName*/ T1 /*TypeSpecAfterType*/
			}`,
		},
		{
			name: "GenDecl",
			code: `package main
			
			func main() {
				// GenDecl
				const /*GenDeclAfterTok*/ ( /*GenDeclAfterLparen*/
					a, b = 1, 2
					c    = 3
				) /*GenDeclAfterRparen*/
			}`,
		},
		{
			name: "FuncDecl",
			code: `package main
			
			// FuncDecl
			func /*FuncDeclAfterDoc*/ (a *b) /*FuncDeclAfterRecv*/ c /*FuncDeclAfterName*/ (d, e int) (f, g int) /*FuncDeclAfterType*/ {
			}`,
		},
	}
	var solo bool
	for _, test := range tests {
		if test.solo {
			solo = true
			break
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if solo && !test.solo {
				t.Skip()
			}
			if test.skip {
				t.Skip()
			}

			// format code and check it hasn't changed
			b, err := format.Source([]byte(test.code))
			if err != nil {
				t.Fatal(err)
			}
			if normalize(string(b)) != normalize(test.code) {
				t.Fatalf("code changed after gofmt. before: \n%s\nafter:\n%s", test.code, string(b))
			}

			// use the formatted version (correct indents etc.)
			test.code = string(b)

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "main.go", test.code, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			file := Decorate(f, fset)

			restoredFile, restoredFset := Restore(file)

			buf := &bytes.Buffer{}
			if err := format.Node(buf, restoredFset, restoredFile); err != nil {
				t.Fatal(err)
			}

			formatted, err := format.Source(buf.Bytes())
			if err != nil {
				t.Fatal(err)
			}

			if string(formatted) != test.code {
				t.Errorf("diff: %s", diff.LineDiff(test.code, string(formatted)))
			}
		})
	}
}
