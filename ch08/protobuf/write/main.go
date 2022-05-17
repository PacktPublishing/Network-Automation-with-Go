package main

import (
	"log"
	"os"

	"github.com/PacktPublishing/Network-Automation-with-Go/ch08/protobuf/pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
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
	f := "../router.data"
	router := &pb.Router{}

	// Read data with new input router
	src, err := os.Open("../input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	check(err)

	uplinks := input.Uplinks
	for _, uplink := range uplinks {
		router.Uplinks = append(router.GetUplinks(),
			&pb.Uplink{
				Name:   uplink.Name,
				Prefix: uplink.Prefix,
			},
		)
	}

	peers := input.Peers
	for _, peer := range peers {
		router.Peers = append(router.GetPeers(),
			&pb.Peer{
				Ip:  peer.IP,
				Asn: int32(peer.ASN),
			},
		)
	}

	router.Asn = int32(input.ASN)

	router.Loopback = &pb.Addr{Ip: input.Loopback.IP}

	// Write the new router back to disk.
	out, err := proto.Marshal(router)
	if err != nil {
		log.Fatalln("Failed to encode router:", err)
	}
	if err := os.WriteFile(f, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
	// fmt.Printf("%X", out)
}
