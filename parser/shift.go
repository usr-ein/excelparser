package parser

import (
	"github.com/pkg/errors"
)

func MoveNode(n Node, origin Cell, dest Cell) (Node, error) {
	if origin.Sheet != dest.Sheet {
		return nil, errors.New("cannot move formula between sheets")
	}
	return ShiftNode(
		n,
		int(dest.Row)-int(origin.Row),
		int(dest.Col)-int(origin.Col),
	)
}

func ShiftNode(n Node, shiftRow int, shiftCol int) (Node, error) {
	switch n.Type() {
	case NodeTypeNumber, NodeTypeText, NodeTypeLogical:
		return n, nil
	case NodeTypeFunction:
		fNode := n.(FunctionNode)
		args := make([]Node, len(fNode.Arguments))
		for i, arg := range fNode.Arguments {
			shiftedArg, err := ShiftNode(arg, shiftRow, shiftCol)
			if err != nil {
				return nil, err
			}
			args[i] = shiftedArg
		}
		return FunctionNode{
			Name:      fNode.Name,
			Arguments: args,
		}, nil
	case NodeTypeBinaryExpression:
		bNode := n.(BinaryExpressionNode)
		shiftedLeft, err := ShiftNode(bNode.Left, shiftRow, shiftCol)
		if err != nil {
			return nil, err
		}
		shiftedRight, err := ShiftNode(bNode.Right, shiftRow, shiftCol)
		if err != nil {
			return nil, err
		}
		return BinaryExpressionNode{
			Left:     shiftedLeft,
			Operator: bNode.Operator,
			Right:    shiftedRight,
		}, nil
	case NodeTypeUnaryExpression:
		uNode := n.(UnaryExpressionNode)
		shiftedOperand, err := ShiftNode(uNode.Operand, shiftRow, shiftCol)
		if err != nil {
			return nil, err
		}
		return UnaryExpressionNode{
			Operator: uNode.Operator,
			Operand:  shiftedOperand,
		}, nil
	case NodeTypeCell:
		cNode := n.(CellNode)
		shiftedCell, err := cNode.Cell.ShiftIfRel(shiftRow, shiftCol)
		if err != nil {
			return nil, err
		}
		return CellNode{
			Cell: shiftedCell,
		}, nil
	case NodeTypeCellRange:
		rNode := n.(CellRangeNode)
		shiftedRange, err := rNode.Range().ShiftIfRel(shiftRow, shiftCol)
		if err != nil {
			return nil, err
		}
		return CellRangeNode{
			Start: CellNode{Cell: shiftedRange.Start},
			End:   CellNode{Cell: shiftedRange.End},
		}, nil
	}
	return nil, errors.New("unknown node type")

}

// Shifts a formula from one cell to another, inside the same sheet.
//
// Deprecated: This is very slow, since we now do it in three steps:
// 1. Parse the formula into a tree.
// 2. Shift the tree.
// 3. Serialize the tree back into a formula.
// Use Parse/ShiftNode/MoveNode/StringifyNode instead depending on your needs.
func ShiftFormula(f Formula, rowDiff int, colDiff int, sheetName string) (Formula, error) {
	node, err := Parse(string(f), sheetName)
	if err != nil {
		return "", err
	}
	shiftedNode, err := ShiftNode(node, rowDiff, colDiff)
	if err != nil {
		return "", err
	}
	return StringifyNode(shiftedNode, sheetName), nil
}
