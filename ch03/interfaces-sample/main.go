package main

import (
	"fmt"
	"os"
)

type CiscoIOS struct {
	Hostname string
	Platform string
}

func (r CiscoIOS) getUptime() int {
	return 1
}

type CiscoNXOS struct {
	Hostname string
	Platform string
	ACI      bool
}

func (s CiscoNXOS) getUptime() int {
	return 2
}

type NetworkDevice interface {
	getUptime() int
}

func LastToReboot(r1, r2 NetworkDevice) bool {
	return r1.getUptime() < r2.getUptime()
}

func main() {
	ios := CiscoIOS{}
	nexus := CiscoNXOS{}

	if LastToReboot(ios, nexus) {
		fmt.Println("IOS-XE has been running for less time, so it was the last to be rebooted")
		os.Exit(0)
	}
	fmt.Println("NXOS was the last one to reboot")
}