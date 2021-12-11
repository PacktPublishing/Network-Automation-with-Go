package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	server := flag.String("server", "127.0.0.1:8080", "HTTP server URL")
	check := flag.Bool("check", false, "healthcheck flag")
	lookup := flag.String("lookup", "domain", "lookup data [mac, ip, domain]")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Println("must provide exactly one query argyment")
		return
	}

	path := "/lookup"
	if *check {
		path = "/check"
	}

	url := fmt.Sprintf("http://%s%s", *server, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add(*lookup, flag.Arg(0))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

}
