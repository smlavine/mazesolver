package main

import (
	"bufio"
	"fmt"
	"os"
)

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
	zeroes, ones := StdinZeroesAndOnes()

	sumdec := float64(zeroes + ones)

	fmt.Printf("0s: %d\n1s: %d\n0s ratio: %.2f\n1s ratio: %.2f\n",
		zeroes, ones, float64(zeroes)/sumdec, float64(ones)/sumdec)
}
