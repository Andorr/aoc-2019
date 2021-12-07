package main

import (
	"aoc11/intcode"
	"aoc11/painter"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	program := readInputFile()

	// Part 1
	fmt.Printf("Part 1:\nNum panels painted: %d\n", paintShip(program, painter.Black, true))

	// Part 2
	fmt.Printf("Part 2:\nNum panels painted: %d\n", paintShip(program, painter.White, true))
}

func paintShip(program []int, startColor int, print bool) int {
	robot := painter.NewPainter()
	ship := painter.NewShip()
	ship.Paint(0, 0, startColor)
	shouldMove := false

	computer := intcode.NewIntCode()
	programCopy := make([]int, len(program))
	copy(programCopy, program)

	computer.SetMemory(programCopy)
	computer.InputHandler = func() int {
		return robot.CurrentColor(ship)
	}

	computer.OnOutputRecieved = func(output int) {
		if shouldMove {
			robot.Rotate(output)
			robot.Move(ship)
			shouldMove = false
		} else {
			robot.Paint(ship, output)
			shouldMove = true
		}

	}

	computer.Execute()

	if print {
		fmt.Println()
		ship.Draw(robot)
	}

	return ship.NumPanelsPainted()
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
