package decorator

import (
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/dstutil"
)

// Print uses format.Node to print a *dst.File to stdout
func Print(f *dst.File) error {
	return Fprint(os.Stdout, f)
}

// Fprint uses format.Node to print a *dst.File to a writer
func Fprint(w io.Writer, f *dst.File) error {
	fset, af := Restore(f)
	return format.Node(w, fset, af)
}

// Restore restores a *dst.File to a *token.FileSet and a *ast.File
func Restore(file *dst.File) (*token.FileSet, *ast.File) {
	r := NewRestorer()
	return r.Fset, r.RestoreFile("", file)
}

// NewRestorer creates a new Restorer
func NewRestorer() *Restorer {
	return &Restorer{
		Fset: token.NewFileSet(),
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
	}
}

type Restorer struct {
	Map
	Fset *token.FileSet // Fset is the *token.FileSet in use. Set this to use a pre-existing FileSet.

	// If a Resolver is provided, the names of all imported packages are resolved, and the imports
	// block is updated. All remote identifiers are updated (sometimes this involves changing
	// SelectorExpr.X.Name, or even swapping between Ident and SelectorExpr). To force specific
	// import alias names, use the FileRestorer.Alias map.
	Resolver resolver.PackageResolver
}

// NewPackageRestorer returns a restorer which resolves package names relative to a specific dir
// and local path.
func (r *Restorer) NewPackageRestorer(path, dir string) *PackageRestorer {
	return &PackageRestorer{
		Restorer: r,
		Path:     path,
		Dir:      dir,
	}
}

type PackageRestorer struct {
	*Restorer
	Dir  string // source dir for package name resolution
	Path string // local package path for identifier resolution
}

type Map struct {
	Ast AstMap
	Dst DstMap
}

type AstMap struct {
	Nodes   map[dst.Node]ast.Node       // Mapping from dst to ast Nodes
	Objects map[*dst.Object]*ast.Object // Mapping from dst to ast Objects
	Scopes  map[*dst.Scope]*ast.Scope   // Mapping from dst to ast Scopes
}

type DstMap struct {
	Nodes   map[ast.Node]dst.Node       // Mapping from ast to dst Nodes
	Objects map[*ast.Object]*dst.Object // Mapping from ast to dst Objects
	Scopes  map[*ast.Scope]*dst.Scope   // Mapping from ast to dst Scopes
}

type FileRestorer struct {
	*PackageRestorer
	Alias           map[string]string // map of package path -> package alias for imports
	name            string
	file            *dst.File
	lines           []int
	comments        []*ast.CommentGroup
	base            int
	cursor          token.Pos
	nodeDecl        map[*ast.Object]dst.Node // Objects that have a ast.Node Decl (look up after file has been rendered)
	nodeData        map[*ast.Object]dst.Node // Objects that have a ast.Node Data (look up after file has been rendered)
	cursorAtNewLine token.Pos                // The cursor position directly after adding a newline decoration (or a line comment which ends in a "\n"). If we're still at this cursor position when we add a line space, reduce the "\n" by one.
}

// RestoreFile restores a *dst.File to an *ast.File
func (r *Restorer) RestoreFile(name string, file *dst.File) *ast.File {
	return r.NewPackageRestorer("", "").RestoreFile(name, file)
}

// RestoreFile restores a *dst.File to an *ast.File
func (pr *PackageRestorer) RestoreFile(name string, file *dst.File) *ast.File {
	return pr.NewFileRestorer(name, file).RestoreFile(context.Background())
}

func (pr *PackageRestorer) NewFileRestorer(name string, file *dst.File) *FileRestorer {
	return &FileRestorer{
		PackageRestorer: pr,
		Alias:           map[string]string{},
		name:            name,
		file:            file,
		lines:           []int{0}, // initialise with the first line at Pos 0
		nodeDecl:        map[*ast.Object]dst.Node{},
		nodeData:        map[*ast.Object]dst.Node{},
	}
}

