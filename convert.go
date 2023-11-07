package c2p

import (
	cparse "c2p/c/parse"
	pparse "c2p/python/parse"
	"fmt"
)

func CToP(cNodes []*cparse.Node) ([]*pparse.Node, error) {
	var pNodes []*pparse.Node

	for _, cn := range cNodes {
		pn, err := cpToplevel(cn)
		if err != nil {
			return nil, err
		}
		// cにはあるけどpyにはないものはnilを返す
		if pn == nil {
			continue
		}
		pNodes = append(pNodes, pn)
	}

	return pNodes, nil
}

func cpToplevel(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.FunctionDefine:
		return cpFunctionDefine(cNode)
	}
	return nil, nil
}

func cpFunctionDefine(cNode *cparse.Node) (*pparse.Node, error) {
	pField := cNode.FunctionDefineField
	ident := pField.Ident.S

	pParams, err := cpFunctionDefineParams(pField.Parameters)
	if err != nil {
		return nil, err
	}

	pStatements, err := cpStmt(pField.Block)
	if err != nil {
		return nil, err
	}

	return pparse.NewFunctionDefineFieldNode(
		pparse.NewIdentNode(ident),
		pParams,
		pStatements,
	), nil
}

func cpFunctionDefineParams(cNode *cparse.Node) (*pparse.Node, error) {
	return nil, nil
}

func cpStmt(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.Block:
		var pyStmts []*pparse.Node
		for _, cStmt := range cNode.BlockField.Stmt {
			pyStmt, err := cpStmt(cStmt)
			if err != nil {
				return nil, err
			}
			if pyStmt == nil {
				continue
			}
			pyStmts = append(pyStmts, pyStmt)
		}
		return pparse.NewBlockFieldNode(pyStmts), nil
	case cparse.Return:
		rv, err := cpExpr(cNode.ReturnField.Value)
		if err != nil {
			return nil, err
		}
		return pparse.NewReturnFieldNode(
			[]*pparse.Node{rv}, // cは多分１つしか戻り値ないけど、pyは複数を許容する
		), nil
	case cparse.IfElse:
		cond, err := cpExpr(cNode.IfElseField.Cond)
		if err != nil {
			return nil, err
		}
		ifBlock, err := cpStmt(cNode.IfElseField.IfBlock)
		if err != nil {
			return nil, err
		}
		if cNode.IfElseField.ElseBlock == nil {
			return pparse.NewIfElseFieldNode(cond, ifBlock, nil), nil
		}
		elseBlock, err := cpStmt(cNode.IfElseField.ElseBlock)
		if err != nil {
			return nil, err
		}
		return pparse.NewIfElseFieldNode(cond, ifBlock, elseBlock), nil

	case cparse.While:
	case cparse.For:
		init, err := cpExpr(cNode.ForField.Init)
		if err != nil {
			return nil, err
		}
		cond, err := cpExpr(cNode.ForField.Cond)
		if err != nil {
			return nil, err
		}
		loop, err := cpExpr(cNode.ForField.Loop)
		if err != nil {
			return nil, err
		}
		block, err := cpStmt(cNode.ForField.Block)
		if err != nil {
			return nil, err
		}
		return pparse.NewWhileFieldNode(init, cond, loop, block), nil
	}
	return cpExpr(cNode)
}

func cpExpr(cNode *cparse.Node) (*pparse.Node, error) {
	return cpAssign(cNode)
}

func cpAssign(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.VariableDeclare:
		// pythonには存在しないので、空を返しておくとりあえず
		return nil, nil
	case cparse.Assign:
		to, err := cpExpr(cNode.AssignField.To)
		if err != nil {
			return nil, err
		}
		val, err := cpExpr(cNode.AssignField.Value)
		if err != nil {
			return nil, err
		}
		return pparse.NewAssignFieldNode(to, val), nil
	}
	return cpAndor(cNode)
}

func cpAndor(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.And:
		lhs, err := cpAndor(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpAndor(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.And, lhs, rhs), nil
	case cparse.Or:
		lhs, err := cpAndor(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpAndor(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Or, lhs, rhs), nil
	}
	return cpEquality(cNode)
}

func cpEquality(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.Eq:
		lhs, err := cpEquality(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpEquality(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Eq, lhs, rhs), nil
	case cparse.Ne:
		lhs, err := cpEquality(cNode.RHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpEquality(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Ne, lhs, rhs), nil
	}
	return cpRelational(cNode)
}

func cpRelational(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.Lt:
		lhs, err := cpRelational(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpRelational(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Lt, lhs, rhs), nil
	case cparse.Le:
		lhs, err := cpRelational(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpRelational(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Le, lhs, rhs), nil
	case cparse.Gt:
		lhs, err := cpRelational(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpRelational(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Gt, lhs, rhs), nil
	case cparse.Ge:
		lhs, err := cpRelational(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpRelational(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Ge, lhs, rhs), nil
	}
	return cpAdd(cNode)
}

func cpAdd(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.Add:
		lhs, err := cpAdd(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpAdd(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Add, lhs, rhs), nil
	case cparse.Sub:
		lhs, err := cpAdd(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpAdd(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Sub, lhs, rhs), nil
	}
	return cpMul(cNode)
}

func cpMul(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.Mul:
		lhs, err := cpMul(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpMul(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Mul, lhs, rhs), nil
	case cparse.Div:
		lhs, err := cpMul(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpMul(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Div, lhs, rhs), nil
	case cparse.Mod:
		lhs, err := cpMul(cNode.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := cpMul(cNode.RHS)
		if err != nil {
			return nil, err
		}
		return pparse.NewBinaryFieldNode(pparse.Mod, lhs, rhs), nil
	}
	return cpUnary(cNode)
}

func cpUnary(cNode *cparse.Node) (*pparse.Node, error) {
	return cpPrimary(cNode)
}

func cpPrimary(cNode *cparse.Node) (*pparse.Node, error) {
	switch cNode.Kind {
	case cparse.Ident:
		return pparse.NewIdentNode(cNode.S), nil
	case cparse.Call:
		return cpCall(cNode)
	case cparse.Int:
		return pparse.NewLiteralFieldNode(pparse.LInt, cNode.I, 0, ""), nil
	case cparse.Float:
		return pparse.NewLiteralFieldNode(pparse.LFloat, 0, cNode.F, ""), nil
	case cparse.String:
		return pparse.NewLiteralFieldNode(pparse.LString, 0, 0, cNode.S), nil
	}
	return nil, fmt.Errorf("unexpected error: primary: %v", cNode)
}

func cpPolynomial(cNode *cparse.Node) (*pparse.Node, error) {
	var values []*pparse.Node
	for _, cv := range cNode.PolynomialField.Values {
		pv, err := cpExpr(cv)
		if err != nil {
			return nil, err
		}
		values = append(values, pv)
	}
	return pparse.NewPolynomialFieldNode(pparse.CallArgs, values), nil
}

func cpCall(cNode *cparse.Node) (*pparse.Node, error) {
	ident := cNode.CallField.Ident.S
	args, err := cpPolynomial(cNode.CallField.Args)
	if err != nil {
		return nil, err
	}
	return pparse.NewCallFieldNode(
		pparse.NewIdentNode(ident),
		args,
	), nil
}
