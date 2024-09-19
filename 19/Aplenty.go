package main

import (
	"bufio"
	"fmt"
	"os"
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

var keyToIndex = map[string]int{"x": 0, "m": 1, "a": 2, "s": 3}

type Statement struct {
	Index         int // Index == -1 => default/else statement
	OperatorIsLT  bool
	CompareValue  int
	ConsequentKey string
}

type Workflow struct {
	Key        string
	Statements []Statement
}

func parseStatement(str string) Statement {
	var statement Statement

	if !strings.Contains(str, ":") {
		statement.Index = -1
		statement.ConsequentKey = str
		return statement
	}

	var strParts []string
	statement.OperatorIsLT = strings.Contains(str, "<")
	if statement.OperatorIsLT {
		strParts = strings.Split(str, "<")
	} else {
		strParts = strings.Split(str, ">")
	}

	statement.Index = keyToIndex[strParts[0]]

	strParts = strings.Split(strParts[1], ":")
	statement.ConsequentKey = strParts[1]

	value, _ := strconv.Atoi(strParts[0])
	statement.CompareValue = value

	return statement
}

func parseWorkflow(str string) (string, []Statement) {
	str = strings.Trim(str, "}")
	strParts := strings.Split(str, "{")

	key := strParts[0]

	strStatements := strings.Split(strParts[1], ",")
	var statements []Statement
	for _, strStatement := range strStatements {
		statement := parseStatement(strStatement)
		statements = append(statements, statement)
	}

	return key, statements
}

func parsePart(str string) [4]int {
	str = strings.Trim(str, "{}")
	strParts := strings.Split(str, ",")

	var part [4]int
	for _, strPart := range strParts {
		pair := strings.Split(strPart, "=")
		key := pair[0]
		value, _ := strconv.Atoi(pair[1])

		idx := keyToIndex[key]
		part[idx] = value
	}

	return part
}

func parseInput(lines []string) (map[string][]Statement, [][4]int) {
	inputSplitIndex := 0
	for idx, line := range lines {
		if line == "" {
			inputSplitIndex = idx
		}
	}

	workflowLines := lines[:inputSplitIndex]
	workflows := make(map[string][]Statement)
	for _, line := range workflowLines {
		key, statements := parseWorkflow(line)
		workflows[key] = statements
	}

	var parts [][4]int
	partLines := lines[inputSplitIndex+1:]
	for _, line := range partLines {
		part := parsePart(line)
		parts = append(parts, part)
	}

	return workflows, parts
}

func evaluateStatement(statement Statement, part [4]int) string {
	if statement.Index == -1 {
		return statement.ConsequentKey
	}

	partValue := part[statement.Index]

	var statementIsTrue bool
	if statement.OperatorIsLT {
		statementIsTrue = partValue < statement.CompareValue
	} else {
		statementIsTrue = partValue > statement.CompareValue
	}

	if statementIsTrue {
		return statement.ConsequentKey
	} else {
		return ""
	}
}

func evaluateWorkflow(workflow []Statement, part [4]int) string {
	key := ""
	statementIdx := 0
	for key == "" {
		key = evaluateStatement(workflow[statementIdx], part)
		statementIdx++
	}

	return key
}

func evaluatePart(workflows map[string][]Statement, part [4]int) bool {
	key := "in"

	for !(key == "A" || key == "R") {
		key = evaluateWorkflow(workflows[key], part)
	}

	return key == "A"
}

func sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func partOne() {
	lines, _ := readLines("input.txt")
	workflows, parts := parseInput(lines)

	acceptedPartsSum := 0
	for _, part := range parts {
		isAccepted := evaluatePart(workflows, part)

		if isAccepted {
			acceptedPartsSum += sum(part[:])
		}
	}

	fmt.Println(acceptedPartsSum)
}

type PartRange struct {
	currentKey string
	min        [4]int
	max        [4]int
}

func NewPartRange() PartRange {
	partRange := PartRange{}
	partRange.currentKey = "in"
	partRange.min = [4]int{1, 1, 1, 1}
	partRange.max = [4]int{4000, 4000, 4000, 4000}

	return partRange
}

func numCombinations(partRange PartRange) int {
	combinations := 1
	for idx := 0; idx < 4; idx++ {
		combinations *= partRange.max[idx] - partRange.min[idx]
	}

	return combinations
}

func evaluateStatementRange(statement Statement, partRange PartRange) []PartRange {
	if statement.Index == -1 {
		partRange.currentKey = statement.ConsequentKey
		return []PartRange{partRange}
	}

	abovePartRange := partRange
	belowPartRange := partRange

	if statement.OperatorIsLT {
		belowPartRange.currentKey = statement.ConsequentKey
		belowPartRange.max[statement.Index] = statement.CompareValue - 1
		abovePartRange.min[statement.Index] = statement.CompareValue
	} else {
		abovePartRange.currentKey = statement.ConsequentKey
		abovePartRange.min[statement.Index] = statement.CompareValue + 1
		belowPartRange.max[statement.Index] = statement.CompareValue
	}

	return []PartRange{abovePartRange, belowPartRange}
}

func evaluateWorkflowRange(workflow []Statement, partRange PartRange) []PartRange {
	partRanges := []PartRange{}
	feedForwardPartRange := partRange
	for _, statement := range workflow {
		newPartRanges := evaluateStatementRange(statement, feedForwardPartRange)

		for _, newPartRange := range newPartRanges {
			if newPartRange.currentKey == partRange.currentKey {
				feedForwardPartRange = newPartRange
			} else {
				partRanges = append(partRanges, newPartRange)
			}
		}
	}

	return partRanges
}

func evaluatePartRange(workflows map[string][]Statement, partRange PartRange) {
	partRanges := []PartRange{partRange}
	var finishedPartRanges []PartRange

	for len(partRanges) > 0 {
		partRange := partRanges[0]
		partRanges = partRanges[1:]

		newPartRanges := evaluateWorkflowRange(workflows[partRange.currentKey], partRange)
		for _, newPartRange := range newPartRanges {
			if newPartRange.currentKey == "A" || newPartRange.currentKey == "R" {
				finishedPartRanges = append(finishedPartRanges, newPartRange)
			} else {
				partRanges = append(partRanges, newPartRange)
			}
		}
	}

	acceptedPartRangeSum := 0
	for _, partRange := range finishedPartRanges {
		if partRange.currentKey == "A" {
			acceptedPartRangeSum = numCombinations(partRange)
		}
	}

	fmt.Println(acceptedPartRangeSum)
}

func partTwo() {
	lines, _ := readLines("smallInput.txt")
	workflows, _ := parseInput(lines)

	partRange := NewPartRange()
	evaluatePartRange(workflows, partRange)
}

func main() {
	// partOne()
	partTwo()
}
