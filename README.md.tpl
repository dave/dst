# Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree with high fidelity. Decorations (e.g. 
comments and newlines) remain attached to the correct nodes as the tree is modified.

## Where does `go/ast` break?

See [this golang issue](https://github.com/golang/go/issues/20744) for more information.

Consider this example where we want to reverse the order of the two declarations. As you can see the 
comments don't remain attached to the correct nodes:

{{ "ExampleAstBroken" | example }}

Here's the same example using `dst`:

{{ "ExampleDstFixed" | example }}

## Examples

#### Adding comments 

{{ "ExampleComment" | example }}

See [generated-decorations.go](https://github.com/dave/dst/blob/master/generated-decorations.go) for a full 
list of decoration attachment points.

#### Adjusting line-spacing

{{ "ExampleSpace" | example }}

## Status

This is an experimental package under development, so the API and behaviour is expected to change 
frequently. However I'm now inviting people to try it out and give feedback. 

## Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
