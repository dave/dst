package simple

import "github.com/dave/dst/decorator/resolver"

// PackageResolver is a map of package path -> package name. Names are resolved from this map, and
// if a name doesn't exist in the map, an error is returned. Note that Guess is not a NodeResolver,
// so can't properly resolve identifiers in dot import packages.
type PackageResolver map[string]string

func (r PackageResolver) ResolvePackage(importPath string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	return "", resolver.PackageNotFoundError
}
