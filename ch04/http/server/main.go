package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func getWhois(input []string) string {
	if len(input) != 1 {
		return fmt.Sprintf("incorrect query %v", input)
	}

	var res string
	query := input[0]
	whoisServer := whoisIANA
	for {
		response, err := whoisLookup(query, whoisServer)
		if err != nil {
			log.Fatalf("lookup failed: %s", err)
		}

		res = response.String()

		refer, found := findRefer(response)
		if found {
			whoisServer = refer
			continue
		}

		break
	}
	return res
}

func getMAC(input []string) string {
	if len(input) != 1 {
		return fmt.Sprintf("incorrect query %v", input)
	}

	mac, err := net.ParseMAC(input[0])
	if err != nil {
		return fmt.Sprintf("Failed to parse MAC")
	}

	oui := mac[:3].String()
	oui = strings.ToUpper(oui)

	res, ok := macDB[oui]
	if !ok {
		return fmt.Sprintf("result not found\n")
	}
	return res
}

func lookup(w http.ResponseWriter, req *http.Request) {

	log.Printf("Incoming %+v", req.URL.Query())
	var response string

	for k, v := range req.URL.Query() {
		switch k {
		case "ip":
			response = getWhois(v)
		case "mac":
			response = getMAC(v)
		case "domain":
			response = getWhois(v)
		default:
			response = fmt.Sprintf("query %q not recognized", k)
		}

	}
	fmt.Fprintf(w, response)
}

func check(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

// Run with 'go run *.go'
func main() {
	http.HandleFunc("/lookup", lookup)

	http.HandleFunc("/check", check)

	log.Println("Starting web server at 0.0.0.0:8080")
	//srv := http.Server{Addr: "0.0.0.0:8080", Handler: nil}
	srv := http.Server{Addr: "0.0.0.0:8080"}
	log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(":8080", nil))

}
