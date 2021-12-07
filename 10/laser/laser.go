package laser

import (
	"aoc10/util"
	"fmt"
	"math"
	"sort"
)

type StationaryLaser struct {
	rotAngles   []*util.Position
	pos         *util.Position
	curRotIndex int
}

func New(m *util.Map, x, y int) *StationaryLaser {
	pos := &util.Position{
		X: x,
		Y: y,
	}
	rotAngles, rotStartIndex := calcRotAngles(m, pos)
	return &StationaryLaser{
		pos:         pos,
		rotAngles:   rotAngles,
		curRotIndex: rotStartIndex,
	}
}

func (l *StationaryLaser) Shoot(m *util.Map) *util.Position {
	curAngle := l.rotAngles[l.curRotIndex]

	x := l.pos.X
	y := l.pos.Y
	for {
		x += curAngle.X
		y += curAngle.Y

		if m.InBounds(x, y) {
			if m.M[y][x] == '#' {
				m.M[y][x] = '.'
				return &util.Position{X: x, Y: y}
			}
		} else {
			break
		}
	}
	return nil
}

func (l *StationaryLaser) Rotate() {
	l.curRotIndex = (l.curRotIndex + 1) % len(l.rotAngles)
}

func calcRotAngles(m *util.Map, pos *util.Position) ([]*util.Position, int) {
	// Calculate the angle between each asteroid from the laser, sort by angle, and filter out duplicate pairs of dx, dy
	type object struct {
		Pos   *util.Position
		Angle float64
	}
	asteroids := make([]*object, 0)
	for _, asteroid := range m.GetAsteroids() {
		asteroids = append(asteroids, &object{
			Pos:   &util.Position{X: asteroid.X, Y: asteroid.Y},
			Angle: math.Atan2(float64(asteroid.Y-pos.Y), float64(asteroid.X-pos.X)),
		})
	}
	sort.SliceStable(asteroids, func(i, j int) bool {
		a := asteroids[i].Angle
		b := asteroids[j].Angle
		return a < b
	})

	// Calculate dx, dy pairs for each asteroid and filter duplicates
	rotations := map[string]bool{}
	deltas := make([]*util.Position, 0)
	startIndex := 0
	for _, asteroid := range asteroids {
		dx, dy := util.Delta(pos.X, asteroid.Pos.X, pos.Y, asteroid.Pos.Y)
		if _, ok := rotations[fmt.Sprintf("%d,%d", dx, dy)]; !ok {
			rotations[fmt.Sprintf("%d,%d", dx, dy)] = true
			deltas = append(deltas, &util.Position{X: dx, Y: dy})

			// Making sure the laser starts by pointing up
			if dx == 0 && dy == -1 {
				startIndex = len(deltas) - 1
			}
		}
	}

	return deltas, startIndex
}
