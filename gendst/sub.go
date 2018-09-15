package main

type FragmentInfo struct {
	Node           string
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
