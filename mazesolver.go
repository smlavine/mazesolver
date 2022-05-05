package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// A Tile represents a location in a Maze.
type Tile = rune

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
	boxTopBottom := "|" + strings.Repeat("-", 2*maze.columns+1) + "|\n"

	builder.WriteString(boxTopBottom)

	for i, row := range maze.grid {
		if i == 0 {
			builder.WriteRune(' ')
		} else {
			builder.WriteRune('|')
		}

		for _, tile := range row {
			builder.WriteRune(' ')
			builder.WriteRune(tile)
		}

		builder.WriteRune(' ')

		if i == len(maze.grid)-1 {
			builder.WriteRune(' ')
		} else {
			builder.WriteRune('|')
		}

		builder.WriteRune('\n')
	}

	builder.WriteString(boxTopBottom)

	return builder.String()
}

func main() {
	var (
		printStart          bool
		printSolutionLength bool
		printSolution       bool
		infile, outfile     string
	)

	flag.BoolVar(&printStart, "d", false,
		"Pretty-print (display) the maze after reading.")
	flag.BoolVar(&printSolutionLength, "s", false,
		"Print length of shortest path or 'No solution'.")
	flag.BoolVar(&printSolution, "p", false,
		"Pretty-print maze with the path, if one exists.")

	flag.StringVar(&infile, "i", "",
		"Read maze from infile. (default: stdin)")
	flag.StringVar(&outfile, "o", "",
		"Write all output to outfile. (default: stdout)")

	flag.Parse()

	var (
		in  io.Reader
		out io.Writer
	)
	if infile == "" {
		in = os.Stdin
	} else {
		var err error
		in, err = os.Open(infile)
		if err != nil {
			log.Fatalf("failed to open '%s': %v\n", infile, err)
		}
		f := in.(*os.File)
		defer f.Close()
	}
	if outfile == "" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create(outfile)
		if err != nil {
			log.Fatalf("failed to open '%s': %v\n", outfile, err)
		}
		f := in.(*os.File)
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("failed to close outfile '%s': %v\n",
					err)
			}
		}()
	}

	maze := NewMaze(in)

	if printStart {
		fmt.Fprint(out, maze.String())
	}
	// TODO: handle the rest of the options
}
