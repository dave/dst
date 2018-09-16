package main

type NodeInfo struct {
	Name          string
	Fragments     []FragmentInfo
	Data          []FieldInfo
	FragmentOrder []string
	FromToLength  bool // From, To -> Length (BadXXX)
}

type FragmentInfo struct {
	Name             string
	AstType          string
	AstTypeActual    string
	AstTypePointer   bool
	DstType          string
	DstTypeActual    string
	DstTypePointer   bool
	AstPositionField string
	DstExistsField   string
}

type FieldInfo struct {
	Name    string
	Type    string
	Actual  string
	Pointer bool
}
