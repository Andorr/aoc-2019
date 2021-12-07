package painter

import (
	"fmt"
	"log"
	"math"
)

const Black = 0
const White = 1
const TurnLeft = 0
const TurnRight = 1

type Painter struct {
	x, y   int // Position
	dx, dy int // Rotation
}

func NewPainter() *Painter {
	return &Painter{
		dx: 0, dy: -1,
		x: 0, y: 0,
	}
}

func (p *Painter) Paint(ship *Ship, color int) {
	ship.Paint(p.x, p.y, color)
}

func (p *Painter) CurrentColor(ship *Ship) int {

	color, ok := ship.panels[p.position()]
	if !ok {
		ship.Paint(p.x, p.y, Black)
		return Black
	}
	return color
}

func (p *Painter) Rotate(rotDir int) {

	var dir float64
	if rotDir == TurnLeft {
		dir = -1.0
	} else if rotDir == TurnRight {
		dir = 1.0
	} else {
		log.Fatalf("Invalid rotation direction %d", rotDir)
	}

	angle := (math.Pi / 2.0) * dir
	dx := float64(p.dx)*math.Cos(angle) - float64(p.dy)*math.Sin(angle)
	dy := float64(p.dx)*math.Sin(angle) + float64(p.dy)*math.Cos(angle)

	p.dx = int(dx)
	p.dy = int(dy)
}

func (p *Painter) Move(ship *Ship) {
	p.x += p.dx
	p.y += p.dy
}

func (p *Painter) String() string {
	if p.dx == 0 && p.dy == -1 {
		return "^"
	} else if p.dx == 0 && p.dy == 1 {
		return "v"
	} else if p.dx == 1 && p.dy == 0 {
		return ">"
	} else if p.dx == -1 && p.dy == 0 {
		return "<"
	} else {
		return "P"
	}
}

func (p *Painter) position() string {
	return posToString(p.x, p.y)
}

type Ship struct {
	panels map[string]int
	minX   int
	maxX   int
	minY   int
	maxY   int

	paintCount map[string]int
}

func NewShip() *Ship {
	return &Ship{
		panels:     make(map[string]int, 0),
		paintCount: make(map[string]int, 0),
	}
}

func (s *Ship) Paint(x, y, color int) {
	s.panels[posToString(x, y)] = color
	count := s.paintCount[posToString(x, y)]
	s.paintCount[posToString(x, y)] = count + 1

	if y < s.minY {
		s.minY = y
	} else if y > s.maxY {
		s.maxY = y
	}
	if x < s.minX {
		s.minX = x
	} else if x > s.maxX {
		s.maxX = x
	}
}

func (s *Ship) Draw(p *Painter) {
	for y := s.minY; y <= s.maxY; y++ {
		for x := s.minX; x <= s.maxX; x++ {
			if x == p.x && y == p.y {
				fmt.Print(p.String())
				continue
			}
			color, ok := s.panels[posToString(x, y)]
			if !ok {
				color = Black
			}
			if color == Black {
				fmt.Print(".")
			} else if color == White {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func (s *Ship) NumPanelsPainted() int {
	count := 0
	for _, v := range s.paintCount {
		if v >= 1 {
			count++
		}
	}
	return count
}

func posToString(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}
