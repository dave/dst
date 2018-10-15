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
			name: "net-hook",
			code: `package a

				var a = func(
					b int,
					c int,
				) int {
					return 1
				}`,
			expect: `GenDecl [Empty line space]
Field [New line space] [New line after]
Field [New line space] [New line after]
ReturnStmt [New line space] [New line after]`,
		},
		{
			name: "multi-line-string",
			code: `package a

				var a = b{
					c: ` + "`" + `
` + "`" + `,
				}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [New line after]`,
		},
		{
			name: "case clause",
			code: `package a
			
				func main() {
					switch a {
					case 1:
						// a
					case 2:
					// b
					case 3:
					}
				}`,
			expect: `FuncDecl [Empty line space]
SwitchStmt [New line space] [New line after]
CaseClause [New line space] [Colon "\n" "// a"]
CaseClause [New line space]
CaseClause [New line space] [Start "// b"] [Colon "\n"]`,
		},
		{
			name: "block comment",
			code: `package a
				
				/*
					foo
				*/
				var i int`,
			expect: `GenDecl [Empty line space] [Start "/*\n\tfoo\n*/" "\n"]`,
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
			expect: `FuncDecl [Empty line space]
SwitchStmt [New line space] [New line after]
CaseClause [New line space] [Colon "\n" "// b" "// c"]
DeclStmt [Empty line space] [New line after]`,
		},
		{
			name: "file",
			code: `/*Start*/ package /*Package*/ postests /*Name*/

			var i int`,
			expect: `File [Start "/*Start*/"] [Package "/*Package*/"] [Name "/*Name*/"]
GenDecl [Empty line space]`,
		},
		{

			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprX*/ ( /*TypeAssertExprLparen*/ int /*TypeAssertExprType*/)`,
			expect: `GenDecl [Empty line space] [Start "// TypeAssertExpr"] [New line after]
GenDecl [New line space]
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
			expect: `FuncDecl [Empty line space]
RangeStmt [New line space] [Start "/*Start*/" "\n"] [For "/*For*/"] [Key "/*Key*/"] [Value "/*Value*/"] [Range "/*Range*/"] [X "/*X*/"] [End "/*End*/"] [New line after]
ExprStmt [New line space] [New line after]`,
		},
		{
			name: "value spec",
			code: `package main

				func main() {
					var foo int
				}`,
			expect: `FuncDecl [Empty line space]
DeclStmt [New line space] [New line after]`,
		},
		{
			name: "chan type",
			code: `package main

				type Y /*Start*/ chan /*Begin*/ <- /*Arrow*/ int /*End*/`,
			expect: `GenDecl [Empty line space] [End "/*End*/"]
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
			expect: `FuncDecl [Empty line space]
IfStmt [New line space] [New line after]
BlockStmt [Lbrace "\n" "// a"]`,
		},
		{
			name: "simple",
			code: `package main
			
			func main() {
				i // foo
			}`,
			expect: `FuncDecl [Empty line space]
ExprStmt [New line space] [End "// foo"] [New line after]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `FuncDecl [Empty line space]
IncDecStmt [New line space] [X "/* foo */"] [New line after]`,
		},
		{
			name: "comment statement spaced",
			code: `package main
			
			func main() {

				// foo

				i
			}`,
			expect: `FuncDecl [Empty line space]
ExprStmt [Empty line space] [Start "// foo" "\n"] [New line after]`,
		},
		{
			name: "comment statement",
			code: `package main
			
			func main() {
				// foo
				i
			}`,
			expect: `FuncDecl [Empty line space]
ExprStmt [New line space] [Start "// foo"] [New line after]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `FuncDecl [Empty line space]
BlockStmt [Lbrace "// foo"]
ExprStmt [New line space] [New line after]`,
		},
		{
			name: "comment after func",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `FuncDecl [Empty line space] [Func "/* foo */"]
ExprStmt [New line space] [New line after]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*FieldName*/ int /*FieldType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `GenDecl [Empty line space]
Field [New line space] [Names "/*FieldName*/"] [Type "/*FieldType*/"] [New line after]`,
		},
		{
			name: "composite literal",
			code: `package main

			var A = B{
				"a": "b",
				"c": "d", // foo
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [New line after]
KeyValueExpr [New line space] [End "// foo"] [New line after]`,
		},
		{
			name: "composite literal 1",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [New line after]
KeyValueExpr [New line space] [Start "// foo"] [New line after]`,
		},
		{
			name: "composite literal 2",
			code: `package main

			var A = B{
				"a": "b",

				// foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [Empty line after]
KeyValueExpr [Empty line space] [Start "// foo"] [New line after]`,
		},
		{
			name: "composite literal 3",
			code: `package main

			var A = B{
				"a": "b",

				// foo

				"c": "d",
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [Empty line after]
KeyValueExpr [Empty line space] [Start "// foo" "\n"] [New line after]`,
		},
		{
			name: "composite literal 4",
			code: `package main

			var A = B{
				"a": "b",
				// foo

				"c": "d",
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [End "\n" "// foo"] [Empty line after]
KeyValueExpr [Empty line space] [New line after]`,
		},
		{
			name: "composite literal 4a",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [New line after]
KeyValueExpr [New line space] [Start "// foo"] [New line after]`,
		},
		{
			name: "composite literal 5",
			code: `package main

			var A = B{
				"a": "b", // foo
				"c": "d",
			}`,
			expect: `GenDecl [Empty line space]
KeyValueExpr [New line space] [End "// foo"] [New line after]
KeyValueExpr [New line space] [New line after]`,
		},
		{
			name: "FuncDecl",
			code: `package main
			
			// FuncDecl
			func /*FuncDeclDoc*/ (a *b) /*FuncDeclRecv*/ c /*FuncDeclName*/ (d, e int) /*FuncDeclParams*/ (f, g int) /*FuncDeclType*/ {
			}`,
			expect: `FuncDecl [Empty line space] [Start "// FuncDecl"] [Func "/*FuncDeclDoc*/"] [Recv "/*FuncDeclRecv*/"] [Name "/*FuncDeclName*/"] [Params "/*FuncDeclParams*/"] [Results "/*FuncDeclType*/"]
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
