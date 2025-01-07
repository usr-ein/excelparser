package xl

import (
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestParseRange_Bad(t *testing.T) {
	badAddresses := []string{
		"A1:",
		":B1",
		// "A1:B1",
		"B:C", // We don't support column ranges yet/ever
	}
	for _, strAddress := range badAddresses {
		cellRange, err := ParseRange(strAddress, "Sheet1")
		if err != nil {
			continue
		}
		t.Errorf("ParseRange(%s, \"Sheet1\") returned \"%s\" and succeeded; want error", strAddress, cellRange.String())
	}
}

type rangeTest struct {
	Start Cell
	End   Cell
}

func TestParseRange_Good(t *testing.T) {
	goodAddresses := []string{
		"A1:B1",
		"B1:A1",
		"Sheet1!A1:B1",
		"Sheet2!A1:B1",
		"Sheet2!A1:Sheet2!B1", // weird but allowed and implicitly fixed to Sheet2!A1:B1
		"Sheet2!A1:Sheet1!B1", // weird but allowed and implicitly fixed to Sheet2!A1:B1
	}
	for _, strAddress := range goodAddresses {
		testName := strings.ReplaceAll(strAddress, ":", "_")
		testName = strings.ReplaceAll(testName, "!", "_")
		t.Run("TestParseRange_Good_"+testName, func(tt *testing.T) {
			cellRange, err := ParseRange(strAddress, "Sheet1")
			if err != nil {
				tt.Errorf("ParseRange(%s, \"Sheet1\") failed with %s error", strAddress, err.Error())
			}

			cupaloy.SnapshotT(tt, rangeTest(cellRange))
		})
	}
}
