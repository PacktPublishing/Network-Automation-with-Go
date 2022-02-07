package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	resty "github.com/go-resty/resty/v2"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
)

func getCVXRoutes(hostname, username, password string, out chan string) {
	log.Print("Collecting CVX routes")

	client := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	client.SetBaseURL(fmt.Sprintf("https://%s:8765", hostname))
	client.SetBasicAuth(username, password)

	var routes map[string]interface{}
	_, err := client.R().
		SetResult(&routes).
		SetQueryParams(map[string]string{
			"rev": "operational",
		}).
		Get("/nvue_v1/vrf/default/router/rib/ipv4/route")

	if err != nil {
		log.Fatal(err)
	}

	for route := range routes {
		out <- route
	}

	close(out)
}

func getSRLRoutes(hostname, username, password string, out chan string) {
	log.Print("Collecting SRL routes")

	lookupCmd := "show network-instance default route-table ipv4-unicast summary"

	conn, err := core.NewSROSClassicDriver(
		hostname,
		base.WithAuthStrictKey(false),
		base.WithAuthUsername(username),
		base.WithAuthPassword(password),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := conn.SendCommand(lookupCmd)
	if err != nil {
		log.Fatal(err)
	}

	ipv4Prefix := regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}\/\d{1,2}`)

	for _, match := range ipv4Prefix.FindAll(resp.RawResult, -1) {
		out <- string(match)
	}

	close(out)
}

func getCEOSRoutes(hostname, username, password string, out chan string) {
	log.Print("Collecting CEOS routes")

	template := "https://raw.githubusercontent.com/networktocode/ntc-templates/master/ntc_templates/templates/arista_eos_show_ip_route.textfsm"

	lookupCmd := "sh ip route"

	conn, err := core.NewEOSDriver(
		hostname,
		base.WithAuthStrictKey(false),
		base.WithAuthUsername(username),
		base.WithAuthPassword(password),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := conn.SendCommand(lookupCmd)
	if err != nil {
		log.Fatal(err)
	}

	parsed, err := resp.TextFsmParse(template)
	if err != nil {
		log.Fatal(err)
	}

	for _, match := range parsed {
		out <- fmt.Sprintf("%s/%s", match["NETWORK"], match["MASK"])
	}

	close(out)
}

func checkRoutes(device string, in chan string) {
	expectedRoutes := map[string]struct{}{
		"198.51.100.0/32": struct{}{},
		"198.51.100.1/32": struct{}{},
		"198.51.100.2/32": struct{}{},
	}

	for route := range in {
		if _, ok := expectedRoutes[route]; ok {
			log.Print("Route ", route, " found on ", device)
			delete(expectedRoutes, route)
		}
	}

	for left := range expectedRoutes {
		log.Print("! Route ", left, " NOT found on ", device)
	}
}

func main() {
	cvxHost := flag.String("cvx-host", "clab-netgo-cvx", "CVX Hostname")
	cvxUser := flag.String("cvx-user", "cumulus", "CVX Username")
	cvxPass := flag.String("cvx-pass", "cumulus", "CVX password")
	srlHost := flag.String("srl-host", "clab-netgo-srl", "SRL Hostname")
	srlUser := flag.String("srl-user", "admin", "SRL Username")
	srlPass := flag.String("srl-pass", "admin", "SRL password")
	ceosHost := flag.String("ceos-host", "clab-netgo-ceos", "CEOS Hostname")
	ceosUser := flag.String("ceos-user", "admin", "CEOS Username")
	ceosPass := flag.String("ceos-pass", "admin", "CEOS password")
	flag.Parse()

	ceosRouteCh := make(chan string)
	cvxRouteCh := make(chan string)
	srlRouteCh := make(chan string)

	log.Printf("Checking reachability...")

	go getCEOSRoutes(*ceosHost, *ceosUser, *ceosPass, ceosRouteCh)
	go getCVXRoutes(*cvxHost, *cvxUser, *cvxPass, cvxRouteCh)
	go getSRLRoutes(*srlHost, *srlUser, *srlPass, srlRouteCh)

	checkRoutes("ceos", ceosRouteCh)
	checkRoutes("cvx", cvxRouteCh)
	checkRoutes("srl", srlRouteCh)

}
