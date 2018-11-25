package resolver

import (
	"errors"
	"go/ast"
)

// PackageResolver resolves a package path to a package name.
type PackageResolver interface {
	ResolvePackage(path string) (string, error)
}

// IdentResolver resolves an identifier to a package path. Returns an empty string if the node is
// not an identifier.
type IdentResolver interface {
	ResolveIdent(file *ast.File, parent ast.Node, id *ast.Ident) (string, error)
}

var PackageNotFoundError = errors.New("package not found")
