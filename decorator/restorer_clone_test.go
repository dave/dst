package decorator

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/gobuild"
)

func TestClone(t *testing.T) {
	testPackageRestoresCorrectlyWithClone(
		t,
		"github.com/dave/dst/gendst/data",
		"fmt",
		"bytes",
		"io",
	)
}

func testPackageRestoresCorrectlyWithClone(t *testing.T, path ...string) {
	t.Helper()
	pkgs, err := Load(nil, path...)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pkgs {

		t.Run(p.PkgPath, func(t *testing.T) {

			// must use go/build package resolver for standard library because of https://github.com/golang/go/issues/26924
			r := NewRestorer()
			r.Path = p.PkgPath
			r.Resolver = &gobuild.RestorerResolver{Dir: p.Dir}

			for _, file := range p.Syntax {

				fpath := p.Decorator.Filenames[file]
				_, fname := filepath.Split(fpath)

				t.Run(fname, func(t *testing.T) {

					cloned := dst.Clone(file).(*dst.File)

					buf := &bytes.Buffer{}
					if err := r.Fprint(buf, cloned); err != nil {
						t.Fatal(err)
					}

					existing, err := ioutil.ReadFile(fpath)
					if err != nil {
						t.Fatal(err)
					}
					expect, err := format.Source(existing)
					if err != nil {
						t.Fatal(err)
					}
					if string(expect) != buf.String() {
						t.Errorf("diff:\n%s", diff(string(expect), buf.String()))
					}
				})
			}
		})
	}
}
