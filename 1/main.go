package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputFile = "input.txt"

func main() {
	fmt.Println(fuelAdditional(12))
	fmt.Println(fuelAdditional(14))
	fmt.Println(fuelAdditional(1969))
	fmt.Println(fuelAdditional(100756))

	fmt.Printf("Total: %d\n", calculateFuel(fuel))
	fmt.Printf("Total: %d\n", calculateFuel(fuelAdditional))
}

func fuel(mass int) int {
	return int(mass/3.0) - 2
}

func fuelAdditional(mass int) int {
	f := fuel(mass)
	if f <= 0 {
		return 0
	}
	return f + fuelAdditional(f)
}

func calculateFuel(f func(m int) int) int {
	sum := 0

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		mass, err := strconv.Atoi(input)
		if err != nil {
			log.Fatalf("Unable to read int: %s\nErr: %+v", input, err)
		}
		sum += f(mass)
	}
	return sum
}
