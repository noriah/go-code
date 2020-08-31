package prime

// GenerateCount finds X number of primes by trial division, starting at 3
// checking the value against an array of primes we fill with previously
// found primes
//
// NOTE: This is NOT a good way to find big prime numbers. As we find primes,
// we increase the max number of iterations of the inner loop for each
// following candidate.
//
// Since primes we find will not be in our array, and won't be divisble by
// any previous prime, we have to go through all at-time known primes to
// find a single new prime.
//
// Time complexity: O(n**2)
// Space complexity: O(n)
func GenerateCount(count int) []int {
	if count < 1 {
		return nil
	}

	primes := make([]int, count)
	// Insert 2 because we already know its prime
	primes[0] = 2

	var candidate, idx, jdx, divisor int
	var isPrime bool

	// var iterations uint64

	// Candidate starts at 3 and idx at 1 beacuse we want to skip 2 as we already know it
	for candidate, idx = 3, 1; idx < count; candidate += 2 {
		// reset the flag
		isPrime = true

		// For all the primes found so far
		for jdx = 1; jdx < idx; jdx++ {
			divisor = primes[jdx]
			// Compare the candidate against it with MOD
			if candidate%divisor == 0 {
				// Oh no. not a prime. set a flag and break the loop
				isPrime = false
				break
			}
		}

		// iterations += uint64(jdx)

		// Did we find a prime
		if isPrime {
			// Add it to the list
			primes[idx] = candidate
			// increment the number of found primes
			idx++
		}
	}

	// fmt.Println(count, iterations)

	return primes
}
