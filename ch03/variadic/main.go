package main

import (
	"fmt"
	"strings"
)

func printOctets(octets ...string) {
	fmt.Println(strings.Join(octets, "."))
}

func main() {
	// prints "127.1"
	printOctets("127", "1")

	ip := []string{"192", "0", "2", "1"}

	// prints "192.0.2.1"
	printOctets(ip...)
}
