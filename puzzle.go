package main

import (
	"fmt"
	"strconv"
	"time"
)

// 変更前: 6.625秒 25個
// [53 ↓ ↓ → ↑ ← ← ← ↓ ← ↑ ↑ ↑ ↑ → ↓ → ↑ → ↓ ↓ ← ← ← ←]

const depth int = 24
var bestMoves BestMoves

var (
	minPoint  int = 30
	minPoint2 int = 30
	minPoint3 int = 30
	minMoves  [depth+2]string
	minMoves2 [depth+2]string
	minMoves3 [depth+2]string
)

func main() {
	var board Board

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
	board.printBoard()
	fmt.Println("----------------------------")
	start := time.Now()

	// ルートを探索
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			bestMoves[i][j].Point = 30
			calcMoves(board, j, i)
		}
	}

	bestMove := BestMove{
		Point: 30,
	}
	// 30の候補の中から最善手を計算
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			if bestMoves[i][j].Point < bestMove.Point {
				bestMove.Point = bestMoves[i][j].Point
				bestMove.Moves = bestMoves[i][j].Moves
			}
		}
	}

	end := time.Now()
	result := fmt.Sprintf("解析時間: %f秒\n", (end.Sub(start)).Seconds())
	fmt.Println(bestMove.Moves)
	printMoves(bestMove.Moves)
	fmt.Println("消える数: " + strconv.Itoa(30 - bestMove.Point))
	fmt.Println(result)
}

// BestMoves 最善手のリスト
type BestMoves [5][6]BestMove

// BestMove 最善手の情報を格納
type BestMove struct {
	Moves [depth+2]string
	Point int
}

// Board 盤面
type Board [5][6]int

// calcMoves どの移動方法が最適か計算する。開始位置は指定。
func calcMoves(board Board, x int, y int) {
	var moves [depth+2]string
	moves[0] = strconv.Itoa(x+1)
	moves[1] = strconv.Itoa(y+1)
	move(1, board, x, y, moves)
}

// move ドロップを移動させる。手数、盤面、x座標、y座標
func move(n int, board Board, x int, y int, moves [depth+2]string) {
	point := 0
	firstXPlus, _ := strconv.Atoi(moves[0])
	firstYPlus, _ := strconv.Atoi(moves[1])
	firstX := firstXPlus - 1
	firstY := firstYPlus - 1
	// 指定の手数-1以上になった場合、ポイントを算出して最高得点だったらmovesを記録。
	if n > depth - 1 {
		point = calcPoint(board)
		if point < bestMoves[firstY][firstX].Point {
			bestMoves[firstY][firstX].Point = point
			bestMoves[firstY][firstX].Moves = moves
		}
		if n > depth {
			return
		}
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
	}

	drop := board[y][x]
	// 高速化のため。
	nPlus := n+1
	xPlus := x+1
	xMinus := x-1
	yPlus := y+1
	yMinus := y-1

	// Right
	if moves[n] != "←" && x != 5 {
		dropR := board[y][xPlus]
		board[y][x] = dropR
		board[y][xPlus] = drop
		moves[nPlus] = "→"
		move(nPlus, board, xPlus, y, moves)
		board[y][x] = drop
		board[y][xPlus] = dropR
		moves[nPlus] = ""
	}

	// Down
	if moves[n] != "↑" && y != 4 {
		dropD := board[yPlus][x]
		board[y][x] = dropD
		board[yPlus][x] = drop
		moves[nPlus] = "↓"
		move(nPlus, board, x, yPlus, moves)
		board[y][x] = drop
		board[yPlus][x] = dropD
		moves[nPlus] = ""
	}

	// Left
	if moves[n] != "→" && x != 0 {
		dropL := board[y][xMinus]
		board[y][x] = dropL
		board[y][xMinus] = drop
		moves[nPlus] = "←"
		move(nPlus, board, xMinus, y, moves)
		board[y][x] = drop
		board[y][xMinus] = dropL
		moves[nPlus] = ""
	}

	// Up
	if moves[n] != "↓" && y != 0 {
		dropU := board[yMinus][x]
		board[y][x] = dropU
		board[yMinus][x] = drop
		moves[nPlus] = "↑"
		move(nPlus, board, x, yMinus, moves)
		board[y][x] = drop
		board[yMinus][x] = dropU
		moves[nPlus] = ""
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
	var fallenBoard Board

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
	var isWork bool

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

// printBoard 配列をきれいに出力。テスト用かも。
func (board Board) printBoard() {
	for _, b := range board {
		fmt.Println(b)
	}
}

// printMoves ルートをきれいに表示
func printMoves(moves [depth+2]string) {
	root := [][]string{
		{"○", " ", "○", " ", "○", " ", "○", " ", "○", " ", "○"},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{"○", " ", "○", " ", "○", " ", "○", " ", "○", " ", "○"},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{"○", " ", "○", " ", "○", " ", "○", " ", "○", " ", "○"},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{"○", " ", "○", " ", "○", " ", "○", " ", "○", " ", "○"},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{"○", " ", "○", " ", "○", " ", "○", " ", "○", " ", "○"},
	}

	startX, _ := strconv.Atoi(moves[0])
	startY, _ := strconv.Atoi(moves[1])
	x := startX * 2 - 2
	y := startY * 2 - 2
	root[y][x] = "S"

	for i := 2; i < len(moves); i++ {
		switch moves[i] {
			case "→":
				arrow := root[y][x+1]
				if arrow == " " || arrow == "→" {
					root[y][x+1] = "→"
				} else if arrow == "←" || arrow == "⇄" {
					root[y][x+1] = "⇄"
				}
				x += 2
			case "↓":
				arrow := root[y+1][x]
				if arrow == " " || arrow == "↓" {
					root[y+1][x] = "↓"
				} else if arrow == "↑" || arrow == "⇅" {
					root[y+1][x] = "⇅"
				}
				y += 2
			case "←":
				arrow := root[y][x-1]
				if arrow == " " || arrow == "←" {
					root[y][x-1] = "←"
				} else if arrow == "→" || arrow == "⇄" {
					root[y][x-1] = "⇄"
				}
				x -= 2
			case "↑":
				arrow := root[y-1][x]
				if arrow == " " || arrow == "↑" {
					root[y-1][x] = "↑"
				} else if arrow == "↓" || arrow == "⇅" {
					root[y-1][x] = "⇅"
				}
				y -= 2
		}
	}

	root[y][x] = "G"

	for _, r := range root {
		fmt.Println(r)
	}
}
