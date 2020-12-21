package main

import (
	"testing"
)

func TestManMoves(t *testing.T) {
	board := [8][8]string{}
	board[3][3] = "BM"

	moves := getMoves(black, down, &board)

	if len(moves) != 2 {
		t.Errorf("len(moves) is %d, should be 2", len(moves))
	}

	board = [8][8]string{}
	board[3][0] = "BM"

	moves = getMoves(black, down, &board)
	if len(moves) != 1 {
		t.Errorf("len(moves) is %d, should be 1", len(moves))
	}

	board = getBoard()
	moves = getMoves(white, up, &board)

	if len(moves) != 7 {
		t.Errorf("len(moves) is %d, should be 7", len(moves))
	}
}

func TestManJumps(t *testing.T) {
	board := [8][8]string{}
	board[3][3] = "BM"
	board[4][4] = "WM"

	jumps := getManJumps(white, up, Position{4, 4}, &board)

	if len(jumps) != 1 {
		t.Errorf("length(jumps) == %d, should be 1", len(jumps))
	}

	jump := jumps[0]

	target := Move{Position{4, 4}, Position{2, 2}}
	if jump != target {
		t.Error("Expected", target, "got", jump)
	}

	board = [8][8]string{}
	board[0][0] = "BM"
	board[1][1] = "WM"

	jumps = getManJumps(white, up, Position{4, 4}, &board)

	if len(jumps) != 0 {
		t.Errorf("length(jumps) == %d, should be 0", len(jumps))
	}

	board = [8][8]string{}
	board[6][6] = "BM"
	board[7][7] = "WM"

	jumps = getManJumps(black, down, Position{4, 4}, &board)

	if len(jumps) != 0 {
		t.Errorf("length(jumps) == %d, should be 0", len(jumps))
	}
}

func TestKingMoves(t *testing.T) {
	board := [8][8]string{}

	board[3][3] = "WK"
	board[6][6] = "BM"

	moves, jumps := getKingMoves(white, Position{3, 3}, &board)

	if len(moves)+len(jumps) != 12 {
		t.Errorf("len(moves) len(jumps) == %d, should be 12", len(moves)+len(jumps))
	}

}

func TestEvaluateBoard(t *testing.T) {
	board := getBoard()

	score := evaluateBoard(white, &board)
	target := 0
	if score != target {
		t.Errorf("Score is %d, but should be %d", score, target)
	}
}

func TestMakeMove(t *testing.T) {
	board := [8][8]string{}
	board[2][2] = "BK"
	board[6][6] = "WM"

	move := Move{Position{2, 2}, Position{7, 7}}

	newBoard := makeMove(move, board)

	if newBoard[2][2] != "" {
		t.Errorf("Origin contains %s, but should be empty", newBoard[2][2])
	}

	if newBoard[6][6] != "" {
		t.Errorf("Field before destination contains %s, but should be empty", newBoard[6][6])
	}

	if newBoard[7][7] != "BK" {
		t.Errorf("Destination contains %s, but should contain 'BK'", newBoard[7][7])
	}

}
