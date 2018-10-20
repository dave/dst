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

func NewFragger(fset *token.FileSet) *Fragger {
	return &Fragger{
		fset:         fset,
		IndentsStart: map[ast.Node]int{},
		IndentsEnd:   map[ast.Node]int{},
	}
}

type Fragger struct {
	cursor       int
	Fragments    []Fragment
	IndentsStart map[ast.Node]int
	IndentsEnd   map[ast.Node]int
	fset         *token.FileSet
}

func (f *Fragger) AddDecoration(n ast.Node, name string, pos token.Pos) {
	f.Fragments = append(f.Fragments, &DecorationFragment{Node: n, Name: name, Pos: token.Pos(f.cursor)})
}

func (f *Fragger) AddToken(n ast.Node, t token.Token, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, &TokenFragment{Node: n, Token: t, Pos: token.Pos(f.cursor)})
	f.cursor += len(t.String())
}

func (f *Fragger) AddString(n ast.Node, s string, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, &StringFragment{Node: n, String: s, Pos: token.Pos(f.cursor)})
	f.cursor += len(s)
}

func (f *Fragger) AddComment(text string, pos token.Pos) {
	// Don't need to worry about the cursor with comments - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, &CommentFragment{Text: text, Pos: pos})
}

func (f *Fragger) AddNewline(pos token.Pos, empty bool) {
	// Don't need to worry about the cursor with newlines - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, &NewlineFragment{Pos: pos, Empty: empty})
}

func (f *Fragger) Fragment(node ast.Node) {

	f.ProcessNode(node)

	if f.fset != nil {
		processFile := func(astf *ast.File) {
			// we will avoid adding a newline decoration that is inside a comment
			avoid := map[int]bool{}
			for _, cg := range astf.Comments {
				for _, c := range cg.List {

					// Add the comment to the fragment list.
					f.AddComment(c.Text, c.Slash)

					// Avoid newlines in multi-line comments
					if strings.HasPrefix(c.Text, "/*") {
						startLine := f.fset.Position(c.Pos()).Line
						endLine := f.fset.Position(c.End()).Line

						// multi line comment
						if endLine > startLine {
							for i := startLine; i < endLine; i++ {
								avoid[i+1] = true // we avoid the lines that follow the lines in the comment
							}
						}
					}
				}
			}

			// avoid newlines inside multi-line (back-quoted) strings
			for _, frag := range f.Fragments {
				switch frag := frag.(type) {
				case *StringFragment:
					if !strings.HasPrefix(frag.String, "`") {
						continue
					}

					startLine := f.fset.Position(frag.Pos).Line
					endLine := f.fset.Position(frag.Pos + token.Pos(len(frag.String))).Line

					// multi line string
					if endLine > startLine {
						for i := startLine; i < endLine; i++ {
							avoid[i+1] = true // we avoid the lines that follow the lines in the string
						}
					}
				}
			}

			line := 1
			tokenf := f.fset.File(astf.Pos())
			max := tokenf.Base() + tokenf.Size()
			for i := tokenf.Base(); i < max; i++ {
				pos := f.fset.Position(token.Pos(i))
				if pos.Line != line {

					line = pos.Line

					if avoid[line] {
						continue
					}

					// Peek ahead to the next position in the fset. If we're on another new line, we have
					// an empty line:
					nextLine := line
					if i < max-1 {
						// can't peek forward at the end of the file
						nextLine = f.fset.Position(token.Pos(i + 1)).Line
					}

					if nextLine != line {
						f.AddNewline(token.Pos(i-1), true)
						line = nextLine
						i++
					} else {
						f.AddNewline(token.Pos(i-1), false)
					}

				}
			}
		}

		switch val := node.(type) {
		case *ast.File:
			processFile(val)
		case *ast.Package:
			for _, file := range val.Files {
				processFile(file)
			}
		}

	}

	sort.SliceStable(f.Fragments, func(i, j int) bool {
		return f.Fragments[i].Position() < f.Fragments[j].Position()
	})

	// Search for nodes and comments that start directly after newlines. We note their indent.
	currentIndent := 0
	for i, frag := range f.Fragments {
		if i == 0 || f.Fragments[i-1].HasNewline() {
			currentIndent = f.fset.Position(frag.Position()).Column
		}
		switch frag := frag.(type) {
		case *DecorationFragment:
			switch frag.Name {
			case "Start":
				f.IndentsStart[frag.Node] = currentIndent
			case "End":
				f.IndentsEnd[frag.Node] = currentIndent
			}
		case *CommentFragment:
			frag.Indent = currentIndent
		}
	}
}

