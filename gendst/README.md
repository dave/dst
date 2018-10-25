# gendst

The `gendst` package is used to create the generated portions of the `dst` and `decorator` packages.
The manually compiled input data is in [data.go](https://github.com/dave/dst/blob/master/gendst/data/data.go). 
In addition the code in [positions.go](https://github.com/dave/dst/blob/master/gendst/data/positions.go)
is sliced up automatically to make the documentation for the [decoration holder classes](https://github.com/dave/dst/blob/master/decorations-types-generated.go).

The following files are generated: 

### dst
* [decorations-node-generated.go](https://github.com/dave/dst/blob/master/decorations-node-generated.go)
* [decorations-types-generated.go](https://github.com/dave/dst/blob/master/decorations-types-generated.go)
* [clone-generated.go](https://github.com/dave/dst/blob/master/clone-generated.go)

### decorator
* [decorator-fragment-generated.go](https://github.com/dave/dst/blob/master/decorator/decorator-fragment-generated.go)
* [decorator-node-generated.go](https://github.com/dave/dst/blob/master/decorator/decorator-node-generated.go)
* [decorator-info-generated.go](https://github.com/dave/dst/blob/master/decorator/decorator-info-generated.go)
* [restorer-generated.go](https://github.com/dave/dst/blob/master/decorator/restorer-generated.go)