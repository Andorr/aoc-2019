package main

import (
	"aoc15/intcode"
	"fmt"
)

func main() {

	computer := intcode.NewIntCode()
	program := intcode.ParseFromFile("input.txt")
	computer.SetMemory(program)

	computer.Execute()

	output := computer.Output()
	for _, o := range output {
		fmt.Printf("%c", o)
	}
}
