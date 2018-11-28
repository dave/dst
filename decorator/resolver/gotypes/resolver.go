package gotypes

import (
	"errors"
	"go/ast"
	"go/types"
)

func New(uses map[*ast.Ident]types.Object) *DecoratorResolver {
	return &DecoratorResolver{Uses: uses}
}

type DecoratorResolver struct {
	Uses map[*ast.Ident]types.Object // Types info - must include Uses
}

func (r *DecoratorResolver) ResolveIdent(file *ast.File, parent ast.Node, id *ast.Ident) (string, error) {

	if r.Uses == nil {
		return "", errors.New("gotypes.DecoratorResolver needs Uses in types info")
	}

	se, ok := parent.(*ast.SelectorExpr)
	if ok {
		// if the parent is a SelectorExpr, only return the path if X is an ident and a package
		xid, ok := se.X.(*ast.Ident)
		if !ok {
			return "", nil // x is not an ident -> not a qualified identifier
		}
		obj, ok := r.Uses[xid]
		if !ok {
			return "", nil // not found in uses -> not a qualified identifier
		}
		pn, ok := obj.(*types.PkgName)
		if !ok {
			return "", nil // not a pkgname -> not a remote identifier
		}
		return pn.Imported().Path(), nil
	}

	obj, ok := r.Uses[id]
	if !ok {
		return "", nil // not found in uses -> not a remote identifier
	}

	if v, ok := obj.(*types.Var); ok && v.IsField() {
		return "", nil // field ident -> doesn't need qualified ident
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return "", nil // pre-defined idents in the universe scope - e.g. "byte"
	}

	return pkg.Path(), nil
}
