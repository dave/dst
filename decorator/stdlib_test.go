package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"testing"

	"path/filepath"

	"github.com/andreyvit/diff"
	"golang.org/x/tools/go/loader"
)

func TestStdLib(t *testing.T) {

	/*
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
	*/
	all := []string{"github.com/dave/dst/gendst/postests", "archive/tar", "archive/zip", "bufio", "io", "os", "fmt"}

	ignore := map[string]bool{
		"builtin": true,
	}
	skip := map[string]bool{
		"archive/tar:format.go": true,
		"archive/zip:struct.go": true,
		"os:stat_darwin.go":     true,
	}

	for _, pkgPath := range all {

		if ignore[pkgPath] {
			continue
		}

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

			t.Run(name, func(t *testing.T) {
				if skip[name] {
					t.Skip()
				}

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
					t.Errorf("diff: %s", diff.LineDiff(expected.String(), output.String()))
				}
			})
		}
	}
}
