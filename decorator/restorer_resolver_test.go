package decorator

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"golang.org/x/tools/go/packages"
)

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
		src   map[string]string
		cases []tc
	}{
		{
			name: "simple",
			root: "a.b",
			src: map[string]string{
				"main/main.go": `package main 

					func main(){}`,
				"a/a.go":       "package a \n\n func A(){}",
				"b/b.go":       "package b \n\n func B(){}",
				"fmt/a/fmt.go": "package fmt \n\n func A(){}",
				"fmt/b/fmt.go": "package fmt \n\n func B(){}",
				"fmt/c/fmt.go": "package fmt \n\n func C(){}",
				"go.mod":       "module a.b",
			},
			cases: []tc{
				{
					name: "add-anon",
					desc: "adding an anonymous import to a file that has no imports creates a new import block",
					restorer: func(r *FileRestorer) {
						r.Alias["a.b/a"] = "_"
					},
					expect: `package main

						import _ "a.b/a"

						func main() {}`,
				},
				{
					name: "add-one",
					desc: "adding an import to a file that has no imports creates a new import block",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/a", Name: "A"}}})
					},
					expect: `package main

						import "a.b/a"

						func main() { a.A() }`,
				},
				{
					name: "add-one-alias",
					desc: "manually added alias work as expected",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/a", Name: "A"}}})
					},
					restorer: func(r *FileRestorer) {
						r.Alias["a.b/a"] = "a1"
					},
					expect: `package main

						import a1 "a.b/a"

						func main() { a1.A() }`,
				},
				{
					name: "add-two",
					desc: "adding two imports to a file that has no imports creates a new import block with parens",
					mutate: func(f *dst.File) {
						b := f.Decls[0].(*dst.FuncDecl).Body
						b.List = append(
							b.List,
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/a", Name: "A"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/b", Name: "B"}}},
						)
					},
					expect: `package main

						import (
							"a.b/a"
							"a.b/b"
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
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/fmt/c", Name: "C"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/fmt/a", Name: "A"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/fmt/b", Name: "B"}}},
							&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "fmt", Name: "Print"}}},
						)
					},
					expect: `package main
			
						import (
							"fmt"

							fmt1 "a.b/fmt/a"
							fmt2 "a.b/fmt/b"
							fmt3 "a.b/fmt/c"
						)
			
						func main() { fmt3.C(); fmt1.A(); fmt2.B(); fmt.Print() }`,
				},
				{
					name: "cgo",
					desc: "if cgo import is found (manually added here), and import is added, it will create a new block below cgo",
					mutate: func(f *dst.File) {

						cgo := &dst.GenDecl{
							Tok: token.IMPORT,
							Specs: []dst.Spec{
								&dst.ImportSpec{Path: &dst.BasicLit{Kind: token.STRING, Value: strconv.Quote("C")}},
							},
						}
						f.Decls = append([]dst.Decl{cgo}, f.Decls...)

						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/a", Name: "A"}}})
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/b", Name: "B"}}})
					},
					expect: `package main
			
						import "C"
			
						import (
							"a.b/a"
							"a.b/b"
						)
			
						func main() { a.A(); b.B() }`,
				},
			},
		},
		{
			name: "single-existing-import-ab",
			root: "a.b",
			src: map[string]string{
				"main/main.go": `package main

					import "a.b/a"

					func main() { a.A() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module a.b",
			},
			cases: []tc{
				{
					name: "change-to-local",
					desc: "change to the local path, should remove selector",
					mutate: func(f *dst.File) {
						sel := f.Decls[1].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
						sel.Path = "a.b/main"
					},
					expect: `package main

						func main() { A() }
					`,
				},
				{
					name: "change-to-dot",
					desc: "change to a dot-import, should remove selector",
					restorer: func(r *FileRestorer) {
						r.Alias["a.b/a"] = "."
					},
					expect: `package main

						import . "a.b/a"

						func main() { A() }
					`,
				},
				{
					name: "change-to-empty",
					desc: "change to empty path, should remove selector",
					mutate: func(f *dst.File) {
						sel := f.Decls[1].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr).Fun.(*dst.Ident)
						sel.Path = ""
					},
					expect: `package main

						func main() { A() }
					`,
				},
				{
					name: "add-c",
					desc: "if C import is found as part of another block, it is ignored and ordered first",
					mutate: func(f *dst.File) {

						firstBlock := f.Decls[0].(*dst.GenDecl)
						firstBlock.Specs = append(firstBlock.Specs, &dst.ImportSpec{Path: &dst.BasicLit{Kind: token.STRING, Value: strconv.Quote("C")}})

						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/b", Name: "B"}}})
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "bufio", Name: "NewReader"}}})
					},
					expect: `package main
			
						import (
							"C"
							"bufio"
			
							"a.b/a"
							"a.b/b"
						)
			
						func main() { a.A(); b.B(); bufio.NewReader() }`,
				},
			},
		},
		{
			name: "single-existing-import",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

					func main() { a.A() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module root",
			},
			cases: []tc{
				{
					name: "add-anon",
					desc: "adding an new anonymous import to a file that already has imports creates a new import",
					restorer: func(r *FileRestorer) {
						r.Alias["root/b"] = "_"
					},
					expect: `package main

						import (
							"root/a"
							_ "root/b"
						)

						func main() { a.A() }`,
				},
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
					name: "change-to-alias",
					desc: "changing a current import to an alias",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "a1"
					},
					expect: `package main
			
						import a1 "root/a"
			
						func main() { a1.A() }`,
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
					name: "add-single-import",
					desc: "adding a simple import should work as expected",
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
					desc: "deleting all the imports should also delete the import block",
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
			name: "existing-anon",
			src: map[string]string{
				"main/main.go": `package main

					import _ "root/a"

					func main() { }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module root",
			},
			cases: []tc{
				{
					name: "change-to-no-alias",
					desc: "changing an anonymous alias to standard should remove the import if it's not used",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = ""
					},
					expect: `package main

						func main() { }`,
				},
				{
					name: "change-to-no-alias-add-usage",
					desc: "changing an anonymous alias to standard should not remove the import if it's used",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = ""
					},
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a", Name: "A"}}})
					},
					expect: `package main

						import "root/a"

						func main() { a.A() }`,
				},
				{
					name: "change-to-alias-add-usage",
					desc: "changing an anonymous alias to a custom alias should not remove the import if it's used",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "b"
					},
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a", Name: "A"}}})
					},
					expect: `package main

						import b "root/a"

						func main() { b.A() }`,
				},
				{
					name: "add-standard",
					desc: "adding a standard import to a file with an anon import, the anon import stays anon",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b", Name: "B"}}})
					},
					expect: `package main

						import (
							_ "root/a"
							"root/b"
						)

						func main() { b.B() }`,
				},
				{
					name: "convert-to-standard",
					desc: "convert anonymous import to standard works as intended",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/a", Name: "A"}}})
					},
					expect: `package main

						import "root/a"

						func main() { a.A() }`,
				},
			},
		},
		{
			name: "block-not-rearranged",
			root: "a.b",
			src: map[string]string{
				"main/main.go": `package main

					import (
						"a.b/a"
						"a.b/b"
						"a.b/c"
						"fmt"
					)

					func main() { a.A(); b.B(); c.C(); fmt.Print() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"c/c.go": "package c \n\n func C(){}",
				"d/d.go": "package d \n\n func D(){}",
				"go.mod": "module a.b",
			},
			cases: []tc{
				{
					name: "no-addition",
					desc: "doesn't re-arrange first import block if no additions",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = b.List[1:4]
					},
					expect: `package main

						import (
							"a.b/b"
							"a.b/c"
							"fmt"
						)

						func main() { b.B(); c.C(); fmt.Print() }`,
				},
				{
					name: "one-added-re-arrange",
					desc: "re-arrange first import block if additions",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/d", Name: "D"}}})
					},
					expect: `package main

						import (
							"fmt"

							"a.b/a"
							"a.b/b"
							"a.b/c"
							"a.b/d"
						)

						func main() { a.A(); b.B(); c.C(); fmt.Print(); d.D() }`,
				},
			},
		},
		{
			name: "two-blocks-not-arranged",
			root: "a.b",
			src: map[string]string{
				"main/main.go": `package main

					import (
						"a.b/a"
						"fmt"
					)

					import (
						"a.b/b"
						"io"
					)

					func main() { a.A(); b.B(); io.Copy(nil, nil); fmt.Print() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"c/c.go": "package c \n\n func C(){}",
				"go.mod": "module a.b",
			},
			cases: []tc{
				{
					name: "add-one",
					desc: "only re-arrange first block",
					mutate: func(f *dst.File) {
						b := f.Decls[2].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/c", Name: "C"}}})
					},
					expect: `package main

						import (
							"fmt"
							
							"a.b/a"
							"a.b/c"
						)

						import (
							"a.b/b"
							"io"
						)

						func main() { a.A(); b.B(); io.Copy(nil, nil); fmt.Print(); c.C() }`,
				},
			},
		},
		{
			name: "two-existing",
			src: map[string]string{
				"main/main.go": `package main

					import (
						"root/a"
						"root/b"
					)

					func main() { a.A(); b.B() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module root",
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
			src: map[string]string{
				"main/main.go": `package main

					// first-import-block
					import (
						"root/a"
						"root/b"
					)

					// second-import-block
					import (
						"root/c"
						"root/d"
					)

					func main() { a.A(); b.B(); c.C(); d.D(); }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"c/c.go": "package c \n\n func C(){}",
				"d/d.go": "package d \n\n func D(){}",
				"e/e.go": "package e \n\n func E(){}",
				"go.mod": "module root",
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

						// first-import-block
						import (
							"root/a"
							"root/b"
						)

						func main() { a.A(); b.B() }`,
				},
				{
					name: "delete-from-second-block",
					desc: "imports can be deleted from non-primary block",
					mutate: func(f *dst.File) {
						b := f.Decls[2].(*dst.FuncDecl).Body
						b.List = b.List[0:3]
					},
					expect: `package main

						// first-import-block
						import (
							"root/a"
							"root/b"
						)

						// second-import-block
						import "root/c"

						func main() { a.A(); b.B(); c.C() }`,
				},
				{
					name: "block-deleted-and-added",
					desc: "if all imports are removed from first block and one added, it's ok",
					mutate: func(f *dst.File) {
						b := f.Decls[2].(*dst.FuncDecl).Body
						b.List = b.List[2:4]
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/e", Name: "E"}}})
					},
					expect: `package main
						
						// first-import-block
						import "root/e"

						// second-import-block
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
			root: "a.b",
			src: map[string]string{
				"main/main.go": `package main

					import (
						// before c
						"a.b/c" // after c
						// before a
						"a.b/a" // after a
						// before fmt
						"fmt" // after fmt
					)

					func main() { a.A(); c.C(); fmt.Print() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"c/c.go": "package c \n\n func C(){}",
				"go.mod": "module a.b",
			},
			cases: []tc{
				{
					name: "decorations-retained",
					desc: "decorations in re-arranged block are retained",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "a.b/b", Name: "B"}}})
					},
					expect: `package main

						import (
							// before fmt
							"fmt" // after fmt

							// before a
							"a.b/a" // after a
							"a.b/b"
							// before c
							"a.b/c" // after c
						)

						func main() { a.A(); c.C(); fmt.Print(); b.B() }`,
				},
			},
		},
		{
			name: "first-block-spacing",
			root: "foo.bar",
			src: map[string]string{
				"main/main.go": `package main

					import (

						"foo.bar/a"

						"fmt"
							
						"bytes"

					)

					func main() { a.A(); fmt.Print(); bytes.Title([]byte{}); }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module foo.bar",
			},
			cases: []tc{
				{
					name: "block-rearranged-spacing-fixed",
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
			
						func main() { a.A(); fmt.Print(); bytes.Title([]byte{}); io.Copy(); b.B() }`,
				},
				{

					name: "block-rearranged-spacing-fixed-delete-first-non-std",
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
			
						func main() { fmt.Print(); bytes.Title([]byte{}); b.B() }`,
				},
			},
		},
		{
			name: "existing-alias",
			src: map[string]string{
				"main/main.go": `package main

					import a1 "root/a"

					func main() { a1.A() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module root",
			},
			cases: []tc{
				{
					name: "alias-retained",
					desc: "alias from import block works correctly",
					mutate: func(f *dst.File) {
						b := f.Decls[1].(*dst.FuncDecl).Body
						b.List = append(b.List, &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Path: "root/b", Name: "B"}}})
					},
					expect: `package main

						import (
							a1 "root/a"
							"root/b"
						)

						func main() { a1.A(); b.B() }`,
				},
				{
					name: "manually-added-overrides",
					desc: "manually added alias take priority over alias in imports block",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "a2"
					},
					expect: `package main

						import a2 "root/a"

						func main() { a2.A() }`,
				},
			},
		},
		{
			name: "existing-alias-package-name",
			src: map[string]string{
				"main/main.go": `package main

					import a "root/a"

					func main() { a.A() }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"go.mod": "module root",
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
			src: map[string]string{
				"main/main.go": `package main

					import (
						aa "root/a"
						bb "root/b"
					)

					import (
						cc "root/c"
						dd "root/d"
					)

					func main() { aa.A(); bb.B(); cc.C(); dd.D(); }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"c/c.go": "package c \n\n func C(){}",
				"d/d.go": "package d \n\n func D(){}",
				"go.mod": "module root",
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
			src: map[string]string{
				"main/main.go": `package main

					import (
						. "root/a"
						"root/b"
						cc "root/c"
					)

					func main() { A(); b.B(); cc.C(); }`,
				"a/a.go": "package a \n\n func A(){}",
				"b/b.go": "package b \n\n func B(){}",
				"c/c.go": "package c \n\n func C(){}",
				"go.mod": "module root",
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
			src: map[string]string{
				"main/main.go": `package main

					import (
						"root/a"
					)

					func main() { a.A(); }`,
				"a/a.go":   "package a \n\n func A(){}",
				"b/a/a.go": "package a \n\n func AA(){}",
				"go.mod":   "module root",
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
		{
			name: "binary-bug",
			src: map[string]string{
				"main/main.go": `package main

					import "encoding/binary"

					func main() {
						_ = binary.LittleEndian.Uint16(nil)
					}`,
				"go.mod": "module root",
			},
			cases: []tc{
				{
					name: "noop",
					expect: `package main

						import "encoding/binary"

						func main() {
							_ = binary.LittleEndian.Uint16(nil)
						}
					`,
				},
			},
		},
		{
			name: "dot-imported-ident-in-x-of-sel",
			src: map[string]string{
				"main/main.go": `package main

					import . "root/a"

					func main() { A.B() }`,
				"a/a.go": "package a \n\n type T struct{} \n\n func (T)B(){} \n\n var A T",
				"go.mod": "module root",
			},
			cases: []tc{
				{
					name: "change-to-normal",
					desc: "ensure Ident in X position of SelectorExpr has been resolved",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = ""
					},
					expect: `package main
			
						import "root/a"
			
						func main() { a.A.B() }`,
				},
			},
		},
		{
			name: "decorated-sel",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

					func main() { /*1*/ a. /*2*/ A. /*3*/ B() /*4*/ }`,
				"a/a.go": "package a \n\n type T struct{} \n\n func (T)B(){} \n\n var A T",
				"go.mod": "module root",
			},
			cases: []tc{
				{
					name: "change-to-dot",
					desc: "ensure X decoration is not lost when converting from SelectorExpr to Ident",
					restorer: func(r *FileRestorer) {
						r.Alias["root/a"] = "."
					},
					expect: `package main
			
						import . "root/a"
			
						func main() { /*1*/ /*2*/ A. /*3*/ B() /*4*/ }`,
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

				rootDir, err := tempDir(test.src)
				if err != nil {
					t.Fatal(err)
				}
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

				if len(pkg.Errors) > 0 {
					for _, v := range pkg.Errors {
						t.Error(v.Error())
					}
					t.Fatal("errors loading package")
				}

				d := NewDecoratorFromPackage(pkg)

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

				if c.mutate != nil {
					c.mutate(file)
				}

				r := NewRestorerWithImports(mainPkg, gopackages.New(mainDir))
				fr := r.FileRestorer()

				if c.restorer != nil {
					c.restorer(fr)
				}

				restoredFile, err := fr.RestoreFile(file)
				if err != nil {
					t.Fatal(err)
				}

				buf := &bytes.Buffer{}
				if err := format.Node(buf, fr.Fset, restoredFile); err != nil {
					t.Fatal(err)
				}

				expected, err := format.Source([]byte(c.expect))
				if err != nil {
					panic(err)
				}

				if buf.String() != string(expected) {
					t.Errorf("expect: %s \n\n found: %s \n\n diff:\n%s", string(expected), buf.String(), diff(string(expected), buf.String()))
				}

			})
		}
	}
}

func TestPackageOrder(t *testing.T) {
	paths := []string{"a.b/d", "a.b/c", "fmt", "bytes", "a/b", "C"}
	sort.Slice(paths, func(i, j int) bool {
		return packagePathOrderLess(paths[i], paths[j])
	})
	expect := "[C a/b bytes fmt a.b/c a.b/d]"
	found := fmt.Sprint(paths)
	if found != expect {
		t.Errorf("expect %s, found %s", expect, found)
	}
}

func TestRestorerDecorationResolver(t *testing.T) {
	tests := []struct {
		skip, solo bool
		name       string
		src        map[string]string
	}{
		{
			name: "sel",
			src: map[string]string{
				"main/main.go": `package main
            
					import "root/a"
            
            		func main() {
            			a.A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space",
			src: map[string]string{
				"main/main.go": `package main
            
					import "root/a"
            
            		func main() {
            			a.
							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-decoration",
			src: map[string]string{
				"main/main.go": `package main
            
					import "root/a"
            
            		func main() {
            			a. /*a*/ A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration",
			src: map[string]string{
				"main/main.go": `package main
            
					import "root/a"
            
            		func main() {
            			a. /*a*/
							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-1",
			src: map[string]string{
				"main/main.go": `package main
            
					import "root/a"
            
            		func main() {
            			a. /*a*/
            
            				A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-2",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

					func main() {
						a.
							/*a*/ A()
					}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-3",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
            				/*a*/
            				A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-4",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
            				/*a*/

            				A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-5",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
            				
							/*a*/ A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-6",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
            				
							/*a*/
							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-decoration-7",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
            				
							/*a*/

							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-line-comment",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a. // a
							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-line-comment-1",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
							// a
							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-line-comment-2",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.

							// a
							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-line-comment-3",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.
							// a

							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
			},
		},
		{
			name: "sel-space-line-comment-4",
			src: map[string]string{
				"main/main.go": `package main

					import "root/a"

            		func main() {
            			a.

							// a

							A()
            		}`,
				"a/a.go": "package a \n\n func A(){}",
				"go.mod": "module root",
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

			// format code and check it hasn't changed
			expected, err := format.Source([]byte(test.src["main/main.go"]))
			if err != nil {
				t.Fatal(err)
			}
			if normalize(string(expected)) != normalize(test.src["main/main.go"]) {
				t.Fatalf("code changed after gofmt. before: \n%s\nafter:\n%s", test.src["main/main.go"], string(expected))
			}

			rootDir, err := tempDir(test.src)
			if err != nil {
				t.Fatal(err)
			}
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

			if len(pkg.Errors) > 0 {
				for _, v := range pkg.Errors {
					t.Error(v.Error())
				}
				t.Fatal("errors loading package")
			}

			d := NewDecoratorFromPackage(pkg)

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

			r := NewRestorerWithImports(mainPkg, gopackages.New(mainDir))
			fr := r.FileRestorer()

			restoredFile, err := fr.RestoreFile(file)
			if err != nil {
				t.Fatal(err)
			}

			buf := &bytes.Buffer{}
			if err := format.Node(buf, fr.Fset, restoredFile); err != nil {
				t.Fatal(err)
			}

			if buf.String() != string(expected) {
				t.Errorf("expect: %s \n\n found: %s \n\n diff:\n%s", string(expected), buf.String(), diff(string(expected), buf.String()))
			}
		})
	}
}
