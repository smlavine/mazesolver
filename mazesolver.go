package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Tile rune

const (
	OPENING     = Tile('.')
	OBSTRUCTION = Tile('#')
	ROUTE       = Tile('+')
)

type Maze struct {
	grid    [][]Tile
	rows    int
	columns int
}

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

func StdinZeroesAndOnes() (zeroes uint, ones uint) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		switch scanner.Text() {
		case "0":
			zeroes++
		case "1":
			ones++
		}
	}

	return zeroes, ones
}

func main() {
	maze := NewMaze(os.Stdin)
	fmt.Print(maze.String())
	//zeroes, ones := StdinZeroesAndOnes()

	//sumdec := float64(zeroes + ones)

	//fmt.Printf("0s: %d\n1s: %d\n0s ratio: %.2f\n1s ratio: %.2f\n",
	//	zeroes, ones, float64(zeroes)/sumdec, float64(ones)/sumdec)
}
