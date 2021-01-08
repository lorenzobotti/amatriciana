package amatriciana

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type color bool

const white = true
const black = false

func (c color) String() string {
	if c == white {
		return "white"
	} else {
		return "black"
	}
}

//Board describes the state of the game at any point. Has all the values of FEN notation
type Board struct {
	pieces         []piece
	turn           color
	whiteCanCastle [2]bool
	blackCanCastle [2]bool
	enPassant      xy
	moveNumber     int
	halfMoves      int
}

func (b Board) Turn() string {
	return b.turn.String()
}

func (b Board) pieceAtPosition(pos xy) (piece, bool) {
	for _, p := range b.pieces {
		if p.position == pos {
			return p, true
		}
	}

	return piece{}, false
}

func (b Board) pointerPieceAtPosition(pos xy) (*piece, bool) {
	for _, p := range b.pieces {
		if p.position == pos {
			return &p, true
		}
	}

	return nil, false
}

func (b Board) piecesOfColor(col color) []piece {
	pieces := make([]piece, 0)

	for _, piece := range b.pieces {
		if piece.color == col {
			pieces = append(pieces, piece)
		}
	}

	return pieces
}

func (b Board) moves(col color) []move {
	pieces := b.piecesOfColor(col)
	moves := make([]move, 0)

	for _, piece := range pieces {
		pieceMoves := piece.moves(b)
		for _, move := range pieceMoves {
			if b.isLegal(move) {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

//IsCheckmate tells you if the current player to move is in checkmate
func (b Board) IsCheckmate() bool {
	if !b.isKingInCheck(b.turn) {
		return false
	}

	moves := b.moves(b.turn)

	if len(moves) == 1 {
		return true
	}

	for _, move := range moves {
		dummyBoard := b.Clone()
		dummyBoard.move(move)

		if !dummyBoard.isKingInCheck(b.turn) {
			return false
		}
	}

	return true
}

func (b Board) isSquareInCheck(square xy, col color) bool {
	diagonals := [...]xy{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, direction := range diagonals {
		piece, something := b.pieceAtPosition(b.lastSquareInDirection(square, direction, col))
		if something && piece.color != col {
			if piece.pieceType == queen || piece.pieceType == bishop {
				return true
			}
		}
	}

	straightLines := [...]xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, direction := range straightLines {
		piece, something := b.pieceAtPosition(b.lastSquareInDirection(square, direction, col))
		if something && piece.color != col {
			if piece.pieceType == queen || piece.pieceType == rook {
				return true
			}
		}
	}

	knightMoves := [...]xy{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {2, -1}, {2, 1}, {1, -2}, {1, 2}}
	for _, direction := range knightMoves {
		piece, something := b.pieceAtPosition(square.plus(direction))
		if something && piece.color != col && piece.pieceType == knight {
			return true
		}
	}
	//      oooooooo
	//     o       ooooooooooooooooooooooooooo
	//     ooooooooo                     o    o
	//     o       o                     o    o
	//     o       ooooooooooooooooooooooooooo
	//     oooooooo

	var otherSide xy
	if col == white {
		otherSide = xy{0, 1}
	} else {
		otherSide = xy{0, -1}
	}

	pawnCaptures := [2]xy{xy{1, 0}.plus(otherSide), xy{-1, 0}.plus(otherSide)}
	for _, direction := range pawnCaptures {
		piece, something := b.pieceAtPosition(square.plus(direction))
		if something && piece.color != col && piece.pieceType == pawn {
			return true
		}
	}

	return false
}

func (b Board) isKingInCheck(col color) bool {
	kingPos, err := b.kingPosition(col)
	if err != nil {
		return false
	}

	return b.isSquareInCheck(kingPos, col)
}

func isInBounds(pos xy) bool {
	return pos.x <= 8 && pos.x >= 1 && pos.y <= 8 && pos.y >= 1
}

func (b Board) kingPosition(col color) (xy, error) {
	for _, p := range b.pieces {
		if p.pieceType == king && p.color == col {
			return p.position, nil
		}
	}

	var err error

	if col == white {
		err = fmt.Errorf("there's no white king, lmao")
	} else {
		err = fmt.Errorf("there's no black king, lmao")
	}

	return xy{}, err
}

//gives you a list of possible squares you can move to in a certain direction
//it stops when it finds a piece of the same color
func (b Board) squaresInDirection(from, dir xy, col color) []xy {
	squares := make([]xy, 0)
	for {
		from = from.plus(dir)

		piece, occupied := b.pieceAtPosition(from)
		if occupied && (piece.color == col || piece.pieceType == king) {
			return squares
		}
		if from.x > 8 || from.x < 1 || from.y > 8 || from.y < 1 {
			return squares
		}

		squares = append(squares, from)
		if occupied {
			return squares
		}
	}
}

func (b Board) lastSquareInDirection(from, dir xy, col color) xy {
	squares := b.squaresInDirection(from, dir, col)

	if len(squares) == 0 {
		return from
	}
	return squares[len(squares)-1]
}

//gives you a list of moves you can make in a certain direction
func (b Board) movesInDirection(from, dir xy, col color) []move {
	piece, _ := b.pieceAtPosition(from)

	squares := b.squaresInDirection(from, dir, col)
	moves := make([]move, len(squares))

	for i, square := range squares {
		moves[i] = move{piece.pieceType, col, from, square, normalMove, pawn}
	}

	return moves
}

//Draw the board in ASCII art
func (b Board) Draw() string {
	var board [8][8]byte

	var output bytes.Buffer

	for _, piece := range b.pieces {
		pos := piece.position

		if !isInBounds(pos) {
			continue
		}

		board[pos.x-1][pos.y-1] = piece.fenLetter()
	}

	for rank := 7; rank >= 0; rank-- {
		output.WriteByte('|')
		for file := 0; file < 8; file++ {
			piece := board[file][rank]
			if piece != 0 {
				output.WriteByte(piece)
			} else {
				output.WriteByte(' ')
			}
			output.WriteByte('|')
		}
		output.WriteString("\n|-+-+-+-+-+-+-+-|\n")
	}

	return output.String()
}

//FEN exports the board in FEN notation
func (b Board) FEN() string {
	var board [8][8]byte
	for _, piece := range b.pieces {
		pos := piece.position

		if !isInBounds(pos) {
			continue
		}

		board[pos.x-1][pos.y-1] = piece.fenLetter()
	}

	var output bytes.Buffer

	for i := 7; i >= 0; i-- {
		var rank [8]byte
		for j := 0; j < 8; j++ {
			rank[j] = board[j][i]
		}

		var emptySquares byte = 0
		for idx, char := range rank {
			if char == 0 {
				emptySquares++
				if idx == 7 {
					output.WriteByte('0' + emptySquares)
				}

				continue
			}

			if emptySquares > 0 {
				output.WriteByte('0' + emptySquares)
			}

			output.WriteByte(char)
			emptySquares = 0
		}

		if i != 0 {
			output.WriteByte('/')
		}
	}

	output.WriteByte(' ')
	if b.turn == white {
		output.WriteByte('w')
	} else {
		output.WriteByte('b')
	}

	output.WriteByte(' ')

	if !b.whiteCanCastle[0] && !b.whiteCanCastle[1] &&
		!b.blackCanCastle[0] && !b.blackCanCastle[1] {
		output.WriteByte('-')
	}
	if b.whiteCanCastle[0] {
		output.WriteByte('K')
	}
	if b.whiteCanCastle[1] {
		output.WriteByte('Q')
	}
	if b.blackCanCastle[0] {
		output.WriteByte('k')
	}
	if b.blackCanCastle[1] {
		output.WriteByte('q')
	}

	output.WriteByte(' ')
	if (b.enPassant == xy{0, 0}) {
		output.WriteByte('-')
	} else {
		output.WriteString(b.enPassant.String())
	}

	output.WriteByte(' ')
	output.WriteString(strconv.Itoa(b.halfMoves))

	output.WriteByte(' ')
	output.WriteString(strconv.Itoa(b.moveNumber))

	return output.String()
}

func (b Board) EightByEight() [8][8]byte {
	var board [8][8]byte
	for _, piece := range b.pieces {
		pos := piece.position

		if !isInBounds(pos) {
			continue
		}

		board[pos.x-1][pos.y-1] = piece.fenLetter()
	}

	return board
}

//BoardFromFEN creates a new Board from a FEN string
func BoardFromFEN(fen string) (Board, error) {
	elements := strings.Split(fen, " ")
	if len(elements) != 6 {
		return Board{}, fmt.Errorf("incorrect number of fields")
	}

	fenBoard := elements[0]

	ranks := strings.Split(fenBoard, "/")
	if len(ranks) != 8 {
		return Board{}, fmt.Errorf("incorrect number of ranks in board")
	}

	board := Board{}
	board.pieces = make([]piece, 0)

	//generating the board
	for rankNum, rankStr := range ranks {
		rank := []byte(rankStr)
		currentFile := 0

		for _, char := range rank {
			if char > '0' && char < '9' {
				currentFile += int(char - '0')
				continue
			}

			piece, err := pieceFromFen(char)
			if err != nil {
				return Board{}, err
			}

			piece.position = xy{currentFile + 1, 8 - rankNum}
			board.pieces = append(board.pieces, piece)
			currentFile++

		}
	}

	//whose turn is it
	if elements[1] == "w" {
		board.turn = white
	} else if elements[1] == "b" {
		board.turn = black
	} else {
		return board, fmt.Errorf("active color isn't \"w\" or \"b\"")
	}

	castling := elements[2]
	if strings.Contains(castling, "K") {
		board.whiteCanCastle[0] = true
	}
	if strings.Contains(castling, "k") {
		board.blackCanCastle[0] = true
	}
	if strings.Contains(castling, "Q") {
		board.whiteCanCastle[0] = true
	}
	if strings.Contains(castling, "q") {
		board.blackCanCastle[0] = true
	}

	enPassant := elements[3]
	_ = enPassant

	halfMoves, err := strconv.Atoi(elements[4])
	if err != nil {
		return board, err
	}
	board.halfMoves = halfMoves

	moveNumber, err := strconv.Atoi(elements[4])
	if err != nil {
		return board, err
	}
	board.moveNumber = moveNumber

	return board, nil
}

//NewBoard creates a new board with the default configuration from scratch
func NewBoard() Board {
	board := Board{
		pieces:         make([]piece, 32),
		turn:           white,
		whiteCanCastle: [2]bool{true, true},
		blackCanCastle: [2]bool{true, true},
	}

	board.pieces[0] = piece{xy{1, 1}, white, rook}
	board.pieces[1] = piece{xy{2, 1}, white, knight}
	board.pieces[2] = piece{xy{3, 1}, white, bishop}
	board.pieces[3] = piece{xy{4, 1}, white, queen}
	board.pieces[4] = piece{xy{5, 1}, white, king}
	board.pieces[5] = piece{xy{6, 1}, white, bishop}
	board.pieces[6] = piece{xy{7, 1}, white, knight}
	board.pieces[7] = piece{xy{8, 1}, white, rook}

	board.pieces[8] = piece{xy{1, 8}, black, rook}
	board.pieces[9] = piece{xy{2, 8}, black, knight}
	board.pieces[10] = piece{xy{3, 8}, black, bishop}
	board.pieces[11] = piece{xy{4, 8}, black, queen}
	board.pieces[12] = piece{xy{5, 8}, black, king}
	board.pieces[13] = piece{xy{6, 8}, black, bishop}
	board.pieces[14] = piece{xy{7, 8}, black, knight}
	board.pieces[15] = piece{xy{8, 8}, black, rook}

	board.pieces[16] = piece{xy{1, 2}, white, pawn}
	board.pieces[17] = piece{xy{2, 2}, white, pawn}
	board.pieces[18] = piece{xy{3, 2}, white, pawn}
	board.pieces[19] = piece{xy{4, 2}, white, pawn}
	board.pieces[20] = piece{xy{5, 2}, white, pawn}
	board.pieces[21] = piece{xy{6, 2}, white, pawn}
	board.pieces[22] = piece{xy{7, 2}, white, pawn}
	board.pieces[23] = piece{xy{8, 2}, white, pawn}

	board.pieces[24] = piece{xy{1, 7}, black, pawn}
	board.pieces[25] = piece{xy{2, 7}, black, pawn}
	board.pieces[26] = piece{xy{3, 7}, black, pawn}
	board.pieces[27] = piece{xy{4, 7}, black, pawn}
	board.pieces[28] = piece{xy{5, 7}, black, pawn}
	board.pieces[29] = piece{xy{6, 7}, black, pawn}
	board.pieces[30] = piece{xy{7, 7}, black, pawn}
	board.pieces[31] = piece{xy{8, 7}, black, pawn}

	return board
}
