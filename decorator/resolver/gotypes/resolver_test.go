package gotypes_test

import (
	"os"
	"testing"

	"path/filepath"

	"go/ast"

	"github.com/dave/dst/decorator/resolver/gotypes"
	"github.com/dave/dst/dstutil/dummy"
	"golang.org/x/tools/go/packages"
)

func TestNodeResolver(t *testing.T) {
	type tc struct{ id, expect string }
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
						)

						func main(){
							a.A()
						}
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{"A", "root/a"},
			},
		},
		{
			name: "more",
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
				{"A", "root/a"},
				{"B", "root/b"},
				{"C", ""},
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

			res := &gotypes.IdentResolver{
				Path:       pkg.PkgPath,
				Uses:       pkg.TypesInfo.Uses,
				Selections: pkg.TypesInfo.Selections,
			}

			nodes := map[string]*ast.Ident{}
			for _, f := range pkg.Syntax {
				ast.Inspect(f, func(n ast.Node) bool {
					// TODO: Only handles idents in CallExpr - extend to any node?
					switch n := n.(type) {
					case *ast.CallExpr:
						switch n := n.Fun.(type) {
						case *ast.SelectorExpr:
							nodes[n.Sel.Name] = n.Sel
						case *ast.Ident:
							if _, ok := nodes[n.Name]; !ok {
								nodes[n.Name] = n
							}
						}
					}
					return true
				})
			}

			for _, c := range test.cases {
				n, ok := nodes[c.id]
				if !ok {
					t.Errorf("node not found for %q", c.id)
				}
				found := res.ResolveIdent(n)
				if found != c.expect {
					t.Errorf("expect %q, found %q", c.expect, found)
				}
			}

		})
	}
}
