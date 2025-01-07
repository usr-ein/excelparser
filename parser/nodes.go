package parser

import (
	"fmt"
	"math"
	"strconv"
)

// Node is the interface that all nodes in the AST implement.
// It represent a node such as a SUM function, a cell reference, a + binary expression, etc.
type Node interface {
	Type() NodeType
	IsEq(Node) bool
	Children() []Node
}

type ValueNode interface {
	Node
	String() string
}

type NodeType uint8

const (
	NodeTypeCell NodeType = iota
	NodeTypeCellRange
	NodeTypeFunction
	NodeTypeBinaryExpression
	NodeTypeUnaryExpression
	NodeTypeNumber
	NodeTypeText
	NodeTypeLogical
)

func (NodeType NodeType) IsTerminal() bool {
	return NodeType == NodeTypeNumber || NodeType == NodeTypeText || NodeType == NodeTypeLogical || NodeType == NodeTypeCell || NodeType == NodeTypeCellRange
}

func (nodeType NodeType) String() string {
	switch nodeType {
	case NodeTypeNumber:
		return "num"
	case NodeTypeText:
		return "txt"
	case NodeTypeLogical:
		return "bool"
	case NodeTypeCell:
		return "cell"
	case NodeTypeCellRange:
		return "range"
	case NodeTypeFunction:
		return "func"
	case NodeTypeBinaryExpression:
		return "binExp"
	case NodeTypeUnaryExpression:
		return "unaExp"
	default:
		return "Unknown"
	}
}

type CellNode struct {
	Cell Cell `json:"cell"`
}

func (c CellNode) Type() NodeType {
	return NodeTypeCell
}

func (c CellNode) IsEq(n Node) bool {
	if n.Type() != NodeTypeCell {
		return false
	}
	cellNode, ok := n.(CellNode)
	if !ok {
		return false
	}
	return c.Cell.IsEq(cellNode.Cell)
}

func (c CellNode) Children() []Node {
	return []Node{}
}

type CellRangeNode struct {
	// Make sure that start and end are in the same sheet!
	Start CellNode `json:"startCell"`
	End   CellNode `json:"endCell"`
}

func (c CellRangeNode) Type() NodeType {
	return NodeTypeCellRange
}

func (c CellRangeNode) IsEq(n Node) bool {
	if n.Type() != NodeTypeCellRange {
		return false
	}
	cellNode, ok := n.(CellRangeNode)
	if !ok {
		return false
	}
	return c.Start.IsEq(cellNode.Start) && c.End.IsEq(cellNode.End)
}

func (c CellRangeNode) Range() Range {
	return Range{
		Start: c.Start.Cell,
		// End was already shifted by the parser
		// so we're good to go here.
		End: c.End.Cell,
	}
}

func (c CellRangeNode) Children() []Node {
	return []Node{}
}

type FunctionNode struct {
	Name      string `json:"name"`
	Arguments []Node `json:"arguments"`
}

func (f FunctionNode) Type() NodeType {
	return NodeTypeFunction
}

func (f FunctionNode) IsEq(n Node) bool {
	if n.Type() != NodeTypeFunction {
		return false
	}
	if f.Name != n.(FunctionNode).Name {
		return false
	}
	if len(f.Arguments) != len(n.(FunctionNode).Arguments) {
		return false
	}
	for i, arg := range f.Arguments {
		if !arg.IsEq(n.(FunctionNode).Arguments[i]) {
			return false
		}
	}
	return true
}

func (c FunctionNode) Children() []Node {
	args := make([]Node, len(c.Arguments))
	copy(args, c.Arguments)
	return args
}

type NumberNode struct {
	Value float64 `json:"value"`
}

func (n NumberNode) Type() NodeType {
	return NodeTypeNumber
}

func (n NumberNode) IsEq(node Node) bool {
	if node.Type() != NodeTypeNumber {
		return false
	}
	return n.Value == node.(NumberNode).Value
}

func (n NumberNode) Children() []Node {
	return []Node{}
}

func (n NumberNode) String() string {
	num := n.Value
	if num == math.Trunc(num) {
		return strconv.FormatInt(int64(num), 10)
	}
	return fmt.Sprintf("%.4f", num)
}

type TextNode struct {
	Value string `json:"value"`
}

func (t TextNode) Type() NodeType {
	return NodeTypeText
}

func (t TextNode) IsEq(node Node) bool {
	if node.Type() != NodeTypeText {
		return false
	}
	return t.Value == node.(TextNode).Value
}

func (t TextNode) Children() []Node {
	return []Node{}
}

func (t TextNode) String() string {
	return "\"" + t.Value + "\""
}

type LogicalNode struct {
	Value bool `json:"value"`
}

func (l LogicalNode) Type() NodeType {
	return NodeTypeLogical
}

func (l LogicalNode) IsEq(node Node) bool {
	if node.Type() != NodeTypeLogical {
		return false
	}
	return l.Value == node.(LogicalNode).Value
}

func (l LogicalNode) Children() []Node {
	return []Node{}
}

func (l LogicalNode) String() string {
	if l.Value {
		return "TRUE"
	}
	return "FALSE"
}

type BinaryExpressionNode struct {
	Operator string `json:"operator"`
	Left     Node   `json:"left"`
	Right    Node   `json:"right"`
}

func (b BinaryExpressionNode) Type() NodeType {
	return NodeTypeBinaryExpression
}

func (b BinaryExpressionNode) IsEq(node Node) bool {
	if node.Type() != NodeTypeBinaryExpression {
		return false
	}
	return b.Operator == node.(BinaryExpressionNode).Operator && b.Left.IsEq(node.(BinaryExpressionNode).Left) && b.Right.IsEq(node.(BinaryExpressionNode).Right)
}

func (b BinaryExpressionNode) Children() []Node {
	return []Node{b.Left, b.Right}
}

type UnaryExpressionNode struct {
	Operator string `json:"operator"`
	Operand  Node   `json:"operand"`
}

func (u UnaryExpressionNode) Type() NodeType {
	return NodeTypeUnaryExpression
}

func (u UnaryExpressionNode) IsEq(node Node) bool {
	if node.Type() != NodeTypeUnaryExpression {
		return false
	}
	return u.Operator == node.(UnaryExpressionNode).Operator && u.Operand.IsEq(node.(UnaryExpressionNode).Operand)
}

func (u UnaryExpressionNode) Children() []Node {
	return []Node{u.Operand}
}
