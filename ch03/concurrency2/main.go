package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"gopkg.in/yaml.v2"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}

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

func getVersion(r Router, wg *sync.WaitGroup) {
	defer wg.Done()
	p, err := platform.NewPlatform(
		r.Platform,
		r.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(r.Username),
		options.WithAuthPassword(r.Password),
		options.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		fmt.Printf("failed to create platform for %s: %+v\n", r.Hostname, err)
		return
	}

	d, err := p.GetNetworkDriver()
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

	fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
				r.Hostname, parsedOut[0]["HARDWARE"],
				parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

}

func main() {
	// To time this process
	defer timeTrack(time.Now())

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

	var wg sync.WaitGroup

	for _, v := range inv.Routers {
		wg.Add(1)
		go getVersion(v, &wg)
	}
	wg.Wait()
}
