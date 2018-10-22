# gendst

The `gendst` package is used to create the generated portions of the `dst` and `decorator` packages.
The input data is in the [data](https://github.com/dave/dst/blob/master/gendst/data/data.go) package, 
which was painstakingly manually compiled. The code in [positions.go](https://github.com/dave/dst/blob/master/gendst/data/positions.go)
is sliced up automatically to make the documentation for the [decoration holder classes](https://github.com/dave/dst/blob/master/decorations-types-generated.go).

Everything else is generated from this: 

### dst
* [decorations-node-generated.go](https://github.com/dave/dst/blob/master/decorations-node-generated.go)
* [decorations-types-generated.go](https://github.com/dave/dst/blob/master/decorations-types-generated.go)

### decorator
* [fragger-generated.go](https://github.com/dave/dst/blob/master/decorator/fragger-generated.go)
* [decorator-generated.go](https://github.com/dave/dst/blob/master/decorator/decorator-generated.go)
* [decorator-info-generated.go](https://github.com/dave/dst/blob/master/decorator/decorator-info-generated.go)
* [restorer-generated.go](https://github.com/dave/dst/blob/master/decorator/restorer-generated.go)