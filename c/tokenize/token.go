package tokenize

type Token struct {
	Kind TokenKind
	S    string
	I    int
	F    float64
	Next *Token
}

func NewToken(kind TokenKind, s string, i int, f float64) *Token {
	return &Token{
		Kind: kind,
		S:    s,
		I:    i,
		F:    f,
		Next: nil,
	}
}

func NewLiteralToken[T int | float64 | string](v T) *Token {
	switch any(v).(type) {
	case int:
		return NewToken(Int, "", any(v).(int), 0)
	case float64:
		return NewToken(Float, "", 0, any(v).(float64))
	case string:
		return NewToken(String, any(v).(string), 0, 0)
	}
	return nil
}

func NewSymbolToken(sym string) *Token {
	switch sym {
	case "(":
		return NewToken(Lrb, "", 0, 0)
	case ")":
		return NewToken(Rrb, "", 0, 0)
	case "[":
		return NewToken(Lsb, "", 0, 0)
	case "]":
		return NewToken(Rsb, "", 0, 0)
	case "{":
		return NewToken(Lcb, "", 0, 0)
	case "}":
		return NewToken(Rcb, "", 0, 0)
	case ".":
		return NewToken(Dot, "", 0, 0)
	case ",":
		return NewToken(Comma, "", 0, 0)
	case ":":
		return NewToken(Colon, "", 0, 0)
	case ";":
		return NewToken(Semi, "", 0, 0)

	case "+":
		return NewToken(Add, "", 0, 0)
	case "-":
		return NewToken(Sub, "", 0, 0)
	case "*":
		return NewToken(Mul, "", 0, 0)
	case "/":
		return NewToken(Div, "", 0, 0)
	case "%":
		return NewToken(Mod, "", 0, 0)

	case "==":
		return NewToken(Eq, "", 0, 0)
	case "!=":
		return NewToken(Ne, "", 0, 0)
	case ">":
		return NewToken(Gt, "", 0, 0)
	case "<":
		return NewToken(Lt, "", 0, 0)
	case ">=":
		return NewToken(Ge, "", 0, 0)
	case "<=":
		return NewToken(Le, "", 0, 0)

	case "=":
		return NewToken(Assign, "", 0, 0)

	case "&&":
		return NewToken(And, "", 0, 0)
	case "||":
		return NewToken(Or, "", 0, 0)
	case "!":
		return NewToken(Not, "", 0, 0)
	}
	return nil
}
