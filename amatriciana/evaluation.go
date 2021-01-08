package amatriciana

import "fmt"

type Evaluation struct {
	Advantage float32
	MateIn    int
}

func (b Board) bestMove() move {
	moves := b.moves(b.turn)

	if len(moves) == 0 {
		return move{}
	}

	bestMove := moves[0]
	var biggestAdvantage float32 = -2000.0

	for _, move := range moves {
		newBoard := b.Clone()
		rootingFor := newBoard.turn
		newBoard.move(move)

		eval := newBoard.Evaluate()
		advantage := eval
		if rootingFor == black {
			advantage *= -1
		}

		if advantage > biggestAdvantage {
			bestMove = move
			biggestAdvantage = advantage
		}
	}

	return bestMove
}

func (b Board) BestMoveString() string {
	bestMove := b.bestMove()

	return bestMove.String()
}

//possible things to help evaluation:
//protector overloading
//doubled pawns
//king exposed (only until middlegame)
//piece mobility (unless protecting something)
//pawn chains
//passed pawns

func (b Board) Evaluate() float32 {
	if b.IsCheckmate() {
		if b.turn == white {
			return -1000
		} else {
			return 1000
		}
	}

	whiteMaterial := b.Material(white)
	blackMaterial := b.Material(black)

	material := whiteMaterial - blackMaterial

	centerControlWhite := float32(b.HowManyAttack(xy{4, 4}, white)+
		b.HowManyAttack(xy{4, 5}, white)+
		b.HowManyAttack(xy{5, 4}, white)+
		b.HowManyAttack(xy{5, 5}, white)) / 4

	centerControlBlack := float32(b.HowManyAttack(xy{4, 4}, black)+
		b.HowManyAttack(xy{4, 5}, black)+
		b.HowManyAttack(xy{5, 4}, black)+
		b.HowManyAttack(xy{5, 5}, black)) / 4

	centerControl := centerControlWhite - centerControlBlack

	doubledPawnsWhite := float32(b.FilesWithDoubledPawns(white)) * 0.3
	doubledPawnsBlack := float32(b.FilesWithDoubledPawns(black)) * 0.3

	doubledPawns := doubledPawnsWhite - doubledPawnsBlack

	return material + float32(centerControl)*2 + float32(doubledPawns)*0.3
}

func (b Board) EvaluateVerbose() float32 {
	if b.IsCheckmate() {
		if b.turn == white {
			return -1000
		} else {
			return 1000
		}
	}

	whiteMaterial := b.Material(white)
	blackMaterial := b.Material(black)

	material := whiteMaterial - blackMaterial
	fmt.Println("material:", material)

	centerControlWhite := (b.HowManyAttack(xy{4, 4}, white) +
		b.HowManyAttack(xy{4, 5}, white) +
		b.HowManyAttack(xy{5, 4}, white) +
		b.HowManyAttack(xy{5, 5}, white)) / 4

	centerControlBlack := (b.HowManyAttack(xy{4, 4}, black) +
		b.HowManyAttack(xy{4, 5}, black) +
		b.HowManyAttack(xy{5, 4}, black) +
		b.HowManyAttack(xy{5, 5}, black)) / 4

	centerControl := centerControlWhite - centerControlBlack
	fmt.Println("centerControl:", centerControl)

	doubledPawnsWhite := float32(b.FilesWithDoubledPawns(white)) * 0.3
	doubledPawnsBlack := float32(b.FilesWithDoubledPawns(black)) * 0.3

	doubledPawns := doubledPawnsWhite - doubledPawnsBlack
	fmt.Println("doubledPawns:", doubledPawns)

	return material + float32(centerControl)*2 + float32(doubledPawns)*0.3
}

func (b Board) Material(col color) float32 {
	pieces := b.piecesOfColor(col)
	var output float32 = 0.0

	howManyBishops := 0

	for _, piece := range pieces {
		switch piece.pieceType {
		case pawn:
			output++
			chain := float32(b.PawnChain(piece.position, piece.color))
			if chain > 1 {
				output += chain * 0.2
			}
		case knight:
			output += 3.0
		case bishop:
			output += 3.0
			howManyBishops++
		case rook:
			output += 5.0

			//check if it's on an open or semiopen file
			if b.IsFileOpen(piece.position.x) {
				output += 0.5
			} else if b.IsFileSemiOpen(piece.position.x, piece.color) {
				output += 0.2
			}
		case queen:
			output += 9.0
		}
	}

	if howManyBishops == 2.0 {
		output += 1.0
	}

	return output
}

