package gotypes

import (
	"errors"
	"go/ast"
	"go/types"
)

type IdentResolver struct {
	Info *types.Info
}

func (r *IdentResolver) ResolveIdent(file *ast.File, id *ast.Ident) (string, error) {

	if r.Info == nil || r.Info.Uses == nil || r.Info.Selections == nil {
		return "", errors.New("gotypes.IdentResolver needs Uses and Selections in types info")
	}

	obj, ok := r.Info.Uses[id]
	if !ok {
		return "", nil // not found in uses -> not a remote identifier
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return "", nil // pre-defined idents in the universe scope - e.g. "byte"
	}
	return pkg.Path(), nil
}
