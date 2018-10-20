package decorator

import (
	"bytes"
	"go/format"
	"go/parser"
	"strings"
	"testing"

	"path/filepath"

	"github.com/andreyvit/diff"
	"golang.org/x/tools/go/loader"
)

func TestStdLibExtra(t *testing.T) {

	t.Skip()

	broken := `cmd/compile/internal/ssa:loopbce.go
cmd/compile/internal/syntax:tokens.go
cmd/internal/obj/arm64:a.out.go
cmd/internal/obj/ppc64:a.out.go
encoding/json:scanner.go
text/template:option.go
time:format.go
unicode:graphic.go`
	fields := strings.Fields(broken)
	packages := map[string]map[string]bool{}
	for _, v := range fields {
		parts := strings.Split(v, ":")
		if packages[parts[0]] == nil {
			packages[parts[0]] = map[string]bool{}
		}
		packages[parts[0]][parts[1]] = true
	}

	//all := []string{"github.com/dave/dst/gendst/postests", "archive/tar", "archive/zip", "bufio", "io", "os", "fmt"}

	for pkgPath, files := range packages {

		//fmt.Println(pkgPath)

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

			if !files[filename] {
				continue
			}
			name := pkgPath + ":" + filename
			//fmt.Println(name)

			t.Run(name, func(t *testing.T) {
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

				/*
					expected1, err := format.Source(expected.Bytes())
					if err != nil {
						t.Fatal(err)
					}

					output1, err := format.Source(output.Bytes())
					if err != nil {
						t.Fatal(err)
					}

					if string(expected1) != string(output1) {
						fmt.Println("expected:", string(expected1))
						fmt.Println("found:", string(output1))
						t.Errorf("diff: %s", diff.LineDiff(string(expected1), string(output1)))
					}
				*/

			})
		}
	}
}
