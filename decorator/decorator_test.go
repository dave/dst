package decorator

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/dave/dst"
	"golang.org/x/tools/go/loader"
)

func TestDecorator(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
		expect     string
	}{
		{
			name: "chan type",
			code: `package main

				type Y /*Start*/ chan /*AfterBegin*/ <- /*AfterArrow*/ int /*End*/`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
            BlockStmt [AfterLbrace "\n"]
            IfStmt [End "\n"]
            BlockStmt [AfterLbrace "\n"] [AfterLbrace "// a"]`,
		},
		{
			name: "simple",
			code: `package main
			
			func main() {
				i // foo
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
            BlockStmt [AfterLbrace "\n"]
            ExprStmt [End "// foo"]`,
		},
		{
			name: "inline comment inside node",
			code: `package main
			
			func main() {
				i /* foo */ ++
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
            BlockStmt [AfterLbrace "\n"] [AfterLbrace "\n"]
            ExprStmt [Start "// foo"] [Start "\n"] [End "\n"]`,
		},
		{
			name: "comment statement",
			code: `package main
			
			func main() {
				// foo
				i
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
            BlockStmt [AfterLbrace "\n"]
            ExprStmt [Start "// foo"] [End "\n"]`,
		},
		{
			name: "comment after lbrace",
			code: `package main
			
			func main() { // foo
				i
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
            BlockStmt [AfterLbrace "// foo"]
            ExprStmt [End "\n"]`,
		},
		{
			name: "comment after func",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
            CompositeLit [AfterLbrace "\n"]
            KeyValueExpr [End "\n"] [End "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
            CompositeLit [AfterLbrace "\n"]
            KeyValueExpr [End "\n"] [End "\n"]
            KeyValueExpr [Start "// foo"] [Start "\n"] [End "\n"]`,
		},
		{
			name: "composite literal 4",
			code: `package main

			var A = B{
				"a": "b",
				// foo

				"c": "d",
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
            CompositeLit [AfterLbrace "\n"]
            KeyValueExpr [End "\n"]
            KeyValueExpr [Start "// foo"] [Start "\n"] [End "\n"]`,
			// TODO: Should "foo" be attached to the end of the first KeyValueExpr?
		},
		{
			name: "composite literal 5",
			code: `package main

			var A = B{
				"a": "b", // foo
				"c": "d",
			}`,
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
			expect: `File [AfterName "\n"] [AfterName "\n"]
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
						result += fmt.Sprintf(" [%s%q]", pos, v.Text)
					}
					result += "\n"
				}
				return true
			})

			if normalize(test.expect) != normalize(result) {
				t.Errorf("diff: %s", diff.LineDiff(normalize(test.expect), normalize(result)))
			}

		})
	}
}

func TestPositions(t *testing.T) {
	path := "github.com/dave/dst/gendst/postests"
	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		t.Fatal(err)
	}
	astFile := prog.Package(path).Files[0]

	d := New()
	file := d.Decorate(astFile, prog.Fset)

	r1 := regexp.MustCompile(`// ([a-zA-Z]+)\(([0-9])\)`)
	r2 := regexp.MustCompile(`// ([a-zA-Z]+)`)
	var currentNodeType string
	var currentTestIndex int
	var done bool

	dst.Inspect(file, func(n dst.Node) bool {
		if n == nil {
			return false
		}
		for _, d := range n.Decorations() {
			if r1.MatchString(d.Text) || r2.MatchString(d.Text) {
				if currentNodeType != "" && !done {
					t.Fatalf("missed %s %d", currentNodeType, currentTestIndex)
				}
				if matches := r1.FindStringSubmatch(d.Text); matches != nil {
					currentNodeType = "*dst." + matches[1]
					currentTestIndex, _ = strconv.Atoi(matches[2])
				} else if matches := r2.FindStringSubmatch(d.Text); matches != nil {
					currentNodeType = "*dst." + matches[1]
					currentTestIndex = 0
				}
				done = false
				break
			}
		}
		if fmt.Sprintf("%T", n) == currentNodeType {
			//fmt.Printf("*** Testing %s (%d)\n", currentNodeType, currentTestIndex)
			for _, d := range n.Decorations() {
				if !strings.HasPrefix(d.Text, "/*") {
					continue
				}
				text := strings.TrimSuffix(strings.TrimPrefix(d.Text, "/*"), "*/")
				if text != d.Position {
					t.Errorf("incorrect position in %s (%d) - expected %s, got %s", currentNodeType, currentTestIndex, text, d.Position)
				}
			}
			done = true
		} else {
			for _, d := range n.Decorations() {
				if !strings.HasPrefix(d.Text, "/*") {
					continue
				}
				text := strings.TrimSuffix(strings.TrimPrefix(d.Text, "/*"), "*/")
				if text != "Start" && text != "End" {
					// Only tolerate comments moved to adjacent decorations for Start and End
					t.Errorf("comment on wrong decoration: %s (%d) %s -> %T %s\n", currentNodeType, currentTestIndex, text, n, d.Position)
				}
			}
		}
		return true
	})
}

var multiSpaces = regexp.MustCompile(" {2,}")

func normalize(s string) string {
	s = multiSpaces.ReplaceAllString(s, "")
	s = strings.Replace(s, "\t", "", -1)
	s = strings.TrimSpace(s)
	return s
}
