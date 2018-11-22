package gobuild

import (
	"go/build"

	"github.com/dave/dst/decorator/resolver"
)

type PackageResolver struct {
	// FindPackage is called during Load to create the build.Package for a given import path from a
	// given directory. If FindPackage is nil, (*build.Context).Import is used. A client may use
	// this hook to adapt to a proprietary build system that does not follow the "go build" layout
	// conventions, for example. It must be safe to call concurrently from multiple goroutines.
	//
	// It should be noted that Manager only uses the Name from the returned *build.Package, so all
	// other fields can be left empty (as in SimpleFinder).
	FindPackage func(ctxt *build.Context, importPath, fromDir string, mode build.ImportMode) (*build.Package, error)
	Context     *build.Context
	Dir         string
}

func (r *PackageResolver) ResolvePackage(importPath string) (string, error) {

	fp := r.FindPackage
	if fp == nil {
		fp = (*build.Context).Import
	}

	bc := r.Context
	if bc == nil {
		bc = &build.Default
	}

	p, err := fp(bc, importPath, r.Dir, 0)
	if err != nil {
		return "", err
	}

	if p == nil {
		return "", resolver.PackageNotFoundError
	}

	return p.Name, nil
}
