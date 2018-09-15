# dst

### Decorated Syntax Tree

The `dst` package attempts to provide a work-arround for [go/ast: Free-floating comments are 
single-biggest issue when manipulating the AST](https://github.com/golang/go/issues/20744).

### Progress as of 15th September

[github.com/dave/dst](https://github.com/dave/dst) is a fork of the `go/ast` package with a few changes. The [decorator](https://github.com/dave/dst/tree/master/decorator) package converts from `*ast.File + *token.FileSet` to `*dst.File` and back again.

All the position fields have been removed from `dst` so it's just the location in the tree that determines the position of the tokens. Decorations (e.g. comments and newlines) are stored along with each node, and attached to the node at various points. The intention is that any place `gofmt` will allow a comment / new-line to be attached, `dst` will allow this.

I've finished a very rough prototype that works pretty well. (Take a look at [restorer_test.go](https://github.com/dave/dst/blob/master/decorator/restorer_test.go#L11) - all the tests pass apart from `FuncDecl` now).

There's several special cases that it doesn't currently handle. Right now I'm generating much of the code, so the special cases are non-trivial to implement. (e.g. Look at [FuncDecl](https://github.com/golang/go/blob/master/src/go/ast/ast.go#L927-L934) - the `func` token from the `Type` field is rendered before `Recv` and `Name`). Over the next few weeks I'll refactor and handle the special cases.

As @griesemer points out a big problem is where to attach the decorations so as you manipulate the tree they remain attached to the node you were expecting. My algorithm probably needs improvement here too (see [decorator_test.go](https://github.com/dave/dst/blob/master/decorator/decorator_test.go)), but I think it currently works well enough to be useful.

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the [#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
