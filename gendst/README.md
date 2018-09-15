# FragmentInfo

```go
package main

// This file serves as documentation to explain how the code generation works

var _ = map[string]NodeInfo{
	"ArrayType": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    1,
				Name:      "Lbrack",
				PosField:  "Lbrack",
				Type:      "Pos",
			},
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "Len",
				SuffixLength: 1,
				Type:         "Expr",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Elt",
				Type:   "Expr",
			},
		},
		Name: "ArrayType",
	},
	"AssignStmt": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Lhs",
				Slice:  true,
				Type:   "Expr",
			},
			{
				LenFieldToken: "Tok",
				Name:          "Tok",
				PosField:      "TokPos",
				Type:          "Token",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Rhs",
				Slice:  true,
				Type:   "Expr",
			},
		},
		Name: "AssignStmt",
	},
	"BadDecl": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Name:      "From",
				PosField:  "From",
				Type:      "Pos",
			},
			{
				HasLength: true,
				Name:      "To",
				PosField:  "To",
				Type:      "Pos",
			},
		},
		Name: "BadDecl",
	},
	"BadExpr": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Name:      "From",
				PosField:  "From",
				Type:      "Pos",
			},
			{
				HasLength: true,
				Name:      "To",
				PosField:  "To",
				Type:      "Pos",
			},
		},
		Name: "BadExpr",
	},
	"BadStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Name:      "From",
				PosField:  "From",
				Type:      "Pos",
			},
			{
				HasLength: true,
				Name:      "To",
				PosField:  "To",
				Type:      "Pos",
			},
		},
		Name: "BadStmt",
	},
	"BasicLit": {
		Fragments: []FragmentInfo{
			{
				LenFieldString: "Value",
				Name:           "Value",
				PosField:       "ValuePos",
				Type:           "String",
			},
		},
		Name: "BasicLit",
	},
	"BinaryExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
			{
				LenFieldToken: "Op",
				Name:          "Op",
				PosField:      "OpPos",
				Type:          "Token",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Y",
				Type:   "Expr",
			},
		},
		Name: "BinaryExpr",
	},
	"BlockStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    1,
				Name:      "Lbrace",
				PosField:  "Lbrace",
				Type:      "Pos",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "List",
				Slice:  true,
				Type:   "Stmt",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rbrace",
				PosField:  "Rbrace",
				Type:      "Pos",
			},
		},
		Name: "BlockStmt",
	},
	"BranchStmt": {
		Fragments: []FragmentInfo{
			{
				LenFieldToken: "Tok",
				Name:          "Tok",
				PosField:      "TokPos",
				Type:          "Token",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Label",
				Type:   "Ident",
			},
		},
		Name: "BranchStmt",
	},
	"CallExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Fun",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Lparen",
				PosField:  "Lparen",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Args",
				Slice:  true,
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    3,
				Name:      "Ellipsis",
				PosField:  "Ellipsis",
				Type:      "Pos",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rparen",
				PosField:  "Rparen",
				Type:      "Pos",
			},
		},
		Name: "CallExpr",
	},
	"CaseClause": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    4,
				Name:      "Case",
				PosField:  "Case",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "List",
				Slice:  true,
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Colon",
				PosField:  "Colon",
				Type:      "Pos",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Slice:  true,
				Type:   "Stmt",
			},
		},
		Name: "CaseClause",
	},
	"ChanType": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    4,
				Name:      "Begin",
				PosField:  "Begin",
				Type:      "Pos",
			},
			{
				HasLength: true,
				Length:    2,
				Name:      "Arrow",
				PosField:  "Arrow",
				Special:   true,
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Value",
				Type:   "Expr",
			},
		},
		Name: "ChanType",
	},
	"CommClause": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    4,
				Name:      "Case",
				PosField:  "Case",
				Special:   true,
				Type:      "Pos",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Comm",
				Type:   "Stmt",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Colon",
				PosField:  "Colon",
				Type:      "Pos",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Slice:  true,
				Type:   "Stmt",
			},
		},
		Name: "CommClause",
	},
	"Comment": {
		Fragments: []FragmentInfo{
			{
				LenFieldString: "Text",
				Name:           "Text",
				PosField:       "Slash",
				Type:           "String",
			},
		},
		Name: "Comment",
	},
	"CommentGroup": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "List",
				Slice:  true,
				Type:   "Comment",
			},
		},
		Name: "CommentGroup",
	},
	"CompositeLit": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Type",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Lbrace",
				PosField:  "Lbrace",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Elts",
				Slice:  true,
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rbrace",
				PosField:  "Rbrace",
				Type:      "Pos",
			},
		},
		Name: "CompositeLit",
	},
	"DeclStmt": {
		Fragments: []FragmentInfo{
			{
				IsDecl: true,
				IsNode: true,
				Name:   "Decl",
				Type:   "Decl",
			},
		},
		Name: "DeclStmt",
	},
	"DeferStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    5,
				Name:      "Defer",
				PosField:  "Defer",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Call",
				Type:   "CallExpr",
			},
		},
		Name: "DeferStmt",
	},
	"Ellipsis": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    3,
				Name:      "Ellipsis",
				PosField:  "Ellipsis",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Elt",
				Type:   "Expr",
			},
		},
		Name: "Ellipsis",
	},
	"EmptyStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Name:      "Semicolon",
				PosField:  "Semicolon",
				Special:   true,
				Type:      "Pos",
			},
		},
		Name: "EmptyStmt",
	},
	"ExprStmt": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
		},
		Name: "ExprStmt",
	},
	"Field": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Names",
				Slice:  true,
				Type:   "Ident",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Type",
				Type:   "Expr",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Tag",
				Type:   "BasicLit",
			},
			{
				IsNode: true,
				Name:   "Comment",
				Type:   "CommentGroup",
			},
		},
		Name: "Field",
	},
	"FieldList": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    1,
				Name:      "Opening",
				PosField:  "Opening",
				Type:      "Pos",
			},
			{
				IsNode: true,
				Name:   "List",
				Slice:  true,
				Type:   "Field",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Closing",
				PosField:  "Closing",
				Type:      "Pos",
			},
		},
		Name: "FieldList",
	},
	"File": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				HasLength: true,
				Length:    7,
				Name:      "Package",
				PosField:  "Package",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Name",
				Type:   "Ident",
			},
			{
				IsDecl: true,
				IsNode: true,
				Name:   "Decls",
				Slice:  true,
				Type:   "Decl",
			},
		},
		Name: "File",
	},
	"ForStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    3,
				Name:      "For",
				PosField:  "For",
				Type:      "Pos",
			},
			{
				IsNode:       true,
				IsStmt:       true,
				Name:         "Init",
				SuffixLength: 1,
				Type:         "Stmt",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Cond",
				Type:   "Expr",
			},
			{
				IsNode:       true,
				IsStmt:       true,
				Name:         "Post",
				PrefixLength: 1,
				Type:         "Stmt",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "ForStmt",
	},
	"FuncDecl": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				IsNode: true,
				Name:   "Recv",
				Type:   "FieldList",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Name",
				Type:   "Ident",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Type",
				Type:   "FuncType",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "FuncDecl",
	},
	"FuncLit": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Type",
				Type:   "FuncType",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "FuncLit",
	},
	"FuncType": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    4,
				Name:      "Func",
				PosField:  "Func",
				Type:      "Pos",
			},
			{
				IsNode: true,
				Name:   "Params",
				Type:   "FieldList",
			},
			{
				IsNode: true,
				Name:   "Results",
				Type:   "FieldList",
			},
		},
		Name: "FuncType",
	},
	"GenDecl": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				LenFieldToken: "Tok",
				Name:          "Tok",
				PosField:      "TokPos",
				Type:          "Token",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Lparen",
				PosField:  "Lparen",
				Type:      "Pos",
			},
			{
				IsNode: true,
				Name:   "Specs",
				Slice:  true,
				Type:   "Spec",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rparen",
				PosField:  "Rparen",
				Type:      "Pos",
			},
		},
		Name: "GenDecl",
	},
	"GoStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    2,
				Name:      "Go",
				PosField:  "Go",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Call",
				Type:   "CallExpr",
			},
		},
		Name: "GoStmt",
	},
	"Ident": {
		Fragments: []FragmentInfo{
			{
				LenFieldString: "Name",
				Name:           "Name",
				PosField:       "NamePos",
				Type:           "String",
			},
		},
		Name: "Ident",
	},
	"IfStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    2,
				Name:      "If",
				PosField:  "If",
				Type:      "Pos",
			},
			{
				IsNode:       true,
				IsStmt:       true,
				Name:         "Init",
				SuffixLength: 1,
				Type:         "Stmt",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Cond",
				Type:   "Expr",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
			{
				IsNode:       true,
				IsStmt:       true,
				Name:         "Else",
				SuffixLength: 4,
				Type:         "Stmt",
			},
		},
		Name: "IfStmt",
	},
	"ImportSpec": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Name",
				Type:   "Ident",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Path",
				Type:   "BasicLit",
			},
			{
				IsNode: true,
				Name:   "Comment",
				Type:   "CommentGroup",
			},
		},
		Name: "ImportSpec",
	},
	"IncDecStmt": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
			{
				LenFieldToken: "Tok",
				Name:          "Tok",
				PosField:      "TokPos",
				Type:          "Token",
			},
		},
		Name: "IncDecStmt",
	},
	"IndexExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Lbrack",
				PosField:  "Lbrack",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Index",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rbrack",
				PosField:  "Rbrack",
				Type:      "Pos",
			},
		},
		Name: "IndexExpr",
	},
	"InterfaceType": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    9,
				Name:      "Interface",
				PosField:  "Interface",
				Type:      "Pos",
			},
			{
				IsNode: true,
				Name:   "Methods",
				Type:   "FieldList",
			},
		},
		Name: "InterfaceType",
	},
	"KeyValueExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Key",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Colon",
				PosField:  "Colon",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Value",
				Type:   "Expr",
			},
		},
		Name: "KeyValueExpr",
	},
	"LabeledStmt": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Label",
				Type:   "Ident",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Colon",
				PosField:  "Colon",
				Type:      "Pos",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Stmt",
				Type:   "Stmt",
			},
		},
		Name: "LabeledStmt",
	},
	"MapType": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    3,
				Name:      "Map",
				PosField:  "Map",
				Type:      "Pos",
			},
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "Key",
				PrefixLength: 1,
				SuffixLength: 1,
				Type:         "Expr",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Value",
				Type:   "Expr",
			},
		},
		Name: "MapType",
	},
	"ParenExpr": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    1,
				Name:      "Lparen",
				PosField:  "Lparen",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rparen",
				PosField:  "Rparen",
				Type:      "Pos",
			},
		},
		Name: "ParenExpr",
	},
	"RangeStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    3,
				Name:      "For",
				PosField:  "For",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Key",
				Type:   "Expr",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Value",
				Type:   "Expr",
			},
			{
				LenFieldToken: "Tok",
				Name:          "Tok",
				PosField:      "TokPos",
				Type:          "Token",
			},
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "X",
				PrefixLength: 5,
				Type:         "Expr",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "RangeStmt",
	},
	"ReturnStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    6,
				Name:      "Return",
				PosField:  "Return",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Results",
				Slice:  true,
				Type:   "Expr",
			},
		},
		Name: "ReturnStmt",
	},
	"SelectStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    6,
				Name:      "Select",
				PosField:  "Select",
				Type:      "Pos",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "SelectStmt",
	},
	"SelectorExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "X",
				SuffixLength: 1,
				Type:         "Expr",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Sel",
				Type:   "Ident",
			},
		},
		Name: "SelectorExpr",
	},
	"SendStmt": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Chan",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    2,
				Name:      "Arrow",
				PosField:  "Arrow",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Value",
				Type:   "Expr",
			},
		},
		Name: "SendStmt",
	},
	"SliceExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Lbrack",
				PosField:  "Lbrack",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Low",
				Type:   "Expr",
			},
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "High",
				PrefixLength: 1,
				Type:         "Expr",
			},
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "Max",
				PrefixLength: 1,
				Special:      true,
				Type:         "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rbrack",
				PosField:  "Rbrack",
				Type:      "Pos",
			},
		},
		Name: "SliceExpr",
	},
	"StarExpr": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    1,
				Name:      "Star",
				PosField:  "Star",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
		},
		Name: "StarExpr",
	},
	"StructType": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    6,
				Name:      "Struct",
				PosField:  "Struct",
				Type:      "Pos",
			},
			{
				IsNode: true,
				Name:   "Fields",
				Type:   "FieldList",
			},
		},
		Name: "StructType",
	},
	"SwitchStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    6,
				Name:      "Switch",
				PosField:  "Switch",
				Type:      "Pos",
			},
			{
				IsNode:       true,
				IsStmt:       true,
				Name:         "Init",
				SuffixLength: 1,
				Type:         "Stmt",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Tag",
				Type:   "Expr",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "SwitchStmt",
	},
	"TypeAssertExpr": {
		Fragments: []FragmentInfo{
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Lparen",
				PosField:  "Lparen",
				Type:      "Pos",
			},
			{
				IsExpr:  true,
				IsNode:  true,
				Name:    "Type",
				Special: true,
				Type:    "Expr",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Rparen",
				PosField:  "Rparen",
				Type:      "Pos",
			},
		},
		Name: "TypeAssertExpr",
	},
	"TypeSpec": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Name",
				Type:   "Ident",
			},
			{
				HasLength: true,
				Length:    1,
				Name:      "Assign",
				PosField:  "Assign",
				Type:      "Pos",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Type",
				Type:   "Expr",
			},
			{
				IsNode: true,
				Name:   "Comment",
				Type:   "CommentGroup",
			},
		},
		Name: "TypeSpec",
	},
	"TypeSwitchStmt": {
		Fragments: []FragmentInfo{
			{
				HasLength: true,
				Length:    6,
				Name:      "Switch",
				PosField:  "Switch",
				Type:      "Pos",
			},
			{
				IsNode:       true,
				IsStmt:       true,
				Name:         "Init",
				SuffixLength: 1,
				Type:         "Stmt",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Assign",
				Type:   "Stmt",
			},
			{
				IsNode: true,
				IsStmt: true,
				Name:   "Body",
				Type:   "BlockStmt",
			},
		},
		Name: "TypeSwitchStmt",
	},
	"UnaryExpr": {
		Fragments: []FragmentInfo{
			{
				LenFieldToken: "Op",
				Name:          "Op",
				PosField:      "OpPos",
				Type:          "Token",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "X",
				Type:   "Expr",
			},
		},
		Name: "UnaryExpr",
	},
	"ValueSpec": {
		Fragments: []FragmentInfo{
			{
				IsNode: true,
				Name:   "Doc",
				Type:   "CommentGroup",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Names",
				Slice:  true,
				Type:   "Ident",
			},
			{
				IsExpr: true,
				IsNode: true,
				Name:   "Type",
				Type:   "Expr",
			},
			{
				IsExpr:       true,
				IsNode:       true,
				Name:         "Values",
				PrefixLength: 1,
				Slice:        true,
				Type:         "Expr",
			},
			{
				IsNode: true,
				Name:   "Comment",
				Type:   "CommentGroup",
			},
		},
		Name: "ValueSpec",
	},
}

```