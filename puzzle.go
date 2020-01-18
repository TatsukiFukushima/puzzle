package main

import (
	"fmt"
	"strconv"
)

const depth int = 24
var	rate float32 = 0.0

var (
	moves     [depth+1]string
	minPoint  int = 30
	minPoint2 int = 30
	minPoint3 int = 30
	minMoves  [depth+1]string
	minMoves2 [depth+1]string
	minMoves3 [depth+1]string
)

func main() {
	board := Board {
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}

	// 盤面を入力
	for i := 0; i < 5; i++ {
		var b1, b2, b3, b4, b5, b6 string
		fmt.Print("> ")
		fmt.Scan(&b1, &b2, &b3, &b4, &b5, &b6)

		board[i][0], _ = strconv.Atoi(b1)
		board[i][1], _ = strconv.Atoi(b2)
		board[i][2], _ = strconv.Atoi(b3)
		board[i][3], _ = strconv.Atoi(b4)
		board[i][4], _ = strconv.Atoi(b5)
		board[i][5], _ = strconv.Atoi(b6)
	}
	fmt.Println("")
	board.print()

	// ルートを探索
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			moves[0] = strconv.Itoa(j+1) + strconv.Itoa(i+1)
			move(1, board, j, i)
			rate += 10/3.0
			log := strconv.Itoa(int(rate)) + "% 最大値:" + strconv.Itoa(30 - minPoint)
			fmt.Printf("\r%s", log)
		}
	}

	fmt.Printf("\r")
	fmt.Println("---------------------")
	fmt.Println("候補１:")
	fmt.Println(minMoves)
	fmt.Println("消える数: " + strconv.Itoa(30 - minPoint))
	fmt.Println("")
	fmt.Println("候補２:")
	fmt.Println(minMoves2)
	fmt.Println("消える数: " + strconv.Itoa(30 - minPoint2))
	fmt.Println("")
	fmt.Println("候補３:")
	fmt.Println(minMoves3)
	fmt.Println("消える数: " + strconv.Itoa(30 - minPoint3))
}

// Board 盤面
type Board [][]int

// move ドロップを移動させる。手数、盤面、x座標、y座標
func move(n int, board Board, x int, y int) {
	point := 0
	// 指定の手数以上になった場合、ポイントを算出して最高得点だったらmovesを記録。
	if n > depth {
		point = calcPoint(board)
		if point < minPoint {
			minPoint = point
			minMoves = moves
			log := strconv.Itoa(int(rate)) + "% 最大値:" + strconv.Itoa(30 - minPoint)
			fmt.Printf("\r%s", log)
		} else if point < minPoint2 {
			minPoint2 = point
			minMoves2 = moves
		} else if point < minPoint3 {
			minPoint3 = point
			minMoves3 = moves
		}
		return
	} else if n == 9 {
		// 枝刈り。8回移動して24ポイントより大きいルートは探索しない。
		point = calcPoint(board)
		if point > 24 {
			return
		}
	} else if n == 13 {
		// 枝刈り2。12回移動して21ポイントより大きいルートは探索しない。
		point = calcPoint(board)
		if point > 21 {
			return
		}
	} else if n == 17 {
		// 枝刈り3。16回移動して18ポイントより大きいルートは探索しない。
		point = calcPoint(board)
		if point > 18 {
			return
		}
	} else if n == 21 {
		// 枝刈り4。20回移動して18ポイントより大きいルートは探索しない。
		point = calcPoint(board)
		if point > 18 {
			return
		}
	} else if n > depth - 3 {
		// 探索深さの3手以内なら全てのルートを評価。
		if point == 0 {
			point = calcPoint(board)
		}
		if point < minPoint {
			minPoint = point
			minMoves = moves
		} else if point < minPoint2 {
			minPoint2 = point
			minMoves2 = moves
		} else if point < minPoint3 {
			minPoint3 = point
			minMoves3 = moves
		}
	}

	drop := board[y][x]
	// 高速化のため。
	nPlus := n+1
	nMinus := n-1
	xPlus := x+1
	xMinus := x-1
	yPlus := y+1
	yMinus := y-1

	// Right
	if moves[nMinus] != "←" && x != 5 {
		dropR := board[y][xPlus]
		board[y][x] = dropR
		board[y][xPlus] = drop
		moves[n] = "→"
		move(nPlus, board, xPlus, y)
		board[y][x] = drop
		board[y][xPlus] = dropR
		moves[n] = ""
	}

	// Down
	if moves[nMinus] != "↑" && y != 4 {
		dropD := board[yPlus][x]
		board[y][x] = dropD
		board[yPlus][x] = drop
		moves[n] = "↓"
		move(nPlus, board, x, yPlus)
		board[y][x] = drop
		board[yPlus][x] = dropD
		moves[n] = ""
	}

	// Left
	if moves[nMinus] != "→" && x != 0 {
		dropL := board[y][xMinus]
		board[y][x] = dropL
		board[y][xMinus] = drop
		moves[n] = "←"
		move(nPlus, board, xMinus, y)
		board[y][x] = drop
		board[y][xMinus] = dropL
		moves[n] = ""
	}

	// Up
	if moves[nMinus] != "↓" && y != 0 {
		dropU := board[yMinus][x]
		board[y][x] = dropU
		board[yMinus][x] = drop
		moves[n] = "↑"
		move(nPlus, board, x, yMinus)
		board[y][x] = drop
		board[yMinus][x] = dropU
		moves[n] = ""
	}
}

