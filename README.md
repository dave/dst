[![Build Status](https://travis-ci.org/dave/dst.svg?branch=master)](https://travis-ci.org/dave/dst) [![Documentation](https://img.shields.io/badge/godoc-documentation-brightgreen.svg)](https://godoc.org/github.com/dave/dst/decorator) <!--[![Go Report Card](https://goreportcard.com/badge/github.com/dave/dst)](https://goreportcard.com/report/github.com/dave/dst)--> <!--[![codecov](https://codecov.io/gh/dave/dst/branch/master/graph/badge.svg)](https://codecov.io/gh/dave/dst)--> ![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg) <a href="https://patreon.com/davebrophy" title="Help with my hosting bills using Patreon"><img src="https://img.shields.io/badge/patreon-donate-yellow.svg" style="max-width:100%;"></a>

# Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree with high fidelity. Decorations (e.g. 
comments and line spacing) remain attached to the correct nodes as the tree is modified.

### Where does `go/ast` break?

The `go/ast` package wasn't created with source manipulation as an intended use-case. Comments are 
stored by their byte offset instead of attached to nodes. Because of this, re-arranging nodes breaks 
the output. See [this golang issue](https://github.com/golang/go/issues/20744) for more information.

Consider this example where we want to reverse the order of the two statements. As you can see the 
comments don't remain attached to the correct nodes:

```go
code := `package a

func main(){
	var a int    // foo
	var b string // bar
}
`
fset := token.NewFileSet()
f, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
if err != nil {
	panic(err)
}

list := f.Decls[0].(*ast.FuncDecl).Body.List
list[0], list[1] = list[1], list[0]

if err := format.Node(os.Stdout, fset, f); err != nil {
	panic(err)
}

//Output:
//package a
//
//func main() {
//	// foo
//	var b string
//	var a int
//	// bar
//}
```

Here's the same example using `dst`:

```go
code := `package a

func main(){
	var a int    // foo
	var b string // bar
}
`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

list := f.Decls[0].(*dst.FuncDecl).Body.List
list[0], list[1] = list[1], list[0]

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package a
//
//func main() {
//	var b string // bar
//	var a int    // foo
//}
```

### Usage

Parsing a source file to `dst` and printing the results after modification can be accomplished with 
several `Parse` and `Print` convenience functions in the [decorator](https://godoc.org/github.com/dave/dst/decorator) 
package. 

For more fine-grained control you can use [Decorator](https://godoc.org/github.com/dave/dst/decorator#Decorator) 
to convert from `ast` to `dst`, and [Restorer](https://godoc.org/github.com/dave/dst/decorator#Restorer) 
to convert back again. See the `go/types` section below for a demonstration.  

#### Comments

Comments are added at decoration attachment points. See [decorations-types-generated.go](https://github.com/dave/dst/blob/master/decorations-types-generated.go) 
for a full list of these points, along with demonstration code of where they are rendered in the output.

The decoration attachment points have convenience functions `Append`, `Prepend`, `Replace`, `Clear` 
and `All` to accomplish common tasks. Use the full text of your comment including the `//` or `/**/` 
markers. When adding a line comment, a newline is automatically rendered.

```go
code := `package main

func main() {
	println("Hello World!")
}`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

call := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

call.Decs.Start.Append("// you can add comments at the start...")
call.Decs.Fun.Append("/* ...in the middle... */")
call.Decs.End.Append("// or at the end.")

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//	// you can add comments at the start...
//	println /* ...in the middle... */ ("Hello World!") // or at the end.
//}
```

#### Line spacing

The `Space` property marks the node as having a line space (new line or empty line) before the node. 
These spaces are rendered before any decorations attached to the `Start` decoration point. The `After`
property is similar but rendered after the node (and after any `End` decorations).

```go
code := `package main

func main() {
	println(a, b, c)
}`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

call := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

call.Decs.Space = dst.EmptyLine
call.Decs.After = dst.EmptyLine

for _, v := range call.Args {
	v := v.(*dst.Ident)
	v.Decs.Space = dst.NewLine
	v.Decs.After = dst.NewLine
}

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//
//	println(
//		a,
//		b,
//		c,
//	)
//
//}
```

#### Common properties

The common decoration properties (`Start`, `End`, `Space` and `After`) occur on all nodes, and can be 
accessed with the `Decorations()` method on the `Node` interface:

```go
code := `package main

func main() {
	var i int
	i++
	println(i)
}`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

list := f.Decls[0].(*dst.FuncDecl).Body.List

list[0].Decorations().Space = dst.EmptyLine
list[0].Decorations().End.Append("// the Decorated interface allows access to the common")
list[1].Decorations().End.Append("// decoration properties (Space, Start, End and After)")
list[2].Decorations().End.Append("// for all Expr, Stmt and Decl nodes.")
list[2].Decorations().After = dst.EmptyLine

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//
//	var i int  // the Decorated interface allows access to the common
//	i++        // decoration properties (Space, Start, End and After)
//	println(i) // for all Expr, Stmt and Decl nodes.
//
//}
```

#### Newlines as decorations

The `Space` and `After` properties cover the majority of cases, but occasionally a newline needs to 
be rendered inside a node. Simply add a `\n` decoration to accomplish this. 

#### Apply function from astutil

The [dstutil](https://github.com/dave/dst/tree/master/dstutil) package is a fork of `golang.org/x/tools/go/ast/astutil`, 
and provides the `Apply` function with similar semantics.

#### Integrating with go/types

Forking the `go/types` package to use a `dst` tree as input is non-trivial because `go/types` uses 
position information in several places. A work-around is to convert `ast` to `dst` using a 
[Decorator](https://github.com/dave/dst/blob/master/decorator/decorator.go). After conversion, this 
exposes the `DstNodes` and `AstNodes` properties which map between `ast.Node` and `dst.Node`. This 
way the `go/types` package can be used:

```go
code := `package main

func main() {
	var i int
	i++
	println(i)
}`

// Parse the code to AST
fset := token.NewFileSet()
astFile, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
if err != nil {
	panic(err)
}

// Invoke the type checker using AST as input
typesInfo := types.Info{
	Defs:	make(map[*ast.Ident]types.Object),
	Uses:	make(map[*ast.Ident]types.Object),
}
conf := &types.Config{}
if _, err := conf.Check("a", fset, []*ast.File{astFile}, &typesInfo); err != nil {
	panic(err)
}

// Decorate the *ast.File to give us a *dst.File
dec := decorator.New()
f := dec.Decorate(fset, astFile).(*dst.File)

// Find the *dst.Ident for the definition of "i"
dstDef := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.DeclStmt).Decl.(*dst.GenDecl).Specs[0].(*dst.ValueSpec).Names[0]

// Find the *ast.Ident using the AstNodes mapping
astDef := dec.AstNodes[dstDef].(*ast.Ident)

// Find the types.Object corresponding to "i"
obj := typesInfo.Defs[astDef]

// Find all the uses of that object
var uses []*dst.Ident
for id, ob := range typesInfo.Uses {
	if ob != obj {
		continue
	}
	uses = append(uses, dec.DstNodes[id].(*dst.Ident))
}

// Change the name of all uses
dstDef.Name = "foo"
for _, id := range uses {
	id.Name = "foo"
}

// Print the DST
if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//	var foo int
//	foo++
//	println(foo)
//}
```

If you would like to help create a fully `dst` compatible version of `go/types`, feel free to 
continue my work in the [types branch](https://github.com/dave/dst/tree/types).

### Status

This is an experimental package under development, but the API is not expected to change much going 
forward. Please try it out and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
