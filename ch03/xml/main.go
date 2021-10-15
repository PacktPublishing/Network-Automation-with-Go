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
	d := xml.NewDecoder(file)

	var inv Inventory

	// Decode XML from the source and store it in the value pointed to by inv.
	err = d.Decode(&inv)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", inv)
}
