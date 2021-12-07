package main

import (
	"aoc9/intcode"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	// input := []int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99} // Should produce copy of itself
	// input := []int{1102, 34915192, 34915192, 7, 4, 7, 99, 0} // Should output a 16-bit number
	// input := []int{104, 1125899906842624, 99}
	input := readInputFile()

	computer := intcode.NewIntCode()
	computer.SetMemory(input)
	computer.InputHandler = func() int { return 2 } // Input value should be 2 for part 2, 1 for testing in part 1

	computer.Execute()
	fmt.Println(computer.Output())
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
