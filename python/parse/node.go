package parse

type Node struct {
	Kind NodeKind

	*FunctionDefineField
	*FunctionDefineParamsField
	*BlockField
	*PolynomialField
	*ReturnField
	*IfElseField
	*WhileField
	*ForField
	*AssignField
	*BinaryField
	*LiteralField
	*NotFiled
	*CallField
}

func NewIdentNode(id string) *Node {
	return NewLiteralFieldNode(LIdent, 0, 0, id)
}
func NewIntNode(i int) *Node {
	return NewLiteralFieldNode(LInt, i, 0, "")
}
func NewFloatNode(f float64) *Node {
	return NewLiteralFieldNode(LFloat, 0, f, "")
}
func NewStringNode(s string) *Node {
	return NewLiteralFieldNode(LString, 0, 0, s)
}
func NewNoneNode() *Node {
	return NewLiteralFieldNode(LNone, 0, 0, "")
}
