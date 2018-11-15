package decorator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/gotypes"
	"github.com/dave/dst/dstutil/dummy"
	"golang.org/x/tools/go/packages"
)

func TestRestorerResolver(t *testing.T) {
	type tc struct {
		skip, solo bool
		name, desc string
		mutate     func(f *dst.File)
		expect     string
	}
	tests := []struct {
		name  string
		src   dummy.Dir
		cases []tc
	}{
		{
			name: "simple",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

						func main(){}
					`),
					"c.go": dummy.Src("package main\n\nfunc C(){}"),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "add-one",
					desc: "adding an import to a file that has no imports creates a new import block",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a", Name: "A"}}})
					},
					expect: `package main
            		import "root/a"
            		func main() { a.A() }`,
				},
				{
					name: "add-two",
					desc: "adding two imports to a file that has no imports creates a new import block with parens",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(
							b.List,
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a", Name: "A"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b", Name: "B"}}},
						)
					},
					expect: `package main
					import (
            			"root/a"
            			"root/b"
            		)
            		func main() { a.A(); b.B() }`,
				},
				{
					name: "conflict",
					desc: "adding a conflicting import renames correctly",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(
							b.List,
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a/c", Name: "A"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b/c", Name: "B"}}},
						)
					},
					expect: `package main
            			import (
            				"root/a/c"
            				c1 "root/b/c"
            			)
            			func main() { c.A(); c1.B() }`,
				},
				{desc: "package names are corrected"},
				{desc: "alias imports are retainied correctly"},
				{desc: "don't remove anon imports"},
				{desc: "don't remove C import"},
				{desc: "only adds to first import block"},
				{desc: "removes from all imports block"},
				{desc: "only re-orders first block"},
				{desc: "doesn't re-order first import block if no additions (check with a deletion)"},
				{desc: "re-orders first import block correctly when adding"},
				{desc: "convert anonymous import to standard"},
				{desc: "anon imports manually added with Alias"},
				{desc: "conflicts are resolved in correct order"},
				{desc: "manually added alias work as expected"},
				{desc: "manually added alias take priority over alias in imports block"},
				{desc: "alias from import block works correctly"},
				{desc: "line-feed formatting in re-agganged first block is correctly modified"},
				{desc: "no conflict checking for dot-imports"},
				{desc: "changing between standard and dot-import correctly changes import block and all idents (and decorations are merged)"},
				{desc: "changing alias correctly changes import block and all idents"},
				{desc: "changing from standard to anonymous import works"},
				{desc: "blocks changing from >1 to 1 imports correctly lose parens"},
			},
		},
		{
			name: "existing",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main
            			import "root/a"
            			func main() { a.A() }
					`),
					"c.go": dummy.Src("package main\n\nfunc C(){}"),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "add",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b", Name: "B"}}})
					},
					expect: `package main
            			import (
            				"root/a"
            				"root/b"
            			)
            			func main() { a.A(); b.B() }`,
				},
				{
					name: "delete",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = nil
					},
					expect: `package main
            			func main() {}`,
				},
			},
		},
		{
			name: "two-blocks",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main
            			import (
							"root/a"
							"root/b"
						)
						import (
							"root/c"
							"root/d"
						)
            			func main() { a.A(); b.B(); c.C(); d.D(); }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"c":      dummy.Dir{"c.go": dummy.Src("package c \n\n func C(){}")},
				"d":      dummy.Dir{"d.go": dummy.Src("package d \n\n func D(){}")},
				"e":      dummy.Dir{"d.go": dummy.Src("package d \n\n func D(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "block-deleted",
					desc: "blocks are deleted ok",
					mutate: func(f *dst.File) {
						b := f.Decls[2].(*dst.FuncDecl).Body
						b.List = b.List[0:2]
					},
					expect: `package main
						import (
							"root/a"
							"root/b"
						)
						func main() { a.A(); b.B() }`,
				},
				{
					name: "block-deleted-ad-added",
					desc: "if all imports are removed from first block and one added, it's ok",
					mutate: func(f *dst.File) {
						b := f.Decls[2].(*dst.FuncDecl).Body
						b.List = b.List[0:2]
					},
					expect: `package main
						import (
							"root/a"
							"root/b"
						)
						func main() { a.A(); b.B() }`,
				},
			},
		},
		{
			name: "first-block-decorated",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main
						import (
							// before c
							"root/c" // after c
							// before a
							"root/a" // after a
						)
            			func main() { a.A(); c.C(); }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"c":      dummy.Dir{"c.go": dummy.Src("package c \n\n func C(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "decorations-retained",
					desc: "decorations in re-ordered block are retained",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b", Name: "B"}}})
					},
					expect: `package main
            			import (
            				// before a
            				"root/a" // after a
            				"root/b"
            				// before c
            				"root/c" // after c
            			)
            			func main() { a.A(); c.C(); b.B() }`,
				},
			},
		},
		{
			name: "existing-alias-package-name",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main
            			import a "root/a"
            			func main() { a.A() }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "alias-retained",
					desc: "alias is retained even when alias == package name",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b", Name: "B"}}})
					},
					expect: `package main
            			import (
            				a "root/a"
            				"root/b"
            			)
            			func main() { a.A(); b.B() }`,
				},
			},
		},
	}
	var solo bool
	for _, test := range tests {
		for _, c := range test.cases {
			if c.solo {
				solo = true
				break
			}
			if solo {
				break
			}
		}
	}
	for _, test := range tests {
		for _, c := range test.cases {
			t.Run(test.name+"/"+c.name, func(t *testing.T) {
				if solo && !c.solo {
					t.Skip()
				}
				if c.skip {
					t.Skip()
				}
				if c.expect == "" {
					t.Skip()
				}

				rootDir := dummy.TempDir(test.src)
				defer os.RemoveAll(rootDir)
				mainDir := filepath.Join(rootDir, "main")
				mainPkg := "root/main"

				cfg := &packages.Config{
					Mode: packages.LoadSyntax,
					Dir:  mainDir,
				}

				pkgs, err := packages.Load(cfg, mainPkg)
				if err != nil {
					t.Fatal(err)
				}
				if len(pkgs) != 1 {
					t.Fatalf("expected 1 package, found %d", len(pkgs))
				}
				pkg := pkgs[0]

				d := New(pkg.Fset)
				d.Resolver = &gotypes.IdentResolver{
					Path:       pkg.PkgPath,
					Uses:       pkg.TypesInfo.Uses,
					Selections: pkg.TypesInfo.Selections,
				}

				var file *dst.File
				for _, sf := range pkg.Syntax {
					if _, name := filepath.Split(pkg.Fset.File(sf.Pos()).Name()); name == "main.go" {
						file = d.Decorate(sf).(*dst.File)
						break
					}
				}

				c.mutate(file)

				r := NewRestorer()
				r.Resolver = &resolver.Guess{}
				pr := r.NewPackageRestorer(mainPkg, mainDir)
				restoredFile := pr.RestoreFile("main.go", file)

				buf := &bytes.Buffer{}
				if err := format.Node(buf, r.Fset, restoredFile); err != nil {
					t.Fatal(err)
				}

				expected, err := format.Source([]byte(c.expect))
				if err != nil {
					panic(err)
				}

				if buf.String() != string(expected) {
					t.Errorf("expected: %s\ngot: %s", string(expected), buf.String())
					t.Errorf("diff: %s", diff.LineDiff(string(expected), buf.String()))
				}

			})
		}
	}
}

func TestPackageOrder(t *testing.T) {
	paths := []string{"C", "a.b/d", "a.b/c", "fmt", "bytes", "a/b"}
	sort.Slice(paths, func(i, j int) bool {
		return packagePathOrderLess(paths[i], paths[j])
	})
	expect := "[a/b bytes fmt a.b/c a.b/d C]"
	found := fmt.Sprint(paths)
	if found != expect {
		t.Errorf("expect %s, found %s", expect, found)
	}
}
