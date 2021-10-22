package main

import (
	"fmt"
)

func main() {
	s1 := "Net"

	s2 := `work`

	if s1 != s2 {
		fmt.Println(s1 + s2 + " Automation")
	}
}
