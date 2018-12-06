package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestDecorator(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
		expect     string
	}{
		{
			name: "import-blocks",
			code: `package main

				// first-import-block
            	import (
					"root/a"
					"root/b"
				)

				// second-import-block
				import (
					"root/c"
					"root/d"
				)`,
			expect: `GenDecl [Empty line before] [Start "// first-import-block"] [Empty line after]
ImportSpec [New line before] [New line after]
ImportSpec [New line before] [New line after]
GenDecl [Empty line before] [Start "// second-import-block"]
ImportSpec [New line before] [New line after]
ImportSpec [New line before] [New line after]`,
		},
		{
			name: "comment-alignment",
			code: `package a

const (
	a = 1 // a
	b     // b
)`,
			expect: `GenDecl [Empty line before]
ValueSpec [New line before] [End "// a"] [New line after]
ValueSpec [New line before] [End "// b"] [New line after]`,
		},
		{
			name: "labelled-statement-hanging-indent",
			code: `package a

func main() {
A:
	var a int
	// a

	// b
	var b int
}`,
			expect: `FuncDecl [Empty line before]
LabeledStmt [New line before] [End "\n" "// a"] [Empty line after]
DeclStmt [New line before]
DeclStmt [Empty line before] [Start "// b"] [New line after]`,
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
			expect: `FuncDecl [Empty line before]
SwitchStmt [New line before] [New line after]
CaseClause [New line before] [End "// a" "// b"] [New line after]
CaseClause [New line before] [Start "// c"] [New line after]`,
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
			expect: `GenDecl [Empty line before] [End "\n" "// a1" "\n" "// a2"] [New line after]
BasicLit [New line before]
GenDecl [New line before] [Empty line after]
GenDecl [Empty line before] [Empty line after]
BasicLit [New line before]
GenDecl [Empty line before] [Start "// d1" "\n" "// d2"]`,
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
			expect: `GenDecl [Empty line before]
Field [New line before] [New line after]
Field [New line before] [New line after]
ReturnStmt [New line before] [New line after]`,
		},
		{
			name: "multi-line-string",
			code: `package a

				var a = b{
					c: ` + "`" + `
` + "`" + `,
				}`,
			expect: `GenDecl [Empty line before]
KeyValueExpr [New line before] [New line after]`,
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
			expect: `FuncDecl [Empty line before]
SwitchStmt [New line before] [New line after]
CaseClause [New line before] [End "\n" "// a"] [New line after]
CaseClause [New line before] [Start "// b"] [New line after]
CaseClause [New line before] [Start "// c"] [New line after]`,
		},
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
CaseClause [New line before] [Colon "\n" "// b" "// c"] [New line after]
DeclStmt [Empty line before]`,
		},
		{
			name: "file",
			code: `/*Start*/ package /*Package*/ postests /*Name*/

			var i int`,
			expect: `File [Start "/*Start*/"] [Package "/*Package*/"] [Name "/*Name*/"]
GenDecl [Empty line before]`,
		},
		{

			name: "TypeAssertExpr",
			code: `package main
			
			// TypeAssertExpr
			var I interface{}
			var J = I. /*TypeAssertExprX*/ ( /*TypeAssertExprLparen*/ int /*TypeAssertExprType*/)`,
			expect: `GenDecl [Empty line before] [Start "// TypeAssertExpr"] [New line after]
GenDecl [New line before]
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
			expect: `FuncDecl [Empty line before]
RangeStmt [New line before] [Start "/*Start*/" "\n"] [For "/*For*/"] [Key "/*Key*/"] [Value "/*Value*/"] [Range "/*Range*/"] [X "/*X*/"] [End "/*End*/"] [New line after]
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

				type Y /*Start*/ chan /*Begin*/ <- /*Arrow*/ int /*End*/`,
			expect: `GenDecl [Empty line before] [End "/*End*/"]
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
			expect: `FuncDecl [Empty line before]
IfStmt [New line before] [New line after]
BlockStmt [Lbrace "\n" "// a"]`,
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
IncDecStmt [New line before] [X "/* foo */"] [New line after]`,
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
BlockStmt [Lbrace "// foo"]
ExprStmt [New line before] [New line after]`,
		},
		{
			name: "comment after func",
			code: `package main
			
			func /* foo */ main() {
				i
			}`,
			expect: `FuncDecl [Empty line before] [Func "/* foo */"]
ExprStmt [New line before] [New line after]`,
		},
		{
			name: "field",
			code: `package main

			type A struct {
				A /*IdentEnd*/ int /*FieldType*/ ` + "`" + `a:"a"` + "`" + `
			}`,
			expect: `GenDecl [Empty line before]
Field [New line before] [Type "/*FieldType*/"] [New line after]
Ident [End "/*IdentEnd*/"]`,
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
			func /*FuncDeclDoc*/ (a *b) /*FuncDeclRecv*/ c /*FuncDeclName*/ (d, e int) /*FuncDeclParams*/ (f, g int) /*FuncDeclType*/ {
			}`,
			expect: `FuncDecl [Empty line before] [Start "// FuncDecl"] [Func "/*FuncDeclDoc*/"] [Recv "/*FuncDeclRecv*/"] [Name "/*FuncDeclName*/"] [Params "/*FuncDeclParams*/"] [Results "/*FuncDeclType*/"]
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
			file, err := Decorate(fset, f)
			if err != nil {
				t.Fatal(err)
			}

			buf := &bytes.Buffer{}
			debug(buf, file)

			if normalize(test.expect) != normalize(buf.String()) {
				t.Errorf("diff:\n%s", diff(normalize(test.expect), normalize(buf.String())))
			}

		})
	}
}

func TestParseFile_Comments(t *testing.T) {
	code := `package a

		// a
		func main(){}`

	f, err := ParseFile(token.NewFileSet(), "", code, 0)
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	if err := Fprint(buf, f); err != nil {
		panic(err)
	}
	compareSrc(t, code, buf.String())
}

func TestBad(t *testing.T) {
	tests := []struct {
		name, code string
	}{
		{
			name: "decl",
			code: `package a

				%BADDECL%
			`,
		},
		{
			name: "stmt",
			code: `package a

				func a() {
					%BADSTMT%
				}
			`,
		},
		{
			name: "expr",
			code: `package a

				func a() {
					var a = %BADEXPR%
				}
			`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fset := token.NewFileSet()
			af, err := parser.ParseFile(fset, "", test.code, parser.ParseComments)
			if err == nil {
				t.Fatal("expected error, found none")
			}
			abuf := &bytes.Buffer{}
			if err := format.Node(abuf, fset, af); err != nil {
				t.Fatal(err)
			}

			df, err := ParseFile(token.NewFileSet(), "", test.code, parser.ParseComments)
			if err == nil {
				t.Fatal("expected error, found none")
			}

			dbuf := &bytes.Buffer{}
			if err := Fprint(dbuf, df); err != nil {
				t.Fatal(err)
			}

			compare(t, abuf.String(), dbuf.String())
		})
	}
}

func TestDecorator_ParseDir(t *testing.T) {

	code := map[string]string{
		"a.go": `package a

		// a
		func a(){}`,
		"b.go": `package a

		// b
		func b(){}`,
	}
	dir, err := tempDir(code)

	pkg, err := ParseDir(token.NewFileSet(), dir, nil, 0)
	if err != nil {
		panic(err)
	}
	p := pkg["a"]

	if len(pkg) != 1 {
		t.Fatalf("expected 1 package, found %d", len(pkg))
	}

	actual := map[string]string{}
	for fpath, file := range p.Files {
		_, fname := filepath.Split(fpath)
		buf := &bytes.Buffer{}
		if err := Fprint(buf, file); err != nil {
			t.Fatal(err)
		}
		actual[fname] = buf.String()
	}

	compareDir(t, dir, actual)

}

var multiSpaces = regexp.MustCompile(" {2,}")

func normalize(s string) string {
	s = multiSpaces.ReplaceAllString(s, "")
	s = strings.Replace(s, "\t", "", -1)
	s = strings.TrimSpace(s)
	return s
}
