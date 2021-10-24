package main

import (
	"fmt"
)

type Router struct {
	Hostname  string
	Platform  string
	Username  string
	Password  string
	StrictKey bool
}

type Inventory struct {
	Routers []Router
}

func main() {
	var r1 Router
	r1.Hostname = "router1.example.com"

	r2 := new(Router)
	r2.Hostname = "router2.example.com"

	r3 := Router{
		Hostname:  "router3.example.com",
		Platform:  "cisco_iosxr",
		Username:  "user",
		Password:  "secret",
		StrictKey: false,
	}

	fmt.Printf("r1: %v\n", r1.Hostname)
	fmt.Printf("r2: %v\n", r2.Hostname)
	fmt.Printf("r3: %v\n", r3.Hostname)

	inv := Inventory{
		Routers: []Router{r1, *r2, r3},
	}

	fmt.Printf("Inventory: %+v\n", inv)
}
