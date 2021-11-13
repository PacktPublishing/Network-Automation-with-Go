package main

import (
	"fmt"
)

func suffixGenerator() func() string {
	i := 0
	return func() string {
		i++
		return fmt.Sprintf("%02d", i)
	}
}

func main() {
	generator1 := suffixGenerator()

	// prints "device-01"
	fmt.Printf("%s-%s\n", "device", generator1())

	// prints "device-02"
	fmt.Printf("%s-%s\n", "device", generator1())

	generator2 := suffixGenerator()

	// prints "device-01"
	fmt.Printf("%s-%s\n", "device", generator2())
}
