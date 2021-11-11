package main

import (
	"fmt"
)

func main() {
	condition := true

	if condition {
		fmt.Printf("Type: %T, Value: %t \n", condition, condition)
	}
}
