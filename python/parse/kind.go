package parse

type NodeKind int

const (
	_ NodeKind = iota

	FunctionDefine
	FunctionDefineParams

	Block
	Return
	IfElse
	While
	For

	Assign
	And
	Or

	Eq
	Ne

	Lt
	Le
	Gt
	Ge

	Add
	Sub
	Mul
	Div

	Not

	Ident
	Int
	Float
	String
	None
	Call
	CallArgs
)
