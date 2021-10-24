package main

import (
	"fmt"
)

func main() {
	dc := make(map[string]string)

	dc["spine"] = "192.168.100.1"

	ip := dc["spine"]
	ip, exists := dc["spine"]

	if exists {
		fmt.Println(ip)
	}

	inv := map[string]string{
		"router1.example.com": "192.0.2.1/32",
		"router2.example.com": "198.51.100.1/32",
	}

	fmt.Printf("inventory: length: %d, %v\n", len(inv), inv)
	
	delete(inv, "router1.example.com")
	
	fmt.Printf("inventory: length: %d, %v\n", len(inv), inv)	
}
