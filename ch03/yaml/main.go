package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Router struct {
	Hostname string `yaml:"hostname"`
	IP       string `yaml:"ip"`
	ASN      uint16 `yaml:"asn"`
}

type Inventory struct {
	Routers []Router `yaml:"router"`
}

func main() {
	file, err := os.Open("input.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	d := yaml.NewDecoder(file)

	var inv Inventory

	// Decode YAML from the source and store it in the value pointed to by inv.
	err = d.Decode(&inv)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", inv)
}
