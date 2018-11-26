package gotypes_test

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/dst/decorator/resolver/gotypes"
	"github.com/dave/dst/dstutil/dummy"
	"golang.org/x/tools/go/packages"
)

func TestIdentResolver(t *testing.T) {
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
			name: "non-qualified-ident",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

						import (
							"root/a"
						)

						func main(){
							t.A()
						}

						var t a.T
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n type T struct{} \n\n func (T)A(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{"A", ""},
			},
		},
		{
			name: "field",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

						import (
							"root/a"
						)

						func main(){
							t := a.T{
								B: 0,
							}
						}
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n type T struct{B int}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{"B", ""},
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
				Path: "root/main",
				Info: pkg.TypesInfo,
			}

			parents := map[string]ast.Node{}
			nodes := map[string]*ast.Ident{}
			for _, f := range pkg.Syntax {
				_, fname := filepath.Split(pkg.Fset.File(f.Pos()).Name())
				if fname != "main.go" {
					continue
				}
				ast.Inspect(f, func(n ast.Node) bool {
					switch n := n.(type) {
					case *ast.SelectorExpr:
						nodes[n.Sel.Name] = n.Sel
						parents[n.Sel.Name] = n
					case *ast.Ident:
						if _, ok := nodes[n.Name]; !ok {
							nodes[n.Name] = n
							parents[n.Name] = nil
						}
					}
					return true
				})
			}

			for _, c := range test.cases {
				//ast.Print(pkg.Fset, parents[c.id])
				//ast.Print(pkg.Fset, nodes[c.id])
				found, err := res.ResolveIdent(nil, parents[c.id], nodes[c.id])
				if err != nil {
					t.Error(err)
				}
				if found != c.expect {
					t.Errorf("expect %q, found %q", c.expect, found)
				}
			}

		})
	}
}
