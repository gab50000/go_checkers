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

// Move holding start and end position
type Move struct {
	From Position
	To   Position
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

func getPositions(
	color playerColor,
	tok tokenType,
	board *[boardSize][boardSize]string,
) []Position {
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

func getMoves(
	color playerColor,
	dir direction,
	board *[boardSize][boardSize]string,
) (moves []Move) {

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

func getManMoves(
	color playerColor,
	dir direction,
	pos Position,
	board *[boardSize][boardSize]string,
) []Move {
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
	color playerColor,
	dir direction,
	pos Position,
	board *[boardSize][boardSize]string) []Move {

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

func getKingMoves(
	color playerColor,
	pos Position,
	board *[boardSize][boardSize]string,
) (moves []Move, jumps []Move) {
	i, j := pos.I, pos.J
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

func evaluateCurrentBoard(color playerColor, board *[boardSize][boardSize]string) float64 {
	nWhite, nBlack := countTokens(board)

	var enemyCount, selfCount int
	switch color {
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

func makeMove(move Move, dir direction, board [boardSize][boardSize]string) [boardSize][boardSize]string {
	noMove := Move{Position{0, 0}, Position{0, 0}}
	if move == noMove {
		return board
	}

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
	return board
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
	color playerColor,
	dir direction,
	board *[boardSize][boardSize]string,
	depthRemaining int,
	bestScoreBlack float64,
	bestScoreWhite float64,
) (score float64) {

	var bestOpponentScore *float64
	var bestSelfScore *float64
	switch color {
	case black:
		bestOpponentScore = &bestScoreWhite
		bestSelfScore = &bestScoreBlack
	case white:
		bestOpponentScore = &bestScoreBlack
		bestSelfScore = &bestScoreWhite

	}

	if depthRemaining == 0 {
		return evaluateCurrentBoard(color, board)
	}
	scores := make([]float64, 0)

	moves := getMoves(color, dir, board)

	if len(moves) == 0 {
		return evaluateCurrentBoard(color, board)
	}

	for _, move := range moves {
		newBoard := makeMove(move, dir, *board)
		newOppScore := evaluateBoard(
			oppositeColor(color),
			oppositeDirection(dir),
			&newBoard,
			depthRemaining-1,
			bestScoreBlack,
			bestScoreWhite,
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
	color playerColor,
	dir direction,
	board *[boardSize][boardSize]string,
	maxDepth int,
) Move {
	bestScoreBlack, bestScoreWhite := math.Inf(-1), math.Inf(-1)
	var bestMove Move
	moves := getMoves(color, dir, board)

	var bestSelfScore *float64
	switch color {
	case black:
		bestSelfScore = &bestScoreBlack
	case white:
		bestSelfScore = &bestScoreWhite
	}

	for _, move := range moves {
		var newScore float64
		log.Print("Evaluate move", move, ":")
		newBoard := makeMove(move, dir, *board)
		newScore = -evaluateBoard(
			oppositeColor(color),
			oppositeDirection(dir),
			&newBoard,
			maxDepth,
			bestScoreBlack,
			bestScoreWhite,
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

func aiVsAi(maxDepth int) {
	log.SetOutput(ioutil.Discard)
	board := getBoard()
	printBoard(&board)

	color := white
	dir := up

	for true {
		move := chooseBestMove(color, dir, &board, maxDepth)
		log.Println("Make move", move)
		board = makeMove(move, dir, board)
		clear()
		printBoard(&board)
		time.Sleep(300 * time.Millisecond)
		color = oppositeColor(color)
		dir = oppositeDirection(dir)

		nWhite, nBlack := countTokens(&board)
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

func gameAgainstAI(maxDepth int) {
	log.SetOutput(ioutil.Discard)
	reader := bufio.NewReader(os.Stdin)

	board := getBoard()
	printBoard(&board)

	color := white
	dir := up

	for true {
		playerInput, err := reader.ReadString('\n')
		if err != nil {
			log.Println("No valid move!")
			continue
		}
		move, err := parseMove(playerInput)
		if err != nil {
			log.Println("No valid move!")
			continue
		}

		possibleMoves := getMoves(color, dir, &board)
		if !contains(possibleMoves, move) {
			log.Println("Invalid move!")
			continue
		}
		board = makeMove(move, dir, board)

		// Computer move
		color = oppositeColor(color)
		dir = oppositeDirection(dir)

		move = chooseBestMove(color, dir, &board, maxDepth)
		board = makeMove(move, dir, board)
		clear()
		printBoard(&board)
		color = oppositeColor(color)
		dir = oppositeDirection(dir)

		nWhite, nBlack := countTokens(&board)
		if nWhite == 0 || nBlack == 0 {
			break
		}
	}
}

func main() {

	var searchDepth int

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
				},
				Action: func(c *cli.Context) error {
					aiVsAi(searchDepth)
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
				},
				Action: func(c *cli.Context) error {
					gameAgainstAI(searchDepth)
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
