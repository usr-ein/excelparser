package parser

import (
	"fmt"
	"strings"
)

func StringifyNode(n Node, sheetName string) Formula {
	return Formula("=" + stringifyNode(n, -1, sheetName))
}

func stringifyNode(n Node, parentPrecedence int, sheetName string) string {
	// To solve the "excessive parenthesis" problem, see this:
	// https://stackoverflow.com/a/58679340/5989906
	switch n.Type() {
	case NodeTypeNumber, NodeTypeLogical, NodeTypeText:
		return n.(ValueNode).String()
	case NodeTypeFunction:
		fNode := n.(FunctionNode)
		args := make([]string, len(fNode.Arguments))
		for i, arg := range fNode.Arguments {
			args[i] = stringifyNode(arg, -1, sheetName)
		}
		argsWithCommas := strings.Join(args, ", ")
		return fmt.Sprintf("%s(%s)", fNode.Name, argsWithCommas)
	case NodeTypeBinaryExpression:
		bNode := n.(BinaryExpressionNode)
		// Deals with the precedence of the operators here
		return stringifyBinaryExp(bNode, parentPrecedence, sheetName)
	case NodeTypeUnaryExpression:
		uNode := n.(UnaryExpressionNode)
		if uNode.Operand.Type().IsTerminal() {
			return fmt.Sprintf(
				"%s%s",
				uNode.Operator,
				stringifyNode(uNode.Operand, -1, sheetName),
			)
		}
		return fmt.Sprintf(
			"%s(%s)",
			uNode.Operator,
			stringifyNode(uNode.Operand, -1, sheetName),
		)
	case NodeTypeCell:
		cNode := n.(CellNode)
		return string(cNode.Cell.ToAddressRel(sheetName))
	case NodeTypeCellRange:
		rNode := n.(CellRangeNode)
		return rNode.Range().StringRel(sheetName)
	default:
		// I know, not great, not terrible...
		return "ERROR_STRINGIFYING_NODE"
	}
}

func stringifyBinaryExp(b BinaryExpressionNode, parentPrecedence int, sheetName string) string {
	opPrecedence, ok := PrecedenceMap[b.Operator]
	if !ok {
		return "ERROR_STRINGIFYING_BINARY_EXP"
	}
	commu, ok := IsCommutative[b.Operator]
	if !ok {
		return "ERROR_STRINGIFYING_BINARY_EXP"
	}
	if !commu {
		left := stringifyNode(b.Left, opPrecedence, sheetName)
		right := stringifyNode(b.Right, opPrecedence+1, sheetName)
		res := left + b.Operator + right
		if parentPrecedence > opPrecedence {
			return "(" + res + ")"
		}
		return res
	} else {
		left := stringifyNode(b.Left, opPrecedence, sheetName)
		right := stringifyNode(b.Right, opPrecedence, sheetName)
		res := left + b.Operator + right
		if parentPrecedence > opPrecedence {
			return "(" + res + ")"
		}
		return res
	}
}
