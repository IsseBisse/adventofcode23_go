package main

import (
	"bufio"
	"fmt"
	"os"
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

func parseMap(lines []string) ([]int, int, [2]int) {
	mapSize := [2]int{len(lines[0]), len(lines)}

	mapString := strings.Join(lines, "")
	var blockedIndicies []int
	var startIndex int
	for idx, char := range mapString {
		if char == 'S' {
			startIndex = idx
		} else if char == '#' {
			blockedIndicies = append(blockedIndicies, idx)
		}
	}

	return blockedIndicies, startIndex, mapSize
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type indexToIndexFunction func(int, [2]int) []int

func generateNextIndexMap(blockedIndicies []int, mapSize [2]int) [][]int {
	var nextIndexMap [][]int

	for idx := 0; idx < mapSize[0]*mapSize[1]; idx++ {
		newIndicies := indexToPossibleNextIndex(idx, mapSize)

		var nonBlockedNewIndicies []int
		for _, newIdx := range newIndicies {
			if !contains(blockedIndicies, newIdx) {
				nonBlockedNewIndicies = append(nonBlockedNewIndicies, newIdx)
			}
		}

		nextIndexMap = append(nextIndexMap, nonBlockedNewIndicies)
	}

	return nextIndexMap
}

func indexToCoord(idx int, mapSize [2]int) [2]int {
	return [2]int{idx % mapSize[0], idx / mapSize[0]}
}

func coordToIndex(coord [2]int, mapSize [2]int) int {
	return coord[0] + coord[1]*mapSize[0]
}

func indexToPossibleNextIndex(index int, mapSize [2]int) []int {
	coord := indexToCoord(index, mapSize)

	var nextIndicies []int
	for axis := 0; axis < 2; axis++ {
		for step := -1; step < 2; step += 2 {
			newCoord := coord
			newCoord[axis] += step

			if newCoord[0] >= 0 && newCoord[0] < mapSize[0] && newCoord[1] >= 0 && newCoord[1] < mapSize[1] {
				nextIndicies = append(nextIndicies, coordToIndex(newCoord, mapSize))
			}
		}
	}

	return nextIndicies
}

func walk(startIndex int, nextIndexMap [][]int, numSteps int) []int {
	currentIndicies := map[int]bool{startIndex: true}

	for step := 0; step < numSteps; step++ {
		newIndicies := map[int]bool{}

		for idx := range currentIndicies {
			nextIndicies := nextIndexMap[idx]

			for _, nextIdx := range nextIndicies {
				newIndicies[nextIdx] = true
			}
		}

		currentIndicies = newIndicies
	}

	var finalIndicies []int
	for idx := range currentIndicies {
		finalIndicies = append(finalIndicies, idx)
	}
	return finalIndicies
}

func partOne() {
	lines, _ := readLines("input.txt")
	blockedIndicies, startIndex, mapSize := parseMap(lines)

	nextIndex := generateNextIndexMap(blockedIndicies, mapSize)
	finalIndicies := walk(startIndex, nextIndex, 64)

	fmt.Println(len(finalIndicies))
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func indexToPossibleNextIndexWrapAround(index int, mapSize [2]int) ([]int, []int) {
	coord := indexToCoord(index, mapSize)

	var nextIndicies []int
	var wrapAroundIndicies []int
	for axis := 0; axis < 2; axis++ {
		for step := -1; step < 2; step += 2 {
			newCoord := coord
			newCoord[axis] = newCoord[axis] + step

			if newCoord[0] >= 0 && newCoord[0] < mapSize[0] && newCoord[1] >= 0 && newCoord[1] < mapSize[1] {
				nextIndicies = append(nextIndicies, coordToIndex(newCoord, mapSize))

			} else {
				modCoord := [2]int{mod(newCoord[0], mapSize[0]), mod(newCoord[1], mapSize[1])}
				wrapAroundIndicies = append(wrapAroundIndicies, coordToIndex(modCoord, mapSize))
			}
		}
	}

	return nextIndicies, wrapAroundIndicies
}

func generateNextIndexMapWrapAround(blockedIndicies []int, mapSize [2]int) []map[int]bool {
	var nextIndexMap []map[int]bool

	for idx := 0; idx < mapSize[0]*mapSize[1]; idx++ {
		newIndicies, wrapAroundIndicies := indexToPossibleNextIndexWrapAround(idx, mapSize)

		nonBlockedNewIndicies := map[int]bool{}
		for _, newIdx := range newIndicies {
			if !contains(blockedIndicies, newIdx) {
				nonBlockedNewIndicies[newIdx] = false
			}
		}

		for _, newIdx := range wrapAroundIndicies {
			if !contains(blockedIndicies, newIdx) {
				nonBlockedNewIndicies[newIdx] = true
			}
		}

		nextIndexMap = append(nextIndexMap, nonBlockedNewIndicies)
	}

	return nextIndexMap
}

func walkWrapAround(startIndex int, nextIndexMap []map[int]bool, numSteps int) int {
	currentIndicies := map[int]int{startIndex: 1}

	for step := 0; step < numSteps; step++ {
		newIndicies := map[int]int{}

		for currentIndex, currentWalkers := range currentIndicies {
			nextIndicies := nextIndexMap[currentIndex]

			for nextIdx, addExtraWalker := range nextIndicies {
				_, hasCurrentWalkers := newIndicies[nextIdx]
				if addExtraWalker && hasCurrentWalkers {
					newIndicies[nextIdx] += currentWalkers
				} else {
					newIndicies[nextIdx] = currentWalkers
				}
			}
		}

		currentIndicies = newIndicies
	}

	finalWalkers := 0
	for _, walkers := range currentIndicies {
		finalWalkers += walkers
	}
	return finalWalkers
}

func partTwo() {
	lines, _ := readLines("smallInput.txt")
	blockedIndicies, startIndex, mapSize := parseMap(lines)

	nextIndex := generateNextIndexMapWrapAround(blockedIndicies, mapSize)
	finalWalkers := walkWrapAround(startIndex, nextIndex, 50)

	fmt.Println(finalWalkers)
}

func main() {
	// partOne()
	partTwo()
}
