package decorator

import (
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
)

// NewRestorer returns a restorer.
func NewRestorer() *Restorer {
	return &Restorer{
		Map:  newMap(),
		Fset: token.NewFileSet(),
	}
}

// NewRestorerWithImports returns a restorer with import management attributes set.
func NewRestorerWithImports(path string, resolver resolver.PackageResolver) *Restorer {
	return &Restorer{
		Map:      newMap(),
		Fset:     token.NewFileSet(),
		Path:     path,
		Resolver: resolver,
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
	// Local package path - required if Resolver is set.
	Path string
}

// Print uses format.Node to print a *dst.File to stdout
func (pr *Restorer) Print(f *dst.File) error {
	return pr.Fprint(os.Stdout, f)
}

// Fprint uses format.Node to print a *dst.File to a writer
func (pr *Restorer) Fprint(w io.Writer, f *dst.File) error {
	af, err := pr.RestoreFile(f)
	if err != nil {
		return err
	}
	return format.Node(w, pr.Fset, af)
}

// RestoreFile restores a *dst.File to an *ast.File
func (pr *Restorer) RestoreFile(file *dst.File) (*ast.File, error) {
	return pr.FileRestorer().RestoreFile(file)
}

func (pr *Restorer) FileRestorer() *FileRestorer {
	return &FileRestorer{
		Restorer: pr,
		Alias:    map[string]string{},
	}
}

type FileRestorer struct {
	*Restorer
	Alias           map[string]string // Map of package path -> package alias for imports
	Name            string            // The name of the restored file in the FileSet. Can usually be left empty.
	file            *dst.File
	lines           []int
	comments        []*ast.CommentGroup
	base            int
	cursor          token.Pos
	nodeDecl        map[*ast.Object]dst.Node // Objects that have a ast.Node Decl (look up after file has been rendered)
	nodeData        map[*ast.Object]dst.Node // Objects that have a ast.Node Data (look up after file has been rendered)
	cursorAtNewLine token.Pos                // The cursor position directly after adding a newline decoration (or a line comment which ends in a "\n"). If we're still at this cursor position when we add a line space, reduce the "\n" by one.
	packageNames    map[string]string        // names in the code of all imported packages ("." for dot-imports)
}

// Print uses format.Node to print a *dst.File to stdout
func (r *FileRestorer) Print(f *dst.File) error {
	return r.Fprint(os.Stdout, f)
}

// Fprint uses format.Node to print a *dst.File to a writer
func (r *FileRestorer) Fprint(w io.Writer, f *dst.File) error {
	af, err := r.RestoreFile(f)
	if err != nil {
		return err
	}
	return format.Node(w, r.Fset, af)
}

// RestoreFile restores a *dst.File to *ast.File
func (r *FileRestorer) RestoreFile(file *dst.File) (*ast.File, error) {

	if r.Resolver == nil && r.Path != "" {
		panic("Restorer Path should be empty when Resolver is nil")
	}

	if r.Resolver != nil && r.Path == "" {
		panic("Restorer Path should be set when Resolver is set")
	}

	if r.Fset == nil {
		r.Fset = token.NewFileSet()
	}

	// reset the FileRestorer, but leave Name and the Alias map unchanged

	r.file = file
	r.lines = []int{0} // initialise with the first line at Pos 0
	r.nodeDecl = map[*ast.Object]dst.Node{}
	r.nodeData = map[*ast.Object]dst.Node{}
	r.packageNames = map[string]string{}
	r.comments = []*ast.CommentGroup{}
	r.cursorAtNewLine = 0
	r.packageNames = map[string]string{}

	r.base = r.Fset.Base() // base is the pos that the file will start at in the fset
	r.cursor = token.Pos(r.base)

	if err := r.updateImports(); err != nil {
		return nil, err
	}

	// restore the file, populate comments and lines
	f := r.restoreNode(r.file, "", "", "", false).(*ast.File)

	for _, cg := range r.comments {
		f.Comments = append(f.Comments, cg)
	}

	size := r.fileSize()

	ff := r.Fset.AddFile(r.Name, r.base, size)
	if !ff.SetLines(r.lines) {
		panic("ff.SetLines failed")
	}

	// Sometimes new nodes are created here (e.g. in RangeStmt the "Object" is an AssignStmt which
	// never occurs in the actual code). These shouldn't have position information but perhaps it
	// doesn't matter?
	// TODO: Disable all position information on these nodes?
	for o, dn := range r.nodeDecl {
		o.Decl = r.restoreNode(dn, "", "", "", true)
	}
	for o, dn := range r.nodeData {
		o.Data = r.restoreNode(dn, "", "", "", true)
	}

	return f, nil
}

func (r *FileRestorer) updateImports() error {

	if r.Resolver == nil {
		return nil
	}

	// list of the import block(s)
	var blocks []*dst.GenDecl
	var hasCgoBlock bool

	// map of package path -> alias for all packages currently in the imports block(s)
	imports := map[string]string{}

	// list of package paths that occur in the source
	packages := map[string]bool{}

	all := map[string]bool{}
	var allOrdered []string

	dst.Inspect(r.file, func(n dst.Node) bool {
		switch n := n.(type) {
		case *dst.Ident:
			if n.Path == "" {
				return true
			}
			if n.Path == r.Path {
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
			// if this block has 1 spec and it's the "C" import, ignore it.
			if len(n.Specs) == 1 && mustUnquote(n.Specs[0].(*dst.ImportSpec).Path.Value) == "C" {
				hasCgoBlock = true
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
		name, err := r.Resolver.ResolvePackage(path)
		if err != nil {
			return err
		}
		resolved[path] = name
	}

	// make a list of packages we should remove from the import block(s)
	deletions := map[string]bool{}
	for path, alias := range imports {
		if alias == "_" || path == "C" || r.Alias[path] == "_" {
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
	for path, alias := range r.Alias {
		if alias == "_" && !all[path] {
			additions[path] = true
			all[path] = true
			allOrdered = append(allOrdered, path)
		}
	}

	sort.Slice(allOrdered, func(i, j int) bool { return packagePathOrderLess(allOrdered[i], allOrdered[j]) })

	// work out the actual aliases for all packages (and rename conflicts)
	aliases := map[string]string{}       // alias in the package block
	r.packageNames = map[string]string{} // name in the code

	conflict := func(name string) bool {
		for _, n := range r.packageNames {
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

		// if we set the alias to "_" but the package is still in use, we should ignore the alias
		ignoreAlias := r.Alias[path] == "_" && packages[path]

		var alias string
		if a, ok := r.Alias[path]; ok && !ignoreAlias {
			// If we have provided a custom alias, use this
			alias = a
		} else if a, ok := imports[path]; ok {
			// ... otherwise use the alias from the existing imports block
			alias = a
		}

		if alias == "." {
			// no conflict checking for dot-imports
			aliases[path] = "."
			r.packageNames[path] = ""
			continue
		}
		if alias == "_" {
			if unanon[path] {
				// for anonymous imports that we are converting to regular imports...
				r.packageNames[path], aliases[path] = findAlias(path, "")
			} else {
				// no conflict checking for anonymous imports
				aliases[path] = "_"
				r.packageNames[path] = ""
			}
			continue
		}
		r.packageNames[path], aliases[path] = findAlias(path, alias)
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
			gd := &dst.GenDecl{
				Tok: token.IMPORT,
				// make sure it has an empty line before and after
				Decs: dst.GenDeclDecorations{
					NodeDecs: dst.NodeDecs{Before: dst.EmptyLine, After: dst.EmptyLine},
				},
			}
			if hasCgoBlock {
				// special case for if we have the "C" import
				r.file.Decls = append([]dst.Decl{r.file.Decls[0], gd}, r.file.Decls[1:]...)
			} else {
				r.file.Decls = append([]dst.Decl{gd}, r.file.Decls...)
			}
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

		blocks[0].Specs = specs
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

	if len(additions) > 0 {
		// imports with a period in the path are assumed to not be standard library packages, so
		// get a newline separating them from standard library packages. We remove any other
		// newlines found in this block. We do this after the deletions because the first non-stdlib
		// import might be deleted.
		var foundDomainImport bool
		for _, spec := range blocks[0].Specs {
			path := mustUnquote(spec.(*dst.ImportSpec).Path.Value)
			if strings.Contains(path, ".") && !foundDomainImport {
				// first non-std-lib import -> empty line above
				spec.Decorations().Before = dst.EmptyLine
				spec.Decorations().After = dst.NewLine
				foundDomainImport = true
				continue
			}
			// all other specs, just newlines
			spec.Decorations().Before = dst.NewLine
			spec.Decorations().After = dst.NewLine
		}

		if len(blocks[0].Specs) == 1 {
			blocks[0].Lparen = false
			blocks[0].Rparen = false
		} else {
			blocks[0].Lparen = true
			blocks[0].Rparen = true
		}
	}

	if len(deleteBlocks) > 0 {
		decls := make([]dst.Decl, 0, len(r.file.Decls))
		for _, decl := range r.file.Decls {
			if deleteBlocks[decl] {
				continue
			}
			decls = append(decls, decl)
		}
		r.file.Decls = decls
	}

	return nil
}

func (r *FileRestorer) restoreIdent(n *dst.Ident, parentName, parentField, parentFieldType string, allowDuplicate bool) ast.Node {

	var name string
	if r.Resolver != nil && n.Path != "" {

		if avoid[parentName+"."+parentField] {
			panic(fmt.Sprintf("Path %s set on illegal Ident %s: parentName %s, parentField %s, parentFieldType %s", n.Path, n.Name, parentName, parentField, parentFieldType))
		}

		if n.Path != r.Path {
			name = r.packageNames[n.Path]
		}

		if name == "." {
			name = ""
		}
	}

	if name == "" {
		// continue to run standard Ident restore
		return nil
	}

	// restore to a SelectorExpr
	out := &ast.SelectorExpr{}
	r.Ast.Nodes[n] = out
	r.Dst.Nodes[out] = n
	r.Dst.Nodes[out.Sel] = n
	r.Dst.Nodes[out.X] = n
	r.applySpace(n.Decs.Before)

	// Decoration: Start
	r.applyDecorations(out, n.Decs.Start, false)

	// Node: X
	out.X = r.restoreNode(dst.NewIdent(name), "SelectorExpr", "X", "Expr", allowDuplicate).(ast.Expr)

	// Token: Period
	r.cursor += token.Pos(len(token.PERIOD.String()))

	// Decoration: X
	r.applyDecorations(out, n.Decs.X, false)

	// Node: Sel
	out.Sel = r.restoreNode(dst.NewIdent(n.Name), "SelectorExpr", "Sel", "Ident", allowDuplicate).(*ast.Ident)

	// Decoration: End
	r.applyDecorations(out, n.Decs.End, true)
	r.applySpace(n.Decs.After)

	return out

}

func packagePathOrderLess(pi, pj string) bool {
	// package paths with a . should be ordered after those without
	idot := strings.Contains(pi, ".")
	jdot := strings.Contains(pj, ".")
	if idot != jdot {
		return jdot
	}

	return pi < pj
}

func (r *FileRestorer) fileSize() int {

	// If a comment is at the end of a file, it will extend past the current cursor position...

	end := int(r.cursor) // end pos of file

	// check that none of the comments or newlines extend past the file end position. If so, increment.
	for _, cg := range r.comments {
		if int(cg.End()) >= end {
			end = int(cg.End()) + 1
		}
	}
	for _, lineOffset := range r.lines {
		pos := lineOffset + r.base // remember lines are relative to the file base
		if pos >= end {
			end = pos + 1
		}
	}

	return end - r.base
}

func (r *FileRestorer) applyLiteral(text string) {
	isMultiLine := strings.HasPrefix(text, "`") && strings.Contains(text, "\n")
	if !isMultiLine {
		return
	}
	for charIndex, char := range text {
		if char == '\n' {
			lineOffset := int(r.cursor) - r.base + charIndex // remember lines are relative to the file base
			r.lines = append(r.lines, lineOffset)
		}
	}
}

func (r *FileRestorer) hasCommentField(n ast.Node) bool {
	switch n.(type) {
	case *ast.Field, *ast.ValueSpec, *ast.TypeSpec, *ast.ImportSpec:
		return true
	}
	return false
}

func (r *FileRestorer) addCommentField(n ast.Node, slash token.Pos, text string) {
	c := &ast.Comment{Slash: slash, Text: text}
	switch n := n.(type) {
	case *ast.Field:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			r.comments = append(r.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	case *ast.ImportSpec:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			r.comments = append(r.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	case *ast.ValueSpec:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			r.comments = append(r.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	case *ast.TypeSpec:
		if n.Comment == nil {
			n.Comment = &ast.CommentGroup{}
			r.comments = append(r.comments, n.Comment)
		}
		n.Comment.List = append(n.Comment.List, c)
	}
}

func (r *FileRestorer) applyDecorations(node ast.Node, decorations dst.Decorations, end bool) {
	firstLine := true
	for _, d := range decorations {

		isNewline := d == "\n"
		isLineComment := strings.HasPrefix(d, "//")
		isInlineComment := strings.HasPrefix(d, "/*")
		isComment := isLineComment || isInlineComment
		isMultiLineComment := isInlineComment && strings.Contains(d, "\n")

		if end && r.cursorAtNewLine == r.cursor {
			r.cursor++ // indent all comments in "End" decorations
		}

		// for multi-line comments, add a newline for each \n
		if isMultiLineComment {
			for charIndex, char := range d {
				if char == '\n' {
					lineOffset := int(r.cursor) - r.base + charIndex // remember lines are relative to the file base
					r.lines = append(r.lines, lineOffset)
				}
			}
		}

		// if the decoration is a comment, add it and advance the cursor
		if isComment {
			if firstLine && end && r.hasCommentField(node) {
				// for comments on the same line as the end of a node that has a Comment field, we
				// add the comment to the node instead of the file.
				r.addCommentField(node, r.cursor, d)
			} else {
				r.comments = append(r.comments, &ast.CommentGroup{List: []*ast.Comment{{Slash: r.cursor, Text: d}}})
			}
			r.cursor += token.Pos(len(d))
		}

		// for newline decorations and also line-comments, add a newline
		if isLineComment || isNewline {
			lineOffset := int(r.cursor) - r.base // remember lines are relative to the file base
			r.lines = append(r.lines, lineOffset)
			r.cursor++

			r.cursorAtNewLine = r.cursor
		}

		if isNewline || isLineComment {
			firstLine = false
		}
	}
}

func (r *FileRestorer) applySpace(space dst.SpaceType) {
	var newlines int
	switch space {
	case dst.NewLine:
		newlines = 1
	case dst.EmptyLine:
		newlines = 2
	}
	if r.cursor == r.cursorAtNewLine {
		newlines--
	}
	for i := 0; i < newlines; i++ {

		// Advance the cursor one more byte for all newlines, so we step over any required
		// separator char - e.g. comma. See net-hook test
		r.cursor++

		lineOffset := int(r.cursor) - r.base // remember lines are relative to the file base
		r.lines = append(r.lines, lineOffset)
		r.cursor++
		r.cursorAtNewLine = r.cursor
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
