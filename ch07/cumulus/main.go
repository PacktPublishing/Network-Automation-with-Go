package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

var defaultNVUEPort = 8765

type cvx struct {
	url   string
	token string
	httpC http.Client
}

type IPAddress struct {
	Address map[string]struct{} `json:"address,omitempty"`
}

type Interface struct {
	IP   *IPAddress `json:"ip,omitempty"`
	Type string     `json:"type"`
}

type redistribute struct {
	Enabled string `json:"enable,omitempty"`
}

type addressFamily struct {
	Enabled      string                  `json:"enable,omitempty"`
	Redistribute map[string]redistribute `json:"redistribute,omitempty"`
}

type neighbor struct {
	RemoteAS int    `json:"remote-as,omitempty"`
	Type     string `json:"type,omitempty"`
}

type bgp struct {
	ASN      int                      `json:"autonomous-system,omitempty"`
	RouterID string                   `json:"router-id,omitempty"`
	AF       map[string]addressFamily `json:"address-family,omitempty"`
	Enabled  string                   `json:"enable,omitempty"`
	Neighbor map[string]neighbor      `json:"neighbor,omitempty"`
}

type router struct {
	Bgp bgp `json:"bgp"`
}

type vrf struct {
	Router router `json:"router"`
}

type nvue struct {
	Interface map[string]Interface `json:"interface"`
	Router    router               `json:"router"`
	Vrf       map[string]vrf       `json:"vrf"`
}

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
	var cfg nvue
	cfg.Interface = map[string]Interface{
		"lo": {
			Type: "loopback",
			IP: &IPAddress{
				Address: map[string]struct{}{
					fmt.Sprintf("%s/32", in.Loopback.IP): {},
				},
			},
		},
	}
	for _, uplink := range in.Uplinks {
		cfg.Interface[uplink.Name] = Interface{
			Type: "swp",
			IP: &IPAddress{
				Address: map[string]struct{}{
					uplink.Prefix: {},
				},
			},
		}
	}
	cfg.Router = router{
		Bgp: bgp{
			RouterID: in.Loopback.IP,
			ASN:      in.ASN,
		},
	}

	var peers = make(map[string]neighbor)
	for _, peer := range in.Peers {
		peers[peer.IP] = neighbor{
			RemoteAS: peer.ASN,
			Type:     "numbered",
		}
	}

	my_bgp := bgp{
		AF: map[string]addressFamily{
			"ipv4-unicast": {
				Enabled: "on",
				Redistribute: map[string]redistribute{
					"connected": {
						Enabled: "on",
					},
				},
			},
		},
		Enabled:  "on",
		Neighbor: peers,
	}

	cfg.Vrf = map[string]vrf{
		"default": {
			Router: router{
				Bgp: my_bgp,
			},
		},
	}
	view, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return b, fmt.Errorf("failed to create json indented view of config: %w", err)
	}
	log.Print("Generated config: ", string(view))

	err = json.NewEncoder(&b).Encode(cfg)
	if err != nil {
		return b, fmt.Errorf("failed json encode the config: %w", err)
	}

	return b, nil
}

func createRevision(c cvx) (string, error) {
	revisionPath := "/nvue_v1/revision"
	addr, err := url.Parse(c.url + revisionPath)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", addr.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+c.token)

	res, err := c.httpC.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var response map[string]interface{}

	json.NewDecoder(res.Body).Decode(&response)

	// Check this out, is this the intended behavior?
	for key := range response {
		return key, nil
	}

	return "", fmt.Errorf("unexpected createRevision error")
}

func applyRevision(c cvx, id string) error {
	applyPath := "/nvue_v1/revision/" + url.PathEscape(id)

	body := []byte("{\"state\": \"apply\", \"auto-prompt\": {\"ays\": \"ays_yes\", \"ignore_fail\": \"ignore_fail_yes\"}} ")

	req, err := http.NewRequest("PATCH", c.url+applyPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+c.token)

	res, err := c.httpC.Do(req)
	if err != nil {
		return err
	}

	return res.Body.Close()
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

	cfg, err := devConfig(input)
	r.check(err, "Couldn't create device specific configuration: "+moduleArgs.Input)

	// Store all HTTP device-related parameters
	device := cvx{
		url:   fmt.Sprintf("https://%s:%d", moduleArgs.Host, defaultNVUEPort),
		token: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", moduleArgs.User, moduleArgs.Password))),
		httpC: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	revisionID, err := createRevision(device)
	r.check(err, "Couldn't create a new candidate configuration revision: ")

	addr, err := url.Parse(device.url + "/nvue_v1/")
	r.check(err, "Couldn't parse the device API URL: "+device.url+"/nvue_v1/")
	params := url.Values{}
	params.Add("rev", revisionID)
	addr.RawQuery = params.Encode()

	req, err := http.NewRequest("PATCH", addr.String(), &cfg)
	r.check(err, "Couldn't create a request to add the desired configuration in the candidate configuration store")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+device.token)

	res, err := device.httpC.Do(req)
	r.check(err, "Couldn't make the request to the device")

	if err := applyRevision(device, revisionID); err != nil {
		r.Msg = "Couldn't apply candidate revision"
		FailJSON(r)
	}
	err = res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Msg = "Config applied with revisionID: " + revisionID
	r.Changed = true
	r.Failed = false
	returnResponse(r)

}
