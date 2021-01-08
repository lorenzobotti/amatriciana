package amatriciana

import "fmt"

type piece struct {
	position xy
	color
	pieceType
}

type pieceType int

const (
	pawn pieceType = iota
	knight
	bishop
	rook
	queen
	king
)

func (pt pieceType) String() string {
	switch pt {
	case pawn:
		return "pawn"
	case knight:
		return "knight"
	case bishop:
		return "bishop"
	case rook:
		return "rook"
	case queen:
		return "queen"
	case king:
		return "king"
	default:
		return "???"
	}
}

func (pt pieceType) letter() byte {
	switch pt {
	case pawn:
		return 'p'
	case knight:
		return 'n'
	case bishop:
		return 'b'
	case rook:
		return 'r'
	case queen:
		return 'q'
	case king:
		return 'k'
	default:
		return '?'
	}
}

func pieceTypeFromFen(char byte) (pieceType, error) {
	switch char {
	case 'P':
		return pawn, nil
	case 'p':
		return pawn, nil
	case 'N':
		return knight, nil
	case 'n':
		return knight, nil
	case 'B':
		return bishop, nil
	case 'b':
		return bishop, nil
	case 'R':
		return rook, nil
	case 'r':
		return rook, nil
	case 'Q':
		return queen, nil
	case 'q':
		return queen, nil
	case 'K':
		return king, nil
	case 'k':
		return king, nil
	default:
		panic(fmt.Errorf("invalid piece : %c", char))
		return 0, fmt.Errorf("invalid piece : %c", char)

	}
}

func pieceFromFen(char byte) (piece, error) {
	switch char {
	case 'P':
		return piece{xy{}, white, pawn}, nil
	case 'p':
		return piece{xy{}, black, pawn}, nil
	case 'N':
		return piece{xy{}, white, knight}, nil
	case 'n':
		return piece{xy{}, black, knight}, nil
	case 'B':
		return piece{xy{}, white, bishop}, nil
	case 'b':
		return piece{xy{}, black, bishop}, nil
	case 'R':
		return piece{xy{}, white, rook}, nil
	case 'r':
		return piece{xy{}, black, rook}, nil
	case 'Q':
		return piece{xy{}, white, queen}, nil
	case 'q':
		return piece{xy{}, black, queen}, nil
	case 'K':
		return piece{xy{}, white, king}, nil
	case 'k':
		return piece{xy{}, black, king}, nil
	default:
		return piece{}, fmt.Errorf("invalid piece : %c", char)

	}
}

func (p piece) String() string {
	return fmt.Sprintf("%s %s in %s", p.color.String(), p.pieceType.String(), p.position.String())
}

