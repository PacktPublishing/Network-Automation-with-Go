package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	input := []byte{0xE9, 0xFB, 0x03}
	asn, n := binary.Uvarint(input)
	if n != len(input) {
		fmt.Println("Varint did not consume all the input")
	}
	fmt.Printf("ASN: %v \n", asn)
}
