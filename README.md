# excelparser: An Abstract Syntax Tree parser for Excel formulas

This repo is a port from JS to Go of the code from [psalaets/excel-formula-ast](https://github.com/psalaets/excel-formula-ast/tree/master).

I was trying to parse Excel formulas into ASTs to process them further, but no Go repo existed and did that in a simple and transparent way.

1. I started by porting the JS code line by line.
2. I then added unit tests and end to end tests on everything.
3. Finally, I removed/refactored shitty looking code until I was pleased with the structure.

## Other Excel/Go libraries

I found the following other useful repos:

- [xuri/xfp](https://github.com/xuri/efp) I use this to tokenize formulas here. From the doc, it claims to "*get an Abstract Syntax Tree (AST) from Excel formula*", but this is wrong, and all it does is pretty print the list of tokens.
- E. W. Bachtal's Excel formula parser, see the [/bachtal folder](bachtal/)
- The excelent [qax-os/excelize](https://github.com/qax-os/excelize). Great library for general Excel things, but not so much for formula manipulation
- [tealeg/xlsx](https://github.com/tealeg/xlsx) gave me a few headaches, but overall good library. Some limitation regarding styling of cells though, and underdocumented in many places, but higher level than excelize.
