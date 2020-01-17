package main

import (
	"fmt"
	"strconv"
)

const depth int = 18

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
	fmt.Println("")

	// ルートを探索
	for i := 1; i < 6; i++ {
		for j := 1; j < 7; j++ {
			moves[0] = strconv.Itoa(j) + strconv.Itoa(i)
			move(1, board, j, i)
		}
	}

	board.print()
	fmt.Println(maxMoves)
	fmt.Println("消える数: " + strconv.Itoa(maxPoint))
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
	} else if n == 9 {
		// 枝刈り。8回移動して6ポイント以下のルートは探索しない。
		point := calcPoint(board)
		if point < 6 {
			return
		}
	} else if n == 13 {
		// 枝刈り2。12回移動して9ポイント以下のルートは探索しない。
		point := calcPoint(board)
		if point < 9 {
			return
		}
	}

	drop := board[y][x]
	dropR := board[y][x+1]
	dropD := board[y+1][x]
	dropL := board[y][x-1]
	dropU := board[y-1][x]

	// Right
	if moves[n-1] != "←" && dropR != 9 {
		board[y][x] = dropR
		board[y][x+1] = drop
		moves[n] = "→"
		move(n+1, board, x+1, y)
		board[y][x] = drop
		board[y][x+1] = dropR
		moves[n] = ""
	}

	// Down
	if moves[n-1] != "↑" && dropD != 9 {
		board[y][x] = dropD
		board[y+1][x] = drop
		moves[n] = "↓"
		move(n+1, board, x, y+1)
		board[y][x] = drop
		board[y+1][x] = dropD
		moves[n] = ""
	}

	// Left
	if moves[n-1] != "→" && dropL != 9 {
		board[y][x] = dropL
		board[y][x-1] = drop
		moves[n] = "←"
		move(n+1, board, x-1, y)
		board[y][x] = drop
		board[y][x-1] = dropL
		moves[n] = ""
	}

	// Up
	if moves[n-1] != "↓" && dropU != 9 {
		board[y][x] = dropU
		board[y-1][x] = drop
		moves[n] = "↑"
		move(n+1, board, x, y-1)
		board[y][x] = drop
		board[y-1][x] = dropU
		moves[n] = ""
	}
}

// delete ドロップを消す
func delete(board Board) (Board, bool) {
	isWork := false
	deletedBoard := Board {
		{9, 9, 9, 9, 9, 9, 9, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 0, 0, 0, 0, 0, 0, 9},
		{9, 9, 9, 9, 9, 9, 9, 9},
	}

	for i := 1; i < 6; i++ {
		for j := 1; j < 7; j++ {
			deletedBoard[i][j] = board[i][j]
		}
	}

	// 横方向の削除
	for i := 1; i < 6; i++ {
		for j := 2; j < 6; j++ {
			drop := board[i][j]

			if drop == 0 {
				continue
			}

			dropL := board[i][j-1]
			dropR := board[i][j+1]

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
			drop := board[i][j]

			if drop == 0 {
				continue
			}

			dropU := board[i-1][j]
			dropD := board[i+1][j]

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
	var isWork       bool

	// 消えるドロップがなくなるまでloop
	point := 0
	for {
		board, isWork = delete(board)
		if isWork {
			board = fall(board)
			continue
		}

		// 消えなくなったらpointを計算
		for i := 1; i < 6; i++ {
			for j := 1; j < 7; j++ {
				if board[i][j] == 0 {
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
