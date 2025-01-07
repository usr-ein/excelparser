package xl

import (
	"errors"
	"regexp"
	"strings"
)

type Shiftable interface {
	Shift(rowDiff int, colDiff int) (Shiftable, error)
}

func (r Range) Shift(rowDiff int, colDiff int) (Range, error) {
	start, err := r.Start.Shift(rowDiff, colDiff)
	if err != nil {
		return Range{}, err
	}
	end, err := r.End.Shift(rowDiff, colDiff)
	if err != nil {
		return Range{}, err
	}
	return Range{start, end}, nil
}

func (r Range) ShiftIfRel(rowDiff int, colDiff int) (Range, error) {
	start, err := r.Start.ShiftIfRel(rowDiff, colDiff)
	if err != nil {
		return Range{}, err
	}
	end, err := r.End.ShiftIfRel(rowDiff, colDiff)
	if err != nil {
		return Range{}, err
	}
	return Range{start, end}, nil
}

func (c Cell) Shift(rowDiff int, colDiff int) (Cell, error) {
	rowUnderflow := rowDiff < 0 && rowDiff < -int(c.Row)
	rowOverflow := (rowDiff > 0 && rowDiff > int(MAX_ROWS-c.Row))
	colUnderflow := (colDiff < 0 && colDiff < -int(c.Col))
	colOverflow := (colDiff > 0 && colDiff > int(MAX_COLS-c.Col))

	if rowOverflow || rowUnderflow || colOverflow || colUnderflow {
		return Cell{}, errors.New("cannot shift cell outside of sheet")
	}
	row := int(c.Row)
	col := int(c.Col)

	row += rowDiff
	col += colDiff

	return Cell{
		Sheet:  c.Sheet,
		Row:    uint32(row),
		Col:    uint16(col),
		RowRel: c.RowRel,
		ColRel: c.ColRel,
	}, nil
}

func (c Cell) ShiftIfRel(rowDiff int, colDiff int) (shifted Cell, err error) {
	shifted = c
	if c.RowRel {
		if shifted, err = shifted.Shift(rowDiff, 0); err != nil {
			return Cell{}, err
		}
	}
	if c.ColRel {
		if shifted, err = shifted.Shift(0, colDiff); err != nil {
			return Cell{}, err
		}
	}
	return shifted, nil
}

var reCommaWithSpaces = regexp.MustCompile(`\s*,\s*`)

func (f Formula) RemoveCommaSpaces() Formula {
	// Remove spaces around commas in formulas.
	// E.g. =SUM(A1, B1) -> =SUM(A1,B1)
	// May completely fuck up the formula, e.g.
	// ='Live, love, laugh'!A1+A2 becomes ='Live love laugh'!A1+A2
	// but the 'Live love laugh' sheet doesn't exist. There can be more pesky
	// edge cases like that.
	return Formula(reCommaWithSpaces.ReplaceAllString(string(f), ","))
}

func (f Formula) RemoveDollars() Formula {
	// Removes all dollar signs from the formula.
	// May completely fuck it up, but it's useful for some no-op checking.
	return Formula(strings.ReplaceAll(string(f), "$", ""))
}
