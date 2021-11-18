package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
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

func getVersion(r Router, out chan map[string]interface{}, wg *sync.WaitGroup, isAlive map[string]bool) {
	defer wg.Done()
	
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

	parsedOut[0]["HOSTNAME"] = r.Hostname
	out <- parsedOut[0]
}

func printer(in chan map[string]interface{}) {
	for out := range in {
		fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
			out["HOSTNAME"], out["HARDWARE"],
			out["VERSION"], out["UPTIME"])
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

	ch := make(chan map[string]interface{})

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
