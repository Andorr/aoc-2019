package util

import "math"

type Position struct {
	X int
	Y int
}

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Delta(fromX, toX, fromY, toY int) (int, int) {
	fX := toX - fromX
	fY := toY - fromY
	absFX := Abs(fX)
	absFY := Abs(fY)
	lowest := absFX
	if absFY > 0 && absFY < lowest {
		lowest = absFY
	} else if lowest == 0 {
		lowest = absFY
	}

	// Shorten
	for i := 2; i <= lowest; i++ {
		for fX%i == 0 && fY%i == 0 {
			fX /= i
			fY /= i
		}
	}
	return fX, fY
}

func Magnitude(fromX, toX, fromY, toY int) float64 {
	return math.Sqrt(math.Pow((float64(toX-fromX)), 2) + math.Pow((float64(toY-fromY)), 2))
}
