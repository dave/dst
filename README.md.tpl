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

{{ "ExampleAstBroken" | example }}

Here's the same example using `dst`:

{{ "ExampleDstFixed" | example }}

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

{{ "ExampleComment" | example }}

#### Line spacing

The `Space` property marks the node as having a line space (new line or empty line) before the node. 
These spaces are rendered before any decorations attached to the `Start` decoration point. The `After`
property is similar but rendered after the node (and after any `End` decorations).

{{ "ExampleSpace" | example }}

#### Common properties

The common decoration properties (`Start`, `End`, `Space` and `After`) occur on all nodes, and can be 
accessed with the `Decorations()` method on the `Node` interface:

{{ "ExampleDecorated" | example }}

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

{{ "ExampleTypes" | example }}

If you would like to help create a fully `dst` compatible version of `go/types`, feel free to 
continue my work in the [types branch](https://github.com/dave/dst/tree/types).

### Status

This is an experimental package under development, but the API is not expected to change much going 
forward. Please try it out and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