func appendDecoration(m map[ast.Node]map[string][]string, n ast.Node, pos, text string) {
	if m[n] == nil {
		m[n] = map[string][]string{}
	}
	m[n][pos] = append(m[n][pos], text)
}

func appendNewLine(m map[ast.Node]map[string][]string, n ast.Node, pos string, empty bool) {
	if m[n] == nil {
		m[n] = map[string][]string{}
	}
	num := 1
	if empty {
		num = 2
	}
	decs := m[n][pos]
	if len(decs) > 0 && strings.HasPrefix(decs[len(decs)-1], "//") {
		num--
	}
	for i := 0; i < num; i++ {
		m[n][pos] = append(m[n][pos], "\n")
	}
}

func (f *Fragger) Link() (space, after map[ast.Node]dst.SpaceType, decorations map[ast.Node]map[string][]string) {

	space = map[ast.Node]dst.SpaceType{}
	after = map[ast.Node]dst.SpaceType{}
	decorations = map[ast.Node]map[string][]string{}

	// Pass 1: associate comment groups with decorations. Sweep up any other comments / new-lines /
	// empty-lines and associate with the same decoration.
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case *DecorationFragment:

			// Special case for hanging indent (See https://github.com/dave/dst/issues/18)
			// If we're on the End decoration of a Stmt or Decl, and the start / end indents
			// are not the same (OR it's a case / comm clause), then search forward over empty lines
			// for all comments with the same indent as the End decoration.
			//
			// These should be attached to the end node. We also search for subsequent comments that
			// have the same indent as the Start. If the next decoration node is the start of a Stmt
			// or Decl with the same indent as the original node, these are attached there.

			if frag.Name != "End" {
				continue
			}
			_, stmt := frag.Node.(ast.Stmt)
			_, decl := frag.Node.(ast.Decl)
			if !stmt && !decl {
				continue
			}

			start := f.IndentsStart[frag.Node]
			end := f.IndentsEnd[frag.Node]

			_, labeledStmt := frag.Node.(*ast.LabeledStmt)

			if labeledStmt {
				// Special case: labeled statements shouldn't be treated in the same way.
				continue
			}

			_, caseClause := frag.Node.(*ast.CaseClause)
			_, commClause := frag.Node.(*ast.CommClause)
			if start == end && (caseClause || commClause) {
				// special case for case / comm clause with no items... the clause node starts and
				// ends on the same line, but comments can still be hanging. We spoof an indented
				// end position:
				end++
			}

			if start == end {
				continue
			}

			frags, next := f.findIndentedComments(i+1, [2]int{end, start})
			endFrags := frags[0]
			nextFrags := frags[1]
			if len(endFrags) > 0 {
				// if endFrags ends with a newline, don't attach it because it was in between the
				// two groups, so should be left unattached so we can attach it as spacing in the
				// second pass
				_, nl := endFrags[len(endFrags)-1].(*NewlineFragment)
				if nl {
					f.attachToDecoration(endFrags[0:len(endFrags)-1], decorations, frag)
				} else {
					f.attachToDecoration(endFrags, decorations, frag)
				}
			}
			if len(nextFrags) > 0 && next != nil {
				_, nextStmt := next.Node.(ast.Stmt)
				_, nextDecl := next.Node.(ast.Decl)
				nextStart := f.IndentsStart[next.Node]
				if (nextStmt || nextDecl) && nextStart == start {
					f.attachToDecoration(nextFrags, decorations, next)
				}
			}

		case *CommentFragment:

			if frag.Attached != nil {
				continue
			}

			// Comments (or comment groups) attach to decoration points in this precedence:
			//
			// 1) Before the comment on the same line
			// 2) After the comment on the same line
			// 3) After the comment on subsequent lines (but stopping at empty lines)
			// 4) Before the comment on previous lines (but stopping at empty lines)
			// 5) After the comment on subsequent lines
			// 6) Before the comment on previous lines
			//
			// We always stop at tokens, strings. If we get to the end without finding a decoration point we panic.

			var frags []Fragment // comment / new-line / empty-line
			var dec *DecorationFragment
			var found bool
			var try int
			for !found {
				try++
				switch try {
				case 1:
					// Before the comment on the same line (search backwards and stop at any newline)
					frags, dec, found = f.findDecoration(true, true, i, -1, false)
				case 2:
					// After the comment on the same line
					// After the comment on line+1 (search forwards and stop at any empty line)
					frags, dec, found = f.findDecoration(false, true, i, 1, false)
				case 3:
					// Before the comment on line-1 (search backwards and stop at any empty line)
					frags, dec, found = f.findDecoration(false, true, i, -1, false)
				case 4:
					// After the comment on line+2 (search forwards)
					frags, dec, found = f.findDecoration(false, false, i, 1, false)
				case 5:
					// After the comment on line-2 (search backwards)
					frags, dec, found = f.findDecoration(false, false, i, -1, false)
				default:
					panic("no decoration found for " + frag.Text)
				}
			}
			f.attachToDecoration(frags, decorations, dec)
		}
	}

	// Pass 2: associate any new-lines / empty-lines that have not been added to decorations to node
	// spacing. If they can't be attached as node spacing, attach them as decorations.
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case *NewlineFragment:

			if frag.Attached != nil {
				continue
			}

			// If the newline is directly before / after a node, we can set the Before / After spacing
			// of the node decoration instead of adding the newline as a decoration.
			nodeSpace, _, foundSpace := f.findNode(i, 1)
			nodeAfter, _, foundAfter := f.findNode(i, -1)
			if foundSpace || foundAfter {
				spaceType := dst.NewLine
				if frag.Empty {
					spaceType = dst.EmptyLine
				}
				if foundSpace {
					space[nodeSpace] = spaceType
				}
				if foundAfter {
					after[nodeAfter] = spaceType
				}
				continue
			}

			// If this newline can't be associated with a node, attach it to the next / previous
			// decoration location:
			var dec *DecorationFragment
			var found bool
			var try int
			for !found {
				try++
				switch try {
				case 1:
					// search backwards but stop at any token
					_, dec, found = f.findDecoration(false, false, i, -1, false)
				case 2:
					// search forwards but stop at any token
					_, dec, found = f.findDecoration(false, false, i, 1, false)
				default:
					panic("no decoration found for newline")
				}
			}
			appendNewLine(decorations, dec.Node, dec.Name, frag.Empty)
		}
	}

	return
}

