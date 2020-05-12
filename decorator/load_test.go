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
