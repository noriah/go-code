package fizzbuzz

import "fmt"

// Branchless does FizzBuzz by shifting, without branches
func Branchless() {
	// declare our vars
	var num, start, end int

	// Our string of FizzBuzz
	word := "FizzBuzz"

	// Keys for printf
	key := []string{"%[2]s\n", "%[1]d\n"}

	// loop through the numbers 1 to 100
	// assign some variables used every loop
	for num = 1; num <= 100; num++ {
		// If we look at the string FizzBuzz as an array of charactersm we want
		// either the first half, the second half, or all of it.
		// F I Z Z B U Z Z
		// 0 1 2 3 4 5 6 7

		// Calculate the start
		// If we set the start value to 0x18 (00011000), then shift it right by
		// num mod3, it will be 0 when num is divisble by 3, and either
		// 12 (00001100) or 6 (00000110) when not divisble. We can mask the result
		// with 0x04 (00000100) and now start can only be 0 or 4
		start = (0x18 >> (num % 3)) & 0x04

		// Calculate the end
		// First shift the value 0x01 (00000001) right by num mod5. This will result in
		// 1 when divisble by 5 and 0 otherwise. Then shift 0x04 (00000100) to the left
		// By either 1 or 0 to get either 8 or 4 (because we need to be not-inclusive on
		// the array)
		end = 0x04 << (0x01 >> (num % 5))

		// Print it out
		// if we divide start by end (integer math only), we can select the right key
		fmt.Printf(key[start/end], num, word[start:end])
	}
}
