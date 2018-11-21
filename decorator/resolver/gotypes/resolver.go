package gotypes

import (
	"errors"
	"go/ast"
	"go/types"
)

type IdentResolver struct{}

func (r *IdentResolver) ResolveIdent(id *ast.Ident, info *types.Info, file *ast.File, dir string) (string, error) {

	if info == nil || info.Uses == nil || info.Selections == nil {
		return "", errors.New("gotypes.IdentResolver needs Uses and Selections in types info")
	}

	obj, ok := info.Uses[id]
	if !ok {
		return "", nil // not found in uses -> not a remote identifier
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return "", nil // pre-defined idents in the universe scope - e.g. "byte"
	}
	return pkg.Path(), nil
}
