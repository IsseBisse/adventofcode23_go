package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

type Round struct {
	red   int
	green int
	blue  int
}

type Game struct {
	number int
	rounds []Round
}

func parseGame(description string) int {
	gameRe := regexp.MustCompile("Game ([0-9]+)")
	number := gameRe.FindAll([]byte(description), 0)

	return Game(1, nil)
}

func partOne() {
	lines, _ := readLines("smallInput.txt")

	var games []Game
	for _, line := range lines {
		game := parseGame(line)
		// games = append(games, game)
	}
	fmt.Println(lines)
}

func partTwo() {
	lines, _ := readLines("smallInput.txt")
	fmt.Println(lines)
}

func main() {
	partOne()
	// partTwo()
}
