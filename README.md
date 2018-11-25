[![Build Status](https://travis-ci.org/dave/dst.svg?branch=master)](https://travis-ci.org/dave/dst) [![Documentation](https://img.shields.io/badge/godoc-documentation-brightgreen.svg)](https://godoc.org/github.com/dave/dst/decorator) [![Go Report Card](https://goreportcard.com/badge/github.com/dave/dst)](https://goreportcard.com/report/github.com/dave/dst) <!--[![codecov](https://codecov.io/gh/dave/dst/branch/master/graph/badge.svg)](https://codecov.io/gh/dave/dst)--> ![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg) <a href="https://patreon.com/davebrophy" title="Help with my hosting bills using Patreon"><img src="https://img.shields.io/badge/patreon-donate-yellow.svg" style="max-width:100%;"></a>

# Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree with high fidelity. Decorations (e.g. 
comments and line spacing) remain attached to the correct nodes as the tree is modified.

### Where does `go/ast` break?

The `go/ast` package wasn't created with source manipulation as an intended use-case. Comments are 
stored by their byte offset instead of attached to nodes, so re-arranging nodes breaks the output. 
See [this Go issue](https://github.com/golang/go/issues/20744) for more information.

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
f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
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
to convert back again. 

#### Comments

Comments are added at decoration attachment points. [See here](https://github.com/dave/dst/blob/master/decorations-types-generated.go) 
for a full list of these points, along with demonstration code of where they are rendered in the 
output.

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

The `Before` property marks the node as having a line space (new line or empty line) before the node. 
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

call.Decs.Before = dst.EmptyLine
call.Decs.After = dst.EmptyLine

for _, v := range call.Args {
	v := v.(*dst.Ident)
	v.Decs.Before = dst.NewLine
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

The common decoration properties (`Start`, `End`, `Before` and `After`) occur on all nodes, and can be 
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

list[0].Decorations().Before = dst.EmptyLine
list[0].Decorations().End.Append("// the Decorations method allows access to the common")
list[1].Decorations().End.Append("// decoration properties (Before, Start, End and After)")
list[2].Decorations().End.Append("// for all nodes.")
list[2].Decorations().After = dst.EmptyLine

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//
//	var i int  // the Decorations method allows access to the common
//	i++        // decoration properties (Before, Start, End and After)
//	println(i) // for all nodes.
//
//}
```

#### Newlines as decorations

The `Before` and `After` properties cover the majority of cases, but occasionally a newline needs to 
be rendered inside a node. Simply add a `\n` decoration to accomplish this. 

#### Clone

Re-using an existing node elsewhere in the tree will panic when the tree is restored to `ast`. Instead,
use the `Clone` function to make a deep copy of the node before re-use:

```go
code := `package main

var i /* a */ int`

f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

cloned := dst.Clone(f.Decls[0]).(*dst.GenDecl)

cloned.Decs.Before = dst.NewLine
cloned.Specs[0].(*dst.ValueSpec).Names[0].Name = "j"
cloned.Specs[0].(*dst.ValueSpec).Names[0].Decs.End.Replace("/* b */")

f.Decls = append(f.Decls, cloned)

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//var i /* a */ int
//var j /* b */ int
```

#### Apply function from astutil

The [dstutil](https://github.com/dave/dst/tree/master/dstutil) package is a fork of `golang.org/x/tools/go/ast/astutil`, 
and provides the `Apply` function with similar semantics.     

#### Imports

The decorator can automatically manage the `import` block, which is a non-trivial task.

Use [NewWithImports](https://godoc.org/github.com/dave/dst/decorator#NewWithImports) and 
[NewRestorerWithImports](https://godoc.org/github.com/dave/dst/decorator#NewRestorerWithImports) to 
create an import aware decorator / restorer with recommended settings.

When adding a qualified identifier node, there is no need to use `SelectorExpr` - just add an 
`Ident` and set the [Path](https://godoc.org/github.com/dave/dst#Ident) property to the imported 
package path. The restorer will wrap it in a `SelectorExpr` where appropriate when converting back 
to ast, and also update the import block.

The [Load](https://godoc.org/github.com/dave/dst/decorator#Load) convenience function uses 
`go/packages` to load packages and decorate all loaded ast files:

```go
// Create a simple module in a temporary directory
dir, _ := ioutil.TempDir("", "")
defer os.RemoveAll(dir)
ioutil.WriteFile(filepath.Join(dir, "go.mod"), []byte("module root"), 0666)
ioutil.WriteFile(filepath.Join(dir, "main.go"), []byte("package main \n\n func main() {}"), 0666)

// Use the Load convenience function that calls go/packages to load the package. All loaded
// ast files are decorated to dst.
pkgs, err := decorator.Load(&packages.Config{Dir: dir, Mode: packages.LoadSyntax}, "root")
if err != nil {
	panic(err)
}
p := pkgs[0]
f := p.Files[0]

// Add a call expression. Note we don't have to use a SelectorExpr - just adding an Ident with
// the imported package path will do. The restorer will add SelectorExpr where appropriate when
// converting back to ast. Note the new Path field on *dst.Ident. Set this to the package path
// of the imported package, and the restorer will automatically add the import to the import
// block.
b := f.Decls[0].(*dst.FuncDecl).Body
b.List = append(b.List, &dst.ExprStmt{
	X: &dst.CallExpr{
		Fun:	&dst.Ident{Path: "fmt", Name: "Println"},
		Args: []dst.Expr{
			&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote("Hello, World!")},
		},
	},
})

// Create a restorer with the import manager enabled, and print the result. As you can see, the
// import block is automatically managed, and the Println ident is converted to a SelectorExpr:
r := decorator.NewRestorerWithImports("root", dir)
if err := r.Print(p.Files[0]); err != nil {
	panic(err)
}

//Output:
//package main
//
//import "fmt"
//
//func main() { fmt.Println("Hello, World!") }
```

The default resolvers that enable import management may not be suitable for all environments. If 
more control is needed, custom resolvers can be used for both the `Decorator` and `Restorer`. More 
details and several alternative implementations can be found [here](https://github.com/dave/dst/tree/master/decorator/resolver).

Here's an example of manually supplying alternative resolvers for the decorator and resolver:

```go
code := `package main

	import "fmt"

	func main() {
		fmt.Println("a")
	}`

dec := decorator.New(token.NewFileSet())
dec.Resolver = &goast.IdentResolver{PackageResolver: &guess.PackageResolver{}}

f, err := dec.Parse(code)
if err != nil {
	panic(err)
}

f.Decls[1].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr).Args = []dst.Expr{
	&dst.CallExpr{
		Fun: &dst.Ident{Name: "A", Path: "foo.bar/baz"},
	},
}

res := decorator.NewRestorer()
res.Resolver = &guess.PackageResolver{}
if err := res.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//import (
//	"fmt"
//
//	"foo.bar/baz"
//)
//
//func main() {
//	fmt.Println(baz.A())
//}
```

#### Mappings between ast and dst nodes

The decorator exposes `Dst.Nodes` and `Ast.Nodes` which map between `ast.Node` and `dst.Node`. This 
enables systems that refer to `ast` nodes (such as `go/types`) to be used:

```go
code := `package main

func main() {
	var i int
	i++
	println(i)
}`

// Parse the code to AST
fset := token.NewFileSet()
astFile, err := parser.ParseFile(fset, "", code, parser.ParseComments)
if err != nil {
	panic(err)
}

// Invoke the type checker using AST as input
typesInfo := types.Info{
	Defs:	make(map[*ast.Ident]types.Object),
	Uses:	make(map[*ast.Ident]types.Object),
}
conf := &types.Config{}
if _, err := conf.Check("", fset, []*ast.File{astFile}, &typesInfo); err != nil {
	panic(err)
}

// Create a new decorator, which will track the mapping between ast and dst nodes
dec := decorator.New(fset)

// Decorate the *ast.File to give us a *dst.File
f := dec.DecorateFile(astFile)

// Find the *dst.Ident for the definition of "i"
dstDef := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.DeclStmt).Decl.(*dst.GenDecl).Specs[0].(*dst.ValueSpec).Names[0]

// Find the *ast.Ident using the Ast.Nodes mapping
astDef := dec.Ast.Nodes[dstDef].(*ast.Ident)

// Find the types.Object corresponding to "i"
obj := typesInfo.Defs[astDef]

// Find all the uses of that object
var astUses []*ast.Ident
for id, ob := range typesInfo.Uses {
	if ob != obj {
		continue
	}
	astUses = append(astUses, id)
}

// Find each *dst.Ident in the Dst.Nodes mapping
var dstUses []*dst.Ident
for _, id := range astUses {
	dstUses = append(dstUses, dec.Dst.Nodes[id].(*dst.Ident))
}

// Change the name of the original definition and all uses
dstDef.Name = "foo"
for _, id := range dstUses {
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

### Status

This is an experimental package under development, but the API is not expected to change much going 
forward. Please try it out and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
