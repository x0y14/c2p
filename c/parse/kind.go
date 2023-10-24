package parse

type NodeKind int

const (
	_ NodeKind = iota

	// toplevel
	Include
	Define

	VariableDeclare
	FunctionDeclare
	FunctionDeclareParameters

	FunctionDefine
	FunctionDefineParameters

	// stmt
	Block
	Return
	IfElse
	While
	For

	// expr
	Assign

	// assign
	And
	Or

	// equality
	Eq
	Ne

	// relational
	Lt
	Le
	Gt
	Ge

	// add
	Add
	Sub

	// mul
	Mul
	Div
	Mod

	// unary
	Not

	// primary
	Parenthesis
	Ident
	Int
	Float
	String
	Null
	Call
	CallArgs
)
