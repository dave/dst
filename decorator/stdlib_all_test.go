package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"testing"

	"path/filepath"

	"fmt"
	"go/build"
	"os/exec"
	"strings"

	"github.com/andreyvit/diff"
	"golang.org/x/tools/go/loader"
)

func TestStdLibAll(t *testing.T) {

	t.Skip()

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

		fmt.Println(pkgPath)

		conf := loader.Config{
			ParserMode: parser.ParseComments,
		}
		conf.Import(pkgPath)
		prog, err := conf.Load()
		if err != nil {
			panic(err)
		}
		pi := prog.Package(pkgPath)
		for _, astFile := range pi.Files {

			_, filename := filepath.Split(prog.Fset.File(astFile.Pos()).Name())
			name := pkgPath + ":" + filename

			expected := &bytes.Buffer{}
			if err := format.Node(expected, prog.Fset, astFile); err != nil {
				t.Fatal(err)
			}

			dstFile := DecorateFile(prog.Fset, astFile)

			restoredFset, restoredFile := Restore(dstFile)

			output := &bytes.Buffer{}
			if err := format.Node(output, restoredFset, restoredFile); err != nil {
				t.Fatal(err)
			}

			if expected.String() != output.String() {
				t.Errorf("%s: diff: %s", name, diff.LineDiff(expected.String(), output.String()))
			}
		}
	}
}
