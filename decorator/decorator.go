package decorator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver"
)

// Parse uses parser.ParseFile to parse and decorate a Go source file. The src parameter should
// be string, []byte, or io.Reader.
func Parse(src interface{}) (*dst.File, error) {
	fset := token.NewFileSet()

	// If ParseFile returns an error and also a non-nil file, the errors were just parse errors so
	// we should continue decorating the file and return the error.
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil && f == nil {
		return nil, err
	}

	return Decorate(fset, f).(*dst.File), err
}

// ParseFile uses parser.ParseFile to parse and decorate a Go source file.
func ParseFile(fset *token.FileSet, filename string, src interface{}, mode parser.Mode) (*dst.File, error) {
	f, err := parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	return Decorate(fset, f).(*dst.File), nil
}

// ParseExpr uses parser.ParseExpr to parse and decorate a Go expression. It should be noted that
// this is of limited use because comments are not parsed by parser.ParseExpr.
func ParseExpr(x string) (dst.Expr, error) {
	expr, err := parser.ParseExpr(x)
	if err != nil {
		return nil, err
	}
	return Decorate(nil, expr).(dst.Expr), nil
}

// ParseDir uses parser.ParseDir to parse and decorate a directory containing Go source.
func ParseDir(fset *token.FileSet, dir string, filter func(os.FileInfo) bool, mode parser.Mode) (map[string]*dst.Package, error) {
	pkgs, err := parser.ParseDir(fset, dir, filter, mode)
	if err != nil {
		return nil, err
	}
	d := New(fset)
	out := map[string]*dst.Package{}
	for k, v := range pkgs {
		out[k] = d.Decorate(v).(*dst.Package)
	}
	return out, nil
}

