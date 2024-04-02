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

type Card struct {
	id             int
	winningNumbers []int
	numbers        []int
}

func parseNumbers(str string) []int {
	var numbers []int

	re := regexp.MustCompile("[0-9]+")
	for _, numString := range re.FindAllString(str, -1) {
		num, _ := strconv.Atoi(numString)
		numbers = append(numbers, num)
	}

	return numbers
}

func parseLine(line string) Card {
	idRe := regexp.MustCompile("Card.*([0-9]+):")
	id, _ := strconv.Atoi(idRe.FindStringSubmatch(line)[1])

	numbersString := strings.Split(line, ":")[1]
	winningNumbers := parseNumbers(strings.Split(numbersString, "|")[0])
	numbers := parseNumbers(strings.Split(numbersString, "|")[1])

	return Card{id, winningNumbers, numbers}
}

func contains(set []int, element int) bool {
	for _, setElement := range set {
		if element == setElement {
			return true
		}
	}
	return false
}

func countMatches(card Card) int {
	numMatches := 0
	for _, number := range card.numbers {
		if contains(card.winningNumbers, number) {
			numMatches += 1
		}
	}

	return numMatches
}

func calculateScore(card Card) int {
	numMatches := countMatches(card)

	score := 0
	if numMatches > 0 {
		score = 1
		for idx := 0; idx < numMatches-1; idx++ {
			score *= 2
		}
	}

	fmt.Println(numMatches, score)

	return score
}

func partOne() {
	lines, _ := readLines("smallInput.txt")

	var cards []Card
	for _, line := range lines {
		cards = append(cards, parseLine(line))
	}

	var scores []int
	for _, card := range cards {
		scores = append(scores, calculateScore(card))
	}

	scoreSum := 0
	for _, score := range scores {
		scoreSum += score
	}

	fmt.Println(scoreSum)
}

func partTwo() {
	lines, _ := readLines("input.txt")

	var cards []Card
	for _, line := range lines {
		cards = append(cards, parseLine(line))
	}

	var numMatches []int
	for _, card := range cards {
		numMatches = append(numMatches, countMatches(card))
	}

	numCardCopies := make([]int, len(cards))
	for idx, _ := range numCardCopies {
		numCardCopies[idx] = 1
	}

	for idx, numCopies := range numCardCopies {
		for offset := 0; offset < numMatches[idx]; offset++ {
			numCardCopies[idx+offset+1] += numCopies
		}
	}

	totalNumCopies := 0
	for _, numCopies := range numCardCopies {
		totalNumCopies += numCopies
	}

	fmt.Println(totalNumCopies)
}

func main() {
	// partOne()
	partTwo()
}
