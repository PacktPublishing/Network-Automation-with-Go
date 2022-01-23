package main

import (
	"context"
	"flag"
	"log"
	"net/url"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
)

var (
	demoNetboxURL   = "https://demo.netbox.dev/"
	demoNetboxToken = "0123456789abcdef0123456789abcdef01234567"
)

func main() {
	nbUrl := flag.String("netbox", demoNetboxURL, "Netbox URL")
	token := flag.String("token", demoNetboxToken, "Token")
	flag.Parse()

	url, err := url.Parse(*nbUrl)
	if err != nil {
		log.Fatal(err)
	}

	client.DefaultSchemes = []string{url.Scheme}

	nbClient := netbox.NewNetboxWithAPIKey(url.Host, *token)

	res, err := nbClient.Dcim.DcimDevicesList(&dcim.DcimDevicesListParams{
		Context: context.Background(),
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range res.Payload.Results {
		log.Print("device ", *device.Name)
	}

}
