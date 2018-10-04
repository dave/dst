// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"github.com/dave/dst"
)

// labels checks correct label use in body.
func (check *Checker) labels(body *dst.BlockStmt) {
	// set of all labels in this body
	all := NewScope(nil, "label")

	// TODO: deleted code to detect forward jumps (this needs position information)

	// spec: "It is illegal to define a label that is never used."
	for _, obj := range all.elems {
		if lbl := obj.(*Label); !lbl.used {
			check.softErrorf("label %s declared but not used", lbl.name)
		}
	}
}

// A block tracks label declarations in a block and its enclosing blocks.
type block struct {
	parent *block                      // enclosing block
	lstmt  *dst.LabeledStmt            // labeled statement to which this block belongs, or nil
	labels map[string]*dst.LabeledStmt // allocated lazily
}

// insert records a new label declaration for the current block.
// The label must not have been declared before in any block.
func (b *block) insert(s *dst.LabeledStmt) {
	name := s.Label.Name
	if debug {
		assert(b.gotoTarget(name) == nil)
	}
	labels := b.labels
	if labels == nil {
		labels = make(map[string]*dst.LabeledStmt)
		b.labels = labels
	}
	labels[name] = s
}

// gotoTarget returns the labeled statement in the current
// or an enclosing block with the given label name, or nil.
func (b *block) gotoTarget(name string) *dst.LabeledStmt {
	for s := b; s != nil; s = s.parent {
		if t := s.labels[name]; t != nil {
			return t
		}
	}
	return nil
}

// enclosingTarget returns the innermost enclosing labeled
// statement with the given label name, or nil.
func (b *block) enclosingTarget(name string) *dst.LabeledStmt {
	for s := b; s != nil; s = s.parent {
		if t := s.lstmt; t != nil && t.Label.Name == name {
			return t
		}
	}
	return nil
}
