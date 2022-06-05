package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vishvananda/netlink"
)

var backupInterfaces = map[string]string{
	"swp1": "swp2",
}

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

func toggleInterface(name string, state netlink.LinkOperState) error {
	log.Printf("Request to change link %s to %v", name, state)
	bkpIntf, ok := backupInterfaces[name]
	if !ok {
		log.Println("Could not find a backup interface for", name)
		return nil

	}

	intf, err := netlink.LinkByName(bkpIntf)
	if err != nil {
		return err
	}

	if intf.Attrs().OperState != state {
		if state == netlink.OperUp {
			log.Printf("ip link %s set up", bkpIntf)
			if err := netlink.LinkSetUp(intf); err != nil {
				return err
			}
			return nil
		}
		log.Printf("ip link %s set down", bkpIntf)
		if err := netlink.LinkSetDown(intf); err != nil {
			return err
		}
		return nil
	}
	log.Println("Links state is the same as expected")
	return nil
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
			if err := toggleInterface(alert.Labels.InterfaceName, netlink.OperUp); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			continue
		}
		// alert is resolved
		if err := toggleInterface(alert.Labels.InterfaceName, netlink.OperDown); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//fmt.Println("Request body ", alerts)

	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("Alert manager event trigger webhook")
	http.HandleFunc("/alert", alertHandler)

	log.Println("Starting web server at", listenAddr)
	srv := http.Server{Addr: listenAddr}
	log.Fatal(srv.ListenAndServe())

}
