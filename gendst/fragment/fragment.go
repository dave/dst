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
			Elem:      Type{"Ident", true},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterNames",
		},
		Node{
			Name: "Type",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterType",
			Use:  Single{jen.Id("n").Dot("Tag").Op("!=").Nil()},
		},
		Node{
			Name: "Tag",
			Type: Type{"BasicLit", true},
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
		Token{
			Name:  "Opening",
			Token: Single{jen.Qual("go/token", "LPAREN")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Opening").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Opening"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Opening"),
				Set: jen.Id("n").Dot("Opening").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterOpening",
		},
		List{
			Name:      "List",
			Elem:      Type{"Field", true},
			Separator: token.COMMA,
		},
		Token{
			Name:  "Closing",
			Token: Single{jen.Qual("go/token", "RPAREN")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Closing").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Closing"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Closing"),
				Set: jen.Id("n").Dot("Closing").Op("=").Id("pos"),
			},
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
				Ast: jen.Int().Parens(jen.Id("n").Dot("End").Call().Op("-").Id("n").Dot("Pos").Call()),
				Dst: jen.Id("n").Dot("Length"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("From"),
				Set: jen.Id("n").Dot("From").Op("=").Id("pos"),
			},
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
			Name: "Name",
			Position: Position{
				Get: jen.Id("n").Dot("NamePos"),
				Set: jen.Id("n").Dot("NamePos").Op("=").Id("pos"),
			},
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
			Name:  "Ellipsis",
			Token: Single{jen.Qual("go/token", "ELLIPSIS")},
			Position: Position{
				Get: jen.Id("n").Dot("Ellipsis"),
				Set: jen.Id("n").Dot("Ellipsis").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterEllipsis",
			Use:  Single{jen.Id("n").Dot("Elt").Op("!=").Nil()},
		},
		Node{
			Name: "Elt",
			Type: Type{"Expr", false},
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
			Name: "Value",
			Position: Position{
				Get: jen.Id("n").Dot("ValuePos"),
				Set: jen.Id("n").Dot("ValuePos").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name: "Kind",
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
			Name: "Type",
			Type: Type{"FuncType", true},
		},
		Decoration{
			Name: "AfterType",
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name: "Type",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterType",
			Use:  Single{jen.Id("n").Dot("Type").Op("!=").Nil()},
		},
		Token{
			Name:  "Lbrace",
			Token: Single{jen.Qual("go/token", "LBRACE")},
			Position: Position{
				Get: jen.Id("n").Dot("Lbrace"),
				Set: jen.Id("n").Dot("Lbrace").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLbrace",
		},
		List{
			Name:      "Elts",
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterElts",
		},
		Token{
			Name:  "Rbrace",
			Token: Single{jen.Qual("go/token", "RBRACE")},
			Position: Position{
				Get: jen.Id("n").Dot("Rbrace"),
				Set: jen.Id("n").Dot("Rbrace").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name: "Incomplete",
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
			Name:  "Lparen",
			Token: Single{jen.Qual("go/token", "LPAREN")},
			Position: Position{
				Get: jen.Id("n").Dot("Lparen"),
				Set: jen.Id("n").Dot("Lparen").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLparen",
		},
		Node{
			Name: "X",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:  "Rparen",
			Token: Single{jen.Qual("go/token", "RPAREN")},
			Position: Position{
				Get: jen.Id("n").Dot("Rparen"),
				Set: jen.Id("n").Dot("Rparen").Op("=").Id("pos"),
			},
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
			Name: "X",
			Type: Type{"Expr", false},
		},
		Token{
			Name:  "Period",
			Token: Single{jen.Qual("go/token", "PERIOD")},
		},
		Decoration{
			Name: "AfterX",
		},
		Node{
			Name: "Sel",
			Type: Type{"Ident", true},
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
			Name: "X",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:  "Lbrack",
			Token: Single{jen.Qual("go/token", "LBRACK")},
			Position: Position{
				Get: jen.Id("n").Dot("Lbrack"),
				Set: jen.Id("n").Dot("Lbrack").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLbrack",
		},
		Node{
			Name: "Index",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterIndex",
		},
		Token{
			Name:  "Rbrack",
			Token: Single{jen.Qual("go/token", "RBRACK")},
			Position: Position{
				Get: jen.Id("n").Dot("Rbrack"),
				Set: jen.Id("n").Dot("Rbrack").Op("=").Id("pos"),
			},
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
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:  "Lbrack",
			Token: Single{jen.Qual("go/token", "LBRACK")},
			Position: Position{
				Get: jen.Id("n").Dot("Lbrack"),
				Set: jen.Id("n").Dot("Lbrack").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLbrack",
			Use:  Single{jen.Id("n").Dot("Low").Op("!=").Nil()},
		},
		Node{
			Name: "Low",
			Type: Type{"Expr", false},
		},
		Token{
			Name:  "Colon1",
			Token: Single{jen.Qual("go/token", "COLON")},
		},
		Decoration{
			Name: "AfterLow",
		},
		Node{
			Name: "High",
			Type: Type{"Expr", false},
		},
		Token{
			Name:   "Colon2",
			Token:  Single{jen.Qual("go/token", "COLON")},
			Exists: Single{jen.Id("n").Dot("Slice3")},
		},
		Decoration{
			Name: "AfterHigh",
			Use:  Single{jen.Id("n").Dot("High").Op("!=").Nil()},
		},
		Node{
			Name: "Max",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterMax",
			Use:  Single{jen.Id("n").Dot("Max").Op("!=").Nil()}, // TODO - Slice3 in here?
		},
		Token{
			Name:  "Rbrack",
			Token: Single{jen.Qual("go/token", "RBRACK")},
			Position: Position{
				Get: jen.Id("n").Dot("Rbrack"),
				Set: jen.Id("n").Dot("Rbrack").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name: "Slice3",
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
			Name: "X",
			Type: Type{"Expr", false},
		},
		Token{
			Name:  "Period",
			Token: Single{jen.Qual("go/token", "PERIOD")},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:  "Lparen",
			Token: Single{jen.Qual("go/token", "LPAREN")},
			Position: Position{
				Get: jen.Id("n").Dot("Lparen"),
				Set: jen.Id("n").Dot("Lparen").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLparen",
		},
		Node{
			Name: "Type",
			Type: Type{"Expr", false},
		},
		Token{
			Name:   "TypeToken",
			Token:  Single{jen.Qual("go/token", "TYPE")},
			Exists: Single{jen.Id("n").Dot("Type").Op("==").Nil()},
		},
		Decoration{
			Name: "AfterType",
		},
		Token{
			Name:  "Rparen",
			Token: Single{jen.Qual("go/token", "RPAREN")},
			Position: Position{
				Get: jen.Id("n").Dot("Rparen"),
				Set: jen.Id("n").Dot("Rparen").Op("=").Id("pos"),
			},
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
			Name: "Fun",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterFun",
		},
		Token{
			Name:  "Lparen",
			Token: Single{jen.Qual("go/token", "LPAREN")},
			Position: Position{
				Get: jen.Id("n").Dot("Lparen"),
				Set: jen.Id("n").Dot("Lparen").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLparen",
		},
		List{
			Name:      "Args",
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterArgs",
		},
		Token{
			Name:  "Ellipsis",
			Token: Single{jen.Qual("go/token", "ELLIPSIS")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Ellipsis").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Ellipsis"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Ellipsis"),
				Set: jen.Id("n").Dot("Ellipsis").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterEllipsis",
			Use: Double{
				Ast: jen.Id("n").Dot("Ellipsis").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Ellipsis"),
			},
		},
		Token{
			Name:  "Rparen",
			Token: Single{jen.Qual("go/token", "RPAREN")},
			Position: Position{
				Get: jen.Id("n").Dot("Rparen"),
				Set: jen.Id("n").Dot("Rparen").Op("=").Id("pos"),
			},
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
			Name:  "Star",
			Token: Single{jen.Qual("go/token", "MUL")},
			Position: Position{
				Get: jen.Id("n").Dot("Star"),
				Set: jen.Id("n").Dot("Star").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterStar",
		},
		Node{
			Name: "X",
			Type: Type{"Expr", false},
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
			Name:  "Op",
			Token: Single{jen.Id("n").Dot("Op")},
			Position: Position{
				Get: jen.Id("n").Dot("OpPos"),
				Set: jen.Id("n").Dot("OpPos").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterOp",
		},
		Node{
			Name: "X",
			Type: Type{"Expr", false},
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
			Name: "X",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:  "Op",
			Token: Single{jen.Id("n").Dot("Op")},
			Position: Position{
				Get: jen.Id("n").Dot("OpPos"),
				Set: jen.Id("n").Dot("OpPos").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterOp",
		},
		Node{
			Name: "Y",
			Type: Type{"Expr", false},
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
			Name: "Key",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterKey",
		},
		Token{
			Name:  "Colon",
			Token: Single{jen.Qual("go/token", "COLON")},
			Position: Position{
				Get: jen.Id("n").Dot("Colon"),
				Set: jen.Id("n").Dot("Colon").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterColon",
		},
		Node{
			Name: "Value",
			Type: Type{"Expr", false},
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
			Name:  "Lbrack",
			Token: Single{jen.Qual("go/token", "LBRACK")},
			Position: Position{
				Get: jen.Id("n").Dot("Lbrack"),
				Set: jen.Id("n").Dot("Lbrack").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLbrack",
		},
		Node{
			Name: "Len",
			Type: Type{"Expr", false},
		},
		Token{
			Name:  "Rbrack",
			Token: Single{jen.Qual("go/token", "RBRACK")},
		},
		Decoration{
			Name: "AfterLen",
		},
		Node{
			Name: "Elt",
			Type: Type{"Expr", false},
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
			Name:  "Struct",
			Token: Single{jen.Qual("go/token", "STRUCT")},
			Position: Position{
				Get: jen.Id("n").Dot("Struct"),
				Set: jen.Id("n").Dot("Struct").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterStruct",
		},
		Node{
			Name: "Fields",
			Type: Type{"FieldList", true},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name: "Incomplete", // TODO: Remove this and always set to false?
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
			Token: Single{jen.Qual("go/token", "FUNC")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Func").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Func"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Func"),
				Set: jen.Id("n").Dot("Func").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterFunc",
			Use: Double{
				Ast: jen.Id("n").Dot("Func").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Func"),
			},
		},
		Node{
			Name: "Params",
			Type: Type{"FieldList", true},
		},
		Decoration{
			Name: "AfterParams",
			Use:  Single{jen.Id("n").Dot("Results").Op("!=").Nil()},
		},
		Node{
			Name: "Results",
			Type: Type{"FieldList", true},
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
			Name:  "Interface",
			Token: Single{jen.Qual("go/token", "INTERFACE")},
			Position: Position{
				Get: jen.Id("n").Dot("Interface"),
				Set: jen.Id("n").Dot("Interface").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterInterface",
		},
		Node{
			Name: "Methods",
			Type: Type{"FieldList", true},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name: "Incomplete", // TODO: Remove this and always set to false?
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
			Name:  "Map",
			Token: Single{jen.Qual("go/token", "MAP")},
			Position: Position{
				Get: jen.Id("n").Dot("Map"),
				Set: jen.Id("n").Dot("Map").Op("=").Id("pos"),
			},
		},
		Token{
			Name:  "Lbrack",
			Token: Single{jen.Qual("go/token", "LBRACK")},
		},
		Decoration{
			Name: "AfterMap",
		},
		Node{
			Name: "Key",
			Type: Type{"Expr", false},
		},
		Token{
			Name:  "Rbrack",
			Token: Single{jen.Qual("go/token", "RBRACK")},
		},
		Decoration{
			Name: "AfterKey",
		},
		Node{
			Name: "Value",
			Type: Type{"Expr", false},
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
		// This is rather a kludge. In SEND variation, we emit "<-" followed by "chan". Otherwise we
		// just emit "chan".
		Token{
			Name: "Begin",
			Token: Double{
				Ast: jen.Func().Params().Qual("go/token", "Token").Block(jen.If(jen.Id("n").Dot("Dir").Op("==").Qual("go/ast", "SEND")).Block(jen.Return(jen.Qual("go/token", "ARROW"))).Else().Block(jen.Return(jen.Qual("go/token", "CHAN")))).Call(),
				Dst: jen.Func().Params().Qual("go/token", "Token").Block(jen.If(jen.Id("n").Dot("Dir").Op("==").Qual(DSTPATH, "SEND")).Block(jen.Return(jen.Qual("go/token", "ARROW"))).Else().Block(jen.Return(jen.Qual("go/token", "CHAN")))).Call(),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Begin"),
				Set: jen.Id("n").Dot("Begin").Op("=").Id("pos"),
			},
		},
		Token{
			Name:  "Chan",
			Token: Single{jen.Qual("go/token", "CHAN")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Dir").Op("==").Qual("go/ast", "SEND"),
				Dst: jen.Id("n").Dot("Dir").Op("==").Qual(DSTPATH, "SEND"),
			},
		},
		Decoration{
			Name: "AfterBegin",
		},
		Token{
			Name:  "Arrow",
			Token: Single{jen.Qual("go/token", "ARROW")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Dir").Op("==").Qual("go/ast", "RECV"),
				Dst: jen.Id("n").Dot("Dir").Op("==").Qual(DSTPATH, "RECV"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Arrow"),
				Set: jen.Id("n").Dot("Arrow").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterArrow",
			Use: Double{
				Ast: jen.Id("n").Dot("Dir").Op("==").Qual("go/ast", "RECV"),
				Dst: jen.Id("n").Dot("Dir").Op("==").Qual(DSTPATH, "RECV"),
			},
		},
		Node{
			Name: "Value",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "End",
		},
		Value{
			Name: "Dir",
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
				Ast: jen.Int().Parens(jen.Id("n").Dot("End").Call().Op("-").Id("n").Dot("Pos").Call()),
				Dst: jen.Id("n").Dot("Length"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("From"),
				Set: jen.Id("n").Dot("From").Op("=").Id("pos"),
			},
		},
	},
	/*
		// A DeclStmt node represents a declaration in a statement list.
		DeclStmt struct {
			Decl Decl // *GenDecl with CONST, TYPE, or VAR token
		}
	*/
	"DeclStmt": {
		Node{
			Name: "Decl",
			Type: Type{"Decl", false},
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
			Name:   "Semicolon",
			Token:  Single{jen.Qual("go/token", "ARROW")},
			Exists: Single{jen.Op("!").Id("n").Dot("Implicit")},
			Position: Position{
				Get: jen.Id("n").Dot("Semicolon"),
				Set: jen.Id("n").Dot("Semicolon").Op("=").Id("pos"),
			},
		},
		Value{
			Name: "Implicit",
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
			Name: "Label",
			Type: Type{"Ident", true},
		},
		Decoration{
			Name: "AfterLabel",
		},
		Token{
			Name:  "Colon",
			Token: Single{jen.Qual("go/token", "COLON")},
			Position: Position{
				Get: jen.Id("n").Dot("Colon"),
				Set: jen.Id("n").Dot("Colon").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterColon",
		},
		Node{
			Name: "Stmt",
			Type: Type{"Stmt", false},
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
		Node{
			Name: "X",
			Type: Type{"Expr", false},
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
			Name: "Chan",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterChan",
		},
		Token{
			Name:  "Arrow",
			Token: Single{jen.Qual("go/token", "ARROW")},
			Position: Position{
				Get: jen.Id("n").Dot("Arrow"),
				Set: jen.Id("n").Dot("Arrow").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterArrow",
		},
		Node{
			Name: "Value",
			Type: Type{"Expr", false},
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
			Name: "X",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Token{
			Name:  "Tok",
			Token: Single{jen.Id("n").Dot("Tok")},
			Position: Position{
				Get: jen.Id("n").Dot("TokPos"),
				Set: jen.Id("n").Dot("TokPos").Op("=").Id("pos"),
			},
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
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterLhs",
		},
		Token{
			Name:  "Tok",
			Token: Single{jen.Id("n").Dot("Tok")},
			Position: Position{
				Get: jen.Id("n").Dot("TokPos"),
				Set: jen.Id("n").Dot("TokPos").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterTok",
		},
		List{
			Name:      "Rhs",
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
			Name:  "Go",
			Token: Single{jen.Qual("go/token", "GO")},
			Position: Position{
				Get: jen.Id("n").Dot("Go"),
				Set: jen.Id("n").Dot("Go").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterGo",
		},
		Node{
			Name: "Call",
			Type: Type{"CallExpr", true},
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
			Name:  "Defer",
			Token: Single{jen.Qual("go/token", "DEFER")},
			Position: Position{
				Get: jen.Id("n").Dot("Defer"),
				Set: jen.Id("n").Dot("Defer").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterDefer",
		},
		Node{
			Name: "Call",
			Type: Type{"CallExpr", true},
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
			Name:  "Return",
			Token: Single{jen.Qual("go/token", "RETURN")},
			Position: Position{
				Get: jen.Id("n").Dot("Return"),
				Set: jen.Id("n").Dot("Return").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterReturn",
		},
		List{
			Name:      "Results",
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
			Name:  "Tok",
			Token: Single{jen.Id("n").Dot("Tok")},
			Position: Position{
				Get: jen.Id("n").Dot("TokPos"),
				Set: jen.Id("n").Dot("TokPos").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterTok",
			Use:  Single{jen.Id("n").Dot("Label").Op("!=").Nil()},
		},
		Node{
			Name: "Label",
			Type: Type{"Ident", true},
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
			Name:  "Lbrace",
			Token: Single{jen.Qual("go/token", "LBRACE")},
			Position: Position{
				Get: jen.Id("n").Dot("Lbrace"),
				Set: jen.Id("n").Dot("Lbrace").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLbrace",
		},
		List{
			Name:      "List",
			Elem:      Type{"Stmt", false},
			Separator: token.SEMICOLON,
		},
		Token{
			Name:  "Rbrace",
			Token: Single{jen.Qual("go/token", "RBRACE")},
			Position: Position{
				Get: jen.Id("n").Dot("Rbrace"),
				Set: jen.Id("n").Dot("Rbrace").Op("=").Id("pos"),
			},
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
			Name:  "If",
			Token: Single{jen.Qual("go/token", "IF")},
			Position: Position{
				Get: jen.Id("n").Dot("If"),
				Set: jen.Id("n").Dot("If").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterIf",
		},
		Node{
			Name: "Init",
			Type: Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Single{jen.Id("n").Dot("Init").Op("!=").Nil()},
		},
		Node{
			Name: "Cond",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterCond",
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
		},
		Token{
			Name:   "ElseTok",
			Token:  Single{jen.Qual("go/token", "ELSE")},
			Exists: Single{jen.Id("n").Dot("Else").Op("!=").Nil()},
		},
		Decoration{
			Name: "AfterElse",
			Use:  Single{jen.Id("n").Dot("Else").Op("!=").Nil()},
		},
		Node{
			Name: "Else",
			Type: Type{"Stmt", false},
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
			Name:  "Case",
			Token: Single{jen.Func().Params().Qual("go/token", "Token").Block(jen.If(jen.Id("n").Dot("List").Op("==").Nil()).Block(jen.Return(jen.Qual("go/token", "DEFAULT"))).Else().Block(jen.Return(jen.Qual("go/token", "CASE")))).Call()},
			Position: Position{
				Get: jen.Id("n").Dot("Case"),
				Set: jen.Id("n").Dot("Case").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterCase",
		},
		List{
			Name:      "List",
			Elem:      Type{"Expr", false},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterList",
			Use:  Single{jen.Id("n").Dot("List").Op("!=").Nil()},
		},
		Token{
			Name:  "Colon",
			Token: Single{jen.Qual("go/token", "COLON")},
			Position: Position{
				Get: jen.Id("n").Dot("Colon"),
				Set: jen.Id("n").Dot("Colon").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterColon",
		},
		List{
			Name:      "Body",
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
			Name:  "Switch",
			Token: Single{jen.Qual("go/token", "SWITCH")},
			Position: Position{
				Get: jen.Id("n").Dot("Switch"),
				Set: jen.Id("n").Dot("Switch").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterSwitch",
		},
		Node{
			Name: "Init",
			Type: Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Single{jen.Id("n").Dot("Init").Op("!=").Nil()},
		},
		Node{
			Name: "Tag",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterTag",
			Use:  Single{jen.Id("n").Dot("Tag").Op("!=").Nil()},
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name:  "Switch",
			Token: Single{jen.Qual("go/token", "SWITCH")},
			Position: Position{
				Get: jen.Id("n").Dot("Switch"),
				Set: jen.Id("n").Dot("Switch").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterSwitch",
		},
		Node{
			Name: "Init",
			Type: Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Single{jen.Id("n").Dot("Init").Op("!=").Nil()},
		},
		Node{
			Name: "Assign",
			Type: Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterAssign",
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name:  "Case",
			Token: Single{jen.Func().Params().Qual("go/token", "Token").Block(jen.If(jen.Id("n").Dot("Comm").Op("==").Nil()).Block(jen.Return(jen.Qual("go/token", "DEFAULT"))).Else().Block(jen.Return(jen.Qual("go/token", "CASE")))).Call()},
			Position: Position{
				Get: jen.Id("n").Dot("Case"),
				Set: jen.Id("n").Dot("Case").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterCase",
		},
		Node{
			Name: "Comm",
			Type: Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterComm",
			Use:  Single{jen.Id("n").Dot("Comm").Op("!=").Nil()},
		},
		Token{
			Name:  "Colon",
			Token: Single{jen.Qual("go/token", "COLON")},
			Position: Position{
				Get: jen.Id("n").Dot("Colon"),
				Set: jen.Id("n").Dot("Colon").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterColon",
		},
		List{
			Name:      "Body",
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
			Name:  "Select",
			Token: Single{jen.Qual("go/token", "SELECT")},
			Position: Position{
				Get: jen.Id("n").Dot("Select"),
				Set: jen.Id("n").Dot("Select").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterSelect",
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name:  "For",
			Token: Single{jen.Qual("go/token", "FOR")},
			Position: Position{
				Get: jen.Id("n").Dot("For"),
				Set: jen.Id("n").Dot("For").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterFor",
		},
		Node{
			Name: "Init",
			Type: Type{"Stmt", false},
		},
		Token{
			Name:   "InitSemicolon",
			Token:  Single{jen.Qual("go/token", "SEMICOLON")},
			Exists: Single{jen.Id("n").Dot("Init").Op("!=").Nil()},
		},
		Decoration{
			Name: "AfterInit",
			Use:  Single{jen.Id("n").Dot("Init").Op("!=").Nil()},
		},
		Node{
			Name: "Cond",
			Type: Type{"Expr", false},
		},
		Token{
			Name:   "CondSemicolon",
			Token:  Single{jen.Qual("go/token", "SEMICOLON")},
			Exists: Single{jen.Id("n").Dot("Post").Op("!=").Nil()},
		},
		Decoration{
			Name: "AfterCond",
			Use:  Single{jen.Id("n").Dot("Cond").Op("!=").Nil()},
		},
		Node{
			Name: "Post",
			Type: Type{"Stmt", false},
		},
		Decoration{
			Name: "AfterPost",
			Use:  Single{jen.Id("n").Dot("Post").Op("!=").Nil()},
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name:  "For",
			Token: Single{jen.Qual("go/token", "FOR")},
			Position: Position{
				Get: jen.Id("n").Dot("For"),
				Set: jen.Id("n").Dot("For").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterFor",
			Use:  Single{jen.Id("n").Dot("Key").Op("!=").Nil()},
		},
		Node{
			Name: "Key",
			Type: Type{"Expr", false},
		},
		Token{
			Name:   "Comma",
			Exists: Single{jen.Id("n").Dot("Value").Op("!=").Nil()},
		},
		Decoration{
			Name: "AfterKey",
			Use:  Single{jen.Id("n").Dot("Key").Op("!=").Nil()},
		},
		Node{
			Name: "Value",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterValue",
			Use:  Single{jen.Id("n").Dot("Value").Op("!=").Nil()},
		},
		Token{
			Name:   "Tok",
			Exists: Single{jen.Id("n").Dot("Tok").Dot("IsValid").Call()},
			Token:  Single{jen.Id("n").Dot("Tok")},
			Position: Position{
				Get: jen.Id("n").Dot("TokPos"),
				Set: jen.Id("n").Dot("TokPos").Op("=").Id("pos"),
			},
		},
		Token{
			Name:  "Range",
			Token: Single{jen.Qual("go/token", "RANGE")},
		},
		Decoration{
			Name: "AfterRange",
		},
		Node{
			Name: "X",
			Type: Type{"Expr", false},
		},
		Decoration{
			Name: "AfterX",
		},
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name: "Name",
			Type: Type{"Ident", true},
		},
		Decoration{
			Name: "AfterName",
			Use:  Single{jen.Id("n").Dot("Name").Op("!=").Nil()},
		},
		Node{
			Name: "Path",
			Type: Type{"BasicLit", false},
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
			Elem:      Type{"Ident", true},
			Separator: token.COMMA,
		},
		Decoration{
			Name: "AfterNames",
			Use:  Single{jen.Id("n").Dot("Type").Op("!=").Nil()},
		},
		Node{
			Name: "Type",
			Type: Type{"Expr", false},
		},
		Token{
			Name:  "Assign",
			Token: Single{jen.Qual("go/token", "ASSIGN")},
		},
		Decoration{
			Name: "AfterAssign",
		},
		List{
			Name:      "Names",
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
			Name: "Name",
			Type: Type{"Ident", true},
		},
		Token{
			Name:  "Assign",
			Token: Single{jen.Qual("go/token", "ASSIGN")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Assign").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Assign"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Assign"),
				Set: jen.Id("n").Dot("Assign").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterName",
		},
		Node{
			Name: "Type",
			Type: Type{"Expr", false},
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
				Ast: jen.Int().Parens(jen.Id("n").Dot("End").Call().Op("-").Id("n").Dot("Pos").Call()),
				Dst: jen.Id("n").Dot("Length"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("From"),
				Set: jen.Id("n").Dot("From").Op("=").Id("pos"),
			},
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
			Name:  "Tok",
			Token: Single{jen.Id("n").Dot("Tok")},
			Position: Position{
				Get: jen.Id("n").Dot("TokPos"),
				Set: jen.Id("n").Dot("TokPos").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterTok",
		},
		Token{
			Name:  "Lparen",
			Token: Single{jen.Id("n").Dot("LPAREN")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Lparen").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Lparen"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Lparen"),
				Set: jen.Id("n").Dot("Lparen").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterLparen",
			Use: Double{
				Ast: jen.Id("n").Dot("Lparen").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Lparen"),
			},
		},
		List{
			Name:      "Specs",
			Elem:      Type{"Spec", false},
			Separator: token.SEMICOLON,
		},
		Token{
			Name:  "Rparen",
			Token: Single{jen.Id("n").Dot("RPAREN")},
			Exists: Double{
				Ast: jen.Id("n").Dot("Rparen").Dot("IsValid").Call(),
				Dst: jen.Id("n").Dot("Rparen"),
			},
			Position: Position{
				Get: jen.Id("n").Dot("Rparen"),
				Set: jen.Id("n").Dot("Rparen").Op("=").Id("pos"),
			},
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
			Name: "Type",
			Type: Type{"FuncType", true},
		},
		Decoration{
			Name: "Start",
		},
		Token{
			Name:  "Func",
			Token: Single{jen.Qual("go/token", "FUNC")},
			Position: Position{
				Get: jen.Id("n").Dot("Type").Dot("Func"),
				Set: jen.Id("n").Dot("Type").Dot("Func").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterFunc",
		},
		// TODO: render any decorations from n.Type.AfterFunc (but never save them there)
		Node{
			Name: "Recv",
			Type: Type{"FieldList", true},
		},
		Decoration{
			Name: "AfterRecv",
			Use:  Single{jen.Id("n").Dot("Recv").Op("!=").Nil()},
		},
		Node{
			Name: "Name",
			Type: Type{"Ident", true},
		},
		Decoration{
			Name: "AfterName",
		},
		Node{
			Name:  "Params",
			Type:  Type{"FieldList", true},
			Field: "Type",
		},
		Decoration{
			Name: "AfterParams",
		},
		// TODO: render any decorations from n.Type.AfterParams (but never save them there)
		Node{
			Name:  "Results",
			Type:  Type{"FieldList", true},
			Field: "Type",
		},
		Decoration{
			Name: "AfterResults",
			Use:  Single{jen.Id("n").Dot("Type").Dot("Recv").Op("!=").Nil()},
		},
		// TODO: render any decorations from n.Type.AfterResults (but never save them there)
		Node{
			Name: "Body",
			Type: Type{"BlockStmt", true},
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
			Name:  "Package",
			Token: Single{jen.Qual("go/token", "PACKAGE")},
			Position: Position{
				Get: jen.Id("n").Dot("Package"),
				Set: jen.Id("n").Dot("Package").Op("=").Id("pos"),
			},
		},
		Decoration{
			Name: "AfterPackage",
		},
		Node{
			Name: "Name",
			Type: Type{"Ident", true},
		},
		Decoration{
			Name: "AfterName",
		},
		List{
			Name:      "Decls",
			Elem:      Type{"Decl", false},
			Separator: token.SEMICOLON,
		},
		// Never want to attach decorations to the end of a list of declarations - always better to
		// attach to the last statement.
	},
}

type Init struct {
	Name string
	Type Type
}

type Decoration struct {
	Name string
	Use  Code
}

type String struct {
	Name     string
	Position Position
}

type List struct {
	Name      string
	Elem      Type
	Separator token.Token
}

type Node struct {
	Name  string
	Type  Type
	Field string
}

type Token struct {
	Name     string
	Exists   Code
	Token    Code
	Position Position
}

type Ignored struct {
	Length   Code
	Position Position
}

// Value that must be copied from ast.Node to dst.Node but doesn't result in anything rendered to the output.
type Value struct {
	Name string
}

type Whitespace struct{}

type Position struct {
	Get *jen.Statement
	Set *jen.Statement
}

type Code interface {
	GetStatement(ast bool) *jen.Statement
}

type Single struct{ *jen.Statement }

func (s Single) GetStatement(ast bool) *jen.Statement {
	return s.Statement
}

type Double struct {
	Ast *jen.Statement
	Dst *jen.Statement
}

func (d Double) GetStatement(ast bool) *jen.Statement {
	if ast {
		return d.Ast
	}
	return d.Dst
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
