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

type Converter struct {
	destStart   int
	sourceStart int
	len         int
}

func contains(converter Converter, number int) bool {
	if number >= converter.sourceStart && number < (converter.sourceStart+converter.len) {
		return true
	}
	return false
}

func convert(converter Converter, number int) int {
	offset := number - converter.sourceStart
	return converter.destStart + offset
}

func convertNumber(converters []Converter, number int) int {
	for _, converter := range converters {
		if contains(converter, number) {
			return convert(converter, number)
		}
	}
	return number
}

func parseConverter(convrterString string) []Converter {
	lines := strings.Split(convrterString, "\n")

	var converters []Converter
	re := regexp.MustCompile("[0-9]+")
	for _, line := range lines[1:] {
		matches := re.FindAllString(line, -1)

		var converterArgs []int
		for _, match := range matches {
			number, _ := strconv.Atoi(match)
			converterArgs = append(converterArgs, number)
		}

		converters = append(converters, Converter{converterArgs[0], converterArgs[1], converterArgs[2]})
	}

	return converters
}

func parseInput(lines []string) ([]int, [][]Converter) {
	re := regexp.MustCompile("[0-9]+")
	var seeds []int
	for _, seedString := range re.FindAllString(lines[0], -1) {
		seedNumber, _ := strconv.Atoi(seedString)
		seeds = append(seeds, seedNumber)
	}

	var converters [][]Converter
	for _, converterString := range strings.Split(strings.Join(lines[2:], "\n"), "\n\n") {
		converters = append(converters, parseConverter(converterString))
	}

	return seeds, converters
}

func partOne() {
	lines, _ := readLines("input.txt")
	seeds, converters := parseInput(lines)

	var locations []int
	for _, seed := range seeds {
		value := seed
		for _, converter := range converters {
			value = convertNumber(converter, value)
		}
		locations = append(locations, value)
	}

	var minLocation int
	for idx, location := range locations {
		if idx == 0 || location < minLocation {
			minLocation = location
		}
	}

	fmt.Println(minLocation)
}

func partTwo() {
	// lines, _ := readLines("smallInput.txt")
}

func main() {
	partOne()
	partTwo()
}