func (fr *FileRestorer) RestoreFile(ctx context.Context) *ast.File {

	fr.base = fr.Fset.Base() // base is the pos that the file will start at in the fset
	fr.cursor = token.Pos(fr.base)

	fr.updateImports(ctx)

	// restore the file, populate comments and lines
	f := fr.restoreNode(fr.file, false).(*ast.File)

	for _, cg := range fr.comments {
		f.Comments = append(f.Comments, cg)
	}

	size := fr.fileSize()

	ff := fr.Fset.AddFile(fr.name, fr.base, size)
	if !ff.SetLines(fr.lines) {
		panic("SetLines failed")
	}

	// Sometimes new nodes are created here (e.g. in RangeStmt the "Object" is an AssignStmt which
	// never occurs in the actual code). These shouldn't have position information but perhaps it
	// doesn't matter?
	// TODO: Disable all position information on these nodes?
	for o, dn := range fr.nodeDecl {
		o.Decl = fr.restoreNode(dn, true)
	}
	for o, dn := range fr.nodeData {
		o.Data = fr.restoreNode(dn, true)
	}

	return f
}

func (fr *FileRestorer) updateImports(ctx context.Context) {

	if fr.Resolver == nil {
		return
	}

	// list of the import block(s)
	var blocks []*dst.GenDecl

	// map of package path -> alias for all packages currently in the imports block(s)
	imports := map[string]string{}

	// list of package paths that occur in the source
	packages := map[string]bool{}

	all := map[string]bool{}
	var allOrdered []string

	dst.Inspect(fr.file, func(n dst.Node) bool {
		switch n := n.(type) {
		case *dst.Ident:
			if n.Path == "" {
				return true
			}
			if _, ok := packages[n.Path]; !ok {
				packages[n.Path] = true
				all[n.Path] = true
				allOrdered = append(allOrdered, n.Path)
			}

		case *dst.GenDecl:
			if n.Tok != token.IMPORT {
				return true
			}
			blocks = append(blocks, n)

		case *dst.ImportSpec:
			path := mustUnquote(n.Path.Value)
			if n.Name == nil {
				imports[path] = ""
			} else {
				imports[path] = n.Name.Name
			}
		}
		return true
	})

	// resolved names of all packages in use
	resolved := map[string]string{}

	for path := range packages {
		name, err := fr.Resolver.ResolvePackage(ctx, path, fr.Dir)
		if err != nil {
			panic(err)
		}
		resolved[path] = name
	}

	// make a list of packages we should remove from the import block(s)
	deletions := map[string]bool{}
	for path, alias := range imports {
		if alias == "_" || path == "C" {
			// never remove anonymous imports, or the "C" import
			if !all[path] {
				all[path] = true
				allOrdered = append(allOrdered, path)
			}
			continue
		}
		if packages[path] {
			// if the package is still in use, don't remove
			continue
		}
		deletions[path] = true
	}

	// make a list of the packages we should add to the first import block, and the packages we
	// should convert from anonymous import to normal import
	additions := map[string]bool{}
	unanon := map[string]bool{}
	for path := range packages {
		if alias, ok := imports[path]; ok {
			if alias == "_" {
				unanon[path] = true
			}
		} else {
			additions[path] = true
		}
	}

	// any anonymous imports manually added with FileRestorer.Alias should be added too
	for path, alias := range fr.Alias {
		if alias == "_" && !all[path] {
			additions[path] = true
		}
	}

	sort.Slice(allOrdered, func(i, j int) bool { return packagePathOrderLess(allOrdered[i], allOrdered[j]) })

	// work out the actual aliases for all packages (and rename conflicts)
	aliases := map[string]string{} // alias in the package block
	names := map[string]string{}   // name in the code

	conflict := func(name string) bool {
		for _, n := range names {
			if name == n {
				return true
			}
		}
		return false
	}

	// Finds a unique alias given a path and a preferred alias. If preferred == "", we look up the
	// name of the package in the resolved map.
	findAlias := func(path, preferred string) (name, alias string) {

		// if we pass in a preferred alias we should always return an alias even when the alias
		// matches the package name.
		aliased := preferred != ""

		if preferred == "" {
			preferred = resolved[path]
		}

		modifier := 1
		current := preferred
		for conflict(current) {
			current = fmt.Sprintf("%s%d", preferred, modifier)
			modifier++
		}

		if !aliased && current == resolved[path] {
			return current, ""
		}

		return current, current
	}

	for _, path := range allOrdered {
		if deletions[path] {
			// ignore if it's going to be deleted
			continue
		}

		var alias string
		if a, ok := fr.Alias[path]; ok {
			// If we have provided a custom alias, use this
			alias = a
		} else if a, ok := imports[path]; ok {
			// ... otherwise use the alias from the existing imports block
			alias = a
		}

		if alias == "." {
			// no conflict checking for dot-imports
			aliases[path] = "."
			names[path] = ""
			continue
		}
		if alias == "_" {
			if unanon[path] {
				// for anonymous imports that we are converting to regular imports...
				names[path], aliases[path] = findAlias(path, "")
			} else {
				// no conflict checking for anonymous imports
				aliases[path] = "_"
				names[path] = ""
			}
			continue
		}
		names[path], aliases[path] = findAlias(path, alias)
	}

	// convert any anonymous imports to regular imports
	if len(unanon) > 0 {
		for _, blk := range blocks {
			for _, spec := range blk.Specs {
				spec := spec.(*dst.ImportSpec)
				path := mustUnquote(spec.Path.Value)
				if !unanon[path] {
					continue
				}
				if aliases[path] == "" {
					spec.Name = nil
					continue
				}
				spec.Name = &dst.Ident{Name: aliases[path]}
			}
		}
	}

	// update the alias for any that need it
	for _, blk := range blocks {
		for _, spec := range blk.Specs {
			spec := spec.(*dst.ImportSpec)
			path := mustUnquote(spec.Path.Value)
			if spec.Name == nil && aliases[path] != "" {
				// missing alias
				spec.Name = &dst.Ident{Name: aliases[path]}
			} else if spec.Name != nil && aliases[path] == "" {
				// alias needs to be removed
				spec.Name = nil
			} else if spec.Name != nil && aliases[path] != spec.Name.Name {
				// alias wrong
				spec.Name.Name = aliases[path]
			}
		}
	}

	// make any additions
	if len(additions) > 0 {

		// if there's currently no import blocks, we must create one
		if len(blocks) == 0 {
			gd := &dst.GenDecl{Tok: token.IMPORT}
			fr.file.Decls = append([]dst.Decl{gd}, fr.file.Decls...)
			blocks = append(blocks, gd)
		}

		specs := blocks[0].Specs

		for path := range additions {
			is := &dst.ImportSpec{
				Path: &dst.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", path)},
			}
			if aliases[path] != "" {
				is.Name = &dst.Ident{
					Name: aliases[path],
				}
			}
			specs = append(specs, is)
		}

		// rearrange import block
		sort.Slice(specs, func(i, j int) bool {
			return packagePathOrderLess(
				mustUnquote(specs[i].(*dst.ImportSpec).Path.Value),
				mustUnquote(specs[j].(*dst.ImportSpec).Path.Value),
			)
		})

		// imports with a period in the path are assumed to not be standard library packages, so
		// get a newline separating them from standard library packages. We remove any other
		// newlines found in this block.
		var foundDomainImport bool
		for _, spec := range specs {
			path := mustUnquote(spec.(*dst.ImportSpec).Path.Value)
			if strings.Contains(path, ".") && !foundDomainImport {
				// first non-std-lib import -> empty line above
				spec.Decorations().Space = dst.EmptyLine
				spec.Decorations().After = dst.NewLine
				foundDomainImport = true
				continue
			}
			// all other specs, just newlines
			spec.Decorations().Space = dst.NewLine
			spec.Decorations().After = dst.NewLine
		}

		blocks[0].Specs = specs

		if len(specs) == 1 {
			blocks[0].Lparen = false
			blocks[0].Rparen = false
		} else {
			blocks[0].Lparen = true
			blocks[0].Rparen = true
		}
	}

	deleteBlocks := map[dst.Decl]bool{}

	if len(deletions) > 0 {
		for _, blk := range blocks {
			specs := make([]dst.Spec, 0, len(blk.Specs))
			for _, spec := range blk.Specs {
				path := mustUnquote(spec.(*dst.ImportSpec).Path.Value)
				if deletions[path] {
					continue
				}
				specs = append(specs, spec)
			}
			blk.Specs = specs
			if len(specs) == 0 {
				deleteBlocks[blk] = true
			} else if len(specs) == 1 {
				blk.Lparen = false
				blk.Rparen = false
			} else {
				blk.Lparen = true
				blk.Rparen = true
			}
		}
	}

	if len(deleteBlocks) > 0 {
		decls := make([]dst.Decl, 0, len(fr.file.Decls))
		for _, decl := range fr.file.Decls {
			if deleteBlocks[decl] {
				continue
			}
			decls = append(decls, decl)
		}
		fr.file.Decls = decls
	}

	// update the SelectorExpr and Ident in the rest of the file
	dstutil.Apply(fr.file, func(c *dstutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *dst.SelectorExpr:
			if n.Sel.Path == "" {
				return true
			}
			x := n.X.(*dst.Ident)
			sel := n.Sel
			if names[n.Sel.Path] == "" {
				// blank name -> replace this with Ident
				id := &dst.Ident{
					Name: n.Sel.Name,
					Path: n.Sel.Path,
					Obj:  n.Sel.Obj,
				}

				// merge decorations for n, x and sel:

				if n.Decs.Space == dst.EmptyLine || x.Decs.Space == dst.EmptyLine || sel.Decs.Space == dst.EmptyLine {
					id.Decs.Space = dst.EmptyLine
				} else if n.Decs.Space == dst.NewLine || x.Decs.Space == dst.NewLine || sel.Decs.Space == dst.NewLine {
					id.Decs.Space = dst.NewLine
				}

				if n.Decs.After == dst.EmptyLine || x.Decs.After == dst.EmptyLine || sel.Decs.After == dst.EmptyLine {
					id.Decs.After = dst.EmptyLine
				} else if n.Decs.After == dst.NewLine || x.Decs.After == dst.NewLine || sel.Decs.After == dst.NewLine {
					id.Decs.After = dst.NewLine
				}

				// add all decorations
				id.Decs.Start.Append([]string(n.Decs.Start)...)
				id.Decs.Start.Append([]string(x.Decs.Start)...)
				id.Decs.Start.Append([]string(x.Decs.End)...)
				id.Decs.Start.Append([]string(n.Decs.X)...)
				id.Decs.Start.Append([]string(sel.Decs.Start)...)
				id.Decs.End.Append([]string(sel.Decs.End)...)
				id.Decs.End.Append([]string(n.Decs.End)...)

				c.Replace(id)
			} else if names[n.Sel.Path] != x.Name {
				// update name
				x.Name = names[n.Sel.Path]
			}
		case *dst.Ident:
			if n.Path == "" {
				return true
			}
			if _, ok := c.Parent().(*dst.SelectorExpr); ok {
				// skip idents inside SelectorExpr
				return true
			}
			if names[n.Path] != "" {
				// add a SelectorExpr
				sel := &dst.SelectorExpr{
					X:   &dst.Ident{Name: names[n.Path]},
					Sel: &dst.Ident{Name: n.Name, Path: n.Path, Obj: n.Obj},
				}
				sel.Decs.Space = n.Decs.Space
				sel.Decs.After = n.Decs.After
				sel.Decs.Start.Append([]string(n.Decs.Start)...)
				sel.Decs.End.Append([]string(n.Decs.End)...)
				c.Replace(sel)
			}
		}
		return true
	}, nil)

}

