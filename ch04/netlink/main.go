package main

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink/rtnl"
)

func main() {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	links, err := conn.Links()
	if err != nil {
		log.Fatal(err)
	}

	var loopback *net.Interface
	for _, l := range links {

		if l.Name == "lo" {
			loopback = l
			log.Printf("Name: %s, Flags:%s\n", l.Name, l.Flags)
		}
	}

	_ = conn.LinkDown(loopback)
	loopback, _ = conn.LinkByIndex(loopback.Index)
	log.Printf("Name: %s, Flags:%s\n", loopback.Name, loopback.Flags)

	_ = conn.LinkUp(loopback)
	loopback, _ = conn.LinkByIndex(loopback.Index)
	log.Printf("Name: %s, Flags:%s\n", loopback.Name, loopback.Flags)

}
