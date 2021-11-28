package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	whoisIANA = "whois.iana.org"

	whoisPort = 43
)

func whoisLookup(query, server string) (*bytes.Buffer, error) {
	log.Printf("whoisLookup %s@%s", query, server)

	server = fmt.Sprintf("%s:%d", server, whoisPort)

	rAddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr failed: %s", err)
	}

	conn, err := net.DialTCP("tcp4", nil, rAddr)
	if err != nil {
		return nil, fmt.Errorf("DialTCP failed: %s", err)
	}
	defer conn.Close()

	// all queries must end with CRLF
	query += "\r\n"
	_, err = conn.Write([]byte(query))
	if err != nil {
		return nil, fmt.Errorf("Write failed: %s", err)
	}

	var response bytes.Buffer

	_, err = io.Copy(&response, conn)
	if err != nil {
		return nil, fmt.Errorf("Read failed: %s", err)
	}

	return &response, nil
}

func findRefer(input *bytes.Buffer) (string, bool) {
	lineScanner := bufio.NewScanner(input)
	for lineScanner.Scan() {
		if len(lineScanner.Bytes()) < 1 {
			continue
		}
		if lineScanner.Bytes()[0] == '#' {
			continue
		}

		line := lineScanner.Text()

		if strings.Contains(line, "refer: ") {
			parts := strings.Fields(line)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return "", false
			}
			return parts[1], true
		}

	}

	if err := lineScanner.Err(); err != nil {
		return "", false
	}

	return "", false
}

func main() {
	query := flag.String("query", "", "whois query")
	flag.Parse()

	if *query == "" {
		log.Fatal("Please provide -query flag")
	}

	if ip := net.ParseIP(*query); ip != nil {
		log.Println("query recognized as IP")
		//*query = "n + " + *query
	}

	whoisServer := whoisIANA
	for {
		response, err := whoisLookup(*query, whoisServer)
		if err != nil {
			log.Fatalf("lookup failed: %s", err)
		}

		answer := response.String()

		refer, found := findRefer(response)
		if found {
			whoisServer = refer
			continue
		}

		log.Println(answer)
		break
	}

}
