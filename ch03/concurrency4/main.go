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

var m sync.RWMutex = sync.RWMutex{}

func getVersion(r Router, out chan data, wg *sync.WaitGroup, isAlive map[string]bool) {
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
		m.Lock()
		isAlive[r.Hostname] = false
		m.Unlock()
		return
	}

	m.Lock()
	isAlive[r.Hostname] = true
	m.Unlock()

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

func printer(in chan data) {
	for out := range in {
		fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
			out.host, out.hw, out.version, out.uptime)
	}
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

	ch := make(chan data)

	isAlive := make(map[string]bool)

	go printer(ch)

	var wg sync.WaitGroup
	for _, v := range inv.Routers {
		wg.Add(1)
		go getVersion(v, ch, &wg, isAlive)
	}
	wg.Wait()

	m.RLock()
	for name, v := range isAlive {
		fmt.Printf("Router %s is alive: %t\n", name, v)
	}
	m.RUnlock()
	
	close(ch)
}