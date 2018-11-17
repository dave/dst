package decorator

import (
	"bytes"
	"context"
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

/*
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
*/

func TestRestorerResolver(t *testing.T) {
	type tc struct {
		skip, solo bool
		name, desc string
		mutate     func(f *dst.File)
		restorer   func(r *FileRestorer)
		expect     string
	}
	tests := []struct {
		name  string
		root  string // root package path - default "root"
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
					desc: "adding a conflicting import renames correctly, and in correct order",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(
							b.List,
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/c/c", Name: "C"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a/c", Name: "A"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b/c", Name: "B"}}},
						)
					},
					expect: `package main
            
            			import (
			            	"root/a/c"
            				c1 "root/b/c"
            				c2 "root/c/c"
            			)
            
            			func main() { c2.C(); c.A(); c1.B() }`,
				},
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
					name: "change-to-anon-remove-use",
					desc: "changing from standard to anonymous import works",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "_"
					},
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = nil
					},
					expect: `package main
            
            			import _ "root/a"
            
            			func main() {}`,
				},
				{
					name: "change-to-anon-still-in-use",
					desc: "changing from standard to anonymous import has no effect if the package is still in use",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "_"
					},
					expect: `package main
            
            			import "root/a"
            
            			func main() { a.A() }`,
				},
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
					name: "delete-all",
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
			name: "two-existing",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

            			import (
							"root/a"
							"root/b"
						)

            			func main() { a.A(); b.B() }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "delete-one",
					desc: "blocks changing from >1 to 1 imports correctly lose parens",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = b.List[1:2]
					},
					expect: `package main
            
            			import "root/b"
            
            			func main() { b.B() }`,
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
				"e":      dummy.Dir{"e.go": dummy.Src("package e \n\n func E(){}")},
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
						b.List = b.List[2:4]
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/e", Name: "E"}}})
					},
					expect: `package main
						
						import "root/e"

						import (
							"root/c"
							"root/d"
						)

						func main() { c.C(); d.D(); e.E() }`,
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
			name: "first-block-spacing",
			root: "foo.bar",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

						import (

							"foo.bar/a"

							"fmt"
							
							"bytes"

						)

            			func main() { a.A(); fmt.Print(); bytes.Title(); }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"go.mod": dummy.Src("module foo.bar"),
			},
			cases: []tc{
				{
					solo: true,
					name: "block-reordered-spacing-fixed",
					desc: "line-feed formatting in re-arranged first block is correctly modified",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "io", Name: "Copy"}}})
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "foo.bar/b", Name: "B"}}})
					},
					expect: `package main
            
            			import (
            				"bytes"
            				"fmt"
            				"io"
            
            				"foo.bar/a"
            				"foo.bar/b"
            			)
            
            			func main() { a.A(); fmt.Print(); bytes.Title(); io.Copy(); b.B() }`,
				},
				{

					name: "block-reordered-spacing-fixed-delete-first-non-std",
					desc: "when we delete the first non-std-lib import, the line-spacing is correct",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = b.List[1:]
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "foo.bar/b", Name: "B"}}})
					},
					expect: `package main
            
            			import (
            				"bytes"
            				"fmt"
            
            				"foo.bar/b"
            			)
            
            			func main() { fmt.Print(); bytes.Title(); b.B() }`,
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
		{
			name: "two-blocks-with-alias",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

            			import (
							aa "root/a"
							bb "root/b"
						)

						import (
							cc "root/c"
							dd "root/d"
						)

            			func main() { aa.A(); bb.B(); cc.C(); dd.D(); }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"c":      dummy.Dir{"c.go": dummy.Src("package c \n\n func C(){}")},
				"d":      dummy.Dir{"d.go": dummy.Src("package d \n\n func D(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "change-alias",
					desc: "changing alias correctly changes import block and all idents",
					restorer: func(r *FileRestorer) {
						r.Alias["root/b"] = "bbb"
						r.Alias["root/d"] = "ddd"
					},
					expect: `package main
            
            			import (
            				aa "root/a"
            				bbb "root/b"
            			)
            
			            import (
            				cc "root/c"
            				ddd "root/d"
            			)
            
            			func main() { aa.A(); bbb.B(); cc.C(); ddd.D() }`,
				},
			},
		},
		{
			name: "dot-imports",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

            			import (
							. "root/a"
							"root/b"
							cc "root/c"
						)

            			func main() { A(); b.B(); cc.C(); }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"b.go": dummy.Src("package b \n\n func B(){}")},
				"c":      dummy.Dir{"c.go": dummy.Src("package c \n\n func C(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "change-to-normal",
					desc: "ensure changing between dot-import, normal and alias import works correctly",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = ""
						r.Alias["root/b"] = ""
						r.Alias["root/c"] = ""
					},
					expect: `package main
            
            			import (
            				"root/a"
            				"root/b"
            				"root/c"
            			)
            
            			func main() { a.A(); b.B(); c.C() }`,
				},
				{
					name: "change-to-dot",
					desc: "ensure changing between dot-import, normal and alias import works correctly",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "."
						r.Alias["root/b"] = "."
						r.Alias["root/c"] = "."
					},
					expect: `package main
            
            			import (
            				. "root/a"
            				. "root/b"
            				. "root/c"
            			)
            
            			func main() { A(); B(); C() }`,
				},
				{
					name: "change-to-alias",
					desc: "ensure changing between dot-import, normal and alias import works correctly",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "aa"
						r.Alias["root/b"] = "bb"
						r.Alias["root/c"] = "cc"
					},
					expect: `package main
            
            			import (
            				aa "root/a"
            				bb "root/b"
            				cc "root/c"
            			)
            
            			func main() { aa.A(); bb.B(); cc.C() }`,
				},
			},
		},
		{
			name: "conflict",
			src: dummy.Dir{
				"main": dummy.Dir{
					"main.go": dummy.Src(`package main

            			import (
							"root/a"
						)

            			func main() { a.A(); }
					`),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"b":      dummy.Dir{"a": dummy.Dir{"a.go": dummy.Src("package a \n\n func AA(){}")}},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					name: "dot-import-conflict-check-disabled",
					desc: "no conflict checking for dot-imports",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "."
					},
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b/a", Name: "AA"}}})
					},
					expect: `package main
            
            			import (
            				. "root/a"
            				"root/b/a"
            			)
            
            			func main() { A(); a.AA() }`,
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
			if solo && !c.solo {
				continue // TODO: remove
			}
			t.Run(test.name+"/"+c.name, func(t *testing.T) {
				if solo && !c.solo {
					t.Skip()
				}
				if c.skip {
					t.Skip()
				}

				rootDir := dummy.TempDir(test.src)
				defer os.RemoveAll(rootDir)
				mainDir := filepath.Join(rootDir, "main")
				mainPkg := "root/main"
				if test.root != "" {
					mainPkg = fmt.Sprintf("%s/main", test.root)
				}

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
				d.Resolver = gotypes.FromPackage(pkg)

				var file *dst.File
				for _, sf := range pkg.Syntax {
					if _, name := filepath.Split(pkg.Fset.File(sf.Pos()).Name()); name == "main.go" {
						file = d.Decorate(sf).(*dst.File)
						break
					}
				}

				if c.mutate != nil {
					c.mutate(file)
				}

				r := NewRestorer()
				r.Resolver = &resolver.Guess{}
				pr := r.NewPackageRestorer(mainPkg, mainDir)
				fr := pr.NewFileRestorer("main.go", file)

				if c.restorer != nil {
					c.restorer(fr)
				}

				restoredFile := fr.RestoreFile(context.Background())

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
