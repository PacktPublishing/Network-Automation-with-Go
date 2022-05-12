package main

import (
	"context"
	"fmt"
	"os"
	"time"

	api "grpc/pkg/xr"

	"github.com/nleiva/xrgrpc"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

//go:generate bash $PWD/create_models

const (
	xrLoopback = "Loopback0"
	//defaultSubIdx  = 0
	defaultNetInst = "default"
	blue           = "\x1b[34;1m"
	white          = "\x1b[0m"
	// OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP E_OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE = 1
	bgpID = api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP
	// OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST E_OpenconfigBgpTypes_AFI_SAFI_TYPE = 3
	ipv4uniAF = api.OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST
)

type Authentication struct {
	Username string
	Password string
}

type IOSXR struct {
	Hostname string
	Authentication
}

type Config struct {
	Device    string
	Running   string
	Timestamp time.Time
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

func grpcConfig(ctx context.Context, conn *grpc.ClientConn, file string) (Config, error) {
	var response string

	var paths string
	// Get config for the YANG paths
	if file != "" {
		f, err := os.ReadFile(file)
		if err != nil {
			return Config{}, fmt.Errorf("could not read file: %v: %w", file, err)

		}
		paths = string(f)
	}
	var id int64 = 1
	response, err := xrgrpc.GetConfig(ctx, conn, paths, id)
	if err != nil {
		return Config{}, fmt.Errorf("could not get the config from %s: %w", conn.Target(), err)
	}

	return Config{
		Device:    conn.Target(),
		Running:   response,
		Timestamp: time.Now(),
	}, nil
}

func (r IOSXR) Connect() (*grpc.ClientConn, context.Context, error) {
	// Hardcoded, don't do at home
	port := ":57777"

	router, err := xrgrpc.BuildRouter(
		xrgrpc.WithUsername(r.Username),
		xrgrpc.WithPassword(r.Password),
		xrgrpc.WithHost(r.Hostname+port),
		xrgrpc.WithTimeout(5),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not build a router: %w", err)
	}

	conn, ctx, err := xrgrpc.Connect(*router)
	if err != nil {
		return nil, nil,
			fmt.Errorf("could not setup a client connection to %s: %w", router.Host, err)
	}
	return conn, ctx, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (m *Model) buildNetworkInstance(dev *api.Device) error {
	name := defaultNetInst
	nis := &api.OpenconfigNetworkInstance_NetworkInstances{}
	ni, err := nis.NewNetworkInstance(name)
	if err != nil {
		return fmt.Errorf("cannot create new network instance: %w", err)
	}
	ni.Config = &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Config{
		Name: &name,
	}

	ni.Protocols = &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols{}
	bgp, err := ni.Protocols.NewProtocol(bgpID, name)
	if err != nil {
		return fmt.Errorf("cannot create new bgp instance: %w", err)
	}
	bgp.Config = &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Config{
		Name:       &name,
		Identifier: bgpID,
	}

	bgp.Bgp = &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp{
		Global: &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global{
			Config: &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global_Config{
				As:       ygot.Uint32(uint32(m.ASN)),
				RouterId: ygot.String(m.Loopback.IP),
			},
		},
	}
	// Initialize the IPv4 Unicast address family.

	bgp.Bgp.Global.AfiSafis = &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global_AfiSafis{}
	safi, err := bgp.Bgp.Global.AfiSafis.NewAfiSafi(ipv4uniAF)
	if err != nil {
		return fmt.Errorf("cannot enable bgp IPv4 address family: %w", err)
	}
	safi.Config = &api.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global_AfiSafis_AfiSafi_Config{
		AfiSafiName: ipv4uniAF,
		Enabled:     ygot.Bool(true),
	}

	if err := ni.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	dev.NetworkInstances = nis

	return nil
}

func main() {
	///////////
	// Device
	//////////
	iosxr := IOSXR{
		Hostname: "sandbox-iosxr-1.cisco.com",
		Authentication: Authentication{
			Username: "admin",
			Password: "C1sco12345",
		},
	}

	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	check(err)

	device := &api.Device{}

	err = input.buildNetworkInstance(device)
	check(err)

	// Generate the json payload for our message
	json, err := ygot.EmitJSON(device, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	fmt.Printf("%s\n", json)

	check(err)

	/////////////
	// Connect
	////////////
	conn, ctx, err := iosxr.Connect()
	check(err)
	defer conn.Close()

	///////////////////
	// Apply BGP config
	///////////////////

	_, err = xrgrpc.MergeConfig(ctx, conn, json, 18)
	check(err)

	fmt.Printf("\n\n\n%sBGP%s config applied on %s\n\n\n", blue, white, conn.Target())

	///////////////////
	// Read BGP config
	///////////////////
	var out Config
	// Empty paths reads full config
	paths := "bgp.json"

	out, err = grpcConfig(ctx, conn, paths)
	check(err)
	fmt.Printf("Config from %s:\n%s\n", iosxr.Hostname, out.Running)

}
