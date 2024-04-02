package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

type Game struct {
	number         int
	rounds         [][3]int
	minCubesNeeded [3]int
}

func getMinCubesNeeded(rounds [][3]int) [3]int {
	minCubesNeeded := [3]int{0, 0, 0}
	for _, round := range rounds {
		for idx := 0; idx < 3; idx++ {
			if minCubesNeeded[idx] < round[idx] {
				minCubesNeeded[idx] = round[idx]
			}
		}
	}

	return minCubesNeeded
}

func parseGame(description string) Game {
	// number, _ := regexp.Match("Game ([0-9+])", []byte(description))
	gameRe := regexp.MustCompile("Game ([0-9]+)")
	number, _ := strconv.Atoi(gameRe.FindStringSubmatch(description)[1])

	allRoundsString := strings.Split(description, ":")[1]
	splitRoundsStrings := strings.Split(allRoundsString, ";")

	colors := [3]string{"red", "green", "blue"}
	var cubes [3]int
	var rounds [][3]int
	for _, roundString := range splitRoundsStrings {
		for colorIdx, color := range colors {
			re := regexp.MustCompile("([0-9]+) " + color)
			cubeString := re.FindStringSubmatch(roundString)
			if len(cubeString) > 0 {
				cubes[colorIdx], _ = strconv.Atoi(cubeString[1])
			} else {
				cubes[colorIdx] = 0
			}
		}
		rounds = append(rounds, cubes)
	}

	return Game{number, rounds, getMinCubesNeeded(rounds)}
}

func cubesLeq(first, second [3]int) bool {
	for idx := 0; idx < 3; idx++ {
		if first[idx] > second[idx] {
			return false
		}
	}
	return true
}

func gameComparator(threshold [3]int) func(Game) bool {
	return func(cubes Game) bool {
		return cubesLeq(cubes.minCubesNeeded, threshold)
	}
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func partOne() {
	lines, _ := readLines("input.txt")

	var games []Game
	for _, line := range lines {
		game := parseGame(line)
		games = append(games, game)
	}

	isValidGame := gameComparator([3]int{12, 13, 14})
	validGames := filter(games, isValidGame)

	idSum := 0
	for _, game := range validGames {
		idSum += game.number
	}

	fmt.Println(idSum)
}

func partTwo() {
	lines, _ := readLines("input.txt")

	var games []Game
	for _, line := range lines {
		game := parseGame(line)
		games = append(games, game)
	}

	powSum := 0
	for _, game := range games {
		pow := 1
		for _, num := range game.minCubesNeeded {
			pow *= num
		}
		powSum += pow
	}

	fmt.Println(powSum)
}

func main() {
	// partOne()
	partTwo()
}
