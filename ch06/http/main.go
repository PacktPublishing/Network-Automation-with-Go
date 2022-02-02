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

var defaultNvuePort = 8765

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

func populateData(i Input, o *nvue) {
	o.Interface = map[string]Interface{
		"lo": Interface{
			Type: "loopback",
			IP: &IPAddress{
				Address: map[string]struct{}{
					fmt.Sprintf("%s/32", i.Loopback.IP): struct{}{},
				},
			},
		},
	}
	for _, uplink := range i.Uplinks {
		o.Interface[uplink.Name] = Interface{
			Type: "swp",
			IP: &IPAddress{
				Address: map[string]struct{}{
					uplink.Prefix: struct{}{},
				},
			},
		}
	}
	o.Router = router{
		Bgp: bgp{
			RouterID: i.Loopback.IP,
			ASN:      i.ASN,
		},
	}

	var peers = make(map[string]neighbor)
	for _, peer := range i.Peers {
		peers[peer.IP] = neighbor{
			RemoteAS: peer.ASN,
			Type:     "numbered",
		}
	}

	o.Vrf = map[string]vrf{
		"default": vrf{
			Router: router{
				Bgp: bgp{
					AF: map[string]addressFamily{
						"ipv4-unicast": addressFamily{
							Enabled: "on",
							Redistribute: map[string]redistribute{
								"connected": redistribute{
									Enabled: "on",
								},
							},
						},
					},
					Enabled:  "on",
					Neighbor: peers,
				},
			},
		},
	}
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

func main() {
	deviceName := flag.String("device", "clab-netgo-cvx", "Device Hostname")
	username := flag.String("username", "cumulus", "SSH Username")
	password := flag.String("password", "cumulus", "SSH password")
	flag.Parse()

	// store all device-related parameters
	device := cvx{
		url:   fmt.Sprintf("https://%s:%d", *deviceName, defaultNvuePort),
		token: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", *username, *password))),
		httpC: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

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

	// populate the device data model
	var data nvue
	populateData(input, &data)

	view, _ := json.MarshalIndent(data, "", " ")
	log.Print("generated config ", string(view))

	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	// create a new candidate configuration revision
	revisionID, err := createRevision(device)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Created revisionID: ", revisionID)

	addr, err := url.Parse(device.url + "/nvue_v1/")
	if err != nil {
		log.Fatal(err)
	}
	params := url.Values{}
	params.Add("rev", revisionID)
	addr.RawQuery = params.Encode()

	// save the device data model in candidate configuration store
	req, err := http.NewRequest("PATCH", addr.String(), body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+device.token)

	res, err := device.httpC.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// apply candidate revision
	err = applyRevision(device, revisionID)
	if err != nil {
		log.Fatal(err)
	}

}
