package parse

type Node struct {
	Kind NodeKind

	*IncludeField
	*DefineField
	*VariableDeclareField
	*FunctionDeclareField
	*PolynomialField
	*FunctionDefineField
	*BlockField
	*ReturnField
	*IfElseField
	*WhileField
	*ForField
	*AssignField
	*BinaryField
	*LiteralField
	*NotField
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
func NewNullNode() *Node {
	return NewLiteralFieldNode(LNull, 0, 0, "")
}
