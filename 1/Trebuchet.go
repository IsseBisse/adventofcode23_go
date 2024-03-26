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

type digitParser func(string) int

func intStringToInt(str string) int {
	number, _ := strconv.Atoi(str)
	return number
}

func getCalibration(path string, re *regexp.Regexp, parsingFunc digitParser) int {
	lines, _ := readLines(path)

	var calibrationNumbers []int
	for _, line := range lines {
		matches := re.FindAll([]byte(line), -1)

		numberString := string(matches[0])
		if len(matches) == 1 {
			numberString += string(matches[0])
		} else {
			numberString += string(matches[len(matches)-1])
		}

		number := intStringToInt(numberString)
		calibrationNumbers = append(calibrationNumbers, number)
	}

	calibrationSum := 0
	for _, number := range calibrationNumbers {
		calibrationSum += number
	}

	return calibrationSum
}

func partOne() {
	re := regexp.MustCompile("[0-9]")
	calibrationSum := getCalibration("input.txt", re, intStringToInt)
	fmt.Println(calibrationSum)
}

var wordsToDigits = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func wordOrIntToInt(str string) int {
	re := regexp.MustCompile("[a-z]")
	if re.Match([]byte(str)) {
		return wordsToDigits[str]
	} else {
		return intStringToInt(str)
	}
}

func partTwo() {
	words := make([]string, len(wordsToDigits))
	i := 0
	for word := range wordsToDigits {
		words[i] = word
		i++
	}
	wordsString := strings.Join(words[:], "|")
	// TODO: Go doesn't support negative look-ahead so this regex must be redone somehow
	re := regexp.MustCompile("(?=([0-9]|" + wordsString + "))")
	calibrationSum := getCalibration("smallInput.txt", re, wordOrIntToInt)
	fmt.Println(calibrationSum)
}

func main() {
	// partOne()
	partTwo()
}
