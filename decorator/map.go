package decorator

import (
	"go/ast"

	"github.com/dave/dst"
)

func newMap() Map {
	return Map{
		Ast: AstMap{
			Nodes:   map[dst.Node]ast.Node{},
			Scopes:  map[*dst.Scope]*ast.Scope{},
			Objects: map[*dst.Object]*ast.Object{},
		},
		Dst: DstMap{
			Nodes:   map[ast.Node]dst.Node{},
			Scopes:  map[*ast.Scope]*dst.Scope{},
			Objects: map[*ast.Object]*dst.Object{},
		},
	}
}

type Map struct {
	Ast AstMap
	Dst DstMap
}

type AstMap struct {
	Nodes   map[dst.Node]ast.Node       // Mapping from dst to ast Nodes
	Objects map[*dst.Object]*ast.Object // Mapping from dst to ast Objects
	Scopes  map[*dst.Scope]*ast.Scope   // Mapping from dst to ast Scopes
}

type DstMap struct {
	Nodes   map[ast.Node]dst.Node       // Mapping from ast to dst Nodes
	Objects map[*ast.Object]*dst.Object // Mapping from ast to dst Objects
	Scopes  map[*ast.Scope]*dst.Scope   // Mapping from ast to dst Scopes
}
