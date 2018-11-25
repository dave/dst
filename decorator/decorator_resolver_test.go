package decorator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil/dummy"
	"golang.org/x/tools/go/packages"
)

func TestDecoratorResolver(t *testing.T) {
	type tc struct {
		expect string
		get    func(*dst.File) *dst.Ident
	}
	tests := []struct {
		skip, solo bool
		name       string
		src        dummy.Dir
		cases      []tc
	}{
		{
			name: "simple",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

						import (
							"root/a"
							. "root/b"
						)

						func main(){
							a.A()
							B()
							C()
						}
					`),
					"c.go": dummy.Src("package main\n\nfunc C(){}"),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					"root/a",
					func(f *dst.File) *dst.Ident {
						d := f.Decls[1]
						return d.(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.SelectorExpr).Sel
					},
				},
				{
					"root/b",
					func(f *dst.File) *dst.Ident {
						return f.Decls[1].(*dst.FuncDecl).Body.List[1].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
					},
				},
				{
					"",
					func(f *dst.File) *dst.Ident {
						return f.Decls[1].(*dst.FuncDecl).Body.List[2].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
					},
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

			root := dummy.TempDir(test.src)

			pkgs, err := packages.Load(
				&packages.Config{
					Mode: packages.LoadSyntax,
					Dir:  filepath.Join(root, "main"),
				},
				"root/main",
			)
			os.RemoveAll(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(pkgs) != 1 {
				t.Fatalf("expected 1 package, found %d", len(pkgs))
			}
			pkg := pkgs[0]

			d := NewWithImports(pkg)

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

			for _, c := range test.cases {
				id := c.get(file)
				if id.Path != c.expect {
					t.Errorf("expected %q, found %q", c.expect, id.Path)
				}
			}

		})
	}
}
