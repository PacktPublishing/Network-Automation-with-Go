package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	api "restconf/pkg/eos"

	"github.com/openconfig/ygot/ygot"
)

//go:generate go run github.com/openconfig/ygot/generator -path=yang -output_file=pkg/eos/eos.go -compress_paths=true -exclude_modules=ietf-interfaces -package_name=eos yang/openconfig/public/release/models/bgp/openconfig-bgp.yang yang/openconfig/public/release/models/interfaces/openconfig-if-ip.yang yang/openconfig/public/release/models/network-instance/openconfig-network-instance.yang

func main() {
	intf := &api.Interface{
		Name: strToPtr("Ethernet1"),
	}
	subIntf, err := intf.NewSubinterface(0)
	if err != nil {
		log.Fatal(err)
	}

	subIntf.Ipv4 = &api.Interface_Subinterface_Ipv4{}
	addr, err := subIntf.Ipv4.NewAddress("12.12.12.1")
	if err != nil {
		log.Fatal(err)
	}
	addr.PrefixLength = uint8ToPtr(24)

	if err := intf.Validate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(printYgot(intf))

	lo := &api.Interface{
		Name: strToPtr("Loopback0"),
	}
	subIntf, err = lo.NewSubinterface(0)
	if err != nil {
		log.Fatal(err)
	}

	subIntf.Ipv4 = &api.Interface_Subinterface_Ipv4{}
	addr, err = subIntf.Ipv4.NewAddress("1.1.1.1")
	if err != nil {
		log.Fatal(err)
	}
	addr.PrefixLength = uint8ToPtr(32)

	if err := intf.Validate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(printYgot(lo))

	netInst := &api.NetworkInstance{
		Name: strToPtr("default"),
	}
	protocol, _ := netInst.NewProtocol(api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP")

	protocol.Bgp = &api.NetworkInstance_Protocol_Bgp{}

	protocol.Bgp.Global = &api.NetworkInstance_Protocol_Bgp_Global{}
	protocol.Bgp.Global.As = uint32ToPtr(65000)

	n, err := protocol.Bgp.NewNeighbor("23.23.23.2")
	if err != nil {
		log.Fatal(err)
	}
	n.PeerAs = uint32ToPtr(65100)

	_, err = n.NewAfiSafi(api.OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST)
	if err != nil {
		log.Fatal(err)
	}

	//netInst.NewTableConnection(api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_DIRECTLY_CONNECTED, api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, api.OpenconfigTypes_ADDRESS_FAMILY_IPV4)
	//netInst.NewTableConnection(api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_DIRECTLY_CONNECTED, api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, api.OpenconfigTypes_ADDRESS_FAMILY_IPV6)

	fmt.Println(printYgot(netInst))

	value, _ := ygot.Marshal7951(netInst)

	fullQuery := "https://clab-netgo-ceos:6020/restconf/data/network-instances/network-instance=default"
	fmt.Println("full query ", fullQuery)
	req, err := http.NewRequest("POST", fullQuery, bytes.NewBuffer(value))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "admin", "admin"))))

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
	fmt.Println(string(body))

}

func printYgot(s ygot.ValidatedGoStruct) string {
	t, _ := ygot.EmitJSON(s, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	},
	)
	return t
}

func strToPtr(v string) *string    { return &v }
func uint32ToPtr(v uint32) *uint32 { return &v }
func uint8ToPtr(v uint8) *uint8    { return &v }
