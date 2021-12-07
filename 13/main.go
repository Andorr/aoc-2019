package main

import (
	"aoc13/intcode"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Printf("Part one: %d number of blocks\n", numOfBlocks())
	// fmt.Printf("Part two: %d\n", numOfBlocks())
	fmt.Printf("Part two: %d points\n", playGame(false, false))
}

func numOfBlocks() int {
	// Set up Intcode computer
	program := intcode.ParseFromFile("input.txt")
	computer := intcode.NewIntCode()
	computer.SetMemory(program)
	computer.Execute()

	tiles := map[string]int{}
	output := computer.Output()
	for i := 0; i < len(output); i += 3 {
		tiles[posString(output[i], output[i+1])] = output[i+2]
	}

	blockCount := 0
	for _, value := range tiles {
		if value == 2 {
			blockCount++
		}
	}
	return blockCount
}

func playGame(print bool, selfPlay bool) int {
	// Initialize intcode computer
	program := intcode.ParseFromFile("input.txt")
	program[0] = 2 // Inserted quarters - is necessary for starting the game
	computer := intcode.NewIntCode()
	computer.SetMemory(program)

	// Initialize game
	game := map[string]int{}
	isGameInit := false
	score := 0
	ballPosX := 0
	paddlePosX := 0

	// Initialize input
	computer.InputHandler = func() int {
		// If the game has not been initialized yet, setup the tiles based on the output
		if !isGameInit {
			output := computer.Output()
			tiles := output[:len(output)-3]
			for i := 0; i < len(tiles); i += 3 {
				game[posString(tiles[i], tiles[i+1])] = tiles[i+2]
			}
			score = output[len(output)-1] // The last instruction is a -1, 0, 0 - which sets the score to 0
			isGameInit = true
		}

		input := 0
		if !selfPlay {
			// Make decent AI
			if ballPosX > paddlePosX {
				input = 1
			} else if ballPosX < paddlePosX {
				input = -1
			}
		} else {
			// Read input from user
			input = intcode.DefaultInputHandler()
		}

		// Visualize the game
		if print {
			fmt.Printf("Current score: %d\n", score)
			drawGame(game)
			if !selfPlay {
				clearScreenAndWait(20 * time.Millisecond)
			}
		}
		return input
	}

	// Update game state based on output
	outputIterator := 0
	instruction := []int{}
	computer.OnOutputRecieved = func(output int) {
		// Don't do anything until the game is initialized
		if !isGameInit {
			return
		}

		outputIterator++
		instruction = append(instruction, output)

		// Execute game instruction
		if outputIterator == 3 {
			x, y, tile := instruction[0], instruction[1], instruction[2]
			if x == -1 && y == 0 {
				score = tile
			} else {
				game[posString(x, y)] = tile
			}

			// If the tile is the ball, update the x position. Ditto for the paddle
			if tile == 4 {
				ballPosX = x
			} else if tile == 3 {
				paddlePosX = x
			}

			// Reset
			instruction = []int{}
			outputIterator = 0
		}
	}

	computer.Execute()

	if print {
		fmt.Print("End game:\n\n")
		drawGame(game)
	}

	return score
}

func drawGame(game map[string]int) {
	var maxX, maxY int
	var xPos, yPos, tiles []int = []int{}, []int{}, []int{}
	for key, tile := range game {
		// Parse position
		positions := strings.Split(key, ",")
		x, _ := strconv.Atoi(positions[0])
		y, _ := strconv.Atoi(positions[1])
		xPos = append(xPos, x)
		yPos = append(yPos, y)
		tiles = append(tiles, tile)

		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	tileTypes := []string{" ", "W", "B", "P", "O"}
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			tile := game[posString(x, y)]
			if tile > 4 {
				log.Fatalf("Invalid tile: %d\n", tile)
			}
			fmt.Print(tileTypes[tile])
		}
		fmt.Println()
	}
}

func posString(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func clearScreenAndWait(t time.Duration) {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
	time.Sleep(t)
}
