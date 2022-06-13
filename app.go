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

var kingOffsets = [][2]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

var rookOffsets = [][2]int{
	{-1, 0},
	{0, -1},
	{1, 0},
	{0, 1},
}

func (c Coord) String() string {
	return fmt.Sprintf("%c%c", c.x+'a', (8-c.y-1)+'1')
}

func (c Coord) addOffset(offset [2]int) Coord {
	return Coord{c.x + offset[0], c.y + offset[1]}
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
		to := from.addOffset(offset)
		if !isValidCoord(to) {
			continue
		}
		res = append(res, to)
	}
	return res
}

func isValidCoord(to Coord) bool {
	outside := to.x < 0 || to.x > 7 || to.y < 0 || to.y > 7
	return !outside
}

func findCoordsInRangeRook(state State, from Coord) []Coord {
	res := make([]Coord, 0)
	for _, offset := range rookOffsets {
		to := from
		for {
			to = to.addOffset(offset)
			if !isValidCoord(to) {
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
			to = to.addOffset(offset)
			if !isValidCoord(to) {
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

func findChildrenNodes(node BreadthFirstSearchNode) []BreadthFirstSearchNode {
	res := make([]BreadthFirstSearchNode, 0)
	for _, move := range findLegalMoves(node.state) {
		newState := applyMove(node.state, move)
		newNode := BreadthFirstSearchNode{state: newState, history: addHistoryCopy(node.history, move)}
		res = append(res, newNode)
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
	fmt.Println(formatMovesSequence(findWinningMoves(parseNewState())))
}

func parseNewState() State {
	var movingPlayer, whiteKing, whiteRook, blackKing string
	fmt.Scan(&movingPlayer, &whiteKing, &whiteRook, &blackKing)
	return State{
		movingPlayer,
		parseCoord(whiteKing),
		parseCoord(whiteRook),
		parseCoord(blackKing),
	}
}

func formatMovesSequence(moves []Move) string {
	res := ""
	for _, move := range moves {
		res += fmt.Sprintf("%s ", move)
	}
	return res
}

func addHistoryCopy(history []Move, move Move) []Move {
	res := make([]Move, len(history)+1)
	copy(res, history)
	res[len(history)] = move
	return res
}

func isAlreadyVisited(state State, cache map[State]bool) bool {
	_, ok := cache[state]
	return ok
}

func setVisited(state State, cache *map[State]bool) {
	(*cache)[state] = true
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
		for _, childNode := range findChildrenNodes(node) {
			if !isAlreadyVisited(childNode.state, visitedState) {
				setVisited(childNode.state, &visitedState)
				queue = append(queue, childNode)
			}
		}
	}
	return nil
}
