package main

import "fmt"

func doFill(toFind int, print bool) {
	newPrimes := make([]int, toFind)
	var curIdx, newIdx int
	var isP bool

	curIdx = 2

	for {
		isP = true
		for i := 0; i < newIdx; i++ {
			prime := newPrimes[i]

			if curIdx%prime == 0 {
				isP = false
				break
			}
		}

		if isP {
			newPrimes[newIdx] = curIdx
			if print {
				fmt.Printf("%d - %d\n", newIdx, curIdx)
			}
			newIdx++
			if newIdx >= toFind {
				break
			}
		}

		curIdx++
	}

	fmt.Println(newPrimes)
}
