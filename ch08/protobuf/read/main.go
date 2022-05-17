package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PacktPublishing/Network-Automation-with-Go/ch08/protobuf/pb"
	"google.golang.org/protobuf/proto"
)

type Model struct {
	Uplinks  []Link `yaml:"uplinks"`
	Peers    []Peer `yaml:"peers"`
	ASN      int    `yaml:"asn"`
	Loopback Addr   `yaml:"loopback"`
}

type Link struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

type Peer struct {
	IP  string `yaml:"ip"`
	ASN int    `yaml:"asn"`
}

type Addr struct {
	IP string `yaml:"ip"`
}

// Main reads the static routers list and writes out to a file.
func main() {
	//File to save data
	f := "../router.data"
	router := &pb.Router{}

	// Read protobuf file with stored routers
	in, err := os.ReadFile(f)
	if err != nil {
		switch os.IsNotExist(err) {
		case true:
			fmt.Printf("%s: File not found.  Creating new file.\n", f)
		default:
			log.Fatalln("Error reading file:", err)
		}
	}

	// Load file contents in router
	if err := proto.Unmarshal(in, router); err != nil {
		log.Fatalln("Failed to parse the routers file:", err)
	}

	fmt.Printf("%v\n", router)
}
