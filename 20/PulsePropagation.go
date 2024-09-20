package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

type Pulse struct {
	fromName string
	toName   string
	isHigh   bool
}

type Module interface {
	GetName() string
	GetChildrenNames() []string
	Send(Pulse) (Module, []Pulse)
}

type Nop struct {
	name string
}

func (module Nop) GetName() string {
	return module.name
}

func (module Nop) GetChildrenNames() []string {
	return nil
}

func (module Nop) Send(pulse Pulse) (Module, []Pulse) {
	return module, nil
}

type Broadcaster struct {
	name          string
	childrenNames []string
}

func (module Broadcaster) GetName() string {
	return module.name
}

func (module Broadcaster) GetChildrenNames() []string {
	return module.childrenNames
}

func (module Broadcaster) Send(pulse Pulse) (Module, []Pulse) {
	var pulses []Pulse
	for _, childName := range module.childrenNames {
		pulses = append(pulses, Pulse{module.GetName(), childName, pulse.isHigh})
	}
	return module, pulses
}

type FlipFlop struct {
	name          string
	isOn          bool
	childrenNames []string
}

func (module FlipFlop) GetName() string {
	return module.name
}

func (module FlipFlop) GetChildrenNames() []string {
	return module.childrenNames
}

func (module FlipFlop) Send(pulse Pulse) (Module, []Pulse) {
	var pulses []Pulse
	if pulse.isHigh {
		return module, pulses
	} else {
		module.isOn = !module.isOn
	}

	for _, childName := range module.childrenNames {
		pulses = append(pulses, Pulse{module.GetName(), childName, module.isOn})
	}
	return module, pulses
}

type Conjuction struct {
	name            string
	parentLastPulse map[string]bool
	childrenNames   []string
}

func (module Conjuction) GetName() string {
	return module.name
}

func (module Conjuction) GetChildrenNames() []string {
	return module.childrenNames
}

func (module Conjuction) Send(pulse Pulse) (Module, []Pulse) {
	module.parentLastPulse[pulse.fromName] = pulse.isHigh

	allPulsesHigh := true
	for _, lastPulse := range module.parentLastPulse {
		allPulsesHigh = allPulsesHigh && lastPulse
	}

	var pulses []Pulse
	for _, childName := range module.childrenNames {
		pulses = append(pulses, Pulse{module.GetName(), childName, !allPulsesHigh})
	}
	return module, pulses
}

func parse(line string) Module {
	var module Module
	parts := strings.Split(line, " -> ")
	name := parts[0]
	children := strings.Split(parts[1], ", ")
	print(parts)

	if line[0] == '%' {
		name = name[1:]
		module = FlipFlop{name, false, children}

	} else if line[0] == '&' {
		name = name[1:]
		module = Conjuction{name, nil, children}

	} else {
		module = Broadcaster{name, children}
	}

	return module
}

func partOne() {
	lines, _ := readLines("input.txt")

	modules := make(map[string]Module)
	var conjunctionModules []Conjuction
	for _, line := range lines {
		module := parse(line)

		conjunctionModule, ok := module.(Conjuction)
		if ok {
			conjunctionModules = append(conjunctionModules, conjunctionModule)
		} else {
			modules[module.GetName()] = module
		}
	}

	for _, conModule := range conjunctionModules {
		parentLastPulse := make(map[string]bool)
		for _, module := range modules {
			if slices.Contains(module.GetChildrenNames(), conModule.GetName()) {
				parentLastPulse[module.GetName()] = false
			}
		}
		conModule.parentLastPulse = parentLastPulse

		modules[conModule.name] = conModule
	}

	highPulses := 0
	lowPulses := 0
	buttonPresses := 1000
	for i := 0; i < buttonPresses; i++ {
		pulses := []Pulse{{"", "broadcaster", false}}
		for len(pulses) > 0 {
			pulse := pulses[0]
			if pulse.isHigh {
				highPulses += 1
			} else {
				lowPulses += 1
			}

			pulses = pulses[1:]

			module, ok := modules[pulse.toName]
			if ok {
				newModule, newPulses := module.Send(pulse)
				pulses = append(pulses, newPulses...)
				modules[pulse.toName] = newModule
			}
		}
	}

	fmt.Println(buttonPresses, lowPulses, highPulses)
	fmt.Println(buttonPresses * lowPulses * highPulses)
}

func partTwo() {
	// lines, _ := readLines("smallInput.txt")
}

func main() {
	partOne()
	// partTwo()
}
