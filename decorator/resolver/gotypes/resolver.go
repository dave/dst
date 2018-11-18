package gotypes

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func FromPackage(pkg *packages.Package) *IdentResolver {
	return &IdentResolver{
		Path:       pkg.PkgPath,
		Uses:       pkg.TypesInfo.Uses,
		Selections: pkg.TypesInfo.Selections,
	}
}

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
		return "" // pre-defined idents in the universe scope - e.g. "byte"
	}
	if pkg.Path() == r.Path {
		return "" // local package
	}
	return pkg.Path()
}
