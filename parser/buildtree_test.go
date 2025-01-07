package parser

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestBuildtreeNotCrash(t *testing.T) {
	formulas := []string{
		`A1+5`,
		`(A1+5)*B1`,
		`(SUM(A1:B1, A2:B2)+A5+3232-A15)/B90/321.0+MONTH("January")`,
		`SUM(A1,A2:B3, SUM(A1:B1, SUM(A2:B2)), C3)`,
		`SUM(A1:A4)+5`,
	}
	for _, f := range formulas {
		tokens := Tokenize(f)
		_, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
		if err != nil {
			t.Errorf("could not build tree for %s: %v", f, err)
		}
	}
}

func TestBuildtreeUnmatchedParen(t *testing.T) {
	formulas := []string{
		`A(1+5`,
		`(A1+5`,
		`(A1+5)*B1(`,
		`SUM(A1:(A4)+5`,
	}
	for _, f := range formulas {
		tokens := Tokenize(f)
		_, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
		if err == nil {
			t.Errorf("expected error for %s", f)
		}
	}
}

func TestBuildtree_F1(t *testing.T) {
	f := `A1+5`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}
	nodeJson := ToNodeJson(tree)
	cupaloy.SnapshotT(t, nodeJson)
}

func TestBuildtree_F2(t *testing.T) {
	f := `(A1+5)*B1`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}
	nodeJson := ToNodeJson(tree)
	cupaloy.SnapshotT(t, nodeJson)
}

func TestBuildtree_F3(t *testing.T) {
	f := `(SUM(A1:B1, A2:B2)+A5+3232-A15)/B90/321.0+MONTH("January")`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}
	nodeJson := ToNodeJson(tree)
	cupaloy.SnapshotT(t, nodeJson)
}

func TestBuildtree_F4(t *testing.T) {
	f := `SUM(A1,A2:B3, SUM(A1:B1, SUM(A2:B2)), C3)`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}
	nodeJson := ToNodeJson(tree)
	cupaloy.SnapshotT(t, nodeJson)
}

func TestBuildtree_F5(t *testing.T) {
	f := `SUM(A1:A4)+5`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}
	nodeJson := ToNodeJson(tree)
	cupaloy.SnapshotT(t, nodeJson)
}

func TestBuildtree_RangeStartEnd(t *testing.T) {
	f := `SUM(Sheet2!A1:A5)`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}

	cupaloy.SnapshotT(t, tree)
}

func TestBuildtree_RangeStartEndBad(t *testing.T) {
	// This is invalid notation, but we allow it and just ignore the sheet name of the end cell,
	// and replace it with the sheet name of the start cell.
	f := `SUM(Sheet2!A1:Sheet3!A5)`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}

	cupaloy.SnapshotT(t, tree)
}

func TestBuildtree_RangeStartImplicit(t *testing.T) {
	// This is invalid notation, but we allow it and just ignore the sheet name of the end cell,
	// and replace it with the sheet name of the start cell.
	f := `SUM(A1:Sheet2!A5)`
	tokens := Tokenize(f)
	tree, err := BuildTree(Context{CurrentSheet: "Sheet1"}, tokens)
	if err != nil {
		t.Errorf("could not build tree for %s: %v", f, err)
	}

	cupaloy.SnapshotT(t, tree)
}
