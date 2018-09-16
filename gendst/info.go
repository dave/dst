package main

/*
func generateInfo(names []string, nodes map[string]NodeInfo) error {

	customValues := Options{
		Open:      "{",
		Close:     "}",
		Separator: ",",
		Multi:     true,
	}

	f := NewFile("main")
	f.Comment("This file serves as documentation to explain how the code generation works").Line()
	f.Var().Id("_").Op("=").Map(String()).Id("NodeInfo").Values(DictFunc(func(d Dict) {
		for _, name := range names {
			d[Lit(name)] = Values(Dict{
				Id("Name"): Lit(name),
				Id("Fragments"): Index().Id("FragmentInfo").CustomFunc(customValues, func(g *Group) {
					for _, f := range nodes[name].Fragments {
						g.Values(DictFunc(func(d Dict) {
							v := reflect.ValueOf(f)
							for i := 0; i < v.Type().NumField(); i++ {
								field := v.Type().Field(i)
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
*/
