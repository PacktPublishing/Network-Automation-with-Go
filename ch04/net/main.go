package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ipv4 := net.ParseIP("192.0.2.1")
	fmt.Println(ipv4)

	ipv6 := net.ParseIP("FC02:F00D::1")
	fmt.Println(ipv6)

	// prints false
	fmt.Println(ipv4.IsPrivate())
	// prints true
	fmt.Println(ipv6.IsPrivate())

	// This mask corresponds to a /31 subnet for IPv4.
	// prints [11111111 11111111 11111111 11111110]
	fmt.Printf("%b\n", net.CIDRMask(31, 32))

	// This mask corresponds to a /64 subnet for IPv6.
	// prints ffffffffffffffff0000000000000000
	fmt.Printf("%s\n", net.CIDRMask(64, 128))

	ipv4Addr, ipv4Net, err := net.ParseCIDR("192.0.2.1/24")
	if err != nil {
		log.Fatal(err)
	}
	// prints 192.0.2.1
	fmt.Println(ipv4Addr)
	// prints 192.0.2.0/24
	fmt.Println(ipv4Net)

	ipv6Addr, ipv6Net, err := net.ParseCIDR("2001:db8:a0b:12f0::1/32")
	if err != nil {
		log.Fatal(err)
	}
	// prints 2001:db8:a0b:12f0::1
	fmt.Println(ipv6Addr)
	// prints 2001:db8::/32
	fmt.Println(ipv6Net)
}
