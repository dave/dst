package decorator

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"

	"github.com/dave/dst"
)

func TestDecorator(t *testing.T) {
	tests := []struct {
		skip   bool
		name   string
		code   string
		expect string
	}{
		{
			name: "simple",
			code: `package main
			
			func main() {
				i // foo
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            BlockStmt [Lbrace (after) "\n"]
            ExprStmt [(after) "// foo"]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            BlockStmt [Lbrace (after) "\n"]
            IncDecStmt [X (after) "/* foo */"] [(after) "\n"]`,
		},
		{
			name: "comment statement spaced",
			code: `package main
			
			func main() {

				// foo

				i
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            BlockStmt [Lbrace (after) "\n"] [Lbrace (after) "\n"]
            ExprStmt [(before) "// foo"] [(before) "\n"] [(after) "\n"]`,
		},
		{
			name: "comment statement",
			code: `package main
			
			func main() {
				// foo
				i
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            BlockStmt [Lbrace (after) "\n"]
            ExprStmt [(before) "// foo"] [(after) "\n"]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            BlockStmt [Lbrace (after) "// foo"]
            ExprStmt [(after) "\n"]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            FuncType [Func (after) "/* foo */"]
            BlockStmt [Lbrace (after) "\n"]
            ExprStmt [(after) "\n"]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*FieldAfterName*/ int /*FieldAfterType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            FieldList [Opening (after) "\n"]
            Field [Type (after) "/*FieldAfterType*/"] [(after) "\n"]
            Ident [(after) "/*FieldAfterName*/"]`,
			// TODO: Should "FieldAfterName" be attached to Field.Names?
		},
		{
			name: "composite literal",
			code: `package main

			var A = B{
				"a": "b",
				"c": "d", // foo
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            CompositeLit [Lbrace (after) "\n"]
            KeyValueExpr [(after) "\n"]
            KeyValueExpr [(after) "// foo"]`,
		},
		{
			name: "composite literal 1",
			code: `package main

			var A = B{
				"a": "b",
				// foo
				"c": "d",
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            CompositeLit [Lbrace (after) "\n"]
            KeyValueExpr [(after) "\n"]
            KeyValueExpr [(before) "// foo"] [(after) "\n"]`,
		},
		{
			name: "composite literal 2",
			code: `package main

			var A = B{
				"a": "b",

				// foo
				"c": "d",
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            CompositeLit [Lbrace (after) "\n"]
            KeyValueExpr [(after) "\n"] [(after) "\n"]
            KeyValueExpr [(before) "// foo"] [(after) "\n"]`,
		},
		{
			name: "composite literal 3",
			code: `package main

			var A = B{
				"a": "b",

				// foo

				"c": "d",
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            CompositeLit [Lbrace (after) "\n"]
            KeyValueExpr [(after) "\n"] [(after) "\n"]
            KeyValueExpr [(before) "// foo"] [(before) "\n"] [(after) "\n"]`,
		},
		{
			name: "composite literal 4",
			code: `package main

			var A = B{
				"a": "b",
				// foo

				"c": "d",
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            CompositeLit [Lbrace (after) "\n"]
            KeyValueExpr [(after) "\n"]
            KeyValueExpr [(before) "// foo"] [(before) "\n"] [(after) "\n"]`,
			// TODO: Should "foo" be attached to the end of the first KeyValueExpr?
		},
		{
			name: "composite literal 5",
			code: `package main

			var A = B{
				"a": "b", // foo
				"c": "d",
			}`,
			expect: `File [Name (after) "\n"] [Name (after) "\n"]
            CompositeLit [Lbrace (after) "\n"]
            KeyValueExpr [(after) "// foo"]
            KeyValueExpr [(after) "\n"]`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "main.go", test.code, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			d := New()
			file := d.Decorate(f, fset)
			var result string
			dst.Inspect(file, func(n dst.Node) bool {
				if n == nil {
					return false
				}
				if d := n.Decorations(); len(d) > 0 {
					nodeType := func(n dst.Node) string {
						return strings.Replace(fmt.Sprintf("%T", n), "*dst.", "", -1)
					}
					result += fmt.Sprintf("%s", nodeType(n))
					for _, v := range d {
						var pos string
						if v.Position != "" {
							pos = v.Position + " "
						}
						var rel string
						if v.Start {
							rel = "(before) "
						} else {
							rel = "(after) "
						}
						result += fmt.Sprintf(" [%s%s%q]", pos, rel, v.Text)
					}
					result += "\n"
				}
				return true
			})

			if normalize(test.expect) != normalize(result) {
				t.Fatalf("expected:\n%s\n\nresult:\n%s", normalize(test.expect), normalize(result))
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
