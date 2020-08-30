package main

import "fmt"

func fizzBuzzTable() {

	s := ""
	j := 0
	for i, x, y, z, q := 1, 1, 1, 1, 1; i <= 100; i++ {
		// x, y, z = i%15, i%3, i%5
		q = (q + i)

		x = (i % 15)
		y = (i % 3)
		z = (i % 5)

		if y == 0 || z == 0 {
			s = "*"
		} else {
			s = " "
		}

		j = (1 >> y) | (2 >> z)

		fmt.Printf("%s %3d | %08b ||| %2d | %04b ||| %2d | %04b ||| %2d | %08b\n", s, i, i, y, y, z, z, x, x)

		if false {
			fmt.Println(i, "|", x, y, z, "|", y+z, (x % 3), (x % 5), "|", (i%15)%3, "|", fmt.Sprintf("%04b", (i%15)), "|", i, j)
		}
		q = q % 15
	}
}
