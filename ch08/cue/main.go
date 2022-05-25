package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
)

var (
	inputFilename = "input.yaml"
	inputCueFile  = "input.cue"
	inputPackage  = "input"
	inputPath     = "data"
)

func main() {
	// cue import input.yaml -p input
	if err := importInput(); err != nil {
		log.Fatal(err)
	}

	// cue eval network.automation:cvx -c --out json
	instances := load.Instances([]string{"."}, &load.Config{
		Package: "cvx",
	})

	ctx := cuecontext.New()
	i := ctx.BuildInstance(instances[0])
	if i.Err() != nil {
		msg := errors.Details(i.Err(), nil)
		fmt.Printf("Compile Error:\n%s\n", msg)
	}

	if err := i.Validate(
		cue.Final(),
		cue.Concrete(true),
	); err != nil {
		msg := errors.Details(err, nil)
		fmt.Printf("Validate Error:\n%s\n", msg)
	}

	data, err := i.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(string(data))

	if err := sendBytes(data); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully configured the device")
}

func importInput() error {
	// <cue import input.yaml>
	input, err := os.ReadFile(inputFilename)
	if err != nil {
		return err
	}
	ex, err := yaml.Extract(inputFilename, input)
	if err != nil {
		return err
	}

	// <-p input>
	s := &ast.File{Decls: []ast.Decl{
		&ast.Package{
			Name: ast.NewIdent(inputPackage),
		},
	}}

	// <-l '"data"'>
	//f := &ast.Field{}
	//f.Label = ast.NewString(inputPath)
	//args := make([]interface{}, len(ex.Decls))
	//for i, decl := range ex.Decls {
	//	args[i] = decl
	//}
	//f.Value = ast.NewStruct(args...)
	s.Decls = append(s.Decls, ex.Decls...)

	// dump into binary
	bytes, err := format.Node(s, format.Simplify())
	if err != nil {
		return err
	}

	// write inputCueFile
	return ioutil.WriteFile(inputCueFile, bytes, 0644)
}

type cvx struct {
	url   string
	token string
	httpC http.Client
}

func sendBytes(b []byte) error {
	var (
		hostname    = "clab-netgo-cvx"
		defaultPort = 8765
		username    = "cumulus"
		password    = "cumulus"
	)

	device := cvx{
		url:   fmt.Sprintf("https://%s:%d", hostname, defaultPort),
		token: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))),
		httpC: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	// create a new candidate configuration revision
	revisionID, err := createRevision(device)
	if err != nil {
		return err
	}

	log.Print("Created revisionID: ", revisionID)

	addr, err := url.Parse(device.url + "/nvue_v1/")
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Add("rev", revisionID)
	addr.RawQuery = params.Encode()

	// Save the device desired configuration in candidate configuration store
	req, err := http.NewRequest("PATCH", addr.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+device.token)

	res, err := device.httpC.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Apply candidate revision
	if err := applyRevision(device, revisionID); err != nil {
		log.Fatal(err)
	}

	return nil
}

func createRevision(c cvx) (string, error) {
	revisionPath := "/nvue_v1/revision"
	addr, err := url.Parse(c.url + revisionPath)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", addr.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+c.token)

	res, err := c.httpC.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var response map[string]interface{}

	json.NewDecoder(res.Body).Decode(&response)

	for key := range response {
		return key, nil
	}

	return "", fmt.Errorf("unexpected createRevision error")
}

func applyRevision(c cvx, id string) error {
	applyPath := "/nvue_v1/revision/" + url.PathEscape(id)

	body := []byte("{\"state\": \"apply\", \"auto-prompt\": {\"ays\": \"ays_yes\", \"ignore_fail\": \"ignore_fail_yes\"}} ")

	req, err := http.NewRequest("PATCH", c.url+applyPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+c.token)

	res, err := c.httpC.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)

	return nil
}
