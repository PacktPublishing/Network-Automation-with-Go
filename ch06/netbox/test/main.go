package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-openapi/strfmt"
	resty "github.com/go-resty/resty/v2"
	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
)

var (
	demoNetboxURL = "https://demo.netbox.dev/"
)

type Manufacturers struct {
	List []models.Manufacturer
}

func createManufacturer(nb *client.NetBoxAPI, vnd models.Manufacturer) error {
	crd, err := nb.Dcim.DcimManufacturersCreate(&dcim.DcimManufacturersCreateParams{
		Context: context.Background(),
		Data:    &vnd,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create manufacturer %s: %w", vnd.Display, err)
	}
	fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findManufacturer(nb *client.NetBoxAPI, vnd models.Manufacturer) (fnd bool, err error) {
	rsp, err := nb.Dcim.DcimManufacturersList(&dcim.DcimManufacturersListParams{
		Context: context.Background(),
		SlugIe:  vnd.Slug,
	}, nil)
	if err != nil {
		return fnd, fmt.Errorf("failed to find manufacturer %s: %w", vnd.Display, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		fmt.Printf("Vendor: %s \tID: %v \n",
			*rsp.Payload.Results[0].Name, rsp.Payload.Results[0].ID)
	}
	return fnd, nil
}

type DeviceTypes struct {
	List []models.DeviceType
}

func createDeviceType(nb *client.NetBoxAPI, dt models.DeviceType) error {
	ndt := models.WritableDeviceType{
		Manufacturer: &dt.Manufacturer.ID,
		ID:           dt.ID,
		Display:      dt.Display,
		Model:        dt.Model,
		Slug:         dt.Slug,
		Tags:         []*models.NestedTag{},
	}
	f := strfmt.NewFormats()
	err := ndt.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for type %s: %w", *dt.Model, err)
	}

	crd, err := nb.Dcim.DcimDeviceTypesCreate(&dcim.DcimDeviceTypesCreateParams{
		Context: context.Background(),
		Data:    &ndt,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create device type %s: %w", *dt.Model, err)
	}
	fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findDeviceType(nb *client.NetBoxAPI, dt models.DeviceType) (fnd bool, err error) {
	rsp, err := nb.Dcim.DcimDeviceTypesList(&dcim.DcimDeviceTypesListParams{
		Context: context.Background(),
		Model:   dt.Model,
	}, nil)
	if err != nil {
		return fnd, fmt.Errorf("failed to find device type %s: %w", *dt.Model, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		fmt.Printf("Device Type: %q \tID: %v \n",
			strings.TrimSpace(*rsp.Payload.Results[0].Model), rsp.Payload.Results[0].ID)
	}
	return fnd, nil
}

type DeviceRoles struct {
	List []models.DeviceRole
}

func createDeviceRole(nb *client.NetBoxAPI, dr models.DeviceRole) error {
	ndr := models.DeviceRole{
		ID:           dr.ID,
		Display:      dr.Display,
		Slug:         dr.Slug,
	}
	f := strfmt.NewFormats()
	err := ndr.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for role %s: %w", dr.Display, err)
	}

	crd, err := nb.Dcim.DcimDeviceRolesCreate(&dcim.DcimDeviceRolesCreateParams{
		Context: context.Background(),
		Data:    &ndr,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create device role %s: %w", dr.Display, err)
	}
	fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findDeviceRole(nb *client.NetBoxAPI, dr models.DeviceRole) (fnd bool, err error) {
	rsp, err := nb.Dcim.DcimDeviceRolesList(&dcim.DcimDeviceRolesListParams{
		Context: context.Background(),
		NameIe:   dr.Name,
	}, nil)
	if err != nil {
		return fnd, fmt.Errorf("failed to find device role %s: %w", *dr.Name, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		fmt.Printf("Device Role: %q \tID: %v \n",
			strings.TrimSpace(rsp.Payload.Results[0].Display), rsp.Payload.Results[0].ID)
	}
	return fnd, nil
}

func createToken(usr, pwd string, url *url.URL) (string, error) {
	client := resty.New()
	client.SetBaseURL("https://" + url.Host)

	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, usr, pwd)

	result := make(map[string]interface{})
	_, err := client.R().
		SetResult(&result).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/users/tokens/provision/")

	if err != nil {
		return "", fmt.Errorf("error requesting a token: %w", err)
	}

	if val, ok := result["key"]; ok {
		return val.(string), nil
	}

	return "", fmt.Errorf("empty token")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func createResources(nb *client.NetBoxAPI) error {
	////////////////////////////////
	// Manufacturers
	////////////////////////////////
	man, err := os.Open("manufacturer.json")
	if err != nil {
		return fmt.Errorf("cannot open manufacturers file: %w", err)
	}
	defer man.Close()

	d1 := json.NewDecoder(man)

	var manInput Manufacturers
	err = d1.Decode(&manInput.List)
	if err != nil {
		return fmt.Errorf("cannot decode manufacturers data: %w", err)
	}

	for _, vendor := range manInput.List {
		found, err := findManufacturer(nb, vendor)
		if err != nil {
			return fmt.Errorf("error finding manufacturer %s: %w", vendor.Display, err)
		}
		if !found {
			err = createManufacturer(nb, vendor)
			if err != nil {
				return fmt.Errorf("error creating manufacturer %s: %w", vendor.Display, err)
			}
		}
	}
	////////////////////////////////
	// Device Types
	////////////////////////////////
	dev, err := os.Open("device-types.json")
	if err != nil {
		return fmt.Errorf("cannot open device types file: %w", err)
	}
	defer dev.Close()

	d2 := json.NewDecoder(dev)

	var devInput DeviceTypes
	err = d2.Decode(&devInput.List)
	if err != nil {
		return fmt.Errorf("cannot decode device types data: %w", err)
	}

	for _, devType := range devInput.List {
		found, err := findDeviceType(nb, devType)
		if err != nil {
			return fmt.Errorf("error finding device type %s: %w", devType.Display, err)
		}
		if !found {
			err = createDeviceType(nb, devType)
			if err != nil {
				return fmt.Errorf("error creating device type %s: %w", devType.Display, err)
			}
		}
	}

	return nil
}

func main() {
	nbURL := flag.String("netbox", demoNetboxURL, "Netbox URL")
	username := flag.String("username", "admin", "admin username")
	password := flag.String("password", "admin", "admin password")
	flag.Parse()

	url, err := url.Parse(*nbURL)
	check(err)

	client.DefaultSchemes = []string{url.Scheme}

	token, err := createToken(*username, *password, url)
	check(err)

	nbClient := netbox.NewNetboxWithAPIKey(url.Host, token)

	err = createResources(nbClient)
	check(err)

}
