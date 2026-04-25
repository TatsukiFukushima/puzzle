package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const depth int = 26

var intDropMap map[int]string = map[int]string{
	1: "\033[1m\033[31m●\033[0m",
	2: "\033[1m\033[34m●\033[0m",
	3: "\033[1m\033[32m●\033[0m",
	4: "\033[1m\033[33m●\033[0m",
	5: "\033[1m\033[35m●\033[0m",
	6: "\033[1m\033[38;2;255;105;180m■\033[0m",
}

var rate int = 0
var bestMoves BestMoves

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
	fmt.Printf("探索手数: %d\n", depth)
	start := time.Now()

	// ルートを探索
	var wg sync.WaitGroup
	wg.Add(30)
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			bestMoves[i][j].Point = 30
			go calcMoves(board, j, i, &wg)
		}
	}
	wg.Wait()

	bestMove := BestMove{Point: 30}
	bestMove2 := BestMove{Point: 30}
	bestMove3 := BestMove{Point: 30}

	// 30の候補の中から最善手を計算
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			if bestMoves[i][j].Point < bestMove.Point {
				bestMove.Point = bestMoves[i][j].Point
				bestMove.Moves = bestMoves[i][j].Moves
			} else if bestMoves[i][j].Point < bestMove2.Point {
				bestMove2.Point = bestMoves[i][j].Point
				bestMove2.Moves = bestMoves[i][j].Moves
			} else if bestMoves[i][j].Point < bestMove3.Point {
				bestMove3.Point = bestMoves[i][j].Point
				bestMove3.Moves = bestMoves[i][j].Moves
			}
		}
	}

	end := time.Now()
	result := fmt.Sprintf("解析時間: %f秒\n", (end.Sub(start)).Seconds())
	bestArrowMove := intToArrow(bestMove.Moves)
	bestArrowMove2 := intToArrow(bestMove2.Moves)
	bestArrowMove3 := intToArrow(bestMove3.Moves)

	fmt.Printf("\r")
	fmt.Println("---------------------")
	fmt.Print("候補１: ")
	fmt.Println(bestArrowMove)
	fmt.Println("消える数: " + strconv.Itoa(30-bestMove.Point))
	printMoves(bestArrowMove)
	fmt.Println("")
	fmt.Print("候補２: ")
	fmt.Println(bestArrowMove2)
	fmt.Println("消える数: " + strconv.Itoa(30-bestMove2.Point))
	printMoves(bestArrowMove2)
	fmt.Println("")
	fmt.Print("候補３: ")
	fmt.Println(bestArrowMove3)
	fmt.Println("消える数: " + strconv.Itoa(30-bestMove3.Point))
	printMoves(bestArrowMove3)
	fmt.Println(result)
}

// BestMoves 最善手のリスト
type BestMoves [5][6]BestMove

// BestMove 最善手の情報を格納
type BestMove struct {
	Moves [depth + 2]int
	Point int
}

// Board 盤面
type Board [5][6]int

// calcMoves どの移動方法が最適か計算する。開始位置は指定。
func calcMoves(board Board, x int, y int, wg *sync.WaitGroup) {
	var moves [depth + 2]int
	moves[0] = x
	moves[1] = y
	currentDrop := board[y][x]
	move(1, &board, x, y, currentDrop, moves)
	rate++
	log := strconv.Itoa(rate) + " / 30"
	fmt.Printf("\r%s", log)
	wg.Done()
}

