package decorator

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/dave/dst"
)

func TestDecorator(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
		expect     string
	}{
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
            RangeStmt [Start "/*Start*/" "\n"] [AfterFor "/*AfterFor*/"] [AfterKey "/*AfterKey*/"] [AfterValue "/*AfterValue*/"] [AfterRange "/*AfterRange*/"] [AfterX "/*AfterX*/"] [End "/*End*/" "\n"]
            BlockStmt [AfterLbrace "\n"]
            ExprStmt [End "\n"]`,
		},
		{
			name: "value spec",
			code: `package main

				func main() {
					var foo int
				}`,
			expect: `File [AfterName "\n" "\n"]
            BlockStmt [AfterLbrace "\n"]
            DeclStmt [End "\n"]`,
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
            IfStmt [End "\n"]
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
            ExprStmt [End "// foo"]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `File [AfterName "\n" "\n"]
            BlockStmt [AfterLbrace "\n"]
            IncDecStmt [AfterX "/* foo */"] [End "\n"]`,
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
            ExprStmt [Start "// foo" "\n"] [End "\n"]`,
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
            ExprStmt [Start "// foo"] [End "\n"]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `File [AfterName "\n" "\n"]
            BlockStmt [AfterLbrace "// foo"]
            ExprStmt [End "\n"]`,
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
            ExprStmt [End "\n"]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*FieldAfterName*/ int /*FieldAfterType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `File [AfterName "\n" "\n"]
            FieldList [AfterOpening "\n"]
            Field [AfterNames "/*FieldAfterName*/"] [AfterType "/*FieldAfterType*/"] [End "\n"]`,
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
            KeyValueExpr [End "\n"]
            KeyValueExpr [End "// foo"]`,
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
            KeyValueExpr [End "\n"]
            KeyValueExpr [Start "// foo"] [End "\n"]`,
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
            KeyValueExpr [End "\n" "\n"]
            KeyValueExpr [Start "// foo"] [End "\n"]`,
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
            KeyValueExpr [End "\n" "\n"]
            KeyValueExpr [Start "// foo" "\n"] [End "\n"]`,
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
            KeyValueExpr [End "\n" "// foo" "\n"]
            KeyValueExpr [End "\n"]`,
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
            KeyValueExpr [End "\n"]
            KeyValueExpr [Start "// foo"] [End "\n"]`,
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
            KeyValueExpr [End "// foo"]
            KeyValueExpr [End "\n"]`,
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
			var result string
			nodeType := func(n dst.Node) string {
				return strings.Replace(fmt.Sprintf("%T", n), "*dst.", "", -1)
			}
			dst.Inspect(file, func(n dst.Node) bool {
				if n == nil {
					return false
				}
				var out string
				infos := getDecorationInfo(n)
				for _, info := range infos {
					if len(info.decs) > 0 {
						var values string
						for i, dec := range info.decs {
							if i > 0 {
								values += " "
							}
							values += fmt.Sprintf("%q", dec)
						}
						out += fmt.Sprintf(" [%s %s]", info.name, values)
					}
				}
				if out != "" {
					result += nodeType(n) + out + "\n"
				}
				return true
			})

			if normalize(test.expect) != normalize(result) {
				t.Errorf("diff: %s", diff.LineDiff(normalize(test.expect), normalize(result)))
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
