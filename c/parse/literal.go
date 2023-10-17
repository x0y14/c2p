package parse

type LiteralKind int

const (
	_ LiteralKind = iota
	LIdent
	LInt
	LFloat
	LString
	LNull
)
