package main

import (
	"fmt"
	"os"

	"encoding/xml"
)

type Router struct {
	Hostname string `xml:"hostname"`
	IP       string `xml:"ip"`
	ASN      uint16 `xml:"asn"`
}

type Inventory struct {
	Routers []Router `xml:"router"`
}

func main() {
	file, err := os.Open("input.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	src := xml.NewDecoder(file)

	var inv Inventory
	// Decode reads the next XML-encoded value from its
	// input and stores it in the value pointed to by v.
	err = src.Decode(&inv)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", inv)
}
