package simple

import "github.com/dave/dst/decorator/resolver"

func New(m map[string]string) RestorerResolver {
	return RestorerResolver(m)
}

// RestorerResolver is a map of package path -> package name. Names are resolved from this map, and
// if a name doesn't exist in the map, an error is returned. Note that Guess is not a NodeResolver,
// so can't properly resolve identifiers in dot import packages.
type RestorerResolver map[string]string

func (r RestorerResolver) ResolvePackage(importPath string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	return "", resolver.ErrPackageNotFound
}
