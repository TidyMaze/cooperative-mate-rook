package main

import (
	"fmt"
	"os"
)

/**
 * Find the shortest sequence of cooperative moves to checkmate the black king.
 **/

type Coord struct {
	x int
	y int
}

type State struct {
	movingPlayer string
	whiteKing    Coord
	whiteRook    Coord
	blackKing    Coord
}

type Piece int8

const (
	whiteKing = iota
	whiteRook
	blackKing
)

type Move struct {
	piece Piece
	from  Coord
	to    Coord
}

var kingOffsets = [][]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

var rookOffsets = [][]int{
	{-1, 0},
	{0, -1},
	{1, 0},
	{0, 1},
}

func (c Coord) String() string {
	return fmt.Sprintf("%c%c", c.x+'a', (8-c.y-1)+'1')
}

func (s State) String() string {
	return fmt.Sprintf("%s %s %s %s", s.movingPlayer, s.whiteKing, s.whiteRook, s.blackKing)
}

func (m Move) String() string {
	return fmt.Sprintf("%s%s", m.from, m.to)
}

// "a8" should be parsed as Coord{x: 0, y: 0}
// "a1" should be parsed as Coord{x: 0, y: 7}
// "h8" should be parsed as Coord{x: 7, y: 0}
// "h1" should be parsed as Coord{x: 7, y: 7}
func parseCoord(s string) Coord {
	columnStr := string(s[0])
	rowStr := string(s[1])
	column := int(columnStr[0] - 'a')
	row := int(rowStr[0] - '1')
	return Coord{x: column, y: 8 - row - 1}
}

func findCoordsInRangeKing(from Coord) []Coord {
	res := make([]Coord, 0)
	for _, offset := range kingOffsets {
		to := Coord{x: from.x + offset[0], y: from.y + offset[1]}
		if to.x < 0 || to.x > 7 || to.y < 0 || to.y > 7 {
			continue
		}
		res = append(res, to)
	}
	return res
}

func findCoordsInRangeRook(state State, from Coord) []Coord {
	res := make([]Coord, 0)
	for _, offset := range rookOffsets {
		to := from

		for {
			to = Coord{x: to.x + offset[0], y: to.y + offset[1]}
			if to.x < 0 || to.x > 7 || to.y < 0 || to.y > 7 {
				break
			}

			if to == state.whiteKing || to == state.blackKing {
				break
			}

			res = append(res, to)
		}
	}
	return res
}

func isChecked(state State) bool {
	from := state.whiteRook

	for _, offset := range rookOffsets {
		to := from

		for {
			to = Coord{x: to.x + offset[0], y: to.y + offset[1]}
			if to.x < 0 || to.x > 7 || to.y < 0 || to.y > 7 {
				break
			}

			if to == state.blackKing {
				return true
			}

			if to == state.whiteKing {
				break
			}
		}
	}
	return false
}

// all coords that are in coordsA but not in coordsB
func coordsDifference(coordsA []Coord, coordsB []Coord) []Coord {
	res := make([]Coord, 0)
	for _, coordA := range coordsA {
		found := false
		for _, coordB := range coordsB {
			if coordA == coordB {
				found = true
				break
			}
		}
		if !found {
			res = append(res, coordA)
		}
	}
	return res
}

func findLegalMoves(state State) []Move {
	coordsInRangeBlackKing := findCoordsInRangeKing(state.blackKing)
	coordsInRangeWhiteKing := findCoordsInRangeKing(state.whiteKing)
	coordsInRangeWhiteRook := findCoordsInRangeRook(state, state.whiteRook)

	// debug("coords", fmt.Sprintf("blackKing: %v, whiteKing: %v, whiteRook: %v", coordsInRangeBlackKing, coordsInRangeWhiteKing, coordsInRangeWhiteRook))

	legalMoves := make([]Move, 0)

	if state.movingPlayer == "white" {
		// white king moves (range white king, minus range black king)
		legalWhiteKingCoords := coordsDifference(coordsInRangeWhiteKing, coordsInRangeBlackKing)

		// add white king moves
		for _, coord := range legalWhiteKingCoords {
			legalMoves = append(legalMoves, Move{from: state.whiteKing, to: coord, piece: whiteKing})
		}

		// add white rook moves
		for _, coord := range coordsInRangeWhiteRook {
			legalMoves = append(legalMoves, Move{from: state.whiteRook, to: coord, piece: whiteRook})
		}
	} else {
		// black king moves (range black king, minus range white king, minus range white rook)
		legalBlackKingCoords := coordsDifference(coordsInRangeBlackKing, coordsInRangeWhiteKing)

		for _, coord := range legalBlackKingCoords {
			move := Move{from: state.blackKing, to: coord, piece: blackKing}
			newState := applyMove(state, move)
			if !isChecked(newState) {
				legalMoves = append(legalMoves, move)
			}
		}
	}

	return legalMoves
}

