# excelparser: An Abstract Syntax Tree parser for Excel formulas

This repo is a port from JS to Go of the code from [psalaets/excel-formula-ast](https://github.com/psalaets/excel-formula-ast/tree/master).

I was trying to parse Excel formulas into ASTs to process them further, but no Go repo existed and did that in a simple and transparent way.

1. I started by porting the JS code line by line.
2. I then added unit tests and end to end tests on everything.
3. Finally, I removed/refactored shitty looking code until I was pleased with the structure.

## How to get it

```sh
go get github.com/usr-ein/excelparser
```

[Doc is available online here.](https://pkg.go.dev/github.com/usr-ein/excelparser@latest#section-readme)

## What you get using this repo

Here is an example of the kind of data structure you get by using this lib:

```go
import (
    "fmt"
    "github.com/usr-ein/excelparser/parser"
)
// Node is the interface that all nodes in the AST implement.
// It represent a node such as a SUM function, a cell reference, a + binary expression, etc.
type Node interface {
	Type() NodeType
	IsEq(Node) bool
	Children() []Node
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


func main(){
    node, err := parser.Parse(`(SUM(A1:B1, A2:B2)+A5+3232-A15)/B90/321.0+MONTH("January")`, `Sheet1`)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", n)
    // The large formula is broken down into nodes with fields like left/right/operator etc.
    /**
    {
        Operator:+
        Left:{
            Operator:/ 
            Left:{
                Operator:/ 
                Left:{
                    Operator:- 
                    Left:{
                        Operator:+ 
                        Left:{
                            Operator:+ 
                            Left:{
                                Name:SUM 
                                Arguments:[
                                    {
                                        Start:
                                            {Cell:{Sheet:Sheet1 Row:0 Col:0 RowRel:true ColRel:true}} 
                                        End:{Cell:{Sheet:Sheet1 Row:1 Col:2 RowRel:true ColRel:true}}
                                    }
                                    {
                                        Start:
                                            {Cell:{Sheet:Sheet1 Row:1 Col:0 RowRel:true ColRel:true}}
                                        End:{Cell:{Sheet:Sheet1 Row:2 Col:2 RowRel:true ColRel:true}}
                                    }
                                ]
                            }
                            Right:{
                                Cell:{Sheet:Sheet1 Row:4 Col:0 RowRel:true ColRel:true}
                            }
                        }
                        Right:3232
                    }
                    Right:{
                        Cell:{Sheet:Sheet1 Row:14 Col:0 RowRel:true ColRel:true}
                    }
                }
                Right:{
                    Cell:{Sheet:Sheet1 Row:89 Col:1 RowRel:true ColRel:true}
                }
            }
            Right:321
        }
        Right: {
            Name:MONTH
            Arguments:["January"]
        }
    }
    */
}

```

## Other Excel/Go libraries

I found the following other useful repos:

- [xuri/xfp](https://github.com/xuri/efp) I use this to tokenize formulas here. From the doc, it claims to "*get an Abstract Syntax Tree (AST) from Excel formula*", but this is wrong, and all it does is pretty print the list of tokens.
- E. W. Bachtal's Excel formula parser, see the [/_docs/bachtal folder](_docs/bachtal/)
- The excelent [qax-os/excelize](https://github.com/qax-os/excelize). Great library for general Excel things, but not so much for formula manipulation
- [tealeg/xlsx](https://github.com/tealeg/xlsx) gave me a few headaches, but overall good library. Some limitation regarding styling of cells though, and underdocumented in many places, but higher level than excelize.
