package decorator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
	"github.com/dave/dst/dstutil"
	"golang.org/x/tools/go/loader"
)

func TestPositions(t *testing.T) {
	path := "github.com/dave/dst/gendst/data"
	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		t.Fatal(err)
	}
	var astFile *ast.File
	for _, v := range prog.Package(path).Files {
		_, name := filepath.Split(prog.Fset.File(v.Pos()).Name())
		if name == "positions.go" {
			astFile = v
			break
		}
	}

	dec := NewDecorator(prog.Fset)
	dec.Path = path
	dec.Resolver = &goast.DecoratorResolver{RestorerResolver: &guess.RestorerResolver{}}

	file, err := dec.DecorateFile(astFile)
	if err != nil {
		t.Fatal(err)
	}

	r1 := regexp.MustCompile(`// ([a-zA-Z]+)\(([0-9])\)`)
	r2 := regexp.MustCompile(`// ([a-zA-Z]+)`)
	var currentNodeType string
	var currentTestIndex int
	var done bool

	dst.Inspect(file, func(n dst.Node) bool {
		if n == nil {
			return false
		}
		_, _, points := dstutil.Decorations(n)
		for _, point := range points {
			for _, text := range point.Decs {
				if text == "// notest" {
					continue
				}
				if r1.MatchString(text) || r2.MatchString(text) {
					if currentNodeType != "" && !done {
						t.Fatalf("missed %s %d", currentNodeType, currentTestIndex)
					}
					if matches := r1.FindStringSubmatch(text); matches != nil {
						currentNodeType = "*dst." + matches[1]
						currentTestIndex, _ = strconv.Atoi(matches[2])
					} else if matches := r2.FindStringSubmatch(text); matches != nil {
						currentNodeType = "*dst." + matches[1]
						currentTestIndex = 0
					}
					done = false
					break
				}
			}
		}
		if fmt.Sprintf("%T", n) == currentNodeType {
			//fmt.Printf("*** Testing %s (%d)\n", currentNodeType, currentTestIndex)
			_, _, points := dstutil.Decorations(n)
			for _, point := range points {
				for _, text := range point.Decs {
					if !strings.HasPrefix(text, "/*") {
						continue
					}
					text := strings.TrimSuffix(strings.TrimPrefix(text, "/*"), "*/")
					if text != point.Name {
						t.Errorf("incorrect position in %s (%d) - expected %s, got %s", currentNodeType, currentTestIndex, text, point.Name)
					}
				}
			}
			done = true
		} else {
			_, _, points := dstutil.Decorations(n)
			for _, point := range points {
				for _, text := range point.Decs {
					if !strings.HasPrefix(text, "/*") {
						continue
					}
					text := strings.TrimSuffix(strings.TrimPrefix(text, "/*"), "*/")
					if text != "Start" && text != "End" {
						// Only tolerate comments moved to adjacent decorations for Start and End
						t.Errorf("comment on wrong decoration: %s (%d) %s -> %T %s\n", currentNodeType, currentTestIndex, text, n, point.Name)
					}
				}
			}
		}
		return true
	})
}