func (p piece) moves(b Board) []move {
	moves := make([]move, 0)

	switch p.pieceType {
	case rook:
		moves = b.rookMoves(p.position, p.color)
	case bishop:
		moves = b.bishopMoves(p.position, p.color)
	case knight:
		moves = b.knightMoves(p.position, p.color)
	case queen:
		moves = b.queenMoves(p.position, p.color)
	case pawn:
		moves = b.pawnMoves(p.position, p.color)
	case king:
		moves = b.kingMoves(p.position, p.color)
	}

	if len(moves) == 0 {
		return moves
	}

	/*for _, moveToCheck := range moves {
		if !b.isLegal(moveToCheck) {
			for i, moveToRemove := range moves {
				if moveToRemove == moveToCheck {
					moves[i] = moves[len(moves)-1]
					moves = moves[:len(moves)-1]
					break
				}
			}
			break
		}
	}*/

	legalMoves := make([]move, 0, len(moves))
	for _, move := range moves {
		if b.isLegal(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

func (b Board) rookMoves(pos xy, col color) []move {
	moves := b.movesInDirection(pos, xy{1, 0}, col)
	moves = append(moves, b.movesInDirection(pos, xy{-1, 0}, col)...)
	moves = append(moves, b.movesInDirection(pos, xy{0, 1}, col)...)
	moves = append(moves, b.movesInDirection(pos, xy{0, -1}, col)...)
	return moves
}

func (b Board) bishopMoves(pos xy, col color) []move {
	moves := b.movesInDirection(pos, xy{1, 1}, col)
	moves = append(moves, b.movesInDirection(pos, xy{-1, -1}, col)...)
	moves = append(moves, b.movesInDirection(pos, xy{-1, 1}, col)...)
	moves = append(moves, b.movesInDirection(pos, xy{1, -1}, col)...)
	return moves
}

func (b Board) knightMoves(pos xy, col color) []move {
	moves := make([]move, 0)

	directions := [...]xy{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {2, -1}, {2, 1}, {1, -2}, {1, 2},
	}

	for _, square := range directions {
		if !isInBounds(pos.plus(square)) {
			continue
		}

		piece, occupied := b.pieceAtPosition(pos.plus(square))
		if occupied && piece.color == col {
			continue
		}

		move := move{knight, col, pos, pos.plus(square), normalMove, pawn}
		moves = append(moves, move)
	}

	return moves
}

func (b Board) queenMoves(pos xy, col color) []move {
	moves := b.bishopMoves(pos, col)
	moves = append(moves, b.rookMoves(pos, col)...)
	return moves
}

func (b Board) pawnMoves(pos xy, col color) []move {
	var otherSide xy
	moves := make([]move, 0)

	if col == white {
		otherSide = xy{0, 1}
	} else {
		otherSide = xy{0, -1}
	}

	moveType := normalMove
	if (pos.plus(otherSide).y > 7 && col == white) || (pos.plus(otherSide).y < 2 && col == black) {
		moveType = promotion
	}

	//check for captures
	right, canCaptureRight := b.pieceAtPosition(pos.plus(otherSide).plus(xy{1, 0}))
	if canCaptureRight && right.color != col {
		moves = append(moves, move{pawn, col, pos, right.position, moveType, pawn})
	}

	left, canCaptureLeft := b.pieceAtPosition(pos.plus(otherSide).plus(xy{-1, 0}))
	if canCaptureLeft && left.color != col {
		moves = append(moves, move{pawn, col, pos, left.position, moveType, pawn})
	}

	//check if it can move forwards by one
	_, cantPush := b.pieceAtPosition(pos.plus(otherSide))
	if !cantPush {
		moves = append(moves, move{pawn, col, pos, pos.plus(otherSide), moveType, pawn})

		if (pos.y == 2 && col == white) || (pos.y == 7 && col == black) {
			pushTwoPos := pos.plus(otherSide).plus(otherSide)
			_, cantPush := b.pieceAtPosition(pushTwoPos)
			if !cantPush {
				moves = append(moves, move{pawn, col, pos, pushTwoPos, moveType, pawn})
			}
		}
	}

	return moves
}

func (b Board) pawnCaptures(pos xy, col color) []move {
	var otherSide xy

	if col == white {
		otherSide = xy{0, 1}
	} else {
		otherSide = xy{0, -1}
	}

	//check for captures
	right := pos.plus(otherSide).plus(xy{1, 0})
	left := pos.plus(otherSide).plus(xy{-1, 0})

	moves := make([]move, 0, 2)
	if isInBounds(right) {
		moves = append(moves, move{pawn, col, pos, right, normalMove, pawn})
	}
	if isInBounds(left) {
		moves = append(moves, move{pawn, col, pos, left, normalMove, pawn})
	}

	return moves
}

func (b Board) kingMoves(pos xy, col color) []move {

	var canCastle [2]bool
	moves := make([]move, 0)
	kingPos, err := b.kingPosition(col)

	if err != nil {
		return moves
	}

	switch col {
	case white:
		canCastle = b.whiteCanCastle
	case black:
		canCastle = b.blackCanCastle
	}

	if canCastle[0] {
		_, occupied1 := b.pieceAtPosition(kingPos.plus(xy{1, 0}))
		_, occupied2 := b.pieceAtPosition(kingPos.plus(xy{2, 0}))
		if !b.isKingInCheck(col) &&
			!b.isSquareInCheck(kingPos.plus(xy{1, 0}), col) &&
			!occupied1 &&
			!b.isSquareInCheck(kingPos.plus(xy{2, 0}), col) &&
			!occupied2 {
			moves = append(moves, move{king, col, kingPos, kingPos.plus(xy{2, 0}), shortCastle, pawn})
		}
	}
	if canCastle[1] {
		_, occupied1 := b.pieceAtPosition(kingPos.plus(xy{-1, 0}))
		_, occupied2 := b.pieceAtPosition(kingPos.plus(xy{-2, 0}))
		if !b.isKingInCheck(col) &&
			!b.isSquareInCheck(kingPos.plus(xy{-1, 0}), col) &&
			!b.isSquareInCheck(kingPos.plus(xy{-2, 0}), col) &&
			!occupied1 && occupied2 {
			moves = append(moves, move{king, col, kingPos, kingPos.plus(xy{-2, 0}), longCastle, pawn})
		}
	}

	directions := [...]xy{
		{1, 1},
		{1, 0},
		{1, -1},
		{0, 1},
		{0, -1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
	}

	for _, direction := range directions {
		square := kingPos.plus(direction)

		piece, occupied := b.pieceAtPosition(square)
		_ = piece
		if occupied && piece.color == col {
			continue
		}
		if !b.isSquareInCheck(square, col) && isInBounds(square) {
			moves = append(moves, move{king, col, kingPos, square, normalMove, pawn})
		}
	}

	return moves
}

func (p piece) fenLetter() byte {
	letter := p.pieceType.letter()

	if p.color == white {
		return letter + 'A' - 'a'
	} else {
		return letter
	}
}
