package parser

import (
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/usr-ein/excelparser/parser/shuntingyard"
	"github.com/usr-ein/excelparser/xl"
)

type ShuntingYard = shuntingyard.ShuntingYardState[Node]

// PrecedenceMap is a map of binary operators to their precedence.
// It lets us know which operators should be evaluated first.
var PrecedenceMap = map[string]int{
	// cell range union and intersect
	" ": 8,
	",": 8,
	// raise to power
	"^": 5,
	// multiply, divide
	"*": 4,
	"/": 4,
	// add, subtract
	"+": 3,
	"-": 3,
	// string concat
	"&": 2,
	// comparison
	"=":  1,
	"<>": 1,
	"<=": 1,
	">=": 1,
	">":  1,
	"<":  1,
}

// IsCommutative is a map of binary operators to whether they are commutative.
// True if A OP B == B OP A
// False if A OP B != B OP A
var IsCommutative = map[string]bool{
	// cell range union and intersect
	" ": true,
	",": false,
	// raise to power
	"^": false,
	// multiply, divide
	"*": true,
	"/": false,
	// add, subtract
	"+": true,
	"-": false,
	// string concat
	"&": false,
	// comparison
	"=":  true,
	"<>": true,
	"<=": false,
	">=": false,
	">":  false,
	"<":  false,
}

// BuildTree takes a slice of tokens and returns a Node representing
// the formula root.
// The context parameter is used to resolve relative cell references,
// and contains the current sheet name, and other context information necessary.
func BuildTree(ctx Context, tokens []Token) (Node, error) {
	// named parseFormula in original
	stream := NewTokenStream(tokens)
	shuntingYard := shuntingyard.NewShuntingYardState[Node]()

	if err := parseExpression(ctx, stream, shuntingYard); err != nil {
		return nil, err
	}

	retVal, ok := shuntingYard.Operands.Top()
	if !ok {
		return nil, errors.New("no top operand found after parsing formula")
	}
	return retVal, nil
}

func parseExpression(ctx Context, stream TokenStream, shuntingYard ShuntingYard) error {
	if err := parseOperandExpression(ctx, stream, shuntingYard); err != nil {
		return err
	}

	var pos int
	for {
		if !stream.NextIsBinaryOperator() {
			break
		}
		if pos == stream.Position() {
			return errors.New("position didn't change after parsing binary operator")
		}
		pos = stream.Position()
		binaryOperator, err := createBinaryOperator(stream.GetNext().Value)
		if err != nil {
			return err
		}
		if err := pushOperator(binaryOperator, shuntingYard); err != nil {
			return err
		}
		if err := stream.Consume(); err != nil {
			return err
		}
		if err := parseOperandExpression(ctx, stream, shuntingYard); err != nil {
			return err
		}
	}

	for {
		top, _ := shuntingYard.Operators.Top()
		if top == shuntingyard.Sentinel {
			break
		}
		if err := popOperator(shuntingYard); err != nil {
			return err
		}
	}
	return nil
}

func parseOperandExpression(ctx Context, stream TokenStream, shuntingYard ShuntingYard) error {
	if stream.NextIsTerminal() {
		operand, err := parseTerminal(ctx, stream)
		if err != nil {
			return errors.Wrap(err, "failed to parse terminal")
		}
		shuntingYard.Operands.Push(operand)
		// parseTerminal already consumes once so don't need to consume on line below
		// stream.consume()
	} else if stream.NextIsOpenParen() {
		// open paren
		if err := stream.Consume(); err != nil {
			return errors.Wrap(err, "failed to consume open paren")
		}
		if err := withinSentinel(shuntingYard, func() error {
			return parseExpression(ctx, stream, shuntingYard)
		}); err != nil {
			return errors.Wrap(err, "failed to parse subexpression within sentinel")
		}
		// close paren
		if err := stream.Consume(); err != nil {
			return errors.Wrap(err, "failed to consume close paren")
		}
	} else if stream.NextIsPrefixOperator() {
		unaryOperator, err := createUnaryOperator(stream.GetNext().Value)
		if err != nil {
			return errors.Wrap(err, "failed to create unary operator")
		}
		if err := pushOperator(unaryOperator, shuntingYard); err != nil {
			return errors.Wrap(err, "failed to push unary operator")
		}
		if err := stream.Consume(); err != nil {
			return errors.Wrap(err, "failed to consume prefix operator")
		}
		if err := parseOperandExpression(ctx, stream, shuntingYard); err != nil {
			return errors.Wrap(err, "failed to parse operand expression")
		}
	} else if stream.NextIsFunctionCall() {
		if err := parseFunctionCall(ctx, stream, shuntingYard); err != nil {
			return errors.Wrap(err, "failed to parse function call")
		}
	}
	// TODO: Should be error here or should we do nothing?
	// Original implem does nothing...
	return nil
}

