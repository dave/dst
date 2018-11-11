package decorator

import (
	"path/filepath"
	"testing"

	"bytes"
	"fmt"
	"go/format"

	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/dave/dst/decorator/resolver/gotypes"
	"github.com/dave/dst/dstutil/dummy"
	"golang.org/x/tools/go/packages"
)

func TestRestorerResolver(t *testing.T) {
	type tc struct {
		mutation func(f *dst.File)
		expect   string
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

						func main(){}
					`),
					"c.go": dummy.Src("package main\n\nfunc C(){}"),
				},
				"a":      dummy.Dir{"a.go": dummy.Src("package a \n\n func A(){}")},
				"go.mod": dummy.Src("module root"),
			},
			cases: []tc{
				{
					func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a", Name: "A"}}})
					},
					``,
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

			for _, c := range test.cases {

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

				c.mutation(file)

				r := NewRestorer()
				r.Resolver = &gopackages.PackageResolver{
					Config: cfg,
				}
				pr := r.NewPackageRestorer(mainPkg, mainDir)
				restoredFile := pr.RestoreFile("main.go", file)

				buf := &bytes.Buffer{}
				if err := format.Node(buf, r.Fset, restoredFile); err != nil {
					t.Fatal(err)
				}

				fmt.Println(buf.String())

			}

		})
	}
}
