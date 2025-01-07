package xl

import (
	"encoding/json"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestToSheet_Good(t *testing.T) {
	// 2**63 ±= 9.22e+18 but we convert to float32 and back so it moves a bit
	sheetStr := `{
        "name": "Sheet1",
        "content": [
            [0, 1, 1000000000000000000, 3, 4, 5, 6, 7, 8, "=SUM(A1:I1)"],
            [11, 12, null, 14, 15, 16, 17, 18, 19, "="],
            [
                "=SUM(A1:A2)",
                "=SUM(B1:B2)",
                "=SUM(C1:C2)",
                "=SUM(D1:D2)",
                "=SUM(E1:E2)",
                "=SUM(F1:F2)",
                "=SUM(G1:G2)",
                "=SUM(H1:H2)",
                "=SUM(I1:I2)",
                "=SUM(J1:J2)"
            ]
        ]
    }`
	var rawSheet RawSheet
	err := json.Unmarshal([]byte(sheetStr), &rawSheet)
	if err != nil {
		t.Errorf("json.Unmarshal failed with %s", err)
	}
	parsedSheet, err := rawSheet.ToSheet()
	if err != nil {
		t.Errorf("rawSheet.ToSheet failed with %s", err)
	}
	cupaloy.SnapshotT(t, parsedSheet)
}

func TestToSheet_NotRect(t *testing.T) {
	nonRectangleSheetStr := `{
        "name": "Sheet1",
        "content": [
            [0, 1, 2, 3, 4, 5, 6, 7, 8, "=SUM(A1:I1)"],
            [11, 12, 13, 14, 15, 16, 17, 18, "=SUM(A2:I2)"],
            [
                "=SUM(A1:A2)",
                "=SUM(B1:B2)",
                "=SUM(C1:C2)",
                "=SUM(D1:D2)",
                "=SUM(E1:E2)",
                "=SUM(F1:F2)",
                "=SUM(G1:G2)",
                "=SUM(H1:H2)",
                "=SUM(I1:I2)",
                "=SUM(J1:J2)"
            ]
        ]
    }`
	var rawSheet RawSheet
	err := json.Unmarshal([]byte(nonRectangleSheetStr), &rawSheet)
	if err != nil {
		t.Errorf("json.Unmarshal failed with %s", err)
	}
	_, err = rawSheet.ToSheet()
	if err == nil || err.Error() != "content is not rectangular" {
		t.Error("rawSheet.ToSheet should have failed")
	}
}

func TestToSheet_Computed(t *testing.T) {
	sheetStr := `{
        "name": "Sheet1",
        "content": [
            [0, 1, 1000000000000000000, 3, 4, 5, 6, 7, 8, "=SUM(A1:I1)"],
            [11, 12, null, 14, 15, 16, 17, 18, 19, "="],
            [
                "=SUM(A1:A2)",
                "=SUM(B1:B2)",
                "=SUM(C1:C2)",
                "=SUM(D1:D2)",
                "=SUM(E1:E2)",
                "=SUM(F1:F2)",
                "=SUM(G1:G2)",
                "=SUM(H1:H2)",
                "=SUM(I1:I2)",
                "=SUM(J1:J2)"
            ]
        ],
        "computed": [
            [0, 1, 1000000000000000000, 3, 4, 5, 6, 7, 8, "42"],
            [11, 12, null, 14, 15, 16, 17, 18, 19, "="],
            [
                "42",
                "42",
                "42",
                "42",
                "42",
                "42",
                "42",
                "42",
                "42",
                "42"
            ]
        ]
    }`
	var rawSheet RawSheet
	err := json.Unmarshal([]byte(sheetStr), &rawSheet)
	if err != nil {
		t.Errorf("json.Unmarshal failed with %s", err)
	}
	parsedSheet, err := rawSheet.ToSheet()
	if err != nil {
		t.Errorf("rawSheet.ToSheet failed with %s", err)
	}
	cupaloy.SnapshotT(t, parsedSheet)
}