func (b Board) OpenFiles() []int {
	matrix := b.EightByEight()
	openFiles := make([]int, 0)

	for i, file := range matrix {
		for _, piece := range file {
			if piece == 'p' || piece == 'P' {
				continue
			}
			openFiles = append(openFiles, i+1)
		}
	}

	return openFiles
}

func (b Board) IsFileOpen(file int) bool {
	if file < 1 || file > 8 {
		return false
	}

	matrix := b.EightByEight()

	for _, piece := range matrix[file-1] {
		if piece == 'p' || piece == 'P' {
			return false
		}
	}

	return true
}

func (b Board) SemiOpenFiles(col color) []int {
	matrix := b.EightByEight()
	openFiles := make([]int, 0)

	var pawnLetter byte
	if col == white {
		pawnLetter = 'P'
	} else {
		pawnLetter = 'p'
	}

	for i, file := range matrix {
		for _, piece := range file {
			if piece == pawnLetter {
				continue
			}
			openFiles = append(openFiles, i+1)
		}
	}

	return openFiles
}

func (b Board) IsFileSemiOpen(file int, col color) bool {
	if file < 1 || file > 8 {
		return false
	}

	semiOpenFiles := b.SemiOpenFiles(col)

	for _, possibleFile := range semiOpenFiles {
		if file == possibleFile {
			return true
		}
	}

	return false
}

func (b Board) PawnsInFile(file int, col color) int {
	if file < 1 || file > 8 {
		return 0
	}

	var pawnLetter byte
	if col == white {
		pawnLetter = 'P'
	} else {
		pawnLetter = 'p'
	}

	matrix := b.EightByEight()
	pawns := 0
	for _, piece := range matrix[file-1] {
		if piece == pawnLetter {
			pawns++
		}
	}

	return pawns
}

func (b Board) FilesWithDoubledPawns(col color) int {
	output := 0

	for i := 1; i <= 8; i++ {
		if b.PawnsInFile(i, col) > 1 {
			output++
		}
	}

	return output
}

//recursive function to find the length of a pawn chain from the base
func (b Board) PawnChain(pawnPos xy, col color) int {
	var otherSide xy

	if col == white {
		otherSide = xy{0, 1}
	} else {
		otherSide = xy{0, -1}
	}

	captureRight := xy{1, 0}.plus(otherSide)
	captureLeft := xy{-1, 0}.plus(otherSide)

	chainRight := 0
	chainLeft := 0

	pieceRight, isThereSomething := b.pieceAtPosition(pawnPos.plus(captureRight))
	if isThereSomething && pieceRight.pieceType == pawn {
		chainRight = 1 + b.PawnChain(pawnPos.plus(captureRight), col)
	}

	pieceLeft, isThereSomething := b.pieceAtPosition(pawnPos.plus(captureLeft))
	if isThereSomething && pieceLeft.pieceType == pawn {
		chainLeft = 1 + b.PawnChain(pawnPos.plus(captureLeft), col)
	}

	if chainRight > chainLeft {
		return chainRight
	} else {
		return chainLeft
	}
}

func (b Board) HowManyAttack(square xy, col color) int {
	moves := b.moves(col)

	attackingPieces := 0
	for _, move := range moves {
		if move.to == square {
			attackingPieces++
		}
	}

	for _, piece := range b.pieces {
		if piece.pieceType == pawn {
			captures := b.pawnCaptures(piece.position, col)
			for _, capture := range captures {
				if capture.to == square {
					attackingPieces++
				}
			}
		}
	}

	return attackingPieces
}

func (b Board) CenterControl(col color) float32 {
	var output float32 = 0.0
	centralSquares := [...]string{"e4", "e5", "d4", "d5"}

	for _, squareCoord := range centralSquares {
		square, _ := parsexy(squareCoord)
		attackers := b.HowManyAttack(square, col)

		output += float32(attackers)

	}

	return output
}
