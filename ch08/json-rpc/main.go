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

	"github.com/openconfig/ygot/ygot"
)

// docker exec -it clab-netgo-srl /opt/srlinux/bin/sr_cli

//go:generate go run github.com/openconfig/ygot/generator -path=yang -output_file=pkg/srl/srl.go -package_name=srl yang/srl_nokia/models/network-instance/srl_nokia-bgp.yang yang/srl_nokia/models/routing-policy/srl_nokia-routing-policy.yang

var (
	hostname = "http://clab-netgo-srl/jsonrpc"
	username = "admin"
	password = "admin"
)

type RPC struct {
	Version string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	Commands []*Command `json:"commands"`
}

type Command struct {
	Action string      `json:"action"`
	Path   string      `json:"path"`
	Value  interface{} `json:"value"`
}

func buildL3Interface(name, prefix string) (*Command, error) {
	intf := api.SrlNokiaInterfaces_Interface{}
	subintf, err := intf.NewSubinterface(0)
	if err != nil {
		return nil, err
	}

	subintf.Ipv4 = &api.SrlNokiaInterfaces_Interface_Subinterface_Ipv4{}
	subintf.Ipv4.NewAddress(prefix)

	if err := intf.Validate(); err != nil {
		return nil, err
	}

	value, err := ygot.ConstructIETFJSON(&intf, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n/interface[name=%s]>\n", name)
	printYgot(&intf)

	return &Command{
		Action: "replace",
		Path:   fmt.Sprintf("/interface[name=%s]", name),
		Value:  value,
	}, nil
}

func buildBGPConfig(instance, id, group, policy, neighbor string, lasn, rasn uint32) (*Command, error) {

	bgp := &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp{
		AutonomousSystem: uintToPtr(lasn),
		RouterId:         strToPtr(id),
		Ipv4Unicast: &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp_Ipv4Unicast{
			AdminState: api.SrlNokiaBgp_AdminState_enable,
		},
	}

	g, err := bgp.NewGroup(group)
	if err != nil {
		return nil, err
	}
	g.ExportPolicy = strToPtr(policy)
	g.ImportPolicy = strToPtr(policy)

	n, err := bgp.NewNeighbor(neighbor)
	if err != nil {
		return nil, err
	}
	n.PeerAs = uintToPtr(rasn)
	n.PeerGroup = strToPtr(group)

	if err := bgp.Validate(); err != nil {
		return nil, err
	}

	value, err := ygot.ConstructIETFJSON(bgp, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n/network-instance[name=%s]/protocols/bgp>\n", instance)
	printYgot(bgp)

	return &Command{
		Action: "update",
		Path:   fmt.Sprintf("/network-instance[name=%s]/protocols/bgp", instance),
		Value:  value,
	}, nil

}

func moveIntfsToInstance(intfs []string, instance string) ([]*Command, error) {
	var cmds []*Command
	for _, intf := range intfs {
		path := fmt.Sprintf("/network-instance[name=%s]/interface[name=%s]", instance, intf)
		fmt.Println(path)
		cmds = append(cmds, &Command{
			Action: "update",
			Path:   path,
		})
	}

	return cmds, nil
}

func buildDefaultPolicy(name string) (*Command, error) {
	rp := api.SrlNokiaRoutingPolicy_RoutingPolicy{}
	p, err := rp.NewPolicy(name)
	if err != nil {
		return nil, err
	}

	// populating LocalPreference due to YGOT not supporting presence containers, see "ygot/issues/329"
	p.DefaultAction = &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction{
		Accept: &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction_Accept{
			Bgp: &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction_Accept_Bgp{
				LocalPreference: &api.SrlNokiaRoutingPolicy_RoutingPolicy_Policy_DefaultAction_Accept_Bgp_LocalPreference{
					Set: uintToPtr(uint32(100)),
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

	fmt.Printf("\n/routing-policy>\n")
	printYgot(&rp)

	return &Command{
		Action: "replace",
		Path:   "/routing-policy",
		Value:  value,
	}, nil
}

func buildSetRPC(cmds []*Command) RPC {
	return RPC{
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

	var cmds []*Command

	eth1, err := buildL3Interface("ethernet-1/1", "192.0.2.0/31")
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, eth1)

	lo0, err := buildL3Interface("system0", "198.51.100.0/32")
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, lo0)

	policy, err := buildDefaultPolicy("all")
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, policy)

	insts, err := moveIntfsToInstance([]string{"ethernet-1/1.0", "system0.0"}, "default")
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, insts...)

	bgp, err := buildBGPConfig("default", "198.51.100.0", "EBGP", "all", "192.0.2.1", uint32(65000), uint32(65001))
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status: %s", resp.Status)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func strToPtr(v string) *string  { return &v }
func uintToPtr(v uint32) *uint32 { return &v }
