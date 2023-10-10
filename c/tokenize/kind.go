package tokenize

type TokenKind int

const (
	_ TokenKind = iota
	Eof

	Pre

	Ident
	Int
	Float
	String

	Lrb
	Rrb
	Lsb
	Rsb
	Lcb
	Rcb
	Dot
	Comma
	Colon
	Semi

	Add
	Sub
	Mul
	Div
	Mod

	Eq
	Ne
	Gt
	Lt
	Ge
	Le

	Assign

	And
	Or
	Not
)
