package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
)

func main() {
	// Part 1
	signal := readSignal("input.txt")
	decoded := fft(signal, []int{0, 1, 0, -1}, 100, 8)
	fmt.Printf("Part 1: %s\n", sliceToString(decoded))

	// Part 2
	signal = readSignal("input.txt")
	signal = repeatSignal(signal, 10000)
	decoded = smartFFT(signal, 100, 8, 7)
	fmt.Printf("Part 2: %s\n", sliceToString(decoded))
}

func fft(input, pattern []int, numOfPhases int, n int) []int {
	curInput := input

	// For each phase
	for i := 0; i < numOfPhases; i++ {
		output := []int{}

		// For each output element
		for j := range curInput {

			// Create new pattern based on the position
			newPattern := make([]int, 0)
			for _, value := range pattern {
				for k := 0; k < j+1; k++ {
					newPattern = append(newPattern, value)
				}
			}

			// Calculate new output element
			sum := 0
			patternIndex := 1 // The pattern is always shifted one step to the left
			for _, value := range curInput {
				sum += value * newPattern[patternIndex]
				patternIndex = (patternIndex + 1) % len(newPattern)
			}
			output = append(output, abs(sum)%10)
		}
		curInput = output
	}
	return curInput[:n]
}

func smartFFT(input []int, numOfPhases, n, offsetLength int) []int {

	// Calculate the message offset
	offsetDigits := input[:offsetLength]
	messageOffset := 0
	for i := len(offsetDigits) - 1; i >= 0; i-- {
		messageOffset += offsetDigits[i] * int(math.Pow10(len(offsetDigits)-i-1))
	}

	// Since every number in the pattern above the messageOffset is just 1s and the ones before it is just 0s, we just have
	// to the take the sum of the input. For every element
	curInput := input[messageOffset:]
	for phase := 0; phase < numOfPhases; phase++ {

		// Start from the last element in the list and calculate the sum for that element.
		// For example: The last element's pattern contains only [0, 0, 0 ... 0, 0 , 1], but the second last
		// element has the pattern of [0, 0 , 0 ... 0, 0, 1, 1].
		sum := 0
		for j := len(curInput) - 1; j >= 0; j-- {
			sum += curInput[j]
			curInput[j] = abs(sum) % 10 // Store the first digit of the sum
		}

	}
	return curInput[:n]
}

func repeatSignal(signal []int, n int) []int {
	newSignal := make([]int, 0)
	for i := 0; i < n; i++ {
		for _, value := range signal {
			newSignal = append(newSignal, value)
		}
	}
	return newSignal
}

func readSignal(fileName string) []int {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	signalRaw := string(bytes)
	signal := make([]int, 0)
	for _, c := range signalRaw {
		value, err := strconv.Atoi(string(c))
		if err == nil {
			signal = append(signal, value)
		}
	}
	return signal
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sliceToString(a []int) string {
	var output string
	for _, value := range a {
		output += strconv.Itoa(value)
	}
	return output
}
