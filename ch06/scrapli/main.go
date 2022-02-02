package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
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

type Input struct {
	Uplinks []struct {
		Name   string `yaml:"name"`
		Prefix string `yaml:"prefix"`
	} `yaml:"uplinks"`
	Loopback struct {
		IP string `yaml:"ip"`
	} `yaml:"loopback"`
	ASN   int `yaml:"asn"`
	Peers []struct {
		IP  string `yaml:"ip"`
		ASN int    `yaml:"asn"`
	} `yaml:"peers"`
}

func main() {
	deviceName := flag.String("device", "clab-netgo-ceos", "Device Hostname")
	username := flag.String("username", "admin", "SSH Username")
	password := flag.String("password", "admin", "SSH password")
	flag.Parse()

	// read and parse the input file
	src, err := os.Open("input.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Input
	err = d.Decode(&input)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("config").Parse(ceosTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, input)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("config ", b.String())

	conn, err := core.NewEOSDriver(
		*deviceName,
		base.WithAuthStrictKey(false),
		base.WithAuthUsername(*username),
		base.WithAuthPassword(*password),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conf, err := cfg.NewEOSCfg(conn)
	if err != nil {
		log.Fatal(err)
	}

	err = conf.Prepare()
	if err != nil {
		log.Fatal(err)
	}

	_, err = conf.LoadConfig(b.String(), false)
	if err != nil {
		log.Fatal(err)
	}

}
