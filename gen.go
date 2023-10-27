package c2p

import (
	"c2p/c/parse"
	"fmt"
	"strconv"
	"strings"
)

var nest int
var assignMode bool

func init() {
	nest = 0
	assignMode = false
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
	code := fmt.Sprintf("%s%s = None\n", genIndent(), node.VariableDeclareField.Ident.S)
	return code, nil
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
			//code += c
			if assignMode {
				code += fmt.Sprintf("%s%s\n", genIndent(), c)
				assignMode = false
			} else {
				code += c
			}
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
		var code string
		cond, err := expr(node.IfElseField.Cond)
		if err != nil {
			return "", err
		}
		ifBlock, err := stmt(node.IfElseField.IfBlock)
		if err != nil {
			return "", err
		}

		code += fmt.Sprintf("%sif (%s):\n", genIndent(), cond)
		code += ifBlock

		// only if-block
		if node.IfElseField.ElseBlock == nil {
			return code, nil
		}

		// has else-block
		elseBlock, err := stmt(node.IfElseField.ElseBlock)
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%selse:\n", genIndent())
		code += elseBlock

		return code, nil
	case parse.While:
	case parse.For:
		var code string
		init_, err := expr(node.ForField.Init)
		if err != nil {
			return "", err
		}
		cond, err := expr(node.ForField.Cond)
		if err != nil {
			return "", err
		}
		loop, err := expr(node.ForField.Loop)
		if err != nil {
			return "", err
		}
		body, err := stmt(node.ForField.Block)
		if err != nil {
			return "", err
		}

		code += init_
		code += fmt.Sprintf("%swhile True:\n", genIndent())
		nest++
		code += fmt.Sprintf("%sif (not (%s)):\n", genIndent(), cond)
		nest++
		code += fmt.Sprintf("%sbreak\n", genIndent())
		nest--
		code += body
		code += loop
		nest--

		return code, nil
	}
	return expr(node)
}

func expr(node *parse.Node) (string, error) {
	return assign(node)
}

func assign(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.VariableDeclare:
		var code string
		if assignMode {
			code = fmt.Sprintf("%s", node.VariableDeclareField.Ident.S)
		} else {
			code = fmt.Sprintf("%s%s = None\n", genIndent(), node.VariableDeclareField.Ident.S)
		}
		return code, nil
	case parse.Assign:
		assignMode = true // 変数宣言のインデントなしモード
		to, err := expr(node.AssignField.To)
		if err != nil {
			return "", err
		}
		assignMode = true
		val, err := expr(node.AssignField.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s = %s\n", genIndent(), to, val), nil
	}
	return andor(node)
}

func andor(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.And:
		lhs, err := andor(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := andor(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s and %s", lhs, rhs), nil
	case parse.Or:
		lhs, err := andor(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := andor(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s or %s", lhs, rhs), nil
	}
	return equality(node)
}

func equality(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Eq:
		lhs, err := equality(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := equality(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s == %s", lhs, rhs), nil
	case parse.Ne:
		lhs, err := equality(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := equality(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s != %s", lhs, rhs), nil
	}
	return relational(node)
}

func relational(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Lt:
		lhs, err := relational(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := relational(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s < %s", lhs, rhs), nil
	case parse.Le:
		lhs, err := relational(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := relational(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s <= %s", lhs, rhs), nil
	case parse.Gt:
		lhs, err := relational(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := relational(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s > %s", lhs, rhs), nil
	case parse.Ge:
		lhs, err := relational(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := relational(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s >= %s", lhs, rhs), nil
	}

	return add(node)
}

func add(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Add:
		lhs, err := add(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := add(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s + %s", lhs, rhs), nil
	case parse.Sub:
		lhs, err := add(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := add(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s - %s", lhs, rhs), nil
	}
	return mul(node)
}

func mul(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Mul:
		lhs, err := mul(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := mul(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s * %s", lhs, rhs), nil
	case parse.Div:
		lhs, err := mul(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := mul(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s / %s", lhs, rhs), nil
	case parse.Mod:
		lhs, err := mul(node.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := mul(node.RHS)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %% %s", lhs, rhs), nil
	}
	return unary(node)
}

func unary(node *parse.Node) (string, error) {
	// todo
	return primary(node)
}

func primary(node *parse.Node) (string, error) {
	// todo
	switch node.Kind {
	case parse.Ident:
		return node.S, nil
	case parse.Call:
		return call(node)
	case parse.Int:
		return strconv.Itoa(node.I), nil
	}
	return "", nil
}

func call(node *parse.Node) (string, error) {
	ident := node.CallField.Ident.S
	args := node.CallField.Args.PolynomialField
	var code string
	if ident == "printf" {

	} else {
		var arguments string
		for i, arg := range args.Values {
			a, err := expr(arg)
			if err != nil {
				return "", err
			}
			if i != 0 {
				arguments += ", "
			}
			arguments += a
		}
		code = fmt.Sprintf("%s(%s)", ident, arguments)
	}
	assignMode = true
	return code, nil
}
