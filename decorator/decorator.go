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

// New returns a new decorator
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

	// If a Resolver is provided, it is used to resolve remote identifiers from *ast.Ident and
	// *ast.SelectorExpr nodes. Usually a remote identifier is a SelectorExpr qualified identifier,
	// but in the case of dot-imports they can be simply Ident nodes. During decoration, remote
	// identifiers are replaced with *dst.Ident with Path set to the path of imported package.
	Resolver resolver.IdentResolver
	// Local package path - only required if Resolver is set.
	Path string
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

	if d.Resolver == nil && d.Path != "" {
		panic("Decorator Path should be empty when Resolver is nil")
	}

	if d.Resolver != nil && d.Path == "" {
		panic("Decorator Path should be set when Resolver is set")
	}

	fd := d.newFileDecorator()
	if f, ok := n.(*ast.File); ok {
		fd.file = f
	}
	fd.fragment(n)
	fd.link()

	out, err := fd.decorateNode(nil, "", "", "", n)
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

// Never need to resolve idents that are in these fields (decorateSelectorExpr will force
// SelectorExpr.Sel to be resolved)
var avoid = map[string]bool{
	"Field.Names":       true,
	"LabeledStmt.Label": true,
	"BranchStmt.Label":  true,
	"ImportSpec.Name":   true,
	"ValueSpec.Names":   true,
	"TypeSpec.Name":     true,
	"FuncDecl.Name":     true,
	"File.Name":         true,
	"SelectorExpr.Sel":  true,
}

var decorateAvoid = map[string]bool{
	"SelectorExpr.X": true, // we avoid this in the decorator, but not the restorer
}

func (f *fileDecorator) decorateSelectorExpr(parent ast.Node, parentName, parentField, parentFieldType string, n *ast.SelectorExpr) (dst.Node, error) {

	var path string
	if f.Resolver != nil {
		p, err := f.resolvePath(true, n, "", "", "", n.Sel)
		if err != nil {
			return nil, err
		}
		path = p
	}

	if path == "" {
		// continue to default logic in decorateNode
		return nil, nil
	}

	// replace *ast.SelectorExpr with *dst.Ident and merge decorations
	out := &dst.Ident{}
	f.Dst.Nodes[n] = out
	f.Dst.Nodes[n.X] = out
	f.Dst.Nodes[n.Sel] = out
	f.Ast.Nodes[out] = n

	out.Decs.Before = mergeLineSpace(f.before[n], f.before[n.Sel], f.before[n.X])
	out.Decs.After = mergeLineSpace(f.after[n], f.after[n.Sel], f.after[n.X])

	// String: Name
	out.Name = n.Sel.Name

	// Object: Obj
	ob, err := f.decorateObject(n.Sel.Obj)
	if err != nil {
		return nil, err
	}
	out.Obj = ob

	// Path: Path
	out.Path = path

	nd, nok := f.decorations[n]
	xd, xok := f.decorations[n.X]
	sd, sok := f.decorations[n.Sel]

	if nok {
		if decs, ok := nd["Start"]; ok {
			out.Decs.Start.Append(decs...)
		}
	}

	if xok {
		if decs, ok := xd["Start"]; ok {
			out.Decs.Start.Append(decs...)
		}
		if decs, ok := xd["End"]; ok {
			out.Decs.X.Append(decs...)
		}
	}

	if nok {
		if decs, ok := nd["X"]; ok {
			out.Decs.X.Append(decs...)
		}
	}

	if sok {
		if decs, ok := sd["Start"]; ok {
			out.Decs.X.Append(decs...)
		}
		if decs, ok := sd["End"]; ok {
			out.Decs.End.Append(decs...)
		}
	}

	if nok {
		if decs, ok := nd["End"]; ok {
			out.Decs.End.Append(decs...)
		}
	}

	return out, nil

}

func (f *fileDecorator) resolvePath(force bool, parent ast.Node, parentName, parentField, parentFieldType string, id *ast.Ident) (string, error) {

	if !force {
		if f.Resolver == nil {
			panic("resolvePath needs a Resolver")
		}

		key := parentName + "." + parentField
		if avoid[key] || decorateAvoid[key] {
			return "", nil
		}

		if parentFieldType != "Expr" {
			panic(fmt.Sprintf("decorateIdent: unsupported parentName %s, parentField %s, parentFieldType %s", parentName, parentField, parentFieldType))
		}
	}

	path, err := f.Resolver.ResolveIdent(f.file, parent, id)
	if err != nil {
		return "", err
	}

	path = stripVendor(path)

	if path == stripVendor(f.Path) {
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
		n, err := f.decorateNode(nil, "", "", "", decl)
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
		n, err := f.decorateNode(nil, "", "", "", data)
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

func mergeLineSpace(spaces ...dst.SpaceType) dst.SpaceType {
	var hasNewLine bool
	for _, v := range spaces {
		switch v {
		case dst.EmptyLine:
			return dst.EmptyLine
		case dst.NewLine:
			hasNewLine = true
		}
	}
	if hasNewLine {
		return dst.NewLine
	}
	return dst.None
}
