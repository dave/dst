package decorator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/dst"
	"golang.org/x/tools/go/packages"
)

func TestDecoratorResolver(t *testing.T) {
	type tc struct {
		expect           string
		get              func(*dst.File) *dst.Ident
		resolveLocalPath bool
	}
	tests := []struct {
		skip, solo bool
		name       string
		src        map[string]string
		cases      []tc
	}{
		{
			name: "simple",
			src: map[string]string{
				"main/main.go": `package main

					import (
						"root/a"
						. "root/b"
					)

					func main(){
						a.A()
						B()
						C()
					}`,
				"main/c.go": "package main \n\n func C(){}",
				"a/a.go":    "package a \n\n func A(){}",
				"b/b.go":    "package b \n\n func B(){}",
				"go.mod":    "module root",
			},
			cases: []tc{
				{
					"root/a",
					func(f *dst.File) *dst.Ident {
						d := f.Decls[1]
						return d.(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
					},
					false,
				},
				{
					"root/b",
					func(f *dst.File) *dst.Ident {
						return f.Decls[1].(*dst.FuncDecl).Body.List[1].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
					},
					false,
				},
				{
					"",
					func(f *dst.File) *dst.Ident {
						return f.Decls[1].(*dst.FuncDecl).Body.List[2].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
					},
					false,
				},
				{
					"root/main",
					func(f *dst.File) *dst.Ident {
						return f.Decls[1].(*dst.FuncDecl).Body.List[2].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
					},
					true,
				},
			},
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

			root, err := tempDir(test.src)
			if err != nil {
				t.Fatal(err)
			}

			pkgs, err := packages.Load(
				&packages.Config{
					Mode: packages.LoadSyntax,
					Dir:  filepath.Join(root, "main"),
				},
				"root/main",
			)
			_ = os.RemoveAll(root) // ignore error
			if err != nil {
				t.Fatal(err)
			}
			if len(pkgs) != 1 {
				t.Fatalf("expected 1 package, found %d", len(pkgs))
			}
			pkg := pkgs[0]

			for _, c := range test.cases {
				d := NewDecoratorFromPackage(pkg)
				d.ResolveLocalPath = c.resolveLocalPath

				var file *dst.File
				for _, sf := range pkg.Syntax {
					if _, name := filepath.Split(pkg.Fset.File(sf.Pos()).Name()); name == "main.go" {
						var err error
						file, err = d.DecorateFile(sf)
						if err != nil {
							t.Fatal(err)
						}
						break
					}
				}

				id := c.get(file)
				if id.Path != c.expect {
					t.Errorf("expected %q, found %q", c.expect, id.Path)
				}
			}

		})
	}
}
