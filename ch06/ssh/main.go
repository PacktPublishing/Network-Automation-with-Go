package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

const srlTemplate = `
enter candidate
{{- range $uplink := .Uplinks }}
set / interface {{ $uplink.Name }} subinterface 0 ipv4 address {{ $uplink.Prefix }}
set / network-instance default interface {{ $uplink.Name }}.0
{{- end }}
set / interface system0 subinterface 0 ipv4 address {{ .Loopback.IP }}/32
set / network-instance default interface system0.0
set / routing-policy policy all default-action accept
set / network-instance default protocols bgp autonomous-system {{ .ASN }}
set / network-instance default protocols bgp router-id {{ .Loopback.IP }}
set / network-instance default protocols bgp group EBGP
set / network-instance default protocols bgp group EBGP export-policy all
set / network-instance default protocols bgp group EBGP import-policy all
{{- range $peer := .Peers }}
set / network-instance default protocols bgp neighbor {{ $peer.IP }} peer-as {{ $peer.ASN }}
set / network-instance default protocols bgp neighbor {{ $peer.IP }} peer-group EBGP
{{- end }}
set / network-instance default protocols bgp ipv4-unicast admin-state enable
commit now
quit
`

var (
	sshPort = 22
)

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
	device := flag.String("device", "clab-ch06-srl", "Device Hostname")
	username := flag.String("username", "admin", "SSH Username")
	password := flag.String("password", "admin", "SSH password")
	flag.Parse()

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

	t, err := template.New("config").Parse(srlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, input)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", *device, sshPort),
		config,
	)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal("failed to allocate stdin: ", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal("failed to allocate stdout: ", err)
	}
	defer func() {
		log.Printf("disconnected. dumping output...")
		io.Copy(log.Writer(), stdout)
	}()

	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	log.Print("connected. configuring...")
	b.WriteTo(stdin)

}
