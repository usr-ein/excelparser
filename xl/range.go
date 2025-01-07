package xl

import (
	"errors"
	"fmt"
	"strings"
)

// Make sure that start and end are in the same sheet!
type Range struct {
	Start Cell
	End   Cell
}

func (r Range) String() string {
	// Turns the range into a string.
	// Always includes the sheet name at the beginning address.
	end := string(r.End.ToAddress())
	endNoSheet := strings.TrimPrefix(end, r.End.Sheet)
	return string(r.Start.ToAddress()) + ":" + endNoSheet
}

func (r Range) StringRel(relativeSheetName string) string {
	// Turns the range into a string, but with the sheet name
	// removed from the start address if it matches the relativeSheetName.
	unshiftedEnd, err := r.End.Shift(-1, -1)
	if err != nil {
		return "ERROR_STRINGIFYING_RANGE"
	}
	endNoSheet := unshiftedEnd.ToAddressNoSheet()
	startRel := r.Start.ToAddressRel(relativeSheetName)
	return fmt.Sprintf("%s:%s", startRel, endNoSheet)
}

func ParseRange(s string, sheetName string) (Range, error) {
	startEnd := strings.Split(s, ":")
	if len(startEnd) != 2 {
		return Range{}, errors.New("incompatible format")
	}
	startAddr, errStart := ParseAddress(startEnd[0])
	endAddr, errEnd := ParseAddress(startEnd[1])
	if errors.Join(errStart, errEnd) != nil {
		return Range{}, errors.New("incompatible format")
	}

	startCell, err := ToCell(startAddr, sheetName)
	if err != nil {
		return Range{}, err
	}

	endCell, err := ToCell(endAddr, startCell.Sheet)
	if err != nil {
		return Range{}, err
	}
	endCell.Sheet = startCell.Sheet

	endCell, err = endCell.Shift(1, 1)
	if err != nil {
		return Range{}, err
	}

	return Range{startCell, endCell}, nil
}

// Returns a list of all cells in the range.
func (r Range) Cells() (cells []Cell) {
	for i := r.Start.Row; i < r.End.Row; i++ {
		for j := r.Start.Col; j < r.End.Col; j++ {
			cells = append(cells, Cell{
				Sheet:  r.Start.Sheet,
				Row:    i,
				Col:    j,
				RowRel: true,
				ColRel: true,
			})
		}
	}
	return
}
