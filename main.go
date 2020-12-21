package main

import "fmt"

type playerColor int

const (
	white playerColor = iota
	black
)

type Position struct {
	i int
	j int
}

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

func main() {
	board := getBoard()
	printBoard(&board)
	fmt.Println(getPositions(white, &board))
}
