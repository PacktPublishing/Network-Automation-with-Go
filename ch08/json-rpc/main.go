package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	api "json-rpc/pkg/srl"
	"log"
	"net/http"
	"os"

	"github.com/openconfig/ygot/ygot"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/openconfig/ygot/generator -path=yang -generate_fakeroot -fakeroot_name=device -output_file=pkg/srl/srl.go -package_name=srl yang/srl_nokia/models/network-instance/srl_nokia-bgp.yang yang/srl_nokia/models/routing-policy/srl_nokia-routing-policy.yang yang/srl_nokia/models/network-instance/srl_nokia-ip-route-tables.yang

const (
	srlLoopback    = "system0"
	defaultSubIdx  = 0
	defaultNetInst = "default"
)

var (
	hostname          = "http://clab-netgo-srl/jsonrpc"
	username          = "admin"
	password          = "admin"
	defaultPolicyName = "all"
	defaultBGPGroup   = "EBGP"
)

// SRL JSON-RPC request
type RpcRequest struct {
	Version string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

// SRL JSON-RPC response
type RpcResponse struct {
	Version string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

// SRL JSON-RPC Params
type Params struct {
	Commands []*Command `json:"commands"`
}

// SRL JSON-RPC Command
type Command struct {
	Action string      `json:"action"`
	Path   string      `json:"path"`
	Value  interface{} `json:"value"`
}

// Input Data Model
type Model struct {
	Uplinks  []Link `yaml:"uplinks"`
	Peers    []Peer `yaml:"peers"`
	ASN      int    `yaml:"asn"`
	Loopback Addr   `yaml:"loopback"`
}

// Input Data Model L3 link
type Link struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

// Input Data Model BGP Peer
type Peer struct {
	IP  string `yaml:"ip"`
	ASN int    `yaml:"asn"`
}

// Input Data Model IPv4 addr
type Addr struct {
	IP string `yaml:"ip"`
}

func (m *Model) buildL3Interfaces(dev *api.Device) error {
	links := m.Uplinks
	links = append(
		links,
		Link{
			Name:   srlLoopback,
			Prefix: fmt.Sprintf("%s/32", m.Loopback.IP),
		},
	)

	for _, link := range links {
		intf, err := dev.NewInterface(link.Name)
		if err != nil {
			return err
		}
		subintf, err := intf.NewSubinterface(defaultSubIdx)
		if err != nil {
			return err
		}

		subintf.Ipv4 = &api.SrlNokiaInterfaces_Interface_Subinterface_Ipv4{}
		subintf.Ipv4.NewAddress(link.Prefix)

		if err := intf.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) buildNetworkInstance(dev *api.Device) error {
	ni, err := dev.NewNetworkInstance(defaultNetInst)
	if err != nil {
		return err
	}

	links := m.Uplinks
	links = append(
		links,
		Link{
			Name:   srlLoopback,
			Prefix: fmt.Sprintf("%s/32", m.Loopback.IP),
		},
	)
	for _, link := range links {
		linkName := fmt.Sprintf("%s.%d", link.Name, defaultSubIdx)
		ni.NewInterface(linkName)
	}

	ni.Protocols = &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols{
		Bgp: &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp{
			AutonomousSystem: ygot.Uint32(uint32(m.ASN)),
			RouterId:         ygot.String(m.Loopback.IP),
			Ipv4Unicast: &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp_Ipv4Unicast{
				AdminState: api.SrlNokiaBgp_AdminState_enable,
			},
		},
	}

	g, err := ni.Protocols.Bgp.NewGroup(defaultBGPGroup)
	if err != nil {
		return err
	}
	g.ExportPolicy = ygot.String(defaultPolicyName)
	g.ImportPolicy = ygot.String(defaultPolicyName)

	for _, peer := range m.Peers {
		n, err := ni.Protocols.Bgp.NewNeighbor(peer.IP)
		if err != nil {
			return err
		}
		n.PeerAs = ygot.Uint32(uint32(peer.ASN))
		n.PeerGroup = ygot.String(defaultBGPGroup)
	}

	if err := ni.Validate(); err != nil {
		return err
	}

	return nil

}

func (m *Model) buildDefaultPolicy(dev *api.Device) error {
	dev.RoutingPolicy = &api.SrlNokiaRoutingPolicy_RoutingPolicy{}

	p, err := dev.RoutingPolicy.NewPolicy(defaultPolicyName)
	if err != nil {
		return err
	}

	// populating LocalPreference due to YGOT not supporting presence containers, see "ygot/issues/329"
	p.DefaultAction = &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction{
		Accept: &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction_Accept{
			Bgp: &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction_Accept_Bgp{
				LocalPreference: &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction_Accept_Bgp_LocalPreference{
					Set: ygot.Uint32(uint32(100)),
				},
			},
		},
	}

	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}

func printYgot(s ygot.ValidatedGoStruct) {
	t, _ := ygot.EmitJSON(s, &ygot.EmitJSONConfig{
		Format: ygot.Internal,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	},
	)
	fmt.Println(t)
}

func main() {

	src, err := os.Open("input.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	if err != nil {
		log.Fatal(err)
	}

	device := &api.Device{}

	if err := input.buildDefaultPolicy(device); err != nil {
		log.Fatal(err)
	}

	if err := input.buildL3Interfaces(device); err != nil {
		log.Fatal(err)
	}

	if err := input.buildNetworkInstance(device); err != nil {
		log.Fatal(err)
	}

	printYgot(device)
	v, err := ygot.ConstructIETFJSON(device, nil)
	if err != nil {
		log.Fatal(err)
	}

	value, err := json.Marshal(RpcRequest{
		Version: "2.0",
		ID:      0,
		Method:  "set",
		Params: Params{
			Commands: []*Command{
				{
					Action: "update",
					Path:   "/",
					Value:  v,
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(
		"POST",
		hostname,
		bytes.NewBuffer(value),
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", username, password)),
		),
	)

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body))

	var result RpcResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	if result.Error != nil {
		log.Fatalf(
			"failed to configure the device: %v",
			result.Error,
		)
	}

	log.Println("Successfully configured the device")
}
