package parse

type IncludeField struct {
}

type DefineField struct {
}

type VariableDeclareField struct {
	Types *Node
	Ident *Node
}

func NewVariableDeclareFieldNode(typ, id *Node) *Node {
	return &Node{
		Kind: VariableDeclare,
		VariableDeclareField: &VariableDeclareField{
			Types: typ,
			Ident: id,
		},
	}
}

type FunctionDeclareField struct {
	Types      *Node
	Ident      *Node
	Parameters *Node
}

func NewFunctionDeclareFieldNode(typ, id, params *Node) *Node {
	return &Node{
		Kind: FunctionDeclare,
		FunctionDeclareField: &FunctionDeclareField{
			Types:      typ,
			Ident:      id,
			Parameters: params,
		},
	}
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

type FunctionDefineField struct {
	Types      *Node
	Ident      *Node
	Parameters *Node
	Block      *Node
}

func NewFunctionDefineFieldNode(typ, id, params, block *Node) *Node {
	return &Node{
		Kind: FunctionDefine,
		FunctionDefineField: &FunctionDefineField{
			Types:      typ,
			Ident:      id,
			Parameters: params,
			Block:      block,
		},
	}
}

type BlockField struct {
	Stmt []*Node
}

func NewBlockFieldNode(stmt []*Node) *Node {
	return &Node{
		Kind:       Block,
		BlockField: &BlockField{Stmt: stmt},
	}
}

type ReturnField struct {
	Value *Node
}

func NewReturnFieldNode(val *Node) *Node {
	return &Node{
		Kind:        Return,
		ReturnField: &ReturnField{Value: val},
	}
}

type IfElseField struct {
	Cond      *Node
	IfBlock   *Node
	ElseBlock *Node
}

func NewIfElseFieldNode(cond, ifBlock, elseBlock *Node) *Node {
	return &Node{
		Kind: IfElse,
		IfElseField: &IfElseField{
			Cond:      cond,
			IfBlock:   ifBlock,
			ElseBlock: elseBlock,
		},
	}
}

type WhileField struct {
	Cond  *Node
	Block *Node
}

func NewWhileFieldNode(cond, block *Node) *Node {
	return &Node{
		Kind: While,
		WhileField: &WhileField{
			Cond:  cond,
			Block: block,
		},
	}
}

type ForField struct {
	Init  *Node
	Cond  *Node
	Loop  *Node
	Block *Node
}

func NewForFieldNode(init, cond, loop, block *Node) *Node {
	return &Node{
		Kind: For,
		ForField: &ForField{
			Init:  init,
			Cond:  cond,
			Loop:  loop,
			Block: block,
		},
	}
}

type AssignField struct {
	To    *Node
	Value *Node
}

func NewAssignFieldNode(to, val *Node) *Node {
	return &Node{
		Kind: Assign,
		AssignField: &AssignField{
			To:    to,
			Value: val,
		},
	}
}

type BinaryField struct {
	LHS *Node
	RHS *Node
}

func NewBinaryFieldNode(kind NodeKind, l, r *Node) *Node {
	return &Node{
		Kind: kind,
		BinaryField: &BinaryField{
			LHS: l,
			RHS: r,
		},
	}
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
	case LNull:
		nk = Null
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

type NotField struct {
	V *Node
}

func NewNotFieldNode(v *Node) *Node {
	return &Node{Kind: Not, NotField: &NotField{V: v}}
}

type CallField struct {
	Ident *Node
	Args  *Node
}

func NewCallFieldNode(id, args *Node) *Node {
	return &Node{
		Kind: Call, CallField: &CallField{
			Ident: id,
			Args:  args,
		},
	}
}
