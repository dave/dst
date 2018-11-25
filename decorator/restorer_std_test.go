package decorator

import (
	"bytes"
	"testing"

	"path/filepath"

	"fmt"
	"go/build"
	"go/format"
	"go/parser"
	"os/exec"
	"strings"

	"github.com/andreyvit/diff"
	"golang.org/x/tools/go/loader"
)

func TestStdLibAll(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping standard library test in short mode.")
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
			testPackageRestoresCorrectly(t, pkgPath)
		})

	}
}

func testPackageRestoresCorrectly(t *testing.T, path string) {
	t.Helper()
	conf := loader.Config{
		ParserMode: parser.ParseComments,
	}
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}
	pi := prog.Package(path)
	for _, astFile := range pi.Files {

		_, fname := filepath.Split(prog.Fset.File(astFile.Pos()).Name())

		t.Run(fname, func(t *testing.T) {
			expected := &bytes.Buffer{}
			if err := format.Node(expected, prog.Fset, astFile); err != nil {
				t.Fatal(err)
			}

			dstFile, err := DecorateFile(prog.Fset, astFile)
			if err != nil {
				t.Fatal(err)
			}

			restoredFset, restoredFile, err := Restore(dstFile)
			if err != nil {
				t.Fatal(err)
			}

			output := &bytes.Buffer{}
			if err := format.Node(output, restoredFset, restoredFile); err != nil {
				t.Fatal(err)
			}

			if expected.String() != output.String() {
				t.Error(diff.LineDiff(expected.String(), output.String()))
			}
		})

	}
}
