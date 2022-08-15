package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

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
}

type SRL struct {
	Hostname string
	Authentication
}

type CEOS struct {
	Hostname string
	Authentication
}

type Router interface {
	GetRoutes(wg *sync.WaitGroup)
}

func (r CVX) GetRoutes(wg *sync.WaitGroup) {
	log.Print("Collecting CVX routes")

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
		log.Printf("failed to send request for %s: %s", r.Hostname, err.Error())
		return
	}

	out := []string{}
	for route := range routes {
		out = append(out, route)
	}
	go checkRoutes(r.Hostname, out, wg)
}

func (r SRL) GetRoutes(wg *sync.WaitGroup) {
	log.Print("Collecting SRL routes")

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
		log.Printf("failed to create platform for %s: %s", r.Hostname, err.Error())
		return
	}
	driver, err := conn.GetNetworkDriver()
	if err != nil {
		log.Printf("failed to create driver for %s: %s", r.Hostname, err.Error())
		return
	}
	err = driver.Open()
	if err != nil {
		log.Printf("failed to open driver for %s: %s", r.Hostname, err.Error())
		return
	}
	defer driver.Close()

	resp, err := driver.SendCommand(lookupCmd)
	if err != nil {
		log.Printf("failed to send command for %s: %s", r.Hostname, err.Error())
		return
	}

	ipv4Prefix := regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}\/\d{1,2}`)

	out := []string{}
	for _, match := range ipv4Prefix.FindAll(resp.RawResult, -1) {
		out = append(out, string(match))
	}
	go checkRoutes(r.Hostname, out, wg)
}

func (r CEOS) GetRoutes(wg *sync.WaitGroup) {
	log.Print("Collecting CEOS routes")

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
		log.Printf("failed to create platform for %s: %s", r.Hostname, err.Error())
		return
	}
	driver, err := conn.GetNetworkDriver()
	if err != nil {
		log.Printf("failed to create driver for %s: %s", r.Hostname, err.Error())
		return
	}
	err = driver.Open()
	if err != nil {
		log.Printf("failed to open driver for %s: %s", r.Hostname, err.Error())
		return
	}
	defer driver.Close()

	resp, err := driver.SendCommand(lookupCmd)
	if err != nil {
		log.Printf("failed to send command for %s: %s", r.Hostname, err.Error())
		return
	}

	parsed, err := resp.TextFsmParse(template)
	if err != nil {
		log.Printf("failed to parse command for %s: %s", r.Hostname, err.Error())
		return
	}

	out := []string{}
	for _, match := range parsed {
		out = append(out, fmt.Sprintf("%s/%s", match["NETWORK"], match["MASK"]))
	}
	go checkRoutes(r.Hostname, out, wg)
}

func checkRoutes(device string, in []string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Checking %s routes", device)

	expectedRoutes := map[string]bool{
		"198.51.100.0/32": false,
		"198.51.100.1/32": false,
		"198.51.100.2/32": false,
	}

	for _, route := range in {
		if _, ok := expectedRoutes[route]; ok {
			log.Print("Route ", route, " found on ", device)
			expectedRoutes[route] = true
		}
	}

	for route, found := range expectedRoutes {
		if !found {
			log.Print("! Route ", route, " NOT found on ", device)
		}
	}
}

func main() {
	////////////////////////////////
	// Devices
	////////////////////////////////
	cvx := CVX{
		Hostname: "clab-netgo-cvx",
		Authentication: Authentication{
			Username: "cumulus",
			Password: "cumulus",
		},
	}
	srl := SRL{
		Hostname: "clab-netgo-srl",
		Authentication: Authentication{
			Username: "admin",
			Password: "admin",
		},
	}
	ceos := CEOS{
		Hostname: "clab-netgo-ceos",
		Authentication: Authentication{
			Username: "admin",
			Password: "admin",
		},
	}

	log.Printf("Checking reachability...")

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
}