func packagePathOrderLess(pi, pj string) bool {
	// "C" import should be last
	ic := pi == "C"
	jc := pj == "C"
	if ic != jc {
		return jc
	}

	// package paths with a . should be ordered after those without
	idot := strings.Contains(pi, ".")
	jdot := strings.Contains(pj, ".")
	if idot != jdot {
		return jdot
	}

	return pi < pj
}

func (f *FileRestorer) fileSize() int {

	// If a comment is at the end of a file, it will extend past the current cursor position...

	end := int(f.cursor) // end pos of file

	// check that none of the comments or newlines extend past the file end position. If so, increment.
	for _, cg := range f.comments {
		if int(cg.End()) >= end {
			end = int(cg.End()) + 1
		}
	}
	for _, lineOffset := range f.lines {
		pos := lineOffset + f.base // remember lines are relative to the file base
		if pos >= end {
			end = pos + 1
		}
	}

	return end - f.base
}

func (f *FileRestorer) applyLiteral(text string) {
	isMultiLine := strings.HasPrefix(text, "`") && strings.Contains(text, "\n")
	if !isMultiLine {
		return
	}
	for charIndex, char := range text {
		if char == '\n' {
			lineOffset := int(f.cursor) - f.base + charIndex // remember lines are relative to the file base
			f.lines = append(f.lines, lineOffset)
		}
	}
}

