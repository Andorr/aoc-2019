package main

import (
	intcode "aoc7/intcode"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	program := readInputFile()

	/* program := []int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33,
	1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0} */
	fmt.Printf("Part 1: %d\n", getMaxThrust(program))

	/* program = []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26,
	27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5} */
	fmt.Printf("Part 2 (Feedback loop): %d\n", getMaxThrustWithLoop(program))
}

func getMaxThrust(program []int) int {
	maxThrust := 0
	phases := calcPermutations([][]int{}, []int{0, 1, 2, 3, 4}, 5)
	for _, phase := range phases {
		thrust := calcThrust(program, phase)
		if thrust > maxThrust {
			maxThrust = thrust
		}
	}
	return maxThrust
}

func getMaxThrustWithLoop(program []int) int {
	maxThrust := 0
	phases := calcPermutations([][]int{}, []int{5, 6, 7, 8, 9}, 5)
	for _, phase := range phases {
		thrust := calcThrustWithLoop(program, phase)
		if thrust > maxThrust {
			maxThrust = thrust
		}
	}
	return maxThrust
}

func calcThrust(program, phaseSetting []int) int {
	computer := intcode.NewIntCode()

	inputValue := 0
	for i := 0; i < len(phaseSetting); i++ {
		inputs := []int{phaseSetting[i], inputValue}
		computer.InputHandler = func() int {
			v := inputs[0]
			inputs = inputs[1:]
			return v
		}
		programCopy := make([]int, len(program))
		copy(programCopy, program)
		computer.ResetPointer()
		computer.SetMemory(programCopy)
		computer.Execute()
		inputValue = computer.Outputs()[i]
	}

	output := computer.Outputs()
	return output[len(output)-1]
}

func calcThrustWithLoop(program, phaseSetting []int) int {
	inputValue := 0 // Amp A starts with taking 0 as input
	hasReceivedPhaseSetting := make([]bool, 5)
	shouldPause := make([]bool, 5)

	// Create amps
	amps := make([]*intcode.IntCode, 5)
	for i := 0; i < 5; i++ {
		amps[i] = intcode.NewIntCode()

		// Set memory to the corresponding program
		programCopy := make([]int, len(program))
		copy(programCopy, program)
		amps[i].SetMemory(programCopy)

		// Initialize input handler and output listener
		amps[i].InputHandler = func(ampIndex int) func() int {
			return func() int {
				input := 0
				if !hasReceivedPhaseSetting[ampIndex] {
					hasReceivedPhaseSetting[ampIndex] = true
					input = phaseSetting[ampIndex]
				} else {
					input = inputValue
				}

				if shouldPause[ampIndex] {
					amps[ampIndex].Pause()
				}

				return input
			}
		}(i)
		amps[i].OnOutputRecieved = func(ampIndex int) func(output int) {
			return func(output int) {
				// On output created the program should pause on next input call
				inputValue = output
				shouldPause[ampIndex] = true
			}
		}(i)
	}

	// Execute amps
	currentAmp := 0
	for {
		finished := amps[currentAmp].Execute()
		if finished && currentAmp == 4 {
			return inputValue
		}
		shouldPause[currentAmp] = false
		currentAmp = (currentAmp + 1) % len(amps)
	}
}

func calcPermutations(phases [][]int, phase []int, s int) [][]int {

	if s == 1 {
		p := make([]int, 5)
		copy(p, phase)
		phases = append(phases, p)
		return phases
	}

	for i := 0; i < s; i++ {
		phases = calcPermutations(phases, phase, s-1)
		index := 0
		if s&1 == 0 {
			index = i
		}
		temp := phase[index]
		phase[index] = phase[s-1]
		phase[s-1] = temp
	}
	return phases
}

func readInputFile() []int {
	bytes, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	strValues := strings.Split(string(bytes), ",")
	values := make([]int, 0)
	for _, v := range strValues {
		i, err := strconv.Atoi(strings.Trim(v, "\n"))
		if err != nil {
			log.Fatal(err)
		}
		values = append(values, i)
	}

	return values
}
