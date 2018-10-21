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

// This example shows what an AST looks like when printed for debugging.
func ExamplePrint() {
	// src is the input for which we want to print the AST.
	src := `
package main
func main() {
	println("Hello, World!")
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := decorator.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	dst.Print(f)

	//Output:
	//      0  *dst.File {
	//      1  .  Name: *dst.Ident {
	//      2  .  .  Name: "main"
	//      3  .  .  Decs: dst.IdentDecorations {
	//      4  .  .  .  NodeDecs: dst.NodeDecs {
	//      5  .  .  .  .  Space: None
	//      6  .  .  .  .  After: None
	//      7  .  .  .  }
	//      8  .  .  }
	//      9  .  }
	//     10  .  Decls: []dst.Decl (len = 1) {
	//     11  .  .  0: *dst.FuncDecl {
	//     12  .  .  .  Name: *dst.Ident {
	//     13  .  .  .  .  Name: "main"
	//     14  .  .  .  .  Obj: *dst.Object {
	//     15  .  .  .  .  .  Kind: func
	//     16  .  .  .  .  .  Name: "main"
	//     17  .  .  .  .  .  Decl: *(obj @ 11)
	//     18  .  .  .  .  }
	//     19  .  .  .  .  Decs: dst.IdentDecorations {
	//     20  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     21  .  .  .  .  .  .  Space: None
	//     22  .  .  .  .  .  .  After: None
	//     23  .  .  .  .  .  }
	//     24  .  .  .  .  }
	//     25  .  .  .  }
	//     26  .  .  .  Type: *dst.FuncType {
	//     27  .  .  .  .  Func: true
	//     28  .  .  .  .  Params: *dst.FieldList {
	//     29  .  .  .  .  .  Opening: true
	//     30  .  .  .  .  .  Closing: true
	//     31  .  .  .  .  .  Decs: dst.FieldListDecorations {
	//     32  .  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     33  .  .  .  .  .  .  .  Space: None
	//     34  .  .  .  .  .  .  .  After: None
	//     35  .  .  .  .  .  .  }
	//     36  .  .  .  .  .  }
	//     37  .  .  .  .  }
	//     38  .  .  .  .  Decs: dst.FuncTypeDecorations {
	//     39  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     40  .  .  .  .  .  .  Space: None
	//     41  .  .  .  .  .  .  After: None
	//     42  .  .  .  .  .  }
	//     43  .  .  .  .  }
	//     44  .  .  .  }
	//     45  .  .  .  Body: *dst.BlockStmt {
	//     46  .  .  .  .  List: []dst.Stmt (len = 1) {
	//     47  .  .  .  .  .  0: *dst.ExprStmt {
	//     48  .  .  .  .  .  .  X: *dst.CallExpr {
	//     49  .  .  .  .  .  .  .  Fun: *dst.Ident {
	//     50  .  .  .  .  .  .  .  .  Name: "println"
	//     51  .  .  .  .  .  .  .  .  Decs: dst.IdentDecorations {
	//     52  .  .  .  .  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     53  .  .  .  .  .  .  .  .  .  .  Space: None
	//     54  .  .  .  .  .  .  .  .  .  .  After: None
	//     55  .  .  .  .  .  .  .  .  .  }
	//     56  .  .  .  .  .  .  .  .  }
	//     57  .  .  .  .  .  .  .  }
	//     58  .  .  .  .  .  .  .  Args: []dst.Expr (len = 1) {
	//     59  .  .  .  .  .  .  .  .  0: *dst.BasicLit {
	//     60  .  .  .  .  .  .  .  .  .  Kind: STRING
	//     61  .  .  .  .  .  .  .  .  .  Value: "\"Hello, World!\""
	//     62  .  .  .  .  .  .  .  .  .  Decs: dst.BasicLitDecorations {
	//     63  .  .  .  .  .  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     64  .  .  .  .  .  .  .  .  .  .  .  Space: None
	//     65  .  .  .  .  .  .  .  .  .  .  .  After: None
	//     66  .  .  .  .  .  .  .  .  .  .  }
	//     67  .  .  .  .  .  .  .  .  .  }
	//     68  .  .  .  .  .  .  .  .  }
	//     69  .  .  .  .  .  .  .  }
	//     70  .  .  .  .  .  .  .  Ellipsis: false
	//     71  .  .  .  .  .  .  .  Decs: dst.CallExprDecorations {
	//     72  .  .  .  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     73  .  .  .  .  .  .  .  .  .  Space: None
	//     74  .  .  .  .  .  .  .  .  .  After: None
	//     75  .  .  .  .  .  .  .  .  }
	//     76  .  .  .  .  .  .  .  }
	//     77  .  .  .  .  .  .  }
	//     78  .  .  .  .  .  .  Decs: dst.ExprStmtDecorations {
	//     79  .  .  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     80  .  .  .  .  .  .  .  .  Space: NewLine
	//     81  .  .  .  .  .  .  .  .  After: NewLine
	//     82  .  .  .  .  .  .  .  }
	//     83  .  .  .  .  .  .  }
	//     84  .  .  .  .  .  }
	//     85  .  .  .  .  }
	//     86  .  .  .  .  Decs: dst.BlockStmtDecorations {
	//     87  .  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     88  .  .  .  .  .  .  Space: None
	//     89  .  .  .  .  .  .  After: None
	//     90  .  .  .  .  .  }
	//     91  .  .  .  .  }
	//     92  .  .  .  }
	//     93  .  .  .  Decs: dst.FuncDeclDecorations {
	//     94  .  .  .  .  NodeDecs: dst.NodeDecs {
	//     95  .  .  .  .  .  Space: NewLine
	//     96  .  .  .  .  .  After: None
	//     97  .  .  .  .  }
	//     98  .  .  .  }
	//     99  .  .  }
	//    100  .  }
	//    101  .  Scope: *dst.Scope {
	//    102  .  .  Objects: map[string]*dst.Object (len = 1) {
	//    103  .  .  .  "main": *(obj @ 14)
	//    104  .  .  }
	//    105  .  }
	//    106  .  Decs: dst.FileDecorations {
	//    107  .  .  NodeDecs: dst.NodeDecs {
	//    108  .  .  .  Space: NewLine
	//    109  .  .  .  After: None
	//    110  .  .  }
	//    111  .  }
	//    112  }
}
