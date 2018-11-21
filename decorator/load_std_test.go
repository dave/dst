package decorator

import (
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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

		testPackageRestoresCorrectlyWithImports(t, pkgPath)

	}
}

func testPackageRestoresCorrectlyWithImports(t *testing.T, path string) {
	t.Helper()
	t.Run(path, func(t *testing.T) {
		pkgs, err := Load(nil, path)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range pkgs {
			err := v.save(func(filename string, data []byte, perm os.FileMode) error {
				existing, err := ioutil.ReadFile(filename)
				if err != nil {
					t.Fatal(err)
				}
				expect, err := format.Source(existing)
				if err != nil {
					t.Fatal(err)
				}
				if string(expect) != string(data) {
					_, fname := filepath.Split(filename)
					t.Fatalf("%s: %s", fname, diff.LineDiff(string(expect), string(data)))
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}
