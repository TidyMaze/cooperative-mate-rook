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

func (c Coord) String() string {
	return fmt.Sprintf("%c%c", c.x+'a', (8-c.y-1)+'1')
}

func (s State) String() string {
	return fmt.Sprintf("%s %s %s %s", s.movingPlayer, s.whiteKing, s.whiteRook, s.blackKing)
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

func debug(values ...interface{}) {
	// print to stderr
	fmt.Fprintln(os.Stderr, values...)
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

	debug(state)

	// Write a sequence of moves (a single move is, e.g. a2b1) separated by spaces
	fmt.Println("h5h1 a1a2 b3b8 h1f1")
}
