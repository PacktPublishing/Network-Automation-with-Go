package main

import (
	"fmt"
	"net"
	"net/netip"
)

func main() {
	ipv4 := net.ParseIP("192.0.2.1")
	IPv4s, ok := netip.AddrFromSlice(ipv4)
	if !ok {
		fmt.Println("couldn't parse the v4 address")
		return
	}
	fmt.Println(IPv4s.String())
	fmt.Println(IPv4s.Unmap().Is4())

	ipv6 := net.ParseIP("FC02:F00D::1")
	IPv6s, ok := netip.AddrFromSlice(ipv6)
	if !ok {
		fmt.Println("couldn't parse the v6 address")
		return
	}
	fmt.Println(IPv6s.String())
	fmt.Println(IPv6s.Unmap().Is6())
}
