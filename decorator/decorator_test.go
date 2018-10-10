package decorator

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func TestDecorator(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
		expect     string
	}{
		{
			name: "block comment",
			code: `package a
				
				/*
					foo
				*/
				var i int`,
			expect: `File [AfterName "\n" "\n"]
				GenDecl [Start "/*\n\tfoo\n*/" "\n"]`,
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
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				SwitchStmt [New line]
				BlockStmt [AfterLbrace "\n"]
				CaseClause [AfterColon "\n" "// b" "// c" "\n"]
				DeclStmt [New line]`,
		},
		{
			name: "file",
			code: `/*Start*/ package /*AfterPackage*/ postests /*AfterName*/

			var i int`,
			expect: `File [Start "/*Start*/"] [AfterPackage "/*AfterPackage*/"] [AfterName "/*AfterName*/" "\n" "\n"]`,
		},
		{

			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprAfterX*/ ( /*TypeAssertExprAfterLparen*/ int /*TypeAssertExprAfterType*/)`,
			expect: `File [AfterName "\n" "\n"]
				GenDecl [Start "// TypeAssertExpr"] [New line]
				TypeAssertExpr [AfterX "/*TypeAssertExprAfterX*/"] [AfterLparen "/*TypeAssertExprAfterLparen*/"] [AfterType "/*TypeAssertExprAfterType*/"]`,
		},
		{
			name: "range bug",
			code: `package main

				func main() {
					/*Start*/
					for /*AfterFor*/ k /*AfterKey*/, v /*AfterValue*/ := range /*AfterRange*/ a /*AfterX*/ {
						print(k, v)
					} /*End*/
				}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				RangeStmt [Start "/*Start*/" "\n"] [AfterFor "/*AfterFor*/"] [AfterKey "/*AfterKey*/"] [AfterValue "/*AfterValue*/"] [AfterRange "/*AfterRange*/"] [AfterX "/*AfterX*/"] [End "/*End*/"] [New line]
				BlockStmt [AfterLbrace "\n"]
				ExprStmt [New line]`,
		},
		{
			name: "value spec",
			code: `package main

				func main() {
					var foo int
				}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				DeclStmt [New line]`,
		},
		{
			name: "chan type",
			code: `package main

				type Y /*Start*/ chan /*AfterBegin*/ <- /*AfterArrow*/ int /*End*/`,
			expect: `File [AfterName "\n" "\n"]
				GenDecl [End "/*End*/"]
				TypeSpec [AfterName "/*Start*/"]
				ChanType [AfterBegin "/*AfterBegin*/"] [AfterArrow "/*AfterArrow*/"]`,
		},
		{
			name: "inside if block",
			code: `package main

				func main() {
					if true {
						// a
					}
				}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				IfStmt [New line]
				BlockStmt [AfterLbrace "\n" "// a"]`,
		},
		{
			name: "simple",
			code: `package main
			
			func main() {
				i // foo
			}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				ExprStmt [End "// foo"] [New line]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				IncDecStmt [AfterX "/* foo */"] [New line]`,
		},
		{
			name: "comment statement spaced",
			code: `package main
			
			func main() {

				// foo

				i
			}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n" "\n"]
				ExprStmt [Start "// foo" "\n"] [New line]`,
		},
		{
			name: "comment statement",
			code: `package main
			
			func main() {
				// foo
				i
			}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "\n"]
				ExprStmt [Start "// foo"] [New line]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `File [AfterName "\n" "\n"]
				BlockStmt [AfterLbrace "// foo"]
				ExprStmt [New line]`,
		},
		{
			name: "comment after func",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `File [AfterName "\n" "\n"]
				FuncDecl [AfterFunc "/* foo */"]
				BlockStmt [AfterLbrace "\n"]
				ExprStmt [New line]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*FieldAfterName*/ int /*FieldAfterType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `File [AfterName "\n" "\n"]
				FieldList [AfterOpening "\n"]
				Field [AfterNames "/*FieldAfterName*/"] [AfterType "/*FieldAfterType*/"] [New line]`,
		},
		{
			name: "composite literal",
			code: `package main

			var A = B{
				"a": "b",
				"c": "d", // foo
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [New line]
				KeyValueExpr [End "// foo"] [New line]`,
		},
		{
			name: "composite literal 1",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [New line]
				KeyValueExpr [Start "// foo"] [New line]`,
		},
		{
			name: "composite literal 2",
			code: `package main

			var A = B{
				"a": "b",

				// foo
				"c": "d",
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [Empty line]
				KeyValueExpr [Start "// foo"] [New line]`,
		},
		{
			name: "composite literal 3",
			code: `package main

			var A = B{
				"a": "b",

				// foo

				"c": "d",
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [Empty line]
				KeyValueExpr [Start "// foo" "\n"] [New line]`,
		},
		{
			name: "composite literal 4",
			code: `package main

			var A = B{
				"a": "b",
				// foo

				"c": "d",
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [End "\n" "// foo"] [Empty line]
				KeyValueExpr [New line]`,
		},
		{
			name: "composite literal 4a",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [New line]
				KeyValueExpr [Start "// foo"] [New line]`,
		},
		{
			name: "composite literal 5",
			code: `package main

			var A = B{
				"a": "b", // foo
				"c": "d",
			}`,
			expect: `File [AfterName "\n" "\n"]
				CompositeLit [AfterLbrace "\n"]
				KeyValueExpr [End "// foo"] [New line]
				KeyValueExpr [New line]`,
		},
		{
			name: "FuncDecl",
			code: `package main
			
			// FuncDecl
			func /*FuncDeclAfterDoc*/ (a *b) /*FuncDeclAfterRecv*/ c /*FuncDeclAfterName*/ (d, e int) /*FuncDeclAfterParams*/ (f, g int) /*FuncDeclAfterType*/ {
			}`,
			expect: `File [AfterName "\n" "\n"]
				FuncDecl [Start "// FuncDecl"] [AfterFunc "/*FuncDeclAfterDoc*/"] [AfterRecv "/*FuncDeclAfterRecv*/"] [AfterName "/*FuncDeclAfterName*/"] [AfterParams "/*FuncDeclAfterParams*/"] [AfterResults "/*FuncDeclAfterType*/"]
				BlockStmt [AfterLbrace "\n"]`,
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
			file := Decorate(fset, f)

			buf := &bytes.Buffer{}
			debug(buf, file)

			if normalize(test.expect) != normalize(buf.String()) {
				t.Errorf("diff: %s", diff.LineDiff(normalize(test.expect), normalize(buf.String())))

				fmt.Println(buf.String())
			}

		})
	}
}

var multiSpaces = regexp.MustCompile(" {2,}")

func normalize(s string) string {
	s = multiSpaces.ReplaceAllString(s, "")
	s = strings.Replace(s, "\t", "", -1)
	s = strings.TrimSpace(s)
	return s
}
