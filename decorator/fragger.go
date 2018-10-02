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
	cursor    int
	Fragments []Fragment
}

func (f *Fragger) AddStart(n ast.Node, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, DecorationFragment{Node: n, Name: "Start", Pos: token.Pos(f.cursor)})
}

func (f *Fragger) AddDecoration(n ast.Node, name string) {
	f.Fragments = append(f.Fragments, DecorationFragment{Node: n, Name: name, Pos: token.Pos(f.cursor)})
}

func (f *Fragger) AddToken(n ast.Node, t token.Token, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, TokenFragment{Node: n, Token: t, Pos: token.Pos(f.cursor)})
	f.cursor += len(t.String())
}

func (f *Fragger) AddString(n ast.Node, s string, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, StringFragment{Node: n, String: s, Pos: token.Pos(f.cursor)})
	f.cursor += len(s)
}

func (f *Fragger) AddComment(text string, pos token.Pos) {
	// Don't need to worry about the cursor with comments - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, CommentFragment{Text: text, Pos: pos})
}

func (f *Fragger) AddNewline(pos token.Pos) {
	// Don't need to worry about the cursor with newlines - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, NewlineFragment{Pos: pos})
}

func (f *Fragger) Fragment(file *ast.File, fset *token.FileSet) {

	f.cursor = fset.File(file.Pos()).Base()

	f.ProcessNode(file)

	newlinesInsideComments := map[int]bool{}
	for _, cg := range file.Comments {
		for _, c := range cg.List {

			// Add the comment to the fragment list.
			f.AddComment(c.Text, c.Slash)

			// Add any newlines.
			if strings.HasPrefix(c.Text, "//") {
				newlinesInsideComments[fset.Position(c.Pos()).Line] = true
			} else if strings.HasPrefix(c.Text, "/*") && fset.Position(c.End()).Line > fset.Position(c.Pos()).Line {
				// multi line comment
				for i := fset.Position(c.Pos()).Line; i < fset.Position(c.End()).Line; i++ {
					newlinesInsideComments[i] = true
				}
			}
		}
	}

	fsetFile := fset.File(file.Pos())
	line := 1
	for i := fsetFile.Base(); i < fsetFile.Base()+fsetFile.Size(); i++ {
		pos := fset.Position(token.Pos(i))
		if pos.Line != line {
			if !newlinesInsideComments[line] {
				f.AddNewline(token.Pos(i - 1))
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
	var lastIndex int
	var lastDecoration DecorationFragment
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case CommentFragment:

			if i < lastIndex {
				// If we've already moved something forward past this index, always move this to the
				// same decoration
				out[lastDecoration.Node] = append(out[lastDecoration.Node], dst.Decoration{Position: lastDecoration.Name, Text: frag.Text})
				continue
			}

			var dec DecorationFragment
			var found bool
			var try int
			var index int
			for !found {
				try++
				switch try {
				case 1:
					// search backwards but stop at any newline or token
					dec, index, found = f.nextDecoration(i, -1, true, true)
				case 2:
					// search forwards but stop at any token
					dec, index, found = f.nextDecoration(i, 1, false, true)
				case 3:
					// search backwards but stop at any token
					dec, index, found = f.nextDecoration(i, -1, false, true)
				case 4:
					// search forwards
					dec, index, found = f.nextDecoration(i, 1, false, false)
				case 5:
					// search backwards
					dec, index, found = f.nextDecoration(i, -1, false, false)
				default:
					panic("no decoration found for " + frag.Text)
				}
			}
			out[dec.Node] = append(out[dec.Node], dst.Decoration{Position: dec.Name, Text: frag.Text})
			lastIndex = index
			lastDecoration = dec

		case NewlineFragment:

			if i < lastIndex {
				// If we've already moved something forward past this index, always move this to the
				// same decoration
				out[lastDecoration.Node] = append(out[lastDecoration.Node], dst.Decoration{Position: lastDecoration.Name, Text: "\n"})
				continue
			}

			var dec DecorationFragment
			var found bool
			var try int
			var index int
			for !found {
				try++
				switch try {
				case 1:
					// search backwards but stop at any token
					dec, index, found = f.nextDecoration(i, -1, false, true)
				case 2:
					// search forwards but stop at any token
					dec, index, found = f.nextDecoration(i, 1, false, true)
				case 3:
					// search backwards
					dec, index, found = f.nextDecoration(i, -1, false, false)
				case 4:
					// search backwards
					dec, index, found = f.nextDecoration(i, 1, false, false)
				default:
					panic("no decoration found for newline")
				}
			}
			out[dec.Node] = append(out[dec.Node], dst.Decoration{Position: dec.Name, Text: "\n"})
			lastIndex = index
			lastDecoration = dec
		}
	}
	return out
}

func (f *Fragger) nextDecoration(from int, direction int, stopAtNewline bool, stopAtToken bool) (frag DecorationFragment, index int, found bool) {
	for i := from; i < len(f.Fragments) && i >= 0; i += direction {
		switch f := f.Fragments[i].(type) {
		case DecorationFragment:
			return f, i, true
		case NewlineFragment:
			if stopAtNewline {
				return
			}
		case TokenFragment, StringFragment:
			if stopAtToken {
				return
			}
		}
	}
	return
}

type Fragment interface {
	Position() token.Pos
}

type TokenFragment struct {
	Node  ast.Node
	Token token.Token
	Pos   token.Pos
}

type StringFragment struct {
	Node   ast.Node
	String string
	Pos    token.Pos
}

type CommentFragment struct {
	Text string
	Pos  token.Pos
}

type NewlineFragment struct {
	Pos token.Pos
}

type DecorationFragment struct {
	Node ast.Node
	Name string
	Pos  token.Pos
}

func (v TokenFragment) Position() token.Pos      { return v.Pos }
func (v StringFragment) Position() token.Pos     { return v.Pos }
func (v CommentFragment) Position() token.Pos    { return v.Pos }
func (v NewlineFragment) Position() token.Pos    { return v.Pos }
func (v DecorationFragment) Position() token.Pos { return v.Pos }

func (f Fragger) debug(w io.Writer, fset *token.FileSet) {
	formatPos := func(s token.Position) string {
		return s.String()[strings.Index(s.String(), ":")+1:]
	}
	nodeType := func(n ast.Node) string {
		return strings.Replace(fmt.Sprintf("%T", n), "*ast.", "", -1)
	}
	for _, v := range f.Fragments {
		switch v := v.(type) {
		case TokenFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.Token, formatPos(fset.Position(v.Pos)))
		case StringFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.String, formatPos(fset.Position(v.Pos)))
		case DecorationFragment:
			fmt.Fprintf(w, "%s %s %s\n", nodeType(v.Node), v.Name, formatPos(fset.Position(v.Pos)))
		case CommentFragment:
			fmt.Fprintf(w, "%q %s\n", v.Text, formatPos(fset.Position(v.Pos)))
		case NewlineFragment:
			fmt.Fprintf(w, "\"\\n\" %s\n", formatPos(fset.Position(v.Pos)))
		}
	}
}
