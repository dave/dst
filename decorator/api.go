package decorator

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/dave/dst/decorator/resolver/gotypes"
)

// New returns a new DecoratorRestorer with no resolvers set, so the import block won't be managed.
func New() *DecoratorRestorer {
	return &DecoratorRestorer{
		Decorator: &Decorator{
			Map:       newMap(),
			Filenames: map[*dst.File]string{},
		},
		Restorer: &Restorer{
			Map: newMap(),
		},
	}
}

// WithImports returns a new DecoratorRestorer with go/types and go/packages resolvers set, so
// the import block will be managed automatically.
func WithImports() *DecoratorRestorer {
	return &DecoratorRestorer{
		Decorator: &Decorator{
			Map:       newMap(),
			Filenames: map[*dst.File]string{},
			Resolver:  &gotypes.IdentResolver{},
		},
		Restorer: &Restorer{
			Map:      newMap(),
			Resolver: &gopackages.PackageResolver{},
		},
	}
}

type DecoratorRestorer struct {
	*Decorator
	*Restorer
}

type Decorator struct {
	Map
	Filenames map[*dst.File]string // Source file names

	// If a Resolver is provided, it is used to resolve Ident nodes. During decoration, remote
	// identifiers (e.g. usually part of a qualified identifier SelectorExpr, but sometimes on
	// their own for dot-imported packages) are updated with the path of the package they are
	// imported from.
	Resolver resolver.IdentResolver
}

type Restorer struct {
	Map

	// If a Resolver is provided, the names of all imported packages are resolved, and the imports
	// block is updated. All remote identifiers are updated (sometimes this involves changing
	// SelectorExpr.X.Name, or even swapping between Ident and SelectorExpr). To force specific
	// import alias names, use the FileRestorer.Alias map.
	Resolver resolver.PackageResolver
}
