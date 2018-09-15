package main

import (
	"reflect"

	"bytes"
	"fmt"

	"io/ioutil"

	. "github.com/dave/jennifer/jen"
)

func generateInfo(typeNames []string, nodeInfos map[string]NodeInfo) error {

	customValues := Options{
		Open:      "{",
		Close:     "}",
		Separator: ",",
		Multi:     true,
	}

	f := NewFile("main")
	f.Comment("This file serves as documentation to explain how the code generation works").Line()
	f.Var().Id("_").Op("=").Map(String()).Id("NodeInfo").Values(DictFunc(func(d Dict) {
		for _, name := range typeNames {
			d[Lit(name)] = Values(Dict{
				Id("Name"): Lit(name),
				Id("Fragments"): Index().Id("FragmentInfo").CustomFunc(customValues, func(g *Group) {
					for _, f := range nodeInfos[name].Fragments {
						g.Values(DictFunc(func(d Dict) {
							v := reflect.ValueOf(f)
							for i := 0; i < v.Type().NumField(); i++ {
								field := v.Type().Field(i)
								if field.Name == "Node" {
									continue
								}
								if v.Field(i).Interface() != reflect.Zero(field.Type).Interface() {
									d[Id(field.Name)] = Lit(v.Field(i).Interface())
								}
							}
						}))
					}
				}),
			})
		}
	}))
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return err
	}
	md := fmt.Sprintf("# FragmentInfo\n\n```go\n%s\n```", buf.String())
	return ioutil.WriteFile("./gendst/README.md", []byte(md), 0666)
}
