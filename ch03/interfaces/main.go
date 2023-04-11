package main

import (
	"fmt"
	//"strconv"
	"strings"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
)

type NetworkDevice struct {
	Hostname  string
	Platform  string
	Username  string
	Password  string
	StrictKey bool
}

type Inventory struct {
	Devices []NetworkDevice
}

func getUptime(r NetworkDevice) (string, error) {
	p, err := platform.NewPlatform(
		r.Platform,
		r.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(r.Username),
		options.WithAuthPassword(r.Password),
		options.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create platform for %s: %w", r.Hostname, err)
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		return "", fmt.Errorf("failed to create driver for %s: %w", r.Hostname, err)
	}

	err = d.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open driver for %s: %w", r.Hostname, err)
	}
	defer d.Close()

	rs, err := d.SendCommand("show version")
	if err != nil {
		return "", fmt.Errorf("failed to send command for %s: %w", r.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse(r.Platform + "_show_version.textfsm")
	if err != nil {
		return "", fmt.Errorf("failed to parse command for %s: %w", r.Hostname, err)
	}

	symbol := "s"

	switch r.Platform {
	case "cisco_iosxe":
		symbol = "s"
	case "cisco_nxos":
		symbol = "(s)"
	default:
	}

	// Split uptime string in Hour, Day, Minute, etc. on a slice
	slc := strings.Split(parsedOut[0]["UPTIME"].(string), ",")

	m := make(map[string]string)
	for _, item := range slc {
		// Divide the number from the item
		spl := strings.Split(strings.TrimSpace(item), " ")

		// Remove 's', '(s)', tec.
		nos := strings.TrimRight(spl[1], symbol)

		m[nos] = spl[0]
	}

	// day, _ := strconv.Atoi(m["day"])
	// hour, _ := strconv.Atoi(m["hour"])
	// min, _ := strconv.Atoi(m["minute"])
	// un := 24 * 60 * day + 60 * hour + min
	// fmt.Printf("INT: %v\n", un)

	uptime := fmt.Sprintf("Day: %v, Hour: %v, Minute: %v\n",
		m["day"], m["hour"], m["minute"])

	return fmt.Sprintf("Hostname: %s\nUptime: %s\n", r.Hostname, uptime), nil
}

func main() {
	ios := NetworkDevice{
		Hostname:  "sandbox-iosxe-latest-1.cisco.com",
		Platform:  "cisco_iosxe",
		Username:  "admin",
		Password:  "C1sco12345",
		StrictKey: false,
	}
	nxos := NetworkDevice{
		Hostname:  "sandbox-nxos-1.cisco.com",
		Platform:  "cisco_nxos",
		Username:  "admin",
		Password:  "Admin_1234!",
		StrictKey: false,
	}

	inv := Inventory{
		Devices: []NetworkDevice{ios, nxos},
	}

	for _, v := range inv.Devices {
		s, err := getUptime(v)
		if err != nil {
			fmt.Printf("[ERROR]: %s\n", err.Error())
		}
		fmt.Println(s)
	}
}
