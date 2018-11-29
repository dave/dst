package decorator

import (
	"bytes"
	"go/format"
	"testing"
)

func TestRestorer(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
	}{
		{
			name: "comment-alignment",
			code: `package a

const (
	a = 1 // a
	b     // b
)`,
		},
		{
			name: "import-comment-alignment",
			code: `package a
            
import (
	"bytes"     // a
	"go/format" // b
)`,
		},
		{
			name: "comment-alignment-1",
			code: `package a

func main() {
	switch {
	case true:
		a := 1 // a
		a++    // b
	}
}`,
		},
		{
			name: "labelled-statement",
			code: `package a

func main() {
		/*Start*/
A /*Label*/ : /*Colon*/
	print("Stmt") /*End*/
}`,
		},
		{
			name: "hanging-indent-same-line",
			code: `package a

func a() {
	switch {
	case true: // a
		// b
	// c
	case false:
	}
}`,
		},
		{
			name: "hanging-indent",
			code: `package a

const a = 1 +
	1
	// a1

	// a2
const b = 1

const c = 1 +
	1

// d1

// d2
const d = 1
`,
		},
		{
			name: "net-hook",
			code: `package a

				var a = func(
					b int,
					c int,
				) int {
					return 1
				}`,
		},
		{
			name: "multi-line-string",
			code: `package a

				var a = b{
					c: ` + "`" + `
` + "`" + `,
				}`,
		},
		{
			name: "case clause",
			code: `package a
			
				func main() {
					switch a {
					case 1:
						// a
					// b
					case 2:
					// c
					case 3:
					}
				}`,
		},
		{
			name: "block comment",
			code: `package a
				
				/*
					foo
				*/
				var i int`,
		},
		{
			name: "case comment",
			code: `package a

				func main() {
					switch {
					default:
						// b
						// c

						var i int
					}
				}`,
		},
		{
			name: "file",
			code: `/*Start*/ package /*Package*/ postests /*Name*/

			var i int`,
		},
		{
			name: "RangeStmt",
			code: `package main
			
			func main() {	
				/*Start*/
				for /*For*/ k /*Key*/ := range /*Range*/ a /*X*/ {
				} /*End*/
			}`,
		},
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
    			A /*FieldName*/ int /*FieldType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
		},
		{
			name: "Ellipsis",
			code: `package main
			
			// Ellipsis
			func B(a ... /*EllipsisEllipsis*/ int) {}`,
		},
		{
			name: "FuncLit",
			code: `package main
			
			// FuncLit
			var C = func(a int, b ...int) (c int) /*FuncLitType*/ { return 0 }`,
		},
		{
			name: "CompositeLit",
			code: `package main
			
			// CompositeLit
			var D = A /*CompositeLitType*/ { /*CompositeLitLbrace*/ A: 0 /*CompositeLitElts*/}`,
		},
		{
			name: "ParenExpr",
			code: `package main
			
			// ParenExpr
			var E = ( /*ParenExprLparen*/ 1 + 1 /*ParenExprX*/) / 2`,
		},
		{
			name: "SelectorExpr",
			code: `package main
			
			// SelectorExpr
			var F = fmt. /*SelectorExprX*/ Sprint(0)`,
		},
		{
			name: "IndexExpr",
			code: `package main
			
			// IndexExpr
			var G = []int{0} /*IndexExprX*/ [ /*IndexExprLbrack*/ 0 /*IndexExprIndex*/]`,
		},
		{
			name: "SliceExpr",
			code: `package main
			
			// SliceExpr
			var H = []int{0} /*SliceExprX*/ [ /*SliceExprLbrack*/ 1: /*SliceExprLow*/ 2: /*SliceExprHigh*/ 3 /*SliceExprMax*/]`,
		},
		{
			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprX*/ ( /*TypeAssertExprLparen*/ int /*TypeAssertExprType*/)`,
		},
		{
			name: "CallExpr",
			code: `package main
			
			// CallExpr
			var L = C /*CallExprFun*/ ( /*CallExprLparen*/ 0, []int{} /*CallExprArgs*/ ... /*CallExprEllipsis*/)`,
		},
		{
			name: "StarExpr",
			code: `package main
			
			// StarExpr
			var M = * /*StarExprStar*/ (&I)`,
		},
		{
			name: "UnaryExpr",
			code: `package main
			
			// UnaryExpr
			var N = ^ /*UnaryExprOp*/ 1`,
		},
		{
			name: "BinaryExpr",
			code: `package main
			
			// BinaryExpr
			var O = 1 /*BinaryExprX*/ & /*BinaryExprOp*/ 2`,
		},
		{
			name: "KeyValueExpr",
			code: `package main
			
			// KeyValueExpr
			var P = map[string]string{"a" /*KeyValueExprKey*/ : /*KeyValueExprColon*/ "a"}`,
		},
		{
			name: "ArrayType",
			code: `package main
			
			// ArrayType
			type Q [ /*ArrayTypeLbrack*/ 1] /*ArrayTypeLen*/ int`,
		},
		{
			name: "StructType",
			code: `package main
			
			// StructType
			type R struct /*StructTypeStruct*/ {
				A int
			}`,
		},
		{
			name: "FuncType",
			code: `package main
			
			// FuncType
			type S func /*FuncTypeFunc*/ (a int) /*FuncTypeParams*/ (b int)`,
		},
		{
			name: "InterfaceType",
			code: `package main
			
			// InterfaceType
			type T interface /*InterfaceTypeInterface*/ {
				A()
			}`,
		},
		{
			name: "MapType",
			code: `package main
			
			// MapType
			type U map[ /*MapTypeMap*/ int] /*MapTypeKey*/ int`,
		},
		{
			name: "ChanType",
			code: `package main
			
			// ChanType
			type V chan /*ChanTypeBegin*/ int

			type W <-chan /*ChanTypeBegin*/ int

			type X chan /*ChanTypeBegin*/ <- /*ChanTypeArrow*/ int`,
		},
		{
			name: "LabeledStmt, BranchStmt",
			code: `package main
			
			func main() {
				// LabeledStmt TODO: create cmd/gofmt issue for wonky BeforeNode comment positioning
				A /*LabeledStmtLabel*/ : /*LabeledStmtColon*/
				1++

				// BranchStmt
				goto /*BranchStmtTok*/ A
			}`,
		},
		{
			name: "SendStmt",
			code: `package main
			
			func main() {
				// SendStmt
				B := make(chan int)
				B /*SendStmtChan*/ <- /*SendStmtArrow*/ 0
			}`,
		},
		{
			name: "IncDecStmt",
			code: `package main
			
			func main() {
				// IncDecStmt
				var C int
				C /*IncDecStmtX*/ ++
			}`,
		},
		{
			name: "AssignStmt",
			code: `package main
			
			func main() {
				// AssignStmt
				D, E, F /*AssignStmtLhs*/ := /*AssignStmtTok*/ 1, 2, 3

				fmt.Println(D, E, F)
			}`,
		},
		{
			name: "GoStmt",
			code: `package main
			
			func main() {
				// GoStmt
				go /*GoStmtGo*/ func() {}()
			}`,
		},
		{
			name: "DeferStmt",
			code: `package main
			
			func main() {
				// DeferStmt
				defer /*DeferStmtDefer*/ func() {}()
			}`,
		},
		{
			name: "ReturnStmt",
			code: `package main
			
			func main() {
				// ReturnStmt
				func() (int, int, int) {
					return /*ReturnStmtReturn*/ 1, 2, 3
				}()
			}`,
		},
		{
			name: "BlockStmt",
			code: `package main
			
			func main() {
				// BlockStmt
				if true { /*BlockStmtLbrace*/
					1++
				}

				func() { /*BlockStmtLbrace*/ 1++ }()
			}`,
		},
		{
			name: "IfStmt",
			code: `package main
			
			func main() {
				// IfStmt
				if /*IfStmtIf*/ a := true; /*IfStmtInit*/ a /*IfStmtCond*/ {
					1++
				} else /*IfStmtElse*/ {
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
				case /*CaseClauseCase*/ 1, 2, 3 /*CaseClauseList*/ : /*CaseClauseColon*/
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
				switch /*SwitchStmtSwitch*/ C /*SwitchStmtTag*/ {
				}

				switch /*SwitchStmtSwitch*/ a := C; /*SwitchStmtInit*/ a /*SwitchStmtTag*/ {
				}
			}`,
		},
		{
			name: "TypeSwitchStmt",
			code: `package main
			
			func main() {
				// TypeSwitchStmt
				switch /*TypeSwitchStmtSwitch*/ I.(type) /*TypeSwitchStmtAssign*/ {
				}

				switch /*TypeSwitchStmtSwitch*/ j := I.(type) /*TypeSwitchStmtAssign*/ {
				case int:
					fmt.Print(j)
				}

				switch /*TypeSwitchStmtSwitch*/ j := I; /*TypeSwitchStmtInit*/ j := j.(type) /*TypeSwitchStmtAssign*/ {
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
				case /*CommClauseCase*/ b := <-a /*CommClauseComm*/ : /*CommClauseColon*/
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
				select /*SelectStmtSelect*/ {
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
				for /*ForStmtFor*/ {
					i++
				}

				for /*ForStmtFor*/ i < 1 /*ForStmtCond*/ {
					i++
				}

				for /*ForStmtFor*/ i = 0; /*ForStmtInit*/ i < 10; /*ForStmtCond*/ i++ /*ForStmtPost*/ {
					i++
				}
			}`,
		},
		{
			name: "RangeStmt",
			code: `package main
			
			func main() {	
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
			}`,
		},
		{
			name: "ImportSpec",
			code: `package main

			// ImportSpec
			import (
				// Doc
				/*ImportSpecDoc*/ fmt /*ImportSpecName*/ "fmt" /*ImportSpecPath*/ // Comment
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
				var /*ValueSpecDoc*/ a = /*ValueSpecNames*/ 1 /*ValueSpecValues*/

				var /*ValueSpecDoc*/ a, b = /*ValueSpecNames*/ 1, 2 /*ValueSpecValues*/

				var /*ValueSpecDoc*/ a, b /*ValueSpecNames*/ int = /*ValueSpecType*/ 1, 2 /*ValueSpecValues*/

			}`,
		},
		{
			name: "TypeSpec",
			code: `package main
			
			func main() {
				// TypeSpec
				type /*TypeSpecDoc*/ T1 /*TypeSpecName*/ []int /*TypeSpecType*/

				type /*TypeSpecDoc*/ T2 = /*TypeSpecName*/ T1 /*TypeSpecType*/
			}`,
		},
		{
			name: "GenDecl",
			code: `package main
			
			func main() {
				// GenDecl
				const /*GenDeclTok*/ ( /*GenDeclLparen*/
					a, b = 1, 2
					c    = 3
				) /*GenDeclRparen*/
			}`,
		},
		{
			name: "FuncDecl",
			code: `package main
			
			// FuncDecl
			func /*FuncDeclDoc*/ (a *b) /*FuncDeclRecv*/ c /*FuncDeclName*/ (d, e int) (f, g int) /*FuncDeclType*/ {
			}`,
		},
		{
			name: "sel-space-decoration",
			code: `package main
            
				import "root/a"
            
            	func main() {
            		a. /*a*/
            
           				A()
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

			file, err := Parse(test.code)
			if err != nil {
				t.Fatal(err)
			}

			restoredFset, restoredFile, err := RestoreFile(file)
			if err != nil {
				t.Fatal(err)
			}

			buf := &bytes.Buffer{}
			if err := format.Node(buf, restoredFset, restoredFile); err != nil {
				t.Fatal(err)
			}

			if buf.String() != test.code {
				t.Errorf("diff:\n%s", diff(test.code, buf.String()))
			}
		})
	}
}
