package main

import (
	"fmt"
	"os"
	"strings"

	"encoding/xml"
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
	src, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer src.Close()
	d := json.NewDecoder(src)

	var inv Inventory

	// Decode JSON from the source and store it in the value pointed to by inv.
	err = d.Decode(&inv)
	if err != nil {
		panic(err)
	}

	// Now we can write inv to a destination in a different format.
	var dest strings.Builder

	e := xml.NewEncoder(&dest)
	err = e.Encode(&inv)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", dest.String())

}
