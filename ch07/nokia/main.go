package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/karimra/gnmic/api"
	"github.com/karimra/gnmic/target"
	"google.golang.org/protobuf/encoding/prototext"
	"gopkg.in/yaml.v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Model struct {
	Uplinks  []Link `yaml:"uplinks"`
	Peers    []Peer `yaml:"peers"`
	ASN      int    `yaml:"asn"`
	Loopback Addr   `yaml:"loopback"`
}

type Link struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

type Peer struct {
	IP  string `yaml:"ip"`
	ASN int    `yaml:"asn"`
}

type Addr struct {
	IP string `yaml:"ip"`
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
		r.Msg = msg + ": " + err.Error()
		FailJSON(r)
	}
}

type Data struct {
	Prefix   string `yaml:"prefix,omitempty"`
	Path     string `yaml:"path"`
	Encoding string `yaml:"encoding,omitempty"`
	Value    string `yaml:"value,omitempty"`
}

func (r ModuleArgs) createTarget() (*target.Target, error) {
	return api.NewTarget(
		api.Name("gnmi"),
		api.Address(r.Host+":"+"57400"),
		api.Username(r.User),
		api.Password(r.Password),
		api.SkipVerify(true),
	)
}

// export ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.18
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

	src, err := base64.StdEncoding.DecodeString(moduleArgs.Input)
	r.check(err, "Couldn't decode the configuration inputs file: "+moduleArgs.Input)
	reader := bytes.NewReader(src)

	d := yaml.NewDecoder(reader)

	var input Model
	err = d.Decode(&input)
	r.check(err, "Couldn't decode configuration inputs: "+string(src))

	////////////////////////////////
	// Create a target
	////////////////////////////////
	tg, err := moduleArgs.createTarget()
	check(err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	////////////////////////////////
	// Create a gNMI client
	////////////////////////////////
	err = tg.CreateGNMIClient(ctx)
	check(err)
	defer tg.Close()

	////////////////////////////////
	// Send a gNMI capabilities request to the created target
	////////////////////////////////
	// capResp, err := tg.Capabilities(ctx)
	// check(err)

	// fmt.Println(prototext.Format(capResp))

	////////////////////////////////
	// Read input data for gNMI request
	////////////////////////////////
	gdata, err := os.Open("api-srl.tpl")
	check(err)
	defer gdata.Close()

	d = yaml.NewDecoder(gdata)

	var info []Data
	err = d.Decode(&info)
	check(err)

	var pathBuffer bytes.Buffer
	var valueBuffer bytes.Buffer

	for _, data := range info {
		////////////////////////////////
		// Create a GetRequest
		////////////////////////////////
		// getReq, err := api.NewGetRequest(
		// 	api.Path(data.Prefix+data.Path),
		// 	api.Encoding("json_ietf"))
		// check(err)

		// fmt.Println(prototext.Format(getReq))

		////////////////////////////////
		// Send the created gNMI GetRequest to the created target
		////////////////////////////////
		// getResp, err := tg.Get(ctx, getReq)
		// check(err)

		// fmt.Println(prototext.Format(getResp))

		////////////////////////////////
		// Create an Update gNMI SetRequest
		////////////////////////////////
		// setReq, err := api.NewSetRequest(
		// 	api.Update(
		// 		api.Path(data.Prefix+data.Path),
		// 		api.Value(data.Value, "json_ietf")),
		// )

		////////////////////////////////
		// Create a Delete gNMI SetRequest
		////////////////////////////////
		// clean, err := api.NewSetRequest(
		// 	api.Delete(data.Prefix + data.Path))

		// check(err)
		// fmt.Println(prototext.Format(clean))

		////////////////////////////////
		// Send the Delete gNMI SetRequest to the target
		////////////////////////////////
		// cleanResp, err := tg.Set(ctx, clean)
		// check(err)

		// fmt.Println(prototext.Format(cleanResp))

		////////////////////////////////
		// Template Path
		////////////////////////////////
		dt, err := template.New("path").Parse(data.Path)
		check(err)
		err = dt.Execute(&pathBuffer, input)
		check(err)

		////////////////////////////////
		// Template Value to assign
		////////////////////////////////

		vt, err := template.New("value").Parse(data.Value)
		check(err)

		err = vt.Execute(&valueBuffer, input)
		check(err)

		////////////////////////////////
		// Create a Replace gNMI SetRequest
		////////////////////////////////
		setReq, err := api.NewSetRequest(
			api.Replace(
				api.Path(data.Prefix+pathBuffer.String()),
				api.Value(valueBuffer.String(), "json_ietf")),
		)

		check(err)
		fmt.Println(prototext.Format(setReq))

		////////////////////////////////
		// Send the Replace gNMI SetRequest to the target
		////////////////////////////////
		configResp, err := tg.Set(ctx, setReq)
		check(err)

		fmt.Println(prototext.Format(configResp))
		pathBuffer.Reset()
		valueBuffer.Reset()
	}
}
