package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"strings"
)

const wiresharkDB = "https://gitlab.com/wireshark/wireshark/-/raw/master/manuf"

var macDB map[string]string

func init() {
	var err error
	macDB, err = download()
	if err != nil {
		log.Fatalf("Failed to download mac DB: %s", err)
	}
	log.Printf("macDB initialized")
}

func parse(db io.Reader, out map[string]string) map[string]string {
	lineScanner := bufio.NewScanner(db)
	for lineScanner.Scan() {
		if len(lineScanner.Bytes()) < 1 {
			continue
		}
		if lineScanner.Bytes()[0] == '#' {
			continue
		}

		parts := strings.Split(lineScanner.Text(), "\t")

		if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
			continue
		}

		out[parts[0]] = parts[2]
	}

	if err := lineScanner.Err(); err != nil {
		return out
	}
	return out
}

func download() (map[string]string, error) {
	result := make(map[string]string)

	resp, err := http.Get(wiresharkDB)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	parse(resp.Body, result)

	return result, nil
}
