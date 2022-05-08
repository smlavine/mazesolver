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
type Tile struct {
	char rune
	row  int
	col  int
}

const (
	// OPENING is a rune to represent a Tile that can be traversed.
	OPENING = '.'

	// OBSTRUCTION is a rune to represent a Tile that cannot be traversed.
	OBSTRUCTION = '#'

	// ROUTE is a rune to represent a Tile that is known to be part of a
	// route to solve the Maze.
	ROUTE = '+'
)

// A Maze is a field of Tiles that can be traversed through.
type Maze struct {
	grid [][]Tile
	rows int
	cols int
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
		grid: make([][]Tile, 0),
		rows: 0,
		cols: 0,
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		row := make([]Tile, 0)
		col := 0
		for _, char := range scanner.Text() {
			var t rune
			switch char {
			case '0':
				t = OPENING
			case '1':
				t = OBSTRUCTION
			default:
				continue // Only append to row if 0 or 1 found
			}
			row = append(row, Tile{
				char: t,
				row:  maze.rows,
				col:  col,
			})
			col++
		}
		maze.grid = append(maze.grid, row)
		maze.rows++
	}

	// We can assume that all rows are of the same length. This gets the
	// actual length of the row, but will still work if there are none --
	// maze.cols will remain 0.
	for _, row := range maze.grid {
		maze.cols = len(row)
		break
	}

	return maze
}

// Solve solves the Maze using a Breadth-First-Search algorithm. The amount of
// steps it took to solve the maze is returned, or zero if the maze could not
// be solved.
func (maze *Maze) Solve() (steps uint) {
	start := &maze.grid[0][0]
	if start.char != OPENING {
		// Let's not waste our time.
		// We also can rely on this assumption that the start of the
		// maze is valid at the end of this method, when we are
		// assigning the ROUTE tiles.
		return
	}

	end := &maze.grid[maze.rows-1][maze.cols-1]

	visited := make(map[*Tile]bool)

	predecessors := make(map[*Tile]*Tile)

	var queue []*Tile
	queue = append(queue, start)

	for len(queue) > 0 {
		t := queue[0]
		queue = queue[1:]

		if t == end {
			break // At maze exit
		}

		nbrs := make([]*Tile, 0, 4)

		// We must check the bounds of t.row and t.col to make sure
		// we aren't trying to access a neighbor that doesn't exist.
		if t.row > 0 {
			nbrs = append(nbrs, &maze.grid[t.row-1][t.col])
		}
		if t.row < maze.rows-1 {
			nbrs = append(nbrs, &maze.grid[t.row+1][t.col])
		}
		if t.col > 0 {
			nbrs = append(nbrs, &maze.grid[t.row][t.col-1])
		}
		if t.col < maze.cols-1 {
			nbrs = append(nbrs, &maze.grid[t.row][t.col+1])
		}

		for _, neighbor := range nbrs {
			// If neighbor is not traversable, skip
			if neighbor.char == OBSTRUCTION {
				continue
			}

			// Don't visit the neighbor if we already did.
			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			predecessors[neighbor] = t
			queue = append(queue, neighbor)
		}
	}

	// If we didn't visit the end of the maze, then there was no solution.
	if !visited[end] {
		return
	}

	for t := end; t != start; t = predecessors[t] {
		t.char = ROUTE
		steps++
	}
	start.char = ROUTE
	steps++

	return
}

func (maze *Maze) String() string {
	var builder strings.Builder

	// Bar to be placed at the top and bottom of the box around the maze.
	boxTopBottom := "|" + strings.Repeat("-", 2*maze.cols+1) + "|\n"

	builder.WriteString(boxTopBottom)

	for i, row := range maze.grid {
		if i == 0 {
			builder.WriteRune(' ')
		} else {
			builder.WriteRune('|')
		}

		for _, tile := range row {
			builder.WriteRune(' ')
			builder.WriteRune(tile.char)
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

	if steps := maze.Solve(); printSolutionLength {
		if steps == 0 {
			fmt.Fprintln(out, "No solution.")
		} else {
			fmt.Fprintf(out, "Solution in %v steps.\n", steps)
		}
	}

	if printSolution {
		fmt.Fprint(out, maze.String())
	}
}
