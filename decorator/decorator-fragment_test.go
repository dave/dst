package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestFragment(t *testing.T) {
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
			expect: `File Start 1:1
            File "package" 1:1
            File Package 1:8
            Ident Start 1:9
			Ident X 1:9
            Ident "a" 1:9
            Ident End 1:10
            File Name 1:10
            Empty line 1:10
            "/*\n\tfoo\n*/" 3:1
            New line 5:3
            GenDecl Start 6:1
            GenDecl "var" 6:1
            GenDecl Tok 6:4
            ValueSpec Start 6:5
            Ident Start 6:5
			Ident X 6:5
            Ident "i" 6:5
            Ident End 6:6
            Ident Start 6:7
			Ident X 6:7
            Ident "int" 6:7
            Ident End 6:10
            ValueSpec End 6:10
            GenDecl End 6:10`,
		},
		{
			name: "case clause",
			code: `package a

				func main() {
					switch {
					default:
						// b
						// c

						var i int
					}
				}`,
			expect: `File Start 1:1
            File "package" 1:1
            File Package 1:8
            Ident Start 1:9
			Ident X 1:9
            Ident "a" 1:9
            Ident End 1:10
            File Name 1:10
            Empty line 1:10
            FuncDecl Start 3:1
            FuncDecl "func" 3:1
            FuncDecl Func 3:5
            Ident Start 3:6
			Ident X 3:6
            Ident "main" 3:6
            Ident End 3:10
            FuncDecl Name 3:10
            FieldList Start 3:10
            FieldList "(" 3:10
            FieldList Opening 3:11
            FieldList ")" 3:11
            FieldList End 3:12
            FuncDecl Params 3:12
            BlockStmt Start 3:13
            BlockStmt "{" 3:13
            BlockStmt Lbrace 3:14
            New line 3:14
            SwitchStmt Start 4:2
            SwitchStmt "switch" 4:2
            SwitchStmt Switch 4:8
            BlockStmt Start 4:9
            BlockStmt "{" 4:9
            BlockStmt Lbrace 4:10
            New line 4:10
            CaseClause Start 5:2
            CaseClause "default" 5:2
            CaseClause Case 5:9
            CaseClause ":" 5:9
            CaseClause Colon 5:10
            New line 5:10
            "// b" 6:3
            New line 6:7
            "// c" 7:3
            Empty line 7:7
            DeclStmt Start 9:3
            GenDecl Start 9:3
            GenDecl "var" 9:3
            GenDecl Tok 9:6
            ValueSpec Start 9:7
            Ident Start 9:7
			Ident X 9:7
            Ident "i" 9:7
            Ident End 9:8
            Ident Start 9:9
			Ident X 9:9
            Ident "int" 9:9
            Ident End 9:12
            ValueSpec End 9:12
            GenDecl End 9:12
            DeclStmt End 9:12
			CaseClause End 9:12
            New line 9:12
            BlockStmt "}" 10:2
            BlockStmt End 10:3
            SwitchStmt End 10:3
            New line 10:3
            BlockStmt "}" 11:1
            BlockStmt End 11:2
            FuncDecl End 11:2`,
		},
		{
			name: "empty func",
			code: `package a

				func b() {
					var d int
					// c
					var e int
				}
				`,
			expect: `File Start 1:1
            File "package" 1:1
            File Package 1:8
            Ident Start 1:9
			Ident X 1:9
            Ident "a" 1:9
            Ident End 1:10
            File Name 1:10
            Empty line 1:10
            FuncDecl Start 3:1
            FuncDecl "func" 3:1
            FuncDecl Func 3:5
            Ident Start 3:6
			Ident X 3:6
            Ident "b" 3:6
            Ident End 3:7
            FuncDecl Name 3:7
            FieldList Start 3:7
            FieldList "(" 3:7
            FieldList Opening 3:8
            FieldList ")" 3:8
            FieldList End 3:9
            FuncDecl Params 3:9
            BlockStmt Start 3:10
            BlockStmt "{" 3:10
            BlockStmt Lbrace 3:11
            New line 3:11
            DeclStmt Start 4:2
            GenDecl Start 4:2
            GenDecl "var" 4:2
            GenDecl Tok 4:5
            ValueSpec Start 4:6
            Ident Start 4:6
			Ident X 4:6
            Ident "d" 4:6
            Ident End 4:7
            Ident Start 4:8
			Ident X 4:8
            Ident "int" 4:8
            Ident End 4:11
            ValueSpec End 4:11
            GenDecl End 4:11
            DeclStmt End 4:11
            New line 4:11
            "// c" 5:2
            New line 5:6
            DeclStmt Start 6:2
            GenDecl Start 6:2
            GenDecl "var" 6:2
            GenDecl Tok 6:5
            ValueSpec Start 6:6
            Ident Start 6:6
			Ident X 6:6
            Ident "e" 6:6
            Ident End 6:7
            Ident Start 6:8
			Ident X 6:8
            Ident "int" 6:8
            Ident End 6:11
            ValueSpec End 6:11
            GenDecl End 6:11
            DeclStmt End 6:11
            New line 6:11
            BlockStmt "}" 7:1
            BlockStmt End 7:2
            FuncDecl End 7:2`,
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

			d := NewDecorator(fset)

			fd := d.newFileDecorator()
			fd.fragment(f)

			buf := &bytes.Buffer{}
			fd.debug(buf)

			if test.expect == "" {
				t.Error(buf.String())
			} else if normalize(buf.String()) != normalize(test.expect) {
				t.Errorf("diff:\n%s", diff(normalize(test.expect), normalize(buf.String())))
			}
		})
	}
}

func diff(expect, found string) string {
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(dmp.DiffMain(expect, found, false))
}
