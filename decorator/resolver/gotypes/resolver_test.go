package gotypes_test

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/dst/decorator/resolver/gotypes"
	"golang.org/x/tools/go/packages"
)

func TestDecoratorResolver(t *testing.T) {
	type tc struct{ id, expect string }
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
					)

					func main(){
						a.A()
					}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
			cases: []tc{
				{"A", "root/a"},
			},
		},
		{
			name: "non-qualified-ident",
			src: map[string]string{
				"main/main.go": `package main

					import (
						"root/a"
					)

					func main(){
						t.A()
					}

					var t a.T`,
				"a/a.go": "package a \n\n type T struct{} \n\n func (T)A(){}",
				"go.mod": "module root",
			},
			cases: []tc{
				{"A", ""},
			},
		},
		{
			name: "field",
			src: map[string]string{
				"main/main.go": `package main

					import (
						"root/a"
					)

					func main(){
						t := a.T{
							B: 0,
						}
					}`,
				"a/a.go": "package a \n\n type T struct{B int}",
				"go.mod": "module root",
			},
			cases: []tc{
				{"B", ""},
			},
		},
		{
			name: "more",
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
				{"A", "root/a"},
				{"B", "root/b"},
				{"C", "root/main"},
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
			os.RemoveAll(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(pkgs) != 1 {
				t.Fatalf("expected 1 package, found %d", len(pkgs))
			}
			pkg := pkgs[0]

			res := gotypes.New(pkg.TypesInfo.Uses)

			parents := map[string]ast.Node{}
			parentFields := map[string]string{}
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
						parentFields[n.Sel.Name] = "Sel"
					case *ast.Ident:
						if _, ok := nodes[n.Name]; !ok {
							nodes[n.Name] = n
							parents[n.Name] = nil
							parentFields[n.Name] = ""
						}
					}
					return true
				})
			}

			for _, c := range test.cases {
				//ast.Print(pkg.Fset, parents[c.id])
				//ast.Print(pkg.Fset, nodes[c.id])
				path, err := res.ResolveIdent(nil, parents[c.id], parentFields[c.id], nodes[c.id])
				if err != nil {
					t.Error(err)
				}
				if path != c.expect {
					t.Errorf("expect %q, found %q", c.expect, path)
				}
			}

		})
	}
}
