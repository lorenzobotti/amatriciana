package amatriciana

import (
	"errors"
	"fmt"
)

/*type move interface {
	performMove(*Board) error
	isCapture() bool
}*/

type moveType int

const (
	normalMove moveType = iota
	enPassant
	shortCastle
	longCastle
	promotion
)

type move struct {
	piece pieceType
	color
	from xy
	to   xy
	moveType
	promotesInto pieceType
}

func (m move) String() string {
	if m.moveType == normalMove {
		return fmt.Sprint(m.piece.String(), " from ", m.from.String(), " to ", m.to.String())
	}

	if m.moveType == longCastle {
		return "castles queen side"
	}
	if m.moveType == shortCastle {
		return "castles king side"
	}

	return "some kind of other move i dunno"
}

func (b Board) parseMove(inputStr string) (move, error) {
	input := []byte(inputStr)

	if len(input) < 4 {
		return move{}, fmt.Errorf("a move should have at least 4 characters")
	}

	from, err := parsexy(string(input[0:2]))
	if err != nil {
		return move{}, err
	}
	to, err := parsexy(string(input[2:4]))
	if err != nil {
		return move{}, err
	}
	piece, isThereAPiece := b.pieceAtPosition(from)
	if !isThereAPiece {
		return move{}, fmt.Errorf("there is no piece at %s", from.String())
	}

	moveType := normalMove
	if piece.pieceType == king && from.x == 5 && to.x == 3 {
		moveType = longCastle
	}
	if piece.pieceType == king && from.x == 5 && to.x == 7 {
		moveType = shortCastle
	}

	outputMove := move{piece.pieceType, piece.color, from, to, moveType, pawn}

	if piece.pieceType == pawn && (to.y == 8 || to.y == 1) {
		moveType = promotion
		if len(input) < 5 {
			return move{}, fmt.Errorf("pawn promotes to an unknown piece")
		}

		promotesInto, err := pieceTypeFromFen(input[4])
		if err != nil {
			return outputMove, err
		}

		outputMove.promotesInto = promotesInto
	}

	return outputMove, nil
}

//PerformMove takes a uci-style move (es. e2e4) and performs it if it's legal
func (b *Board) PerformMove(input string) error {
	m, err := b.parseMove(input)
	if err != nil {
		return err
	}

	if !b.isLegal(m) || !b.isPossible(m) {
		return errors.New("illegal move")
	}

	succeeded := b.move(m)
	if !succeeded {
		return errors.New("couldn't perform move")
	}

	return nil
}

type xy struct {
	x, y int
}

func (a xy) equals(b xy) bool {
	return a.x == b.x && a.y == b.y
}

func (a xy) plus(b xy) xy {
	return xy{a.x + b.x, a.y + b.y}
}

func (a xy) String() string {
	var output = make([]byte, 2)
	output[0] = 'a' + byte(a.x) - 1
	output[1] = '1' + byte(a.y) - 1

	return string(output)
}

func parsexy(inputStr string) (xy, error) {
	input := []byte(inputStr)
	output := xy{}

	if len(input) < 2 {
		return xy{}, fmt.Errorf("there should be at least two characters")
	}

	file := input[0]
	if file >= 'a' && file <= 'h' {
		output.x = int(file-'a') + 1
	} else if file >= 'A' && file <= 'H' {
		output.x = int(file-'A') + 1
	} else {
		return xy{}, fmt.Errorf("file coordinate is invalid")
	}

	rank := input[1]
	if rank >= '1' && rank <= '8' {
		output.y = int(rank-'1') + 1
	} else {
		return xy{}, fmt.Errorf("rank coordinate is invalid")
	}

	return output, nil
}

//performs a move. it doesn't check if the move is legal
func (b *Board) move(m move) bool {
	for i, piece := range b.pieces {
		//if the move is a capture, remove the captured piece
		//and also reset the halfMoves field

		//TODO: this just straight up crashes
		if piece.position == m.to {
			b.pieces[len(b.pieces)-1], b.pieces[i] = b.pieces[i], b.pieces[len(b.pieces)-1]
			b.pieces = b.pieces[:len(b.pieces)-1]

			b.halfMoves = 0

			break
		}
	}

	for i, piece := range b.pieces {

		//finds the piece we're trying to move
		if piece.position == m.from {
			b.pieces[i].position = m.to

			if piece.pieceType == pawn {
				b.halfMoves = 0
				if m.moveType == promotion {
					b.pieces[i].pieceType = m.promotesInto
				}
			}
		}
	}

	if m.moveType == longCastle {
		rook, isThereAPiece := b.pointerPieceAtPosition(xy{1, 1})
		if isThereAPiece {
			rook.position = xy{4, 1}
		}
	}

	if m.moveType == shortCastle {
		rook, isThereAPiece := b.pointerPieceAtPosition(xy{8, 1})
		if isThereAPiece {
			rook.position = xy{6, 1}
		}
	}

	b.turn = !b.turn
	return true
}

func (b *Board) testMove(m move) bool {
	b.turn = !b.turn

	return b.move(m)
}

//Clone creates a Board that is identical to the input one
func (b Board) Clone() Board {
	newBoard := Board{
		pieces:         make([]piece, len(b.pieces)),
		turn:           b.turn,
		whiteCanCastle: b.whiteCanCastle,
		blackCanCastle: b.blackCanCastle,
		enPassant:      b.enPassant,
		halfMoves:      b.halfMoves,
		moveNumber:     b.moveNumber,
	}

	copy(newBoard.pieces, b.pieces)
	return newBoard
}

func (b Board) eliminateIllegalMoves(moves []move) []move {
	legalMoves := make([]move, 0, len(moves))
	for _, move := range moves {
		if b.isLegal(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

//to check if a move doesn't put one's own king in check,
//we just perform it and make sure the king isn't in check
//is this performant? no. is it easy? no lmao
func (b Board) isLegal(m move) bool {
	dummyBoard := b.Clone()
	dummyBoard.testMove(m)

	if dummyBoard.isKingInCheck(dummyBoard.turn) {
		return false
	}

	return true
}

//checks if a move is actually doable
//doesn't have anything to do with checks or pins
//just tells you if there's a piece in a that can go from a to b
func (b Board) isPossible(m move) bool {
	possibleMoves := b.moves(m.color)

	for _, move := range possibleMoves {
		if m.from == move.from && m.to == move.to {
			return true
		}
	}

	return false
}
