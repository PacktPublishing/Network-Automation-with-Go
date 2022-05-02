package main

import (
	"fmt"
)

func main() {
	empty := []string{}
	words := []string{"zero", "one", "two", "three", "four", "five", "six"}
	three := make([]string, 3)

	fmt.Printf("empty: length: %d, capacity: %d, %v\n", len(empty), cap(empty), empty)
	fmt.Printf("words: length: %d, capacity: %d, %v\n", len(words), cap(words), words)
	fmt.Printf("three: length: %d, capacity: %d, %v\n", len(three), cap(three), three)

	// Creating a slice
	mySlice := words[1:3]
	fmt.Printf("mySlice: length: %d, capacity: %d, %v\n", len(mySlice), cap(mySlice), mySlice)

	mySlice = append(mySlice, "seven")
	fmt.Printf("mySlice: length: %d, capacity: %d, %v\n", len(mySlice), cap(mySlice), mySlice)

	mySlice = append(mySlice, "eight", "nine", "ten", "eleven")
	fmt.Printf("mySlice: length: %d, capacity: %d, %v\n", len(mySlice), cap(mySlice), mySlice)

	
}
