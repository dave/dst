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

// NewDecorator returns a new decorator. If fset is nil, a new FileSet is created.
func NewDecorator(fset *token.FileSet) *Decorator {
	if fset == nil {
		fset = token.NewFileSet()
	}
	return &Decorator{
		Map:       newMap(),
		Filenames: map[*dst.File]string{},
		Fset:      fset,
	}
}

// NewDecoratorWithImports returns a new decorator with import management enabled.
func NewDecoratorWithImports(fset *token.FileSet, path string, resolver resolver.DecoratorResolver) *Decorator {
	dec := NewDecorator(fset)
	dec.Path = path
	dec.Resolver = resolver
	return dec
}

// NewDecoratorFromPackage returns a new decorator configured to decorate files in pkg.
func NewDecoratorFromPackage(pkg *packages.Package) *Decorator {
	return NewDecoratorWithImports(pkg.Fset, pkg.PkgPath, gotypes.New(pkg.TypesInfo.Uses))
}

// Decorator converts ast nodes into dst nodes, and converts decoration info from the ast fileset
// to the dst nodes. Create a new Decorator for each package you need to decorate.
type Decorator struct {
	Map                            // Mapping between ast and dst Nodes, Objects and Scopes
	Filenames map[*dst.File]string // Source file names
	Fset      *token.FileSet       // The ast FileSet containing ast decoration info for the files

	// If a Resolver is provided, it is used to resolve remote identifiers from *ast.Ident and
	// *ast.SelectorExpr nodes. Usually a remote identifier is a SelectorExpr qualified identifier,
	// but in the case of dot-imports they can be simply Ident nodes. During decoration, remote
	// identifiers are replaced with *dst.Ident with Path set to the path of imported package.
	Resolver resolver.DecoratorResolver
	// Local package path - required if Resolver is set.
	Path string
}

// Parse uses parser.ParseFile to parse and decorate a Go source file. The src parameter should
// be string, []byte, or io.Reader.
func (d *Decorator) Parse(src interface{}) (*dst.File, error) {
	return d.ParseFile("", src, parser.ParseComments)
}

// ParseFile uses parser.ParseFile to parse and decorate a Go source file. The ParseComments flag
// is added to mode if it doesn't exist.
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

// DecorateFile decorates *ast.File and returns *dst.File
func (d *Decorator) DecorateFile(f *ast.File) (*dst.File, error) {
	file, err := d.DecorateNode(f)
	if err != nil {
		return nil, err
	}
	return file.(*dst.File), nil
}

// DecorateNode decorates ast.Node and returns dst.Node
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

// We never need to resolve idents that are in these fields (decorateSelectorExpr will override
// this check with the force parameter for SelectorExpr.Sel when needed).
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

