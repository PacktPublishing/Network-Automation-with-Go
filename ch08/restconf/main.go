package main

import (
	api "restconf/pkg/eos"
)

//go:generate go run github.com/openconfig/ygot/generator -path=yang -output_file=pkg/eos/eos.go -compress_paths=true -exclude_modules=ietf-interfaces -package_name=eos yang/openconfig/public/release/models/bgp/openconfig-bgp.yang yang/openconfig/public/release/models/interfaces/openconfig-if-ip.yang

func main() {
	intf := &api.Interface{
		Name: strToPtr("Ethernet1"),
	}
	subIntf, _ := intf.NewSubinterface(0)
	subIntf.Ipv4 = &api.Interface_Subinterface_Ipv4{}
	subIntf.Ipv4.NewAddress("12.12.12.1/24")

}

func strToPtr(v string) *string  { return &v }
func uintToPtr(v uint32) *uint32 { return &v }
