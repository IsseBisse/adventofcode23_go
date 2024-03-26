package main

import (
	"bufio"
	"os"
)

func readLines(path string) ([]string, error) {
	file, _ := os.Open(path)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func partOne() {
	lines, _ := readLines("smallInput.txt")
}

func partTwo() {
	lines, _ := readLines("smallInput.txt")
}

func main() {
	partOne()
	partTwo()
}
