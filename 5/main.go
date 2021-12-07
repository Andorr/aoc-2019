package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	positionMode = iota
	parameterMode
)
const supportedParamCount = 3

// INPUT: 1 for part 1, 5 for part 2
func main() {
	/* program := []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
	1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
	999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99} */
	// program := []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}
	program := readInputFile()
	execute(program, 0)
}

func getParam(program []int, pos int, mode int) int {

	if mode == positionMode {
		return program[program[pos]]
	} else if mode == parameterMode {
		return program[pos]
	}
	log.Fatalf("Invalid mode %d", mode)
	return -1
}

func execute(program []int, i int) {
	// Get the current program
	digits := getDigits(program[i])
	digits = append(digits, make([]int, (2+supportedParamCount)-len(digits))...) // Add unprovided modes as 0s
	opcode := digits[0] + digits[1]*10
	modes := digits[2:]

	switch opcode {
	case 1:
		{
			program[program[i+3]] = getParam(program, i+1, modes[0]) + getParam(program, i+2, modes[1])
			execute(program, i+4)
		}
	case 2:
		{
			program[program[i+3]] = getParam(program, i+1, modes[0]) * getParam(program, i+2, modes[1])
			execute(program, i+4)
		}
	case 3:
		{
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter input number: ")
			numberText, err := reader.ReadString('\n')
			input, err := strconv.Atoi(strings.Trim(numberText, "\r\n"))
			if err != nil {
				log.Fatal(err)
			}
			program[getParam(program, i+1, 1)] = input
			execute(program, i+2)
		}
	case 4:
		{
			fmt.Println(getParam(program, i+1, modes[0]))
			execute(program, i+2)
		}
	case 5:
		{
			pointer := i + 3
			if getParam(program, i+1, modes[0]) != 0 {
				pointer = getParam(program, i+2, modes[1])
			}
			execute(program, pointer)
		}
	case 6:
		{
			pointer := i + 3
			if getParam(program, i+1, modes[0]) == 0 {
				pointer = getParam(program, i+2, modes[1])
			}
			execute(program, pointer)
		}
	case 7:
		{
			value := 0
			if getParam(program, i+1, modes[0]) < getParam(program, i+2, modes[1]) {
				value = 1
			}
			program[getParam(program, i+3, 1)] = value
			execute(program, i+4)
		}
	case 8:
		{
			value := 0
			if getParam(program, i+1, modes[0]) == getParam(program, i+2, modes[1]) {
				value = 1
			}
			program[getParam(program, i+3, 1)] = value
			execute(program, i+4)
		}
	case 99:
		{
			return
		}
	default:
		{
			log.Fatalf("Invalid opcode: %d at %d\n%+v", opcode, i, program)
		}
	}

	return
}

func readInputFile() []int {
	bytes, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	strValues := strings.Split(string(bytes), ",")
	values := make([]int, 0)
	for _, v := range strValues {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		values = append(values, i)
	}

	return values
}

func getDigits(value int) []int {
	digits := make([]int, 0)
	for value > 0 {
		digits = append(digits, value%10)
		value /= 10
	}
	return digits
}
