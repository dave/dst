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

{{ "ExampleAstBroken" | example }}

Here's the same example using `dst`:

{{ "ExampleDstFixed" | example }}

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

{{ "ExampleComment" | example }}

#### Line spacing

The `Before` property marks the node as having a line space (new line or empty line) before the node. 
These spaces are rendered before any decorations attached to the `Start` decoration point. The `After`
property is similar but rendered after the node (and after any `End` decorations).

{{ "ExampleSpace" | example }}

#### Common properties

The common decoration properties (`Start`, `End`, `Before` and `After`) occur on all nodes, and can be 
accessed with the `Decorations()` method on the `Node` interface:

{{ "ExampleDecorated" | example }}

#### Newlines as decorations

The `Before` and `After` properties cover the majority of cases, but occasionally a newline needs to 
be rendered inside a node. Simply add a `\n` decoration to accomplish this. 

#### Clone

Re-using an existing node elsewhere in the tree will panic when the tree is restored to `ast`. Instead,
use the `Clone` function to make a deep copy of the node before re-use:

{{ "ExampleClone" | example }}

#### Apply function from astutil

The [dstutil](https://github.com/dave/dst/tree/master/dstutil) package is a fork of `golang.org/x/tools/go/ast/astutil`, 
and provides the `Apply` function with similar semantics.     

#### Imports

The decorator can automatically manage the `import` block, which is a non-trivial task.

Use [NewWithImports](https://godoc.org/github.com/dave/dst/decorator#NewWithImports) and 
[NewRestorerWithImports](https://godoc.org/github.com/dave/dst/decorator#NewRestorerWithImports) to 
create an import aware decorator / restorer with default settings.

During decoration, remote identifiers are normalised: `*ast.SelectorExpr` qualified identifiers are 
replaced with `*dst.Ident` with the `Path` field set to the path of the imported package. 

When adding a qualified identifier node, there is no need to use `SelectorExpr` - just add an 
`Ident` and set `Path` to the imported package path. The restorer will wrap it in a `SelectorExpr` 
where appropriate when converting back to ast, and also update the import block.

The [Load](https://godoc.org/github.com/dave/dst/decorator#Load) convenience function uses 
`go/packages` to load packages and decorate all loaded ast files:

{{ "ExampleImports" | example }}

The default resolvers that enable import management may not be suitable for all environments. If 
more control is needed, custom resolvers can be used for both the `Decorator` and `Restorer`. More 
details and several alternative implementations can be found [here](https://github.com/dave/dst/tree/master/decorator/resolver).

Here's an example of manually supplying alternative resolvers for the decorator and resolver:

{{ "ExampleManualImports" | example }}

#### Mappings between ast and dst nodes

The decorator exposes `Dst.Nodes` and `Ast.Nodes` which map between `ast.Node` and `dst.Node`. This 
enables systems that refer to `ast` nodes (such as `go/types`) to be used:

{{ "ExampleTypes" | example }}

### Status

This is an experimental package under development, but the API is not expected to change much going 
forward. Please try it out and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
