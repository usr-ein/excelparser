package xl

import (
	"testing"
)

func TestToCellWithSheet_Success(t *testing.T) {
	addressesAndCells := map[Address]Cell{
		"D1":            {Sheet: "default", Col: 3, Row: 0, RowRel: true, ColRel: true},
		"$C$1":          {Sheet: "default", Col: 2, Row: 0, RowRel: false, ColRel: false},
		"She$et1!$B1":   {Sheet: "She$et1", Col: 1, Row: 0, RowRel: true, ColRel: false},
		"Sheet1!$B1":    {Sheet: "Sheet1", Col: 1, Row: 0, RowRel: true, ColRel: false},
		"Z$1":           {Sheet: "default", Col: 25, Row: 0, RowRel: false, ColRel: true},
		"'Sheet1'!AC$1": {Sheet: "Sheet1", Col: 28, Row: 0, RowRel: false, ColRel: true},
		"A11":           {Sheet: "default", Col: 0, Row: 10, RowRel: true, ColRel: true},
		"AZ1":           {Sheet: "default", Col: 51, Row: 0, RowRel: true, ColRel: true},
		"XFD1":          {Sheet: "default", Col: MAX_COLS - 1, Row: 0, RowRel: true, ColRel: true},
		"A1048576":      {Sheet: "default", Col: 0, Row: MAX_ROWS - 1, RowRel: true, ColRel: true},
		"XFD1048576":    {Sheet: "default", Col: MAX_COLS - 1, Row: MAX_ROWS - 1, RowRel: true, ColRel: true},
	}
	for address, expectedCell := range addressesAndCells {
		cell, err := ToCell(address, "default")
		if err != nil {
			t.Errorf("ToCellWithSheet(%s, default) failed with %s", address, err)
		}
		if cell != expectedCell {
			t.Errorf("ToCellWithSheet(%s, default) = %v; want %v", address, cell, expectedCell)
		}
	}

}

func TestToCellWithSheet_Bad(t *testing.T) {
	badAddresses := []string{
		"!$B1",
		"Sheet1!",
		"Sheet1!$",
		"'Sheet1!$B",
		"'Sheet1'!$B",
		"'Sheet1'!B",
		"B1000000000000000",
		"XFE1",
		"A1048577",
		"XFD1048577",
		"XFE1048576",
	}
	for _, strAddress := range badAddresses {
		addr, convErr := ParseAddress(strAddress)
		if convErr != nil {
			continue
		}
		_, convErr = ToCell(addr, "default")
		if convErr == nil {
			t.Errorf("ToCellWithSheet(%s, \"default\") succeeded; want error", strAddress)
		}
	}
}
