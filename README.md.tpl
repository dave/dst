# Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree with high fidelity. Decorations (e.g. 
comments and newlines) remain attached to the correct nodes as the tree is modified.

## Where does `go/ast` break?

See [this golang issue](https://github.com/golang/go/issues/20744) for more information.

Consider this example where we want to reverse the order of the two statements. As you can see the 
comments don't remain attached to the correct nodes:

{{ "ExampleAstBroken" | example }}

Here's the same example using `dst`:

{{ "ExampleDstFixed" | example }}

## Examples

#### Line spacing

The `Space` property marks the node as having a line space (new-line or empty-line) before the node. 
These spaces are rendered before any decorations attached to the `Start` decoration point. The `After`
property is similar but rendered after the node (and after any `End` decorations).

{{ "ExampleSpace" | example }}

#### Comments

Comments are added at decoration attachment points. See [generated-decorations.go](https://github.com/dave/dst/blob/master/generated-decorations.go) 
for a full list of these points, along with demonstration code of where they are rendered in the output.

The the decoration points have convenience functions `Add`, `Replace`, `Clear` and `All` to accomplish 
common tasks. 

{{ "ExampleComment" | example }}

#### Common properties

The common decoration properties (`Start`, `End`, `Space` and `After`) occur on all `Expr`, `Stmt` 
and `Decl` nodes, so are available on those interfaces:

{{ "ExampleDecorated" | example }}

#### Newlines as decorations

The `Space` and `After` properties cover the vast majority of cases, but occasionally a newline needs 
to be rendered inside a node. Simply add a `\n` decoration to accomplish this. 

## Status

This is an experimental package under development, so the API and behaviour is expected to change 
frequently. However I'm now inviting people to try it out and give feedback. 

## Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
