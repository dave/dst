package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"sort"

	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
func run() error {

	var conf loader.Config

	conf.Import("go/ast")
	prog, err := conf.Load()
	if err != nil {
		return err
	}
	astPkg := prog.Package("go/ast")

	// Find the "Node" interface so we can match later
	node := astPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)
	stmt := astPkg.Pkg.Scope().Lookup("Stmt").Type().Underlying().(*types.Interface)
	expr := astPkg.Pkg.Scope().Lookup("Expr").Type().Underlying().(*types.Interface)
	decl := astPkg.Pkg.Scope().Lookup("Decl").Type().Underlying().(*types.Interface)

	// Order the types so we get reproducible output
	var ordered []*ast.Ident
	for k := range astPkg.Defs {
		ordered = append(ordered, k)
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Name < ordered[j].Name })

	var names []string
	astTypes := map[string]*types.TypeName{}

	nodes := map[string]NodeInfo{}

	for _, k := range ordered {
		v := astPkg.Defs[k]

		typeName := k.Name

		if ignoredTypes[typeName] {
			continue
		}
		if v == nil {
			continue
		}
		if !v.Exported() {
			continue
		}
		tn, ok := v.(*types.TypeName)
		if !ok {
			continue
		}
		if !types.Implements(types.NewPointer(tn.Type()), node) {
			continue
		}

		ts := tn.Type().Underlying().(*types.Struct)

		names = append(names, tn.Name())
		astTypes[tn.Name()] = tn

		ni := NodeInfo{
			Name: tn.Name(),
		}

		for i := 0; i < ts.NumFields(); i++ {

			fieldName := ts.Field(i).Name()
			fullName := fmt.Sprintf("%s.%s", typeName, fieldName)

			if ignoredFields[fullName] || ignoredFields[fieldName] {
				continue
			}

			sn := FragmentInfo{
				Node: &ni,
				Name: fieldName,
			}

			if ph, ok := positionHelpers[fullName]; ok {
				sn.PosField = ph
			} else if ph, ok := positionHelpers[fieldName]; ok {
				sn.PosField = ph
			}

			fieldType := ts.Field(i).Type()
			if st, ok := fieldType.(*types.Slice); ok {
				fieldType = st.Elem()
				sn.Slice = true
			}
			fieldType = unwrap(fieldType)

			switch fieldType := fieldType.(type) {
			case *types.Named:
				sn.Type = fieldType.Obj().Id()
				sn.IsNode = types.Implements(fieldType.Obj().Type(), node) || types.Implements(types.NewPointer(fieldType.Obj().Type()), node)
				sn.IsStmt = types.Implements(fieldType.Obj().Type(), stmt) || types.Implements(types.NewPointer(fieldType.Obj().Type()), stmt)
				sn.IsDecl = types.Implements(fieldType.Obj().Type(), decl) || types.Implements(types.NewPointer(fieldType.Obj().Type()), decl)
				sn.IsExpr = types.Implements(fieldType.Obj().Type(), expr) || types.Implements(types.NewPointer(fieldType.Obj().Type()), expr)
			case *types.Basic:
				switch fieldType.Name() {
				case "bool":
					continue // TODO SliceExpr.Slice3, EmptyStmt.Implicit
				case "string":
					sn.Type = "String"
					sn.LenFieldString = fieldName
				}
			default:
				fmt.Printf("  %s %T ***\n", ts.Field(i).Name(), fieldType)
			}

			switch sn.Type {
			case "Pos":
				if fl, ok := fixedLength[fieldName]; ok {
					sn.HasLength = true
					sn.Length = fl
				}
				sn.PosField = fieldName
			case "Token":
				sn.LenFieldToken = fieldName
			}

			if pl := prefixLengths[fullName]; pl > 0 {
				sn.PrefixLength = pl
			}

			if sl := suffixLengths[fullName]; sl > 0 {
				sn.SuffixLength = sl
			}

			if specialCases[fullName] {
				sn.Special = true
			}

			ni.Fragments = append(ni.Fragments, sn)
		}
		nodes[ni.Name] = ni

	}

	var dstConf loader.Config
	dstConf.AllowErrors = true
	dstConf.TypeChecker.Error = func(err error) {}
	dstConf.Import("github.com/dave/dst")
	dstProg, err := dstConf.Load()
	if err != nil {
		return err
	}
	dstPkg := dstProg.Package("github.com/dave/dst")
	dstTypes := map[string]*types.TypeName{}
	for _, typ := range astTypes {
		dstTypes[typ.Name()] = dstPkg.Pkg.Scope().Lookup(typ.Name()).(*types.TypeName)
	}

	if err := generateProcessor(names, nodes); err != nil {
		return err
	}
	if err := generateDecorator(names, astPkg, astTypes, dstPkg, dstTypes); err != nil {
		return err
	}
	if err := generateDst(names); err != nil {
		return err
	}
	if err := generateRestorer(names, nodes, astPkg, astTypes, dstPkg, dstTypes); err != nil {
		return err
	}
	if err := generateInfo(names, nodes); err != nil {
		return err
	}

	return nil
}

