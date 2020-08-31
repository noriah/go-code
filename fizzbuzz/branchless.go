package fizzbuzz

import "fmt"

// Branchless does FizzBuzz by shifting, without branches
func Branchless() {
	// declare our vars
	var num, start, end int

	// Our string of FizzBuzz
	fizzBuzz := "FizzBuzz"

	// Keys for printf
	key := []string{"%[2]s\n", "%[1]d\n"}

	// loop through the numbers 1 to 100
	// assign some variables used every loop
	for num = 1; num <= 100; num++ {
		// Calculate the start
		start = (0x18 >> (num % 3)) & 0x04

		// Calculate the end
		end = 0x04 << (0x01 >> (num % 5))

		// Print it out
		fmt.Printf(key[start/end], num, fizzBuzz[start:end])
	}
}
