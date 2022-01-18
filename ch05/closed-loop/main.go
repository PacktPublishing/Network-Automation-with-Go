package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
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

type Service struct {
	Name     string
	Port     string
	AF       string
	Insecure bool
}

func (r Router) getConfig() (c Config, err error) {
	d, err := core.NewCoreDriver(
		r.Hostname,
		r.Platform,
		base.WithAuthStrictKey(r.StrictKey),
		base.WithAuthUsername(r.Username),
		base.WithAuthPassword(r.Password),
		base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return c, fmt.Errorf("failed to create driver for %s: %w", r.Hostname, err)
	}

	err = d.Open()
	if err != nil {
		return c, fmt.Errorf("failed to open driver for %s: %w", r.Hostname, err)
	}
	defer d.Close()

	rs, err := d.SendCommand("show run")
	if err != nil {
		return c, fmt.Errorf("failed to send 'show run' for %s: %w", r.Hostname, err)
	}

	c = Config{
		Device:    r.Hostname,
		Running:   rs.Result,
		Timestamp: time.Now(),
	}

	return c, nil
}

func (c Config) save() error {
	layout := "01-02-2006_15-04_EST"

	f, err := os.Create("backups/" + c.Device + "_" + c.Timestamp.Format(layout) + ".cfg")
	if err != nil {
		return fmt.Errorf("failed to create 'show run' file for %s: %w", c.Device, err)
	}
	defer f.Close()

	_, err = io.WriteString(f, c.Running)
	if err != nil {
		return fmt.Errorf("failed to create write 'show run' for %s: %w", c.Device, err)
	}
	return f.Sync()
}

func check(err error){
	if err != nil {
		panic(err)
	}
}

func (s Service) genConfig() (string, error) {
	base, err := os.ReadFile(s.Name + ".template")
	if err != nil {
		return "", fmt.Errorf("failed to read template file for %s: %w", s.Name, err)
	}

	t, err := template.New("service").Parse(string(base))
	if err != nil {
		return "", fmt.Errorf("failed to parse template for %s: %w", s.Name, err)
	}
	var b strings.Builder
	err = t.Execute(&b, s)
	if err != nil {
		return "", fmt.Errorf("failed to parse template for %s: %w", s.Name, err)
	}
	return b.String(), nil
}

func main() {
	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var inv Inventory
	err = d.Decode(&inv)
	check(err)
	iosxr := inv.Routers[0]

	// Backup config
	config, err := iosxr.getConfig()
	check(err)
	
	err = config.save()
	check(err)

	// Generate config
	svc := Service{"grpc", "57777", "ipv4", false}
	cfg, err := svc.genConfig()
	check(err)
	fmt.Println(cfg)

	
}