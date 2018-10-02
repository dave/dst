package dst

func (d *ArrayTypeDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterLbrack":
		return &d.AfterLbrack
	case "AfterLen":
		return &d.AfterLen
	case "End":
		return &d.End
	}
	return nil
}
func (d *AssignStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterLhs":
		return &d.AfterLhs
	case "AfterTok":
		return &d.AfterTok
	case "End":
		return &d.End
	}
	return nil
}
func (d *BadDeclDecorations) Get(name string) *Decorations {
	switch name {
	}
	return nil
}
func (d *BadExprDecorations) Get(name string) *Decorations {
	switch name {
	}
	return nil
}
func (d *BadStmtDecorations) Get(name string) *Decorations {
	switch name {
	}
	return nil
}
func (d *BasicLitDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "End":
		return &d.End
	}
	return nil
}
func (d *BinaryExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterX":
		return &d.AfterX
	case "AfterOp":
		return &d.AfterOp
	case "End":
		return &d.End
	}
	return nil
}
func (d *BlockStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterLbrace":
		return &d.AfterLbrace
	case "End":
		return &d.End
	}
	return nil
}
func (d *BranchStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterTok":
		return &d.AfterTok
	case "End":
		return &d.End
	}
	return nil
}
func (d *CallExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterFun":
		return &d.AfterFun
	case "AfterLparen":
		return &d.AfterLparen
	case "AfterArgs":
		return &d.AfterArgs
	case "AfterEllipsis":
		return &d.AfterEllipsis
	case "End":
		return &d.End
	}
	return nil
}
func (d *CaseClauseDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterCase":
		return &d.AfterCase
	case "AfterList":
		return &d.AfterList
	case "AfterColon":
		return &d.AfterColon
	}
	return nil
}
func (d *ChanTypeDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterBegin":
		return &d.AfterBegin
	case "AfterArrow":
		return &d.AfterArrow
	case "End":
		return &d.End
	}
	return nil
}
func (d *CommClauseDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterCase":
		return &d.AfterCase
	case "AfterComm":
		return &d.AfterComm
	case "AfterColon":
		return &d.AfterColon
	}
	return nil
}
func (d *CompositeLitDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterType":
		return &d.AfterType
	case "AfterLbrace":
		return &d.AfterLbrace
	case "End":
		return &d.End
	}
	return nil
}
func (d *DeclStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "End":
		return &d.End
	}
	return nil
}
func (d *DeferStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterDefer":
		return &d.AfterDefer
	case "End":
		return &d.End
	}
	return nil
}
func (d *EllipsisDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterEllipsis":
		return &d.AfterEllipsis
	case "End":
		return &d.End
	}
	return nil
}
func (d *EmptyStmtDecorations) Get(name string) *Decorations {
	switch name {
	}
	return nil
}
func (d *ExprStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "End":
		return &d.End
	}
	return nil
}
func (d *FieldDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterNames":
		return &d.AfterNames
	case "AfterType":
		return &d.AfterType
	case "End":
		return &d.End
	}
	return nil
}
func (d *FieldListDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterOpening":
		return &d.AfterOpening
	case "End":
		return &d.End
	}
	return nil
}
func (d *FileDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterPackage":
		return &d.AfterPackage
	case "AfterName":
		return &d.AfterName
	}
	return nil
}
func (d *ForStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterFor":
		return &d.AfterFor
	case "AfterInit":
		return &d.AfterInit
	case "AfterCond":
		return &d.AfterCond
	case "AfterPost":
		return &d.AfterPost
	case "End":
		return &d.End
	}
	return nil
}
func (d *FuncDeclDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterFunc":
		return &d.AfterFunc
	case "AfterRecv":
		return &d.AfterRecv
	case "AfterName":
		return &d.AfterName
	case "AfterParams":
		return &d.AfterParams
	case "AfterResults":
		return &d.AfterResults
	case "End":
		return &d.End
	}
	return nil
}
func (d *FuncLitDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterType":
		return &d.AfterType
	case "End":
		return &d.End
	}
	return nil
}
func (d *FuncTypeDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterFunc":
		return &d.AfterFunc
	case "AfterParams":
		return &d.AfterParams
	case "End":
		return &d.End
	}
	return nil
}
func (d *GenDeclDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterTok":
		return &d.AfterTok
	case "AfterLparen":
		return &d.AfterLparen
	case "End":
		return &d.End
	}
	return nil
}
func (d *GoStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterGo":
		return &d.AfterGo
	case "End":
		return &d.End
	}
	return nil
}
func (d *IdentDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "End":
		return &d.End
	}
	return nil
}
func (d *IfStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterIf":
		return &d.AfterIf
	case "AfterInit":
		return &d.AfterInit
	case "AfterCond":
		return &d.AfterCond
	case "AfterElse":
		return &d.AfterElse
	case "End":
		return &d.End
	}
	return nil
}
func (d *ImportSpecDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterName":
		return &d.AfterName
	case "End":
		return &d.End
	}
	return nil
}
func (d *IncDecStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterX":
		return &d.AfterX
	case "End":
		return &d.End
	}
	return nil
}
func (d *IndexExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterX":
		return &d.AfterX
	case "AfterLbrack":
		return &d.AfterLbrack
	case "AfterIndex":
		return &d.AfterIndex
	case "End":
		return &d.End
	}
	return nil
}
func (d *InterfaceTypeDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterInterface":
		return &d.AfterInterface
	case "End":
		return &d.End
	}
	return nil
}
func (d *KeyValueExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterKey":
		return &d.AfterKey
	case "AfterColon":
		return &d.AfterColon
	case "End":
		return &d.End
	}
	return nil
}
func (d *LabeledStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterLabel":
		return &d.AfterLabel
	case "AfterColon":
		return &d.AfterColon
	case "End":
		return &d.End
	}
	return nil
}
func (d *MapTypeDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterMap":
		return &d.AfterMap
	case "AfterKey":
		return &d.AfterKey
	case "End":
		return &d.End
	}
	return nil
}
func (d *ParenExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterLparen":
		return &d.AfterLparen
	case "AfterX":
		return &d.AfterX
	case "End":
		return &d.End
	}
	return nil
}
func (d *RangeStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterFor":
		return &d.AfterFor
	case "AfterKey":
		return &d.AfterKey
	case "AfterValue":
		return &d.AfterValue
	case "AfterRange":
		return &d.AfterRange
	case "AfterX":
		return &d.AfterX
	case "End":
		return &d.End
	}
	return nil
}
func (d *ReturnStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterReturn":
		return &d.AfterReturn
	case "End":
		return &d.End
	}
	return nil
}
func (d *SelectStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterSelect":
		return &d.AfterSelect
	case "End":
		return &d.End
	}
	return nil
}
func (d *SelectorExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterX":
		return &d.AfterX
	case "End":
		return &d.End
	}
	return nil
}
func (d *SendStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterChan":
		return &d.AfterChan
	case "AfterArrow":
		return &d.AfterArrow
	case "End":
		return &d.End
	}
	return nil
}
func (d *SliceExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterX":
		return &d.AfterX
	case "AfterLbrack":
		return &d.AfterLbrack
	case "AfterLow":
		return &d.AfterLow
	case "AfterHigh":
		return &d.AfterHigh
	case "AfterMax":
		return &d.AfterMax
	case "End":
		return &d.End
	}
	return nil
}
func (d *StarExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterStar":
		return &d.AfterStar
	case "End":
		return &d.End
	}
	return nil
}
func (d *StructTypeDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterStruct":
		return &d.AfterStruct
	case "End":
		return &d.End
	}
	return nil
}
func (d *SwitchStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterSwitch":
		return &d.AfterSwitch
	case "AfterInit":
		return &d.AfterInit
	case "AfterTag":
		return &d.AfterTag
	case "End":
		return &d.End
	}
	return nil
}
func (d *TypeAssertExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterX":
		return &d.AfterX
	case "AfterLparen":
		return &d.AfterLparen
	case "AfterType":
		return &d.AfterType
	case "End":
		return &d.End
	}
	return nil
}
func (d *TypeSpecDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterName":
		return &d.AfterName
	case "End":
		return &d.End
	}
	return nil
}
func (d *TypeSwitchStmtDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterSwitch":
		return &d.AfterSwitch
	case "AfterInit":
		return &d.AfterInit
	case "AfterAssign":
		return &d.AfterAssign
	case "End":
		return &d.End
	}
	return nil
}
func (d *UnaryExprDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterOp":
		return &d.AfterOp
	case "End":
		return &d.End
	}
	return nil
}
func (d *ValueSpecDecorations) Get(name string) *Decorations {
	switch name {
	case "Start":
		return &d.Start
	case "AfterNames":
		return &d.AfterNames
	case "AfterAssign":
		return &d.AfterAssign
	case "End":
		return &d.End
	}
	return nil
}
