// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dst_test

import (
	"fmt"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// This example demonstrates how to inspect the AST of a Go program.
func ExampleInspect() {
	// src is the input for which we want to inspect the AST.
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := decorator.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	dst.Inspect(f, func(n dst.Node) bool {
		var s string
		switch x := n.(type) {
		case *dst.BasicLit:
			s = x.Value
		case *dst.Ident:
			s = x.Name
		}
		if s != "" {
			fmt.Println(s)
		}
		return true
	})

	// Output:
	// p
	// c
	// 1.0
	// X
	// f
	// 3.14
	// 2
	// c
}
