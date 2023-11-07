package Gen

import (
	"c2p/python/parse"
	"fmt"
	"strconv"
	"strings"
)

var nest int
var code string

func init() {
	nest = 0
	code = ""
}

func indent() string {
	return strings.Repeat("\t", nest)
}

func list(s string) []string {
	return []string{s}
}

func Gen(nodes []*parse.Node) (string, error) {
	err := program(nodes)
	if err != nil {
		return "", err
	}
	code += "if __name__ == \"__main__\":\n\tmain()"
	return code, nil
}

func program(nodes []*parse.Node) error {
	for _, node := range nodes {
		err := toplevel(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func toplevel(node *parse.Node) error {
	switch node.Kind {
	case parse.FunctionDefine:
		return functionDefine(node)
	}
	return nil
}

func functionDefine(node *parse.Node) error {
	field := node.FunctionDefineField

	ident := field.Ident.S

	params, err := functionDefineParams(field.Params)
	if err != nil {
		return err
	}

	bl, err := stmt(field.Block)
	if err != nil {
		return err
	}

	c := fmt.Sprintf("def %s(%s):", ident, params) + "\n"
	c += bl[0]
	c += "\n"

	code += c

	return nil
}

func functionDefineParams(node *parse.Node) (string, error) {
	return "", nil
}

func stmt(node *parse.Node) ([]string, error) {
	switch node.Kind {
	case parse.Block:
		var c string
		nest++
		for _, statementNode := range node.BlockField.Stmts {
			statements, err := stmt(statementNode)
			if err != nil {
				return nil, err
			}
			for _, line := range statements {
				c += fmt.Sprintf("%s%s", indent(), line) + "\n"
			}
		}
		nest--
		return list(c), nil
	case parse.Return:
		var values []string
		for _, valueNode := range node.ReturnField.Values {
			v, err := expr(valueNode)
			if err != nil {
				if err != nil {
					return nil, err
				}
			}
			values = append(values, v)
		}
		return list(fmt.Sprintf("return %s", strings.Join(values, ", "))), nil
	case parse.IfElse:
		field := node.IfElseField
		cond, err := expr(field.Cond)
		if err != nil {
			return nil, err
		}
		ifBlock, err := stmt(field.IfBlock)
		if err != nil {
			return nil, err
		}

		if field.ElseBlock == nil {
			// ifだけ
			rt := []string{
				fmt.Sprintf("if (%s):", cond),
			}
			rt = append(rt, ifBlock...)
			return rt, nil
		}

		elseBlock, err := stmt(field.ElseBlock)
		if err != nil {
			return nil, err
		}

		rt := []string{
			fmt.Sprintf("if (%s):", cond),
		}
		rt = append(rt, ifBlock...)
		rt = append(rt, []string{
			"else:",
		}...)
		rt = append(rt, elseBlock...)
		return rt, nil
	case parse.While:
		field := node.WhileField
		wInit, err := expr(field.Init)
		if err != nil {
			return nil, err
		}
		wCond, err := expr(field.Cond)
		if err != nil {
			return nil, err
		}
		wLoop, err := expr(field.Loop)
		if err != nil {
			return nil, err
		}
		wBlock, err := stmt(field.Block)
		if err != nil {
			return nil, err
		}
		// blockの最後にloopを入れてあげる
		wBlock = append(wBlock, wLoop)

		rt := []string{
			wInit,
			fmt.Sprintf("while (%s):", wCond),
		}
		rt = append(rt, wBlock...)
		return rt, nil
	case parse.For:
	}

	e, err := expr(node)
	if err != nil {
		return nil, err
	}
	return list(e), nil
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
			return "", nil
		}
		rhs, err := add(node.RHS)
		if err != nil {
			return "", nil
		}
		return fmt.Sprintf("%s + %s", lhs, rhs), nil
	case parse.Sub:
		lhs, err := add(node.LHS)
		if err != nil {
			return "", nil
		}
		rhs, err := add(node.RHS)
		if err != nil {
			return "", nil
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
	return primary(node)
}

func primary(node *parse.Node) (string, error) {
	switch node.Kind {
	case parse.Ident:
		return node.S, nil
	case parse.Call:
		return call(node)
	case parse.Int:
		return strconv.Itoa(node.I), nil
	case parse.Float:
		return strconv.FormatFloat(node.F, 'f', -1, 64), nil
	case parse.String:
		return node.S, nil
	}
	return "", nil
}

func call(node *parse.Node) (string, error) {
	return "call()", nil
}
