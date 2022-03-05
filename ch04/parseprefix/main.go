package main

import (
	"fmt"
	"net/netip"
)

func main() {
	addr1 := "192.0.2.18"
	addr2 := "198.51.100.3"

	network4 := "192.0.2.0/24"
	pf := netip.MustParsePrefix(network4)
	fmt.Printf("Prefix address: %v, length: %v\n", pf.Addr(), pf.Bits())

	ip1 := netip.MustParseAddr(addr1)
	if pf.Contains(ip1) {
		fmt.Println(addr1, " is in ", network4)
	}

	ip2 := netip.MustParseAddr(addr2)
	if pf.Contains(ip2) {
		fmt.Println(addr2, " is in ", network4)
	}

	addr3 := "2600::"
	addr4 := "2001:db8:F00D::CAFE"

	network6 := "2001:db8::/32"
	pf = netip.MustParsePrefix(network6)
	fmt.Printf("Prefix address: %v, length: %v\n", pf.Addr(), pf.Bits())

	ip3 := netip.MustParseAddr(addr3)
	if pf.Contains(ip3) {
		fmt.Println(addr3, " is in ", network6)
	}

	ip4 := netip.MustParseAddr(addr4)
	if pf.Contains(ip4) {
		fmt.Println(addr4, " is in ", network6)
	}

}
