package main

import (
	"aoc10/laser"
	"aoc10/util"
	"fmt"
)

const EMPTY = '.'
const ASTROID = '#'

func main() {
	/* m := util.NewMap([][]rune{
		[]rune{'.', '#', '.', '.', '#'},
		[]rune{'.', '.', '.', '.', '.'},
		[]rune{'#', '#', '#', '#', '#'},
		[]rune{'.', '.', '.', '.', '#'},
		[]rune{'.', '.', '.', '#', '#'},
	})
	fmt.Println(findBestLocation(m))

	m = util.ReadMap("test_02.txt")
	fmt.Println(findBestLocation(m)) */

	m := util.ReadMap("input.txt")
	bestPos, count, _ := findBestLocation(m)
	fmt.Printf("Part 1: %d\n\n", count)

	target := vaporizeAsteroids(m, bestPos.X, bestPos.Y)
	fmt.Printf("Part 2\nPos: %+v\nScore: %d\n", target, 100*target.X+target.Y)
}

func findBestLocation(m *util.Map) (*util.Position, int, []int) {
	asteroidPositions := m.GetAsteroids()
	var bestPos *util.Position
	maxDA := 0 // Max detectable asteroids
	dACounts := make([]int, len(asteroidPositions))

	for i, pos := range asteroidPositions {
		num := numDetectableAstroids(m, asteroidPositions, pos.X, pos.Y)
		if num > maxDA {
			maxDA = num
			bestPos = pos
		}
		dACounts[i] = num
	}
	return bestPos, maxDA, dACounts
}

func numDetectableAstroids(m *util.Map, positions []*util.Position, x, y int) int {

	detected := make(map[string]bool)
	for _, pos := range positions {
		if pos.X == x && pos.Y == y {
			continue
		}
		if _, ok := detected[fmt.Sprintf("%d,%d", pos.X, pos.Y)]; ok {
			continue
		}

		fX, fY := util.Delta(x, pos.X, y, pos.Y)
		curX := x
		curY := y
		for {
			curX += fX
			curY += fY
			if m.InBounds(curX, curY) {
				if m.M[curY][curX] == '#' {
					detected[fmt.Sprintf("%d,%d", curX, curY)] = true
					break
				}
			} else {
				break
			}
		}
	}

	a := len(detected)
	return a
}

func vaporizeAsteroids(m *util.Map, x, y int) *util.Position {
	m.M[y][x] = 'X'
	l := laser.New(m, x, y)
	wantedKills := util.Min(200, len(m.GetAsteroids()))
	asteroidCount := 0
	var target *util.Position
	for asteroidCount < wantedKills {
		target = l.Shoot(m)
		if target != nil {
			/* fmt.Println()
			m.Draw() */
			asteroidCount++
		}
		l.Rotate()
	}

	return target
}
