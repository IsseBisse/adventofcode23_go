package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	lines, _ := readLines("input.txt")
	// lines, _ := readLines("smallInput.txt")

	re := regexp.MustCompile("[0-9]")
	var calibrationNumbers []int
	for _, line := range lines {
		matches := re.FindAll([]byte(line), -1)

		numberString := string(matches[0])
		if len(matches) == 1 {
			numberString += string(matches[0])
		} else {
			numberString += string(matches[len(matches)-1])
		}

		number, _ := strconv.Atoi(numberString)
		calibrationNumbers = append(calibrationNumbers, number)
	}

	calibrationSum := 0
	for _, number := range calibrationNumbers {
		calibrationSum += number
	}

	fmt.Println(calibrationSum)
}

func partTwo() {

}

func main() {
	partOne()
	partTwo()
}