// ParseExprFrom uses parser.ParseExprFrom to parse and decorate a Go expression. It should be noted
// that this is of limited use because comments are not parsed by parser.ParseExprFrom.
func ParseExprFrom(fset *token.FileSet, filename string, src interface{}, mode parser.Mode) (dst.Expr, error) {
	expr, err := parser.ParseExprFrom(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	return Decorate(fset, expr).(dst.Expr), nil
}

// Decorate decorates an ast.Node and returns a dst.Node.
func Decorate(fset *token.FileSet, n ast.Node) dst.Node {
	return New(fset).Decorate(n)
}

// Decorate decorates a *ast.File and returns a *dst.File.
func DecorateFile(fset *token.FileSet, f *ast.File) *dst.File {
	return New(fset).Decorate(f).(*dst.File)
}

// New returns a new decorator.
func New(fset *token.FileSet) *Decorator {
	return &Decorator{
		Fset: fset,
		Map: Map{
			Ast: AstMap{
				Nodes:   map[dst.Node]ast.Node{},
				Scopes:  map[*dst.Scope]*ast.Scope{},
				Objects: map[*dst.Object]*ast.Object{},
			},
			Dst: DstMap{
				Nodes:   map[ast.Node]dst.Node{},
				Scopes:  map[*ast.Scope]*dst.Scope{},
				Objects: map[*ast.Object]*dst.Object{},
			},
		},
		Filenames: map[*dst.File]string{},
	}
}

type Decorator struct {
	Map
	Fset      *token.FileSet
	Filenames map[*dst.File]string // Source file names

	// If a Resolver is provided, it is used to resolve Ident nodes. During decoration, remote
	// identifiers (e.g. usually part of a qualified identifier SelectorExpr, but sometimes on
	// their own for dot-imported packages) are updated with the path of the package they are
	// imported from.
	Resolver resolver.IdentResolver
}

// Decorate decorates an ast.Node and returns a dst.Node
func (d *Decorator) Decorate(n ast.Node) dst.Node {

	fd := newFileDecorator(d)
	fd.fragment(n)
	fd.link()
	out := fd.decorateNode(n)

	//fmt.Println("\Fragments:")
	//fd.debug(os.Stdout)

	//fmt.Println("\nDecorator:")
	//debug(os.Stdout, out)

	// Populate Info with filenames if we're decorating a File or Package.
	switch n := n.(type) {
	case *ast.Package:
		for k, v := range n.Files {
			d.Filenames[d.Dst.Nodes[v].(*dst.File)] = k
		}
	case *ast.File:
		d.Filenames[out.(*dst.File)] = d.Fset.File(n.Pos()).Name()
	}

	return out
}

type decorationInfo struct {
	name string
	decs []string
}

func (f *fileDecorator) resolvePath(id *ast.Ident) string {
	if f.Resolver == nil {
		return ""
	}
	return f.Resolver.ResolveIdent(id)
}

func (f *fileDecorator) decorateObject(o *ast.Object) *dst.Object {
	if o == nil {
		return nil
	}
	if do, ok := f.Dst.Objects[o]; ok {
		return do
	}
	/*
		// An Object describes a named language entity such as a package,
		// constant, type, variable, function (incl. methods), or label.
		//
		// The Data fields contains object-specific data:
		//
		//	Kind    Data type         Data value
		//	Pkg     *Scope            package scope
		//	Con     int               iota for the respective declaration
		//
		type Object struct {
			Kind ObjKind
			Name string      // declared name
			Decl interface{} // corresponding Field, XxxSpec, FuncDecl, LabeledStmt, AssignStmt, Scope; or nil
			Data interface{} // object-specific data; or nil
			Type interface{} // placeholder for type information; may be nil
		}
	*/

	out := &dst.Object{}
	f.Dst.Objects[o] = out
	f.Ast.Objects[out] = o
	out.Kind = dst.ObjKind(o.Kind)
	out.Name = o.Name

	switch decl := o.Decl.(type) {
	case *ast.Scope:
		out.Decl = f.decorateScope(decl)
	case ast.Node:
		out.Decl = f.decorateNode(decl)
	case nil:
	default:
		panic(fmt.Sprintf("o.Decl is %T", o.Decl))
	}

	// TODO: I believe Data is either a *Scope or an int. We will support both and panic if something else if found.
	switch data := o.Data.(type) {
	case int:
		out.Data = data
	case *ast.Scope:
		out.Data = f.decorateScope(data)
	case ast.Node:
		out.Data = f.decorateNode(data)
	case nil:
	default:
		panic(fmt.Sprintf("o.Data is %T", o.Data))
	}

	return out
}

func (f *fileDecorator) decorateScope(s *ast.Scope) *dst.Scope {
	if s == nil {
		return nil
	}
	if ds, ok := f.Dst.Scopes[s]; ok {
		return ds
	}
	/*
		// A Scope maintains the set of named language entities declared
		// in the scope and a link to the immediately surrounding (outer)
		// scope.
		//
		type Scope struct {
			Outer   *Scope
			Objects map[string]*Object
		}
	*/
	out := &dst.Scope{}

	f.Dst.Scopes[s] = out
	f.Ast.Scopes[out] = s

	out.Outer = f.decorateScope(s.Outer)
	out.Objects = map[string]*dst.Object{}
	for k, v := range s.Objects {
		out.Objects[k] = f.decorateObject(v)
	}

	return out
}

func debug(w io.Writer, file dst.Node) {
	var result string
	nodeType := func(n dst.Node) string {
		return strings.Replace(fmt.Sprintf("%T", n), "*dst.", "", -1)
	}
	dst.Inspect(file, func(n dst.Node) bool {
		if n == nil {
			return false
		}
		var out string
		space, after, infos := getDecorationInfo(n)
		switch space {
		case dst.NewLine:
			out += " [New line space]"
		case dst.EmptyLine:
			out += " [Empty line space]"
		}
		for _, info := range infos {
			if len(info.decs) > 0 {
				var values string
				for i, dec := range info.decs {
					if i > 0 {
						values += " "
					}
					values += fmt.Sprintf("%q", dec)
				}
				out += fmt.Sprintf(" [%s %s]", info.name, values)
			}
		}
		switch after {
		case dst.NewLine:
			out += " [New line after]"
		case dst.EmptyLine:
			out += " [Empty line after]"
		}
		if out != "" {
			result += nodeType(n) + out + "\n"
		}
		return true
	})
	fmt.Fprint(w, result)
}
