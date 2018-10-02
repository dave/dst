package decorator

import (
	"go/ast"
	"go/token"

	"github.com/dave/dst"
)

func New() *Decorator {
	return &Decorator{
		nodes:       map[ast.Node]dst.Node{},
		decorations: map[ast.Node][]dst.Decoration{},
	}
}

type Decorator struct {
	nodes       map[ast.Node]dst.Node
	decorations map[ast.Node][]dst.Decoration
}

func (d *Decorator) Decorate(f *ast.File, fset *token.FileSet) *dst.File {
	p := &Fragger{}
	p.Fragment(f, fset)

	//p.debug(os.Stdout, fset)

	d.decorations = p.Link()
	return d.DecorateNode(f).(*dst.File)
}
