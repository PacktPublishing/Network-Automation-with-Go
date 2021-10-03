package main

import (
	"fmt"
	"github.com/PacktPublishing/Network-Automation-with-Go/ch01/ping"
)

func main() {
	s := ping.Send()
	fmt.Println(s)
}
