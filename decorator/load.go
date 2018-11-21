package decorator

import (
	"errors"

	"github.com/dave/dst"
	"golang.org/x/tools/go/packages"
)

func Load(cfg *packages.Config, patterns ...string) ([]*Package, error) {

	if cfg.Mode != packages.LoadSyntax && cfg.Mode != packages.LoadAllSyntax {
		return nil, errors.New("config mode should be LoadSyntax or LoadAllSyntax")
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	d := WithImports()

	var out []*Package
	for _, pkg := range pkgs {

		p := &Package{
			Package: pkg,
		}

		if pkg.Syntax != nil {
			p.Decorator = d.PackageDecoratorFromPackage(pkg)
			for _, f := range pkg.Syntax {
				p.Files = append(p.Files, p.Decorator.DecorateFile(f))
			}
		}

		out = append(out, p)

	}
	return out, nil
}

type Package struct {
	*packages.Package
	Decorator *PackageDecorator
	Files     []*dst.File
}
