package main

import (
	"fmt"
	"os"
	"time"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"gopkg.in/yaml.v2"
)

type Router struct {
	Hostname  string `yaml:"hostname"`
	Platform  string `yaml:"platform"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	StrictKey bool   `yaml:"strictkey"`
}

type Inventory struct {
	Routers []Router `yaml:"router"`
}

func getVersion(r Router, out chan data) {
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

	out <- data{
		host: r.Hostname,
		hw: fmt.Sprintf("%s", parsedOut[0]["HARDWARE"]),
		version: fmt.Sprintf("%s", parsedOut[0]["VERSION"]),
		uptime: fmt.Sprintf("%s", parsedOut[0]["UPTIME"]),
	}

}

type data struct{
	host string
	hw string
	version string
	uptime string
}

func main() {
	src, err := os.Open("input.yml")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	d := yaml.NewDecoder(src)

	var inv Inventory
	err = d.Decode(&inv)
	if err != nil {
		panic(err)
	}

	ch := make(chan data)

	for _, v := range inv.Routers {
		go getVersion(v, ch)
	}

	for {
		select {
		case out := <-ch:
			fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
				out.host, out.hw, out.version, out.uptime)
		case <-time.After(5 * time.Second):
			close(ch)
			fmt.Println("Timeout: 5 seconds")
			return
		}
	}
}