func (f *Fragger) attachToDecoration(frags []Fragment, decorations map[ast.Node]map[string][]string, dec *DecorationFragment) {
	for _, fr := range frags {
		switch fr := fr.(type) {
		case *CommentFragment:
			appendDecoration(decorations, dec.Node, dec.Name, fr.Text)
			fr.Attached = dec
		case *NewlineFragment:
			appendNewLine(decorations, dec.Node, dec.Name, fr.Empty)
			fr.Attached = dec
		}
	}
}

func (f *Fragger) findDecoration(stopAtNewline, stopAtEmptyLine bool, from int, direction int, onlyClause bool) (swept []Fragment, dec *DecorationFragment, found bool) {
	var frags []Fragment
	for i := from; i < len(f.Fragments) && i >= 0; i += direction {
		switch current := f.Fragments[i].(type) {
		case *DecorationFragment:
			if onlyClause {
				switch current.Node.(type) {
				case *ast.CommClause, *ast.CaseClause:
					if current.Name == "Start" {
						return frags, current, true
					}
					return
				default:
					return
				}
			}
			return frags, current, true
		case *NewlineFragment:
			if stopAtNewline {
				return
			}
			if stopAtEmptyLine && current.Empty {
				return
			}
			if current.Attached != nil {
				continue
			}
			if direction == 1 {
				frags = append(frags, current)
			} else {
				frags = append([]Fragment{current}, frags...)
			}
		case *CommentFragment:
			if current.Attached != nil {
				continue
			}
			if direction == 1 {
				frags = append(frags, current)
			} else {
				frags = append([]Fragment{current}, frags...)
			}
		case *TokenFragment, *StringFragment:
			return
		}
	}
	return
}

