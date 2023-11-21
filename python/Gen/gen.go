package Gen

import (
	"c2p/python/parse"
	"fmt"
	"strconv"
	"strings"
)

var nest int

type Line struct {
	C string
	N int
}

func NewLine(c string, n int) *Line {
	return &Line{C: c, N: n}
}

func genLine(lines []*Line) string {
	var c string
	for _, l := range lines {
		c += fmt.Sprintf("%s%s\n", strings.Repeat("    ", l.N), l.C)
	}
	return c
}

func Gen(nodes []*parse.Node) (string, error) {
	lines, err := program(nodes)
	if err != nil {
		return "", err
	}
	code := genLine(lines)
	code += "if __name__ == \"__main__\":\n    main()"
	return code, nil
}

func program(nodes []*parse.Node) ([]*Line, error) {
	var lines []*Line
	for _, node := range nodes {
		l, err := toplevel(node)
		if err != nil {
			return nil, err
		}
		lines = append(lines, l...)
	}

	return lines, nil
}

func toplevel(node *parse.Node) ([]*Line, error) {
	switch node.Kind {
	case parse.FunctionDefine:
		return functionDefine(node)
	}
	return nil, nil
}

func functionDefine(node *parse.Node) ([]*Line, error) {
	field := node.FunctionDefineField

	ident := field.Ident.S

	params, err := functionDefineParams(field.Params)
	if err != nil {
		return nil, err
	}

	bl, err := stmt(field.Block)
	if err != nil {
		return nil, err
	}

	//c := fmt.Sprintf("def %s(%s):", ident, params) + "\n"
	//c += bl[0]
	//c += "\n"
	//
	//code += c
	var lines []*Line
	lines = append(lines, NewLine(fmt.Sprintf("def %s(%s):", ident, params), 0))
	lines = append(lines, bl...)

	return lines, nil
}

func functionDefineParams(node *parse.Node) (string, error) {
	_ = node
	return "", nil
}

func stmt(node *parse.Node) ([]*Line, error) {
	switch node.Kind {
	case parse.Block:
		//var c string
		var lines []*Line
		nest++
		for _, statementNode := range node.BlockField.Stmts {
			statements, err := stmt(statementNode)
			if err != nil {
				return nil, err
			}
			//for _, line := range statements {
			//	//c += fmt.Sprintf("%s%s", indent(), line) + "\n"
			//	lines = append(lines, line...)
			//}
			lines = append(lines, statements...)
		}
		nest--
		return lines, nil
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
		return []*Line{NewLine(fmt.Sprintf("return %s", strings.Join(values, ", ")), nest)}, nil
	case parse.IfElse:
		field := node.IfElseField
		var lines []*Line

		// 条件式
		cond, err := expr(field.Cond)
		if err != nil {
			return nil, err
		}
		lines = append(lines, NewLine(fmt.Sprintf("if (%s):", cond), nest))

		// IFの中身の作成
		ifBlock, err := stmt(field.IfBlock)
		if err != nil {
			return nil, err
		}
		lines = append(lines, ifBlock...)

		// もしElseがなかったら
		if field.ElseBlock == nil {
			return lines, nil
		}

		// Elseがあったら
		lines = append(lines, NewLine(fmt.Sprintf("else:"), nest))
		nest++
		elseBlock, err := stmt(field.ElseBlock)
		if err != nil {
			return nil, err
		}
		nest--
		lines = append(lines, elseBlock...)

		return lines, nil
	case parse.While:
		field := node.WhileField
		var lines []*Line

		// 初期
		wInit, err := expr(field.Init)
		if err != nil {
			return nil, err
		}
		lines = append(lines, NewLine(wInit, nest))

		// 条件式
		wCond, err := expr(field.Cond)
		if err != nil {
			return nil, err
		}
		lines = append(lines, NewLine(fmt.Sprintf("while (%s):", wCond), nest))

		// 中身
		wBlock, err := stmt(field.Block)
		if err != nil {
			return nil, err
		}
		lines = append(lines, wBlock...)

		// ループ
		wLoop, err := expr(field.Loop)
		if err != nil {
			return nil, err
		}
		// whileの中に入れてあげる
		lines = append(lines, NewLine(wLoop, nest+1))

		return lines, nil
	case parse.For:
	}

	e, err := expr(node)
	if err != nil {
		return nil, err
	}
	return []*Line{NewLine(e, nest)}, nil
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
	ident := node.CallField.Ident.S
	args := node.CallField.Args.PolynomialField

	var code string
	if ident == "printf" {
		var arguments string
		for i, arg := range args.Values[1:] {
			a, err := expr(arg)
			if err != nil {
				return "", err
			}
			if i != 0 {
				arguments += ", "
			}
			arguments += a
		}

		f, err := expr(args.Values[0])
		if err != nil {
			return "", err
		}
		if arguments == "" {
			code = fmt.Sprintf("print(%#v)", f)
		} else {
			code = fmt.Sprintf("print(%#v.format(%s))", formatting(f), arguments)
		}
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

	return code, nil
}

func formatting(s string) string {
	s = strings.ReplaceAll(s, "%d", "{}")
	s = strings.ReplaceAll(s, "%s", "{}")
	return s
}
