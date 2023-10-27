package c2p

import (
	"c2p/c/parse"
	"fmt"
	"strconv"
	"strings"
)

var nest int

func init() {
	nest = 0
}

func genIndent() string {
	return strings.Repeat("\t", nest)
}

func Gen(nodes []*parse.Node) (string, error) {
	code, err := program(nodes)
	if err != nil {
		return "", err
	}
	// run main function
	code += "if __name__ == \"__main__\":\n\tmain()"
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
	case parse.VariableDeclare:
		return variableDeclare(node)
	case parse.FunctionDeclare:
		return funcDeclare(node)
	case parse.FunctionDefine:
		return funcDefine(node)
	}
	return "", nil
}

func variableDeclare(node *parse.Node) (string, error) {
	return "", nil
}

func funcDeclare(node *parse.Node) (string, error) {
	return "", nil
}

func funcDefine(node *parse.Node) (string, error) {
	field := node.FunctionDefineField

	// 戻り値の型
	_ = field.Types

	// 関数名
	funcName := field.Ident.S

	// 引数
	params, err := funcDefineParameters(field.Parameters)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%sdef %s(%s):\n", genIndent(), funcName, params)
	statement, err := stmt(field.Block)
	if err != nil {
		return "", err
	}
	code += statement

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
		nest++
		for _, statement := range node.BlockField.Stmt {
			c, err := stmt(statement)
			if err != nil {
				return "", err
			}
			code += c
		}
		nest--
		return code, nil
	case parse.Return:
		returnVal, err := expr(node.ReturnField.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%sreturn %s\n", genIndent(), returnVal), nil
	case parse.IfElse:
	case parse.While:
	case parse.For:
	}
	return expr(node)
}

//func block(node *parse.Node) (string, error) {
//	var code string
//	nest++
//	for _, statement := range node.BlockField.Stmt {
//
//	}
//	nest--
//	return code, nil
//}

func expr(node *parse.Node) (string, error) {
	return assign(node)
}

func assign(node *parse.Node) (string, error) {
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
	switch node.Kind {
	case parse.Ident:
		return node.S, nil
	case parse.Int:
		return strconv.Itoa(node.I), nil
	}
	return "", nil
}
