package main

import (
	"fmt"
	"github.com/PacktPublishing/Network-Automation-with-Go/ch02/ping"
)

func main() {
	s := ping.Send()
	fmt.Println(s)
}
