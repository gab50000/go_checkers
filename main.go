package main

import (
	"fmt"
	"strings"
)

type playerColor int
type direction int
type tokenType int

const (
	white playerColor = iota
	black
)

const (
	up direction = iota
	down
)

const (
	man tokenType = iota
	king
)

const boardSize = 8

// Position of a token
type Position struct {
	i int
	j int
}

// Move holding start and end position
type Move struct {
	from Position
	to   Position
}

func getBoard() [boardSize][boardSize]string {
	board := [boardSize][boardSize]string{}
	for i := 0; i < 4; i++ {
		board[0][2*i] = "BM"
		board[1][2*i+1] = "BM"
		board[2][2*i] = "BM"
		board[5][2*i+1] = "WM"
		board[6][2*i] = "WM"
		board[7][2*i+1] = "WM"
	}

	return board
}

func printBoard(board *[boardSize][boardSize]string) {
	const letters = "   A   B   C   D   E   F   G   H"
	fmt.Println(letters)
	for i, row := range board {
		fmt.Print(i + 1)
		for _, elem := range row {
			fmt.Print("|")
			switch elem {
			case "":
				fmt.Print("   ")
			case "BM":
				fmt.Print(" o ")
			case "BK":
				fmt.Print(" ♔ ")
			case "WM":
				fmt.Print(" ● ")
			case "WK":
				fmt.Print(" ♚ ")
			}
		}
		fmt.Printf("|%d\n", i+1)
	}
	fmt.Println(letters)
}

func getPositions(color playerColor, tok tokenType, board *[boardSize][boardSize]string) []Position {
	positions := make([]Position, 0, 12)
	var prefix string
	if color == black {
		prefix = "B"
	} else {
		prefix = "W"
	}
	var postfix string
	if tok == man {
		postfix = "M"
	} else {
		postfix = "K"
	}
	target := prefix + postfix

	for i, row := range board {
		for j, elem := range row {
			if elem == target {
				positions = append(positions, Position{i, j})
			}
		}
	}
	return positions
}

func getMoves(color playerColor, dir direction, board *[boardSize][boardSize]string) (moves []Move) {

	manPositions := getPositions(color, man, board)
	kingPositions := getPositions(color, king, board)

	for _, pos := range manPositions {
		moves = append(moves, getManJumps(color, dir, pos, board)...)
	}

	for _, pos := range kingPositions {
		_, kingJumps := getKingMoves(color, pos, board)
		moves = append(moves, kingJumps...)
	}

	if len(moves) > 0 {
		return moves
	}

	for _, pos := range manPositions {
		moves = append(moves, getManMoves(color, dir, pos, board)...)
	}

	for _, pos := range kingPositions {
		kingMoves, _ := getKingMoves(color, pos, board)
		moves = append(moves, kingMoves...)

	}
	return moves
}

func getManMoves(color playerColor, dir direction, pos Position, board *[boardSize][boardSize]string) []Move {
	moves := make([]Move, 0, 2)
	i, j := pos.i, pos.j
	var ii int
	switch {
	case dir == up && i > 0:
		ii = i - 1
	case dir == down && i < 7:
		ii = i + 1
	default:
		return []Move{}
	}
	for dj := -1; dj <= 1; dj += 2 {
		jj := j + dj
		if jj < 0 || jj > 7 {
			continue
		}

		if board[ii][jj] == "" {
			moves = append(moves, Move{pos, Position{ii, jj}})
		}
	}

	return moves
}

func oppositeColor(color playerColor) playerColor {
	var oppColor playerColor
	switch color {
	case white:
		oppColor = black
	case black:
		oppColor = white
	}
	return oppColor
}

func colorPrefix(color playerColor) (prefix string) {
	switch color {
	case black:
		prefix = "B"
	case white:
		prefix = "W"
	}
	return prefix
}

func withinBounds(indices ...int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= boardSize {
			return false
		}
	}
	return true
}

func getManJumps(
	color playerColor,
	dir direction,
	pos Position,
	board *[boardSize][boardSize]string) []Move {

	moves := make([]Move, 0)
	i, j := pos.i, pos.j

	enemyColor := oppositeColor(color)
	enemyPrefix := colorPrefix(enemyColor)

	var iEnemy, iDestination, jEnemy, jDestination int
	switch {
	case dir == up && i > 1:
		{
			iEnemy = i - 1
			iDestination = i - 2
		}
	case dir == down && i < 6:
		{
			iEnemy = i + 1
			iDestination = i + 2
		}
	default:
		return []Move{}
	}

	for dj := -1; dj <= 1; dj += 2 {
		jEnemy = j + dj
		if jEnemy == 0 || jEnemy == 7 || !withinBounds(jEnemy) {
			continue
		}
		jDestination = j + 2*dj

		enemyElem := board[iEnemy][jEnemy]
		if enemyElem != "" &&
			enemyElem[:1] == enemyPrefix &&
			board[iDestination][jDestination] == "" {
			moves = append(moves, Move{pos, Position{iDestination, jDestination}})
		}
	}
	return moves
}

func getKingMoves(color playerColor, pos Position, board *[boardSize][boardSize]string) (moves []Move, jumps []Move) {
	i, j := pos.i, pos.j
	enemyColor := oppositeColor(color)
	enemyPrefix := colorPrefix(enemyColor)

	for di := -1; di < 2; di += 2 {
		for dj := -1; dj < 2; dj += 2 {
			ii, jj := i+di, j+dj
			for withinBounds(ii, jj) && board[ii][jj] == "" {
				moves = append(moves, Move{Position{i, j}, Position{ii, jj}})
				ii, jj = ii+di, jj+dj
			}

			iDest, jDest := ii+di, jj+dj

			if !withinBounds(iDest, jDest) {
				continue
			}

			if board[ii][jj][:1] == enemyPrefix && board[iDest][jDest] == "" {
				jumps = append(jumps, Move{Position{i, j}, Position{iDest, jDest}})
			}

		}

	}
	return moves, jumps
}

func countTokens(board *[boardSize][boardSize]string) (counter map[playerColor]int) {
	counter = make(map[playerColor]int)
	for _, row := range board {
		for _, elem := range row {
			if strings.HasPrefix(elem, "B") {
				counter[black]++
			} else if strings.HasPrefix(elem, "W") {
				counter[white]++
			}
		}
	}
	return counter
}

func evaluateBoard(color playerColor, board *[boardSize][boardSize]string) int {
	tokenCounter := countTokens(board)
	enemyColor := oppositeColor(color)

	enemyCount := tokenCounter[enemyColor]
	selfCount := tokenCounter[color]

	if enemyCount == 0 {
		return 10
	} else if selfCount == 0 {
		return -10
	}
	return selfCount - enemyCount
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func makeMove(move Move, board [boardSize][boardSize]string) [boardSize][boardSize]string {
	token := board[move.from.i][move.from.j]
	board[move.to.i][move.to.j] = token
	board[move.from.i][move.from.j] = ""
	dI := move.from.i - move.to.i
	dI /= abs(dI)
	dJ := move.from.j - move.to.j
	dJ /= abs(dJ)
	board[move.to.i+dI][move.to.j+dJ] = ""
	return board
}

func main() {
	board := getBoard()
	printBoard(&board)
	fmt.Println(getPositions(white, man, &board))
}
