package resolver

import (
	"errors"
	"go/ast"
)

// PackageResolver resolves a package path to a package name.
type PackageResolver interface {
	ResolvePackage(path string) (string, error)
}

// RefResolver resolves an identifier to a local or remote reference.
//
// Returns local == false, path == "" if the node is not a local or remote reference (e.g. a field
// in a composite literal, the selector in a selector expression etc.).
//
// Returns local == true, path == "" is the node is a local reference.
//
// Returns local == false, path != "" is the node is a remote reference.
type RefResolver interface {
	ResolveIdent(file *ast.File, parent ast.Node, id *ast.Ident) (local bool, path string, err error)
}

var PackageNotFoundError = errors.New("package not found")
