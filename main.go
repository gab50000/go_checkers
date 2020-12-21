package main

import (
	"fmt"
)

type playerColor int
type direction int

const (
	white playerColor = iota
	black
)

const (
	up direction = iota
	down
)

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

func getBoard() [8][8]string {
	board := [8][8]string{}
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

func printBoard(board *[8][8]string) {
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

func getPositions(color playerColor, board *[8][8]string) []Position {
	positions := make([]Position, 0, 12)
	var prefix string
	if color == black {
		prefix = "B"
	} else {
		prefix = "W"
	}

	for i, row := range board {
		for j, elem := range row {
			if elem != "" && elem[:1] == prefix {
				positions = append(positions, Position{i, j})
			}
		}
	}
	return positions
}

func getMoves(color playerColor, dir direction, board *[8][8]string) []Move {
	moves := make([]Move, 0)

	positions := getPositions(color, board)

	for _, pos := range positions {
		moves = append(moves, getManMoves(color, dir, pos, board)...)
	}
	return moves
}

func getManMoves(color playerColor, dir direction, pos Position, board *[8][8]string) []Move {
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

func getManJumps(
	color playerColor,
	dir direction,
	pos Position,
	board *[8][8]string) []Move {

	moves := make([]Move, 0)
	i, j := pos.i, pos.j

	enemyColor := oppositeColor(color)
	var enemyPrefix string
	switch enemyColor {
	case black:
		enemyPrefix = "B"
	case white:
		enemyPrefix = "W"
	}

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
		if jEnemy == 0 || jEnemy == 7 {
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

// func getKingMoves(color playerColor, dir direction, pos Position) {}

// func getKingJumps(color playerColor, dir direction, pos Position) {}

func getJumps(color playerColor, dir direction, board *[8][8]string) []Move {
	moves := make([]Move, 0)
	// positions := getPositions(color, board)

	return moves
}

func main() {
	board := getBoard()
	printBoard(&board)
	fmt.Println(getPositions(white, &board))
}
