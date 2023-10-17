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
	if semi := consumeKind(tokenize.Semi); semi != nil {
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
	if semi := consumeKind(tokenize.Semi); semi != nil {
		return NewFunctionDeclareFieldNode(typeNode, identNode, paramsNode), nil
	}

	// expect { -> func def
	// todo
	return nil, nil
}

func varDeclare() (*Node, error) {
	return nil, nil
}

func funcDeclare() (*Node, error) {
	return nil, nil
}

func funcDefine() (*Node, error) {
	return nil, nil
}

func funcParams() (*Node, error) {
	return nil, nil
}
