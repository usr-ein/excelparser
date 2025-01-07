package xl

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Address string
type LocalAddress string

type SplitAddress struct {
	// From the address Sheet1!A1

	Sheet        string       // e.g. Sheet1
	LocalAddress LocalAddress // e.g. A1
}

var /* const */ isAddressRegex = regexp.MustCompile(`[$]?[A-Z]+[$]?[1-9][0-9]*`)

func IsAddress(s string) bool {
	/**
	 * Checks if a string is a valid address.
	 * E.g. A1, $A$1, $A1, A$1, AA1, A11, AZ1, ...
	 */
	return isAddressRegex.MatchString(s)
}

func ParseAddress(address string) (Address, error) {
	/**
	 * Parses an address.
	 * Very fast but not very strict.
	 * Will not check:
	 * - if the sheet name is valid
	 * - if the column is valid (not out of bound)
	 * - if the row is valid (not out of bound)
	 * For full validation, try ToCell or ToCellWithSheet.
	 */
	if !IsAddress(address) {
		return "", errors.New("incompatible format")
	}

	return Address(address), nil
}

func (c Cell) ToAddress() Address {
	/**
	 * Converts a cell to an address.
	 * E.g. A1, $A$1, $A1, A$1, AA1, A11, AZ1, ...
	 */
	col, err := columnToLetters(c.Col)
	if err != nil {
		return ""
	}

	row := strconv.Itoa(int(c.Row + 1))

	if !c.RowRel {
		row = "$" + row
	}
	if !c.ColRel {
		col = "$" + col
	}
	sheetName := c.Sheet
	if shouldQuoteSheetName(sheetName) {
		sheetName = "'" + sheetName + "'"
	}
	address := sheetName + "!" + col + row

	// Leap of faith
	return Address(address)
}

// https://stackoverflow.com/a/53483274/5989906
var goodSheetName = regexp.MustCompile(`^[^\d.][a-zA-Z0-9_]*$`)

func shouldQuoteSheetName(sheetName string) bool {
	return !goodSheetName.MatchString(sheetName)
}

func splitAddress(address Address) (SplitAddress, error) {
	split := strings.Split(string(address), "!")

	if len(split) > 1 {
		sheet := split[0]
		sheet = strings.TrimPrefix(sheet, "'")
		sheet = strings.TrimSuffix(sheet, "'")

		if sheet == "" {
			return SplitAddress{}, errors.New("missing sheet name")
		}

		return SplitAddress{
			Sheet:        sheet,
			LocalAddress: LocalAddress(split[1]),
		}, nil
	}

	return SplitAddress{
		Sheet:        "",
		LocalAddress: LocalAddress(address),
	}, nil
}

func columnToLetters(col uint16) (string, error) {
	/**
	 * Converts a column number to the corresponding letters.
	 * E.g. 0 -> A, 1 -> B, 25 -> Z, 26 -> AA, 27 -> AB
	 */
	if col > MAX_COLS-1 {
		return "", errors.ErrUnsupported
	}
	var dividend int = int(col + 1)

	columnName := ""

	for dividend > 0 {
		modulo := (dividend - 1) % 26
		// ignore the linter here
		columnName = string(rune(65+modulo)) + columnName
		dividend = (dividend - modulo) / 26
	}

	return columnName, nil
}

func lettersToColumn(letters string) (uint16, error) {
	/**
	 * Converts a string of letters to a column number.
	 * A: 0, B: 1, ..., Z: 25, AA: 26, AB: 27, AZ: 51, BA: 52, ..., ZZ: 701, AAA: 702, ...
	 */
	letters = strings.ToUpper(letters)
	column := 0
	for _, char := range letters {
		if char > unicode.MaxASCII {
			return 0, errors.ErrUnsupported
		}
		column = column*26 + int(char-'A') + 1
	}
	column -= 1

	if column >= MAX_COLS || column > math.MaxUint16 || column < 0 {
		return 0, errors.ErrUnsupported
	}
	return uint16(column), nil
}

func (c Cell) ToAddressRel(relativeSheetName string) Address {
	// Turns the cell into an address string, but with the sheet name
	// removed if it matches the relativeSheetName.
	if c.Sheet == relativeSheetName {
		return c.ToAddressNoSheet()
	}
	return c.ToAddress()
}

func (c Cell) ToAddressNoSheet() Address {
	addr := string(c.ToAddress())
	addrSplit := strings.Split(addr, "!")
	if len(addrSplit) == 2 {
		return Address(addrSplit[1])
	}
	return Address(addrSplit[0])
}
