package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"grpc/pkg/oc"
	xr "grpc/proto/ems"
	"grpc/proto/telemetry"

	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

//go:generate bash $PWD/create_models

const (
	// xrLoopback = "Loopback0"
	//defaultSubIdx  = 0
	defaultNetInst = "default"
	blue           = "\x1b[34;1m"
	white          = "\x1b[0m"
	// OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP E_OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE = 1
	bgpID = oc.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP
	// OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST E_OpenconfigBgpTypes_AFI_SAFI_TYPE = 3
	ipv4uniAF = oc.OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST
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

// Provides the user/password for the connection. It implements
// the PerRPCCredentials interface.
type loginCreds struct {
	Username, Password string
	requireTLS         bool
}

// Method of the PerRPCCredentials interface.
func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}, nil
}

// Method of the PerRPCCredentials interface.
func (c *loginCreds) RequireTransportSecurity() bool {
	return c.requireTLS
}

func (r IOSXR) Connect() (xr xrgrpc, err error) {
	// Hardcoded. Don't do at home.
	port := ":57777"

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	creds := credentials.NewTLS(config)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// Add gRPC overall timeout to the config options array.
	// Hardcoded at 10 seconds. Don't do at home.
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(10))

	// Add user/password to config options array.
	opts = append(opts, grpc.WithPerRPCCredentials(&loginCreds{
		Username:   r.Username,
		Password:   r.Password,
		requireTLS: true}))

	conn, err := grpc.DialContext(ctx, r.Hostname+port, opts...)
	if err != nil {
		return xr, fmt.Errorf("could not build a router: %w", err)
	}
	xr.conn = conn
	xr.ctx = ctx

	return xr, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type xrgrpc struct {
	conn *grpc.ClientConn
	ctx  context.Context
}

func (x *xrgrpc) ReplaceConfig(json string) error {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()

	// 'g' is the gRPC stub.
	g := xr.NewGRPCConfigOperClient(x.conn)

	// 'a' is the object we send to the router via the stub.
	a := xr.ConfigArgs{ReqId: id, Yangjson: json}

	// 'r' is the result that comes back from the target.
	r, err := g.ReplaceConfig(x.ctx, &a)
	if err != nil {
		return fmt.Errorf("cannot replace the config: %w", err)
	}
	if len(r.GetErrors()) != 0 {
		return fmt.Errorf("error triggered by remote host for ReqId: %v; %s", id, r.GetErrors())
	}
	return nil
}

func (x *xrgrpc) GetConfig(file string) (cfg Config, err error) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()
	cfg.Device = x.conn.Target()
	cfg.Timestamp = time.Now()

	var paths string
	// Get config for the YANG paths
	if file != "" {
		f, err := os.ReadFile(file)
		if err != nil {
			return cfg, fmt.Errorf("could not read file: %v: %w", file, err)

		}
		paths = string(f)
	}
	// 'g' is the gRPC stub.
	g := xr.NewGRPCConfigOperClient(x.conn)

	// 'a' is the object we send to the router via the stub.
	a := xr.ConfigGetArgs{ReqId: id, Yangpathjson: paths}

	// 'st' is the streamed result that comes back from the target.
	st, err := g.GetConfig(x.ctx, &a)
	if err != nil {
		return cfg, fmt.Errorf("could not get the config from %s: %w", x.conn.Target(), err)
	}
	for {
		// Loop through the responses in the stream until there is nothing left.
		r, err := st.Recv()
		if err == io.EOF {
			return cfg, nil
		}
		if len(r.GetErrors()) != 0 {
			return cfg, fmt.Errorf("error triggered by remote host for ReqId: %v; %s", id, r.GetErrors())
		}
		if len(r.GetYangjson()) > 0 {
			cfg.Running += r.GetYangjson()
		}
	}
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

	device := &oc.Device{}

	err = input.buildNetworkInstance(device)
	check(err)

	// Generate the payload payload for our message
	payload, err := ygot.EmitJSON(device, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	fmt.Printf("%s\n", payload)

	check(err)

	/////////////
	// Connect
	////////////
	router, err := iosxr.Connect()
	check(err)
	defer router.conn.Close()

	///////////////////
	// Replace BGP config
	///////////////////
	err = router.ReplaceConfig(payload)
	check(err)

	fmt.Printf("\n\n\n%sBGP%s config applied on %s\n\n\n", blue, white, router.conn.Target())

	///////////////////
	// Read BGP config
	///////////////////
	var out Config
	paths := "bgp.json"

	out, err = router.GetConfig(paths)
	check(err)

	fmt.Printf("Config from %s:\n%s\n", iosxr.Hostname, out.Running)

	///////////////////
	// Stream Telemetry
	///////////////////

	ctx, cancel := context.WithCancel(router.ctx)
	defer cancel()
	router.ctx = ctx

	ch, ech, err := router.GetSubscription("BGP", "gpbkv")
	check(err)

	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			fmt.Printf("\nmanually cancelled the session to %v\n\n", router.conn.Target())
			cancel()
			return
		case <-ctx.Done():
			// Timeout: "context deadline exceeded"
			err = ctx.Err()
			fmt.Printf("\ngRPC session timed out after %s seconds: %v\n\n", "10", err.Error())
			return
		case err = <-ech:
			// Session canceled: "context canceled"
			fmt.Printf("\ngRPC session to %v failed: %v\n\n", router.conn.Target(), err.Error())
			return
		}
	}()

	for msg := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(msg, message)
		check(err)

		fmt.Printf("Time %v, Path: %v\n", message.GetMsgTimestamp(), message.GetEncodingPath())

		b, err := json.Marshal(message)
		check(err)

		bjs, err := prettyprint(b)
		check(err)

		fmt.Println(string(bjs))
	}

}
