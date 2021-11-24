package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

const wiresharkDB = "https://gitlab.com/wireshark/wireshark/-/raw/master/manuf"

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

func lookup(mac net.HardwareAddr) (string, error) {
	db, err := download()
	if err != nil {
		return "", err
	}

	oui := mac[:3].String()
	oui = strings.ToUpper(oui)

	result, ok := db[oui]
	if !ok {
		return "", fmt.Errorf("OUI %s not found in the DB", oui)
	}

	return result, nil
}

func main() {
	macStr := flag.String("mac", "", "MAC address")
	flag.Parse()

	if *macStr == "" {
		log.Fatal("Please provide -mac flag")
	}

	mac, err := net.ParseMAC(*macStr)
	if err != nil {
		log.Fatalf("Failed to parse MAC: %s", err)
	}

	log.Printf("MAC: %s", mac.String())

	vendor, err := lookup(mac)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Vendor: %s", vendor)
}
