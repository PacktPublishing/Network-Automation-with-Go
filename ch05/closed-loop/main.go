package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/driver/network"
	cfg "github.com/scrapli/scrapligocfg"
	"gopkg.in/yaml.v2"
)

type Router struct {
	Hostname  string `yaml:"hostname"`
	Platform  string `yaml:"platform"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	StrictKey bool   `yaml:"strictkey"`
	Conn      *network.Driver
}

type Inventory struct {
	Routers []Router `yaml:"router"`
}

type DeviceInfo struct {
	Device    string
	Output    string
	Timestamp time.Time
}

type Service struct {
	Name     string
	Port     string
	AF       string
	Insecure bool
	CLI      string
	Config   string
}

func (r Router) getConfig() (c DeviceInfo, err error) {
	rs, err := r.Conn.SendCommand("show run")
	if err != nil {
		return c, fmt.Errorf("failed to send 'show run' for %s: %w", r.Hostname, err)
	}
	c = DeviceInfo{
		Device:    r.Hostname,
		Output:    rs.Result,
		Timestamp: time.Now(),
	}
	return c, nil
}

func (r Router) getOper(s Service) (o DeviceInfo, err error) {
	rs, err := r.Conn.SendCommand(s.CLI)
	if err != nil {
		return o, fmt.Errorf("failed to send %s for %s: %w", s.CLI, r.Hostname, err)
	}
	o = DeviceInfo{
		Device:    r.Hostname,
		Output:    rs.Result,
		Timestamp: time.Now(),
	}
	return o, nil
}

func (r Router) sendConfig(conf string) error {
	c, err := cfg.NewCfg(r.Conn, r.Platform)
	if err != nil {
		return fmt.Errorf("failed create config driver for %s: %w", r.Hostname, err)
	}
	err = c.Prepare()
	if err != nil {
		return fmt.Errorf("failed to prepare config for %s: %w", r.Hostname, err)
	}
	_, err = c.LoadConfig(conf, false)
	if err != nil {
		return fmt.Errorf("failed to load config for %s: %w", r.Hostname, err)
	}
	_, err = c.CommitConfig()
	if err != nil {
		return fmt.Errorf("failed to commit the config for %s: %w", r.Hostname, err)
	}
	return nil
}

func (c DeviceInfo) save() error {
	layout := "01-02-2006_15-04_EST"

	f, err := os.Create("backups/" + c.Device + "_" + c.Timestamp.Format(layout) + ".cfg")
	if err != nil {
		return fmt.Errorf("failed to create 'show run' file for %s: %w", c.Device, err)
	}
	defer f.Close()

	_, err = io.WriteString(f, c.Output)
	if err != nil {
		return fmt.Errorf("failed to create write 'show run' for %s: %w", c.Device, err)
	}
	return f.Sync()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (s Service) genConfig() (string, error) {
	if s.Config != "" {
		return s.Config, nil
	}

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

func (s Service) parseOper(input string) (Service, error) {
	if s.Name != "grpc" {
		return Service{}, fmt.Errorf("service %s not supported", s.Name)
	}
	n := regexp.MustCompile(`transport.*`)
	n.Find([]byte(input))
	nsl := strings.Split(string(n.Find([]byte(input))), " ")

	p := regexp.MustCompile(`listening-port.*`)
	p.Find([]byte(input))
	psl := strings.Split(string(p.Find([]byte(input))), " ")

	i := regexp.MustCompile(`TLS.*`)
	i.Find([]byte(input))
	isl := strings.Split(string(i.Find([]byte(input))), " ")
	noTLS := true
	if isl[len(isl)-1] == "enabled" {
		noTLS = false
	}

	a := regexp.MustCompile(`access-family.*`)
	a.Find([]byte(input))
	asl := strings.Split(string(a.Find([]byte(input))), " ")
	af := "ipv6"
	if asl[len(asl)-1] == "tcp4" {
		af = "ipv4"
	}

	o := Service{
		Name:     nsl[len(nsl)-1],
		Port:     psl[len(psl)-1],
		AF:       af,
		Insecure: noTLS,
		CLI:      "show grpc status",
	}
	return o, nil
}

func main() {
	////////////////////////////////
	// Read input data
	////////////////////////////////
	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var inv Inventory
	err = d.Decode(&inv)
	check(err)
	iosxr := inv.Routers[0]

	////////////////////////////////
	// Service Definition
	////////////////////////////////
	intent := Service{
		Name:     "grpc",
		Port:     "57777",
		AF:       "ipv4",
		Insecure: false,
		CLI:      "show grpc status",
	}
	intentHash, err := hashstructure.Hash(intent,
		hashstructure.FormatV2, nil)
	check(err)

	////////////////////////////////////////
	// Open connection to the network device
	///////////////////////////////////////
	conn, err := platform.NewPlatform(
		iosxr.Platform,
		iosxr.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(iosxr.Username),
		options.WithAuthPassword(iosxr.Password),
		options.WithSSHConfigFile("ssh_config"),
	)
	check(err)

	driver, err := conn.GetNetworkDriver()
	check(err)

	iosxr.Conn = driver
	err = driver.Open()
	check(err)
	defer driver.Close()

	////////////////////////////////
	// Backup config
	////////////////////////////////
	config, err := iosxr.getConfig()
	check(err)

	err = config.save()
	check(err)

	////////////////////////////////
	// Continuous/Enforcement loop
	////////////////////////////////
	fmt.Println("Entering to continuous loop ====>")
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	go func() {
		for ; true; <-ticker.C {
			h, m, s := time.Now().Clock()
			fmt.Printf(" Loop at %02d:%02d:%02d\n", h, m, s)
			////////////////////////////////
			// Get Operational Data
			////////////////////////////////
			opr, err := iosxr.getOper(intent)
			check(err)
			// fmt.Println(opr.Output)

			////////////////////////////////
			// Parse Operational Data
			////////////////////////////////
			parsed, err := intent.parseOper(opr.Output)
			check(err)
			fmt.Printf("  Operational state from device:\n   service: %v\n   addr-family: %v\n   port: %v\n   TLS: %v\n\n",
				parsed.Name, parsed.AF, parsed.Port, !(parsed.Insecure))

			oprHash, err := hashstructure.Hash(parsed, hashstructure.FormatV2, nil)
			check(err)

			////////////////////////////////
			// Validate State
			////////////////////////////////
			if oprHash == intentHash {
				continue
			}

			////////////////////////////////
			// Generate config
			////////////////////////////////
			conf, err := intent.genConfig()
			check(err)
			intent.Config = conf

			////////////////////////////////
			// Send config if necessary
			////////////////////////////////
			fmt.Printf("Configuring device ====>\n\n")
			err = iosxr.sendConfig(conf)
			check(err)
		}
	}()

	time.Sleep(time.Second * 200)
	fmt.Println("End of the program")
}
