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
	"github.com/dave/dst/decorator/resolver/gotypes"
	"golang.org/x/tools/go/packages"
)

// New returns a new package decorator
func New(fset *token.FileSet) *Decorator {
	return &Decorator{
		Map:       newMap(),
		Filenames: map[*dst.File]string{},
		Fset:      fset,
	}
}

// NewWithImports returns a new package decorator with import management attributes set.
func NewWithImports(pkg *packages.Package) *Decorator {
	return &Decorator{
		Map:       newMap(),
		Filenames: map[*dst.File]string{},
		Fset:      pkg.Fset,
		Path:      pkg.PkgPath,
		Resolver: &gotypes.IdentResolver{
			Info: pkg.TypesInfo,
		},
	}
}

type Decorator struct {
	Map
	Filenames map[*dst.File]string // Source file names
	Fset      *token.FileSet

	Path      string // local package path, used to ensure the local path is not set in idents
	canonical string // local package path, de-vendored, set in DecorateNode

	// If a Resolver is provided, it is used to resolve Ident nodes. During decoration, remote
	// identifiers (e.g. usually part of a qualified identifier SelectorExpr, but sometimes on
	// their own for dot-imported packages) are updated with the path of the package they are
	// imported from.
	Resolver resolver.IdentResolver
}

// Parse uses parser.ParseFile to parse and decorate a Go source file. The src parameter should
// be string, []byte, or io.Reader.
func (d *Decorator) Parse(src interface{}) (*dst.File, error) {
	return d.ParseFile("", src, parser.ParseComments)
}

// ParseFile uses parser.ParseFile to parse and decorate a Go source file. The ParseComments flag is
// added to mode if it doesn't exist.
func (d *Decorator) ParseFile(filename string, src interface{}, mode parser.Mode) (*dst.File, error) {

	// If ParseFile returns an error and also a non-nil file, the errors were just parse errors so
	// we should continue decorating the file and return the error.
	f, perr := parser.ParseFile(d.Fset, filename, src, mode|parser.ParseComments)
	if perr != nil && f == nil {
		return nil, perr
	}

	file, err := d.DecorateFile(f)
	if err != nil {
		return nil, err
	}

	return file, perr
}

