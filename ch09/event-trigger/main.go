package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func alertHandler(w http.ResponseWriter, req *http.Request) {

	log.Printf("Incoming alert for \n %+v", req.URL)
	alerts := new(Alerts)

	err := json.NewDecoder(req.Body).Decode(alerts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Request body ", alerts)

	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	fmt.Println("Alert manager event trigger webhook")
	http.HandleFunc("/alert", alertHandler)

	log.Println("Starting web server at", listenAddr)
	srv := http.Server{Addr: listenAddr}
	log.Fatal(srv.ListenAndServe())

}
