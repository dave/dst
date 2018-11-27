package goast

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"sync"

	"github.com/dave/dst/decorator/resolver"
)

// RefResolver is a simple ident resolver that parses the imports block of the file and resolves
// qualified identifiers using resolved package names. It is not possible to resolve identifiers in
// dot-imported packages without the full export data of the imported package, so this resolver will
// return an error if it encounters a dot-import. See gotypes.RefResolver for a dot-imports
// capable ident resolver.
type RefResolver struct {
	PackageResolver resolver.PackageResolver
	filesM          sync.Mutex
	files           map[*ast.File]map[string]string
}

func (r *RefResolver) ResolveIdent(file *ast.File, parent ast.Node, id *ast.Ident) (bool, string, error) {

	if r.PackageResolver == nil {
		return false, "", errors.New("goast.RefResolver should have PackageResolver set")
	}

	imports, err := r.imports(file)
	if err != nil {
		return false, "", err
	}

	se, ok := parent.(*ast.SelectorExpr)
	if !ok {
		return false, "", nil
	}

	xid, ok := se.X.(*ast.Ident)
	if !ok {
		return false, "", nil
	}

	if xid.Obj != nil {
		// Obj != nil -> not a qualified ident
		return false, "", nil
	}

	path, ok := imports[xid.Name]
	if !ok {
		return false, "", nil
	}

	// This ident resolver doesn't ever need to know the local package path because it will not
	// attempt to resolve idents that are not inside SelectorExpr nodes.

	return false, path, nil
}

func (r *RefResolver) imports(file *ast.File) (map[string]string, error) {
	r.filesM.Lock()
	defer r.filesM.Unlock()

	if r.files == nil {
		r.files = map[*ast.File]map[string]string{}
	}

	imports, ok := r.files[file]
	if ok {
		return imports, nil
	}

	imports = map[string]string{}
	var done bool
	var outer error
	ast.Inspect(file, func(node ast.Node) bool {
		if done || outer != nil {
			return false
		}
		switch node := node.(type) {
		case *ast.FuncDecl:
			// Import decls must come before all other decls, so as soon as we find a func decl, we
			// can finish.
			done = true
			return false
		case *ast.GenDecl:
			if node.Tok != token.IMPORT {
				// Import decls must come before all other decls, so as soon as we find a non-import
				// gen decl, we can finish.
				done = true
				return false
			}
			return true
		case *ast.ImportSpec:
			path := mustUnquote(node.Path.Value)
			if path == "C" {
				return false
			}
			var name string
			if node.Name != nil {
				name = node.Name.Name
			}
			switch name {
			case ".":
				// We can't resolve "." imports, so throw an error
				outer = fmt.Errorf("goast.RefResolver unsupported dot-import found for %s", path)
				return false
			case "_":
				// Don't need to worry about _ imports
				return false
			case "":
				var err error
				name, err = r.PackageResolver.ResolvePackage(path)
				if err != nil {
					outer = err
					return false
				}
			}
			if p, ok := imports[name]; ok {
				outer = fmt.Errorf("goast.RefResolver found multiple packages using name %s: %s and %s", name, p, path)
				return false
			}
			imports[name] = path
		}
		return true
	})
	if outer != nil {
		return nil, outer
	}

	r.files[file] = imports

	return imports, nil
}

func mustUnquote(s string) string {
	out, err := strconv.Unquote(s)
	if err != nil {
		panic(err)
	}
	return out
}
