package gopackages

import (
	"fmt"

	"github.com/dave/dst/decorator/resolver"
	"golang.org/x/tools/go/packages"
)

type PackageResolver struct {
	Config packages.Config
	Dir    string
}

func (r *PackageResolver) ResolvePackage(path string) (string, error) {

	if r.Dir != "" {
		r.Config.Dir = r.Dir
	}
	r.Config.Mode = packages.LoadTypes
	r.Config.Tests = false

	pkgs, err := packages.Load(&r.Config, "pattern="+path)
	if err != nil {
		return "", err
	}

	if len(pkgs) > 1 {
		return "", fmt.Errorf("%d packages found for %s, %s", len(pkgs), path, r.Config.Dir)
	}
	if len(pkgs) == 0 {
		return "", resolver.PackageNotFoundError
	}

	p := pkgs[0]

	if len(p.Errors) > 0 {
		return "", p.Errors[0]
	}

	return p.Name, nil
}
