package guess

import "strings"

func New() PackageResolver {
	return PackageResolver{}
}

func WithMap(m map[string]string) PackageResolver {
	return PackageResolver(m)
}

// PackageResolver is a map of package path -> package name. Names are resolved from this map, and
// if a name doesn't exist in the map, the package name is guessed from the last part of the path
// (after the last slash).
type PackageResolver map[string]string

func (r PackageResolver) ResolvePackage(importPath string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	if !strings.Contains(importPath, "/") {
		return importPath, nil
	}
	return importPath[strings.LastIndex(importPath, "/")+1:], nil
}
