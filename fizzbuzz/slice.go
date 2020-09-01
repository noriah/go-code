package fizzbuzz

import "fmt"

// Slice does FizzBuzz by slicing a string
func Slice() {

	// declare our vars
	var num, d3, d5, start, end uint8

	// Our string of FizzBuzz
	const word = "FizzBuzz"

	// loop through the numbers 1 to 100
	for num = 1; num <= 100; num++ {
		// if num is divisble by 3, num mod 3 will return 0, otherwise
		// we shift the value to 0
		d3 = (1 >> (num % 3))

		// if num is divisble by 5, num mod 5 will return 0, otherwise
		// we shift the value to 0
		d5 = (1 >> (num % 5))

		// OR the divisbility results, do we have a flag set?
		if d3|d5 > 0 {

			// If we look at the string FizzBuzz as an array of characters we want
			// either the first half, the second half, or all of it.
			// F I Z Z B U Z Z
			// 0 1 2 3 4 5 6 7

			// Calculate the beginning
			start = 4 - (d3 * 4)

			// Calculate the end
			end = 4 + (d5 * 4)

			// Print it out
			fmt.Println(word[start:end])

			// go to the beginning of this loop
			continue
		}

		// If we don't have a flag, just print the number
		fmt.Println(num)
	}
}