func (f *Fragger) findNode(from int, direction int) (node ast.Node, dec *DecorationFragment, found bool) {

	var name string
	switch direction {
	case 1:
		name = "Start"
	case -1:
		name = "End"
	}

	for i := from; i < len(f.Fragments) && i >= 0; i += direction {
		switch frag := f.Fragments[i].(type) {
		case *DecorationFragment:
			if frag.Name == name {
				return frag.Node, frag, true
			}
			return
		case *CommentFragment:
			if frag.Attached != nil && frag.Attached.Name == name {
				return frag.Attached.Node, frag.Attached, true
			}
		case *NewlineFragment:
			if frag.Attached != nil && frag.Attached.Name == name {
				return frag.Attached.Node, frag.Attached, true
			}
		case *TokenFragment, *StringFragment:
			return
		}
	}
	return
}

func (f *Fragger) findIndentedComments(from int, indents [2]int) (frags [2][]Fragment, nextDecoration *DecorationFragment) {
	var stage int
	var pastNewline bool // while this is false, we're on the same line that the stmt ended, so we accept all comments regardless of the indent (e.g. empty clauses) - see "hanging-indent-same-line" test case.
	for i := from; i < len(f.Fragments); i++ {
		switch current := f.Fragments[i].(type) {
		case *DecorationFragment:
			return frags, current
		case *NewlineFragment:
			pastNewline = true
			frags[stage] = append(frags[stage], current)
		case *CommentFragment:
			if !pastNewline {
				frags[stage] = append(frags[stage], current)
				continue
			}
			if stage == 0 {
				// Check indent matches. If not, move to second stage or exit if that doesn't match.
				if current.Indent != indents[0] {
					if current.Indent == indents[1] {
						stage = 1
					} else {
						return
					}
				}
			} else if stage == 1 {
				if current.Indent != indents[1] {
					return
				}
			}
			frags[stage] = append(frags[stage], current)
		case *TokenFragment, *StringFragment:
			return
		}
	}
	return
}

type Fragment interface {
	Position() token.Pos
	HasNewline() bool
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
	Text     string
	Pos      token.Pos
	Attached *DecorationFragment // where did we attach this comment in pass 1?
	Indent   int                 // indent if this comment follows a newline
}

type NewlineFragment struct {
	Pos      token.Pos
	Empty    bool                // true if this newline is an empty line (e.g. follows a "//" comment or "\n")
	Attached *DecorationFragment // where did we attach this comment in pass 1?
}

type DecorationFragment struct {
	Node ast.Node
	Name string
	Pos  token.Pos
}

func (v *TokenFragment) Position() token.Pos      { return v.Pos }
func (v *StringFragment) Position() token.Pos     { return v.Pos }
func (v *CommentFragment) Position() token.Pos    { return v.Pos }
func (v *NewlineFragment) Position() token.Pos    { return v.Pos }
func (v *DecorationFragment) Position() token.Pos { return v.Pos }

func (v *TokenFragment) HasNewline() bool      { return false }
func (v *StringFragment) HasNewline() bool     { return false }
func (v *CommentFragment) HasNewline() bool    { return strings.HasPrefix(v.Text, "//") }
func (v *NewlineFragment) HasNewline() bool    { return true }
func (v *DecorationFragment) HasNewline() bool { return false }

func (f Fragger) debug(fset *token.FileSet, w io.Writer) {
	formatPos := func(s token.Position) string {
		return s.String()[strings.Index(s.String(), ":")+1:]
	}
	nodeType := func(n ast.Node) string {
		return strings.Replace(fmt.Sprintf("%T", n), "*ast.", "", -1)
	}
	for _, v := range f.Fragments {
		switch v := v.(type) {
		case *NewlineFragment:
			if v.Empty {
				fmt.Fprintf(w, "Empty line %s\n", formatPos(fset.Position(v.Pos)))
			} else {
				fmt.Fprintf(w, "New line %s\n", formatPos(fset.Position(v.Pos)))
			}
		case *TokenFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.Token, formatPos(fset.Position(v.Pos)))
		case *StringFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.String, formatPos(fset.Position(v.Pos)))
		case *DecorationFragment:
			fmt.Fprintf(w, "%s %s %s\n", nodeType(v.Node), v.Name, formatPos(fset.Position(v.Pos)))
		case *CommentFragment:
			fmt.Fprintf(w, "%q %s\n", v.Text, formatPos(fset.Position(v.Pos)))
		default:
			fmt.Fprintf(w, "%T %s\n", v, formatPos(fset.Position(v.Position())))
		}
	}
}
