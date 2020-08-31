package fizzbuzz

import "fmt"

// Shift2 does FizzBuzz by shifting bits but with less steps
func Shift2() {

	// declare our vars
	var num, div3, div5, start, end int

	// Our string of FizzBuzz
	fizzBuzz := "FizzBuzz"

	// loop through the numbers 1 to 100
	// assign some variables used every loop
	for num = 1; num <= 100; num++ {
		// if num is divisble by 3, num mod 3 will return 0
		div3 = num % 3

		// if num is divisble by 5, num mod 5 will return 0
		div5 = num % 5

		// Calculate the start
		start = (0x18 >> div3) & 0x04

		// Calculate the end
		end = 0x04 << (0x01 >> div5)

		// The above calculations clamp the start to 0 or 4, and
		// clamp the end to 4 or 8. when both are 4, the slice will be empty
		if start == end {

			// If we don't have a flag, just print the number
			fmt.Println(num)

			// go to the beginning of this loop
			continue
		}

		// Print it out
		fmt.Println(fizzBuzz[start:end])
	}
}
