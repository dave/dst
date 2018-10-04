// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements various error reporters.

package types

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"
)

func assert(p bool) {
	if !p {
		panic("assertion failed")
	}
}

func unreachable() {
	panic("unreachable")
}

func (check *Checker) qualifier(pkg *Package) string {
	if pkg != check.pkg {
		return pkg.path
	}
	return ""
}

func (check *Checker) sprintf(format string, args ...interface{}) string {
	for i, arg := range args {
		switch a := arg.(type) {
		case nil:
			arg = "<nil>"
		case operand:
			panic("internal error: should always pass *operand")
		case *operand:
			arg = operandString(a, check.qualifier)
		case token.Pos:
			arg = check.fset.Position(a).String()
		case dst.Expr:
			arg = ExprString(a)
		case Object:
			arg = ObjectString(a, check.qualifier)
		case Type:
			arg = TypeString(a, check.qualifier)
		}
		args[i] = arg
	}
	return fmt.Sprintf(format, args...)
}

func (check *Checker) trace(format string, args ...interface{}) {
	fmt.Printf("%s%s\n",
		strings.Repeat(".  ", check.indent),
		check.sprintf(format, args...),
	)
}

// dump is only needed for debugging
func (check *Checker) dump(format string, args ...interface{}) {
	fmt.Println(check.sprintf(format, args...))
}

func (check *Checker) err(msg string, soft bool) {
	// Cheap trick: Don't report errors with messages containing
	// "invalid operand" or "invalid type" as those tend to be
	// follow-on errors which don't add useful information. Only
	// exclude them if these strings are not at the beginning,
	// and only if we have at least one error already reported.
	if check.firstErr != nil && (strings.Index(msg, "invalid operand") > 0 || strings.Index(msg, "invalid type") > 0) {
		return
	}

	err := Error{msg, soft}
	if check.firstErr == nil {
		check.firstErr = err
	}

	f := check.conf.Error
	if f == nil {
		panic(bailout{}) // report only first error
	}
	f(err)
}

func (check *Checker) error(msg string) {
	check.err(msg, false)
}

func (check *Checker) errorf(format string, args ...interface{}) {
	check.err(check.sprintf(format, args...), false)
}

func (check *Checker) softErrorf(format string, args ...interface{}) {
	check.err(check.sprintf(format, args...), true)
}

func (check *Checker) invalidAST(format string, args ...interface{}) {
	check.errorf("invalid AST: "+format, args...)
}

func (check *Checker) invalidArg(format string, args ...interface{}) {
	check.errorf("invalid argument: "+format, args...)
}

func (check *Checker) invalidOp(format string, args ...interface{}) {
	check.errorf("invalid operation: "+format, args...)
}
