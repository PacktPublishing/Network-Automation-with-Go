package main

import (
	"fmt"
	"strings"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
)

type Router struct {
	Hostname  string
	Platform  string
	Username  string
	Password  string
	StrictKey bool
}

type Inventory struct {
	Routers []Router
}

func getUptime(r Router) {
	d, err := core.NewCoreDriver(
		r.Hostname,
		r.Platform,
		base.WithAuthStrictKey(r.StrictKey),
		base.WithAuthUsername(r.Username),
		base.WithAuthPassword(r.Password),
		base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		fmt.Printf("failed to create driver for %s: %+v\n", r.Hostname, err)
		return
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver for %s: %+v\n", r.Hostname, err)
		return
	}
	defer d.Close()

	rs, err := d.SendCommand("show version")
	if err != nil {
		fmt.Printf("failed to send command for %s: %+v\n", r.Hostname, err)
		return
	}

	parsedOut, err := rs.TextFsmParse(r.Platform + "_show_version.textfsm")
	if err != nil {
		fmt.Printf("failed to parse command for %s: %+v\n", r.Hostname, err)
		return
	}
	uptime := "N/A"

	switch r.Platform {
	case "cisco_iosxe":
		uptime = parseIOS(parsedOut[0]["UPTIME"].(string))
	case "cisco_nxos":
		uptime = parseNXOS(parsedOut[0]["UPTIME"].(string))
	default:
	}

	fmt.Printf("Hostname: %s\nUptime: %s\n\n", r.Hostname, uptime)
}

func parseIOS(s string) (u string) {
	slc := strings.Split(s, ",")
	m := make(map[string]string)
	for _, item := range slc {
		spl := strings.Split(strings.TrimSpace(item), " ")

		nos := strings.TrimRight(spl[1], "s")

		m[nos] = spl[0]

	}
	return fmt.Sprintf("Day: %v, Hour: %v, Minute: %v\n",
		m["day"], m["hour"], m["minute"])

}

func parseNXOS(s string) (u string) {
	slc := strings.Split(s, ",")
	m := make(map[string]string)
	for _, item := range slc {
		spl := strings.Split(strings.TrimSpace(item), " ")

		nos := strings.TrimRight(spl[1], "(s)")

		m[nos] = spl[0]

	}
	return fmt.Sprintf("Day: %v, Hour: %v, Minute: %v\n",
		m["day"], m["hour"], m["minute"])

}

func main() {
	ios := Router{
		Hostname:  "sandbox-iosxe-latest-1.cisco.com",
		Platform:  "cisco_iosxe",
		Username:  "developer",
		Password:  "C1sco12345",
		StrictKey: false,
	}
	nxos := Router{
		Hostname:  "sandbox-nxos-1.cisco.com",
		Platform:  "cisco_nxos",
		Username:  "admin",
		Password:  "Admin_1234!",
		StrictKey: false,
	}

	inv := Inventory{
		Routers: []Router{ios, nxos},
	}

	for _, v := range inv.Routers {
		getUptime(v)
	}

}
