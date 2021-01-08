package amatriciana

import (
	"testing"
)

func TestCoordString(t *testing.T) {
	one := xy{2, 2}.String()
	if one != "b2" {
		t.Fail()
	}

	two := xy{7, 3}.String()
	if two != "g3" {
		t.Fail()
	}
}

func TestParsexy(t *testing.T) {
	value, err := parsexy("a2")
	expectedResult := xy{1, 2}
	if value != expectedResult || err != nil {
		t.Fail()
	}

	value, err = parsexy("h3")
	expectedResult = xy{8, 3}
	if value != expectedResult || err != nil {
		t.Fail()
	}

	value, err = parsexy("i3")
	if err == nil {
		t.Fail()
	}
}

func TestParseMove(t *testing.T) {
	board1 := NewBoard()

	move, err := board1.parseMove("e2e4")
	if err != nil {
		t.Fail()
	}
	expectedFrom := xy{5, 2}
	expectedTo := xy{5, 4}
	if move.from != expectedFrom || move.to != expectedTo || move.piece != pawn {
		t.Fail()
	}

	board2, err := BoardFromFEN("5kn1/p4p2/8/2qp4/5P2/1PBK2P1/P3r3/8 b - - 1 26")
	if err != nil {
		t.Fail()
	}
	move, err = board2.parseMove("c5e3")
	if err != nil {
		t.Fail()
	}
	expectedFrom = xy{3, 5}
	expectedTo = xy{5, 3}
	if move.from != expectedFrom || move.to != expectedTo || move.piece != queen {
		t.Fail()
	}
}
