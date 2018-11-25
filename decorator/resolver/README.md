# Managing imports

Managing the imports block is non-trivial. There are two separate interfaces defined by this package
which allow the decorator and restorer to automatically manage the imports block.

The decorator uses a `IdentResolver` which resolves the package path of any `*ast.Ident`. This is 
complicated by dot-import syntax (see below).

The restorer uses a `PackageResolver` which resolves the name of any package given the path. This 
is complicated by vendoring and Go modules.

This package provides several implementations of both interfaces that are suitable for different 
environments:

## IdentResolver implementations

#### gotypes/IdentResolver

This is the default implementation, and provides full compatibility with dot-imports. However this 
requires full export data for all imported packages, so a `go/types.Info` is required. There are 
many ways of loading ast and generating `go/types.Info`. Using `golang.org/x/tools/go/packages.Load` 
is recommended for full modules compatibility. See the `decorator.Load` convenience function to 
automate this.

#### goast/IdentResolver

This is a simplified implementation that only scans a single ast file. This is unable to resolve 
dot-import idents, so will panic if a dot-import is encountered in the import block. It uses the 
provided `PackageResolver` to resolve the names of all imported packages.

## PackageResolver implementations

#### gopackages/PackageResolver

This is the default implementation, and provides full compatibility with modules. It uses 
`golang.org/x/tools/go/packages` to load the package data. This may be very slow if the package is 
inside a module that hasn't been loaded before. 

#### gobuild/PackageResolver

This is an alternative implementation that uses the legacy `go/build` package to load the imported 
package data. This may be needed in some circumstances and provides better performance. This will
ignore modules and just searches the system's `GOPATH`.

#### guess/PackageResolver and simple/PackageResolver

These are very simple implementations for testing purposes. `simple/PackageResolver` resolves paths 
only if they occur in a provided map. `guess/PackageResolver` guesses the package name based on the 
last part of the path.

## Why is resolving identifiers hard?

Consider this file...

```go
package main

import (
	. "a"
)

func main() {
	B()
	C()
}
```

`B` and `C` could be local identifiers from a different file in this package,
or from the imported package `a`. If only one is from `a` and it is removed, we should remove the
import when we restore to `ast`. Thus the resolver needs to be able to resolve the package using 
the full info from `go/types`.

