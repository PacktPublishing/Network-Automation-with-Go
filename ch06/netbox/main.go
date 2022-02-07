package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"

	resty "github.com/go-resty/resty/v2"
)

var (
	demoNetboxURL = "https://demo.netbox.dev/"
)

func createToken(username, password string, url *url.URL) (string, error) {
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("https://%s", url.Host))

	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)

	response := make(map[string]interface{})
	resp, err := client.R().
		SetResult(&response).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/users/tokens/provision/")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response %+v", resp)

	// super unsafe
	token := response["key"].(string)

	return token, nil
}

func generateUniqueName() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d", rand.Int())
}

func main() {
	nbUrl := flag.String("netbox", demoNetboxURL, "Netbox URL")
	username := flag.String("username", "admin", "admin username")
	password := flag.String("password", "admin", "admin password")
	flag.Parse()

	url, err := url.Parse(*nbUrl)
	if err != nil {
		log.Fatal(err)
	}

	client.DefaultSchemes = []string{url.Scheme}

	token, err := createToken(*username, *password, url)
	if err != nil {
		log.Fatal(err)
	}

	nbClient := netbox.NewNetboxWithAPIKey(url.Host, token)

	var fakeID int64 = 1
	demoDeviceName := generateUniqueName()

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
