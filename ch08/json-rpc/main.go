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

// docker exec -it clab-netgo-srl /opt/srlinux/bin/sr_cli

//go:generate go run github.com/openconfig/ygot/generator -path=yang -output_file=pkg/srl/srl.go -package_name=srl yang/srl_nokia/models/network-instance/srl_nokia-bgp.yang yang/srl_nokia/models/routing-policy/srl_nokia-routing-policy.yang yang/srl_nokia/models/network-instance/srl_nokia-ip-route-tables.yang

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
	Version string       `json:"jsonrpc"`
	ID      int          `json:"id"`
	Result  *interface{} `json:"result,omitempty"`
	Error   *interface{} `json:"error,omitempty"`
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

func (m *Model) buildL3Interfaces() ([]*Command, error) {
	var cmds []*Command

	links := m.Uplinks
	links = append(links, Link{Name: srlLoopback, Prefix: fmt.Sprintf("%s/32", m.Loopback.IP)})

	for _, link := range links {
		intf := api.SrlNokiaInterfaces_Interface{}
		subintf, err := intf.NewSubinterface(defaultSubIdx)
		if err != nil {
			return nil, err
		}

		subintf.Ipv4 = &api.SrlNokiaInterfaces_Interface_Subinterface_Ipv4{}
		subintf.Ipv4.NewAddress(link.Prefix)

		if err := intf.Validate(); err != nil {
			return nil, err
		}

		value, err := ygot.ConstructIETFJSON(&intf, nil)
		if err != nil {
			return nil, err
		}

		fmt.Printf("\n/interface[name=%s]:\n", link.Name)
		printYgot(&intf)

		cmds = append(cmds, &Command{
			Action: "replace",
			Path:   fmt.Sprintf("/interface[name=%s]", link.Name),
			Value:  value,
		})
	}
	return cmds, nil
}

func (m *Model) buildBGPConfig() (*Command, error) {

	bgp := &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp{
		AutonomousSystem: ygot.Uint32(uint32(m.ASN)),
		RouterId:         ygot.String(m.Loopback.IP),
		Ipv4Unicast: &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp_Ipv4Unicast{
			AdminState: api.SrlNokiaBgp_AdminState_enable,
		},
	}

	g, err := bgp.NewGroup(defaultBGPGroup)
	if err != nil {
		return nil, err
	}
	g.ExportPolicy = ygot.String(defaultPolicyName)
	g.ImportPolicy = ygot.String(defaultPolicyName)

	for _, peer := range m.Peers {
		n, err := bgp.NewNeighbor(peer.IP)
		if err != nil {
			return nil, err
		}
		n.PeerAs = ygot.Uint32(uint32(peer.ASN))
		n.PeerGroup = ygot.String(defaultBGPGroup)
	}

	if err := bgp.Validate(); err != nil {
		return nil, err
	}

	dev := api.DeviceRoot("eos")
	gotPath, _, errs := ygot.ResolvePath(dev.NetworkInstance(defaultNetInst))
	if errs != nil {
		return nil, err
	}

	s, _ := ygot.PathToString(gotPath)
	fmt.Printf("!!! Path %s\n", s)

	value, err := ygot.ConstructIETFJSON(bgp, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n/network-instance[name=%s]/protocols/bgp:\n", defaultNetInst)
	printYgot(bgp)

	return &Command{
		Action: "replace",
		Path:   fmt.Sprintf("/network-instance[name=%s]/protocols/bgp", defaultNetInst),
		Value:  value,
	}, nil

}

func (m *Model) moveIntfsToInstance() ([]*Command, error) {
	var cmds []*Command

	var intfs []string
	for _, link := range m.Uplinks {
		intfs = append(intfs, fmt.Sprintf("%s.%d", link.Name, defaultSubIdx))
	}
	intfs = append(intfs, fmt.Sprintf("%s.%d", srlLoopback, defaultSubIdx))

	for _, intf := range intfs {
		path := fmt.Sprintf("/network-instance[name=%s]/interface[name=%s]", defaultNetInst, intf)
		fmt.Println(path)
		cmds = append(cmds, &Command{
			Action: "update",
			Path:   path,
		})
	}

	return cmds, nil
}

func (m *Model) buildDefaultPolicy() (*Command, error) {
	rp := api.SrlNokiaRoutingPolicy_RoutingPolicy{}
	p, err := rp.NewPolicy(defaultPolicyName)
	if err != nil {
		return nil, err
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

	if err := rp.Validate(); err != nil {
		return nil, err
	}
	value, err := ygot.ConstructIETFJSON(&rp, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n/routing-policy:\n")
	printYgot(&rp)

	return &Command{
		Action: "replace",
		Path:   "/routing-policy",
		Value:  value,
	}, nil
}

func buildSetRPC(cmds []*Command) RpcRequest {
	return RpcRequest{
		Version: "2.0",
		ID:      0,
		Method:  "set",
		Params: Params{
			Commands: cmds,
		},
	}
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

	var cmds []*Command

	l3Intfs, err := input.buildL3Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, l3Intfs...)

	policy, err := input.buildDefaultPolicy()
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, policy)

	insts, err := input.moveIntfsToInstance()
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, insts...)

	bgp, err := input.buildBGPConfig()
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, bgp)

	value, _ := json.Marshal(buildSetRPC(cmds))

	req, err := http.NewRequest("POST", hostname, bytes.NewBuffer(value))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))

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
		log.Fatalf("failed to configure the device: %v", *result.Error)
	}

	log.Println("Successfully configured the device")
}
