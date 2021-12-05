package main

import (
	"fmt"
	"net"
	"sort"

	"github.com/c-robinson/iplib"
)

func main() {
	IP := net.ParseIP("192.0.2.1")
	nextIP := iplib.NextIP(IP)
	incrIP := iplib.IncrementIPBy(nextIP, 19)

	// prints 20
	fmt.Println(iplib.DeltaIP(IP, incrIP))
	// prints -1
	fmt.Println(iplib.CompareIPs(IP, incrIP))

	iplist := []net.IP{incrIP, nextIP, IP}
	// prints [192.0.2.21 192.0.2.2 192.0.2.1]
	fmt.Println(iplist)

	sort.Sort(iplib.ByIP(iplist))
	// prints [192.0.2.1 192.0.2.2 192.0.2.21]
	fmt.Println(iplist)
	fmt.Println("--------------------------------")  

	n4 := iplib.NewNet4(net.ParseIP("198.51.100.0"), 24)
	fmt.Printf("\nIP Net: %s\n", n4.String())
	fmt.Println("Total IP addresses: ", n4.Count())            // 254
	// Enumerate(size, offset int): generates an array of all 
	// usable addresses in Net up to the given size starting at the 
	// given offset
	fmt.Println("First three IPs: ", n4.Enumerate(3, 0))
	fmt.Println("First IP: ", n4.FirstAddress())
	fmt.Println("Last IP: ", n4.LastAddress())
	fmt.Println(n4.Subnet(0))
	fmt.Println("--------------------------------")   

	n6 := iplib.NewNet6(net.ParseIP("2001:db8::"), 32, 0)
	fmt.Printf("\nIP Net: %s\n", n6.String())
	fmt.Printf("Total IP addresses: %d\n", n6.Count())
	fmt.Println("First two IPs, after address #1024: ", n6.Enumerate(2, 1024))
}
