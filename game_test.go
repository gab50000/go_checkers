package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestManMoves(t *testing.T) {
	board := [8][8]string{}
	board[3][3] = "BM"
	gs := GameState{board, black, down}

	moves := getMoves(&gs)

	if len(moves) != 2 {
		t.Errorf("len(moves) is %d, should be 2", len(moves))
	}

	board = [8][8]string{}
	board[3][0] = "BM"
	gs = GameState{board, black, down}

	moves = getMoves(&gs)
	if len(moves) != 1 {
		t.Errorf("len(moves) is %d, should be 1", len(moves))
	}

	gs = createInitialState()
	moves = getMoves(&gs)

	if len(moves) != 7 {
		t.Errorf("len(moves) is %d, should be 7", len(moves))
	}
}

func TestManJumps(t *testing.T) {
	board := [8][8]string{}
	board[3][3] = "BM"
	board[4][4] = "WM"

	gs := GameState{board, white, up}
	fmt.Println(gs)

	jumps := getManJumps(&gs, Position{4, 4})
	fmt.Println(jumps)

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

	gs = GameState{board, white, up}

	jumps = getManJumps(&gs, Position{4, 4})

	if len(jumps) != 0 {
		t.Errorf("length(jumps) == %d, should be 0", len(jumps))
	}

	board = [8][8]string{}
	board[6][6] = "BM"
	board[7][7] = "WM"

	gs = GameState{board, black, down}

	jumps = getManJumps(&gs, Position{4, 4})

	if len(jumps) != 0 {
		t.Errorf("length(jumps) == %d, should be 0", len(jumps))
	}
}

func TestKingMoves(t *testing.T) {
	board := [8][8]string{}

	board[3][3] = "WK"
	board[6][6] = "BM"

	gs := GameState{board, white, up}

	moves, jumps := getKingMoves(&gs, Position{3, 3})

	if len(moves)+len(jumps) != 12 {
		t.Errorf("len(moves) len(jumps) == %d, should be 12", len(moves)+len(jumps))
	}

}

func TestEvaluateBoard(t *testing.T) {
	gs := createInitialState()

	score := evaluateCurrentBoard(&gs)
	target := 0.0
	if score != target {
		t.Errorf("Score is %f, but should be %f", score, target)
	}
}

func TestMakeMove(t *testing.T) {
	board := [8][8]string{}
	board[2][2] = "BK"
	board[6][6] = "WM"

	gs := GameState{board, white, up}

	move := Move{Position{2, 2}, Position{7, 7}}

	newState := gs.makeMove(move)
	newBoard := newState.board

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

func TestMin(t *testing.T) {
	numbers := []float64{3, 4, 1, -123, 44}

	m, _ := min(numbers)
	target := -123.0

	if m != target {
		t.Errorf("m == %f, should be %f", m, target)
	}

	numbers = []float64{}
	_, err := min(numbers)
	if err == nil {
		t.Errorf("An empty slice should throw an error")
	}
}

func TestGetMoves(t *testing.T) {
	board := [boardSize][boardSize]string{}

	board[3][3] = "BM"
	board[4][4] = "WM"

	gs := GameState{board, black, down}

	move := getMoves(&gs)
	target := []Move{{Position{3, 3}, Position{5, 5}}}

	if !cmp.Equal(move, target) {
		t.Error("move ==", move, ", but should be", target)
	}

}

func TestChooseBestMove(t *testing.T) {
	board := [boardSize][boardSize]string{}

	board[3][3] = "BM"
	board[4][4] = "WM"

	gs := GameState{board, black, down}

	move := chooseBestMove(&gs, 5, true)
	target := Move{Position{3, 3}, Position{5, 5}}

	if move != target {
		t.Error("move ==", move, ", but should be", target)
	}

}

func TestParseMove(t *testing.T) {
	inputs := []string{
		"a1b2",
		"c3  d4",
		"a5    b6",
	}
	targets := []Move{
		{Position{0, 0}, Position{1, 1}},
		{Position{2, 2}, Position{3, 3}},
		{Position{4, 0}, Position{5, 1}},
	}

	for i, inp := range inputs {
		trg := targets[i]

		result, _ := parseMove(inp)

		if result != trg {
			t.Error("Result is", result, "but should be", trg)
		}
	}
}

func TestTreeSearchVsAlphaBetaPruning(t *testing.T) {
	searchDepth := 7
	numberOfMoves := 10

	play := func(alphaBetaPruning bool) GameState {
		state := createInitialState()
		color := white
		dir := up

		for i := 0; i < numberOfMoves; i++ {
			move := chooseBestMove(&state, searchDepth, alphaBetaPruning)
			state = state.makeMove(move)
			color = oppositeColor(color)
			dir = oppositeDirection(dir)
		}
		return state
	}

	board1 := play(true)
	board2 := play(false)

	fmt.Println(board1)
	fmt.Println(board2)

	if board1 != board2 {
		t.Error("Alpha-Beta pruning yields different results!")
	}
}

func BenchmarkChooseBestMove(b *testing.B) {
	gs := createInitialState()

	for i := 0; i < b.N; i++ {
		chooseBestMove(&gs, 7, true)
	}

}
