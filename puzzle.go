package main

import (
	"fmt"
	"strconv"
)

const depth int = 7

var (
	moves    [depth+1]string
	maxPoint int = 0
	maxMoves [depth+1]string
)

func main() {
	board := Board {
		{9, 9, 9, 9, 9, 9, 9, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 9, 9, 9, 9, 9, 9, 9},
	}

	// 盤面を入力
	for i := 1; i < 6; i++ {
		var b1, b2, b3, b4, b5, b6 string
		fmt.Print("> ")
		fmt.Scan(&b1, &b2, &b3, &b4, &b5, &b6)

		board[i][1], _ = strconv.Atoi(b1)
		board[i][2], _ = strconv.Atoi(b2)
		board[i][3], _ = strconv.Atoi(b3)
		board[i][4], _ = strconv.Atoi(b4)
		board[i][5], _ = strconv.Atoi(b5)
		board[i][6], _ = strconv.Atoi(b6)
	}

	// ルートを探索
	for i := 1; i < 6; i++ {
		for j := 1; j < 7; j++ {
			moves[0] = strconv.Itoa(j) + strconv.Itoa(i)
			move(1, board, j, i)
		}
	}

	board.print()
	fmt.Println(maxMoves)
}

// Board 盤面
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
	if moves[n-1] != "←" && dropR != 9 {
		movedBoard[y][x] = dropR
		movedBoard[y][x+1] = drop
		moves[n] = "→"
		move(n+1, movedBoard, x+1, y)
		movedBoard[y][x] = drop
		movedBoard[y][x+1] = dropR
		moves[n] = ""
	}

	// Down
	if moves[n-1] != "↑" && dropD != 9 {
		movedBoard[y][x] = dropD
		movedBoard[y+1][x] = drop
		moves[n] = "↓"
		move(n+1, movedBoard, x, y+1)
		movedBoard[y][x] = drop
		movedBoard[y+1][x] = dropD
		moves[n] = ""
	}

	// Left
	if moves[n-1] != "→" && dropL != 9 {
		movedBoard[y][x] = dropL
		movedBoard[y][x-1] = drop
		moves[n] = "←"
		move(n+1, movedBoard, x-1, y)
		movedBoard[y][x] = drop
		movedBoard[y][x-1] = dropL
		moves[n] = ""
	}

	// Up
	if moves[n-1] != "↓" && dropU != 9 {
		movedBoard[y][x] = dropU
		movedBoard[y-1][x] = drop
		moves[n] = "↑"
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
