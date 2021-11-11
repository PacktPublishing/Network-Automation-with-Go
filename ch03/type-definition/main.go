package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// type int, size 8 bytes on a 64-bit machine
	a := -1

	// unsigned 32-bit integer, size 4 bytes
	var b uint32
	b = 4294967295

	// floating point number, size 4 bytes
	var c float32 = 42.1

	fmt.Printf("a: %T, size: %d bytes\n", a, unsafe.Sizeof(a))
	fmt.Printf("b: %T, size: %d bytes\n", b, unsafe.Sizeof(b))
	fmt.Printf("c: %T, size: %d bytes\n", c, unsafe.Sizeof(c))
}
