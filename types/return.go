// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements isTerminating.

package types

import (
	"go/token"

	"github.com/dave/dst"
)

// isTerminating reports if s is a terminating statement.
// If s is labeled, label is the label name; otherwise s
// is "".
func (check *Checker) isTerminating(s dst.Stmt, label string) bool {
	switch s := s.(type) {
	default:
		unreachable()

	case *dst.BadStmt, *dst.DeclStmt, *dst.EmptyStmt, *dst.SendStmt,
		*dst.IncDecStmt, *dst.AssignStmt, *dst.GoStmt, *dst.DeferStmt,
		*dst.RangeStmt:
		// no chance

	case *dst.LabeledStmt:
		return check.isTerminating(s.Stmt, s.Label.Name)

	case *dst.ExprStmt:
		// calling the predeclared (possibly parenthesized) panic() function is terminating
		if call, ok := unparen(s.X).(*dst.CallExpr); ok && check.isPanic[call] {
			return true
		}

	case *dst.ReturnStmt:
		return true

	case *dst.BranchStmt:
		if s.Tok == token.GOTO || s.Tok == token.FALLTHROUGH {
			return true
		}

	case *dst.BlockStmt:
		return check.isTerminatingList(s.List, "")

	case *dst.IfStmt:
		if s.Else != nil &&
			check.isTerminating(s.Body, "") &&
			check.isTerminating(s.Else, "") {
			return true
		}

	case *dst.SwitchStmt:
		return check.isTerminatingSwitch(s.Body, label)

	case *dst.TypeSwitchStmt:
		return check.isTerminatingSwitch(s.Body, label)

	case *dst.SelectStmt:
		for _, s := range s.Body.List {
			cc := s.(*dst.CommClause)
			if !check.isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
				return false
			}

		}
		return true

	case *dst.ForStmt:
		if s.Cond == nil && !hasBreak(s.Body, label, true) {
			return true
		}
	}

	return false
}

func (check *Checker) isTerminatingList(list []dst.Stmt, label string) bool {
	// trailing empty statements are permitted - skip them
	for i := len(list) - 1; i >= 0; i-- {
		if _, ok := list[i].(*dst.EmptyStmt); !ok {
			return check.isTerminating(list[i], label)
		}
	}
	return false // all statements are empty
}

func (check *Checker) isTerminatingSwitch(body *dst.BlockStmt, label string) bool {
	hasDefault := false
	for _, s := range body.List {
		cc := s.(*dst.CaseClause)
		if cc.List == nil {
			hasDefault = true
		}
		if !check.isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
			return false
		}
	}
	return hasDefault
}

// TODO(gri) For nested breakable statements, the current implementation of hasBreak
//	     will traverse the same subtree repeatedly, once for each label. Replace
//           with a single-pass label/break matching phase.

// hasBreak reports if s is or contains a break statement
// referring to the label-ed statement or implicit-ly the
// closest outer breakable statement.
func hasBreak(s dst.Stmt, label string, implicit bool) bool {
	switch s := s.(type) {
	default:
		unreachable()

	case *dst.BadStmt, *dst.DeclStmt, *dst.EmptyStmt, *dst.ExprStmt,
		*dst.SendStmt, *dst.IncDecStmt, *dst.AssignStmt, *dst.GoStmt,
		*dst.DeferStmt, *dst.ReturnStmt:
		// no chance

	case *dst.LabeledStmt:
		return hasBreak(s.Stmt, label, implicit)

	case *dst.BranchStmt:
		if s.Tok == token.BREAK {
			if s.Label == nil {
				return implicit
			}
			if s.Label.Name == label {
				return true
			}
		}

	case *dst.BlockStmt:
		return hasBreakList(s.List, label, implicit)

	case *dst.IfStmt:
		if hasBreak(s.Body, label, implicit) ||
			s.Else != nil && hasBreak(s.Else, label, implicit) {
			return true
		}

	case *dst.CaseClause:
		return hasBreakList(s.Body, label, implicit)

	case *dst.SwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *dst.TypeSwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *dst.CommClause:
		return hasBreakList(s.Body, label, implicit)

	case *dst.SelectStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *dst.ForStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *dst.RangeStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}
	}

	return false
}

func hasBreakList(list []dst.Stmt, label string, implicit bool) bool {
	for _, s := range list {
		if hasBreak(s, label, implicit) {
			return true
		}
	}
	return false
}
