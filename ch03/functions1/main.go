package main

import (
	"fmt"
	"strings"
)

func generateName(base string, suffix string) string {
	parts := []string{base, suffix}
	return strings.Join(parts, "-")
}

func processDevice(getName func(string, string) string, ip string) {
	base := "device"
	name := getName(base, ip)
	fmt.Println(name)
}

func main() {
	s := generateName("device", "01")

	// prints "device-01"
	fmt.Println(s)

	// prints "device-192.0.2.1"
	processDevice(generateName, "192.0.2.1")

}
