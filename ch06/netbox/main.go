package main

import (
	"context"
	"flag"
	"log"
	"net/url"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
)

var (
	demoNetboxURL   = "https://demo.netbox.dev/"
	demoNetboxToken = "0123456789abcdef0123456789abcdef01234567"
	demoDeviceName  = "go-automation"
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

	var fakeID int64 = 1
	created, err := nbClient.Dcim.DcimDevicesCreate(&dcim.DcimDevicesCreateParams{
		Context: context.Background(),
		Data: &models.WritableDeviceWithConfigContext{
			Name:       &demoDeviceName,
			DeviceRole: &fakeID,
			DeviceType: &fakeID,
			Site:       &fakeID,
			Tags:       []*models.NestedTag{},
		},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := nbClient.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
		ID:      created.Payload.ID,
		Context: context.Background(),
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Created device ", *res.Payload.Name)

}
