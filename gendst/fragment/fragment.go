package fragment

import (
	"go/token"

	"github.com/dave/jennifer/jen"
)

const DSTPATH = "github.com/dave/dst"

type Part interface{}

var Info = map[string][]Part{
	/*
		// A Field represents a Field declaration list in a struct type,
		// a method list in an interface type, or a parameter/result declaration
		// in a signature.
		// Field.Names is nil for unnamed parameters (parameter lists which only contain types)
		// and embedded struct fields. In the latter case, the field name is the type name.
		//
		type Field struct {
			Doc     *CommentGroup // associated documentation; or nil
			Names   []*Ident      // field/method/parameter names; or nil
			Type    Expr          // field/method/parameter type
			Tag     *BasicLit     // field tag; or nil
			Comment *CommentGroup // line comments; or nil
		}
	*/
	"Field": {
		Decoration{
			Name: "Start",
		},
		List{
			Name:      "Names",
			Field:     Field{"Names"},
			Elem:      Type{"Ident", true},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterNames",
		},
		Node{
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterType",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tag").Op("!=").Nil() }),
		},
		Node{
			Name:  "Tag",
			Field: Field{"Tag"},
			Type:  Type{"BasicLit", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A FieldList represents a list of Fields, enclosed by parentheses or braces.
		type FieldList struct {
			Opening token.Pos // position of opening parenthesis/brace, if any
			List    []*Field  // field list; or nil
			Closing token.Pos // position of closing parenthesis/brace, if any
		}
	*/
	"FieldList": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:  "Opening",
			Token: Basic{jen.Qual("go/token", "LPAREN")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Opening").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Opening") }),
			},
			ExistsField:   Field{"Opening"},
			PositionField: Field{"Opening"},
		},
		Decoration{
			Name: "AfterOpening",
		},
		List{
			Name:      "List",
			Field:     Field{"List"},
			Elem:      Type{"Field", true},
			Separator: token.COMMA,
		},
		Token{
			Name:  "Closing",
			Token: Basic{jen.Qual("go/token", "RPAREN")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Closing").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Closing") }),
			},
			ExistsField:   Field{"Closing"},
			PositionField: Field{"Closing"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A BadExpr node is a placeholder for expressions containing
		// syntax errors for which no correct expression nodes can be
		// created.
		//
		BadExpr struct {
			From, To token.Pos // position range of bad expression
		}
	*/
	"BadExpr": {
		Ignored{
			Length: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement {
					return jen.Int().Parens(jen.Add(n).Dot("End").Call().Op("-").Add(n).Dot("Pos").Call())
				}),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Length") }),
			},
			PositionField: Field{"From"},
		},
	},
	/*
		// An Ident node represents an identifier.
		Ident struct {
			NamePos token.Pos // identifier position
			Name    string    // identifier name
			Obj     *Object   // denoted object; or nil
		}
	*/
	"Ident": {
		Decoration{
			Name: "Start",
		},
		String{
			Name:          "Name",
			ValueField:    Field{"Name"},
			PositionField: Field{"NamePos"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An Ellipsis node stands for the "..." type in a
		// parameter list or the "..." length in an array type.
		//
		Ellipsis struct {
			Ellipsis token.Pos // position of "..."
			Elt      Expr      // ellipsis element type (parameter lists only); or nil
		}
	*/
	"Ellipsis": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Ellipsis",
			Token:         Basic{jen.Qual("go/token", "ELLIPSIS")},
			PositionField: Field{"Ellipsis"},
		},
		Decoration{
			Name: "AfterEllipsis",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Elt").Op("!=").Nil() }),
		},
		Node{
			Name:  "Elt",
			Field: Field{"Elt"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A BasicLit node represents a literal of basic type.
		BasicLit struct {
			ValuePos token.Pos   // literal position
			Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
			Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
		}
	*/
	"BasicLit": {
		Decoration{
			Name: "Start",
		},
		String{
			Name:          "Value",
			ValueField:    Field{"Value"},
			PositionField: Field{"ValuePos"},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name:  "Kind",
			Field: Field{"Kind"},
		},
	},
	/*
		// A FuncLit node represents a function literal.
		FuncLit struct {
			Type *FuncType  // function type
			Body *BlockStmt // function body
		}
	*/
	"FuncLit": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"FuncType", true},
		},
		Decoration{
			Name: "AfterType",
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A CompositeLit node represents a composite literal.
		CompositeLit struct {
			Type       Expr      // literal type; or nil
			Lbrace     token.Pos // position of "{"
			Elts       []Expr    // list of composite elements; or nil
			Rbrace     token.Pos // position of "}"
			Incomplete bool      // true if (source) expressions are missing in the Elts list
		}
	*/
	"CompositeLit": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterType",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Type").Op("!=").Nil() }),
		},
		Token{
			Name:          "Lbrace",
			Token:         Basic{jen.Qual("go/token", "LBRACE")},
			PositionField: Field{"Lbrace"},
		},
		Decoration{
			Name: "AfterLbrace",
		},
		List{
			Name:      "Elts",
			Field:     Field{"Elts"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		// TODO: removed AfterElts decoration - should we remove all decorations after comma separated lists?
		Token{
			Name:          "Rbrace",
			Token:         Basic{jen.Qual("go/token", "RBRACE")},
			PositionField: Field{"Rbrace"},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name:  "Incomplete",
			Field: Field{"Incomplete"},
		},
	},
	/*
		// A ParenExpr node represents a parenthesized expression.
		ParenExpr struct {
			Lparen token.Pos // position of "("
			X      Expr      // parenthesized expression
			Rparen token.Pos // position of ")"
		}
	*/
	"ParenExpr": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Lparen",
			Token:         Basic{jen.Qual("go/token", "LPAREN")},
			PositionField: Field{"Lparen"},
		},
		Decoration{
			Name: "AfterLparen",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:          "Rparen",
			Token:         Basic{jen.Qual("go/token", "RPAREN")},
			PositionField: Field{"Rparen"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A SelectorExpr node represents an expression followed by a selector.
		SelectorExpr struct {
			X   Expr   // expression
			Sel *Ident // field selector
		}
	*/
	"SelectorExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:  "Period",
			Token: Basic{jen.Qual("go/token", "PERIOD")},
		},
		Decoration{
			Name: "AfterX",
		},
		Node{
			Name:  "Sel",
			Field: Field{"Sel"},
			Type:  Type{"Ident", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An IndexExpr node represents an expression followed by an index.
		IndexExpr struct {
			X      Expr      // expression
			Lbrack token.Pos // position of "["
			Index  Expr      // index expression
			Rbrack token.Pos // position of "]"
		}
	*/
	"IndexExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:          "Lbrack",
			Token:         Basic{jen.Qual("go/token", "LBRACK")},
			PositionField: Field{"Lbrack"},
		},
		Decoration{
			Name: "AfterLbrack",
		},
		Node{
			Name:  "Index",
			Field: Field{"Index"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterIndex",
		},
		Token{
			Name:          "Rbrack",
			Token:         Basic{jen.Qual("go/token", "RBRACK")},
			PositionField: Field{"Rbrack"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An SliceExpr node represents an expression followed by slice indices.
		SliceExpr struct {
			X      Expr      // expression
			Lbrack token.Pos // position of "["
			Low    Expr      // begin of slice range; or nil
			High   Expr      // end of slice range; or nil
			Max    Expr      // maximum capacity of slice; or nil
			Slice3 bool      // true if 3-index slice (2 colons present)
			Rbrack token.Pos // position of "]"
		}
	*/
	// var H = /*Start*/ []int{0} /*AfterX*/ [ /*AfterLbrack*/ 1: /*AfterLow*/ 2: /*AfterHigh*/ 3 /*AfterMax*/] /*End*/
	// TODO: Why Slice3? Why not Max != nil... Can we have Max == nil && Slice3 == true?
	"SliceExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:          "Lbrack",
			Token:         Basic{jen.Qual("go/token", "LBRACK")},
			PositionField: Field{"Lbrack"},
		},
		Decoration{
			Name: "AfterLbrack",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Low").Op("!=").Nil() }),
		},
		Node{
			Name:  "Low",
			Field: Field{"Low"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:  "Colon1",
			Token: Basic{jen.Qual("go/token", "COLON")},
		},
		Decoration{
			Name: "AfterLow",
		},
		Node{
			Name:  "High",
			Field: Field{"High"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:   "Colon2",
			Token:  Basic{jen.Qual("go/token", "COLON")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Slice3") }),
		},
		Decoration{
			Name: "AfterHigh",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("High").Op("!=").Nil() }),
		},
		Node{
			Name:  "Max",
			Field: Field{"Max"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterMax",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Max").Op("!=").Nil() }), // TODO - Slice3 in here?
		},
		Token{
			Name:          "Rbrack",
			Token:         Basic{jen.Qual("go/token", "RBRACK")},
			PositionField: Field{"Rbrack"},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name:  "Slice3",
			Field: Field{"Slice3"},
		},
	},
	/*
		// A TypeAssertExpr node represents an expression followed by a
		// type assertion.
		//
		TypeAssertExpr struct {
			X      Expr      // expression
			Lparen token.Pos // position of "("
			Type   Expr      // asserted type; nil means type switch X.(type)
			Rparen token.Pos // position of ")"
		}
	*/
	"TypeAssertExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:  "Period",
			Token: Basic{jen.Qual("go/token", "PERIOD")},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:          "Lparen",
			Token:         Basic{jen.Qual("go/token", "LPAREN")},
			PositionField: Field{"Lparen"},
		},
		Decoration{
			Name: "AfterLparen",
		},
		Node{
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:   "TypeToken",
			Token:  Basic{jen.Qual("go/token", "TYPE")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Type").Op("==").Nil() }),
		},
		Decoration{
			Name: "AfterType",
		},
		Token{
			Name:          "Rparen",
			Token:         Basic{jen.Qual("go/token", "RPAREN")},
			PositionField: Field{"Rparen"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A CallExpr node represents an expression followed by an argument list.
		CallExpr struct {
			Fun      Expr      // function expression
			Lparen   token.Pos // position of "("
			Args     []Expr    // function arguments; or nil
			Ellipsis token.Pos // position of "..." (token.NoPos if there is no "...")
			Rparen   token.Pos // position of ")"
		}
	*/
	"CallExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Fun",
			Field: Field{"Fun"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterFun",
		},
		Token{
			Name:          "Lparen",
			Token:         Basic{jen.Qual("go/token", "LPAREN")},
			PositionField: Field{"Lparen"},
		},
		Decoration{
			Name: "AfterLparen",
		},
		List{
			Name:      "Args",
			Field:     Field{"Args"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterArgs",
		},
		Token{
			Name:  "Ellipsis",
			Token: Basic{jen.Qual("go/token", "ELLIPSIS")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Ellipsis").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Ellipsis") }),
			},
			ExistsField:   Field{"Ellipsis"},
			PositionField: Field{"Ellipsis"},
		},
		Decoration{
			Name: "AfterEllipsis",
			Use: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Ellipsis").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Ellipsis") }),
			},
		},
		Token{
			Name:          "Rparen",
			Token:         Basic{jen.Qual("go/token", "RPAREN")},
			PositionField: Field{"Rparen"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A StarExpr node represents an expression of the form "*" Expression.
		// Semantically it could be a unary "*" expression, or a pointer type.
		//
		StarExpr struct {
			Star token.Pos // position of "*"
			X    Expr      // operand
		}
	*/
	"StarExpr": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Star",
			Token:         Basic{jen.Qual("go/token", "MUL")},
			PositionField: Field{"Star"},
		},
		Decoration{
			Name: "AfterStar",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A UnaryExpr node represents a unary expression.
		// Unary "*" expressions are represented via StarExpr nodes.
		//
		UnaryExpr struct {
			OpPos token.Pos   // position of Op
			Op    token.Token // operator
			X     Expr        // operand
		}
	*/
	"UnaryExpr": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Op",
			TokenField:    Field{"Op"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Op") }),
			PositionField: Field{"OpPos"},
		},
		Decoration{
			Name: "AfterOp",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A BinaryExpr node represents a binary expression.
		BinaryExpr struct {
			X     Expr        // left operand
			OpPos token.Pos   // position of Op
			Op    token.Token // operator
			Y     Expr        // right operand
		}
	*/
	"BinaryExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:          "Op",
			TokenField:    Field{"Op"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Op") }),
			PositionField: Field{"OpPos"},
		},
		Decoration{
			Name: "AfterOp",
		},
		Node{
			Name:  "Y",
			Field: Field{"Y"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A KeyValueExpr node represents (key : value) pairs
		// in composite literals.
		//
		KeyValueExpr struct {
			Key   Expr
			Colon token.Pos // position of ":"
			Value Expr
		}
	*/
	"KeyValueExpr": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Key",
			Field: Field{"Key"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterKey",
		},
		Token{
			Name:          "Colon",
			Token:         Basic{jen.Qual("go/token", "COLON")},
			PositionField: Field{"Colon"},
		},
		Decoration{
			Name: "AfterColon",
		},
		Node{
			Name:  "Value",
			Field: Field{"Value"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An ArrayType node represents an array or slice type.
		ArrayType struct {
			Lbrack token.Pos // position of "["
			Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
			Elt    Expr      // element type
		}
	*/
	"ArrayType": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Lbrack",
			Token:         Basic{jen.Qual("go/token", "LBRACK")},
			PositionField: Field{"Lbrack"},
		},
		Decoration{
			Name: "AfterLbrack",
		},
		Node{
			Name:  "Len",
			Field: Field{"Len"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:  "Rbrack",
			Token: Basic{jen.Qual("go/token", "RBRACK")},
		},
		Decoration{
			Name: "AfterLen",
		},
		Node{
			Name:  "Elt",
			Field: Field{"Elt"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A StructType node represents a struct type.
		StructType struct {
			Struct     token.Pos  // position of "struct" keyword
			Fields     *FieldList // list of field declarations
			Incomplete bool       // true if (source) fields are missing in the Fields list
		}
	*/
	"StructType": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Struct",
			Token:         Basic{jen.Qual("go/token", "STRUCT")},
			PositionField: Field{"Struct"},
		},
		Decoration{
			Name: "AfterStruct",
		},
		Node{
			Name:  "Fields",
			Field: Field{"Fields"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name:  "Incomplete",
			Field: Field{"Incomplete"}, // TODO: Remove this and always set to false?
		},
	},
	/*
		// A FuncType node represents a function type.
		FuncType struct {
			Func    token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
			Params  *FieldList // (incoming) parameters; non-nil
			Results *FieldList // (outgoing) results; or nil
		}
	*/
	"FuncType": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:  "Func",
			Token: Basic{jen.Qual("go/token", "FUNC")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Func").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Func") }),
			},
			ExistsField:   Field{"Func"},
			PositionField: Field{"Func"},
		},
		Decoration{
			Name: "AfterFunc",
			Use: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Func").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Func") }),
			},
		},
		Node{
			Name:  "Params",
			Field: Field{"Params"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "AfterParams",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Results").Op("!=").Nil() }),
		},
		Node{
			Name:  "Results",
			Field: Field{"Results"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An InterfaceType node represents an interface type.
		InterfaceType struct {
			Interface  token.Pos  // position of "interface" keyword
			Methods    *FieldList // list of methods
			Incomplete bool       // true if (source) methods are missing in the Methods list
		}
	*/
	"InterfaceType": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Interface",
			Token:         Basic{jen.Qual("go/token", "INTERFACE")},
			PositionField: Field{"Interface"},
		},
		Decoration{
			Name: "AfterInterface",
		},
		Node{
			Name:  "Methods",
			Field: Field{"Methods"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name:  "Incomplete",
			Field: Field{"Incomplete"}, // TODO: Remove this and always set to false?
		},
	},
	/*
		// A MapType node represents a map type.
		MapType struct {
			Map   token.Pos // position of "map" keyword
			Key   Expr
			Value Expr
		}
	*/
	"MapType": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Map",
			Token:         Basic{jen.Qual("go/token", "MAP")},
			PositionField: Field{"Map"},
		},
		Token{
			Name:  "Lbrack",
			Token: Basic{jen.Qual("go/token", "LBRACK")},
		},
		Decoration{
			Name: "AfterMap",
		},
		Node{
			Name:  "Key",
			Field: Field{"Key"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:  "Rbrack",
			Token: Basic{jen.Qual("go/token", "RBRACK")},
		},
		Decoration{
			Name: "AfterKey",
		},
		Node{
			Name:  "Value",
			Field: Field{"Value"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A ChanType node represents a channel type.
		ChanType struct {
			Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
			Arrow token.Pos // position of "<-" (token.NoPos if there is no "<-")
			Dir   ChanDir   // channel direction
			Value Expr      // value type
		}
	*/
	"ChanType": {
		Decoration{
			Name: "Start",
		},
		// This is rather a kludge. In RECV variation, we emit "<-" followed by "chan". Otherwise we
		// just emit "chan".
		Token{
			Name: "Begin",
			Token: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement {
					return jen.Func().Params().Qual("go/token", "Token").Block(jen.If(n.Dot("Dir").Op("==").Qual("go/ast", "RECV")).Block(jen.Return(jen.Qual("go/token", "ARROW"))).Else().Block(jen.Return(jen.Qual("go/token", "CHAN")))).Call()
				}),
				Dst: Expr(func(n *jen.Statement) *jen.Statement {
					return jen.Func().Params().Qual("go/token", "Token").Block(jen.If(n.Dot("Dir").Op("==").Qual(DSTPATH, "RECV")).Block(jen.Return(jen.Qual("go/token", "ARROW"))).Else().Block(jen.Return(jen.Qual("go/token", "CHAN")))).Call()
				}),
			},
			PositionField: Field{"Begin"},
		},
		Token{
			Name:  "Chan",
			Token: Basic{jen.Qual("go/token", "CHAN")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Dir").Op("==").Qual("go/ast", "RECV") }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Dir").Op("==").Qual(DSTPATH, "RECV") }),
			},
		},
		Decoration{
			Name: "AfterBegin",
		},
		Token{
			Name:  "Arrow",
			Token: Basic{jen.Qual("go/token", "ARROW")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Dir").Op("==").Qual("go/ast", "SEND") }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Dir").Op("==").Qual(DSTPATH, "SEND") }),
			},
			PositionField: Field{"Arrow"},
		},
		Decoration{
			Name: "AfterArrow",
			Use: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Dir").Op("==").Qual("go/ast", "SEND") }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Dir").Op("==").Qual(DSTPATH, "SEND") }),
			},
		},
		Node{
			Name:  "Value",
			Field: Field{"Value"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name:  "Dir",
			Field: Field{"Dir"},
			Value: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return jen.Qual(DSTPATH, "ChanDir").Parens(n.Dot("Dir")) }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return jen.Qual("go/ast", "ChanDir").Parens(n.Dot("Dir")) }),
			},
		},
	},
	/*
		// A BadStmt node is a placeholder for statements containing
		// syntax errors for which no correct statement nodes can be
		// created.
		//
		BadStmt struct {
			From, To token.Pos // position range of bad statement
		}
	*/
	"BadStmt": {
		Ignored{
			Length: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement {
					return jen.Int().Parens(jen.Add(n).Dot("End").Call().Op("-").Add(n).Dot("Pos").Call())
				}),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Length") }),
			},
			PositionField: Field{"From"},
		},
	},
	/*
		// A DeclStmt node represents a declaration in a statement list.
		DeclStmt struct {
			Decl Decl // *GenDecl with CONST, TYPE, or VAR token
		}
	*/
	"DeclStmt": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Decl",
			Field: Field{"Decl"},
			Type:  Type{"Decl", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An EmptyStmt node represents an empty statement.
		// The "position" of the empty statement is the position
		// of the immediately following (explicit or implicit) semicolon.
		//
		EmptyStmt struct {
			Semicolon token.Pos // position of following ";"
			Implicit  bool      // if set, ";" was omitted in the source
		}
	*/
	"EmptyStmt": {
		Token{
			Name:          "Semicolon",
			Token:         Basic{jen.Qual("go/token", "ARROW")},
			Exists:        Expr(func(n *jen.Statement) *jen.Statement { return jen.Op("!").Add(n).Dot("Implicit") }),
			PositionField: Field{"Semicolon"},
		},
		Value{
			Name:  "Implicit",
			Field: Field{"Implicit"},
		},
	},
	/*
		// A LabeledStmt node represents a labeled statement.
		LabeledStmt struct {
			Label *Ident
			Colon token.Pos // position of ":"
			Stmt  Stmt
		}
	*/
	"LabeledStmt": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Label",
			Field: Field{"Label"},
			Type:  Type{"Ident", true},
		},
		Decoration{
			Name: "AfterLabel",
		},
		Token{
			Name:          "Colon",
			Token:         Basic{jen.Qual("go/token", "COLON")},
			PositionField: Field{"Colon"},
		},
		Decoration{
			Name: "AfterColon",
		},
		Node{
			Name:  "Stmt",
			Field: Field{"Stmt"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An ExprStmt node represents a (stand-alone) expression
		// in a statement list.
		//
		ExprStmt struct {
			X Expr // expression
		}
	*/
	"ExprStmt": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A SendStmt node represents a send statement.
		SendStmt struct {
			Chan  Expr
			Arrow token.Pos // position of "<-"
			Value Expr
		}
	*/
	///*Start*/
	//	c /*AfterChan*/ <- /*AfterArrow*/ 0 /*End*/
	"SendStmt": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Chan",
			Field: Field{"Chan"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterChan",
		},
		Token{
			Name:          "Arrow",
			Token:         Basic{jen.Qual("go/token", "ARROW")},
			PositionField: Field{"Arrow"},
		},
		Decoration{
			Name: "AfterArrow",
		},
		Node{
			Name:  "Value",
			Field: Field{"Value"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An IncDecStmt node represents an increment or decrement statement.
		IncDecStmt struct {
			X      Expr
			TokPos token.Pos   // position of Tok
			Tok    token.Token // INC or DEC
		}
	*/
	"IncDecStmt": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:          "Tok",
			TokenField:    Field{"Tok"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tok") }),
			PositionField: Field{"TokPos"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An AssignStmt node represents an assignment or
		// a short variable declaration.
		//
		AssignStmt struct {
			Lhs    []Expr
			TokPos token.Pos   // position of Tok
			Tok    token.Token // assignment token, DEFINE
			Rhs    []Expr
		}
	*/
	"AssignStmt": {
		Decoration{
			Name: "Start",
		},
		List{
			Name:      "Lhs",
			Field:     Field{"Lhs"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterLhs",
		},
		Token{
			Name:          "Tok",
			TokenField:    Field{"Tok"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tok") }),
			PositionField: Field{"TokPos"},
		},
		Decoration{
			Name: "AfterTok",
		},
		List{
			Name:      "Rhs",
			Field:     Field{"Rhs"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A GoStmt node represents a go statement.
		GoStmt struct {
			Go   token.Pos // position of "go" keyword
			Call *CallExpr
		}
	*/
	"GoStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Go",
			Token:         Basic{jen.Qual("go/token", "GO")},
			PositionField: Field{"Go"},
		},
		Decoration{
			Name: "AfterGo",
		},
		Node{
			Name:  "Call",
			Field: Field{"Call"},
			Type:  Type{"CallExpr", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A DeferStmt node represents a defer statement.
		DeferStmt struct {
			Defer token.Pos // position of "defer" keyword
			Call  *CallExpr
		}
	*/
	"DeferStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Defer",
			Token:         Basic{jen.Qual("go/token", "DEFER")},
			PositionField: Field{"Defer"},
		},
		Decoration{
			Name: "AfterDefer",
		},
		Node{
			Name:  "Call",
			Field: Field{"Call"},
			Type:  Type{"CallExpr", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A ReturnStmt node represents a return statement.
		ReturnStmt struct {
			Return  token.Pos // position of "return" keyword
			Results []Expr    // result expressions; or nil
		}
	*/
	"ReturnStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Return",
			Token:         Basic{jen.Qual("go/token", "RETURN")},
			PositionField: Field{"Return"},
		},
		Decoration{
			Name: "AfterReturn",
		},
		List{
			Name:      "Results",
			Field:     Field{"Results"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A BranchStmt node represents a break, continue, goto,
		// or fallthrough statement.
		//
		BranchStmt struct {
			TokPos token.Pos   // position of Tok
			Tok    token.Token // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
			Label  *Ident      // label name; or nil
		}
	*/
	"BranchStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Tok",
			TokenField:    Field{"Tok"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tok") }),
			PositionField: Field{"TokPos"},
		},
		Decoration{
			Name: "AfterTok",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Label").Op("!=").Nil() }),
		},
		Node{
			Name:  "Label",
			Field: Field{"Label"},
			Type:  Type{"Ident", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A BlockStmt node represents a braced statement list.
		BlockStmt struct {
			Lbrace token.Pos // position of "{"
			List   []Stmt
			Rbrace token.Pos // position of "}"
		}
	*/
	"BlockStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Lbrace",
			Token:         Basic{jen.Qual("go/token", "LBRACE")},
			PositionField: Field{"Lbrace"},
		},
		Decoration{
			Name: "AfterLbrace",
		},
		List{
			Name:      "List",
			Field:     Field{"List"},
			Elem:      Type{"Stmt", false},
			Separator: token.SEMICOLON,
		},
		Token{
			Name:          "Rbrace",
			Token:         Basic{jen.Qual("go/token", "RBRACE")},
			PositionField: Field{"Rbrace"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An IfStmt node represents an if statement.
		IfStmt struct {
			If   token.Pos // position of "if" keyword
			Init Stmt      // initialization statement; or nil
			Cond Expr      // condition
			Body *BlockStmt
			Else Stmt // else branch; or nil
		}
	*/
	"IfStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "If",
			Token:         Basic{jen.Qual("go/token", "IF")},
			PositionField: Field{"If"},
		},
		Decoration{
			Name: "AfterIf",
		},
		Node{
			Name:  "Init",
			Field: Field{"Init"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Init").Op("!=").Nil() }),
		},
		Node{
			Name:  "Cond",
			Field: Field{"Cond"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterCond",
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Token{
			Name:   "ElseTok",
			Token:  Basic{jen.Qual("go/token", "ELSE")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Else").Op("!=").Nil() }),
		},
		Decoration{
			Name: "AfterElse",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Else").Op("!=").Nil() }),
		},
		Node{
			Name:  "Else",
			Field: Field{"Else"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A CaseClause represents a case of an expression or type switch statement.
		CaseClause struct {
			Case  token.Pos // position of "case" or "default" keyword
			List  []Expr    // list of expressions or types; nil means default case
			Colon token.Pos // position of ":"
			Body  []Stmt    // statement list; or nil
		}
	*/
	"CaseClause": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name: "Case",
			Token: Expr(func(n *jen.Statement) *jen.Statement {
				return jen.Func().Params().Qual("go/token", "Token").Block(jen.If(n.Dot("List").Op("==").Nil()).Block(jen.Return(jen.Qual("go/token", "DEFAULT"))).Else().Block(jen.Return(jen.Qual("go/token", "CASE")))).Call()
			}),
			PositionField: Field{"Case"},
		},
		Decoration{
			Name: "AfterCase",
		},
		List{
			Name:      "List",
			Field:     Field{"List"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterList",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("List").Op("!=").Nil() }),
		},
		Token{
			Name:          "Colon",
			Token:         Basic{jen.Qual("go/token", "COLON")},
			PositionField: Field{"Colon"},
		},
		Decoration{
			Name: "AfterColon",
		},
		List{
			Name:      "Body",
			Field:     Field{"Body"},
			Elem:      Type{"Stmt", false},
			Separator: token.SEMICOLON,
		},
		// Never want to attach decorations to the end of a list of statements - always better to
		// attach to the last statement.
	},
	/*
		// A SwitchStmt node represents an expression switch statement.
		SwitchStmt struct {
			Switch token.Pos  // position of "switch" keyword
			Init   Stmt       // initialization statement; or nil
			Tag    Expr       // tag expression; or nil
			Body   *BlockStmt // CaseClauses only
		}
	*/
	"SwitchStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Switch",
			Token:         Basic{jen.Qual("go/token", "SWITCH")},
			PositionField: Field{"Switch"},
		},
		Decoration{
			Name: "AfterSwitch",
		},
		Node{
			Name:  "Init",
			Field: Field{"Init"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Init").Op("!=").Nil() }),
		},
		Node{
			Name:  "Tag",
			Field: Field{"Tag"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterTag",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tag").Op("!=").Nil() }),
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An TypeSwitchStmt node represents a type switch statement.
		TypeSwitchStmt struct {
			Switch token.Pos  // position of "switch" keyword
			Init   Stmt       // initialization statement; or nil
			Assign Stmt       // x := y.(type) or y.(type)
			Body   *BlockStmt // CaseClauses only
		}
	*/
	"TypeSwitchStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Switch",
			Token:         Basic{jen.Qual("go/token", "SWITCH")},
			PositionField: Field{"Switch"},
		},
		Decoration{
			Name: "AfterSwitch",
		},
		Node{
			Name:  "Init",
			Field: Field{"Init"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Init").Op("!=").Nil() }),
		},
		Node{
			Name:  "Assign",
			Field: Field{"Assign"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterAssign",
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A CommClause node represents a case of a select statement.
		CommClause struct {
			Case  token.Pos // position of "case" or "default" keyword
			Comm  Stmt      // send or receive statement; nil means default case
			Colon token.Pos // position of ":"
			Body  []Stmt    // statement list; or nil
		}
	*/
	"CommClause": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name: "Case",
			Token: Expr(func(n *jen.Statement) *jen.Statement {
				return jen.Func().Params().Qual("go/token", "Token").Block(jen.If(n.Dot("Comm").Op("==").Nil()).Block(jen.Return(jen.Qual("go/token", "DEFAULT"))).Else().Block(jen.Return(jen.Qual("go/token", "CASE")))).Call()
			}),
			PositionField: Field{"Case"},
		},
		Decoration{
			Name: "AfterCase",
		},
		Node{
			Name:  "Comm",
			Field: Field{"Comm"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterComm",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Comm").Op("!=").Nil() }),
		},
		Token{
			Name:          "Colon",
			Token:         Basic{jen.Qual("go/token", "COLON")},
			PositionField: Field{"Colon"},
		},
		Decoration{
			Name: "AfterColon",
		},
		List{
			Name:      "Body",
			Field:     Field{"Body"},
			Elem:      Type{"Stmt", false},
			Separator: token.SEMICOLON,
		},
		// Never want to attach decorations to the end of a list of statements - always better to
		// attach to the last statement.
	},
	/*
		// An SelectStmt node represents a select statement.
		SelectStmt struct {
			Select token.Pos  // position of "select" keyword
			Body   *BlockStmt // CommClauses only
		}
	*/
	"SelectStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Select",
			Token:         Basic{jen.Qual("go/token", "SELECT")},
			PositionField: Field{"Select"},
		},
		Decoration{
			Name: "AfterSelect",
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A ForStmt represents a for statement.
		ForStmt struct {
			For  token.Pos // position of "for" keyword
			Init Stmt      // initialization statement; or nil
			Cond Expr      // condition; or nil
			Post Stmt      // post iteration statement; or nil
			Body *BlockStmt
		}
	*/
	"ForStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "For",
			Token:         Basic{jen.Qual("go/token", "FOR")},
			PositionField: Field{"For"},
		},
		Decoration{
			Name: "AfterFor",
		},
		Node{
			Name:  "Init",
			Field: Field{"Init"},
			Type:  Type{"Stmt", false},
		},
		Token{
			Name:   "InitSemicolon",
			Token:  Basic{jen.Qual("go/token", "SEMICOLON")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Init").Op("!=").Nil() }),
		},
		Decoration{
			Name: "AfterInit",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Init").Op("!=").Nil() }),
		},
		Node{
			Name:  "Cond",
			Field: Field{"Cond"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:   "CondSemicolon",
			Token:  Basic{jen.Qual("go/token", "SEMICOLON")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Post").Op("!=").Nil() }),
		},
		Decoration{
			Name: "AfterCond",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Cond").Op("!=").Nil() }),
		},
		Node{
			Name:  "Post",
			Field: Field{"Post"},
			Type:  Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterPost",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Post").Op("!=").Nil() }),
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A RangeStmt represents a for statement with a range clause.
		RangeStmt struct {
			For        token.Pos   // position of "for" keyword
			Key, Value Expr        // Key, Value may be nil
			TokPos     token.Pos   // position of Tok; invalid if Key == nil
			Tok        token.Token // ILLEGAL if Key == nil, ASSIGN, DEFINE
			X          Expr        // value to range over
			Body       *BlockStmt
		}
	*/
	"RangeStmt": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "For",
			Token:         Basic{jen.Qual("go/token", "FOR")},
			PositionField: Field{"For"},
		},
		Decoration{
			Name: "AfterFor",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Key").Op("!=").Nil() }),
		},
		Node{
			Name:  "Key",
			Field: Field{"Key"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:   "Comma",
			Token:  Basic{jen.Qual("go/token", "COMMA")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Value").Op("!=").Nil() }),
		},
		Decoration{
			Name: "AfterKey",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Key").Op("!=").Nil() }),
		},
		Node{
			Name:  "Value",
			Field: Field{"Value"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterValue",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Value").Op("!=").Nil() }),
		},
		Token{
			Name:          "Tok",
			TokenField:    Field{"Tok"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tok") }),
			Exists:        Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tok").Op("!=").Qual("go/token", "ILLEGAL") }),
			PositionField: Field{"TokPos"},
		},
		Token{
			Name:  "Range",
			Token: Basic{jen.Qual("go/token", "RANGE")},
		},
		Decoration{
			Name: "AfterRange",
		},
		Node{
			Name:  "X",
			Field: Field{"X"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// An ImportSpec node represents a single package import.
		ImportSpec struct {
			Doc     *CommentGroup // associated documentation; or nil
			Name    *Ident        // local package name (including "."); or nil
			Path    *BasicLit     // import path
			Comment *CommentGroup // line comments; or nil
			EndPos  token.Pos     // end of spec (overrides Path.Pos if nonzero)
		}
	*/
	// TODO: Do we need EndPos? I think it's a kludge to ensure comments don't move around after re-writing imports, so we should be able to ignore it?
	"ImportSpec": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Name",
			Field: Field{"Name"},
			Type:  Type{"Ident", true},
		},
		Decoration{
			Name: "AfterName",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Name").Op("!=").Nil() }),
		},
		Node{
			Name:  "Path",
			Field: Field{"Path"},
			Type:  Type{"BasicLit", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A ValueSpec node represents a constant or variable declaration
		// (ConstSpec or VarSpec production).
		//
		ValueSpec struct {
			Doc     *CommentGroup // associated documentation; or nil
			Names   []*Ident      // value names (len(Names) > 0)
			Type    Expr          // value type; or nil
			Values  []Expr        // initial values; or nil
			Comment *CommentGroup // line comments; or nil
		}
	*/
	"ValueSpec": {
		Decoration{
			Name: "Start",
		},
		List{
			Name:      "Names",
			Field:     Field{"Names"},
			Elem:      Type{"Ident", true},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterNames",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Type").Op("!=").Nil() }),
		},
		Node{
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"Expr", false},
		},
		Token{
			Name:   "Assign",
			Token:  Basic{jen.Qual("go/token", "ASSIGN")},
			Exists: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Values").Op("!=").Nil() }),
		},
		Decoration{
			Name: "AfterAssign",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Values").Op("!=").Nil() }),
		},
		List{
			Name:      "Values",
			Field:     Field{"Values"},
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A TypeSpec node represents a type declaration (TypeSpec production).
		TypeSpec struct {
			Doc     *CommentGroup // associated documentation; or nil
			Name    *Ident        // type name
			Assign  token.Pos     // position of '=', if any
			Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
			Comment *CommentGroup // line comments; or nil
		}
	*/
	"TypeSpec": {
		Decoration{
			Name: "Start",
		},
		Node{
			Name:  "Name",
			Field: Field{"Name"},
			Type:  Type{"Ident", true},
		},
		Token{
			Name:  "Assign",
			Token: Basic{jen.Qual("go/token", "ASSIGN")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Assign").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Assign") }),
			},
			ExistsField:   Field{"Assign"},
			PositionField: Field{"Assign"},
		},
		Decoration{
			Name: "AfterName",
		},
		Node{
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A BadDecl node is a placeholder for declarations containing
		// syntax errors for which no correct declaration nodes can be
		// created.
		//
		BadDecl struct {
			From, To token.Pos // position range of bad declaration
		}
	*/
	"BadDecl": {
		Ignored{
			Length: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement {
					return jen.Int().Parens(jen.Add(n).Dot("End").Call().Op("-").Add(n).Dot("Pos").Call())
				}),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Length") }),
			},
			PositionField: Field{"From"},
		},
	},
	/*
		// A GenDecl node (generic declaration node) represents an import,
		// constant, type or variable declaration. A valid Lparen position
		// (Lparen.IsValid()) indicates a parenthesized declaration.
		//
		// Relationship between Tok value and Specs element type:
		//
		//	token.IMPORT  *ImportSpec
		//	token.CONST   *ValueSpec
		//	token.TYPE    *TypeSpec
		//	token.VAR     *ValueSpec
		//
		GenDecl struct {
			Doc    *CommentGroup // associated documentation; or nil
			TokPos token.Pos     // position of Tok
			Tok    token.Token   // IMPORT, CONST, TYPE, VAR
			Lparen token.Pos     // position of '(', if any
			Specs  []Spec
			Rparen token.Pos // position of ')', if any
		}
	*/
	"GenDecl": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Tok",
			TokenField:    Field{"Tok"},
			Token:         Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Tok") }),
			PositionField: Field{"TokPos"},
		},
		Decoration{
			Name: "AfterTok",
		},
		Token{
			Name:  "Lparen",
			Token: Basic{jen.Qual("go/token", "LPAREN")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Lparen").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Lparen") }),
			},
			ExistsField:   Field{"Lparen"},
			PositionField: Field{"Lparen"},
		},
		Decoration{
			Name: "AfterLparen",
			Use: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Lparen").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Lparen") }),
			},
		},
		List{
			Name:      "Specs",
			Field:     Field{"Specs"},
			Elem:      Type{"Spec", false},
			Separator: token.SEMICOLON,
		},
		Token{
			Name:  "Rparen",
			Token: Basic{jen.Qual("go/token", "RPAREN")},
			Exists: Double{
				Ast: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Rparen").Dot("IsValid").Call() }),
				Dst: Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Rparen") }),
			},
			ExistsField:   Field{"Rparen"},
			PositionField: Field{"Rparen"},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A FuncDecl node represents a function declaration.
		FuncDecl struct {
			Doc  *CommentGroup // associated documentation; or nil
			Recv *FieldList    // receiver (methods); or nil (functions)
			Name *Ident        // function/method name
			Type *FuncType     // function signature: parameters, results, and position of "func" keyword
			Body *BlockStmt    // function body; or nil for external (non-Go) function
		}

		// A FuncType node represents a function type.
		FuncType struct {
			Func    token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
			Params  *FieldList // (incoming) parameters; non-nil
			Results *FieldList // (outgoing) results; or nil
		}
	*/
	"FuncDecl": {
		Init{
			// Initializer for "Type"
			Name:  "Type",
			Field: Field{"Type"},
			Type:  Type{"FuncType", true},
		},
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Func",
			Token:         Basic{jen.Qual("go/token", "FUNC")},
			Exists:        Basic{jen.Lit(true)},
			ExistsField:   InnerField{"Type", "Func"},
			PositionField: InnerField{"Type", "Func"},
		},
		Decoration{
			Name: "AfterFunc",
		},
		// TODO: render any decorations from n.Type.AfterFunc (but never save them there)
		Node{
			Name:  "Recv",
			Field: Field{"Recv"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "AfterRecv",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Recv").Op("!=").Nil() }),
		},
		Node{
			Name:  "Name",
			Field: Field{"Name"},
			Type:  Type{"Ident", true},
		},
		Decoration{
			Name: "AfterName",
		},
		Node{
			Name:  "Params",
			Field: InnerField{"Type", "Params"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "AfterParams",
		},
		// TODO: render any decorations from n.Type.AfterParams (but never save them there)
		Node{
			Name:  "Results",
			Field: InnerField{"Type", "Results"},
			Type:  Type{"FieldList", true},
		},
		Decoration{
			Name: "AfterResults",
			Use:  Expr(func(n *jen.Statement) *jen.Statement { return n.Dot("Type").Dot("Results").Op("!=").Nil() }),
		},
		// TODO: render any decorations from n.Type.AfterResults (but never save them there)
		Node{
			Name:  "Body",
			Field: Field{"Body"},
			Type:  Type{"BlockStmt", true},
		},
		Decoration{
			Name: "End",
		},
	},
	/*
		// A File node represents a Go source file.
		//
		// The Comments list contains all comments in the source file in order of
		// appearance, including the comments that are pointed to from other nodes
		// via Doc and Comment fields.
		//
		// For correct printing of source code containing comments (using packages
		// go/format and go/printer), special care must be taken to update comments
		// when a File's syntax tree is modified: For printing, comments are interspersed
		// between tokens based on their position. If syntax tree nodes are
		// removed or moved, relevant comments in their vicinity must also be removed
		// (from the File.Comments list) or moved accordingly (by updating their
		// positions). A CommentMap may be used to facilitate some of these operations.
		//
		// Whether and how a comment is associated with a node depends on the
		// interpretation of the syntax tree by the manipulating program: Except for Doc
		// and Comment comments directly associated with nodes, the remaining comments
		// are "free-floating" (see also issues #18593, #20744).
		//
		type File struct {
			Doc        *CommentGroup   // associated documentation; or nil
			Package    token.Pos       // position of "package" keyword
			Name       *Ident          // package name
			Decls      []Decl          // top-level declarations; or nil
			Scope      *Scope          // package scope (this file only)
			Imports    []*ImportSpec   // imports in this file
			Unresolved []*Ident        // unresolved identifiers in this file
			Comments   []*CommentGroup // list of all comments in the source file
		}
	*/
	// TODO: File.Scope?
	// TODO: File.Imports?
	// TODO: File.Unresolved?
	"File": {
		Decoration{
			Name: "Start",
		},
		Token{
			Name:          "Package",
			Token:         Basic{jen.Qual("go/token", "PACKAGE")},
			PositionField: Field{"Package"},
		},
		Decoration{
			Name: "AfterPackage",
		},
		Node{
			Name:  "Name",
			Field: Field{"Name"},
			Type:  Type{"Ident", true},
		},
		Decoration{
			Name: "AfterName",
		},
		List{
			Name:      "Decls",
			Field:     Field{"Decls"},
			Elem:      Type{"Decl", false},
			Separator: token.SEMICOLON,
		},
		// Never want to attach decorations to the end of a list of declarations - always better to
		// attach to the last statement.
	},
}

type Init struct {
	Name  string
	Field FieldSpec
	Type  Type
}

type Decoration struct {
	Name string
	Use  Code
}

type String struct {
	Name          string
	ValueField    FieldSpec
	PositionField FieldSpec
}

type List struct {
	Name      string
	Field     FieldSpec
	Elem      Type
	Separator token.Token
}

type Node struct {
	Name  string
	Field FieldSpec
	Type  Type
}

type Token struct {
	Name          string
	Exists        Code
	Token         Code
	ExistsField   FieldSpec
	PositionField FieldSpec
	TokenField    FieldSpec
}

type Ignored struct {
	Length        Code
	PositionField FieldSpec
}

// Value that must be copied from ast.Node to dst.Node but doesn't result in anything rendered to the output.
type Value struct {
	Name  string
	Field FieldSpec
	Value Code
}

type Code interface {
	Get(id string, ast bool) *jen.Statement
}

type Basic struct {
	*jen.Statement
}

func (b Basic) Get(id string, ast bool) *jen.Statement {
	return b.Statement
}

type Expr func(n *jen.Statement) *jen.Statement

func (e Expr) Get(id string, ast bool) *jen.Statement {
	return e(jen.Id(id))
}

type Double struct {
	Ast Expr
	Dst Expr
}

func (d Double) Get(id string, ast bool) *jen.Statement {
	if ast {
		return d.Ast(jen.Id(id))
	}
	return d.Dst(jen.Id(id))
}

type Type struct {
	Name    string
	Pointer bool
}

func (t Type) Literal(path string) *jen.Statement {
	if t.Pointer {
		return jen.Op("*").Qual(path, t.Name)
	}
	return jen.Qual(path, t.Name)
}

func (t Type) String() string {
	if t.Pointer {
		return "*" + t.Name
	}
	return t.Name
}

type FieldSpec interface {
	Get(id string) *jen.Statement
}

type Field struct {
	Name string
}

func (f Field) Get(id string) *jen.Statement {
	return jen.Id(id).Dot(f.Name)
}

type InnerField struct {
	Inner, Name string
}

func (f InnerField) Get(id string) *jen.Statement {
	return jen.Id(id).Dot(f.Inner).Dot(f.Name)
}
