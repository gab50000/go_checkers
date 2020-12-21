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

}
