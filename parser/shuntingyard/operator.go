package shuntingyard

/**
https://github.com/psalaets/excel-formula-ast/blob/master/lib/shunting-yard/operator.js
*/

type Operator interface {
	IsUnary() bool
	IsBinary() bool
	EvaluatesBefore(Operator) bool

	Precedence() int
	IsSentinel() bool
	Symbol() string
}

type OperatorImpl struct {
	symbol          string
	precedence      int
	operandCount    int
	leftAssociative bool
}

func NewWithDefault(symbol string, precedence int) Operator {
	return New(symbol, precedence, 2, true)
}

func New(symbol string, precedence int, operandCount int, leftAssociative bool) Operator {
	return OperatorImpl{
		symbol:          symbol,
		precedence:      precedence,
		operandCount:    operandCount,
		leftAssociative: leftAssociative,
	}
}

func (o OperatorImpl) IsUnary() bool {
	return o.operandCount == 1
}

func (o OperatorImpl) IsBinary() bool {
	return o.operandCount == 2
}

func (o OperatorImpl) IsSentinel() bool {
	return o == Sentinel
}

func (o OperatorImpl) Precedence() int {
	return o.precedence
}

func (o OperatorImpl) Symbol() string {
	return o.symbol
}

func (thisOp OperatorImpl) EvaluatesBefore(other Operator) bool {
	if thisOp.IsSentinel() {
		return false
	}
	if other.IsSentinel() {
		return true
	}
	if other.IsUnary() {
		return false
	}

	if thisOp.IsUnary() {
		return thisOp.Precedence() >= other.Precedence()
	} else if thisOp.IsBinary() {
		if thisOp.Precedence() == other.Precedence() {
			return thisOp.leftAssociative
		} else {
			return thisOp.Precedence() > other.Precedence()
		}
	}

	return false
}

var Sentinel Operator = NewWithDefault("S", 0)
