package decorator

import (
	"go/ast"
	"go/token"

	"github.com/dave/dst"
)

func New() *Decorator {
	return &Decorator{
		nodes:       map[ast.Node]dst.Node{},
		decorations: map[ast.Node]map[string][]string{},
	}
}

type Decorator struct {
	nodes       map[ast.Node]dst.Node
	decorations map[ast.Node]map[string][]string
}

func (d *Decorator) Decorate(f *ast.File, fset *token.FileSet) *dst.File {
	fragger := &Fragger{}
	fragger.Fragment(f, fset)

	//p.debug(os.Stdout, fset)

	d.decorations = fragger.Link()
	return d.DecorateNode(f).(*dst.File)
}

type decorationInfo struct {
	name string
	decs []string
}
