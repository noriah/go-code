package fizzbuzz

import (
	"fmt"
)

// Table prints a table with binary representation
func Table() {

	// define our variables
	var num, div3, div5, div15, flag int
	var star rune

	var fmtStr string = "%[1]c %3[2]d - %08[2]b |  %2[3]d - %04[3]b | %2[4]d - %04[4]b | %2[5]d - %04[5]b\n"

	for num = 1; num <= 100; num++ {

		// num MOD 3 will return the remainder after dividing
		// by 3 as much as possible.
		// if no remainder, num is divisble by 3
		div3 = num % 3

		// num MOD 5 will return the remainder after dividing
		// by 5 as much as possible.
		// if no remainder, num is divisble by 5
		div5 = num % 5

		// num MOD 15 will return the remainder after dividing
		// by 15 as much as possible.
		// if no remainder, num is divisble by 15
		div15 = num % 15

		// If div3 or div5 are not 0, we right shift their value
		// 0x01 by any amount, making it 0. OR the two shifts and
		// we now have a flag representing divisble by 3 or 5.
		//
		// Alternatively, add the shifts to get number of conditions met
		flag = (1 >> div3) | (1 >> div5)

		// space is 0x20 (32)
		// star is 0x2a
		// difference 0x0a
		// multiply the flag by the difference and add that to space
		// This will return space when there is no divide, and star otherwise
		star = rune(0x20 + flag*0x0a)

		// Print the row
		fmt.Printf(fmtStr, star, num, div3, div5, div15)
	}
}
