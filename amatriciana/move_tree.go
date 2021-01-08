package amatriciana

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"
)

type moveCalculation struct {
	board    Board
	eval     float32
	branches []*moveCalculation
}

var positions map[string]*moveCalculation

func (b Board) treeify(depth int, previousScore float32) (*moveCalculation, error) {
	if b.IsCheckmate() {
		return &moveCalculation{b, 0, nil}, nil
	}

	eval := b.Evaluate()
	color := b.turn
	if (eval < 0 && previousScore < 0 && color == white) ||
		(eval > 0 && previousScore > 0 && color == black) {
		return &moveCalculation{b, 0, nil}, nil
	}

	if depth == 0 {
		return &moveCalculation{b, 0, nil}, nil
	}

	premadeCalculation, alreadyCalculated := positions[b.FEN()]
	if alreadyCalculated {
		//println(b.FEN(), "already calculated")
		return premadeCalculation, nil
	}

	moves := b.moves(b.turn)
	if len(moves) == 0 {
		return &moveCalculation{b, b.Evaluate(), nil}, nil
	}

	calculation := &moveCalculation{
		board:    b,
		branches: make([]*moveCalculation, 0, len(moves)),
	}

	//return calculation, nil

	for _, move := range moves {
		childBoard := b.Clone()
		moved := childBoard.move(move)
		if !moved {
			continue
		}

		newBranch, err := childBoard.treeify(depth-1, eval)
		if err != nil {
			continue
		}

		if newBranch == nil {
			continue
		}
		calculation.branches = append(calculation.branches, newBranch)
	}

	positions[b.FEN()] = calculation

	return calculation, nil

}

func (m moveCalculation) printTree() string {
	var output bytes.Buffer

	output.WriteString("\nboard: ")
	output.WriteString(m.board.FEN())

	output.WriteString("\nself evaluation: ")
	output.WriteString(strconv.FormatFloat(float64(m.eval), 'f', 6, 64))
	output.WriteString("\n")
	return output.String()
}

func (m move) UCIString() string {
	return strings.Join([]string{m.from.String(), m.to.String()}, "")
}

func minimax(calc *moveCalculation, depth int, alpha, beta float32, maximiseFor color) float32 {

	if calc == nil {
		println("not good, this pointer seems to be nil")
		return 0
	}

	if depth == 0 {
		//println("depth 0, returning static evaluation:", FloatToString(calc.board.Evaluate()))
		return calc.board.Evaluate()
	}

	if len(calc.branches) == 0 {
		return calc.board.Evaluate()
	}

	var bestEval float32
	if maximiseFor == white {
		bestEval = -10000.0
		for i := range calc.branches {
			eval := minimax(calc.branches[i], depth-1, alpha, beta, !maximiseFor)
			bestEval = max(bestEval, eval)
			alpha = max(alpha, eval)
			if alpha <= beta {
				break
			}
		}
	} else if maximiseFor == black {
		bestEval = 10000.0
		for i := range calc.branches {
			eval := minimax(calc.branches[i], depth-1, alpha, beta, !maximiseFor)
			bestEval = min(bestEval, eval)
			beta = min(beta, eval)
			if beta <= alpha {
				break
			}
		}
	}

	calc.eval = bestEval

	return bestEval
}

func (b Board) BestMove(depth int) (move, error) {
	println("initializing positions map")
	positions = make(map[string]*moveCalculation, 2000)

	maximiseFor := b.turn
	moves := b.moves(b.turn)
	println("there are", len(moves), "possible moves")

	bestMove := move{}
	var biggestAdvantage float32 = -1000.0

	//tree, err := b.treeify(depth + 1)
	//if err != nil {
	//	return move{}, err
	//}
	for _, move := range moves {
		println("generating tree for", move.String())
		branchBoard := b.Clone()
		branchBoard.move(move)
		//println("branching move", move.String())
		start := time.Now()
		tree, err := branchBoard.treeify(depth, 0.0)
		stop := time.Since(start)
		println("tree took", Float64ToString(stop.Seconds()), "seconds")
		if err != nil {
			continue
		}
		if tree == nil {
			println("tree for", move.String(), "is still nil! :(")
		}
		println("minimaxing tree for", move.String())
		//println("minimaxing move", move.String())
		start = time.Now()
		advantageGained := minimax(tree, depth, -10000, 10000, maximiseFor)
		stop = time.Since(start)
		println("minimax took", Float64ToString(stop.Seconds()), "seconds")
		println("eval for", move.String(), "is", FloatToString(advantageGained))
		if maximiseFor == black {
			advantageGained *= -1
		}

		if advantageGained > biggestAdvantage {
			bestMove = move
			biggestAdvantage = advantageGained
		}
	}

	emptyMove := move{}
	if bestMove == emptyMove {
		return move{}, errors.New("couldn't find the best move no idea why don't @ me")
	}

	println("i calculated", len(positions), "different positions")
	return bestMove, nil
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}

	return b
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}

	return b
}

func FloatToString(input_num float32) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(input_num), 'f', 2, 64)
}

func Float64ToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}