// decorateSelectorExpr is a special case for decorating a SelectorExpr, which might return an
// Ident if the resolver determines that the SelectorExpr represents a qualified ident.
func (f *fileDecorator) decorateSelectorExpr(parent ast.Node, parentName, parentField, parentFieldType string, n *ast.SelectorExpr) (dst.Node, error) {

	if f.Resolver == nil {
		// continue to default logic in decorateNode
		return nil, nil
	}

	// resolve the path with force == true, so we skip the tests that would normally prevent the
	// Sel field of SelectorExpr from being resolved.
	path, err := f.resolvePath(true, n, "SelectorExpr", "Sel", "Ident", n.Sel)
	if err != nil {
		return nil, err
	}

	if path == "" {
		// path == "" -> not a qualified ident -> continue to default logic in decorateNode
		return nil, nil
	}

	// replace *ast.SelectorExpr with *dst.Ident and merge decorations
	out := &dst.Ident{}
	f.Dst.Nodes[n] = out
	f.Dst.Nodes[n.X] = out
	f.Dst.Nodes[n.Sel] = out
	f.Ast.Nodes[out] = n

	/*
		This is rather messy. We must merge the SelectorExpr decorations into an Ident. The Ident
		has an X decoration attachment point, but we don't have a simple place to merge the X.After
		and Sel.Before line-spacing. This is rather an edge case, but we can fix it by converting
		the line-spacing to "\n" decorations before / after the X decoration. This will at least
		mean that decorated / restored code with no mutations should be byte-perfect.

		Here's a list of the decorations we're merging:

		{1}{2}{3}{4}[  X  ].{5}{6}{7}{8}{9}[ Sel ]{10}{11}{12}{13}

		1: SelectorExpr Before Space     - f.before[n]
		2: SelectorExpr Start Decoration - f.decorations[n]["Start"]
		3: X Before Space                - f.before[n.X]
		4: X Start Decoration            - f.decorations[n.X]["Start"]

		5: X End Decoration              - f.decorations[n.X]["End"]
		6: X After Space                 - f.after[n.X]
		7: SelectorExpr X Decoration     - f.decorations[n]["X"]
		8: Sel Before Space              - f.before[n.Sel]
		9: Sel Start Decoration          - f.decorations[n.Sel]["Start"]

		10: Sel End decoration           - f.decorations[n.Sel]["End"]
		11: Sel After Space              - f.after[n.Sel]
		12: SelectorExpr End Decoration  - f.decorations[n]["End"]
		13: SelectorExpr After Space     - f.after[n]

		1-4:   merge into Ident.Before / Ident.Start
		5-9:   merge into Ident.X (convert line spaces to decorations)
		10-13: merge into Ident.End / Ident.After
	*/

	out.Decs.Before = mergeLineSpace(f.before[n], f.before[n.X])
	out.Decs.After = mergeLineSpace(f.after[n], f.after[n.Sel])
	spaceBeforeX := f.after[n.X]
	spaceAfterX := f.before[n.Sel]

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

	var hasXDecorations bool
	if nok {
		if decs, ok := nd["X"]; ok {
			hasXDecorations = len(decs) > 0
		}
	}

	if !hasXDecorations {

		// if there's no x decoration, we should merge the two line spaces because they are
		// adjoining.
		mergedSpace := mergeLineSpace(spaceBeforeX, spaceAfterX)
		if mergedSpace == dst.NewLine {
			out.Decs.X.Append("\n")
		} else if mergedSpace == dst.EmptyLine {
			out.Decs.X.Append("\n", "\n")
		}

	} else {

		if spaceBeforeX == dst.NewLine {
			out.Decs.X.Append("\n")
		} else if spaceBeforeX == dst.EmptyLine {
			out.Decs.X.Append("\n", "\n")
		}

		// we know there's some x decorations, so no need for the checks
		decs := nd["X"]
		out.Decs.X.Append(decs...)

		// does the last x decoration introduce a new-line? (e.g. "//" comment or "\n")
		xDecorationEndsInNewline := decs[len(decs)-1] == "\n" || strings.HasPrefix(decs[len(decs)-1], "//")

		// we reduce the number of "\n" emitted if the last x-decoration adds a line ("//" or "\n")
		if spaceAfterX == dst.NewLine {
			if xDecorationEndsInNewline {
				// nothing to do
			} else {
				out.Decs.X.Append("\n")
			}
		} else if spaceAfterX == dst.EmptyLine {
			if xDecorationEndsInNewline {
				out.Decs.X.Append("\n")
			} else {
				out.Decs.X.Append("\n", "\n")
			}
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

	if f.Resolver == nil {
		panic("resolvePath needs a Resolver")
	}

	if !force {

		key := parentName + "." + parentField
		if avoid[key] {
			return "", nil
		}

		if parentFieldType != "Expr" {
			panic(fmt.Sprintf("decorateIdent: unsupported parentName %s, parentField %s, parentFieldType %s", parentName, parentField, parentFieldType))
		}
	}

	path, err := f.Resolver.ResolveIdent(f.file, parent, parentField, id)
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
