package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
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
	I int
	J int
}

func (p Position) String() string {
	return fmt.Sprintf("Position{%v, %v}", p.I, p.J)
}

// Move holding start and end position
type Move struct {
	From Position
	To   Position
}

func (m Move) String() string {
	return fmt.Sprintf("Move{%v, %v}", m.From, m.To)
}

// GameState holds the board state
type GameState struct {
	board            [8][8]string
	currentPlayer    playerColor
	currentDirection direction
}

func (gs GameState) String() string {
	return boardToString(&gs.board)
}

func createInitialState() GameState {
	return GameState{
		board:            getBoard(),
		currentPlayer:    white,
		currentDirection: up,
	}
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

func boardToString(board *[boardSize][boardSize]string) string {
	const letters = "   A   B   C   D   E   F   G   H"
	var s string
	s += fmt.Sprintf("%v\n", letters)
	for i, row := range board {
		s += fmt.Sprint(i + 1)
		for _, elem := range row {
			s += fmt.Sprint("|")
			switch elem {
			case "":
				s += fmt.Sprint("   ")
			case "BM":
				s += fmt.Sprint(" o ")
			case "BK":
				s += fmt.Sprint(" ♔ ")
			case "WM":
				s += fmt.Sprint(" ● ")
			case "WK":
				s += fmt.Sprint(" ♚ ")
			}
		}
		s += fmt.Sprintf("|%d\n", i+1)
	}
	s += fmt.Sprintf("%v\n", letters)
	return s
}

func printBoard(board *[boardSize][boardSize]string) {
	fmt.Print(boardToString(board))
}

func getPositions(gs *GameState, tok tokenType) []Position {
	positions := make([]Position, 0, 12)
	var prefix string
	if gs.currentPlayer == black {
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

	for i, row := range gs.board {
		for j, elem := range row {
			if elem == target {
				positions = append(positions, Position{i, j})
			}
		}
	}
	return positions
}

func getMoves(gs *GameState) (moves []Move) {

	manPositions := getPositions(gs, man)
	kingPositions := getPositions(gs, king)

	for _, pos := range manPositions {
		moves = append(moves, getManJumps(gs, pos)...)
	}

	for _, pos := range kingPositions {
		_, kingJumps := getKingMoves(gs, pos)
		moves = append(moves, kingJumps...)
	}

	if len(moves) > 0 {
		return moves
	}

	for _, pos := range manPositions {
		moves = append(moves, getManMoves(gs, pos)...)
	}

	for _, pos := range kingPositions {
		kingMoves, _ := getKingMoves(gs, pos)
		moves = append(moves, kingMoves...)

	}
	return moves
}

func getManMoves(gs *GameState, pos Position) []Move {
	dir := gs.currentDirection
	moves := make([]Move, 0, 2)
	i, j := pos.I, pos.J
	var ii int
	switch {
	case dir == up && i > 0:
		ii = i - 1
	case dir == down && i < 7:
		ii = i + 1
	default:
		return []Move{}
	}
	for _, dj := range []int{-1, 1} {
		jj := j + dj
		if jj < 0 || jj > 7 {
			continue
		}

		if gs.board[ii][jj] == "" {
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

func oppositeDirection(dir direction) direction {
	var newDir direction
	switch dir {
	case up:
		newDir = down
	case down:
		newDir = up
	}
	return newDir
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
	gs *GameState,
	pos Position,
) []Move {
	board := gs.board
	color := gs.currentPlayer
	dir := gs.currentDirection

	moves := make([]Move, 0)
	i, j := pos.I, pos.J

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

	for dj := range []int{-1, 1} {
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

func getKingMoves(gs *GameState, pos Position) (moves []Move, jumps []Move) {
	i, j := pos.I, pos.J
	color := gs.currentPlayer
	board := gs.board
	enemyColor := oppositeColor(color)
	enemyPrefix := colorPrefix(enemyColor)

	for _, di := range []int{-1, 1} {
		for _, dj := range []int{-1, 1} {
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

func countTokens(board *[boardSize][boardSize]string) (nWhite, nBlack int) {
	for _, row := range board {
		for _, elem := range row {
			if strings.HasPrefix(elem, "B") {
				nBlack++
			} else if strings.HasPrefix(elem, "W") {
				nWhite++
			}
		}
	}
	return nWhite, nBlack
}

func evaluateCurrentBoard(gs *GameState) float64 {
	nWhite, nBlack := countTokens(&gs.board)

	var enemyCount, selfCount int
	switch gs.currentPlayer {
	case white:
		{
			enemyCount = nBlack
			selfCount = nWhite
		}
	case black:
		{
			enemyCount = nWhite
			selfCount = nBlack
		}
	}

	if enemyCount == 0 {
		return 10
	} else if selfCount == 0 {
		return -10
	}
	return float64(selfCount - enemyCount)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func toKing(token string) string {
	return token[:1] + "K"
}

func makeMove(gs GameState, move Move) GameState {
	noMove := Move{Position{0, 0}, Position{0, 0}}
	if move == noMove {
		return GameState{
			board:            gs.board,
			currentPlayer:    oppositeColor(gs.currentPlayer),
			currentDirection: oppositeDirection(gs.currentDirection),
		}
	}

	board := gs.board
	dir := gs.currentDirection

	destI, destJ := move.To.I, move.To.J
	origI, origJ := move.From.I, move.From.J
	token := board[origI][origJ]
	if (destI == 0 && dir == up) || (destI == boardSize-1 && dir == down) {
		token = toKing(token)
	}
	board[destI][destJ] = token
	board[origI][origJ] = ""
	dI := origI - destI
	dI /= abs(dI)
	dJ := origJ - destJ
	dJ /= abs(dJ)
	board[destI+dI][destJ+dJ] = ""
	return GameState{
		board:            board,
		currentPlayer:    oppositeColor(gs.currentPlayer),
		currentDirection: oppositeDirection(gs.currentDirection),
	}
}

func min(numbers []float64) (m float64, e error) {
	if len(numbers) == 0 {
		return 0, errors.New("slice is empty")
	}
	for i, num := range numbers {
		if i == 0 {
			m = num
		} else if num < m {
			m = num
		}
	}
	return m, nil
}

func evaluateBoard(
	gs *GameState,
	depthRemaining int,
	bestScoreBlack float64,
	bestScoreWhite float64,
	alphaBetaPruning bool,
) (score float64) {

	var bestOpponentScore *float64
	var bestSelfScore *float64
	switch gs.currentPlayer {
	case black:
		bestOpponentScore = &bestScoreWhite
		bestSelfScore = &bestScoreBlack
	case white:
		bestOpponentScore = &bestScoreBlack
		bestSelfScore = &bestScoreWhite

	}

	if depthRemaining == 0 {
		return evaluateCurrentBoard(gs)
	}
	scores := make([]float64, 0)

	moves := getMoves(gs)

	if len(moves) == 0 {
		return evaluateCurrentBoard(gs)
	}

	for _, move := range moves {
		newState := makeMove(*gs, move)
		newOppScore := evaluateBoard(
			&newState,
			depthRemaining-1,
			bestScoreBlack,
			bestScoreWhite,
			alphaBetaPruning,
		)
		newScore := -newOppScore

		if newScore > *bestSelfScore {
			*bestSelfScore = newScore
		}

		if *bestSelfScore > *bestOpponentScore {
			return *bestSelfScore
		}
		scores = append(scores, newOppScore)
	}
	score, err := min(scores)
	if err != nil {
		panic("oops")
	}

	return -score
}

func chooseBestMove(
	gs *GameState,
	maxDepth int,
	alphaBetaPruning bool,
) Move {
	bestScoreBlack, bestScoreWhite := math.Inf(-1), math.Inf(-1)
	var bestMove Move
	moves := getMoves(gs)

	var bestSelfScore *float64
	switch gs.currentPlayer {
	case black:
		bestSelfScore = &bestScoreBlack
	case white:
		bestSelfScore = &bestScoreWhite
	}

	for _, move := range moves {
		var newScore float64
		log.Printf("Evaluate move %v", move)
		newState := makeMove(*gs, move)
		newScore = -evaluateBoard(
			&newState,
			maxDepth,
			bestScoreBlack,
			bestScoreWhite,
			alphaBetaPruning,
		)

		if newScore > *bestSelfScore {
			bestMove = move
			*bestSelfScore = newScore
		}
		log.Println("-> Score:", newScore)

	}
	return bestMove
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func aiVsAi(maxDepth int, alphaBetaPruning bool) {
	log.SetOutput(ioutil.Discard)
	state := createInitialState()
	clear()
	fmt.Print(state)

	color := white
	dir := up

	for true {
		move := chooseBestMove(&state, maxDepth, alphaBetaPruning)
		log.Println("Make move", move)
		state = makeMove(state, move)
		clear()
		fmt.Print(state)
		time.Sleep(300 * time.Millisecond)
		color = oppositeColor(color)
		dir = oppositeDirection(dir)

		nWhite, nBlack := countTokens(&state.board)
		if nWhite == 0 || nBlack == 0 {
			break
		}
	}
}

func parseMove(playerInput string) (Move, error) {
	regex, _ := regexp.Compile(`(?P<letter>[a-hA-H])(?P<number>[1-8])`)
	match := regex.FindAllStringSubmatch(playerInput, 2)
	if len(match) != 2 {
		return Move{}, errors.New("could not parse input")
	}

	fromI, err := strconv.Atoi(match[0][2])
	if err != nil {
		return Move{}, err
	}
	fromI--
	fromJ := int(match[0][1][0]) - 97
	toI, err := strconv.Atoi(match[1][2])
	if err != nil {
		return Move{}, err
	}
	toI--
	toJ := int(match[1][1][0]) - 97

	return Move{Position{fromI, fromJ}, Position{toI, toJ}}, nil
}

func contains(moves []Move, mv Move) bool {
	for _, move := range moves {
		if move == mv {
			return true
		}
	}
	return false
}

func gameAgainstAI(maxDepth int, alphaBetaPruning bool) {
	f, err := os.OpenFile("./log.out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	reader := bufio.NewReader(os.Stdin)

	state := createInitialState()
	clear()
	fmt.Print(state)

	for true {
		playerInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("No valid move!")
			continue
		}
		move, err := parseMove(playerInput)
		if err != nil {
			log.Println("No valid move!")
			continue
		}

		possibleMoves := getMoves(&state)
		if !contains(possibleMoves, move) {
			fmt.Println("Invalid move! Choose between: ", possibleMoves)
			continue
		}
		log.Println("Making move", move)
		state = makeMove(state, move)
		clear()
		fmt.Print(state)

		// Computer move
		startTime := time.Now()
		move = chooseBestMove(&state, maxDepth, alphaBetaPruning)
		duration := int(time.Now().Sub(startTime).Milliseconds())
		delay := 300 // in milliseconds
		slp := math.Max(float64(delay-duration), 0)
		log.Println("Sleep", slp, "duration:", duration)
		time.Sleep(time.Duration(slp) * time.Millisecond)
		state = makeMove(state, move)
		clear()
		fmt.Print(state)

		nWhite, nBlack := countTokens(&state.board)
		if nWhite == 0 || nBlack == 0 {
			break
		}
	}
}

func main() {

	var searchDepth int
	var alphaBetaPruning bool

	app := &cli.App{
		Name:  "Checkers!",
		Usage: "Play checkers on the command line",
		Commands: []*cli.Command{
			{
				Name:  "auto",
				Usage: "AI vs AI!",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "depth",
						Value:       5,
						Usage:       "Tree search depth",
						Destination: &searchDepth,
					},
					&cli.BoolFlag{
						Name:        "alpha_beta_pruning",
						Aliases:     []string{"a"},
						Usage:       "Use alpha-beta pruning",
						Destination: &alphaBetaPruning,
					},
				},
				Action: func(c *cli.Context) error {
					aiVsAi(searchDepth, alphaBetaPruning)
					return nil
				},
			},
			{
				Name:  "play",
				Usage: "Play checkers!",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "depth",
						Value:       5,
						Usage:       "Tree search depth",
						Destination: &searchDepth,
					},
					&cli.BoolFlag{
						Name:        "alpha_beta_pruning",
						Aliases:     []string{"a"},
						Usage:       "Use alpha-beta pruning",
						Destination: &alphaBetaPruning,
					},
				},
				Action: func(c *cli.Context) error {
					gameAgainstAI(searchDepth, alphaBetaPruning)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
