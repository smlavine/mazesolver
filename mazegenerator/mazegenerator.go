package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	var (
		rows, cols int
		density    float64
		tries      int
		seed       int64
	)

	flag.IntVar(&rows, "r", 10, "Rows in the matrix")
	flag.IntVar(&cols, "c", 20, "Columns in the matrix")
	flag.Float64Var(&density, "d", 0.7, "Ratio of zeroes to ones")
	flag.IntVar(&tries, "m", 1, "Amount of tries if maze isn't solvable")
	flag.Int64Var(&seed, "s", time.Now().Unix(), "Random seed")

	flag.Parse()

	rand.Seed(seed)

	// Proportional amount of ones and zeroes
	amt := rows * cols
	zeroAmt := int(float64(amt) * density)
	maze := append(make([]byte, zeroAmt),
		bytes.Repeat([]byte{1}, amt-zeroAmt)...)

	rand.Shuffle(len(maze), func(i, j int) {
		maze[i], maze[j] = maze[j], maze[i]
	})

	str := strings.Trim(fmt.Sprintf("%v", maze), "[]")
	// Replace spaces at end of rows with newlines
	for i := 2*cols - 1; i < len(str); i += 2 * cols {
		str = str[:i] + "\n" + str[i+1:]
	}

	fmt.Println(str)
}
