package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"json-rpc/pkg/srl"
	"log"
	"net/http"

	"github.com/openconfig/ygot/ygot"
)

//go:generate go run github.com/openconfig/ygot/generator -path=yang -output_file=pkg/srl/srl.go -package_name=srl yang/interfaces/srl_nokia-interfaces.yang yang/interfaces/srl_nokia-if-ip.yang

var (
	hostname = "http://clab-netgo-srl/jsonrpc"
	username = "admin"
	password = "admin"
)

type RPC struct {
	Version string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	Commands []Command `json:"commands"`
}

type Command struct {
	Action string      `json:"action"`
	Path   string      `json:"path"`
	Value  interface{} `json:"value"`
}

func buildL3Interface(name, prefix string) (Command, error) {
	intf := srl.SrlNokiaInterfaces_Interface{}
	subintf, _ := intf.NewSubinterface(0)

	subintf.Ipv4 = &srl.SrlNokiaInterfaces_Interface_Subinterface_Ipv4{}
	subintf.Ipv4.NewAddress(prefix)

	if err := intf.Validate(); err != nil {
		log.Fatal(err)
	}

	ietf, _ := ygot.ConstructIETFJSON(&intf, nil)

	return Command{
		Action: "replace",
		Path:   fmt.Sprintf("/interface[name=%s]", name),
		Value:  ietf,
	}, nil
}

func buildSetRPC(cmds []Command) RPC {
	return RPC{
		Version: "2.0",
		ID:      0,
		Method:  "set",
		Params: Params{
			Commands: cmds,
		},
	}
}

func main() {

	var cmds []Command

	eth1, err := buildL3Interface("ethernet-1/1", "192.0.2.0/31")
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, eth1)

	lo0, err := buildL3Interface("system0", "198.51.100.0/32")
	if err != nil {
		log.Fatal(err)
	}
	cmds = append(cmds, lo0)

	value, _ := json.Marshal(buildSetRPC(cmds))

	req, err := http.NewRequest("POST", hostname, bytes.NewBuffer(value))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status: %s", resp.Status)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

}
