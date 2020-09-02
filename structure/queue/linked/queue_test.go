package linked

func generateIntArray(size int) []int {
	var ret = make([]int, size)
	for i := 0; i < size; i++ {
		ret[i] = size - i
	}
	return ret
}

// TODO: Write tests!
