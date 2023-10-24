package parse

import (
	"c2p/c/tokenize"
	"fmt"
)

var token *tokenize.Token

func isEof() bool {
	return token.Kind == tokenize.Eof
}

func peekKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Kind == kind {
		return token
	}
	return nil
}
func peekNextKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Next.Kind == kind {
		return token.Next
	}
	return nil
}
func peekNextNextKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Next.Next.Kind == kind {
		return token.Next.Next
	}
	return nil
}

func consumeKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Kind == kind {
		tok := token
		token = token.Next
		return tok
	}
	return nil
}

func consumeIdent(s string) *tokenize.Token {
	if token.Kind == tokenize.Ident && s == token.S {
		tok := token
		token = token.Next
		return tok
	}
	return nil
}

func expectKind(kind tokenize.TokenKind) (*tokenize.Token, error) {
	if token.Kind == kind {
		tok := token
		token = token.Next
		return tok, nil
	}
	return nil, fmt.Errorf("unexpected kind: want=%v, got=%v", kind, token.Kind)
}

func expectIdent(s string) (*tokenize.Token, error) {
	if token.Kind == tokenize.Ident && s == token.S {
		tok := token
		token = token.Next
		return tok, nil
	}
	return nil, fmt.Errorf("unexpected ident: want=%v, got=%v", s, token.S)
}

func Parse(head *tokenize.Token) ([]*Node, error) {
	token = head
	return program()
}

