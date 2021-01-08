package amatriciana

import (
	"fmt"
	"testing"
)

func TestAttackers(t *testing.T) {
	board, err := BoardFromFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1")
	if err != nil {
		t.Fail()
	}

	coord, err := parsexy("d5")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if board.HowManyAttack(coord, white) != 1 {
		fmt.Println("attackers of d5:", board.HowManyAttack(xy{4, 5}, white))
		t.Fail()
	}

	board, err = BoardFromFEN("r1bqkbnr/3p1pp1/n1pBp2p/pp6/3P1P2/3BPN2/PPP1Q1PP/RN2K2R w KQkq - 0 11")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	coord, err = parsexy("e5")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	attackers := board.HowManyAttack(coord, white)
	if attackers != 4 {
		fmt.Println("attackers of e5:", attackers)
		t.Fail()
	}
}

func TestCheckmate(t *testing.T) {
	board, err := BoardFromFEN("R1k5/6R1/8/8/8/3K4/8/8 b - - 11 6")
	if err != nil {
		t.Fail()
	}

	if !board.IsCheckmate() {
		t.Fail()
	}
}

func TestCheckmateEvaluation(t *testing.T) {
	board, err := BoardFromFEN("R1k5/6R1/8/8/8/3K4/8/8 b - - 11 6")
	if err != nil {
		t.Fail()
	}

	if board.Evaluate() != 1000 {
		t.Fail()
	}
	fmt.Println("board evaluation:", board.Evaluate())
}

func TestMaxMin(t *testing.T) {
	if max(3.0, 5.0) != 5.0 {
		t.Fail()
	}
	if min(3.0, 5.0) != 3.0 {
		t.Fail()
	}
}
