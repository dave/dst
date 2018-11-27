package gotypes

import (
	"errors"
	"go/ast"
	"go/types"
	"strings"
)

type RefResolver struct {
	Path string      // Local package path
	Info *types.Info // Types info - must include Uses
}

func (r *RefResolver) ResolveIdent(file *ast.File, parent ast.Node, id *ast.Ident) (bool, string, error) {

	if r.Info == nil || r.Info.Uses == nil {
		return false, "", errors.New("gotypes.RefResolver needs Uses in types info")
	}

	if r.Path == "" {
		return false, "", errors.New("gotypes.RefResolver needs Path")
	}

	se, ok := parent.(*ast.SelectorExpr)
	if ok {
		// if the parent is a SelectorExpr, only return the path if X is an ident and a package
		xid, ok := se.X.(*ast.Ident)
		if !ok {
			return false, "", nil // x is not an ident -> not a qualified identifier
		}
		obj, ok := r.Info.Uses[xid]
		if !ok {
			return false, "", nil // not found in uses -> not a qualified identifier
		}
		pn, ok := obj.(*types.PkgName)
		if !ok {
			return false, "", nil // not a pkgname -> not a remote identifier
		}
		return false, stripVendor(pn.Imported().Path()), nil
	}

	obj, ok := r.Info.Uses[id]
	if !ok {
		return false, "", nil // not found in uses -> not a remote identifier
	}

	if v, ok := obj.(*types.Var); ok && v.IsField() {
		return false, "", nil // field ident -> doesn't need qualified ident
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false, "", nil // pre-defined idents in the universe scope - e.g. "byte"
	}

	path := stripVendor(pkg.Path())

	if path == stripVendor(r.Path) {
		return true, "", nil // ident in the local package
	}

	return false, path, nil
}

func stripVendor(path string) string {
	findVendor := func(path string) (index int, ok bool) {
		// Two cases, depending on internal at start of string or not.
		// The order matters: we must return the index of the final element,
		// because the final one is where the effective import path starts.
		switch {
		case strings.Contains(path, "/vendor/"):
			return strings.LastIndex(path, "/vendor/") + 1, true
		case strings.HasPrefix(path, "vendor/"):
			return 0, true
		}
		return 0, false
	}
	i, ok := findVendor(path)
	if !ok {
		return path
	}
	return path[i+len("vendor/"):]
}
