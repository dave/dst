package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/andreyvit/diff"
)

func TestFragger(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		code       string
		expect     string
	}{
		{
			name: "empty func",
			code: `package a

				func b() {
				} // c
				`,
			expect: `File Start 1:1
				File "package" 1:1
				File AfterPackage 1:8
				Ident Start 1:9
				Ident "a" 1:9
				Ident End 1:10
				File AfterName 1:10
				"\n" 1:10
				"\n" 2:1
				FuncDecl Start 3:1
				FuncDecl "func" 3:1
				FuncDecl AfterFunc 3:5
				Ident Start 3:6
				Ident "b" 3:6
				Ident End 3:7
				FuncDecl AfterName 3:7
				FieldList Start 3:7
				FieldList "(" 3:7
				FieldList AfterOpening 3:8
				FieldList ")" 3:8
				FieldList End 3:9
				FuncDecl AfterParams 3:9
				BlockStmt Start 3:10
				BlockStmt "{" 3:10
				BlockStmt AfterLbrace 3:11
				"\n" 3:11
				BlockStmt "}" 4:1
				BlockStmt End 4:2
				FuncDecl End 4:2
				"// c" 4:3`,
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

			p := &Fragger{}
			p.Fragment(fset, f)

			buf := &bytes.Buffer{}
			p.debug(fset, buf)

			if test.expect == "" {
				t.Error(buf.String())
			} else if normalize(buf.String()) != normalize(test.expect) {
				t.Errorf("diff: %s", diff.LineDiff(normalize(buf.String()), normalize(test.expect)))
			}
		})
	}
}
