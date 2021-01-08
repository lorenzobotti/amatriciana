package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"../amatriciana"
)

func main() {
	//board := amatriciana.NewBoard()
	//board, _ := amatriciana.BoardFromFEN("3k4/7R/R7/8/8/3K4/8/8 w - - 2 2")
	board, _ := amatriciana.BoardFromFEN("4rkn1/p1Q2p1q/8/2pp4/5P2/1P4P1/PBbKB3/8 b - - 0 20")
	fmt.Println(board.FEN())

	reader := bufio.NewReader(os.Stdin)
	colors := [2]string{"white", "black"}
	_ = colors

	rand.Seed(time.Now().UnixNano())
	for {
		fmt.Println(("---- time to move ------"))
		fmt.Println(board.FEN())

		fmt.Printf("it's %s's turn to move\n", board.Turn())
		bestMove, err := board.BestMove(3)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("best move:", bestMove.String())
		fmt.Println("input a move")

		fmt.Println(board.Draw())
		move, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading input")
		}

		err = board.PerformMove(move)
		if err != nil {
			fmt.Println("couldn't perform move:", err.Error())
		} else {
			fmt.Println("move performed successfully")
		}

		if board.IsCheckmate() {
			fmt.Println(board.Turn(), "has won the game")
			break
		}
	}

}
