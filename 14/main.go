package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type oreHandler func(map[string]int, int) map[string]int

func main() {
	// Part 1
	recipeBook, lookUp := parseFromFile("input.txt")
	ingredients := findOreCount(recipeBook, lookUp, nil, "FUEL", 1)
	fmt.Printf("Part 1: %d ORES needed to produce 1 FUEL\n", ingredients["ORE_TOTAL"])

	// Part 2
	fmt.Printf("Part 2: %d FUEL with 1 trillion ores\n", maxAmountOfFuel(recipeBook, lookUp, 1000000000000))

}

type chemical struct {
	Name   string
	Amount int
}

func findOreCount(recipeBook [][]*chemical, lookUp map[string]int, ingredients map[string]int, chemical string, amount int) map[string]int {
	if ingredients == nil {
		ingredients = make(map[string]int, 0)
	}

	if chemical == "ORE" {
		amountNeeded := amount - ingredients["ORE"]
		if amountNeeded > 0 {
			ingredients["ORE"] += amountNeeded
			ingredients["ORE_TOTAL"] += amountNeeded
		}
		return ingredients
	}

	// Check if necassary ingredients is available
	pageNumber := lookUp[chemical]
	neededIngredients := recipeBook[pageNumber][1:]
	batchSize := recipeBook[pageNumber][0].Amount
	batchesNeeded := 1
	if amount > batchSize {
		batchesNeeded = amount / batchSize
		if amount%batchSize != 0 {
			batchesNeeded++
		}
	}

	// Get all the neceassary ingredients
	hasAllIngredients := false
	for !hasAllIngredients {
		for _, ingredient := range neededIngredients {
			amountNeeded := ingredient.Amount*batchesNeeded - ingredients[ingredient.Name]
			if amountNeeded > 0 {
				ingredients = findOreCount(recipeBook, lookUp, ingredients, ingredient.Name, amountNeeded)
			}
		}

		// Check if all ingredients is stored
		hasAllIngredients = true
		for _, ingredient := range neededIngredients {
			amountNeeded := ingredient.Amount*batchesNeeded - ingredients[ingredient.Name]
			if amountNeeded > 0 {
				hasAllIngredients = false
			}
		}
	}

	// Should now have all the necassary ingredients, time to produce the chemical
	for _, ingredient := range neededIngredients {
		if ingredients[ingredient.Name] < ingredient.Amount*batchesNeeded {
			log.Fatalf("Something is wrong! Have to little of %s to producde %s. Have %d, need %d\n", ingredient.Name, chemical, ingredients[ingredient.Name], ingredient.Amount*batchesNeeded)
		}
		ingredients[ingredient.Name] -= ingredient.Amount * batchesNeeded
	}

	ingredients[chemical] += batchesNeeded * batchSize
	ingredients[chemical+"_TOTAL"] += batchesNeeded * batchSize
	return ingredients
}

func maxAmountOfFuel(recipeBook [][]*chemical, lookUp map[string]int, oreLimit int) int {

	intervals := []int{1000000, 100000, 10000, 1000, 500, 100, 1}
	intervalIndex := 0
	fuelAmount := intervals[0]
	prevFuelAmount := 0
	var u, l bool = false, false
	for {
		prevFuelAmount = fuelAmount

		// Produce fuel equal to the count of "fuelAmount"
		ingredients := findOreCount(recipeBook, lookUp, nil, "FUEL", fuelAmount)

		// Check if there is ore left, if so increase the fuelAmount and check if it is possible to produce more
		if ingredients["ORE_TOTAL"] < oreLimit {
			fuelAmount += intervals[intervalIndex]

			if u == true {
				l = true
			}

		} else {
			// The fuelAmount exceeded the ore limit, decrease the number of fuel to produce
			fuelAmount -= intervals[intervalIndex]
			u = true
		}

		// If the fuelAmount has changed direction twice, decrease the interval step
		if u && l {
			u = false
			l = false
			if intervalIndex < len(intervals)-1 {
				intervalIndex++
			}

			// If the prevFuelAmount is only one step away from the current fuel amount,
			// "fuelAmount -1" is the maximum number of fuel that can be produced
			if prevFuelAmount-fuelAmount == -1 || prevFuelAmount-fuelAmount == 1 {
				break
			}
		}

	}

	return fuelAmount - 1
}

func parseFromFile(fileName string) ([][]*chemical, map[string]int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reactions := make([][]*chemical, 0)
	recipeIndexes := map[string]int{}
	r, _ := regexp.Compile("[0-9]+ [A-Z]+")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chemicals := make([]*chemical, 0)
		line := scanner.Text()
		matches := r.FindAllStringSubmatch(line, -1)
		for i := len(matches) - 1; i >= 0; i-- {
			s := strings.Split(matches[i][0], " ")
			count, err := strconv.Atoi(s[0])
			if err != nil {
				log.Fatalf("Invalid number: %\nErr: %+v", s[0], err)
			}
			chemicals = append(chemicals, &chemical{
				Name:   s[1],
				Amount: count,
			})

		}
		reactions = append(reactions, chemicals)
		recipeIndexes[chemicals[0].Name] = len(reactions) - 1
	}
	return reactions, recipeIndexes
}