// ParseDir uses parser.ParseDir to parse and decorate a directory containing Go source. The
// ParseComments flag is added to mode if it doesn't exist.
func (d *Decorator) ParseDir(dir string, filter func(os.FileInfo) bool, mode parser.Mode) (map[string]*dst.Package, error) {
	pkgs, err := parser.ParseDir(d.Fset, dir, filter, mode|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	out := map[string]*dst.Package{}
	for k, v := range pkgs {
		pkg, err := d.DecorateNode(v)
		if err != nil {
			return nil, err
		}
		out[k] = pkg.(*dst.Package)
	}
	return out, nil
}

func (d *Decorator) DecorateFile(f *ast.File) (*dst.File, error) {
	file, err := d.DecorateNode(f)
	if err != nil {
		return nil, err
	}
	return file.(*dst.File), nil
}

// Decorate decorates an ast.Node and returns a dst.Node
func (d *Decorator) DecorateNode(n ast.Node) (dst.Node, error) {

	d.canonical = stripVendor(d.Path)

	fd := d.newFileDecorator()
	if f, ok := n.(*ast.File); ok {
		fd.file = f
	}
	fd.fragment(n)
	fd.link()

	out, err := fd.decorateNode(nil, "", n)
	if err != nil {
		return nil, err
	}

	//fmt.Println("\nFragments:")
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

	return out, nil
}

func (pd *Decorator) newFileDecorator() *fileDecorator {
	return &fileDecorator{
		Decorator:    pd,
		startIndents: map[ast.Node]int{},
		endIndents:   map[ast.Node]int{},
		before:       map[ast.Node]dst.SpaceType{},
		after:        map[ast.Node]dst.SpaceType{},
		decorations:  map[ast.Node]map[string][]string{},
	}
}

type fileDecorator struct {
	*Decorator
	file          *ast.File // file we're decorating in for import name resolution - can be nil if we're just decorating an isolated node
	cursor        int
	fragments     []fragment
	startIndents  map[ast.Node]int
	endIndents    map[ast.Node]int
	before, after map[ast.Node]dst.SpaceType
	decorations   map[ast.Node]map[string][]string
}

type decorationInfo struct {
	name string
	decs []string
}

func (f *fileDecorator) resolvePath(parent ast.Node, typ string, id *ast.Ident) (string, error) {

	if f.Resolver == nil {
		return "", nil
	}

	// The parent field type (typ) for all idents is either "Ident" or "Expr".
	//
	// If the parent field type is Ident, there is no possibility of this field holding a
	// SelectorExpr, so this ident cannot possibly be a qualified identifier. We avoid resolving
	// the Path for these idents.
	//
	// Inside SelectorExpr is a special case where the logic is reversed. We avoid setting Path for
	// X (Expr) but set it for Sel (Ident).
	if _, sel := parent.(*ast.SelectorExpr); (typ == "Ident" && !sel) || (typ == "Expr" && sel) {
		return "", nil
	}

	path, err := f.Resolver.ResolveIdent(f.file, parent, id)
	if err != nil {
		return "", err
	}

	if path == f.canonical {
		return "", nil
	}

	return path, nil
}

func stripVendor(path string) string {
	findVendor := func(path string) (index int, ok bool) {
		// Two cases, depending on internal at start of string or not.
		// The order matters: we must return the index of the final element,
		// because the final one is where the effective import path starts.
		switch {
		case strings.Contains(path, "/vendor/"):
			return strings.LastIndex(path, "/vendor/") + 1, true
		case strings.HasPrefix(path, "vendor/"):
			return 0, true
		}
		return 0, false
	}
	i, ok := findVendor(path)
	if !ok {
		return path
	}
	return path[i+len("vendor/"):]
}

func (f *fileDecorator) decorateObject(o *ast.Object) (*dst.Object, error) {
	if o == nil {
		return nil, nil
	}
	if do, ok := f.Dst.Objects[o]; ok {
		return do, nil
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
		s, err := f.decorateScope(decl)
		if err != nil {
			return nil, err
		}
		out.Decl = s
	case ast.Node:
		n, err := f.decorateNode(nil, "", decl)
		if err != nil {
			return nil, err
		}
		out.Decl = n
	case nil:
	default:
		panic(fmt.Sprintf("o.Decl is %T", o.Data))
	}

	switch data := o.Data.(type) {
	case int:
		out.Data = data
	case *ast.Scope:
		s, err := f.decorateScope(data)
		if err != nil {
			return nil, err
		}
		out.Data = s
	case ast.Node:
		n, err := f.decorateNode(nil, "", data)
		if err != nil {
			return nil, err
		}
		out.Data = n
	case nil:
	default:
		panic(fmt.Sprintf("o.Data is %T", o.Data))
	}

	return out, nil
}

func (f *fileDecorator) decorateScope(s *ast.Scope) (*dst.Scope, error) {
	if s == nil {
		return nil, nil
	}
	if ds, ok := f.Dst.Scopes[s]; ok {
		return ds, nil
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

	outer, err := f.decorateScope(s.Outer)
	if err != nil {
		return nil, err
	}
	out.Outer = outer
	out.Objects = map[string]*dst.Object{}
	for k, v := range s.Objects {
		ob, err := f.decorateObject(v)
		if err != nil {
			return nil, err
		}
		out.Objects[k] = ob
	}

	return out, nil
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
		before, after, infos := getDecorationInfo(n)
		switch before {
		case dst.NewLine:
			out += " [New line before]"
		case dst.EmptyLine:
			out += " [Empty line before]"
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