func program() ([]*Node, error) {
	var nodes []*Node
	for !isEof() {
		n, err := toplevel()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func toplevel() (*Node, error) {
	typ, err := expectKind(tokenize.Ident)
	if err != nil {
		return nil, err
	}
	id, err := expectKind(tokenize.Ident)
	if err != nil {
		return nil, err
	}
	typeNode := NewIdentNode(typ.S)
	identNode := NewIdentNode(id.S)

	// ; -> var declare
	if semi := peekKind(tokenize.Semi); semi != nil {
		return NewVariableDeclareFieldNode(typeNode, identNode), nil
	}

	// ( params? )
	// "("
	if _, err := expectKind(tokenize.Lrb); err != nil {
		return nil, err
	}
	// params
	paramsNode, err := funcParams()
	if err != nil {
		return nil, err
	}
	// ")"
	if _, err := expectKind(tokenize.Rrb); err != nil {
		return nil, err
	}
	// ; -> func declare
	if semi := peekKind(tokenize.Semi); semi != nil {
		return NewFunctionDeclareFieldNode(typeNode, identNode, paramsNode), nil
	}
	blockNode, err := stmt()
	if err != nil {
		return nil, err
	}
	return NewFunctionDefineFieldNode(typeNode, identNode, paramsNode, blockNode), nil
	//return nil, fmt.Errorf("unexpect token: %v", token)
}

func funcParams() (*Node, error) {
	if void := consumeIdent("void"); void != nil {
		return NewIdentNode(void.S), nil
	}
	return nil, nil
}

func stmt() (*Node, error) {
	// block
	if lcb := consumeKind(tokenize.Lcb); lcb != nil {
		var statements []*Node
		for consumeKind(tokenize.Rcb) == nil {
			statement, err := stmt()
			if err != nil {
				return nil, err
			}
			statements = append(statements, statement)
		}
		return NewBlockFieldNode(statements), nil
	}

	// return
	if return_ := consumeIdent("return"); return_ != nil {
		valNode, err := expr()
		if err != nil {
			return nil, err
		}
		if _, err := expectKind(tokenize.Semi); err != nil {
			return nil, err
		}
		return NewReturnFieldNode(valNode), nil
	}

	// if
	if if_ := consumeIdent("if"); if_ != nil {
		// (
		if _, err := expectKind(tokenize.Lrb); err != nil {
			return nil, err
		}
		// cond
		condNode, err := expr()
		if err != nil {
			return nil, err
		}
		// )
		if _, err := expectKind(tokenize.Rrb); err != nil {
			return nil, err
		}
		// ifBlock
		ifBlockNode, err := stmt()
		if err != nil {
			return nil, err
		}
		// elseBlock
		if els := consumeIdent("else"); els != nil {
			// elseBlock
			elseBlockNode, err := stmt()
			if err != nil {
				return nil, err
			}
			return NewIfElseFieldNode(condNode, ifBlockNode, elseBlockNode), nil
		} else {
			return NewIfElseFieldNode(condNode, ifBlockNode, nil), nil
		}
		//if lrb := consumeKind(tokenize.Lrb); lrb != nil {
		//	// else
		//	if _, err := expectIdent("else"); err != nil {
		//		return nil, err
		//	}
		//	// elseBlock
		//	elseBlockNode, err := stmt()
		//	if err != nil {
		//		return nil, err
		//	}
		//	return NewIfElseFieldNode(condNode, ifBlockNode, elseBlockNode), nil
		//} else {
		//	return NewIfElseFieldNode(condNode, ifBlockNode, nil), nil
		//}
	}

	// while
	if while_ := consumeIdent("while"); while_ != nil {
		// (
		if _, err := expectKind(tokenize.Lrb); err != nil {
			return nil, err
		}
		condNode, err := expr()
		if err != nil {
			return nil, err
		}
		// )
		if _, err := expectKind(tokenize.Rrb); err != nil {
			return nil, err
		}
		bodyBlockNode, err := stmt()
		if err != nil {
			return nil, err
		}
		return NewWhileFieldNode(condNode, bodyBlockNode), nil
	}

	// for
	if for_ := consumeIdent("for"); for_ != nil {
		var initNode *Node
		var condNode *Node
		var loopNode *Node
		// (
		if _, err := expectKind(tokenize.Lrb); err != nil {
			return nil, err
		}
		// init?
		if semi := consumeKind(tokenize.Semi); semi == nil {
			i, err := expr()
			if err != nil {
				return nil, err
			}
			initNode = i
			if _, err := expectKind(tokenize.Semi); err != nil {
				return nil, err
			}
		}
		// cond
		if semi := consumeKind(tokenize.Semi); semi == nil {
			c, err := expr()
			if err != nil {
				return nil, err
			}
			condNode = c
			if _, err := expectKind(tokenize.Semi); err != nil {
				return nil, err
			}
		}
		// loop
		if rrb := consumeKind(tokenize.Rrb); rrb == nil {
			l, err := expr()
			if err != nil {
				return nil, err
			}
			loopNode = l
			if _, err := expectKind(tokenize.Rrb); err != nil {
				return nil, err
			}
		}
		bodyBlock, err := stmt()
		if err != nil {
			return nil, err
		}
		return NewForFieldNode(initNode, condNode, loopNode, bodyBlock), nil
	}

	// expr
	e, err := expr()
	if err != nil {
		return nil, err
	}
	if _, err := expectKind(tokenize.Semi); err != nil {
		return nil, err
	}
	return e, nil
}

func expr() (*Node, error) {
	return assign()
}

func assign() (*Node, error) {
	// type
	if typ_ := peekKind(tokenize.Ident); typ_ != nil {
		// ident
		if id_ := peekNextKind(tokenize.Ident); id_ != nil {
			// declare or define?
			typ := consumeKind(tokenize.Ident)
			id := consumeKind(tokenize.Ident)
			// =?
			if ass := consumeKind(tokenize.Assign); ass != nil {
				// type ident = value;
				valNode, err := andor()
				if err != nil {
					return nil, err
				}
				return NewAssignFieldNode(
					NewVariableDeclareFieldNode(NewIdentNode(typ.S), NewIdentNode(id.S)), valNode), nil
			}
			// type value;
			return NewVariableDeclareFieldNode(NewIdentNode(typ.S), NewIdentNode(id.S)), nil
		}
	}
	// ident = value
	if id_ := peekKind(tokenize.Ident); id_ != nil {
		// =?
		if ass := peekNextKind(tokenize.Assign); ass != nil {
			if eq := peekNextNextKind(tokenize.Assign); eq == nil {
				typ := consumeKind(tokenize.Ident)
				_ = consumeKind(tokenize.Assign)
				// type ident = value;
				valNode, err := andor()
				if err != nil {
					return nil, err
				}
				return NewAssignFieldNode(NewIdentNode(typ.S), valNode), nil
			}
		}
	}

	return andor()
}

func andor() (*Node, error) {
	n, err := equality()
	if err != nil {
		return nil, err
	}
	for {
		if and := consumeKind(tokenize.And); and != nil {
			rhs, err := equality()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(And, n, rhs)
		} else if or := consumeKind(tokenize.Or); or != nil {
			rhs, err := equality()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Or, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func equality() (*Node, error) {
	n, err := relational()
	if err != nil {
		return nil, err
	}
	for {
		if eq := consumeKind(tokenize.Eq); eq != nil {
			rhs, err := relational()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Eq, n, rhs)
		} else if ne := consumeKind(tokenize.Ne); ne != nil {
			rhs, err := relational()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Ne, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func relational() (*Node, error) {
	n, err := add()
	if err != nil {
		return nil, err
	}
	for {
		if lt := consumeKind(tokenize.Lt); lt != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Lt, n, rhs)
		} else if le := consumeKind(tokenize.Le); le != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Le, n, rhs)
		} else if gt := consumeKind(tokenize.Gt); gt != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Gt, n, rhs)
		} else if ge := consumeKind(tokenize.Ge); ge != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Ge, n, rhs)
		} else {
			break
		}
	}

	return n, nil
}

