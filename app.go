package main

import (
	"fmt"
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

type BreadthFirstSearchNode struct {
	state   State
	history []Move
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
	for _, offset := range rookOffsets {
		to := state.whiteRook
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

func findAllLegalBlackKingMoves(state State) []Move {
	coordsInRangeBlackKing := findCoordsInRangeKing(state.blackKing)
	coordsInRangeWhiteKing := findCoordsInRangeKing(state.whiteKing)
	coords := coordsDifference(coordsInRangeBlackKing, coordsInRangeWhiteKing)
	res := make([]Move, 0)
	for _, coord := range coords {
		move := Move{from: state.blackKing, to: coord, piece: blackKing}
		newState := applyMove(state, move)
		if !isChecked(newState) {
			res = append(res, move)
		}
	}

	return res
}

func findAllLegalWhiteKingMoves(state State) []Move {
	coordsInRangeWhiteKing := findCoordsInRangeKing(state.whiteKing)
	coordsInRangeBlackKing := findCoordsInRangeKing(state.blackKing)
	coords := coordsDifference(coordsInRangeWhiteKing, coordsInRangeBlackKing)
	coords = coordsDifference(coords, []Coord{state.whiteRook})
	res := make([]Move, 0)
	for _, coord := range coords {
		res = append(res, Move{from: state.whiteKing, to: coord, piece: whiteKing})

	}
	return res
}

func findAllLegalWhiteRookMoves(state State) []Move {
	coordsInRangeWhiteRook := findCoordsInRangeRook(state, state.whiteRook)
	coords := coordsDifference(coordsInRangeWhiteRook, []Coord{state.whiteKing})
	res := make([]Move, 0)
	for _, coord := range coords {
		res = append(res, Move{from: state.whiteRook, to: coord, piece: whiteRook})
	}
	return res
}

func findLegalMoves(state State) []Move {
	res := make([]Move, 0)
	if state.movingPlayer == "white" {
		res = append(res, findAllLegalWhiteKingMoves(state)...)
		res = append(res, findAllLegalWhiteRookMoves(state)...)
	} else {
		res = append(res, findAllLegalBlackKingMoves(state)...)
	}

	return res
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

func isCheckmate(state State) bool {
	return isChecked(state) && len(findLegalMoves(state)) == 0
}

func main() {
	var movingPlayer, whiteKing, whiteRook, blackKing string
	fmt.Scan(&movingPlayer, &whiteKing, &whiteRook, &blackKing)
	whiteKingCoord := parseCoord(whiteKing)
	whiteRookCoord := parseCoord(whiteRook)
	blackKingCoord := parseCoord(blackKing)
	state := State{movingPlayer, whiteKingCoord, whiteRookCoord, blackKingCoord}
	winningMoves := findWinningMoves(state)
	fmt.Println(formatMovesSequence(winningMoves))
}

func formatMovesSequence(moves []Move) string {
	res := ""
	for _, move := range moves {
		res += fmt.Sprintf("%s ", move)
	}
	return res
}

func findWinningMoves(state State) []Move {
	visitedState := make(map[State]bool)
	queue := []BreadthFirstSearchNode{{state, []Move{}}}
	var node BreadthFirstSearchNode
	for len(queue) > 0 {
		node, queue = queue[0], queue[1:]
		if isCheckmate(node.state) {
			return node.history
		}
		for _, move := range findLegalMoves(node.state) {
			newState := applyMove(node.state, move)
			if _, ok := visitedState[newState]; !ok {
				visitedState[newState] = true
				newHistory := make([]Move, len(node.history), len(node.history)+1)
				copy(newHistory, node.history)
				newHistory = append(newHistory, move)
				queue = append(queue, BreadthFirstSearchNode{newState, newHistory})
			}
		}
	}
	return nil
}
