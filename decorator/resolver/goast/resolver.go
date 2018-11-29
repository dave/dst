package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"sync"

	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/guess"
)

func New() *DecoratorResolver {
	return &DecoratorResolver{}
}

func WithResolver(resolver resolver.RestorerResolver) *DecoratorResolver {
	return &DecoratorResolver{RestorerResolver: resolver}
}

// DecoratorResolver is a simple ident resolver that parses the imports block of the file and resolves
// qualified identifiers using resolved package names. It is not possible to resolve identifiers in
// dot-imported packages without the full export data of the imported package, so this resolver will
// return an error if it encounters a dot-import. See gotypes.DecoratorResolver for a dot-imports
// capable ident resolver.
type DecoratorResolver struct {
	RestorerResolver resolver.RestorerResolver
	filesM           sync.Mutex
	files            map[*ast.File]map[string]string
}

func (r *DecoratorResolver) ResolveIdent(file *ast.File, parent ast.Node, parentField string, id *ast.Ident) (string, error) {

	if r.RestorerResolver == nil {
		r.RestorerResolver = guess.New()
	}

	imports, err := r.imports(file)
	if err != nil {
		return "", err
	}

	se, ok := parent.(*ast.SelectorExpr)
	if !ok || parentField != "Sel" {
		return "", nil
	}

	xid, ok := se.X.(*ast.Ident)
	if !ok {
		return "", nil
	}

	if xid.Obj != nil {
		// Obj != nil -> not a qualified ident
		return "", nil
	}

	path, ok := imports[xid.Name]
	if !ok {
		return "", nil
	}

	return path, nil
}

func (r *DecoratorResolver) imports(file *ast.File) (map[string]string, error) {
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
				outer = fmt.Errorf("goast.DecoratorResolver unsupported dot-import found for %s", path)
				return false
			case "_":
				// Don't need to worry about _ imports
				return false
			case "":
				var err error
				name, err = r.RestorerResolver.ResolvePackage(path)
				if err != nil {
					outer = err
					return false
				}
			}
			if p, ok := imports[name]; ok {
				outer = fmt.Errorf("goast.DecoratorResolver found multiple packages using name %s: %s and %s", name, p, path)
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
