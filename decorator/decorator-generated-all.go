package decorator

import "github.com/dave/dst"

func getDecorationInfo(n dst.Node) (space dst.SpaceType, info []decorationInfo) {
	switch n := n.(type) {
	case *dst.ArrayType:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		info = append(info, decorationInfo{"AfterLen", n.Decs.AfterLen})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.AssignStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLhs", n.Decs.AfterLhs})
		info = append(info, decorationInfo{"AfterTok", n.Decs.AfterTok})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BadDecl:
		space = n.Decs.Space
	case *dst.BadExpr:
		space = n.Decs.Space
	case *dst.BadStmt:
		space = n.Decs.Space
	case *dst.BasicLit:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BinaryExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterOp", n.Decs.AfterOp})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BlockStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLbrace", n.Decs.AfterLbrace})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.BranchStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterTok", n.Decs.AfterTok})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.CallExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFun", n.Decs.AfterFun})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"AfterArgs", n.Decs.AfterArgs})
		info = append(info, decorationInfo{"AfterEllipsis", n.Decs.AfterEllipsis})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.CaseClause:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterCase", n.Decs.AfterCase})
		info = append(info, decorationInfo{"AfterList", n.Decs.AfterList})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
	case *dst.ChanType:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterBegin", n.Decs.AfterBegin})
		info = append(info, decorationInfo{"AfterArrow", n.Decs.AfterArrow})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.CommClause:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterCase", n.Decs.AfterCase})
		info = append(info, decorationInfo{"AfterComm", n.Decs.AfterComm})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
	case *dst.CompositeLit:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"AfterLbrace", n.Decs.AfterLbrace})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.DeclStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.DeferStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterDefer", n.Decs.AfterDefer})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Ellipsis:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterEllipsis", n.Decs.AfterEllipsis})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.EmptyStmt:
		space = n.Decs.Space
	case *dst.ExprStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Field:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterNames", n.Decs.AfterNames})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FieldList:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterOpening", n.Decs.AfterOpening})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.File:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterPackage", n.Decs.AfterPackage})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
	case *dst.ForStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFor", n.Decs.AfterFor})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterCond", n.Decs.AfterCond})
		info = append(info, decorationInfo{"AfterPost", n.Decs.AfterPost})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FuncDecl:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFunc", n.Decs.AfterFunc})
		info = append(info, decorationInfo{"AfterRecv", n.Decs.AfterRecv})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
		info = append(info, decorationInfo{"AfterParams", n.Decs.AfterParams})
		info = append(info, decorationInfo{"AfterResults", n.Decs.AfterResults})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FuncLit:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.FuncType:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFunc", n.Decs.AfterFunc})
		info = append(info, decorationInfo{"AfterParams", n.Decs.AfterParams})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.GenDecl:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterTok", n.Decs.AfterTok})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.GoStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterGo", n.Decs.AfterGo})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Ident:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.IfStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterIf", n.Decs.AfterIf})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterCond", n.Decs.AfterCond})
		info = append(info, decorationInfo{"AfterElse", n.Decs.AfterElse})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.ImportSpec:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.IncDecStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.IndexExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		info = append(info, decorationInfo{"AfterIndex", n.Decs.AfterIndex})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.InterfaceType:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterInterface", n.Decs.AfterInterface})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.KeyValueExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterKey", n.Decs.AfterKey})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.LabeledStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLabel", n.Decs.AfterLabel})
		info = append(info, decorationInfo{"AfterColon", n.Decs.AfterColon})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.MapType:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterMap", n.Decs.AfterMap})
		info = append(info, decorationInfo{"AfterKey", n.Decs.AfterKey})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.Package:
	case *dst.ParenExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.RangeStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterFor", n.Decs.AfterFor})
		info = append(info, decorationInfo{"AfterKey", n.Decs.AfterKey})
		info = append(info, decorationInfo{"AfterValue", n.Decs.AfterValue})
		info = append(info, decorationInfo{"AfterRange", n.Decs.AfterRange})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.ReturnStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterReturn", n.Decs.AfterReturn})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SelectStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterSelect", n.Decs.AfterSelect})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SelectorExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SendStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterChan", n.Decs.AfterChan})
		info = append(info, decorationInfo{"AfterArrow", n.Decs.AfterArrow})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SliceExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		info = append(info, decorationInfo{"AfterLow", n.Decs.AfterLow})
		info = append(info, decorationInfo{"AfterHigh", n.Decs.AfterHigh})
		info = append(info, decorationInfo{"AfterMax", n.Decs.AfterMax})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.StarExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterStar", n.Decs.AfterStar})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.StructType:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterStruct", n.Decs.AfterStruct})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.SwitchStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterSwitch", n.Decs.AfterSwitch})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterTag", n.Decs.AfterTag})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.TypeAssertExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterX", n.Decs.AfterX})
		info = append(info, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		info = append(info, decorationInfo{"AfterType", n.Decs.AfterType})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.TypeSpec:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterName", n.Decs.AfterName})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.TypeSwitchStmt:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterSwitch", n.Decs.AfterSwitch})
		info = append(info, decorationInfo{"AfterInit", n.Decs.AfterInit})
		info = append(info, decorationInfo{"AfterAssign", n.Decs.AfterAssign})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.UnaryExpr:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterOp", n.Decs.AfterOp})
		info = append(info, decorationInfo{"End", n.Decs.End})
	case *dst.ValueSpec:
		space = n.Decs.Space
		info = append(info, decorationInfo{"Start", n.Decs.Start})
		info = append(info, decorationInfo{"AfterNames", n.Decs.AfterNames})
		info = append(info, decorationInfo{"AfterAssign", n.Decs.AfterAssign})
		info = append(info, decorationInfo{"End", n.Decs.End})
	}
	return
}
