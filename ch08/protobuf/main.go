package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
	"github.com/Network-Automation-with-Go/ch08/protobuf/pb" 
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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
	f := "routers.data"
	
	// Read the existing routers
	in, err := os.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found.  Creating new file.\n", f)
		}  
		log.Fatalln("Error reading file:", err)
	}

	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	check(err)


	routers := &pb.Router{}
	// Load file contents in routers
	if err := proto.Unmarshal(in, routers); err != nil {
		log.Fatalln("Failed to parse the routers file:", err)
	}

	router := &pb.Router{}

	routers.Router = append(routers.Router, router)

	// Write the new router back to disk.
	out, err := proto.Marshal(routers)
	if err != nil {
		log.Fatalln("Failed to encode router:", err)
	}
	if err := os.WriteFile(f, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}

}