package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	input := readInputFile()

	program := make([]int, len(input))
	copy(program, input)
	program[1] = 12
	program[2] = 2
	execute(program, 0)
	fmt.Println(program[0])
	fmt.Println(find(input))
}

func find(program []int) (int, int, int) {
	node := 0
	verb := 0

	for verb <= 99 {
		node = 0
		for node <= 99 {
			input := make([]int, len(program))
			copy(input, program)
			input[1] = node
			input[2] = verb
			execute(input, 0)
			if input[0] == 19690720 {
				return node, verb, node*100 + verb
			}
			node++
		}
		verb++
	}
	return -1, -1, -1
}

func execute(program []int, i int) {
	// Get the current program
	opcode := program[i]

	switch opcode {
	case 1:
		{
			program[program[i+3]] = program[program[i+1]] + program[program[i+2]]
			execute(program, i+4)
		}
	case 2:
		{
			program[program[i+3]] = program[program[i+1]] * program[program[i+2]]
			execute(program, i+4)
		}
	case 99:
		{
			return
		}
	default:
		{
			log.Fatalf("Invalid opcode: %d\n%+v", opcode, program)
		}
	}

	return
}

func readInputFile() []int {
	bytes, err := ioutil.ReadFile("data.txt")
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
