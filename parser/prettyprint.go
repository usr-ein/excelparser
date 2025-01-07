package parser

import (
	"fmt"
)

type NodeJSON struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

func ToNodeJson(n Node) NodeJSON {
	switch n.Type() {
	case NodeTypeNumber, NodeTypeText, NodeTypeLogical, NodeTypeCell, NodeTypeCellRange:
		return NodeJSON{
			Type:  n.Type().String(),
			Value: getLabel(n),
		}
	case NodeTypeFunction:
		funcNode := n.(FunctionNode)
		args := make([]NodeJSON, len(funcNode.Arguments))
		for i, arg := range funcNode.Arguments {
			args[i] = ToNodeJson(arg)
		}
		return NodeJSON{
			Type:  n.Type().String() + " " + funcNode.Name,
			Value: args,
		}
	case NodeTypeBinaryExpression:
		binNode := n.(BinaryExpressionNode)
		return NodeJSON{
			Type: n.Type().String() + " " + binNode.Operator,
			Value: []NodeJSON{
				ToNodeJson(binNode.Left),
				ToNodeJson(binNode.Right),
			},
		}
	case NodeTypeUnaryExpression:
		unaryNode := n.(UnaryExpressionNode)
		return NodeJSON{
			Type:  n.Type().String() + " " + unaryNode.Operator,
			Value: ToNodeJson(unaryNode.Operand),
		}
	default:
		return NodeJSON{
			Type:  "Unknown",
			Value: "Unknown",
		}
	}
}

func getLabel(node Node) string {
	switch node.Type() {
	case NodeTypeFunction:
		return node.(FunctionNode).Name
	case NodeTypeBinaryExpression:
		return node.(BinaryExpressionNode).Operator
	case NodeTypeNumber:
		return fmt.Sprintf("%f", node.(NumberNode).Value)
	case NodeTypeText:
		return node.(TextNode).Value
	case NodeTypeLogical:
		return fmt.Sprintf("%t", node.(LogicalNode).Value)
	case NodeTypeCell:
		return string(node.(CellNode).Cell.ToAddress())
	case NodeTypeCellRange:
		rightCell, err := node.(CellRangeNode).End.Cell.Shift(-1, -1)
		if err != nil {
			rightCell = node.(CellRangeNode).End.Cell
		}
		return fmt.Sprintf("%s:%s", node.(CellRangeNode).Start.Cell.ToAddress(), rightCell.ToAddressNoSheet())
	default:
		return "Unknown node type"
	}
}
