package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	cfg "github.com/scrapli/scrapligocfg"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"

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

// ModuleArgs are the module inputs
type ModuleArgs struct {
	Host     string
	User     string
	Password string
	Input    string
}

// Response are the values returned from the module
type Response struct {
	Msg     string `json:"msg"`
	Busy    bool   `json:"busy"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

// ExitJSON is ...
func ExitJSON(responseBody Response) {
	returnResponse(responseBody)
}

// FailJSON is ...
func FailJSON(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(r Response) {
	var response []byte
	var err error
	response, err = json.Marshal(r)
	if err != nil {
		response, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if r.Failed {
		os.Exit(1)
	}
	os.Exit(0)
}

func (r Response) check(err error, msg string) {
	if err != nil {
		r.Msg = msg + ": " + err.Error()
		FailJSON(r)
	}
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
	return b, nil
}

func main() {
	var r Response

	if len(os.Args) != 2 {
		r.Msg = "No argument file provided"
		FailJSON(r)
	}

	argsFile := os.Args[1]

	text, err := os.ReadFile(argsFile)
	r.check(err, "Could not read configuration file: "+argsFile)

	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)
	r.check(err, "Ansible inputs are not valid (JSON): "+argsFile)

	src, err := base64.StdEncoding.DecodeString(moduleArgs.Input)
	r.check(err, "Couldn't decode the configuration inputs file: "+moduleArgs.Input)
	reader := bytes.NewReader(src)

	d := yaml.NewDecoder(reader)

	var input Model
	err = d.Decode(&input)
	r.check(err, "Couldn't decode configuration inputs: "+string(src))

	config, err := devConfig(input)
	r.check(err, "Couldn't create an EOS specific config for: "+string(src))

	conn, err := platform.NewPlatform(
		"arista_eos",
		moduleArgs.Host,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(moduleArgs.User),
		options.WithAuthPassword(moduleArgs.Password),
	)
	r.check(err, "Couldn't create client connection for: "+moduleArgs.Host)
	driver, err := conn.GetNetworkDriver()
	r.check(err, "Couldn't create driver for: "+moduleArgs.Host)

	err = driver.Open()
	r.check(err, "Couldn't connect to: "+moduleArgs.Host)
	defer driver.Close()

	conf, err := cfg.NewCfg(driver, "arista_eos")
	r.check(err, "Couldn't create a config with scrapli for: "+moduleArgs.Host)

	err = conf.Prepare()
	r.check(err, "Couldn't prepare a config with scrapli for: "+moduleArgs.Host)

	_, err = conf.LoadConfig(config.String(), false)
	r.check(err, "Couldn't load the config with scrapli for: "+moduleArgs.Host)

	r.Msg = "Configuration applied on: " + moduleArgs.Host
	r.Changed = true
	r.Failed = false
	returnResponse(r)
}