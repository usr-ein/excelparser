package parser

import (
	"errors"
	"strings"
)

// https://github.com/psalaets/excel-formula-ast/blob/master/lib/token-stream.js

type Token struct {
	Type    string
	Subtype string
	Value   string
}

type TokenStream interface {
	Consume() error
	GetNext() Token
	NextIs(ttype string, tsubtype string) bool
	NextIsOpenParen() bool
	NextIsTerminal() bool
	NextIsFunctionCall() bool
	NextIsFunctionArgumentSeparator() bool
	NextIsEndOfFunctionCall() bool
	NextIsBinaryOperator() bool
	NextIsPrefixOperator() bool
	NextIsPostfixOperator() bool
	NextIsRange() bool
	NextIsCell() bool
	NextIsNumber() bool
	NextIsText() bool
	NextIsLogical() bool
	Position() int
}

type TokenStreamImpl struct {
	tokens   []Token
	position int
}

func NewTokenStream(tokens []Token) TokenStream {
	tokensArr := make([]Token, len(tokens)+1)
	copy(tokensArr, tokens)
	tokensArr[len(tokens)] = Token{}
	return &TokenStreamImpl{
		tokens:   tokensArr,
		position: 0,
	}
}

func (ts *TokenStreamImpl) Consume() error {
	ts.position++
	if ts.position >= len(ts.tokens) {
		return errors.New("invalid syntax")
	}
	return nil
}

func (ts *TokenStreamImpl) GetNext() Token {
	return ts.tokens[ts.position]
}

func (ts *TokenStreamImpl) NextIs(ttype string, tsubtype string) bool {
	if ts.GetNext().Type != ttype {
		return false
	}
	if tsubtype != "" && ts.GetNext().Subtype != tsubtype {
		return false
	}
	return true
}

func (ts *TokenStreamImpl) NextIsOpenParen() bool {
	return ts.NextIs("Subexpression", "Start")
}

func (ts *TokenStreamImpl) NextIsTerminal() bool {
	return ts.NextIsNumber() || ts.NextIsText() || ts.NextIsRange() || ts.NextIsCell() || ts.NextIsLogical()
}

func (ts *TokenStreamImpl) NextIsFunctionCall() bool {
	return ts.NextIs("Function", "Start")
}

func (ts *TokenStreamImpl) NextIsFunctionArgumentSeparator() bool {
	return ts.NextIs("Argument", "")
}

func (ts *TokenStreamImpl) NextIsEndOfFunctionCall() bool {
	return ts.NextIs("Function", "Stop")
}

func (ts *TokenStreamImpl) NextIsBinaryOperator() bool {
	return ts.NextIs("OperatorInfix", "")
}

func (ts *TokenStreamImpl) NextIsPrefixOperator() bool {
	return ts.NextIs("OperatorPrefix", "")
}

func (ts *TokenStreamImpl) NextIsPostfixOperator() bool {
	return ts.NextIs("OperatorPostfix", "")
}

func (ts *TokenStreamImpl) NextIsRange() bool {
	return ts.NextIs("Operand", "Range") && strings.Contains(ts.GetNext().Value, ":")
}

func (ts *TokenStreamImpl) NextIsCell() bool {
	return ts.NextIs("Operand", "Range") && !strings.Contains(ts.GetNext().Value, ":")
}

func (ts *TokenStreamImpl) NextIsNumber() bool {
	return ts.NextIs("Operand", "Number")
}

func (ts *TokenStreamImpl) NextIsText() bool {
	return ts.NextIs("Operand", "Text")
}

func (ts *TokenStreamImpl) NextIsLogical() bool {
	return ts.NextIs("Operand", "Logical")
}

func (ts *TokenStreamImpl) Position() int {
	return ts.position
}