func parseFunctionCall(ctx Context, stream TokenStream, shuntingYard ShuntingYard) error {
	name := stream.GetNext().Value
	// consume start of function call
	if err := stream.Consume(); err != nil {
		return errors.Wrap(err, "failed to consume start of function call")
	}
	args, err := parseFunctionArgList(ctx, stream, shuntingYard)
	if err != nil {
		return errors.Wrap(err, "failed to parse function arg list")
	}
	shuntingYard.Operands.Push(FunctionNode{
		Name:      name,
		Arguments: args,
	})

	// consume end of function call
	if err := stream.Consume(); err != nil {
		return errors.Wrap(err, "failed to consume end of function call")
	}
	return nil
}

func parseFunctionArgList(ctx Context, stream TokenStream, shuntingYard ShuntingYard) ([]Node, error) {
	reverseArgs := make([]Node, 0)
	if err := withinSentinel(shuntingYard, func() error {
		// I don't like this JavaScript-y way of
		// having a closure binding on reverseArgs... oh well
		arity := 0
		var pos int
		for {
			if stream.NextIsEndOfFunctionCall() {
				break
			}
			if pos == stream.Position() {
				return errors.New("position didn't change after parsing function argument")
			}
			pos = stream.Position()
			if err := parseExpression(ctx, stream, shuntingYard); err != nil {
				return errors.Wrap(err, "failed to parse expression in func arg list")
			}

			arity += 1

			if stream.NextIsFunctionArgumentSeparator() {
				if err := stream.Consume(); err != nil {
					return errors.Wrap(err, "failed to consume function argument separator")
				}
			}
		}

		for i := 0; i < arity; i++ {
			arg, ok := shuntingYard.Operands.Pop()
			if !ok {
				return errors.New("failed to pop operand from stack for function argument list")
			}
			reverseArgs = append(reverseArgs, arg)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to parse function argument list")
	}

	slices.Reverse(reverseArgs)
	return reverseArgs, nil
}

type AnyClosure func() error

func withinSentinel(shuntingYard ShuntingYard, closure AnyClosure) error {
	shuntingYard.Operators.Push(shuntingyard.Sentinel)
	if err := closure(); err != nil {
		return errors.Wrap(err, "failed to run closure within sentinel")
	}
	if _, ok := shuntingYard.Operators.Pop(); !ok {
		return errors.New("failed to pop sentinel operator")
	}
	return nil
}

func pushOperator(operator shuntingyard.Operator, shuntingYard ShuntingYard) error {
	// for top, ok := shuntingYard.Operators.Top(); ok && top.EvaluatesBefore(operator); {
	for {
		top, ok := shuntingYard.Operators.Top()
		if !ok {
			return errors.New("failed to get top operator from stack while pushing new operator")
		}
		if !top.EvaluatesBefore(operator) {
			break
		}
		if err := popOperator(shuntingYard); err != nil {
			return errors.Wrap(err, "failed to pop operator while pushing new operator")
		}
	}
	shuntingYard.Operators.Push(operator)
	return nil
}

func popOperator(shuntingYard ShuntingYard) error {
	top, ok := shuntingYard.Operators.Top()
	if !ok {
		return errors.New("failed to get top operator from stack while poping")
	}

	if top.IsBinary() {
		right, ok := shuntingYard.Operands.Pop()
		if !ok {
			return errors.New("failed to pop right operand from stack")
		}
		left, ok := shuntingYard.Operands.Pop()
		if !ok {
			return errors.New("failed to pop left operand from stack")
		}
		operator, ok := shuntingYard.Operators.Pop()
		if !ok {
			return errors.New("failed to pop binary operator from stack")
		}
		shuntingYard.Operands.Push(BinaryExpressionNode{
			Operator: operator.Symbol(),
			Left:     left,
			Right:    right,
		})
	} else if top.IsUnary() {
		operand, ok := shuntingYard.Operands.Pop()
		if !ok {
			return errors.New("failed to pop operand from stack")
		}
		operator, ok := shuntingYard.Operators.Pop()
		if !ok {
			return errors.New("failed to pop unary operator from stack")
		}
		shuntingYard.Operands.Push(UnaryExpressionNode{
			Operator: operator.Symbol(),
			Operand:  operand,
		})
	}
	return nil
}

func parseTerminal(ctx Context, stream TokenStream) (Node, error) {
	if stream.NextIsNumber() {
		return parseNumber(stream)
	}
	if stream.NextIsText() {
		return parseText(stream)
	}
	if stream.NextIsLogical() {
		return parseLogical(stream)
	}
	if stream.NextIsCell() {
		return parseCell(ctx, stream)
	}
	if stream.NextIsRange() {
		return parseRange(ctx, stream)
	}
	return nil, errors.New("invalid terminal")
}

func parseCell(ctx Context, stream TokenStream) (CellNode, error) {
	next := stream.GetNext()
	if err := stream.Consume(); err != nil {
		return CellNode{}, err
	}
	if strings.Contains(next.Value, ":") {
		return CellNode{}, errors.New("cell contains range")
	}
	cell, err := xl.ParseCell(next.Value, ctx.CurrentSheet)
	if err != nil {
		return CellNode{}, errors.Wrap(err, "failed to parse cell")
	}
	return CellNode{Cell: cell}, nil
}

func parseRange(ctx Context, stream TokenStream) (CellRangeNode, error) {
	next := stream.GetNext()
	if err := stream.Consume(); err != nil {
		return CellRangeNode{}, err
	}

	leftRight := strings.Split(next.Value, ":")
	if len(leftRight) != 2 {
		return CellRangeNode{}, errors.New("invalid range")
	}
	start, err := xl.ParseCell(leftRight[0], ctx.CurrentSheet)
	if err != nil {
		return CellRangeNode{}, errors.Wrap(err, "failed to parse start cell")
	}
	end, err := xl.ParseCell(leftRight[1], start.Sheet)
	if err != nil {
		return CellRangeNode{}, errors.Wrap(err, "failed to parse end cell")
	}
	end, err = end.Shift(1, 1)
	if err != nil {
		return CellRangeNode{}, errors.Wrap(err, "failed to shift end cell")
	}
	end.Sheet = start.Sheet
	// TODO: add test that Sheet2!A1:B2 with current sheet=Sheet1 ends with
	// start: {A1, Sheet2}, end: {B2, Sheet2} and not
	// start: {A1, Sheet2}, end: {B2, Sheet1}
	// ---
	// Also, Sheet1!A1:Sheet2!B2 is not valid!!

	return CellRangeNode{
		Start: CellNode{Cell: start},
		End:   CellNode{Cell: end},
	}, nil
}

func parseText(stream TokenStream) (TextNode, error) {
	next := stream.GetNext()
	if err := stream.Consume(); err != nil {
		return TextNode{}, errors.Wrap(err, "failed to consume text token")
	}
	return TextNode{Value: next.Value}, nil
}

func parseLogical(stream TokenStream) (LogicalNode, error) {
	next := stream.GetNext()
	if err := stream.Consume(); err != nil {
		return LogicalNode{}, errors.Wrap(err, "failed to consume logical token")
	}
	return LogicalNode{Value: next.Value == "TRUE"}, nil
}

func parseNumber(stream TokenStream) (NumberNode, error) {
	next := stream.GetNext()
	value, err := strconv.ParseFloat(next.Value, 64)
	if err != nil {
		return NumberNode{}, errors.Wrap(err, "failed to parse number")
	}
	if err = stream.Consume(); err != nil {
		return NumberNode{}, errors.Wrap(err, "failed to consume number token")
	}
	if stream.NextIsPostfixOperator() {
		value *= 0.01
		if err = stream.Consume(); err != nil {
			return NumberNode{}, errors.Wrap(err, "failed to consume postfix number operator")
		}
	}
	return NumberNode{Value: value}, nil
}

func createUnaryOperator(symbol string) (shuntingyard.Operator, error) {
	precedenceMap := map[string]int{
		// negation
		"-": 7,
	}
	precedence, ok := precedenceMap[symbol]
	if !ok {
		return nil, errors.New("invalid unary operator")
	}
	return shuntingyard.New(symbol, precedence, 1, true), nil
}

func createBinaryOperator(symbol string) (shuntingyard.Operator, error) {
	precedence, ok := PrecedenceMap[symbol]
	if !ok {
		return nil, errors.New("invalid binary operator")
	}
	return shuntingyard.New(symbol, precedence, 2, true), nil
}
