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
	af, fset := Restore(f)
	return format.Node(w, fset, af)
}

func Restore(file *dst.File) (*ast.File, *token.FileSet) {
	fset := token.NewFileSet()
	return RestoreNamed("a.go", file, fset), fset
}

func RestoreNamed(name string, file *dst.File, fset *token.FileSet) *ast.File {
	r := &restorer{Fset: fset}
	return r.restore(name, file)
}

type restorer struct {
	Fset *token.FileSet
}

type fileRestorer struct {
	*restorer
	Lines    []int
	Comments []*ast.CommentGroup
	base     token.Pos
	cursor   token.Pos
	nodes    map[dst.Node]ast.Node
}

func (r *restorer) restore(fname string, dstFile *dst.File) *ast.File {
	if r.Fset == nil {
		r.Fset = token.NewFileSet()
	}
	fr := &fileRestorer{
		restorer: r,
		base:     token.Pos(r.Fset.Base()),
		cursor:   token.Pos(r.Fset.Base()),
		nodes:    map[dst.Node]ast.Node{},
		Lines:    []int{0},
	}
	astFile := fr.restoreNode(dstFile).(*ast.File)

	astFileEndPos := int(fr.cursor - fr.base)
	// Check that none of the comments or newlines extend past the file end position. If so, increment.
	for _, cg := range fr.Comments {
		if int(cg.End()) >= astFileEndPos {
			astFileEndPos = int(cg.End()) + 1
		}
	}
	for _, l := range fr.Lines {
		if l >= astFileEndPos {
			astFileEndPos = l + 1
		}
	}

	fsetFile := r.Fset.AddFile(fname, r.Fset.Base(), astFileEndPos)
	for _, cg := range fr.Comments {
		astFile.Comments = append(astFile.Comments, cg)
	}
	success := fsetFile.SetLines(fr.Lines)
	if !success {
		panic("SetLines failed")
	}

	return astFile
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
			for i, c := range d {
				if c == '\n' {
					f.Lines = append(f.Lines, int(f.cursor)+i)
				}
			}
		}

		// if the decoration is a comment, add it and advance the cursor
		if isComment {
			f.Comments = append(f.Comments, &ast.CommentGroup{List: []*ast.Comment{{Slash: f.cursor, Text: d}}})
			f.cursor += token.Pos(len(d))
		}

		// for newline decorations and also line-comments, add a newline
		if isLineComment || isNewline {
			f.Lines = append(f.Lines, int(f.cursor))
			f.cursor++
		}
	}
}
