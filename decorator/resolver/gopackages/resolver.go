package gopackages

import (
	"context"
	"fmt"

	"github.com/dave/dst/decorator/resolver"
	"golang.org/x/tools/go/packages"
)

type PackageResolver struct {
	*packages.Config
	Dir string // default dir for when fromDir == ""
}

func (r *PackageResolver) ResolvePackage(ctx context.Context, importPath, fromDir string) (string, error) {

	if r.Config == nil {
		r.Config = &packages.Config{}
	}

	if fromDir != "" {
		r.Config.Dir = fromDir
	} else if r.Dir != "" {
		r.Config.Dir = r.Dir
	}
	r.Tests = false
	r.Mode = packages.LoadTypes

	pkgs, err := packages.Load(r.Config, "pattern="+importPath)
	if err != nil {
		return "", err
	}

	if len(pkgs) > 1 {
		return "", fmt.Errorf("%d packages found for %s, %s", len(pkgs), importPath, fromDir)
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
