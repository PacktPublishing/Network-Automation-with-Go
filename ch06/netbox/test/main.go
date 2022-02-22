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
	_, err := nb.Dcim.DcimManufacturersCreate(&dcim.DcimManufacturersCreateParams{
		Context: context.Background(),
		Data:    &vnd,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create manufacturer %s: %w", vnd.Display, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findManufacturer(nb *client.NetBoxAPI, vnd models.Manufacturer) (fnd bool, id int64, err error) {
	rsp, err := nb.Dcim.DcimManufacturersList(&dcim.DcimManufacturersListParams{
		Context: context.Background(),
		SlugIe:  vnd.Slug,
	}, nil)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find manufacturer %s: %w", vnd.Display, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id = rsp.Payload.Results[0].ID
		fmt.Printf("Vendor: %s \tID: %v \n",
			*rsp.Payload.Results[0].Name, id)
	}
	return fnd, id, nil
}

type DeviceTypes struct {
	List []models.DeviceType
}

func createDeviceType(nb *client.NetBoxAPI, dt models.DeviceType) error {
	man := models.Manufacturer{
		Display: dt.Manufacturer.Display,
		Name:    dt.Manufacturer.Name,
		Slug:    dt.Manufacturer.Slug,
	}

	found, id, err := findManufacturer(nb, man)
	if err != nil || !found {
		return fmt.Errorf("error finding manufacturer %s: %w", man.Display, err)
	}

	ndt := models.WritableDeviceType{
		Manufacturer: &id,
		Display:      dt.Display,
		Model:        dt.Model,
		Slug:         dt.Slug,
		Tags:         []*models.NestedTag{},
	}
	f := strfmt.NewFormats()
	err = ndt.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for type %s: %w", *dt.Model, err)
	}

	_, err = nb.Dcim.DcimDeviceTypesCreate(&dcim.DcimDeviceTypesCreateParams{
		Context: context.Background(),
		Data:    &ndt,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create device type %s: %w", *dt.Model, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
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
		id := rsp.Payload.Results[0].ID
		fmt.Printf("Device Type: %q \tID: %v \n",
			strings.TrimSpace(*rsp.Payload.Results[0].Model), id)
	}
	return fnd, nil
}

type DeviceRoles struct {
	List []models.DeviceRole
}

func createDeviceRole(nb *client.NetBoxAPI, dr models.DeviceRole) error {
	_, err := nb.Dcim.DcimDeviceRolesCreate(&dcim.DcimDeviceRolesCreateParams{
		Context: context.Background(),
		Data:    &dr,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create device role %s: %w", dr.Display, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findDeviceRole(nb *client.NetBoxAPI, dr models.DeviceRole) (fnd bool, err error) {
	rsp, err := nb.Dcim.DcimDeviceRolesList(&dcim.DcimDeviceRolesListParams{
		Context: context.Background(),
		SlugIe:  dr.Slug,
	}, nil)
	if err != nil {
		return fnd, fmt.Errorf("failed to find device role %s: %w", *dr.Name, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id := rsp.Payload.Results[0].ID
		fmt.Printf("Site: %q \tID: %v \n",
			strings.TrimSpace(rsp.Payload.Results[0].Display), id)
	}
	return fnd, nil
}

type Sites struct {
	List []models.Site
}

func createSite(nb *client.NetBoxAPI, s models.Site) error {
	ns := models.WritableSite{
		Name:    s.Name,
		Display: s.Display,
		Slug:    s.Slug,
	}
	f := strfmt.NewFormats()
	err := ns.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for site %s: %w", ns.Display, err)
	}

	_, err = nb.Dcim.DcimSitesCreate(&dcim.DcimSitesCreateParams{
		Context: context.Background(),
		Data:    &ns,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create site %s: %w", ns.Display, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findSite(nb *client.NetBoxAPI, s models.Site) (fnd bool, err error) {
	rsp, err := nb.Dcim.DcimSitesList(&dcim.DcimSitesListParams{
		Context: context.Background(),
		SlugIe:  s.Slug,
	}, nil)
	if err != nil {
		return fnd, fmt.Errorf("failed to find site %s: %w", *s.Name, err)
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

	d := json.NewDecoder(man)

	var manInput Manufacturers
	err = d.Decode(&manInput.List)
	if err != nil {
		return fmt.Errorf("cannot decode manufacturers data: %w", err)
	}

	for _, vendor := range manInput.List {
		found, _, err := findManufacturer(nb, vendor)
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
	man.Close()

	////////////////////////////////
	// Device Types
	////////////////////////////////
	dev, err := os.Open("device-types.json")
	if err != nil {
		return fmt.Errorf("cannot open device types file: %w", err)
	}

	d = json.NewDecoder(dev)

	var devTypes DeviceTypes
	err = d.Decode(&devTypes.List)
	if err != nil {
		return fmt.Errorf("cannot decode device types data: %w", err)
	}

	for _, devType := range devTypes.List {
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
	dev.Close()

	////////////////////////////////
	// Device Role
	////////////////////////////////
	rol, err := os.Open("device-roles.json")
	if err != nil {
		return fmt.Errorf("cannot open device roles file: %w", err)
	}

	d = json.NewDecoder(rol)

	var devRoles DeviceRoles
	err = d.Decode(&devRoles.List)
	if err != nil {
		return fmt.Errorf("cannot decode device roles data: %w", err)
	}

	for _, devRole := range devRoles.List {
		found, err := findDeviceRole(nb, devRole)
		if err != nil {
			return fmt.Errorf("error finding device role %s: %w", devRole.Display, err)
		}
		if !found {
			err = createDeviceRole(nb, devRole)
			if err != nil {
				return fmt.Errorf("error creating device role %s: %w", devRole.Display, err)
			}
		}
	}
	dev.Close()

	////////////////////////////////
	// Sites
	////////////////////////////////
	sit, err := os.Open("sites.json")
	if err != nil {
		return fmt.Errorf("cannot open sites file: %w", err)
	}
	defer dev.Close()

	d = json.NewDecoder(sit)

	var devSites Sites
	err = d.Decode(&devSites.List)
	if err != nil {
		return fmt.Errorf("cannot decode sites data: %w", err)
	}

	for _, devSite := range devSites.List {
		found, err := findSite(nb, devSite)
		if err != nil {
			return fmt.Errorf("error finding site %s: %w", devSite.Display, err)
		}
		if !found {
			err = createSite(nb, devSite)
			if err != nil {
				return fmt.Errorf("error creating site %s: %w", devSite.Display, err)
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
