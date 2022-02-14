package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

	my_bgp :=  bgp{
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
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)

	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	hostname := flag.String("device", "clab-netgo-cvx", "Device Hostname")
	username := flag.String("username", "cumulus", "SSH Username")
	password := flag.String("password", "cumulus", "SSH password")
	flag.Parse()

	// read and parse the input file
	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	check(err)
     
	cfg, err := devConfig(input)
	check(err)

	// Store all HTTP device-related parameters
	device := cvx{
		url:   fmt.Sprintf("https://%s:%d", *hostname, defaultNVUEPort),
		token: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", *username, *password))),
		httpC: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	// create a new candidate configuration revision
	revisionID, err := createRevision(device)
	check(err)

	log.Print("Created revisionID: ", revisionID)

	addr, err := url.Parse(device.url + "/nvue_v1/")
	check(err)
	params := url.Values{}
	params.Add("rev", revisionID)
	addr.RawQuery = params.Encode()

	// Save the device desired configuration in candidate configuration store
	req, err := http.NewRequest("PATCH", addr.String(), &cfg)
	check(err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+device.token)

	res, err := device.httpC.Do(req)
	check(err)
	defer res.Body.Close()

	// Apply candidate revision
	if err := applyRevision(device, revisionID); err != nil {
		log.Fatal(err)
	}
}
