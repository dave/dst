package decorator

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
)

func Parse(src interface{}) (*dst.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "a.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return Decorate(f, fset), nil
}

func Decorate(f *ast.File, fset *token.FileSet) *dst.File {
	d := &decorator{
		nodes:       map[ast.Node]dst.Node{},
		decorations: map[ast.Node]map[string][]string{},
	}
	return d.decorateFile(f, fset)
}

type decorator struct {
	nodes       map[ast.Node]dst.Node
	decorations map[ast.Node]map[string][]string
}

func (d *decorator) decorateFile(f *ast.File, fset *token.FileSet) *dst.File {
	fragger := &Fragger{}
	fragger.Fragment(f, fset)
	//fragger.debug(os.Stdout, fset)
	d.decorations = fragger.Link()
	return d.decorateNode(f).(*dst.File)
}

type decorationInfo struct {
	name string
	decs []string
}
