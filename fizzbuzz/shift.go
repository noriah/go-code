package fizzbuzz

import "fmt"

// Shift does FizzBuzz by shifting bits
func Shift() {

	// declare our vars
	var num, div3, div5, flags, start, end int

	// Our string of FizzBuzz
	fizzBuzz := "FizzBuzz"

	// loop through the numbers 1 to 100
	// assign some variables used every loop
	for num = 1; num <= 100; num++ {
		// if num is divisble by 3, num mod 3 will return 0
		div3 = num % 3

		// if num is divisble by 5, num mod 5 will return 0
		div5 = num % 5

		flags = ((1 >> div3) << 1) | (1 >> div5)

		// OR the divisbility results, do we have a flag set?
		if flags == 0 {

			// If we don't have a flag, just print the number
			fmt.Println(num)

			// go to the beginning of this loop
			continue
		}

		// Calculate the beginning
		start = (0x0c >> flags) & 0x04

		// Calculate the end
		end = (0x50 >> flags) & 0x0c

		// Print it out
		fmt.Println(fizzBuzz[start:end])
	}
}

// Shift2 does FizzBuzz by shifting bits but its different somehow?
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
