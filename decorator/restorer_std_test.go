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
	if len(pkgs) == 0 {
		t.Fatal("No packages loaded")
	}
	// we skip some packages because they have no source:
	skip := map[string]bool{
		"unsafe":                   true,
		"embed/internal/embedtest": true,
		"os/signal/internal/pty":   true,
	}
	// These files import the same path more than once with different aliases (e.g. "unsafe" and
	// `_ "unsafe"`). The import management block can only hold one alias per path, so the imports
	// are merged and the restored file doesn't match. Supporting this would need a rewrite of the
	// import management block - see https://github.com/dave/dst/issues/45
	skipDuplicateImports := map[string]map[string]bool{
		"crypto/rand":        {"rand.go": true},
		"crypto/x509":        {"x509.go": true},
		"internal/godebug":   {"godebug.go": true},
		"net/http":           {"server.go": true, "request.go": true},
		"reflect":            {"badlinkname.go": true},
		"runtime":            {"rand.go": true},
		"testing/cryptotest": {"rand.go": true},
	}
	for _, p := range pkgs {
		if skip[p.PkgPath] {
			continue
		}
		if len(p.GoFiles) == 0 {
			// nothing to restore in a package with no non-test source - e.g. one that only
			// contains external test files
			continue
		}
		if len(p.Syntax) == 0 && len(p.CompiledGoFiles) != len(p.GoFiles) {
			// Load only decorates the files in GoFiles, so a package where every source file is
			// preprocessed by cgo ends up with no decorated files
			continue
		}
		if len(p.Syntax) == 0 {
			t.Fatalf("Package %s has no syntax", p.PkgPath)
		}
		t.Run(p.PkgPath, func(t *testing.T) {

			// must use go/build package resolver for standard library because of https://github.com/golang/go/issues/26924
			r := NewRestorer()
			r.Path = p.PkgPath
			r.Resolver = &gobuild.RestorerResolver{Dir: p.Dir}

			for _, file := range p.Syntax {

				fpath := p.Decorator.Filenames[file]
				_, fname := filepath.Split(fpath)

				t.Run(fname, func(t *testing.T) {

					if skipDuplicateImports[p.PkgPath][fname] {
						t.Skip("TODO: multiple imports with the same path and different aliases - see https://github.com/dave/dst/issues/45")
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
