package main

import (
	"fmt"
)

func main() {
	board := Board {
		{9, 9, 9, 9, 9, 9, 9, 9},
		{9, 0, 0, 0, 4, 5, 6, 9},
		{9, 6, 1, 2, 3, 0, 5, 9},
		{9, 3, 3, 3, 3, 0, 1, 9},
		{9, 4, 4, 4, 4, 0, 1, 9},
		{9, 5, 1, 5, 2, 5, 1, 9},
		{9, 9, 9, 9, 9, 9, 9, 9},
	}

	fmt.Println(board)
	board = delete(board)
	fmt.Println(board)
}

type Board [][]int

func delete(board Board) Board {
	var deletedBoard [][]int
	for _, b1 := range board {
		var dB []int
		for _, b2 := range b1 {
			dB = append(dB, b2)
		}
		deletedBoard = append(deletedBoard, dB)
	}

	// 横方向の削除
	for i := 1; i < 6; i++ {
		for j := 2; j < 6; j++ {
			dropL := board[i][j-1]
			drop := board[i][j]
			dropR := board[i][j+1]

			if drop == 0 {
				continue
			}

			if dropL == drop && drop == dropR {
				deletedBoard[i][j-1] = 0
				deletedBoard[i][j] = 0
				deletedBoard[i][j+1] = 0
			}
		}
	}

	// 縦方向の削除
	for i := 2; i < 5; i++ {
		for j := 1; j < 7; j++ {
			dropU := board[i-1][j]
			drop := board[i][j]
			dropD := board[i+1][j]

			if drop == 0 {
				continue
			}

			if dropU == drop && drop == dropD {
				deletedBoard[i-1][j] = 0
				deletedBoard[i][j] = 0
				deletedBoard[i+1][j] = 0
			}
		}
	}

	return deletedBoard
}
