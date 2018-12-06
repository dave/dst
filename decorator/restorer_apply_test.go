package decorator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/guess"
	"github.com/dave/dst/dstutil"
)

func TestApply(t *testing.T) {
	testPackageRestoresCorrectlyWithApplyClone(
		t,
		"github.com/dave/dst/gendst/data",
		"fmt",
		"bytes",
		"io",
	)
}

func testPackageRestoresCorrectlyWithApplyClone(t *testing.T, path ...string) {
	t.Helper()
	pkgs, err := Load(nil, path...)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pkgs {

		t.Run(p.PkgPath, func(t *testing.T) {

			r := NewRestorer()
			r.Path = p.PkgPath
			r.Resolver = &guess.RestorerResolver{}

			for _, file := range p.Syntax {

				fpath := p.Decorator.Filenames[file]
				_, fname := filepath.Split(fpath)

				t.Run(fname, func(t *testing.T) {

					cloned1 := dst.Clone(file).(*dst.File)
					cloned2 := dst.Clone(file).(*dst.File)

					cloned1 = dstutil.Apply(cloned1, func(c *dstutil.Cursor) bool {
						switch n := c.Node().(type) {
						case *dst.Ident:
							n1 := dst.Clone(c.Node())
							n1.Decorations().End.Replace(fmt.Sprintf("/* %s */", n.Name))
							c.Replace(n1)
						}
						return true
					}, nil).(*dst.File)

					// same with dst.Inspect
					dst.Inspect(cloned2, func(n dst.Node) bool {
						switch n := n.(type) {
						case *dst.Ident:
							n.Decorations().End.Replace(fmt.Sprintf("/* %s */", n.Name))
						}
						return true
					})

					buf1 := &bytes.Buffer{}
					if err := r.Fprint(buf1, cloned1); err != nil {
						t.Fatal(err)
					}

					buf2 := &bytes.Buffer{}
					if err := r.Fprint(buf2, cloned2); err != nil {
						t.Fatal(err)
					}

					if buf1.String() != buf2.String() {
						t.Errorf("diff:\n%s", diff(buf2.String(), buf1.String()))
					}
				})
			}
		})
	}
}
