package xl

import (
	"errors"
	"fmt"
	"strconv"
)

type RawSheet struct {
	Name    string  `json:"name" binding:"required"`
	Content [][]any `json:"content" binding:"required"`
	// The computed values of formulas in the sheet, if known.
	// If not provided, the computed values will be of size 0.
	Computed [][]any `json:"computed"`
}

type Sheet struct {
	Name    string   `json:"name"`
	Content [][]CVal `json:"content"`
}

type UsedRange struct {
	RowCount int
	ColCount int
}

func (s *Sheet) Get(c Cell) (CVal, error) {
	if !c.IsInBounds(s) {
		return CVal{}, errors.New("cell is out of bounds")
	}
	return s.Content[c.Row][c.Col], nil
}

func (s *Sheet) GetRange(c Range) ([][]CVal, error) {
	vals := make([][]CVal, 0)
	for i := c.Start.Row; i < c.End.Row; i++ {
		row := make([]CVal, 0)
		for j := c.Start.Col; j < c.End.Col; j++ {
			val, err := s.Get(Cell{Row: i, Col: j, Sheet: s.Name})
			if err != nil {
				return [][]CVal{}, err
			}
			row = append(row, val)
		}
		vals = append(vals, row)
	}
	return vals, nil
}

func (s *Sheet) UsedRange() UsedRange {
	if len(s.Content) == 0 {
		return UsedRange{}
	}
	return UsedRange{
		RowCount: len(s.Content),
		ColCount: len(s.Content[0]),
	}
}

func (s *Sheet) ContentValues() [][]any {
	vals := make([][]any, len(s.Content))
	for i, row := range s.Content {
		vals[i] = make([]any, len(row))
		for j, cell := range row {
			vals[i][j] = cell.ToValue()
		}
	}
	return vals
}

func isRectangular(v [][]any) bool {
	if len(v) == 0 {
		return true
	}
	length := len(v[0])
	for _, row := range v {
		if len(row) != length {
			return false
		}
	}
	return true
}

func (s *RawSheet) ToSheet() (Sheet, error) {
	content, err := makeContent(s.Content, s.Computed)
	if err != nil {
		return Sheet{}, err
	}
	return Sheet{
		Name:    s.Name,
		Content: content,
	}, nil
}

func (s Sheet) PrettyPrint(length ...int) {
	var fmtLen int
	if len(length) == 0 {
		fmtLen = 11
	} else {
		fmtLen = length[0]
	}
	print := func(v string) {
		if len(v) >= fmtLen {
			v = v[0:fmtLen]
		}
		fmt.Printf("| %"+strconv.Itoa(fmtLen)+"s ", v)
	}
	for _, row := range s.Content {
		for _, cell := range row {
			print(cell.String())
		}
		fmt.Println("|")
	}
}

// Cells is an iterator over the cells in the sheet, giving their value and position.
// func (s Sheet) Cells() func(yield func(Cell, CVal) bool) {
// 	return func(yield func(Cell, CVal) bool) {
// 		for i, row := range s.Content {
// 			for j, val := range row {
// 				if !yield(Cell{Row: uint32(i), Col: uint16(j), Sheet: s.Name}, val) {
// 					return
// 				}
// 			}
// 		}
// 	}
// }

func (s Sheet) Cells() []Cell {
	cells := make([]Cell, 0)
	for i, row := range s.Content {
		for j := range row {
			cells = append(cells, Cell{Row: uint32(i), Col: uint16(j), Sheet: s.Name})
		}
	}
	return cells
}
