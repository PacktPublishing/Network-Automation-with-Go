package main

import (
	"fmt"
	"os"

	"encoding/json"
)

func main() {
	file, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	d := json.NewDecoder(file)

	var empty map[string]interface{}

	err = d.Decode(&empty)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", empty)

	for _, r := range empty["router"].([]interface{}) {
		fmt.Printf("%v\n", r)
	}

}
