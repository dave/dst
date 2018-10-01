package main

import (
	"go/types"

	"sort"

	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func getPkg(path string) *loader.PackageInfo {
	var conf loader.Config
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}
	return prog.Package(path)
}

func run() error {

	astPkg := getPkg("go/ast")
	dstPkg := getPkg(DSTPATH)

	// Find the "Node" interface so we can match later
	astNode := astPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)
	dstNode := dstPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)

	// Order the types so we get reproducible output
	var ordered []string
	for _, name := range astPkg.Pkg.Scope().Names() {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)

	var names []string

	nodes := map[string]NodeInfo{}

	//fields := map[string][]string{}

	for _, typeName := range ordered {

		astFieldNames, astFields, ok := getFields(typeName, astPkg, astNode)
		if !ok {
			continue
		}
		_, dstFields, _ := getFields(typeName, dstPkg, dstNode)

		names = append(names, typeName)

		ni := NodeInfo{
			Name: typeName,
		}

		for _, fieldName := range astFieldNames {

			if !matchBool(fragmentFieldNames, typeName, fieldName) {
				continue
			}

			sn := FragmentInfo{
				Name: fieldName,
			}

			astField := astFields[fieldName]

			sn.AstType = getTypeName(astField.Type(), astNode)
			sn.AstTypeActual, sn.AstTypePointer = getTypeActual(astField.Type())

			dstField, ok := dstFields[fieldName]
			if ok {
				sn.DstType = getTypeName(dstField.Type(), dstNode)
				sn.DstTypeActual, sn.DstTypePointer = getTypeActual(dstField.Type())
			}

			if v, ok := matchString(fragmentPositionFields, typeName, fieldName); ok {
				sn.AstPositionField = v
			} else if sn.AstType == "Pos" {
				sn.AstPositionField = fieldName
			}

			ni.Fragments = append(ni.Fragments, sn)
		}

		if fromToLengthTypes[typeName] {
			ni.FromToLength = true
		}

		if fo, ok := overrideFragmentOrder[typeName]; ok {
			ni.FragmentOrder = fo
		}

		if fis, ok := dataFields[typeName]; ok {
			ni.Data = fis
		}

		nodes[ni.Name] = ni

	}

	//for k, v := range fields {
	//	fmt.Println(k, v)
	//}

	if err := generateDst(names, nodes); err != nil {
		return err
	}
	if err := generateFragger(names); err != nil {
		return err
	}
	if err := generateDecorator(names); err != nil {
		return err
	}
	if err := generateRestorer(names, nodes); err != nil {
		return err
	}

	//if err := generateInfo(names, nodes); err != nil {
	//	return err
	//}

	/*
		// This was used to generate the basis of restorer-length.go, but special cases are added after
		// generation, so we should not run it again.
		if err := generateLength(names, nodes); err != nil {
			return err
		}
	*/
	return nil
}

func getFields(name string, pkg *loader.PackageInfo, node *types.Interface) (fieldNames []string, fields map[string]*types.Var, found bool) {

	if ignoredTypes[name] {
		return
	}

	v := pkg.Pkg.Scope().Lookup(name)
	if v == nil {
		return
	}
	if !v.Exported() {
		return
	}

	tn, ok := v.(*types.TypeName)
	if !ok {
		return
	}
	if !types.Implements(types.NewPointer(tn.Type()), node) {
		return
	}

	found = true

	typ := tn.Type().Underlying().(*types.Struct)
	fields = map[string]*types.Var{}

	for i := 0; i < typ.NumFields(); i++ {
		f := typ.Field(i)
		fieldNames = append(fieldNames, f.Name())
		fields[f.Name()] = f
	}

	return
}

func unwrap(t types.Type) types.Type {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	return t
}

