package main

import (
	"aoc15/intcode"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	wall = iota + 1
	walkable
	oxygenSystem
	oxygen
)

const north, south, west, east int = 1, 2, 3, 4

func main() {
	m, oxyX, oxyY, distance := exploreAndFindOxygenSystem()
	start := time.Now()
	fmt.Printf("Part 1: %d\n", distance)
	end := time.Now()
	fmt.Println(end.Sub(start))
	// drawMap(m)
	start = time.Now()
	fmt.Printf("Part 2: %d\n", calcOxygenFillTime(m, oxyX, oxyY))
	end = time.Now()
	fmt.Println(end.Sub(start))
	// drawMap(m)
}

func exploreAndFindOxygenSystem() (m map[string]int, oxygenSystemX, oxygenSystemY, distance int) {
	// Initialize computer
	computer := intcode.NewIntCode()
	program := intcode.ParseFromFile("input.txt")
	computer.SetMemory(program)

	// Initialize map
	m = map[string]int{}
	curX, curY := 0, 0
	m[posString(curX, curY)] = walkable
	latestOutput := 1
	stack := [][]int{}

	computer.InputHandler = func() int {

		// Handle output and draw map.
		if latestOutput == 0 {
			// Found wall, reset position
			m[posString(curX, curY)] = wall
			pos := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			curX, curY = pos[0], pos[1]
		} else if latestOutput == 2 {
			// Found oxygen system
			m[posString(curX, curY)] = oxygenSystem
			oxygenSystemX, oxygenSystemY = curX, curY
		} else if latestOutput == 1 {
			// The robot moved, the path was walkable
			m[posString(curX, curY)] = walkable
		}

		nextX, nextY := -1, -1
		direction := -1
		if _, ok := m[posString(curX, curY+1)]; !ok {
			nextX, nextY = curX, curY+1
			direction = north
		} else if _, ok := m[posString(curX, curY-1)]; !ok {
			nextX, nextY = curX, curY-1
			direction = south
		} else if _, ok := m[posString(curX-1, curY)]; !ok {
			nextX, nextY = curX-1, curY
			direction = west
		} else if _, ok := m[posString(curX+1, curY)]; !ok {
			nextX, nextY = curX+1, curY
			direction = east
		}

		// No path were found, pop from stack
		if direction == -1 && len(stack) > 0 {
			pos := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			curX, curY, direction = pos[0], pos[1], pos[2]
		} else if direction > 0 {
			// New position found, Push current position to stack
			stack = append(stack, []int{curX, curY, oppositeDirection(direction)})
			curX, curY = nextX, nextY
		} else {
			// No where else to go, rip
			// Pause the machine
			computer.Pause()
		}

		return direction
	}

	computer.OnOutputRecieved = func(output int) {
		latestOutput = output
	}

	computer.Execute()

	path := djikstra(m, 0, 0, oxygenSystemX, oxygenSystemY) // Can not guarantee that the part of the first approach to the oxygen system is the shortest path
	distance = len(path)
	return
}

func djikstra(m map[string]int, startX, startY, endX, endY int) []string {
	distances := map[string]int{}
	parents := map[string]string{}
	open := [][]int{[]int{startX, startY}}
	distances[posString(startX, startY)] = 0
	pathFound := false

	for len(open) > 0 {
		// Find the position in "open" with the shortest distance
		// and remove it from the "open" slice
		smallestDistance := int(^uint(0) >> 1)
		index := 0
		for i, pos := range open {
			distance := distances[posString(pos[0], pos[1])]
			if distance < smallestDistance {
				smallestDistance = distance
				index = i
			}
		}
		curPos := open[index]
		open[index] = open[len(open)-1] // Overwrite as last element
		open = open[:len(open)-1]       // Remove the last element

		if curPos[0] == endX && curPos[1] == endY {
			pathFound = true
			break
		}

		// Find neighbours with higher distance than curDistance + 1
		for _, neighbour := range getAdjecent(curPos[0], curPos[1]) {
			slotType := m[posString(neighbour[0], neighbour[1])]
			if slotType == walkable || slotType == oxygenSystem {
				newCost := distances[posString(curPos[0], curPos[1])] + 1
				neighbourDistance, ok := distances[posString(neighbour[0], neighbour[1])]
				if !ok {
					// If the neighbour does not have an distance registered, set to max int
					neighbourDistance = int(^uint(0) >> 1)
				}
				if newCost < neighbourDistance {
					distances[posString(neighbour[0], neighbour[1])] = newCost
					parents[posString(neighbour[0], neighbour[1])] = posString(curPos[0], curPos[1])
					open = append(open, []int{neighbour[0], neighbour[1]})
				}
			}
		}

	}

	if !pathFound {
		return nil
	}

	// A path was found
	path := []string{}
	curPos := posString(endX, endY)
	for curPos != "" {
		curPos = parents[curPos]
		if curPos != "" {
			path = append(path, curPos)
		}
	}

	return path
}

func calcOxygenFillTime(m map[string]int, oxyX, oxyY int) int {

	curOxygenPositions := [][]int{[]int{oxyX, oxyY}}
	minutes := 0
	for len(curOxygenPositions) > 0 {
		minutes++
		newOxygenPositions := make([][]int, 0)
		for _, pos := range curOxygenPositions {
			// Add oxygen to adjecent positions
			for _, neighbour := range getAdjecent(pos[0], pos[1]) {
				// Get the slotType and if it is walkable then add oxygen
				slotType := m[posString(neighbour[0], neighbour[1])]
				if slotType == walkable {
					m[posString(neighbour[0], neighbour[1])] = oxygen
					newOxygenPositions = append(newOxygenPositions, []int{neighbour[0], neighbour[1]})
				}
			}
		}
		curOxygenPositions = newOxygenPositions
	}
	return minutes - 1 // Substracting 1 due to the extra check when there is no place for the oxygen to spread
}

func posString(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func oppositeDirection(dir int) int {
	if dir == north {
		return south
	} else if dir == south {
		return north
	} else if dir == west {
		return east
	} else {
		return west
	}
}

func getAdjecent(x, y int) [][]int {
	return [][]int{
		[]int{x, y + 1},
		[]int{x, y - 1},
		[]int{x + 1, y},
		[]int{x - 1, y},
	}
}

func drawMap(m map[string]int) {
	var minX, maxX, minY, maxY int
	for key := range m {
		// Parse position
		positions := strings.Split(key, ",")
		x, _ := strconv.Atoi(positions[0])
		y, _ := strconv.Atoi(positions[1])

		if x > maxX {
			maxX = x
		} else if x < minX {
			minX = x
		}
		if y > maxY {
			maxY = y
		} else if y < minY {
			minY = y
		}
	}

	slotTypes := []string{"#", ".", "O", "0"}
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			slotType, ok := m[posString(x, y)]
			if !ok {
				fmt.Print("*")
				continue
			}
			if slotType > 4 {
				log.Fatalf("Invalid slot: %d\n", slotType)
			}
			fmt.Print(slotTypes[slotType-1])
		}
		fmt.Println()
	}
}
