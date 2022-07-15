package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"os"
	"encoding/json"

	resty "github.com/go-resty/resty/v2"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
)

type Authentication struct {
	Username string
	Password string
}

type CVX struct {
	Hostname string
	Authentication
	Resp *Response
}

type SRL struct {
	Hostname string
	Authentication
	Resp *Response
}

type CEOS struct {
	Hostname string
	Authentication
	Resp *Response
}

type Router interface {
	GetRoutes(wg *sync.WaitGroup)
}

func (r CVX) GetRoutes(wg *sync.WaitGroup) {
	client := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	client.SetBaseURL("https://" + r.Hostname + ":8765" )
	client.SetBasicAuth(r.Username, r.Password)

	var routes map[string]interface{}
	_, err := client.R().
		SetResult(&routes).
		SetQueryParams(map[string]string{
			"rev": "operational",
		}).
		Get("/nvue_v1/vrf/default/router/rib/ipv4/route")

	if err != nil {
		r.Resp.check(err, "failed to send request for: "+r.Hostname)
		return
	}

	out := []string{}
	for route := range routes {
		out = append(out, route)
	}
	go checkRoutes(r.Hostname, out, wg, r.Resp)
}

func (r SRL) GetRoutes(wg *sync.WaitGroup) {
	lookupCmd := "show network-instance default route-table ipv4-unicast summary"

	conn, err := platform.NewPlatform(
		"nokia_srl",
		r.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(r.Username),
		options.WithAuthPassword(r.Password),
		options.WithTermWidth(176),
	)
	if err != nil {
		r.Resp.check(err, "failed to create connection for: "+r.Hostname)
		return
	}
	driver, err := conn.GetNetworkDriver()
	if err != nil {
		r.Resp.check(err, "failed to create driver for: "+r.Hostname)
		return
	}

	err = driver.Open()
	if err != nil {
		r.Resp.check(err, "failed to open driver for: "+r.Hostname)
		return
	}	
	defer driver.Close()

	resp, err := driver.SendCommand(lookupCmd)
	if err != nil {
		r.Resp.check(err, "failed to send command for: "+r.Hostname)
		return
	}

	ipv4Prefix := regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}\/\d{1,2}`)

	out := []string{}
	for _, match := range ipv4Prefix.FindAll(resp.RawResult, -1) {
		out = append(out, string(match))
	}
	go checkRoutes(r.Hostname, out, wg, r.Resp)
}

func (r CEOS) GetRoutes(wg *sync.WaitGroup) {
	template := "https://raw.githubusercontent.com/networktocode/ntc-templates/master/ntc_templates/templates/arista_eos_show_ip_route.textfsm"

	lookupCmd := "sh ip route"

	conn, err := platform.NewPlatform(
		"arista_eos",
		r.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(r.Username),
		options.WithAuthPassword(r.Password),
	)
	if err != nil {
		r.Resp.check(err, "failed to create connection for: "+r.Hostname)
		return
	}
	driver, err := conn.GetNetworkDriver()
	if err != nil {
		r.Resp.check(err, "failed to create driver for: "+r.Hostname)
		return
	}

	err = driver.Open()
	if err != nil {
		r.Resp.check(err, "failed to open driver for: "+r.Hostname)
		return
	}
	defer driver.Close()

	resp, err := driver.SendCommand(lookupCmd)
	if err != nil {
		r.Resp.check(err, "failed to send command for: "+r.Hostname)
		return
	}

	parsed, err := resp.TextFsmParse(template)
	if err != nil {
		r.Resp.check(err, "failed to parse command for: "+r.Hostname)
		return
	}

	out := []string{}
	for _, match := range parsed {
		out = append(out, fmt.Sprintf("%s/%s", match["NETWORK"], match["MASK"]))
	}
	go checkRoutes(r.Hostname, out, wg, r.Resp)
}

func checkRoutes(device string, in []string, wg *sync.WaitGroup, r *Response) {
	defer wg.Done()
	m.Lock()
	r.Msg += fmt.Sprintf("Checking %s routes\n", device)
	m.Unlock()

	expectedRoutes := map[string]bool{
		"198.51.100.0/32": false,
		"198.51.100.1/32": false,
		"198.51.100.2/32": false,
	}

	for _, route := range in {
		if _, ok := expectedRoutes[route]; ok {
			m.Lock()
			r.Msg += fmt.Sprintln("Route ", route, " found on ", device)
			m.Unlock()
			expectedRoutes[route] = true
		}
	}

	for route, found := range expectedRoutes {
		if !found {
			m.Lock()
			r.Msg += fmt.Sprintln("! Route ", route, " NOT found on ", device)
			m.Unlock()
			r.Failed = true
		}
	}
}

// ModuleArgs are the module inputs
type ModuleArgs struct {
	Host     string
	User     string
	Password string
	Input    string
}

// Response are the values returned from the module
type Response struct {
	Msg     string `json:"msg"`
	Busy    bool   `json:"busy"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

// ExitJSON is ...
func ExitJSON(responseBody Response) {
	returnResponse(responseBody)
}

// FailJSON is ...
func FailJSON(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(r Response) {
	var response []byte
	var err error
	response, err = json.Marshal(r)
	if err != nil {
		response, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if r.Failed {
		os.Exit(1)
	}
	os.Exit(0)
}

func (r Response) check(err error, msg string) {
	if err != nil {
		m.Lock()
		r.Msg = msg + ": " + err.Error()
		m.Unlock()
		FailJSON(r)
	}
}

var m sync.RWMutex = sync.RWMutex{}

func main() {
	var r Response

	if len(os.Args) != 2 {
		r.Msg = "No argument file provided"
		FailJSON(r)
	}

	argsFile := os.Args[1]

	text, err := os.ReadFile(argsFile)
	r.check(err, "Could not read configuration file: "+argsFile)

	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)
	r.check(err, "Ansible inputs are not valid (JSON): "+argsFile)

	////////////////////////////////
	// Devices
	////////////////////////////////
	cvx := CVX{
		Hostname: "clab-netgo-cvx",
		Authentication: Authentication{
			Username: "cumulus",
			Password: "cumulus",
		},
		Resp: &r,
	}
	srl := SRL{
		Hostname: "clab-netgo-srl",
		Authentication: Authentication{
			Username: "admin",
			Password: "admin",
		},
		Resp: &r,
	}
	ceos := CEOS{
		Hostname: "clab-netgo-ceos",
		Authentication: Authentication{
			Username: "admin",
			Password: "admin",
		},
		Resp: &r,
	}

	////////////////////////////////
	// Get routes & validate them
	////////////////////////////////

	devices := []Router{cvx, srl, ceos}

	var wg sync.WaitGroup
	for _, router := range devices {
		wg.Add(1)
		go router.GetRoutes(&wg)
	}
	wg.Wait()

	returnResponse(r)
}
