package main

import "fmt"

func getBoard() [8][8]string {
	board := [8][8]string{}
	for i := 0; i < 4; i++ {
		board[0][2*i] = "BM"
		board[1][2*i+1] = "BM"
		board[6][2*i] = "WM"
		board[7][2*i+1] = "WM"
	}

	return board
}

func printBoard(board [8][8]string) {
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

func main() {
	fmt.Println("hello")
	board := getBoard()
	printBoard(board)
}