// delete ドロップを消す
func delete(board Board) (Board, bool) {
	isWork := false

	// 非常に汚い書き方だが、多分これの方が高速。
	deletedBoard := Board {
		{board[0][0], board[0][1], board[0][2], board[0][3], board[0][4], board[0][5]},
		{board[1][0], board[1][1], board[1][2], board[1][3], board[1][4], board[1][5]},
		{board[2][0], board[2][1], board[2][2], board[2][3], board[2][4], board[2][5]},
		{board[3][0], board[3][1], board[3][2], board[3][3], board[3][4], board[3][5]},
		{board[4][0], board[4][1], board[4][2], board[4][3], board[4][4], board[4][5]},
	}

	// 横方向の削除
	for i := 0; i < 5; i++ {
		for j := 1; j < 5; j++ {
			drop := board[i][j]
			jPlus := j+1
			jMinus := j-1

			if drop == 0 {
				continue
			}

			if board[i][jMinus] == drop && drop == board[i][jPlus] {
				deletedBoard[i][jMinus] = 0
				deletedBoard[i][j] = 0
				deletedBoard[i][jPlus] = 0
				isWork = true
			}
		}
	}

	// 縦方向の削除
	for i := 1; i < 4; i++ {
		for j := 0; j < 6; j++ {
			drop := board[i][j]
			iPlus := i+1
			iMinus := i-1

			if drop == 0 {
				continue
			}

			if board[iMinus][j] == drop && drop == board[iPlus][j] {
				deletedBoard[iMinus][j] = 0
				deletedBoard[i][j] = 0
				deletedBoard[iPlus][j] = 0
				isWork = true
			}
		}
	}

	return deletedBoard, isWork
}

// fall ドロップを落とす
func fall(board Board) Board {
	fallenBoard := Board {
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}

	// 列ごとにループを回しているということを明示するためにjとしている。
	for j := 0; j < 6; j++ {
		// ドロップが落ちる先のy座標、的なニュアンス
		nextY := 4
		for i := 4; i >= 0; i-- {
			drop := board[i][j]
			if drop != 0 {
				fallenBoard[nextY][j] = drop
				nextY--
			}
		}
	}

	return fallenBoard
}

// calcPoint どれだけ残るかを計算
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
		for i := 0; i < 5; i++ {
			for j := 0; j < 6; j++ {
				if board[i][j] != 0 {
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
}
