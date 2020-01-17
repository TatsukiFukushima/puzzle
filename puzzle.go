package main

import (
	"fmt"
	"strconv"
)

const depth int = 3

var (
	moves    [depth+1]string
	maxPoint int = 0
	maxMoves [depth+1]string
)

func main() {
	board := Board {
		{9, 9, 9, 9, 9, 9, 9, 9},
		{9, 1, 2, 2, 4, 5, 6, 9},
		{9, 6, 1, 2, 3, 3, 5, 9},
		{9, 3, 1, 3, 3, 2, 1, 9},
		{9, 4, 4, 4, 4, 3, 1, 9},
		{9, 5, 1, 5, 2, 5, 1, 9},
		{9, 9, 9, 9, 9, 9, 9, 9},
	}

	for i := 1; i < 6; i++ {
		for j := 1; j < 7; j++ {
			moves[0] = strconv.Itoa(j) + strconv.Itoa(i)
			move(1, board, j, i)
		}
	}

	fmt.Println(maxMoves)
}

type Board [][]int

// move ドロップを移動させる。手数、盤面、x座標、y座標
func move(n int, board Board, x int, y int) {
	// 指定の手数以上になった場合、ポイントを算出して最高得点だったらmovesを記録。
	if n > depth {
		point := calcPoint(board)
		if point > maxPoint {
			maxPoint = point
			maxMoves = moves
		}
		return
	}

	var movedBoard Board
	for _, b1 := range board {
		var dB []int
		for _, b2 := range b1 {
			dB = append(dB, b2)
		}
		movedBoard = append(movedBoard, dB)
	}

	drop := board[y][x]
	dropR := board[y][x+1]
	dropD := board[y+1][x]
	dropL := board[y][x-1]
	dropU := board[y-1][x]

	// Right
	if dropR != 9 {
		movedBoard[y][x] = dropR
		movedBoard[y][x+1] = drop
		moves[n] = "R"
		move(n+1, movedBoard, x+1, y)
		movedBoard[y][x] = drop
		movedBoard[y][x+1] = dropR
		moves[n] = ""
	}

	// Down
	if dropD != 9 {
		movedBoard[y][x] = dropD
		movedBoard[y+1][x] = drop
		moves[n] = "D"
		move(n+1, movedBoard, x, y+1)
		movedBoard[y][x] = drop
		movedBoard[y+1][x] = dropD
		moves[n] = ""
	}

	// Left
	if dropL != 9 {
		movedBoard[y][x] = dropL
		movedBoard[y][x-1] = drop
		moves[n] = "L"
		move(n+1, movedBoard, x-1, y)
		movedBoard[y][x] = drop
		movedBoard[y][x-1] = dropL
		moves[n] = ""
	}

	// Up
	if dropU != 9 {
		movedBoard[y][x] = dropU
		movedBoard[y-1][x] = drop
		moves[n] = "U"
		move(n+1, movedBoard, x, y-1)
		movedBoard[y][x] = drop
		movedBoard[y-1][x] = dropU
		moves[n] = ""
	}
}

// delete ドロップを消す
func delete(board Board) (Board, bool) {
	var deletedBoard [][]int
	isWork := false
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
				isWork = true
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
				isWork = true
			}
		}
	}

	return deletedBoard, isWork
}

// fall ドロップを落とす
func fall(board Board) Board {
	fallenBoard := Board {
		{9, 9, 9, 9, 9, 9, 9, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 9, 9, 9, 9, 9, 9, 9},
	}

	// 列ごとにループを回しているということを明示するためにjとしている。
	for j := 1; j < 7; j++ {
		// ドロップが落ちる先のy座標、的なニュアンス
		nextY := 5
		for i := 5; i > 0; i-- {
			drop := board[i][j]
			if drop != 0 {
				fallenBoard[nextY][j] = drop
				nextY--
			}
		}
	}

	return fallenBoard
}

// calcPoint どれだけ消えたかを計算
func calcPoint(board Board) int {
	var (
		deletedBoard [][]int
		isWork       bool
	)
	for _, b1 := range board {
		var dB []int
		for _, b2 := range b1 {
			dB = append(dB, b2)
		}
		deletedBoard = append(deletedBoard, dB)
	}

	// 消えるドロップがなくなるまでloop
	point := 0
	for {
		deletedBoard, isWork = delete(deletedBoard)
		if isWork {
			deletedBoard = fall(deletedBoard)
			continue
		}

		// 消えなくなったらpointを計算
		for i := 1; i < 6; i++ {
			for j := 1; j < 7; j++ {
				if deletedBoard[i][j] == 0 {
					point++
				}
			}
		}
		break
	}

	return point
}

// print 配列をきれいに出力。テスト用かも。
func (board Board) print() {
	for _, b := range board {
		fmt.Println(b)
	}
	fmt.Println("---------------------")
}
