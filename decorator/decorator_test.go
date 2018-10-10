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
			expect: `File [Name "\n" "\n"]
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
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				SwitchStmt [New line]
				BlockStmt [Lbrace "\n"]
				CaseClause [Colon "\n" "// b" "// c" "\n"]
				DeclStmt [New line]`,
		},
		{
			name: "file",
			code: `/*Start*/ package /*Package*/ postests /*Name*/

			var i int`,
			expect: `File [Start "/*Start*/"] [Package "/*Package*/"] [Name "/*Name*/" "\n" "\n"]`,
		},
		{

			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprX*/ ( /*TypeAssertExprLparen*/ int /*TypeAssertExprType*/)`,
			expect: `File [Name "\n" "\n"]
				GenDecl [Start "// TypeAssertExpr"] [New line]
				TypeAssertExpr [X "/*TypeAssertExprX*/"] [Lparen "/*TypeAssertExprLparen*/"] [Type "/*TypeAssertExprType*/"]`,
		},
		{
			name: "range bug",
			code: `package main

				func main() {
					/*Start*/
					for /*For*/ k /*Key*/, v /*Value*/ := range /*Range*/ a /*X*/ {
						print(k, v)
					} /*End*/
				}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				RangeStmt [Start "/*Start*/" "\n"] [For "/*For*/"] [Key "/*Key*/"] [Value "/*Value*/"] [Range "/*Range*/"] [X "/*X*/"] [End "/*End*/"] [New line]
				BlockStmt [Lbrace "\n"]
				ExprStmt [New line]`,
		},
		{
			name: "value spec",
			code: `package main

				func main() {
					var foo int
				}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				DeclStmt [New line]`,
		},
		{
			name: "chan type",
			code: `package main

				type Y /*Start*/ chan /*Begin*/ <- /*Arrow*/ int /*End*/`,
			expect: `File [Name "\n" "\n"]
				GenDecl [End "/*End*/"]
				TypeSpec [Name "/*Start*/"]
				ChanType [Begin "/*Begin*/"] [Arrow "/*Arrow*/"]`,
		},
		{
			name: "inside if block",
			code: `package main

				func main() {
					if true {
						// a
					}
				}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				IfStmt [New line]
				BlockStmt [Lbrace "\n" "// a"]`,
		},
		{
			name: "simple",
			code: `package main
			
			func main() {
				i // foo
			}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				ExprStmt [End "// foo"] [New line]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				IncDecStmt [X "/* foo */"] [New line]`,
		},
		{
			name: "comment statement spaced",
			code: `package main
			
			func main() {

				// foo

				i
			}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n" "\n"]
				ExprStmt [Start "// foo" "\n"] [New line]`,
		},
		{
			name: "comment statement",
			code: `package main
			
			func main() {
				// foo
				i
			}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "\n"]
				ExprStmt [Start "// foo"] [New line]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `File [Name "\n" "\n"]
				BlockStmt [Lbrace "// foo"]
				ExprStmt [New line]`,
		},
		{
			name: "comment after func",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `File [Name "\n" "\n"]
				FuncDecl [Func "/* foo */"]
				BlockStmt [Lbrace "\n"]
				ExprStmt [New line]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*FieldName*/ int /*FieldType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `File [Name "\n" "\n"]
				FieldList [Opening "\n"]
				Field [Names "/*FieldName*/"] [Type "/*FieldType*/"] [New line]`,
		},
		{
			name: "composite literal",
			code: `package main

			var A = B{
				"a": "b",
				"c": "d", // foo
			}`,
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
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
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
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
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
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
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
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
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
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
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
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
			expect: `File [Name "\n" "\n"]
				CompositeLit [Lbrace "\n"]
				KeyValueExpr [End "// foo"] [New line]
				KeyValueExpr [New line]`,
		},
		{
			name: "FuncDecl",
			code: `package main
			
			// FuncDecl
			func /*FuncDeclDoc*/ (a *b) /*FuncDeclRecv*/ c /*FuncDeclName*/ (d, e int) /*FuncDeclParams*/ (f, g int) /*FuncDeclType*/ {
			}`,
			expect: `File [Name "\n" "\n"]
				FuncDecl [Start "// FuncDecl"] [Func "/*FuncDeclDoc*/"] [Recv "/*FuncDeclRecv*/"] [Name "/*FuncDeclName*/"] [Params "/*FuncDeclParams*/"] [Results "/*FuncDeclType*/"]
				BlockStmt [Lbrace "\n"]`,
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
