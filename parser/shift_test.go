package parser

import (
	"testing"
)

func TestShiftFormulaStaysSame(t *testing.T) {
	// Very important test: if we shift a formula by 0, it should stay the same,
	// character for character.
	f := Formula(`=('Operations (no battery)'!CI$25+'Operations (no battery)'!CI$26)*INDEX(General!$C$52:$I$56, MATCH(CI$5, General!$C$52:$C$56, 0), MATCH("Yes", General!$C$51:$I$51, 0))/1000`)
	shifted, err := ShiftFormula(f, 0, 0, `Cashflow (no battery)`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	if shifted != f {
		t.Errorf("ShiftFormula failed, expected %s, got %s", f, shifted)
		return
	}
}

func TestShiftFormulaBackForth(t *testing.T) {
	f := Formula(`=('Operations (no battery)'!CI$25+'Operations (no battery)'!CI$26)*INDEX(General!$C$52:$I$56, MATCH(CI$5, General!$C$52:$C$56, 0), MATCH("Yes", General!$C$51:$I$51, 0))/1000`)
	shifted, err := ShiftFormula(f, 2, 3, `Cashflow (no battery)`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}
	doubleShifted, err := ShiftFormula(shifted, -2, -3, `Cashflow (no battery)`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}
	if shifted == f {
		t.Errorf("ShiftFormula failed, shifted formula is the same as the original")
		return
	}

	if doubleShifted != f {
		t.Errorf("ShiftFormula failed, expected %s, got %s", f, doubleShifted)
		return
	}
}

func TestShiftFormulaSingleRef(t *testing.T) {
	f := Formula(`=A1`)
	shifted, err := ShiftFormula(f, 2, 3, `Cashflow (no battery)`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}
	doubleShifted, err := ShiftFormula(shifted, -2, -3, `Cashflow (no battery)`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}
	if shifted == doubleShifted {
		t.Errorf("ShiftFormula failed, shifted formula is the same as the original")
		return
	}

	if doubleShifted != f {
		t.Errorf("ShiftFormula failed, expected %s, got %s", f, doubleShifted)
		return
	}
}

func TestShiftFormulaSimple(t *testing.T) {
	f := Formula(`=SUM(A1:B1)`)
	shifted, err := ShiftFormula(f, 5, 3, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	expected := Formula(`=SUM(D6:E6)`)

	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", expected, shifted)
		return
	}
}

func TestShiftFormulaOnlyCol(t *testing.T) {
	f := Formula(`=SUM(A1:B1)`)
	shifted, err := ShiftFormula(f, 0, 3, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	expected := Formula(`=SUM(D1:E1)`)

	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", expected, shifted)
		return
	}
}

func TestShiftFormulaOnlyRow(t *testing.T) {
	f := Formula(`=SUM(A1:B1)`)
	shifted, err := ShiftFormula(f, 3, 0, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	expected := Formula(`=SUM(A4:B4)`)

	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", expected, shifted)
		return
	}
}

func TestShiftFormulaNoOp(t *testing.T) {
	f := Formula(`=SUM(A1:B1)`)
	shifted, err := ShiftFormula(f, 0, 0, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	expected := Formula(`=SUM(A1:B1)`)

	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", expected, shifted)
		return
	}
}

func TestShiftFormulaPrecendenceExtraParenth(t *testing.T) {
	f := Formula(`=(A1+B2)-C3*D4/E5`)
	shifted, err := ShiftFormula(f, 0, 0, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)

		return
	}

	expected := Formula(`=A1+B2-C3*D4/E5`)
	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", f, shifted)
		return
	}
}

func TestShiftFormulaComplexShift(t *testing.T) {
	f := Formula(`=('Operations (no battery)'!CI$25+'Operations (no battery)'!CI$26)*INDEX(General!$C$52:$I$56, MATCH(CI$5, General!$C$52:$C$56, 0), MATCH("Yes", General!$C$51:$I51, 0))/1000`)
	shifted, err := ShiftFormula(f, 3, 5, `Cashflow (no battery)`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	expected := Formula(`=('Operations (no battery)'!CN$25+'Operations (no battery)'!CN$26)*INDEX(General!$C$52:$I$56, MATCH(CN$5, General!$C$52:$C$56, 0), MATCH("Yes", General!$C$51:$I54, 0))/1000`)

	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", expected, shifted)
		return
	}
}

func TestShiftFormulaUnary(t *testing.T) {
	f := Formula(`=A1*-B2`)
	shifted, err := ShiftFormula(f, 0, 0, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	if shifted != f {
		t.Errorf("ShiftFormula failed, expected %s, got %s", f, shifted)
		return
	}
}

func TestShiftFormulaUnaryParenth(t *testing.T) {
	f := Formula(`=A1*(-B2)`)
	shifted, err := ShiftFormula(f, 0, 0, `Sheet1`)
	if err != nil {
		t.Errorf("ShiftFormula failed with %s", err)
		return
	}

	expected := Formula(`=A1*-B2`)
	if shifted != expected {
		t.Errorf("ShiftFormula failed, expected %s, got %s", f, shifted)
		return
	}
}
