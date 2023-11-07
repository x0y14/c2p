package parse

type FunctionDefineField struct {
	Ident  *Node
	Params *Node
	Block  *Node
}

func NewFunctionDefineFieldNode(id, params, block *Node) *Node {
	return &Node{
		Kind: FunctionDefine,
		FunctionDefineField: &FunctionDefineField{
			Ident:  id,
			Params: params,
			Block:  block,
		},
	}
}

type FunctionDefineParamsField struct {
}

type BlockField struct {
	Stmts []*Node
}

func NewBlockFieldNode(stmts []*Node) *Node {
	return &Node{Kind: Block, BlockField: &BlockField{Stmts: stmts}}
}

type PolynomialField struct {
	Values []*Node
}

func NewPolynomialFieldNode(kind NodeKind, values []*Node) *Node {
	return &Node{
		Kind:            kind,
		PolynomialField: &PolynomialField{Values: values},
	}
}

type ReturnField struct {
	Values []*Node
}

func NewReturnFieldNode(values []*Node) *Node {
	return &Node{
		Kind:        Return,
		ReturnField: &ReturnField{Values: values},
	}
}

type IfElseField struct {
	Cond      *Node
	IfBlock   *Node
	ElseBlock *Node
}

func NewIfElseFieldNode(cond, ifBlock, elseBlock *Node) *Node {
	return &Node{Kind: IfElse, IfElseField: &IfElseField{
		Cond:      cond,
		IfBlock:   ifBlock,
		ElseBlock: elseBlock,
	}}
}

type WhileField struct {
	Init  *Node
	Cond  *Node
	Loop  *Node
	Block *Node
}

func NewWhileFieldNode(init, cond, loop, block *Node) *Node {
	return &Node{Kind: While, WhileField: &WhileField{
		Init:  init,
		Cond:  cond,
		Loop:  loop,
		Block: block,
	}}
}

type ForField struct {
	// for i in range()
	// for i in []
	Idents []*Node
	Values []*Node
	Block  *Node
}

func NewForFieldNode(ids, values []*Node, block *Node) *Node {
	return &Node{Kind: For, ForField: &ForField{
		Idents: ids,
		Values: values,
		Block:  block,
	}}
}

type AssignField struct {
	// x, y = 1, 2...?
	To    *Node
	Value *Node
}

func NewAssignFieldNode(to, value *Node) *Node {
	return &Node{Kind: Assign, AssignField: &AssignField{
		To:    to,
		Value: value,
	}}
}

type BinaryField struct {
	LHS *Node
	RHS *Node
}

func NewBinaryFieldNode(kind NodeKind, lhs, rhs *Node) *Node {
	return &Node{Kind: kind, BinaryField: &BinaryField{
		LHS: lhs,
		RHS: rhs,
	}}
}

type LiteralField struct {
	Kind LiteralKind
	I    int
	F    float64
	S    string
}

func NewLiteralFieldNode(kind LiteralKind, i int, f float64, s string) *Node {
	var nk NodeKind
	switch kind {
	case LIdent:
		nk = Ident
	case LInt:
		nk = Int
	case LFloat:
		nk = Float
	case LString:
		nk = String
	case LNone:
		nk = None
	}
	return &Node{
		Kind: nk,
		LiteralField: &LiteralField{
			Kind: kind,
			I:    i,
			F:    f,
			S:    s,
		},
	}
}

type NotFiled struct {
	Value *Node
}

func NewNotFieldNode(value *Node) *Node {
	return &Node{Kind: Not, NotFiled: &NotFiled{Value: value}}
}

type CallField struct {
	Ident *Node
	Args  *Node
}

func NewCallFieldNode(id, args *Node) *Node {
	return &Node{Kind: Call, CallField: &CallField{
		Ident: id,
		Args:  args,
	}}
}
