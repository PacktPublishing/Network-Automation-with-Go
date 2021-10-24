package main

import (
	"fmt"
)

func main() {
	hostnames := [2]string{"router1.example.com", "router2.example.com"}

	ips := [3]string{
		"192.0.2.1/32",
		"198.51.100.1/32",
		"203.0.113.1/32",
	}

	// Prints router2.example.com
	fmt.Println(hostnames[1])

	// Prints 203.0.113.1/32
	fmt.Println(ips[2])

	// ipv4 is [0000 0000, 0000 0000, 0000 0000, 0000 0000]
	var ipAddr [4]byte

	// ipv4 is [1111 1111, 0000 0000, 0000 0000, 0000 0001]
	var localhost = [4]byte{127, 0, 0, 1}

	// prints 4
	fmt.Println(len(localhost))

	// prints [1111111 0 0 1]
	fmt.Printf("%b\n", localhost)

	// prints false
	fmt.Println(ipAddr == localhost)
}
