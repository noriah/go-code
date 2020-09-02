package fizzbuzz

import "fmt"

// Channel does fizzbuzz with channels
func Channel() {

	var num, div3, div5 int

	var numChannel = make(chan int)
	var doneChannel = make(chan bool)

	go channelPrinter(numChannel, doneChannel)

	for num = 1; num <= 100; num++ {

		div3 = (1 >> (num % 3))
		div5 = (1 >> (num % 5))

		numChannel <- (div3 << 7) | (div5 << 8) | num
	}

	doneChannel <- true
}

// This function is very unsafe
func channelPrinter(numCh <-chan int, doneCh <-chan bool) {

	const word = "FizzBuzz"

	var value int

	for {
		select {
		case value = <-numCh:
			if value > 100 {
				value >>= 7
				fmt.Println(word[4-((value&0x01)<<2) : 4+((value&0x02)<<1)])
				continue
			}

			fmt.Println(value)

		case <-doneCh:
			return
		}
	}
}