var fixedLength = map[string]int{
	"Defer":          len("defer"),     // [DeferStmt]
	"Interface":      len("interface"), // [InterfaceType]
	"From":           0,                // [BadDecl BadStmt BadExpr]
	"OpPos":          0,                // [BinaryExpr UnaryExpr]
	"For":            len("for"),       // [RangeStmt ForStmt]
	"Return":         len("return"),    // [ReturnStmt]
	"Func":           len("func"),      // [FuncType]
	"Case":           len("case"),      // [CaseClause CommClause]
	"SendStmt.Arrow": len("<-"),        // [SendStmt] (ChanType has custom length)
	"If":             len("if"),        // [IfStmt]
	"ValuePos":       0,                // [BasicLit]
	"Lbrack":         len("["),         // [SliceExpr ArrayType IndexExpr]
	"TokPos":         0,                // [RangeStmt BranchStmt AssignStmt IncDecStmt GenDecl]
	"Struct":         len("struct"),    // [StructType]
	"Slash":          0,                // [Comment]
	"Star":           len("*"),         // [StarExpr]
	"Rbrace":         len("}"),         // [BlockStmt CompositeLit]
	"Assign":         len("="),         // [TypeSpec]
	"Opening":        len("("),         // [FieldList]
	"NamePos":        0,                // [Ident]
	"Rbrack":         len("]"),         // [SliceExpr IndexExpr]
	"Semicolon":      0,                // [EmptyStmt]
	"Ellipsis":       len("..."),       // [Ellipsis CallExpr]
	"Closing":        len(")"),         // [FieldList]
	"Lparen":         len("("),         // [TypeAssertExpr ParenExpr GenDecl CallExpr]
	"Rparen":         len(")"),         // [TypeAssertExpr ParenExpr GenDecl CallExpr]
	"Package":        len("package"),   // [File]
	"Switch":         len("switch"),    // [SwitchStmt TypeSwitchStmt]
	"To":             0,                // [BadDecl BadStmt BadExpr]
	"Colon":          len(":"),         // [LabeledStmt KeyValueExpr CaseClause CommClause]
	"Map":            len("map"),       // [MapType]
	"Lbrace":         len("{"),         // [BlockStmt CompositeLit]
	"EndPos":         0,                // [ImportSpec]
	"Go":             len("go"),        // [GoStmt]
	"Begin":          len("chan"),      // [ChanType]
	"Select":         len("select"),    // [SelectStmt]
}

var ignoredTypes = map[string]bool{
	"Package":      true,
	"Comment":      true,
	"CommentGroup": true,
}

// TODO: EmptyStmt?
var hasCustomLength = map[string]bool{
	"SliceExpr.Max":       true,
	"ChanType.Arrow":      true,
	"TypeAssertExpr.Type": true,
	"CommClause.Comm":     true,
	"ArrayType.Len":       true,
}

var prefixLengths = map[string]int{
	"ForStmt.Post":     len(";"),
	"MapType.Key":      len("["),
	"RangeStmt.X":      len("range"),
	"SliceExpr.High":   len(":"),
	"ValueSpec.Values": len("="),
}

var suffixLengths = map[string]int{
	"ForStmt.Init":        len(";"),
	"IfStmt.Init":         len(";"),
	"IfStmt.Else":         len("else"),
	"MapType.Key":         len("]"),
	"SelectorExpr.X":      len("."),
	"SwitchStmt.Init":     len(";"),
	"TypeSwitchStmt.Init": len(";"),
}

