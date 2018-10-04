package decorator

import (
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"

	"strings"

	"github.com/dave/dst"
)

func Print(f *dst.File) error {
	return Fprint(os.Stdout, f)
}

func Fprint(w io.Writer, f *dst.File) error {
	fset, af := Restore(f)
	return format.Node(w, fset, af)
}

func Restore(file *dst.File) (*token.FileSet, *ast.File) {
	r := NewRestorer()
	return r.Fset, r.RestoreFile("", file)
}

func NewRestorer() *Restorer {
	return &Restorer{
		Fset:  token.NewFileSet(),
		Nodes: map[dst.Node]ast.Node{},
	}
}

type Restorer struct {
	Fset  *token.FileSet
	Nodes map[dst.Node]ast.Node
}

type fileRestorer struct {
	*Restorer
	lines    []int
	comments []*ast.CommentGroup
	base     int
	cursor   token.Pos
}

func (r *Restorer) RestoreFile(name string, file *dst.File) *ast.File {

	// Base is the pos that the file will start at in the fset
	base := r.Fset.Base()

	fr := &fileRestorer{
		Restorer: r,
		lines:    []int{0}, // initialise with the first line at Pos 0
		base:     base,
		cursor:   token.Pos(base),
	}

	// restore the file, populate comments and lines
	f := fr.restoreNode(file).(*ast.File)

	for _, cg := range fr.comments {
		f.Comments = append(f.Comments, cg)
	}

	size := fr.fileSize()

	ff := r.Fset.AddFile(name, base, size)
	if !ff.SetLines(fr.lines) {
		panic("SetLines failed")
	}

	return f
}

func (f *fileRestorer) fileSize() int {

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

func (f *fileRestorer) applyDecorations(decorations dst.Decorations) {
	for _, d := range decorations {

		isNewline := d == "\n"
		isLineComment := strings.HasPrefix(d, "//")
		isInlineComment := strings.HasPrefix(d, "/*")
		isComment := isLineComment || isInlineComment
		isMultiLineComment := isInlineComment && strings.Contains(d, "\n")

		// for multi-line comments, add a newline for each \n
		if isMultiLineComment {
			for i, char := range d {
				if char == '\n' {
					lineOffset := int(f.cursor) - f.base + i // remember lines are relative to the file base
					f.lines = append(f.lines, lineOffset)
				}
			}
		}

		// if the decoration is a comment, add it and advance the cursor
		if isComment {
			f.comments = append(f.comments, &ast.CommentGroup{List: []*ast.Comment{{Slash: f.cursor, Text: d}}})
			f.cursor += token.Pos(len(d))
		}

		// for newline decorations and also line-comments, add a newline
		if isLineComment || isNewline {
			lineOffset := int(f.cursor) - f.base // remember lines are relative to the file base
			f.lines = append(f.lines, lineOffset)
			f.cursor++
		}
	}
}
