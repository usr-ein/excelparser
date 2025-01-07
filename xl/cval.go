package xl

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CellVal
// Represents the value of a cell in a sheet.
// The json tags are minimal to reduce file size
//
// The computed value can be a string, number, or bool, but not a formula, and so
// its value can be found in the other fields. E.g. the computed CVal of =SUM(A1:A2) would be
// Before computation:
//
//	{
//		Type: CTFormula,
//		HasComputed: false,
//		ComputedType: CTFormula,
//	    ValFormula: "=SUM(A1:A2)",
//	}
//
// After computation:
//
//	{
//		Type: CTFormula,
//		ValNumber: 3,
//		HasComputed: true,
//		ComputedType: CTNumber,
//	    ValFormula: "=SUM(A1:A2)",
//	}
//
// Another possibility is the spill/empty cells.
// let's say A1:B2 contains [[1,2],[3,4]]; D4 is empty, and C3 contains =A1:B2,
// then D4 will contain the spill value of B2, which is 4. In this case, we will have the following CVal in D4:
//
//	{
//		Type: CTEmpty,
//		ValNumber: 4,
//		HasComputed: true,
//		ComputedType: CTNumber,
//	}
//
// aka, an empty cell has a computed number value of 4! This is because of the formula =A1:B2 spilling into D4.
type CVal struct {
	Type CType `json:"-"`

	ValString string `json:"str,omitempty"`
	// We don't need the precision of float64
	ValNumber float32 `json:"num,omitempty"`
	ValBool   bool    `json:"bool,omitempty"`

	// If a formula, this is the formula string
	ValFormula Formula `json:"formula,omitempty"`
	// If the formula underwent computation, this is true
	HasComputed bool `json:"-"`
	// If the formula underwent computation, this is the type of the computed value.
	ComputedType CType `json:"-"`
}

var CValEmpty = CVal{Type: CTEmpty}

type CType uint8

const (
	CTEmpty CType = iota
	CTString
	CTFormula
	CTNumber
	CTBool
)

func (c CType) String() string {
	switch c {
	case CTString:
		return "string"
	case CTFormula:
		return "formula"
	case CTEmpty:
		return "empty"
	case CTNumber:
		return "number"
	case CTBool:
		return "bool"
	default:
		return "unknown"
	}
}

func (c CVal) Computed() CVal {
	if c.HasComputed {
		return CVal{
			Type:         c.ComputedType,
			ValString:    c.ValString,
			ValNumber:    c.ValNumber,
			ValBool:      c.ValBool,
			HasComputed:  false,
			ComputedType: CTEmpty,
		}
	}
	return c
}

func makeContent(raw [][]any, computed [][]any) ([][]CVal, error) {
	if !isRectangular(raw) {
		return nil, errors.New("content is not rectangular")
	}
	content := make([][]CVal, len(raw))

	hasComputed := true
	if len(computed) == 0 || len(computed[0]) == 0 {
		hasComputed = false
	} else {
		if !isRectangular(computed) {
			return nil, errors.New("computed values are not rectangular")
		}
		if len(computed) != len(raw) || len(computed[0]) != len(raw[0]) {
			return nil, errors.New("computed values are not the same size as content")
		}
	}
	var computedVal any = nil
	for i, row := range raw {
		content[i] = make([]CVal, len(row))
		for j, rawCell := range row {
			if hasComputed {
				computedVal = computed[i][j]
			}
			inputVal, err := makeCellVal(rawCell, computedVal)
			if err != nil {
				return nil, err
			}
			content[i][j] = inputVal
		}
	}
	return content, nil
}

func makeCellVal(raw any, computed any) (CVal, error) {
	if raw == nil {
		return CVal{Type: CTEmpty}, nil
	}
	switch val := raw.(type) {
	case float64:
		return CVal{Type: CTNumber, ValNumber: float32(val)}, nil
	case int:
		return CVal{Type: CTNumber, ValNumber: float32(val)}, nil
	case float32:
		return CVal{Type: CTNumber, ValNumber: val}, nil
	case string:
		if val == "" {
			return CVal{Type: CTEmpty}, nil
		}
		lowerVal := strings.ToLower(val)
		if lowerVal == "true" {
			return CVal{Type: CTBool, ValBool: true}, nil
		}
		if lowerVal == "false" {
			return CVal{Type: CTBool, ValBool: false}, nil
		}
		if formula, err := getFormula(val); err == nil {
			formulaCVal := CVal{Type: CTFormula, ValFormula: formula}
			if computed == nil {
				return formulaCVal, nil
			} else {
				computedCVal, err := makeCellVal(computed, nil)
				if err != nil || computedCVal.Type == CTFormula {
					return formulaCVal, nil
				}
				formulaCVal.HasComputed = true
				formulaCVal.ComputedType = computedCVal.Type
				formulaCVal.ValString = computedCVal.ValString
				formulaCVal.ValNumber = computedCVal.ValNumber
				formulaCVal.ValBool = computedCVal.ValBool

				return formulaCVal, nil
			}
		}
		return CVal{Type: CTString, ValString: val}, nil
	case bool:
		return CVal{Type: CTBool, ValBool: val}, nil
	default:
		return CVal{}, errors.New("unknown cell type")
	}
}

func (cell CVal) String() string {
	switch cell.Type {
	case CTString:
		return cell.ValString
	case CTFormula:
		if !cell.HasComputed {
			return string(cell.ValFormula)
		}
		return fmt.Sprintf("%s -> %s", cell.ValFormula, cell.Computed().String())
	case CTNumber:
		if cell.ValNumber == float32(math.Trunc(float64(cell.ValNumber))) {
			return strconv.Itoa(int(cell.ValNumber))
		} else {
			return fmt.Sprintf("%.3f", cell.ValNumber)
		}
	case CTBool:
		if cell.ValBool {
			return "true"
		} else {
			return "false"
		}
	case CTEmpty:
		return "nil"
	}
	return "unknown"
}

func (cell CVal) ToValue() any {
	switch cell.Type {
	case CTString:
		return cell.ValString
	case CTFormula:
		return cell.ValFormula
	case CTNumber:
		return cell.ValNumber
	case CTBool:
		return cell.ValBool
	case CTEmpty:
		return nil
	}
	return nil
}