func unwrap(t types.Type) types.Type {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	return t
}

var fixedLength = map[string]int{
	"Defer":     len("defer"),     // [DeferStmt]
	"Interface": len("interface"), // [InterfaceType]
	"From":      0,                // [BadDecl BadStmt BadExpr]
	"OpPos":     0,                // [BinaryExpr UnaryExpr]
	"For":       len("for"),       // [RangeStmt ForStmt]
	"Return":    len("return"),    // [ReturnStmt]
	"Func":      len("func"),      // [FuncType]
	"Case":      len("case"),      // [CaseClause CommClause]
	"Arrow":     2,                // [SendStmt ChanType]
	"If":        len("if"),        // [IfStmt]
	"ValuePos":  0,                // [BasicLit]
	"Lbrack":    1,                // [SliceExpr ArrayType IndexExpr]
	"TokPos":    0,                // [RangeStmt BranchStmt AssignStmt IncDecStmt GenDecl]
	"Struct":    len("struct"),    // [StructType]
	"Slash":     0,                // [Comment]
	"Star":      1,                // [StarExpr]
	"Rbrace":    1,                // [BlockStmt CompositeLit]
	"Assign":    1,                // [TypeSpec]
	"Opening":   1,                // [FieldList]
	"NamePos":   0,                // [Ident]
	"Rbrack":    1,                // [SliceExpr IndexExpr]
	"Semicolon": 0,                // [EmptyStmt]
	"Ellipsis":  3,                // [Ellipsis CallExpr]
	"Closing":   1,                // [FieldList]
	"Lparen":    1,                // [TypeAssertExpr ParenExpr GenDecl CallExpr]
	"Rparen":    1,                // [TypeAssertExpr ParenExpr GenDecl CallExpr]
	"Package":   len("package"),   // [File]
	"Switch":    len("switch"),    // [SwitchStmt TypeSwitchStmt]
	"To":        0,                // [BadDecl BadStmt BadExpr]
	"Colon":     1,                // [LabeledStmt KeyValueExpr CaseClause CommClause]
	"Map":       len("map"),       // [MapType]
	"Lbrace":    1,                // [BlockStmt CompositeLit]
	"EndPos":    0,                // [ImportSpec]
	"Go":        len("go"),        // [GoStmt]
	"Begin":     len("chan"),      //[ChanType]
	"Select":    len("select"),    // [SelectStmt]
}
var ignoredTypes = map[string]bool{
	"Package": true,
	//"BadDecl": true,
	//"BadExpr": true,
	//"BadStmt": true,
}
var ignoredFields = map[string]bool{
	"Obj":               true,
	"Incomplete":        true,
	"Kind":              true,
	"TokPos":            true,
	"ValuePos":          true,
	"OpPos":             true,
	"NamePos":           true,
	"ChanDir":           true,
	"Slash":             true,
	"File.Scope":        true,
	"File.Imports":      true,
	"File.Unresolved":   true,
	"File.Comments":     true,
	"ImportSpec.EndPos": true,
	"ChanType.Dir":      true,
}
var positionHelpers = map[string]string{
	"Tok":            "TokPos",
	"Op":             "OpPos",
	"BasicLit.Value": "ValuePos",
	"Ident.Name":     "NamePos",
	"Comment.Text":   "Slash",
}
var prefixLengths = map[string]int{
	"ForStmt.Post":     1,            // ";"
	"MapType.Key":      1,            // "["
	"RangeStmt.X":      len("range"), // "range"
	"SliceExpr.High":   1,            // ":"
	"SliceExpr.Max":    1,            // ":"
	"ValueSpec.Values": 1,            // "="
}
var suffixLengths = map[string]int{
	"ArrayType.Len":       1,           // "]"
	"ForStmt.Init":        1,           // ";"
	"IfStmt.Init":         1,           // ";"
	"IfStmt.Else":         len("else"), // "else"
	"MapType.Key":         1,           // "]"
	"SelectorExpr.X":      1,           // "."
	"SwitchStmt.Init":     1,           // ";"
	"TypeSwitchStmt.Init": 1,           // ";"
}
var specialCases = map[string]bool{
	"SliceExpr.Max":       true,
	"ChanType.Arrow":      true,
	"EmptyStmt.Semicolon": true,
	"TypeAssertExpr.Type": true,
	"CommClause.Case":     true,
}
