package main

import (
	"fmt"
	"strings"

	"aoc6/heap"
	"aoc6/reader"
)

// Note to self: Diksjtra is not the best solution! Think about the problem more next time!

func main() {
	// TEST
	m := buildOrbitMap([]string{
		"COM)B",
		"B)C",
		"C)D",
		"D)E",
		"E)F",
		"B)G",
		"G)H",
		"D)I",
		"E)J",
		"J)K",
		"K)L",
	})
	total, _ := countOrbits(m)
	fmt.Printf("Total: %d\n", total)
	fmt.Println(shortestPath(m, "K", "I"))

	// INPUT
	lines := reader.ReadLines()
	m = buildOrbitMap(lines)
	total, _ = countOrbits(m)
	fmt.Printf("\nTotal: %d\n", total)
	shortestPath := shortestPath(m, m["YOU"].Nodes[0].Id, m["SAN"].Nodes[0].Id)
	fmt.Printf("Shortest Path Length: %d\n%+v\n", len(shortestPath), shortestPath)
}

// CountOrbits solves part 1
func countOrbits(m map[string]*heap.Object) (total int, counts map[string]int) {
	counts = make(map[string]int, 0)
	totalOrbits := 0

	for key, obj := range m {
		orbitCount := 0
		for len(obj.Nodes) > 0 && obj.Nodes[0] != nil {
			orbitCount++
			obj = obj.Nodes[0]
		}
		counts[key] = orbitCount
		totalOrbits += orbitCount
	}
	return totalOrbits, counts
}

func shortestPath(m map[string]*heap.Object, start, end string) []string {
	// Set distances to infinity
	for _, obj := range m {
		obj.Distance = int(^uint(0) >> 1)
	}

	h := heap.NewHeap()
	s := m[start]
	s.Distance = 0
	h.Push(s)
	target := m[end]
	pathFound := false

	for h.Len() > 0 {
		curNode := h.Pop().(*heap.Object)
		if curNode == target {
			pathFound = true
			break
		}

		// Check the neighbours
		for _, neighbour := range curNode.Nodes {

			newCost := curNode.Distance + 1
			if newCost < neighbour.Distance {
				neighbour.Distance = newCost
				neighbour.Parent = curNode
				h.Push(neighbour)
			}
		}
	}

	if !pathFound {
		return nil
	}

	// Find path
	path := []string{}
	curNode := target
	for curNode.Parent != nil {
		path = append(path, curNode.Id)
		curNode = curNode.Parent
	}

	return path
}

func buildOrbitMap(m []string) map[string]*heap.Object {
	objects := map[string]*heap.Object{}

	for _, orbit := range m {
		// Read and create objects
		objectPair := strings.Split(orbit, ")")
		for _, obj := range objectPair {
			if _, ok := objects[obj]; !ok {
				objects[obj] = &heap.Object{
					Id:    obj,
					Nodes: make([]*heap.Object, 0),
				}
			}
		}

		// Set orbit relation
		o1 := objects[objectPair[0]]
		o2 := objects[objectPair[1]]
		o2.Nodes = append([]*heap.Object{o1}, o2.Nodes...) // Makes sure the orbitet object is always first
		o1.Nodes = append(o1.Nodes, o2)
	}
	com := objects["COM"]
	com.Nodes = append([]*heap.Object{nil}, com.Nodes...)
	return objects
}
