package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// A Tile represents a location in a Maze.
type Tile rune

const (
	// OPENING is a Tile in a Maze that can be traversed.
	OPENING = Tile('.')

	// OBSTRUCTION is a Tile in a Maze that cannot be passed.
	OBSTRUCTION = Tile('#')

	// ROUTE is a Tile in a Maze that is known to be part of a route to
	// solve the Maze.
	ROUTE = Tile('+')
)

// A Maze is a field of Tiles that can be traversed through.
type Maze struct {
	grid    [][]Tile
	rows    int
	columns int
}

// NewMaze returns a new Maze, with data read from r. The data expected is
// text in the form of
//
//	0 1 0 1 0 0 1 0 0 0
//	1 0 1 1 0 0 0 1 0 1
//	0 1 0 0 1 0 0 0 0 0
//	0 0 0 0 0 1 0 1 0 1
//	0 0 1 0 0 0 1 0 0 0
func NewMaze(r io.Reader) *Maze {
	maze := &Maze{
		grid:    make([][]Tile, 0),
		rows:    0,
		columns: 0,
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		row := make([]Tile, 0)
		for _, r := range scanner.Text() {
			if r == '0' {
				row = append(row, OPENING)
			} else if r == '1' {
				row = append(row, OBSTRUCTION)
			}
		}
		maze.grid = append(maze.grid, row)
		maze.rows++
	}

	// We can assume that all rows are of the same length. This gets the
	// actual length of the row, but will still work if there are none --
	// maze.columns will remain 0.
	for _, row := range maze.grid {
		maze.columns = len(row)
		break
	}

	return maze
}

func (maze *Maze) String() string {
	var builder strings.Builder

	// Bar to be placed at the top and bottom of the box around the maze.
	boxTopBottom := "|" + strings.Repeat("-", maze.columns) + "|\n"

	builder.WriteString(boxTopBottom)

	// TODO: render maze.

	builder.WriteString(boxTopBottom)

	return builder.String()

}

func main() {
	maze := NewMaze(os.Stdin)
	fmt.Print(maze.String())
}
