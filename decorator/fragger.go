package decorator

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"sort"
	"strings"

	"github.com/dave/dst"
)

type Fragger struct {
	Fragments []Fragment
}

func (f *Fragger) ProcessToken(n ast.Node, name string, pos token.Pos, end bool) {
	f.Fragments = append(f.Fragments, NodeFragment{Node: n, Name: name, End: end, Pos: pos})
}

func (f *Fragger) Fragment(file *ast.File, fset *token.FileSet) {

	f.ProcessNode(file)

	// make a list of the comments that are already included in the AST - we don't need to include
	// them as decorations
	comments := map[*ast.Comment]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if c, ok := n.(*ast.Comment); ok {
			comments[c] = true
		}
		return true
	})

	lineComments := map[int]bool{}
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if comments[c] {
				continue
			}
			f.Fragments = append(f.Fragments, CommentFragment{Pos: c.Pos(), Text: c.Text})
			if strings.HasPrefix(c.Text, "//") {
				lineComments[fset.Position(c.Pos()).Line] = true
			} else if strings.HasPrefix(c.Text, "/*") && fset.Position(c.End()).Line > fset.Position(c.Pos()).Line {
				// multi line comment
				for i := fset.Position(c.Pos()).Line; i < fset.Position(c.End()).Line; i++ {
					lineComments[i] = true
				}
			}
		}
	}

	fsetFile := fset.File(file.Pos())
	line := 1
	for i := fsetFile.Base(); i < fsetFile.Base()+fsetFile.Size(); i++ {
		pos := fset.Position(token.Pos(i))
		if pos.Line != line {
			if !lineComments[line] {
				f.Fragments = append(f.Fragments, NewlineFragment{Pos: token.Pos(i - 1)})
			}
			line = pos.Line
		}
	}

	sort.SliceStable(f.Fragments, func(i, j int) bool {
		return f.Fragments[i].Position() < f.Fragments[j].Position()
	})
}

func (f *Fragger) Link() map[ast.Node][]dst.Decoration {
	out := map[ast.Node][]dst.Decoration{}
	currentNodeFragment, found := f.nextNodeFragment(0)
	if !found {
		return out
	}
	var foundNewlines bool
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case NodeFragment:
			currentNodeFragment = frag
			foundNewlines = false
		case CommentFragment:
			if foundNewlines {
				// if the comment was preceded by newlines, add it to the next fragment
				foundNewlines = false
				if n, found := f.nextNodeFragment(i); found {
					currentNodeFragment = n
				}
			}
			out[currentNodeFragment.Node] = append(out[currentNodeFragment.Node], dst.Decoration{Position: currentNodeFragment.Name, End: currentNodeFragment.End, Text: frag.Text})
			foundNewlines = strings.HasPrefix(frag.Text, "//")
		case NewlineFragment:
			out[currentNodeFragment.Node] = append(out[currentNodeFragment.Node], dst.Decoration{Position: currentNodeFragment.Name, End: currentNodeFragment.End, Text: "\n"})
			foundNewlines = true
		}
	}
	return out
}

func (f *Fragger) nextNodeFragment(from int) (frag NodeFragment, found bool) {
	for i := from; i < len(f.Fragments); i++ {
		switch frag := f.Fragments[i].(type) {
		case NodeFragment:
			return frag, true
		}
	}
	return NodeFragment{}, false
}

type Fragment interface {
	Position() token.Pos
}

type NodeFragment struct {
	Node ast.Node
	Name string
	End  bool
	Pos  token.Pos
}

type CommentFragment struct {
	Pos  token.Pos
	Text string
}

type NewlineFragment struct {
	Pos token.Pos
}

func (v NodeFragment) Position() token.Pos    { return v.Pos }
func (v CommentFragment) Position() token.Pos { return v.Pos }
func (v NewlineFragment) Position() token.Pos { return v.Pos }

func (f Fragger) debug(w io.Writer, fset *token.FileSet) {
	formatPos := func(s token.Position) string {
		return s.String()[strings.Index(s.String(), ":")+1:]
	}
	nodeType := func(n ast.Node) string {
		return strings.Replace(fmt.Sprintf("%T", n), "*ast.", "", -1)
	}
	for i, v := range f.Fragments {
		switch v := v.(type) {
		case NodeFragment:
			var name string
			if v.Name != "" {
				name = ":" + v.Name
			}
			var pos string
			if v.End {
				pos = " (end)"
			} else {
				pos = " (start)"
			}
			fmt.Fprintf(w, "%d %s%s%s %s\n", i, nodeType(v.Node), name, pos, formatPos(fset.Position(v.Pos)))
		case CommentFragment:
			fmt.Fprintf(w, "%d * Comment %s %s\n", i, formatPos(fset.Position(v.Pos)), v.Text)
		case NewlineFragment:
			fmt.Fprintf(w, "%d * Newline %s\n", i, formatPos(fset.Position(v.Pos)))
		}
	}
}

func (f *Fragger) funcDeclOverride(n *ast.FuncDecl) {
	// Doc
	if n.Doc != nil {
		f.ProcessToken(n, "Doc", n.Doc.Pos(), false)
		f.ProcessNode(n.Doc)
		f.ProcessToken(n, "Doc", n.Doc.End(), true)
	}
	// Func
	if n.Type.Func.IsValid() {
		f.ProcessToken(n, "Func", n.Type.Func, true)
	}
	// Recv
	if n.Recv != nil {
		f.ProcessToken(n, "Recv", n.Recv.Pos(), false)
		f.ProcessNode(n.Recv)
		f.ProcessToken(n, "Recv", n.Recv.End(), true)
	}
	// Name
	if n.Name != nil {
		f.ProcessToken(n, "Name", n.Name.Pos(), false)
		f.ProcessNode(n.Name)
		f.ProcessToken(n, "Name", n.Name.End(), true)
	}
	// Params
	if n.Type.Params != nil {
		f.ProcessToken(n, "Params", n.Type.Params.Pos(), false)
		f.ProcessNode(n.Type.Params)
		f.ProcessToken(n, "Params", n.Type.Params.End(), true)
	}
	// Results
	if n.Type.Results != nil {
		f.ProcessToken(n, "Results", n.Type.Results.Pos(), false)
		f.ProcessNode(n.Type.Results)
		f.ProcessToken(n, "Results", n.Type.Results.End(), true)
	}
	// Body
	if n.Body != nil {
		f.ProcessToken(n, "Body", n.Body.Pos(), false)
		f.ProcessNode(n.Body)
		f.ProcessToken(n, "Body", n.Body.End(), true)
	}
}
