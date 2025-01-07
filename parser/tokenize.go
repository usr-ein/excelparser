package parser

import "github.com/xuri/efp"

func Tokenize(formula string) []Token {
	parser := efp.ExcelParser()
	rawTokens := parser.Parse(formula)

	tokens := make([]Token, len(rawTokens))
	for i, rawToken := range rawTokens {
		tokens[i] = Token{
			Type:    rawToken.TType,
			Subtype: rawToken.TSubType,
			Value:   rawToken.TValue,
		}
	}
	return tokens
}
