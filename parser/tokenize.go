package parser

import "github.com/xuri/efp"

// Tokenizes a formula string into a slice of tokens,
// for later parsing into a tree.
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
