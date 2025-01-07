package xl

import (
	"errors"
	"regexp"
	"strconv"
)

const MAX_ROWS = 1_048_576
const MAX_COLS = 16_384

type Cell struct {
	Sheet string `json:"sheet"`

	/* Max 2**20 = 1 048 576 */
	Row uint32 `json:"row"`
	/* Max 2**14 == 16 384 */
	Col uint16 `json:"col"`

	RowRel bool `json:"rowRel"`
	ColRel bool `json:"colRel"`
}

func NewCell(row int, col int, sheet string) (Cell, error) {
	if row < 0 || row >= MAX_ROWS {
		return Cell{}, errors.New("invalid row number")
	}
	if col < 0 || col >= MAX_COLS {
		return Cell{}, errors.New("invalid column number")
	}
	return Cell{
		Sheet:  sheet,
		Row:    uint32(row),
		Col:    uint16(col),
		RowRel: false,
		ColRel: false,
	}, nil
}

func (c Cell) StripDollars() Cell {
	return Cell{
		Sheet:  c.Sheet,
		Row:    c.Row,
		Col:    c.Col,
		RowRel: true,
		ColRel: true,
	}
}

func (c Cell) WithDollars() Cell {
	return Cell{
		Sheet:  c.Sheet,
		Row:    c.Row,
		Col:    c.Col,
		RowRel: false,
		ColRel: false,
	}
}

func (c Cell) IsEq(other Cell) bool {
	return c.Sheet == other.Sheet && c.Row == other.Row && c.Col == other.Col && c.RowRel == other.RowRel && c.ColRel == other.ColRel
}

func (c Cell) IsInBounds(s *Sheet) bool {
	return int(c.Row) < len(s.Content) && int(c.Col) < len(s.Content[int(c.Row)])
}

func ParseCell(s string, defaultSheet string) (Cell, error) {
	addr, err := ParseAddress(s)
	if err != nil {
		return Cell{}, err
	}
	return ToCell(addr, defaultSheet)
}

func ToCell(address Address, defaultSheet string) (Cell, error) {
	/**
	 * Converts an address to a cell.
	 * If the address doesn't contain the sheet name, the default sheet name is used.
	 * If the address contains the sheet name, the default sheet name is ignored,
	 * and the sheet name from the address is used.
	 * If neither have a sheet name, we fail with "missing sheet prefix".
	 *
	 * Test payloads:
	 * addresses := []string{"D1", "$C$1", "She$et1!$B1", "Sheet1!$B1", "Z$1", "'Sheet1'!AC$1", "A11", "AZ1"}
	 */
	split, err := splitAddress(address)
	if err != nil {
		return Cell{}, err
	}

	if split.Sheet == "" {
		if defaultSheet == "" {
			return Cell{}, errors.New("missing sheet prefix and fallback sheet name")
		}
		return toCellFromSplit(SplitAddress{
			Sheet:        defaultSheet,
			LocalAddress: split.LocalAddress,
		})
	}

	return toCellFromSplit(split)
}

var /* const */ addressDigitsRegex = regexp.MustCompile(`\d+`)
var /* const */ addressLettersRegex = regexp.MustCompile(`[A-Z]+`)

func toCellFromSplit(address SplitAddress) (Cell, error) {
	lettersMatching := addressLettersRegex.FindAllString(string(address.LocalAddress), -1)
	if len(lettersMatching) != 1 {
		return Cell{}, errors.New("invalid local address: no letters")
	}
	numbersMatching := addressDigitsRegex.FindAllString(string(address.LocalAddress), -1)
	if len(numbersMatching) != 1 {
		return Cell{}, errors.New("invalid local address: no digits")
	}
	row, err := strconv.Atoi(numbersMatching[0])
	if err != nil {
		return Cell{}, err
	}
	row -= 1
	if row < 0 || row >= MAX_ROWS {
		return Cell{}, errors.New("invalid row number")
	}

	col, err := lettersToColumn(lettersMatching[0])
	if err != nil {
		return Cell{}, err
	}

	rel := getRelativeness(address.LocalAddress)
	return Cell{
		Sheet: address.Sheet,
		Row:   uint32(row),
		Col:   col,

		RowRel: rel.Row,
		ColRel: rel.Col,
	}, nil

}
