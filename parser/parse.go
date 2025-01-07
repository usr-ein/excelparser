package parser

// Parse takes a formula string and returns a Node representing
// the formula root.
// The currentSheet parameter is used to resolve relative cell references.
func Parse(formula string, currentSheet string) (Node, error) {
	tokens := Tokenize(formula)
	return BuildTree(Context{
		CurrentSheet: currentSheet,
	}, tokens)
}
