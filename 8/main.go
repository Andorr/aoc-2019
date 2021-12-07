package main

import (
	"aoc8/imgutil"
	"fmt"
	"io/ioutil"
)

func main() {
	inputBytes, _ := ioutil.ReadFile("input.txt")

	// Part 1
	fmt.Println(layerFewest0s(string(inputBytes), 25, 6))

	// Part 2
	imgViewer := imgutil.Parse(string(inputBytes), 25, 6)
	fmt.Println(imgViewer.Draw(50))
}

// LayerFewest0s get the layer number of the layer with fewest zeros and returns product = (num 1s)*(num 2s) in that layer #Part1
func layerFewest0s(img string, w, h int) (layer, minNumZeros, product int) {
	layerLength := w * h
	numOfLayers := len(img) / layerLength
	minNumZeros = int(^uint(0) >> 1)

	for i := 0; i < numOfLayers; i++ {
		numZeros := 0
		numOnes := 0
		numTwos := 0
		for j := 0; j < layerLength; j++ {
			value := img[i*layerLength+j]
			if value == '0' {
				numZeros++
			} else if value == '1' {
				numOnes++
			} else if value == '2' {
				numTwos++
			}
		}
		if numZeros < minNumZeros {
			layer = i + 1
			minNumZeros = numZeros
			product = numOnes * numTwos
		}
	}
	return
}
