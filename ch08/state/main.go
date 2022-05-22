package main

import (
	"crypto/tls"
	"fmt"
	srlAPI "json-rpc/pkg/srl"
	"log"
	"net/http"
	"reflect"
	eosAPI "restconf/pkg/eos"
	"sync"

	"cuelang.org/go/cue/cuecontext"
	resty "github.com/go-resty/resty/v2"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

var deviceLoopbacks = map[int]string{
	0: "198.51.100.0/32",
	1: "198.51.100.1/32",
	2: "198.51.100.2/32",
}

type Authentication struct {
	Username string
	Password string
}

type CVX struct {
	Hostname string
	ID       int
	Authentication
}

type SRL struct {
	Hostname string
	ID       int
	Authentication
}

type CEOS struct {
	Hostname string
	ID       int
	Authentication
}

type Router interface {
	GetRoutes(wg *sync.WaitGroup)
}

// CVX RIB type
type CVXRIB struct {
	IPv4Route struct {
		Route map[string]interface{} `json:"route"`
	} `json:"ipv4"`
}

func (r CVX) GetRoutes(wg *sync.WaitGroup) {
	log.Print("Collecting CVX routes")

	client := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	client.SetBaseURL("https://" + r.Hostname + ":8765")
	client.SetBasicAuth(r.Username, r.Password)

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"rev": "operational",
		}).
		Get("/nvue_v1/vrf/default/router/rib")

	if err != nil {
		log.Printf("failed to send request for %s: %s", r.Hostname, err.Error())
		wg.Done()
		return
	}

	ctx := cuecontext.New()
	v := ctx.CompileBytes(resp.Body())

	var rib CVXRIB
	err = v.Value().Decode(&rib)
	if err != nil {
		log.Fatal(err)
	}

	out := []string{}
	for prefix := range rib.IPv4Route.Route {
		out = append(out, prefix)
	}

	expectedRoutes := make(map[string]bool)

	for id, loopback := range deviceLoopbacks {
		if id != r.ID {
			expectedRoutes[loopback] = false
		}
	}

	go checkRoutes(r.Hostname, out, expectedRoutes, wg)
}

// SRL JSON-RPC request
type SrlRequest struct {
	Version string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Method  string    `json:"method"`
	Params  SrlParams `json:"params"`
}

// SRL JSON-RPC response
type SrlResponse struct {
	Version string       `json:"jsonrpc"`
	ID      int          `json:"id"`
	Result  *interface{} `json:"result,omitempty"`
	Error   *interface{} `json:"error,omitempty"`
}

// SRL JSON-RPC Params
type SrlParams struct {
	Commands []*SrlCommand `json:"commands"`
}

// SRL JSON-RPC Command
type SrlCommand struct {
	Action    string      `json:"action,omitempty"`
	Path      string      `json:"path"`
	Datastore string      `json:"datastore,omitempty"`
	Value     interface{} `json:"value,omitempty"`
}

// a modified version of srlAPI.Unmarshal
func srlUnmarshal(result interface{}, destStruct ygot.GoStruct) error {
	tn := reflect.TypeOf(&srlAPI.SrlNokiaNetworkInstance_NetworkInstance_RouteTable{}).Elem().Name()
	schema, ok := srlAPI.SchemaTree[tn]
	if !ok {
		return fmt.Errorf("could not find schema for type %s", tn)
	}
	return ytypes.Unmarshal(schema, destStruct, result)
}

func (r SRL) GetRoutes(wg *sync.WaitGroup) {
	log.Print("Collecting SRL routes")
	client := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	client.SetBaseURL("http://" + r.Hostname + "/jsonrpc")
	client.SetBasicAuth(r.Username, r.Password)
	client.SetDisableWarn(true)

	request := SrlRequest{
		Version: "2.0",
		ID:      0,
		Method:  "get",
		Params: SrlParams{
			Commands: []*SrlCommand{
				&SrlCommand{
					Path:      fmt.Sprintf("/network-instance[name=%s]/route-table", "default"),
					Datastore: "state",
				},
			},
		},
	}

	var response SrlResponse
	_, err := client.R().
		SetResult(&response).
		SetBody(request).
		Post("")
	if err != nil {
		log.Printf("failed to send request for %s: %s", r.Hostname, err.Error())
		wg.Done()
		return
	}

	routes := &srlAPI.SrlNokiaNetworkInstance_NetworkInstance_RouteTable{}
	switch v := (*response.Result).(type) {
	case nil:
		log.Printf("No result in JSON-RPC response")
		wg.Done()
		return
	case []interface{}:
		srlUnmarshal(v[0], routes)
	case interface{}:
		srlUnmarshal(v, routes)
	default:
		log.Printf("unexpected type")
		wg.Done()
		return
	}

	if routes.Ipv4Unicast == nil {
		log.Printf("No local routes found")
		wg.Done()
		return
	}
	out := []string{}
	for key := range routes.Ipv4Unicast.Route {
		out = append(out, key.Ipv4Prefix)
	}

	expectedRoutes := make(map[string]bool)

	for id, loopback := range deviceLoopbacks {
		if id != r.ID {
			expectedRoutes[loopback] = false
		}
	}

	go checkRoutes(r.Hostname, out, expectedRoutes, wg)
}

func (r CEOS) GetRoutes(wg *sync.WaitGroup) {
	log.Print("Collecting CEOS routes")

	client := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	client.SetBaseURL("https://" + r.Hostname + ":6020")
	client.SetBasicAuth(r.Username, r.Password)

	resp, err := client.R().
		SetHeader("Accept", "application/yang-data+json").
		Get(fmt.Sprintf("/restconf/data/network-instances/network-instance=%s/afts", "default"))

	if err != nil {
		log.Printf(
			"failed to send request for %s: %s",
			r.Hostname,
			err.Error(),
		)
		wg.Done()
		return
	}

	response := &eosAPI.NetworkInstance_Afts{}
	if err := eosAPI.Unmarshal(resp.Body(), response); err != nil {
		log.Printf("Failed to unmarshal response :%s", err)
		wg.Done()
		return
	}

	out := []string{}
	for key := range response.Ipv4Entry {
		out = append(out, key)
	}

	expectedRoutes := make(map[string]bool)

	for id, loopback := range deviceLoopbacks {
		if id != r.ID {
			expectedRoutes[loopback] = false
		}
	}

	go checkRoutes(r.Hostname, out, expectedRoutes, wg)
}

func checkRoutes(device string, in []string, expected map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Checking %s routes", device)

	for _, route := range in {
		if _, ok := expected[route]; ok {
			log.Print("Route ", route, " found on ", device)
			expected[route] = true
		}
	}

	for route, found := range expected {
		if !found {
			log.Print("! Route ", route, " NOT found on ", device)
		}
	}
}

func main() {
	srl := SRL{
		Hostname: "clab-netgo-srl",
		ID:       0,
		Authentication: Authentication{
			Username: "admin",
			Password: "admin",
		},
	}
	ceos := CEOS{
		Hostname: "clab-netgo-ceos",
		ID:       1,
		Authentication: Authentication{
			Username: "admin",
			Password: "admin",
		},
	}
	cvx := CVX{
		Hostname: "clab-netgo-cvx",
		ID:       2,
		Authentication: Authentication{
			Username: "cumulus",
			Password: "cumulus",
		},
	}

	devices := []Router{srl, cvx, ceos}

	var wg sync.WaitGroup
	for _, router := range devices {
		wg.Add(1)
		go router.GetRoutes(&wg)
	}
	wg.Wait()
}
