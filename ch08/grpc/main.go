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
	"strings"
	"time"

	"grpc/pkg/oc"
	xr "grpc/proto/ems"
	"grpc/proto/telemetry"

	"github.com/openconfig/ygot/ygot"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

//go:generate bash $PWD/generate_code

const (
	blue  = "\x1b[34;1m"
	white = "\x1b[0m"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
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

func (t targetModel) createPayload() (string, error) {
	return ygot.EmitJSON(t.device, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
}

type targetModel struct {
	device *oc.Device
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

	//////////////////////
	// Read config inputs
	/////////////////////
	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var input Model
	err = d.Decode(&input)
	check(err)

	/////////////////////////////////
	//Build OpenConfig configuration
	////////////////////////////////
	tgtModel := targetModel{device: &oc.Device{}}

	err = input.buildNetworkInstance(tgtModel.device)
	check(err)

	payload, err := tgtModel.createPayload()
	check(err)

	///////////////////////////////////////////////////
	// Connect to target device (DevNet IOS XR device)
	//////////////////////////////////////////////////
	router, err := iosxr.Connect()
	check(err)
	defer router.conn.Close()

	/////////////////////
	// Replace BGP config
	/////////////////////
	// TODO: It fails if router is configured on a different ASN
	err = router.ReplaceConfig(payload)
	check(err)

	fmt.Printf("\n\n%sBGP%s config applied on %s\n\n\n", blue, white, router.conn.Target())

	///////////////////
	// Read BGP config
	///////////////////
	out, err := router.GetConfig("bgp.json")
	check(err)

	fmt.Printf("Config from %s:\n%s\n", iosxr.Hostname, out.Running)

	///////////////////
	// Stream Telemetry
	///////////////////

	ctx, cancel := context.WithCancel(router.ctx)
	defer cancel()
	router.ctx = ctx

	ch, errCh, err := router.GetSubscription("BGP", "gpbkv")
	check(err)

	//////////////////////////////////
	// Deal with Session cancellation
	/////////////////////////////////
	sigCh := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sigCh, os.Interrupt)
	defer func() {
		signal.Stop(sigCh)
		cancel()
	}()
	go router.SessionCancel(errCh, sigCh, cancel)

	/////////////////////////////////////////////////////////////
	// Decode Telemetry Protobuf message (payload still a string)
	/////////////////////////////////////////////////////////////

	for msg := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(msg, message)
		check(err)
		fmt.Printf("\n%s\n", strings.Repeat("-", 4))
		t := time.UnixMilli(int64(message.GetMsgTimestamp()))
		fmt.Printf("Time: %v\nPath: %v\n\n", t.Format(time.ANSIC), message.GetEncodingPath())

		b, err := json.Marshal(message.GetDataGpbkv())
		check(err)

		j := string(b)

		// https://go.dev/play/p/uyWenG-1Keu
		data := gjson.Get(j, "0.fields.0.fields.#(name==neighbor-address).ValueByType.StringValue")
		fmt.Println("  Neighbor: ", data)

		data = gjson.Get(j, "0.fields.1.fields.#(name==connection-state).ValueByType.StringValue")
		fmt.Println("  Connection state: ", data)
	}
}
