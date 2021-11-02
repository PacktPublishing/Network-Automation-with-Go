package main

import (
	"fmt"
)



func main() {
	// Header length (measured in 32-bit words) is 5
	var headerWords uint8 = 5

	// Header length in bytes is 20
	headerLen := headerWords * 32 / 8

	// Build a slice of 20 bytes to store the entire TCP header
	b := make([]byte, headerLen)

	// Shift header words bits to the left to fit the Header Length field of the TCP header
    s := headerWords << 4

	// Perform OR operation on byte 13 and store the new value
	b[13] = b[13] | s

	// Print the 13 byte of the TCP header -> [01010000]
	fmt.Printf("%08b\n", b[13])

	// Let's also assume that this is the initial TCP SYN message
	var tcpSyn uint8 = 1

	// SYN flag is the second bit from the right so we shift it by 1 position
	f := tcpSyn << 1
	
	// Perform OR operation on byte 13 and store the new value
	b[14] = b[14] | f

	// Print the 14 byte of the TCP header -> [00000010]
	fmt.Printf("%08b\n", b[14])

	// Full header -> prints [10100000 00000010]
	fmt.Printf("%08b\n", b)

	// We are only interested if a TCP SYN flag has been set
	tcpSynFlag := (b[14] & 0x02) != 0

	// Shift the header length right and drop any low-order bits
	parsedHeaderWords := b[13] >> 4

	// prints "TCP Flag is set: true"
	fmt.Printf("TCP Flag is set: %t\n", tcpSynFlag)

	// prints "TCP header words: 5"
	fmt.Printf("TCP header words: %d\n", parsedHeaderWords)
}
