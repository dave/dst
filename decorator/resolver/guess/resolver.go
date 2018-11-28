package guess

import "strings"

func New() RestorerResolver {
	return RestorerResolver{}
}

func WithMap(m map[string]string) RestorerResolver {
	return RestorerResolver(m)
}

// RestorerResolver is a map of package path -> package name. Names are resolved from this map, and
// if a name doesn't exist in the map, the package name is guessed from the last part of the path
// (after the last slash).
type RestorerResolver map[string]string

func (r RestorerResolver) ResolvePackage(importPath string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	if !strings.Contains(importPath, "/") {
		return importPath, nil
	}
	return importPath[strings.LastIndex(importPath, "/")+1:], nil
}
