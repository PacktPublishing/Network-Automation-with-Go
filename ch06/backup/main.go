package main

import (
	"fmt"
	"io"
	"os"
	"sync"
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

type Config struct {
	Device    string
	Running   string
	Timestamp time.Time
}

func getVersion(r Router, out chan map[string]interface{}, wg *sync.WaitGroup) {
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
		fmt.Printf("failed to send 'show version' for %s: %+v\n", r.Hostname, err)
		return
	}

	parsedOut, err := rs.TextFsmParse(r.Platform + "_show_version.textfsm")
	if err != nil {
		fmt.Printf("failed to parse 'show version' for %s: %+v\n", r.Hostname, err)
		return
	}

	parsedOut[0]["HOSTNAME"] = r.Hostname
	out <- parsedOut[0]

}

func getConfig(r Router, out chan Config, wg *sync.WaitGroup) {
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

	rs, err := d.SendCommand("show run")
	if err != nil {
		fmt.Printf("failed to send 'show run' for %s: %+v\n", r.Hostname, err)
		return
	}

	output := Config{
		Device:    r.Hostname,
		Running:   rs.Result,
		Timestamp: time.Now(),
	}

	out <- output

}

func printer(in chan map[string]interface{}) {
	for out := range in {
		fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
			out["HOSTNAME"], out["HARDWARE"],
			out["VERSION"], out["UPTIME"])
	}
}

func save(in chan Config) {
	layout := "01-02-2006_15-04_EST"

	for out := range in {
		f, err := os.Create("backups/" + out.Device + "_" + out.Timestamp.Format(layout) + ".cfg")
		if err != nil {
			fmt.Printf("failed to create 'show run' file for %s: %+v\n", out.Device, err)
			return
		}

		_, err = io.WriteString(f, out.Running)
		if err != nil {
			fmt.Printf("failed to create write 'show run' for %s: %+v\n", out.Device, err)
		}
		f.Sync()
		f.Close()
	}
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

	ch := make(chan map[string]interface{})

	go printer(ch)

	var wg sync.WaitGroup
	for _, v := range inv.Routers {
		wg.Add(1)
		go getVersion(v, ch, &wg)
	}

	wg.Wait()
	close(ch)

	ch2 := make(chan Config)
	go save(ch2)

	var wg2 sync.WaitGroup
	for _, v := range inv.Routers {
		wg2.Add(1)
		go getConfig(v, ch2, &wg2)
	}

	wg2.Wait()
	close(ch2)
}
