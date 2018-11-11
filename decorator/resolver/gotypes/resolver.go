package gotypes

import (
	"fmt"
	"go/ast"
	"go/types"
)

type IdentResolver struct {
	Path       string // local path
	Uses       map[*ast.Ident]types.Object
	Selections map[*ast.SelectorExpr]*types.Selection
}

func (r *IdentResolver) ResolveIdent(id *ast.Ident) string {
	obj, ok := r.Uses[id]
	if !ok {
		return "" // not found in uses -> not a remote identifier
	}
	pkg := obj.Pkg()
	if pkg == nil {
		panic(fmt.Sprintf("Ident %q Pkg is nil", id.Name))
	}
	if pkg.Path() == r.Path {
		return "" // local package
	}
	return pkg.Path()
}
