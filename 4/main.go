package main

import "fmt"

func main() {
	/* 	fmt.Println(getDigits(1124556789))
	   	fmt.Println(isValidPassword(1124556789))
	   	fmt.Println(isValidPassword(111111))
	   	fmt.Println(isValidPassword(223450))
	   	fmt.Println(isValidPassword(123789))
	   	fmt.Println(isValidPassword(112233))
	   	fmt.Println(isValidPassword(123444))
	   	fmt.Println(isValidPassword(111122)) */

	fmt.Println(combinationCount(245318, 765747))
}

func combinationCount(from, to int) int {
	count := 0
	for i := from; i <= to; i++ {
		if isValidPassword(i) {
			count++
		}
	}
	return count
}

func isValidPassword(password int) bool {
	digits := getDigits(password)
	hasDouble := false

	for i := len(digits) - 1; i > 0; i-- {
		if digits[i] > digits[i-1] {
			return false
		}
		if !hasDouble && digits[i] == digits[i-1] {
			if len(digits)-1 > i && i > 1 {
				// Check when i is between the start and the end
				hasDouble = digits[i+1] != digits[i] && digits[i-1] != digits[i-2]
			} else if i == len(digits)-1 {
				// i == len(digits) - 1. Check if the third digits is equal
				hasDouble = digits[i-1] != digits[i-2]
			} else {
				// i == 1, checks the number before it
				hasDouble = digits[i+1] != digits[i]
			}
		}
	}
	return hasDouble
}

func getDigits(value int) []int {
	digits := make([]int, 0)
	for value > 0 {
		digits = append(digits, value%10)
		value /= 10
	}
	return digits
}