func (f *FileRestorer) hasCommentField(n ast.Node) bool {
	switch n.(type) {
	case *ast.Field, *ast.ValueSpec, *ast.TypeSpec, *ast.ImportSpec:
		return true
	}
	return false
}

func (f *FileRestorer) addCommentField(n ast.Node, slash token.Pos, text string) {
	c := &ast.Comment{Slash: slash, Text: text}
	switch n := n.(type) {
	case *ast.Field:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			f.comments = append(f.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	case *ast.ImportSpec:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			f.comments = append(f.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	case *ast.ValueSpec:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			f.comments = append(f.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	case *ast.TypeSpec:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			f.comments = append(f.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	}
}

func (f *FileRestorer) applyDecorations(node ast.Node, decorations dst.Decorations, end bool) {
	firstLine := true
	for _, d := range decorations {

		isNewline := d == "\n"
		isLineComment := strings.HasPrefix(d, "//")
		isInlineComment := strings.HasPrefix(d, "/*")
		isComment := isLineComment || isInlineComment
		isMultiLineComment := isInlineComment && strings.Contains(d, "\n")

		if end && f.cursorAtNewLine == f.cursor {
			f.cursor++ // indent all comments in "End" decorations
		}

		// for multi-line comments, add a newline for each \n
		if isMultiLineComment {
			for charIndex, char := range d {
				if char == '\n' {
					lineOffset := int(f.cursor) - f.base + charIndex // remember lines are relative to the file base
					f.lines = append(f.lines, lineOffset)
				}
			}
		}

		// if the decoration is a comment, add it and advance the cursor
		if isComment {
			if firstLine && end && f.hasCommentField(node) {
				// for comments on the same line as the end of a node that has a Comment field, we
				// add the comment to the node instead of the file.
				f.addCommentField(node, f.cursor, d)
			} else {
				f.comments = append(f.comments, &ast.CommentGroup{List: []*ast.Comment{{Slash: f.cursor, Text: d}}})
			}
			f.cursor += token.Pos(len(d))
		}

		// for newline decorations and also line-comments, add a newline
		if isLineComment || isNewline {
			lineOffset := int(f.cursor) - f.base // remember lines are relative to the file base
			f.lines = append(f.lines, lineOffset)
			f.cursor++

			f.cursorAtNewLine = f.cursor
		}

		if isNewline || isLineComment {
			firstLine = false
		}
	}
}

func (f *FileRestorer) applySpace(space dst.SpaceType) {
	var newlines int
	switch space {
	case dst.NewLine:
		newlines = 1
	case dst.EmptyLine:
		newlines = 2
	}
	if f.cursor == f.cursorAtNewLine {
		newlines--
	}
	for i := 0; i < newlines; i++ {

		// Advance the cursor one more byte for all newlines, so we step over any required
		// separator char - e.g. comma. See net-hook test
		f.cursor++

		lineOffset := int(f.cursor) - f.base // remember lines are relative to the file base
		f.lines = append(f.lines, lineOffset)
		f.cursor++
		f.cursorAtNewLine = f.cursor
	}
}

func (r *FileRestorer) restoreObject(o *dst.Object) *ast.Object {
	if o == nil {
		return nil
	}
	if ro, ok := r.Ast.Objects[o]; ok {
		return ro
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
	out := &ast.Object{}

	r.Ast.Objects[o] = out
	r.Dst.Objects[out] = o

	out.Kind = ast.ObjKind(o.Kind)
	out.Name = o.Name

	switch decl := o.Decl.(type) {
	case *dst.Scope:
		out.Decl = r.restoreScope(decl)
	case dst.Node:
		// Can't use restoreNode here because we aren't at the right cursor position, so we store a link
		// to the Object and Node so we can look the Nodes up in the cache after the file is fully processed.
		r.nodeDecl[out] = decl
	case nil:
	default:
		panic(fmt.Sprintf("o.Decl is %T", o.Decl))
	}

	// TODO: I believe Data is either a *Scope or an int. We will support both and panic if something else if found.
	switch data := o.Data.(type) {
	case int:
		out.Data = data
	case *dst.Scope:
		out.Data = r.restoreScope(data)
	case dst.Node:
		// Can't use restoreNode here because we aren't at the right cursor position, so we store a link
		// to the Object and Node so we can look the Nodes up in the cache after the file is fully processed.
		r.nodeData[out] = data
	case nil:
	default:
		panic(fmt.Sprintf("o.Data is %T", o.Data))
	}

	return out
}

func (r *FileRestorer) restoreScope(s *dst.Scope) *ast.Scope {
	if s == nil {
		return nil
	}
	if rs, ok := r.Ast.Scopes[s]; ok {
		return rs
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
	out := &ast.Scope{}

	r.Ast.Scopes[s] = out
	r.Dst.Scopes[out] = s

	out.Outer = r.restoreScope(s.Outer)
	out.Objects = map[string]*ast.Object{}
	for k, v := range s.Objects {
		out.Objects[k] = r.restoreObject(v)
	}

	return out
}

func mustUnquote(s string) string {
	out, err := strconv.Unquote(s)
	if err != nil {
		panic(err)
	}
	return out
}
