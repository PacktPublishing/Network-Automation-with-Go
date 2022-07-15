package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	cfg "github.com/scrapli/scrapligocfg"
	"gopkg.in/yaml.v2"
)

const ceosTemplate = `
!
configure
!
ip routing
!
{{- range $uplink := .Uplinks }}
interface {{ $uplink.Name }}
  no switchport
  ip address {{ $uplink.Prefix }}
!
{{- end }}
interface Loopback0
  ip address {{ .Loopback.IP }}/32
!
router bgp {{ .ASN }}
  router-id {{ .Loopback.IP }}
{{- range $peer := .Peers }}  
  neighbor {{ $peer.IP }} remote-as {{ $peer.ASN }}
{{- end }}
  redistribute connected
!
`

type Model struct {
	Uplinks  []Link `yaml:"uplinks"`
	Peers    []Peer `yaml:"peers"`
	ASN      int    `yaml:"asn"`
	Loopback Addr   `yaml:"loopback"`
}

type Link struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

type Peer struct {
	IP  string `yaml:"ip"`
	ASN int    `yaml:"asn"`
}

type Addr struct {
	IP string `yaml:"ip"`
}

func devConfig(in Model) (b bytes.Buffer, err error) {
	t, err := template.New("config").Parse(ceosTemplate)
	if err != nil {
		return b, fmt.Errorf("failed create template: %w", err)
	}

	err = t.Execute(&b, in)
	if err != nil {
		return b, fmt.Errorf("failed create template: %w", err)
	}
	log.Print("Generated config: ", b.String())
	return b, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	hostname := flag.String("device", "clab-netgo-ceos", "Device hostname")
	username := flag.String("username", "admin", "SSH username")
	password := flag.String("password", "admin", "SSH password")
	nos := flag.String("nos", "arista_eos", "Network operating system")
	flag.Parse()

	// read and parse the input file
	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	check(err)

	config, err := devConfig(input)
	check(err)

	conn, err := platform.NewPlatform(
		*nos,
		*hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(*username),
		options.WithAuthPassword(*password),
	)
	check(err)

	driver, err := conn.GetNetworkDriver()
	check(err)

	err = driver.Open()
	check(err)
	defer driver.Close()

	conf, err := cfg.NewCfg(driver, *nos)
	check(err)

	err = conf.Prepare()
	check(err)

	_, err = conf.LoadConfig(config.String(), false)
	check(err)
}
