package decorator

import "github.com/dave/dst"

func getDecorationInfo(n dst.Node) (before, after dst.SpaceType, info []decorationInfo) {
	switch n := n.(type) {
	case *dst.ArrayType:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		info = append(info, decorationInfo{"AfterLen", n.Decs.AfterLen})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.AssignStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLhs", n.Decs.AfterLhs})
		info = append(info, decorationInfo{"AfterTok", n.Decs.AfterTok})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BadDecl:
		before = n.Decs.Before
		after = n.Decs.After
	case *dst.BadExpr:
		before = n.Decs.Before
		after = n.Decs.After
	case *dst.BadStmt:
		before = n.Decs.Before
		after = n.Decs.After
	case *dst.BasicLit:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BinaryExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterOp", n.Decs.AfterOp})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BlockStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLbrace", n.Decs.AfterLbrace})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BranchStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterTok", n.Decs.AfterTok})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.CallExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFun", n.Decs.AfterFun})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"AfterArgs", n.Decs.AfterArgs})
		info = append(info, decorationInfo{"AfterEllipsis", n.Decs.AfterEllipsis})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.CaseClause:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterCase", n.Decs.AfterCase})
		info = append(info, decorationInfo{"AfterList", n.Decs.AfterList})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
	case *dst.ChanType:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterBegin", n.Decs.AfterBegin})
		info = append(info, decorationInfo{"AfterArrow", n.Decs.AfterArrow})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.CommClause:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterCase", n.Decs.AfterCase})
		info = append(info, decorationInfo{"AfterComm", n.Decs.AfterComm})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
	case *dst.CompositeLit:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"AfterLbrace", n.Decs.AfterLbrace})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.DeclStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.DeferStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterDefer", n.Decs.AfterDefer})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Ellipsis:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterEllipsis", n.Decs.AfterEllipsis})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.EmptyStmt:
		before = n.Decs.Before
		after = n.Decs.After
	case *dst.ExprStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Field:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterNames", n.Decs.AfterNames})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FieldList:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterOpening", n.Decs.AfterOpening})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.File:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterPackage", n.Decs.AfterPackage})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
	case *dst.ForStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFor", n.Decs.AfterFor})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterCond", n.Decs.AfterCond})
		info = append(info, decorationInfo{"AfterPost", n.Decs.AfterPost})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FuncDecl:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFunc", n.Decs.AfterFunc})
		info = append(info, decorationInfo{"AfterRecv", n.Decs.AfterRecv})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
		info = append(info, decorationInfo{"AfterParams", n.Decs.AfterParams})
		info = append(info, decorationInfo{"AfterResults", n.Decs.AfterResults})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FuncLit:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FuncType:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFunc", n.Decs.AfterFunc})
		info = append(info, decorationInfo{"AfterParams", n.Decs.AfterParams})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.GenDecl:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterTok", n.Decs.AfterTok})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.GoStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterGo", n.Decs.AfterGo})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Ident:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.IfStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterIf", n.Decs.AfterIf})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterCond", n.Decs.AfterCond})
		info = append(info, decorationInfo{"AfterElse", n.Decs.AfterElse})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.ImportSpec:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.IncDecStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.IndexExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		info = append(info, decorationInfo{"AfterIndex", n.Decs.AfterIndex})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.InterfaceType:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterInterface", n.Decs.AfterInterface})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.KeyValueExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterKey", n.Decs.AfterKey})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.LabeledStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLabel", n.Decs.AfterLabel})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.MapType:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterMap", n.Decs.AfterMap})
		info = append(info, decorationInfo{"AfterKey", n.Decs.AfterKey})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Package:
	case *dst.ParenExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.RangeStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFor", n.Decs.AfterFor})
		info = append(info, decorationInfo{"AfterKey", n.Decs.AfterKey})
		info = append(info, decorationInfo{"AfterValue", n.Decs.AfterValue})
		info = append(info, decorationInfo{"AfterRange", n.Decs.AfterRange})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.ReturnStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterReturn", n.Decs.AfterReturn})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SelectStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterSelect", n.Decs.AfterSelect})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SelectorExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SendStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterChan", n.Decs.AfterChan})
		info = append(info, decorationInfo{"AfterArrow", n.Decs.AfterArrow})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SliceExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		info = append(info, decorationInfo{"AfterLow", n.Decs.AfterLow})
		info = append(info, decorationInfo{"AfterHigh", n.Decs.AfterHigh})
		info = append(info, decorationInfo{"AfterMax", n.Decs.AfterMax})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.StarExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterStar", n.Decs.AfterStar})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.StructType:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterStruct", n.Decs.AfterStruct})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SwitchStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterSwitch", n.Decs.AfterSwitch})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterTag", n.Decs.AfterTag})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.TypeAssertExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.TypeSpec:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.TypeSwitchStmt:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterSwitch", n.Decs.AfterSwitch})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterAssign", n.Decs.AfterAssign})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.UnaryExpr:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterOp", n.Decs.AfterOp})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.ValueSpec:
		before = n.Decs.Before
		after = n.Decs.After
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterNames", n.Decs.AfterNames})
		info = append(info, decorationInfo{"AfterAssign", n.Decs.AfterAssign})
		info = append(info, decorationInfo{"End", n.Decs.End})
	}
	return
}
