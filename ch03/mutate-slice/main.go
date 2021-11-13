package main

import (
	"fmt"
)

func mutateV(input []string) {
	input[0] = "r03"
	input = append(input, "r04")
}

func mutateP(input *[]string) {
	(*input)[0] = "r03"
	*input = append(*input, "r04")
}

func main() {
	d1 := []string{"r01", "r02"}
	mutateV(d1)

	// prints "[r03 r02]"
	fmt.Printf("%v\n", d1)

	d2 := []string{"r01", "r02"}
	mutateP(&d2)

	// prints "[r03 r02 r04]"
	fmt.Printf("%v\n", d2)
}
