package decorator

import "github.com/dave/dst"

func getDecorationInfo(n dst.Node) []decorationInfo {
	var out []decorationInfo
	switch n := n.(type) {
	case *dst.ArrayType:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		out = append(out, decorationInfo{"AfterLen", n.Decs.AfterLen})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.AssignStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterLhs", n.Decs.AfterLhs})
		out = append(out, decorationInfo{"AfterTok", n.Decs.AfterTok})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.BadDecl:
	case *dst.BadExpr:
	case *dst.BadStmt:
	case *dst.BasicLit:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.BinaryExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"AfterOp", n.Decs.AfterOp})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.BlockStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterLbrace", n.Decs.AfterLbrace})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.BranchStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterTok", n.Decs.AfterTok})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.CallExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterFun", n.Decs.AfterFun})
		out = append(out, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		out = append(out, decorationInfo{"AfterArgs", n.Decs.AfterArgs})
		out = append(out, decorationInfo{"AfterEllipsis", n.Decs.AfterEllipsis})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.CaseClause:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterCase", n.Decs.AfterCase})
		out = append(out, decorationInfo{"AfterList", n.Decs.AfterList})
		out = append(out, decorationInfo{"AfterColon", n.Decs.AfterColon})
	case *dst.ChanType:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterBegin", n.Decs.AfterBegin})
		out = append(out, decorationInfo{"AfterArrow", n.Decs.AfterArrow})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.CommClause:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterCase", n.Decs.AfterCase})
		out = append(out, decorationInfo{"AfterComm", n.Decs.AfterComm})
		out = append(out, decorationInfo{"AfterColon", n.Decs.AfterColon})
	case *dst.CompositeLit:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterType", n.Decs.AfterType})
		out = append(out, decorationInfo{"AfterLbrace", n.Decs.AfterLbrace})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.DeclStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.DeferStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterDefer", n.Decs.AfterDefer})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.Ellipsis:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterEllipsis", n.Decs.AfterEllipsis})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.EmptyStmt:
	case *dst.ExprStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.Field:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterNames", n.Decs.AfterNames})
		out = append(out, decorationInfo{"AfterType", n.Decs.AfterType})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.FieldList:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterOpening", n.Decs.AfterOpening})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.File:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterPackage", n.Decs.AfterPackage})
		out = append(out, decorationInfo{"AfterName", n.Decs.AfterName})
	case *dst.ForStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterFor", n.Decs.AfterFor})
		out = append(out, decorationInfo{"AfterInit", n.Decs.AfterInit})
		out = append(out, decorationInfo{"AfterCond", n.Decs.AfterCond})
		out = append(out, decorationInfo{"AfterPost", n.Decs.AfterPost})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.FuncDecl:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterFunc", n.Decs.AfterFunc})
		out = append(out, decorationInfo{"AfterRecv", n.Decs.AfterRecv})
		out = append(out, decorationInfo{"AfterName", n.Decs.AfterName})
		out = append(out, decorationInfo{"AfterParams", n.Decs.AfterParams})
		out = append(out, decorationInfo{"AfterResults", n.Decs.AfterResults})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.FuncLit:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterType", n.Decs.AfterType})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.FuncType:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterFunc", n.Decs.AfterFunc})
		out = append(out, decorationInfo{"AfterParams", n.Decs.AfterParams})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.GenDecl:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterTok", n.Decs.AfterTok})
		out = append(out, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.GoStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterGo", n.Decs.AfterGo})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.Ident:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.IfStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterIf", n.Decs.AfterIf})
		out = append(out, decorationInfo{"AfterInit", n.Decs.AfterInit})
		out = append(out, decorationInfo{"AfterCond", n.Decs.AfterCond})
		out = append(out, decorationInfo{"AfterElse", n.Decs.AfterElse})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.ImportSpec:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterName", n.Decs.AfterName})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.IncDecStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.IndexExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		out = append(out, decorationInfo{"AfterIndex", n.Decs.AfterIndex})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.InterfaceType:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterInterface", n.Decs.AfterInterface})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.KeyValueExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterKey", n.Decs.AfterKey})
		out = append(out, decorationInfo{"AfterColon", n.Decs.AfterColon})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.LabeledStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterLabel", n.Decs.AfterLabel})
		out = append(out, decorationInfo{"AfterColon", n.Decs.AfterColon})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.MapType:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterMap", n.Decs.AfterMap})
		out = append(out, decorationInfo{"AfterKey", n.Decs.AfterKey})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.ParenExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.RangeStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterFor", n.Decs.AfterFor})
		out = append(out, decorationInfo{"AfterKey", n.Decs.AfterKey})
		out = append(out, decorationInfo{"AfterValue", n.Decs.AfterValue})
		out = append(out, decorationInfo{"AfterRange", n.Decs.AfterRange})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.ReturnStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterReturn", n.Decs.AfterReturn})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.SelectStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterSelect", n.Decs.AfterSelect})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.SelectorExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.SendStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterChan", n.Decs.AfterChan})
		out = append(out, decorationInfo{"AfterArrow", n.Decs.AfterArrow})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.SliceExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"AfterLbrack", n.Decs.AfterLbrack})
		out = append(out, decorationInfo{"AfterLow", n.Decs.AfterLow})
		out = append(out, decorationInfo{"AfterHigh", n.Decs.AfterHigh})
		out = append(out, decorationInfo{"AfterMax", n.Decs.AfterMax})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.StarExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterStar", n.Decs.AfterStar})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.StructType:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterStruct", n.Decs.AfterStruct})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.SwitchStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterSwitch", n.Decs.AfterSwitch})
		out = append(out, decorationInfo{"AfterInit", n.Decs.AfterInit})
		out = append(out, decorationInfo{"AfterTag", n.Decs.AfterTag})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.TypeAssertExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterX", n.Decs.AfterX})
		out = append(out, decorationInfo{"AfterLparen", n.Decs.AfterLparen})
		out = append(out, decorationInfo{"AfterType", n.Decs.AfterType})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.TypeSpec:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterName", n.Decs.AfterName})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.TypeSwitchStmt:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterSwitch", n.Decs.AfterSwitch})
		out = append(out, decorationInfo{"AfterInit", n.Decs.AfterInit})
		out = append(out, decorationInfo{"AfterAssign", n.Decs.AfterAssign})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.UnaryExpr:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterOp", n.Decs.AfterOp})
		out = append(out, decorationInfo{"End", n.Decs.End})
	case *dst.ValueSpec:
		out = append(out, decorationInfo{"Start", n.Decs.Start})
		out = append(out, decorationInfo{"AfterNames", n.Decs.AfterNames})
		out = append(out, decorationInfo{"AfterAssign", n.Decs.AfterAssign})
		out = append(out, decorationInfo{"End", n.Decs.End})
	}
	return out
}
