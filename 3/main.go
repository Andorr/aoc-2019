package main

import (
	"aoc/util"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getResult(wires [][]string) (int, *util.Point, int, *util.Point) {
	m := util.NewMap()

	for _, w := range wires {
		m.AddWire(w)
	}

	c, p1 := m.ClosestInterDist(&m.Start)
	s, p2 := m.InterPointLowestSteps()
	return c, p1, s, p2
}

func readWires() [][]string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wires := make([][]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		wires = append(wires, strings.Split(input, ","))
	}
	return wires
}

func main() {
	fmt.Println(getResult([][]string{
		[]string{"U7", "R6", "D4", "L4"},
		[]string{"R8", "U5", "L5", "D3"},
	}))

	fmt.Println(getResult([][]string{
		[]string{"R75", "D30", "R83", "U83", "L12", "D49", "R71", "U7", "L72"},
		[]string{"U62", "R66", "U55", "R34", "D71", "R55", "D58", "R83"},
	}))

	fmt.Println(getResult([][]string{
		[]string{"R98", "U47", "R26", "D63", "R33", "U87", "L62", "D20", "R33", "U53", "R51"},
		[]string{"U98", "R91", "D20", "R16", "D67", "R40", "U7", "R15", "U6", "R7"},
	}))

	fmt.Println(getResult(readWires()))
}
