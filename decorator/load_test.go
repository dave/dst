package decorator

import (
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/simple"
	"golang.org/x/tools/go/packages"
)

func TestLoad(t *testing.T) {
	code := map[string]string{
		"a.go": `package a

			import "fmt"

			func a() {
				fmt.Println("a")
			}
		`,
		"go.mod": "module root\n\ngo 1.14",
	}
	expect := map[string]string{
		"a.go": `package a

			import "fmt"

			func a() {
				fmt.Println("a") // a
			}
		`,
		"go.mod": "module root\n\ngo 1.14",
	}
	dir, err := tempDir(code)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  dir,
	}
	pkgs, err := Load(cfg, "root")
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Syntax {
			dst.Inspect(f, func(n dst.Node) bool {
				switch n.(type) {
				case *dst.CallExpr:
					n.Decorations().End.Append("// a")
				}
				return true
			})
		}
		if err := pkg.Save(); err != nil {
			t.Fatal(err)
		}
	}
	compareDir(t, dir, expect)
}

func TestLoad_IncludeLocalPkg(t *testing.T) {
	code := map[string]string{
		"a/a.go": `package a

			import "fmt"

			type b struct {
				c string
			}

			func a(param b) {
				fmt.Println(b.c)
			}
		`,
		"go.mod": "module root\n\ngo 1.14",
	}
	dir, err := tempDir(code)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  dir,
		Env:  []string{"DST_INCLUDE_LOCAL_PKG=true"},
	}
	pkgs, err := Load(cfg, "root/a")
	if err != nil {
		t.Fatal(err)
	}
	// This line will panic if `code` above changes
	path := pkgs[0].Syntax[0].Decls[2].(*dst.FuncDecl).Type.Params.List[0].Type.(*dst.Ident).Path
	expect := "root/a"
	if path != expect {
		t.Errorf("expected %q, found %q", expect, path)
	}
}

func TestPackage_SaveWithResolver(t *testing.T) {
	code := map[string]string{
		"a.go": `package a

			import "fmt"

			func a() {
				fmt.Println("a")
			}
		`,
		"go.mod": "module root\n\ngo 1.14",
	}
	expect := map[string]string{
		"a.go": `package a

			import "fmt"

			func a() {
				alternate_pkg_name.Println("a")
			}
		`,
		"go.mod": "module root\n\ngo 1.14",
	}
	dir, err := tempDir(code)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  dir,
	}
	pkgs, err := Load(cfg, "root")
	if err != nil {
		t.Fatal(err)
	}
	res := simple.New(map[string]string{"fmt": "alternate_pkg_name"})
	for _, pkg := range pkgs {
		if err := pkg.SaveWithResolver(res); err != nil {
			t.Fatal(err)
		}
	}
	compareDir(t, dir, expect)
}
