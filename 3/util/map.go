package util

import (
	"fmt"
	"strconv"
	"strings"
)

const interPoint = -2

type Point struct {
	x int
	y int
}

func (p *Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

type Map struct {
	m         map[string]int // the map
	Start     Point
	stepCount map[string]map[int]int
	wireCount int
}

func NewMap() *Map {
	m := &Map{
		m:         make(map[string]int, 0),
		stepCount: make(map[string]map[int]int),
		Start: Point{
			x: 0,
			y: 0,
		},
		wireCount: 0,
	}
	m.m[m.Start.String()] = -1
	return m
}

func (m *Map) Intersections() []*Point {
	intersections := make([]*Point, 0)
	for key, value := range m.m {
		if value == interPoint {
			points := strings.Split(key, ",")
			x, _ := strconv.Atoi(points[0])
			y, _ := strconv.Atoi(points[1])
			intersections = append(intersections, &Point{x: x, y: y})
		}
	}
	return intersections
}

func (m *Map) AddWire(wire []string) {
	curPoint := Point{
		x: m.Start.x,
		y: m.Start.y,
	}
	m.wireCount++
	stepCount := 0

	for _, w := range wire {
		dir := rune(w[0])
		steps, _ := strconv.Atoi(w[1:])
		for steps > 0 {
			steps--

			switch dir {
			case 'U':
				curPoint.y++
			case 'R':
				curPoint.x++
			case 'D':
				curPoint.y--
			case 'L':
				curPoint.x--
			}

			// Initialize stepcount
			if _, ok := m.stepCount[curPoint.String()]; !ok {
				m.stepCount[curPoint.String()] = make(map[int]int)
			}
			stepCount++

			_, exists := m.m[curPoint.String()]
			if !exists {
				m.m[curPoint.String()] = m.wireCount
				m.stepCount[curPoint.String()][m.wireCount] = stepCount

			} else {
				_, hasBeenHere := m.stepCount[curPoint.String()][m.wireCount]
				if !hasBeenHere {
					// The wire has not been here before
					m.m[curPoint.String()] = interPoint
					m.stepCount[curPoint.String()][m.wireCount] = stepCount
				}
			}
		}
	}
}

func (m *Map) ClosestInterDist(s *Point) (int, *Point) {
	points := m.Intersections()

	minDist := int(^uint(0) >> 1) // Max int
	var closestPoint *Point

	for _, p := range points {
		dist := abs(p.x-s.x) + abs(p.y-s.y)
		if dist < minDist {
			minDist = dist
			closestPoint = p
		}
	}
	return minDist, closestPoint
}

func (m *Map) InterPointLowestSteps() (int, *Point) {
	points := m.Intersections()

	minSteps := int(^uint(0) >> 1) // Max int
	var bestPoint *Point
	for _, p := range points {
		totalSteps := 0
		for _, steps := range m.stepCount[p.String()] {
			totalSteps += steps
		}
		if totalSteps < minSteps {
			minSteps = totalSteps
			bestPoint = p
		}
	}
	return minSteps, bestPoint
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
