package decorator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/dst/decorator/resolver/gobuild"
)

func TestLoadStdLibAll(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping standard library load test in short mode.")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "list", "./...")
	cmd.Env = []string{
		fmt.Sprintf("GOPATH=%s", build.Default.GOPATH),
		fmt.Sprintf("GOROOT=%s", build.Default.GOROOT),
		fmt.Sprintf("HOME=%s", home),
	}
	cmd.Dir = filepath.Join(build.Default.GOROOT, "src")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v", string(b), err)
	}
	all := strings.Split(strings.TrimSpace(string(b)), "\n")

	testPackageRestoresCorrectlyWithImports(t, all...)

}

func testPackageRestoresCorrectlyWithImports(t *testing.T, path ...string) {
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

					if (p.PkgPath == "net/http" && (fname == "server.go" || fname == "request.go")) || (p.PkgPath == "crypto/x509" && fname == "x509.go") {
						t.Skip("TODO: In net/http/server.go, net/http/request.go, and crypto/x509/x509.go we multiple imports with the same path and different aliases. This edge case would need a complete rewrite of the import management block to support - see see https://github.com/dave/dst/issues/45")
					}

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
						t.Errorf("diff:\n%s", diff(string(expect), buf.String()))
					}
				})
			}
		})
	}
}
