package decorator

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dave/dst"
	"golang.org/x/tools/go/packages"
)

func Load(cfg *packages.Config, patterns ...string) ([]*Package, error) {

	if cfg == nil {
		cfg = &packages.Config{Mode: packages.LoadSyntax}
	}

	if cfg.Mode != packages.LoadSyntax && cfg.Mode != packages.LoadAllSyntax {
		return nil, errors.New("config mode should be LoadSyntax or LoadAllSyntax")
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	var out []*Package
	for _, pkg := range pkgs {

		p := &Package{
			Package: pkg,
		}

		if len(pkg.Syntax) > 0 {

			// Only decorate files in the GoFiles list. Syntax also has preprocessed cgo files which
			// break things.
			goFiles := make(map[string]bool, len(pkg.GoFiles))
			for _, fpath := range pkg.GoFiles {
				goFiles[fpath] = true
			}

			p.Decorator = NewWithImports(pkg)
			for _, f := range pkg.Syntax {
				fpath := pkg.Fset.File(f.Pos()).Name()
				if !goFiles[fpath] {
					continue
				}
				file, err := p.Decorator.DecorateFile(f)
				if err != nil {
					return nil, err
				}
				p.Files = append(p.Files, file)
			}

			dir, _ := filepath.Split(pkg.Fset.File(pkg.Syntax[0].Pos()).Name())
			p.Dir = dir
		}

		out = append(out, p)

	}
	return out, nil
}

type Package struct {
	*packages.Package
	Dir       string
	Decorator *Decorator
	Files     []*dst.File
}

func (p *Package) Save() error {
	return p.save(ioutil.WriteFile)
}

func (p *Package) save(writeFile func(filename string, data []byte, perm os.FileMode) error) error {
	r := NewRestorerWithImports(p.PkgPath, p.Dir)
	for _, file := range p.Files {
		buf := &bytes.Buffer{}
		if err := r.Fprint(buf, file); err != nil {
			return err
		}
		if err := writeFile(p.Decorator.Filenames[file], buf.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}
