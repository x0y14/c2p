package Gen

import (
	"c2p/python/parse"
	"fmt"
)

func Gen(nodes []*parse.Node) (string, error) {
	code, err := program(nodes)
	if err != nil {
		return "", err
	}
	return code, nil
}

func program(nodes []*parse.Node) (string, error) {
	code := ""
	for _, node := range nodes {
		tl, err := toplevel(node)
		if err != nil {
			return "", err
		}
		code += tl + "\n"
	}
	return code, nil
}

func toplevel(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.FunctionDefine:
		return functionDefine(node)
	}
	return "", nil
}

func functionDefine(node *parse.Node) (string, error) {
	field := node.FunctionDefineField

	fnName := field.Ident.S
	fnParams, err := funcDefineParameters(field.Params)
	if err != nil {
		return "", err
	}
	statements, err := stmt(field.Block)
	if err != nil {
		return "", err
	}
	code := fnName + fnParams + statements
	return code, nil
}

func funcDefineParameters(node *parse.Node) (string, error) {
	_ = node
	return "", nil
}

func stmt(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Block:
		var code string
		for _, statement := range node.BlockField.Stmts {
			c, err := stmt(statement)
			if err != nil {
				return "", err
			}
			code += c
		}
		return code, nil
	case parse.Return:
		rv := ""
		for _, v := range node.ReturnField.Values {
			rv_, err := expr(v)
			if err != nil {
				return "", err
			}
			rv += rv_
		}
		return "return " + rv, nil

	case parse.IfElse:
	case parse.While:
	case parse.For:

	}
	return expr(node)
}

func expr(node *parse.Node) (string, error) {
	return assign(node)
}

func assign(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Assign:
		to, err := expr(node.AssignField.To)
		if err != nil {
			return "", err
		}
		val, err := expr(node.AssignField.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", to, val), nil
	}

	return andor(node)
}

func andor(node *parse.Node) (string, error) {
	return equality(node)
}

func equality(node *parse.Node) (string, error) {
	return relational(node)
}

func relational(node *parse.Node) (string, error) {
	return add(node)
}

func add(node *parse.Node) (string, error) {
	return mul(node)
}

func mul(node *parse.Node) (string, error) {
	return unary(node)
}

func unary(node *parse.Node) (string, error) {
	return primary(node)
}

func primary(node *parse.Node) (string, error) {
	return "", nil
}

func call(node *parse.Node) (string, error) {
	return "", nil
}
