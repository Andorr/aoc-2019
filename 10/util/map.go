package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Map struct {
	M [][]rune
	W int
	H int
}

func NewMap(m [][]rune) *Map {
	return &Map{
		M: m,
		W: len(m[0]),
		H: len(m),
	}
}

func (m *Map) InBounds(x, y int) bool {
	return x >= 0 && x < m.W && y >= 0 && y < m.H
}

func (m *Map) Draw() {
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			fmt.Print(string(m.M[y][x]))
		}
		fmt.Println()
	}
}

func (m *Map) GetAsteroids() []*Position {
	asteroids := make([]*Position, 0)
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if m.M[y][x] == '#' {
				asteroids = append(asteroids, &Position{X: x, Y: y})
			}
		}
	}
	return asteroids
}

func ReadMap(fileName string) *Map {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make([][]rune, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		runes := make([]rune, 0)
		line := scanner.Text()
		for _, c := range line {
			runes = append(runes, c)
		}
		lines = append(lines, runes)
	}

	return NewMap(lines)
}
