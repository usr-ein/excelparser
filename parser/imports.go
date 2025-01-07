package parser

import "github.com/usr-ein/excelparser/xl"

type Cell = xl.Cell
type Range = xl.Range
type Sheet = xl.Sheet
type RawSheet = xl.RawSheet
type Workbook = xl.Workbook
type CVal = xl.CVal
type Formula = xl.Formula

type CType = xl.CType

const CTEmpty = xl.CTEmpty
const CTString = xl.CTString
const CTFormula = xl.CTFormula
const CTNumber = xl.CTNumber
const CTBool = xl.CTBool
