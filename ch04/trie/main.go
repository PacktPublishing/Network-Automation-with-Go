package main

import (
	"fmt"
	"log"
	"net"

	"github.com/yl2chen/cidranger"
)

func main() {
	// instantiate NewPCTrieRanger
	ranger := cidranger.NewPCTrieRanger()

	IPs := []string{
		"100.64.0.0/16",
		"127.0.0.0/8",
		"172.16.0.0/16",
		"192.0.2.0/24",
		"192.0.2.0/24",
		"192.0.2.0/25",
		"192.0.2.127/25",
	}

	checkIP := "127.0.0.1"
	netIP := "192.0.2.18"

	fmt.Println("Adding prefixes")
	// Add prefixes to ranger
	for idx, prefix := range IPs {
		ipv4Addr, ipv4Net, err := net.ParseCIDR(prefix)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v: %s\n", idx, ipv4Addr)
		ranger.Insert(cidranger.NewBasicRangerEntry(*ipv4Net))
	}
	fmt.Println("--------------------------------------------")

	// Check if the IP address is within the list of IP ranges.
	ok, err := ranger.Contains(net.ParseIP(checkIP))
	if err != nil {
		log.Fatal("Error running ranger.Contains", err.Error())
	}
	fmt.Printf("Does the range contain %s?: %v\n", checkIP, ok)

	// Request the list of networks that containing the IP address
	nets, err := ranger.ContainingNetworks(net.ParseIP(netIP))
	if err != nil {
		log.Fatal("Error running ranger.ContainingNetworks", err.Error())
	}

	fmt.Printf("\n\nPrint networks that contain IP address %s ->\n", netIP)
	for _, e := range nets {
		n := e.Network()
		fmt.Println("\t", n.String())
	}
}
