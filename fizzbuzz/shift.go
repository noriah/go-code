package main

import "fmt"

func fizzBuzzShift() {

	// declare our vars
	var num, d3, d5, start, end uint8

	// Our string of FizzBuzz
	fizzBuzz := "FizzBuzz"

	// loop through the numbers 1 to 100
	// assign some variables used every loop
	for num = 1; num <= 100; num++ {
		// if num is divisble by 3, num mod 3 will return 0
		d3 = (1 >> (num % 3))

		// if num is divisble by 3, num mod 3 will return 0
		d5 = (1 >> (num % 5))

		// OR the divisbility results, do we have a flag set?
		if d3|d5 != 0 {
			// Calculate the beginning
			start = (1 - d3) * 4

			// Calculate the end
			end = (1 + d5) * 4

			// Print it out
			fmt.Println(fizzBuzz[start:end])

			// go to the beginning of this loop
			continue
		}

		// If we don't have a flag, just print the number
		fmt.Println(num)
	}
}
