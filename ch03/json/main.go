package main

import (
	"fmt"
	"os"

	"encoding/json"
)

type Router struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	ASN      uint16 `json:"asn"`
}

type Inventory struct {
	Routers []Router `json:"router"`
}

func main() {
	file, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	src := json.NewDecoder(file)

	var inv Inventory
	// Decode reads the next JSON-encoded value from its
	// input and stores it in the value pointed to by v.
	err = src.Decode(&inv)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", inv)
}
