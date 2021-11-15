package main

import (
	"fmt"
	"os"
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

func getVersion(r Router, out chan map[string]interface{}, last bool) {
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

	out <- parsedOut[0]

	if last {
		close(out)
	}

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

	for i, v := range inv.Routers {
		go getVersion(v, ch, i == len(inv.Routers)-1)
	}

	printer(ch)

}