// Field names corresponding to fragments
var fragmentFieldNames = map[string]bool{
	"Args":      true, // [CallExpr([]Node)]
	"Arrow":     true, // [ChanType(Pos) SendStmt(Pos)]
	"Assign":    true, // [TypeSpec(Pos) TypeSwitchStmt(Node)]
	"Begin":     true, // [ChanType(Pos)]
	"Body":      true, // [CaseClause([]Node) CommClause([]Node) ForStmt(Node) FuncDecl(Node) FuncLit(Node) IfStmt(Node) RangeStmt(Node) SelectStmt(Node) SwitchStmt(Node) TypeSwitchStmt(Node)]
	"Call":      true, // [DeferStmt(Node) GoStmt(Node)]
	"Case":      true, // [CaseClause(Pos) CommClause(Pos)]
	"Chan":      true, // [SendStmt(Node)]
	"Closing":   true, // [FieldList(Pos)]
	"Colon":     true, // [CaseClause(Pos) CommClause(Pos) KeyValueExpr(Pos) LabeledStmt(Pos)]
	"Comm":      true, // [CommClause(Node)]
	"Cond":      true, // [ForStmt(Node) IfStmt(Node)]
	"Decl":      true, // [DeclStmt(Node)]
	"Decls":     true, // [File([]Node)]
	"Defer":     true, // [DeferStmt(Pos)]
	"Ellipsis":  true, // [CallExpr(Pos) Ellipsis(Pos)]
	"Else":      true, // [IfStmt(Node)]
	"Elt":       true, // [ArrayType(Node) Ellipsis(Node)]
	"Elts":      true, // [CompositeLit([]Node)]
	"Fields":    true, // [StructType(Node)]
	"For":       true, // [ForStmt(Pos) RangeStmt(Pos)]
	"Fun":       true, // [CallExpr(Node)]
	"Func":      true, // [FuncType(Pos)]
	"Go":        true, // [GoStmt(Pos)]
	"High":      true, // [SliceExpr(Node)]
	"If":        true, // [IfStmt(Pos)]
	"Index":     true, // [IndexExpr(Node)]
	"Init":      true, // [ForStmt(Node) IfStmt(Node) SwitchStmt(Node) TypeSwitchStmt(Node)]
	"Interface": true, // [InterfaceType(Pos)]
	"Key":       true, // [KeyValueExpr(Node) MapType(Node) RangeStmt(Node)]
	"Label":     true, // [BranchStmt(Node) LabeledStmt(Node)]
	"Lbrace":    true, // [BlockStmt(Pos) CompositeLit(Pos)]
	"Lbrack":    true, // [ArrayType(Pos) IndexExpr(Pos) SliceExpr(Pos)]
	"Len":       true, // [ArrayType(Node)]
	"Lhs":       true, // [AssignStmt([]Node)]
	"List":      true, // [BlockStmt([]Node) CaseClause([]Node) CommentGroup([]Node) FieldList([]Node)]
	"Low":       true, // [SliceExpr(Node)]
	"Lparen":    true, // [CallExpr(Pos) GenDecl(Pos) ParenExpr(Pos) TypeAssertExpr(Pos)]
	"Map":       true, // [MapType(Pos)]
	"Max":       true, // [SliceExpr(Node)]
	"Methods":   true, // [InterfaceType(Node)]
	"Name":      true, // [File(Node) FuncDecl(Node) Ident(string) ImportSpec(Node) TypeSpec(Node)]
	"Names":     true, // [Field([]Node) ValueSpec([]Node)]
	"Op":        true, // [BinaryExpr(Token) UnaryExpr(Token)]
	"Opening":   true, // [FieldList(Pos)]
	"Package":   true, // [File(Pos)]
	"Params":    true, // [FuncType(Node)]
	"Path":      true, // [ImportSpec(Node)]
	"Post":      true, // [ForStmt(Node)]
	"Rbrace":    true, // [BlockStmt(Pos) CompositeLit(Pos)]
	"Rbrack":    true, // [IndexExpr(Pos) SliceExpr(Pos)]
	"Recv":      true, // [FuncDecl(Node)]
	"Results":   true, // [FuncType(Node) ReturnStmt([]Node)]
	"Return":    true, // [ReturnStmt(Pos)]
	"Rhs":       true, // [AssignStmt([]Node)]
	"Rparen":    true, // [CallExpr(Pos) GenDecl(Pos) ParenExpr(Pos) TypeAssertExpr(Pos)]
	"Sel":       true, // [SelectorExpr(Node)]
	"Select":    true, // [SelectStmt(Pos)]
	"Semicolon": true, // [EmptyStmt(Pos)]
	"Specs":     true, // [GenDecl([]Node)]
	"Star":      true, // [StarExpr(Pos)]
	"Stmt":      true, // [LabeledStmt(Node)]
	"Struct":    true, // [StructType(Pos)]
	"Switch":    true, // [SwitchStmt(Pos) TypeSwitchStmt(Pos)]
	"Tag":       true, // [Field(Node) SwitchStmt(Node)]
	"Text":      true, // [Comment(string)]
	"Tok":       true, // [AssignStmt(Token) BranchStmt(Token) GenDecl(Token) IncDecStmt(Token) RangeStmt(Token)]
	"Type":      true, // [CompositeLit(Node) Field(Node) FuncDecl(Node) FuncLit(Node) TypeAssertExpr(Node) TypeSpec(Node) ValueSpec(Node)]
	"Value":     true, // [BasicLit(string) ChanType(Node) KeyValueExpr(Node) MapType(Node) RangeStmt(Node) SendStmt(Node)]
	"Values":    true, // [ValueSpec([]Node)]
	"X":         true, // [BinaryExpr(Node) ExprStmt(Node) IncDecStmt(Node) IndexExpr(Node) ParenExpr(Node) RangeStmt(Node) SelectorExpr(Node) SliceExpr(Node) StarExpr(Node) TypeAssertExpr(Node) UnaryExpr(Node)]
	"Y":         true, // [BinaryExpr(Node)]
}