func add() (*Node, error) {
	n, err := mul()
	if err != nil {
		return nil, err
	}
	for {
		if plus := consumeKind(tokenize.Add); plus != nil {
			rhs, err := mul()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Add, n, rhs)
		} else if minus := consumeKind(tokenize.Sub); minus != nil {
			rhs, err := mul()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Sub, n, rhs)
		} else {
			break
		}
	}

	return n, nil
}

func mul() (*Node, error) {
	n, err := unary()
	if err != nil {
		return nil, err
	}
	for {
		if star := consumeKind(tokenize.Mul); star != nil {
			rhs, err := unary()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Mul, n, rhs)
		} else if div := consumeKind(tokenize.Div); div != nil {
			rhs, err := unary()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Div, n, rhs)
		} else if mod := consumeKind(tokenize.Mod); mod != nil {
			rhs, err := unary()
			if err != nil {
				return nil, err
			}
			n = NewBinaryFieldNode(Mod, n, rhs)
		} else {
			break
		}
	}

	return n, nil
}

func unary() (*Node, error) {
	if plus := consumeKind(tokenize.Add); plus != nil {
		return primary()
	}
	if minus := consumeKind(tokenize.Sub); minus != nil {
		rhs, err := primary()
		if err != nil {
			return nil, err
		}
		return NewBinaryFieldNode(Sub, NewIntNode(0), rhs), nil
	}
	if not := consumeKind(tokenize.Not); not != nil {
		rhs, err := primary()
		if err != nil {
			return nil, err
		}
		return NewNotFieldNode(rhs), nil
	}
	return primary()
}

func primary() (*Node, error) {
	if lrb := consumeKind(tokenize.Lrb); lrb != nil {
		if lrb := consumeKind(tokenize.Lrb); lrb != nil {
		}
		v, err := expr()
		if err != nil {
			return nil, err
		}
		if _, err := expectKind(tokenize.Rrb); err != nil {
			return nil, err
		}
		return v, nil
	}

	// call-args
	if id := consumeKind(tokenize.Ident); id != nil {
		if lrb := consumeKind(tokenize.Lrb); lrb != nil {
			args, err := callArgs()
			if err != nil {
				return nil, err
			}
			if _, err := expectKind(tokenize.Rrb); err != nil {
				return nil, err
			}
			return NewCallFieldNode(NewIdentNode(id.S), args), nil
		}
		// normal ident
		return NewIdentNode(id.S), nil
	}

	if i := consumeKind(tokenize.Int); i != nil {
		return NewIntNode(i.I), nil
	}
	if f := consumeKind(tokenize.Float); f != nil {
		return NewFloatNode(f.F), nil
	}
	if s := consumeKind(tokenize.String); s != nil {
		return NewStringNode(s.S), nil
	}
	if null := consumeIdent("NULL"); null != nil {
		return NewNullNode(), nil
	}
	return nil, fmt.Errorf("unexpected error: primary: %v", token)
}

func callArgs() (*Node, error) {
	var values []*Node
	v, err := unary()
	if err != nil {
		return nil, err
	}
	values = append(values, v)
	for consumeKind(tokenize.Comma) != nil {
		v, err = unary()
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return NewPolynomialFieldNode(CallArgs, values), nil
}