// move ドロップを移動させる。手数、盤面、x座標、y座標
func move(n int, board *Board, x int, y int, currentDrop int, moves [depth + 2]int) {
	point := 0
	// 指定の手数になった場合、ポイントを算出して最高得点だったらmovesを記録。
	if n > depth {
		point = calcPoint(*board)
		firstX := moves[0]
		firstY := moves[1]
		if point < bestMoves[firstY][firstX].Point {
			bestMoves[firstY][firstX].Point = point
			bestMoves[firstY][firstX].Moves = moves
		}
		return
	} else if n == 9 {
		// 枝刈り。8回移動して27ポイントより大きいルートは探索しない。
		point = calcPoint(*board)
		if point > 27 {
			return
		}
	} else if n == 13 {
		// 枝刈り2。12回移動して24ポイントより大きいルートは探索しない。
		point = calcPoint(*board)
		if point > 24 {
			return
		}
	} else if n == 17 {
		// 枝刈り3。16回移動して24ポイントより大きいルートは探索しない。
		point = calcPoint(*board)
		if point > 24 {
			return
		}
	} else if n == 21 {
		// 枝刈り4。20回移動して24ポイントより大きいルートは探索しない。
		point = calcPoint(*board)
		if point > 24 {
			return
		}
	}

	// 高速化のため。
	nPlus := n + 1
	xPlus := x + 1
	xMinus := x - 1
	yPlus := y + 1
	yMinus := y - 1

	// Right 1
	if (moves[n] != 3 || n == 1) && x != 5 {
		board[y][x] = board[y][xPlus]
		board[y][xPlus] = currentDrop
		moves[nPlus] = 1
		move(nPlus, board, xPlus, y, currentDrop, moves)
		board[y][xPlus] = board[y][x]
		board[y][x] = currentDrop
	}

	// Down 2
	if (moves[n] != 4 || n == 1) && y != 4 {
		board[y][x] = board[yPlus][x]
		board[yPlus][x] = currentDrop
		moves[nPlus] = 2
		move(nPlus, board, x, yPlus, currentDrop, moves)
		board[yPlus][x] = board[y][x]
		board[y][x] = currentDrop
	}

	// Left 3
	if (moves[n] != 1 || n == 1) && x != 0 {
		board[y][x] = board[y][xMinus]
		board[y][xMinus] = currentDrop
		moves[nPlus] = 3
		move(nPlus, board, xMinus, y, currentDrop, moves)
		board[y][xMinus] = board[y][x]
		board[y][x] = currentDrop
	}

	// Up 4
	if (moves[n] != 2 || n == 1) && y != 0 {
		board[y][x] = board[yMinus][x]
		board[yMinus][x] = currentDrop
		moves[nPlus] = 4
		move(nPlus, board, x, yMinus, currentDrop, moves)
		board[yMinus][x] = board[y][x]
		board[y][x] = currentDrop
	}
	moves[nPlus] = 0
}

// delete ドロップを消す
func delete(board Board) (Board, bool) {
	isWork := false
	deletedBoard := board

	// 横方向の削除
	for i := 0; i < 5; i++ {
		for j := 1; j < 5; j++ {
			drop := board[i][j]
			if drop == 0 {
				continue
			}

			jPlus := j + 1
			jMinus := j - 1
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
			if drop == 0 {
				continue
			}

			iPlus := i + 1
			iMinus := i - 1
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
func fall(board *Board) {
	// 列ごとにループを回しているということを明示するためにjとしている。
	for j := 0; j < 6; j++ {
		// ドロップが落ちる先のy座標、的なニュアンス
		nextY := 4
		for i := 4; i >= 0; i-- {
			drop := board[i][j]
			if drop != 0 {
				if i != nextY {
					board[nextY][j] = drop
					board[i][j] = 0
				}
				nextY--
			}
		}
	}
}

// calcPoint どれだけ残るかを計算
func calcPoint(board Board) int {
	var isWork bool

	// 消えるドロップがなくなるまでloop
	point := 0
	for {
		board, isWork = delete(board)
		if isWork {
			fall(&board)
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

// printBoard 配列をきれいに出力
func (board Board) printBoard() {
	for _, b := range board {
		for _, drop := range b {
			if drop == 0 {
				fmt.Print("  ")
			} else {
				fmt.Print(intDropMap[drop] + " ")
			}
		}
		fmt.Println()
	}
}

// printMoves ルートをきれいに表示
func printMoves(moves [depth + 2]string) {
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
	x := startX*2 - 2
	y := startY*2 - 2
	root[y][x] = "S"

	for i := 2; i < len(moves); i++ {
		switch moves[i] {
		case "→":
			switch root[y][x+1] {
			case " ":
				root[y][x+1] = "→"
			case "→":
				root[y][x+1] = "⇉"
			default:
				root[y][x+1] = "⇄"
			}
			x += 2
		case "↓":
			switch root[y+1][x] {
			case " ":
				root[y+1][x] = "↓"
			case "↓":
				root[y+1][x] = "⇊"
			default:
				root[y+1][x] = "⇅"
			}
			y += 2
		case "←":
			switch root[y][x-1] {
			case " ":
				root[y][x-1] = "←"
			case "←":
				root[y][x-1] = "⇇"
			default:
				root[y][x-1] = "⇄"
			}
			x -= 2
		case "↑":
			switch root[y-1][x] {
			case " ":
				root[y-1][x] = "↑"
			case "↑":
				root[y-1][x] = "⇈"
			default:
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

// intToArrow 数字のmoveを矢印に変換
func intToArrow(moves [depth + 2]int) [depth + 2]string {
	var arrowMoves [depth + 2]string
	arrowMoves[0] = strconv.Itoa(moves[0] + 1)
	arrowMoves[1] = strconv.Itoa(moves[1] + 1)

	for i := 2; i < depth+2; i++ {
		switch moves[i] {
		case 1:
			arrowMoves[i] = "→"
		case 2:
			arrowMoves[i] = "↓"
		case 3:
			arrowMoves[i] = "←"
		case 4:
			arrowMoves[i] = "↑"
		default:
			arrowMoves[i] = ""
		}
	}

	return arrowMoves
}