// TODO: ImportSpec.EndPos ???

var fromToLengthTypes = map[string]bool{
	"BadDecl": true,
	"BadExpr": true,
	"BadStmt": true,
}

var nodeStartPositionFieldNames = map[string]string{
	"BadDecl": "From",
	"BadExpr": "From",
	"BadStmt": "From",
}

var nodeEndPositionFieldNames = map[string]string{
	"BadDecl":    "To",
	"BadExpr":    "To",
	"BadStmt":    "To",
	"ImportSpec": "EndPos",
}

// Fields that exist in ast and dst but aren't fragments - we should just copy the values
var dataFields = map[string][]FieldInfo{
	"ChanType": {
		{Name: "Dir", Type: "ChanDir"},
	},
	"EmptyStmt": {
		{Name: "Implicit", Type: "bool"},
	},
	"File": {
		{Name: "Imports", Type: "[]Node", Actual: "ImportSpec", Pointer: true},
		{Name: "Unresolved", Type: "[]Node", Actual: "Ident", Pointer: true},
		{Name: "Scope", Type: "Scope"},
	},
	"CompositeLit": {
		{Name: "Incomplete", Type: "bool"}},
	"InterfaceType": {
		{Name: "Incomplete", Type: "bool"},
	},
	"StructType": {
		{Name: "Incomplete", Type: "bool"},
	},
	"BasicLit": {
		{Name: "Kind", Type: "Token"},
	},
	"Ident": {
		{Name: "Obj", Type: "Object"},
	},
	"SliceExpr": {
		{Name: "Slice3", Type: "bool"},
	},
}

var overrideFragmentOrder = map[string][]string{
	"FuncDecl": {
		"Doc",
		"Type.Func",
		"Recv",
		"Name",
		"Type.Params",
		"Type.Results",
		"Body",
	},
}

/*
var specialNodes = map[string]bool{
	"File":          true,
	"ChanType":      true,
	"EmptyStmt":     true,
	"CompositeLit":  true,
	"InterfaceType": true,
	"StructType":    true,
	"BasicLit":      true,
	"Ident":         true,
	"SliceExpr":     true,
	//TypeAssertExpr
}

var specialCases = map[string]bool{
	"SliceExpr.Max":       true,
	"ChanType.Arrow":      true,
	"EmptyStmt.Semicolon": true,
	"TypeAssertExpr.Type": true,
	"CommClause.Case":     true,
}
*/

var fragmentPositionFields = map[string]string{
	"Ident.Name":     "NamePos",  // [Ident(Pos)]
	"Op":             "OpPos",    // [BinaryExpr(Pos) UnaryExpr(Pos)]
	"Tok":            "TokPos",   // [AssignStmt(Pos) BranchStmt(Pos) GenDecl(Pos) IncDecStmt(Pos) RangeStmt(Pos)]
	"BasicLit.Value": "ValuePos", // [BasicLit(Pos)]
	"Comment.Text":   "Slash",    // [Comment(Pos)]
}

func matchBool(m map[string]bool, node, field string) bool {
	if m[node+"."+field] {
		return true
	}
	if m[field] {
		return true
	}
	return false
}

func matchInt(m map[string]int, node, field string) (int, bool) {
	if v, ok := m[node+"."+field]; ok {
		return v, true
	}
	if v, ok := m[field]; ok {
		return v, true
	}
	return 0, false
}

func matchString(m map[string]string, node, field string) (string, bool) {
	if v, ok := m[node+"."+field]; ok {
		return v, true
	}
	if v, ok := m[field]; ok {
		return v, true
	}
	return "", false
}
