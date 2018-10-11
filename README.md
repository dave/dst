# Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree with high fidelity. Decorations (e.g. 
comments and newlines) remain attached to the correct nodes as the tree is modified.

## Where does `go/ast` break?

See [this golang issue](https://github.com/golang/go/issues/20744) for more information.

Consider this example where we want to reverse the order of the two declarations. As you can see the 
comments don't remain attached to the correct nodes:

```go
code := `package a

var a int    // foo
var b string // bar
`
fset := token.NewFileSet()
f, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
if err != nil {
	panic(err)
}

f.Decls = []ast.Decl{f.Decls[1], f.Decls[0]}

if err := format.Node(os.Stdout, fset, f); err != nil {
	panic(err)
}

//Output:
//package a
//
//// foo
//var b string
//var a int
//
//// bar
```

Here's the same example using `dst`:

```go
code := `package a

var a int    // foo
var b string // bar
`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

f.Decls = []dst.Decl{f.Decls[1], f.Decls[0]}

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package a
//
//var b string // bar
//var a int    // foo
```

## Examples

#### Adding comments 

```go
code := `package main

func main() {
	println("Hello World!")
}`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

callExpr := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

callExpr.Decs.Start.Add("// you can add comments at the start...")
callExpr.Decs.Fun.Add("/* ...in the middle... */")
callExpr.Decs.End.Add("// or at the end.")

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

See [generated-decorations.go](https://github.com/dave/dst/blob/master/generated-decorations.go) for a full 
list of decoration attachment points.

#### Adjusting line-spacing

```go
code := `package main

func main() {
	println(a, b, c)
}`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}

callExpr := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)
callExpr.Decs.Lparen.Add("\n")
for _, v := range callExpr.Args {
	v.SetSpace(dst.NewLine)
}
callExpr.Args[0].End().Add("// you can adjust line-spacing")

if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//	println(
//		a, // you can adjust line-spacing
//		b,
//		c,
//	)
//}
```

## Status

This is an experimental package under development, so the API and behaviour is expected to change 
frequently. However I'm now inviting people to try it out and give feedback. 

## Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
