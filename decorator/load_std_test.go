package decorator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/dst/decorator/resolver/gobuild"

	"github.com/andreyvit/diff"
)

func TestLoadStdLibAll(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping standard library load test in short mode.")
	}

	cmd := exec.Command("go", "list", "./...")
	cmd.Env = []string{
		fmt.Sprintf("GOPATH=%s", build.Default.GOPATH),
		fmt.Sprintf("GOROOT=%s", build.Default.GOROOT),
	}
	cmd.Dir = filepath.Join(build.Default.GOROOT, "src")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	all := strings.Split(strings.TrimSpace(string(b)), "\n")

	ignore := map[string]bool{
		"builtin": true,
	}

	for _, pkgPath := range all {

		if ignore[pkgPath] {
			continue
		}

		t.Run(pkgPath, func(t *testing.T) {
			testPackageRestoresCorrectlyWithImports(t, pkgPath)
		})
	}
}

func testPackageRestoresCorrectlyWithImports(t *testing.T, path string) {
	t.Helper()
	pkgs, err := Load(nil, path)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pkgs {

		// must use go/build package resolver for standard library because of https://github.com/golang/go/issues/26924
		r := NewRestorer()
		r.Path = p.PkgPath
		r.Resolver = &gobuild.PackageResolver{Dir: p.Dir}

		for _, file := range p.Files {

			fpath := p.Decorator.Filenames[file]
			_, fname := filepath.Split(fpath)

			t.Run(fname, func(t *testing.T) {
				buf := &bytes.Buffer{}
				if err := r.Fprint(buf, file); err != nil {
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
					t.Error(diff.LineDiff(string(expect), buf.String()))
				}
			})
		}
	}
}
