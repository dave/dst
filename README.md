# dst

### Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree in high fidelity. Decorations (e.g. 
comments and newlines) remain attached to the correct nodes as the tree is modified.

See: [golang issue](https://github.com/golang/go/issues/20744).

### How is `go/ast` broken?

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
apply := func(c *astutil.Cursor) bool {
	switch n := c.Node().(type) {
	case *ast.File:
		n.Decls = []ast.Decl{n.Decls[1], n.Decls[0]}
	}
	return true
}
f = astutil.Apply(f, apply, nil).(*ast.File)
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
apply := func(c *dstutil.Cursor) bool {
	switch n := c.Node().(type) {
	case *dst.File:
		n.Decls = []dst.Decl{n.Decls[1], n.Decls[0]}
	}
	return true
}
f = dstutil.Apply(f, apply, nil).(*dst.File)
if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package a
//
//var b string // bar
//var a int    // foo
```

### Example:

This would be very difficult using the `go/ast` package:

```go
code := `package main

func main() {
	var a int
	a++
	print(a)
}`
f, err := decorator.Parse(code)
if err != nil {
	panic(err)
}
apply := func(c *dstutil.Cursor) bool {
	switch n := c.Node().(type) {
	case *dst.DeclStmt:
		n.Decs.End.Replace("// foo")
	case *dst.IncDecStmt:
		n.Decs.AfterX.Add("/* bar */")
	case *dst.CallExpr:
		n.Decs.AfterLparen.Add("\n")
		n.Decs.AfterArgs.Add("\n")
	}
	return true
}
f = dstutil.Apply(f, apply, nil).(*dst.File)
if err := decorator.Print(f); err != nil {
	panic(err)
}

//Output:
//package main
//
//func main() {
//	var a int // foo
//	a /* bar */ ++
//	print(
//		a,
//	)
//}
```

### Status

This is an experimental package under development, so the API and behaviour is expected to change 
frequently. However I'm now inviting people to use it and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
