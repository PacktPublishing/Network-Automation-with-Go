package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	plName      = "ADVERTISE"
	backupRules = map[string][]int{
		"swp1": {10, 20},
	}
)

type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Status string `json:"status"`
	Labels struct {
		Alertname          string `json:"alertname"`
		Instance           string `json:"instance"`
		InterfaceName      string `json:"interface_name"`
		Job                string `json:"job"`
		Severity           string `json:"severity"`
		Source             string `json:"source"`
		SubscriptionName   string `json:"subscription_name"`
		SubscriptionTarget string `json:"subscription_target"`
		Target             string `json:"target"`
	} `json:"labels"`
}

var listenAddr = "0.0.0.0:10000"

type Rule struct {
	Action string `json:"action"`
}

type PrefixList struct {
	Rules map[string]Rule `json:"rule"`
}
type Policy struct {
	PrefixLists map[string]PrefixList `json:"prefix-list"`
}

type Router struct {
	Policy Policy `json:"policy"`
}

type nvue struct {
	Router Router `json:"router"`
}

func toggleBackup(intf string, action string) error {
	log.Printf("%s needs to %s backup prefixes", intf, action)
	ruleIDs, ok := backupRules[intf]
	if !ok {
		log.Println("Could not find a backup prefix for", intf)
		return nil

	}

	var pl PrefixList
	pl.Rules = make(map[string]Rule)
	for _, ruleID := range ruleIDs {
		pl.Rules[strconv.Itoa(ruleID)] = Rule{
			Action: action,
		}
	}

	var payload nvue
	payload.Router.Policy.PrefixLists = map[string]PrefixList{
		plName: pl,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return sendBytes(b)
}

func alertHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("Incoming alert")
	var alerts Alerts

	err := json.NewDecoder(req.Body).Decode(&alerts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, alert := range alerts.Alerts {
		if alert.Status == "firing" {
			if err := toggleBackup(alert.Labels.InterfaceName, "permit"); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			continue
		}
		// alert is resolved
		if err := toggleBackup(alert.Labels.InterfaceName, "deny"); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//fmt.Println("Request body ", alerts)

	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("AlertManager event-triggered webhook")
	http.HandleFunc("/alert", alertHandler)

	log.Println("Starting web server at", listenAddr)
	srv := http.Server{Addr: listenAddr}
	log.Fatal(srv.ListenAndServe())

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
