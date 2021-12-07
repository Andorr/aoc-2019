package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

func main() {
	objects := initObjects("input.txt")
	fmt.Println(simulate(objects, 200))
	objects = initObjects("input.txt")
	fmt.Println(simulateUntilRepeat(objects))
}

type Object struct {
	x, y, z          int
	velX, velY, velZ int
}

func (o *Object) potentialEnergy() int {
	return abs(o.x) + abs(o.y) + abs(o.z)
}

func (o *Object) kineticEnergy() int {
	return abs(o.velX) + abs(o.velY) + abs(o.velZ)
}

func (o *Object) move() {
	o.x += o.velX
	o.y += o.velY
	o.z += o.velZ
}

func simulate(objects []*Object, steps int) int {

	for i := 0; i < steps; i++ {
		for _, obj := range objects {
			for _, objB := range objects {
				if obj == objB {
					continue
				}
				applyGravity(obj, objB)
			}
		}

		for _, obj := range objects {
			obj.move()
		}
	}

	// Calculate total energy
	totalEnergy := 0
	for _, obj := range objects {
		totalEnergy += obj.potentialEnergy() * obj.kineticEnergy()
	}
	return totalEnergy // Total energy
}

func simulateUntilRepeat(objects []*Object) int {

	var freqX, freqY, freqZ int                    // The frequency of the positions in each dimensions
	var startX, startY, startZ string = "", "", "" // The start state of the position of each dimension
	for _, obj := range objects {
		startX += fmt.Sprintf("%d,", obj.x)
		startY += fmt.Sprintf("%d,", obj.y)
		startZ += fmt.Sprintf("%d,", obj.z)
	}

	// Start on step 2 because the inital state is step 1
	for step := 2; ; step++ {
		var stateX, stateY, stateZ string = "", "", ""

		// Apply gravity
		for _, obj := range objects {
			for _, objB := range objects {
				if obj == objB {
					continue
				}
				applyGravity(obj, objB)
			}
		}

		// Move objects and initialize state
		for _, obj := range objects {
			obj.move()
			stateX += fmt.Sprintf("%d,", obj.x)
			stateY += fmt.Sprintf("%d,", obj.y)
			stateZ += fmt.Sprintf("%d,", obj.z)
		}

		// Check if the state has been repeated for each dimension
		if stateX == startX && freqX == 0 {
			freqX = step
		}
		if stateY == startY && freqY == 0 {
			freqY = step
		}
		if stateZ == startZ && freqZ == 0 {
			freqZ = step
		}

		// If each dimension in positions has been repeated
		if freqX != 0 && freqY != 0 && freqZ != 0 {
			break
		}
	}

	// Add the max frequency until the step is divisible by all the frequencies
	maxStepFreq := max(max(freqX, freqY), freqZ)
	steps := maxStepFreq
	for steps%freqX != 0 || steps%freqY != 0 || steps%freqZ != 0 {
		steps += maxStepFreq
	}

	return steps
}

func applyGravity(a, b *Object) {
	if a.x > b.x {
		a.velX--
	} else if a.x < b.x {
		a.velX++
	}
	if a.y > b.y {
		a.velY--
	} else if a.y < b.y {
		a.velY++
	}
	if a.z > b.z {
		a.velZ--
	} else if a.z < b.z {
		a.velZ++
	}
}

func initObjects(fileName string) []*Object {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	r, err := regexp.Compile("[-0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	matches := r.FindAllStringSubmatch(string(bytes), 12)

	values := make([]*Object, 0)
	for i := 0; i < len(matches); i += 3 {
		x, _ := strconv.Atoi(matches[i][0])
		y, _ := strconv.Atoi(matches[i+1][0])
		z, _ := strconv.Atoi(matches[i+2][0])

		values = append(values, &Object{
			x: x, y: y, z: z,
		})
	}

	return values
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
