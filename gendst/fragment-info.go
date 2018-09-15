package main

type NodeInfo struct {
	Name      string
	Fragments []FragmentInfo
}

type FragmentInfo struct {
	Node           *NodeInfo
	Name           string
	Type           string
	IsNode         bool
	IsStmt         bool
	IsExpr         bool
	IsDecl         bool
	PosField       string
	Slice          bool
	HasLength      bool
	Length         int
	LenFieldString string
	LenFieldToken  string
	PrefixLength   int
	SuffixLength   int
	Special        bool
}