func applyMove(state State, move Move) State {
	if move.piece == whiteKing {
		return State{
			movingPlayer: "black",
			whiteKing:    move.to,
			whiteRook:    state.whiteRook,
			blackKing:    state.blackKing,
		}
	} else if move.piece == blackKing {
		return State{
			movingPlayer: "white",
			whiteKing:    state.whiteKing,
			whiteRook:    state.whiteRook,
			blackKing:    move.to,
		}
	} else if move.piece == whiteRook {
		return State{
			movingPlayer: "black",
			whiteKing:    state.whiteKing,
			whiteRook:    move.to,
			blackKing:    state.blackKing,
		}
	} else {
		panic("unknown piece")
	}
}

func debug(message string, values ...interface{}) {
	// print to stderr
	fmt.Fprintf(os.Stderr, "%s%v\n", message, values)
}

func isCheckmate(state State) bool {
	return isChecked(state) && len(findLegalMoves(state)) == 0
}

func main() {
	// movingPlayer: Either black or white
	// whiteKing: Position of the white king, e.g. a2
	// whiteRook: Position of the white rook
	// blackKing: Position of the black king
	var movingPlayer, whiteKing, whiteRook, blackKing string
	fmt.Scan(&movingPlayer, &whiteKing, &whiteRook, &blackKing)

	whiteKingCoord := parseCoord(whiteKing)
	whiteRookCoord := parseCoord(whiteRook)
	blackKingCoord := parseCoord(blackKing)

	state := State{movingPlayer, whiteKingCoord, whiteRookCoord, blackKingCoord}

	debug("state", state)

	winningMoves := findWinningMoves(state)

	// Write a sequence of moves (a single move is, e.g. a2b1) separated by spaces
	fmt.Println(formatMovesSequence(winningMoves))
}

func formatMovesSequence(moves []Move) string {
	res := ""
	for _, move := range moves {
		res += fmt.Sprintf("%s ", move)
	}
	return res
}

type BreadthFirstSearchNode struct {
	state   State
	history []Move
}

// Breadth-first search until the first checkmate is found
func findWinningMoves(state State) []Move {

	visitedState := make(map[State]bool)

	queue := make([]BreadthFirstSearchNode, 0)
	queue = append(queue, BreadthFirstSearchNode{state, []Move{}})

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		checkMated := isCheckmate(node.state)
		if checkMated {
			debug("Checkmate found", node.history)
			return node.history
		}

		legalMoves := findLegalMoves(node.state)

		if len(node.history) > 0 {
			lastMove := node.history[len(node.history)-1]
			if lastMove.piece == whiteKing && node.state.whiteKing != lastMove.to {
				debug("Path", node.history, "state", node.state, "legalMoves", legalMoves, "checkMated", checkMated)
				panic("wrong white king")
			} else if lastMove.piece == blackKing && node.state.blackKing != lastMove.to {
				debug("Path", node.history, "state", node.state, "legalMoves", legalMoves, "checkMated", checkMated)
				panic("wrong black king")
			} else if lastMove.piece == whiteRook && node.state.whiteRook != lastMove.to {
				debug("Path", node.history, "state", node.state, "legalMoves", legalMoves, "checkMated", checkMated)
				panic("wrong white rook")
			}
		}

		for iMove := 0; iMove < len(legalMoves); iMove++ {
			newState := applyMove(node.state, legalMoves[iMove])
			if _, ok := visitedState[newState]; !ok {
				visitedState[newState] = true
				newHistory := make([]Move, len(node.history), len(node.history)+1)
				copy(newHistory, node.history)
				newHistory = append(newHistory, legalMoves[iMove])
				queue = append(queue, BreadthFirstSearchNode{newState, newHistory})
			}
		}
	}

	return nil
}
