package main

import (
	"fmt"
)

type Device struct {
	name string
}

func mutate(input Device) {
	input.name += "-suffix"
}

func main() {
	d := Device{name: "myname"}
	mutate(d)

	// prints "myname"
	fmt.Println(d.name)
}
