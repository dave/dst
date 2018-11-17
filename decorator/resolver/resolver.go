package resolver

import (
	"context"
	"errors"
	"go/ast"
	"strings"
)

/*

// Consider this file... B and C could be local identifiers from a different file in this package,
// or from the imported package "a". If only one is from "a" and it is removed, we should remove the
// import when we restore the AST. Thus the node resolver interface needs to be able to resolve the
// package using the full info from go/types.

package main

import (
	. "a"
)

func main() {
	B()
	C()
}
*/

// TODO: combine interfaces into Resolver, combine Decorator and Restorer to a single type
//type Resolver interface {
//	PackageResolver
//	IdentResolver
//}

// PackageResolver resolves a package path to a package name.
type PackageResolver interface {
	ResolvePackage(ctx context.Context, importPath, fromDir string) (string, error)
}

// IdentResolver resolves an identifier to a package path. Returns an empty string if the node is
// not an identifier.
type IdentResolver interface {
	ResolveIdent(node *ast.Ident) string
}

var PackageNotFoundError = errors.New("package not found")

// Guess is a map of package path -> package name. Names are resolved from this map, and if a name
// doesn't exist in the map, the package name is guessed from the last part of the path (after the
// last slash).
type Guess map[string]string

func (r Guess) ResolvePackage(ctx context.Context, importPath, fromDir string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	if !strings.Contains(importPath, "/") {
		return importPath, nil
	}
	return importPath[strings.LastIndex(importPath, "/")+1:], nil
}

// Map is a map of package path -> package name. Names are resolved from this map, and if a name
// doesn't exist in the map, an error is returned. Note that Guess is not a NodeResolver, so can't
// properly resolve identifiers in dot import packages.
type Map map[string]string

func (r Map) ResolvePackage(ctx context.Context, importPath, fromDir string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	return "", PackageNotFoundError
}
