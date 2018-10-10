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
			expect: `GenDecl [Empty line before] [Start "/*\n\tfoo\n*/" "\n"]`,
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
			expect: `FuncDecl [Empty line before]
				SwitchStmt [New line before] [New line after]
				CaseClause [New line before] [AfterColon "\n" "// b" "// c"]
				DeclStmt [Empty line before] [New line after]`,
		},
		{
			name: "file",
			code: `/*Start*/ package /*AfterPackage*/ postests /*AfterName*/

			var i int`,
			expect: `File [Start "/*Start*/"] [AfterPackage "/*AfterPackage*/"] [AfterName "/*AfterName*/"]
				GenDecl [Empty line before]`,
		},
		{

			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprAfterX*/ ( /*TypeAssertExprAfterLparen*/ int /*TypeAssertExprAfterType*/)`,
			expect: `GenDecl [Empty line before] [Start "// TypeAssertExpr"] [New line after]
				GenDecl [New line before]
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
			expect: `FuncDecl [Empty line before]
				RangeStmt [New line before] [Start "/*Start*/" "\n"] [AfterFor "/*AfterFor*/"] [AfterKey "/*AfterKey*/"] [AfterValue "/*AfterValue*/"] [AfterRange "/*AfterRange*/"] [AfterX "/*AfterX*/"] [End "/*End*/"] [New line after]
				ExprStmt [New line before] [New line after]`,
		},
		{
			name: "value spec",
			code: `package main

				func main() {
					var foo int
				}`,
			expect: `FuncDecl [Empty line before]
				DeclStmt [New line before] [New line after]`,
		},
		{
			name: "chan type",
			code: `package main

				type Y /*Start*/ chan /*AfterBegin*/ <- /*AfterArrow*/ int /*End*/`,
			expect: `GenDecl [Empty line before] [End "/*End*/"]
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
			expect: `FuncDecl [Empty line before]
				IfStmt [New line before] [New line after]
				BlockStmt [AfterLbrace "\n" "// a"]`,
		},
		{
			name: "simple",
			code: `package main
			
			func main() {
				i // foo
			}`,
			expect: `FuncDecl [Empty line before]
				ExprStmt [New line before] [End "// foo"] [New line after]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `FuncDecl [Empty line before]
				IncDecStmt [New line before] [AfterX "/* foo */"] [New line after]`,
		},
		{
			name: "comment statement spaced",
			code: `package main
			
			func main() {

				// foo

				i
			}`,
			expect: `FuncDecl [Empty line before]
				ExprStmt [Empty line before] [Start "// foo" "\n"] [New line after]`,
		},
		{
			name: "comment statement",
			code: `package main
			
			func main() {
				// foo
				i
			}`,
			expect: `FuncDecl [Empty line before]
				ExprStmt [New line before] [Start "// foo"] [New line after]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `FuncDecl [Empty line before]
				BlockStmt [AfterLbrace "// foo"]
				ExprStmt [New line before] [New line after]`,
		},
		{
			name: "comment after func",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `FuncDecl [Empty line before] [AfterFunc "/* foo */"]
				ExprStmt [New line before] [New line after]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*FieldAfterName*/ int /*FieldAfterType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `GenDecl [Empty line before]
				Field [New line before] [AfterNames "/*FieldAfterName*/"] [AfterType "/*FieldAfterType*/"] [New line after]`,
		},
		{
			name: "composite literal",
			code: `package main

			var A = B{
				"a": "b",
				"c": "d", // foo
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [New line after]
				KeyValueExpr [New line before] [End "// foo"] [New line after]`,
		},
		{
			name: "composite literal 1",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [New line after]
				KeyValueExpr [New line before] [Start "// foo"] [New line after]`,
		},
		{
			name: "composite literal 2",
			code: `package main

			var A = B{
				"a": "b",

				// foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [Empty line after]
				KeyValueExpr [Empty line before] [Start "// foo"] [New line after]`,
		},
		{
			name: "composite literal 3",
			code: `package main

			var A = B{
				"a": "b",

				// foo

				"c": "d",
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [Empty line after]
				KeyValueExpr [Empty line before] [Start "// foo" "\n"] [New line after]`,
		},
		{
			name: "composite literal 4",
			code: `package main

			var A = B{
				"a": "b",
				// foo

				"c": "d",
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [End "\n" "// foo"] [Empty line after]
				KeyValueExpr [Empty line before] [New line after]`,
		},
		{
			name: "composite literal 4a",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [New line after]
				KeyValueExpr [New line before] [Start "// foo"] [New line after]`,
		},
		{
			name: "composite literal 5",
			code: `package main

			var A = B{
				"a": "b", // foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line before]
				KeyValueExpr [New line before] [End "// foo"] [New line after]
				KeyValueExpr [New line before] [New line after]`,
		},
		{
			name: "FuncDecl",
			code: `package main
			
			// FuncDecl
			func /*FuncDeclAfterDoc*/ (a *b) /*FuncDeclAfterRecv*/ c /*FuncDeclAfterName*/ (d, e int) /*FuncDeclAfterParams*/ (f, g int) /*FuncDeclAfterType*/ {
			}`,
			expect: `FuncDecl [Empty line before] [Start "// FuncDecl"] [AfterFunc "/*FuncDeclAfterDoc*/"] [AfterRecv "/*FuncDeclAfterRecv*/"] [AfterName "/*FuncDeclAfterName*/"] [AfterParams "/*FuncDeclAfterParams*/"] [AfterResults "/*FuncDeclAfterType*/"]
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
