package decorator

import (
	"go/ast"

	"github.com/dave/dst"
)

func newMap() Map {
	return Map{
		Ast: map[dst.Node]ast.Node{},
		Dst: map[ast.Node]dst.Node{},
	}
}

// Map holds a record of the mapping between ast and dst nodes.
type Map struct {
	Ast map[dst.Node]ast.Node
	Dst map[ast.Node]dst.Node
}
