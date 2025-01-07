package parser

func Parse(formula string, currentSheet string) (Node, error) {
	tokens := Tokenize(formula)
	return BuildTree(Context{
		CurrentSheet: currentSheet,
	}, tokens)
}
