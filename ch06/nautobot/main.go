package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/types"
	nb "github.com/nautobot/go-nautobot"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func NewSecurityProviderNautobotToken(t string) (*SecurityProviderNautobotToken, error) {
	return &SecurityProviderNautobotToken{
		token: t,
	}, nil
}

type SecurityProviderNautobotToken struct {
	token string
}

func (s *SecurityProviderNautobotToken) Intercept(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.token))
	return nil
}

func getDeviceIDs(n *nb.ClientWithResponses, in nb.Device) (fnd bool, out *nb.WritableDeviceWithConfigContext, err error) {
	rsp, err := n.DcimDevicesListWithResponse(
		context.TODO(),
		&nb.DcimDevicesListParams{
			NameIe: &[]string{*in.Name},
		})
	if err != nil {
		return fnd, out, fmt.Errorf("failed to find device %s: %w", *in.Name, err)
	}
	d := json.NewDecoder(bytes.NewReader(rsp.Body))
	var r nb.PaginatedDeviceWithConfigContextList
	err = d.Decode(&r)
	check(err)

	if *r.Count != 0 {
		fnd = true
		slc := *r.Results
		fmt.Printf("ID: %v\n", *slc[0].Id)
		fmt.Printf("DeviceRole: %v\n", *slc[0].DeviceRole.Id)
		fmt.Printf("DeviceType: %v\n", *slc[0].DeviceType.Id)
		fmt.Printf("Site: %v\n", *slc[0].Site.Id)

		out = &nb.WritableDeviceWithConfigContext{
			Name:       in.Name,
			Id:         slc[0].Id,
			DeviceRole: *slc[0].DeviceRole.Id,
			DeviceType: *slc[0].DeviceType.Id,
			Site:       *slc[0].Site.Id,
			Status:     nb.WritableDeviceWithConfigContextStatusEnumActive,
			Tags:       slc[0].Tags,
		}
		return fnd, out, nil
	}
	find, dr, err := findDeviceRole(n, *in.DeviceRole.Slug)
	if err != nil || !find {
		return fnd, out, fmt.Errorf("failed to find device role id for %s: %w", *in.Name, err)
	}
	find, dt, err := findDeviceType(n, in.DeviceType.Slug)
	if err != nil || !find {
		return fnd, out, fmt.Errorf("failed to find device type id for %s: %w", *in.Name, err)
	}
	find, st, err := findSite(n, *in.Site.Slug)
	if err != nil || !find {
		return fnd, out, fmt.Errorf("failed to find site id for %s: %w", *in.Name, err)
	}
	out = &nb.WritableDeviceWithConfigContext{
		Name:       in.Name,
		DeviceRole: dr,
		DeviceType: dt,
		Site:       st,
		Status:     nb.WritableDeviceWithConfigContextStatusEnumActive,
		Tags:       &[]nb.TagSerializerField{},
	}
	return fnd, out, nil
}

func findDeviceRole(n *nb.ClientWithResponses, slug string) (fnd bool, id types.UUID, err error) {
	rsp, err := n.DcimDeviceRolesListWithResponse(
		context.TODO(),
		&nb.DcimDeviceRolesListParams{
			SlugIe: &[]string{slug},
		})
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find device role %v: %w", &slug, err)
	}
	d := json.NewDecoder(bytes.NewReader(rsp.Body))
	var r nb.PaginatedDeviceRoleList
	err = d.Decode(&r)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to decode response finding device role %v: %w", &slug, err)
	}

	if *r.Count != 0 {
		fnd = true

		slc := *r.Results
		fmt.Printf("Device-Role ID: %v\n", *slc[0].Id)
		id = *slc[0].Id
	}
	return fnd, id, nil
}

func findDeviceType(n *nb.ClientWithResponses, slug string) (fnd bool, id types.UUID, err error) {
	rsp, err := n.DcimDeviceTypesListWithResponse(
		context.TODO(),
		&nb.DcimDeviceTypesListParams{
			SlugIe: &[]string{slug},
		})
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find device type %v: %w", &slug, err)
	}
	d := json.NewDecoder(bytes.NewReader(rsp.Body))
	var r nb.PaginatedDeviceTypeList
	err = d.Decode(&r)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to decode response finding device type %v: %w", &slug, err)
	}

	if *r.Count != 0 {
		fnd = true

		slc := *r.Results
		fmt.Printf("Device-Type ID: %v\n", *slc[0].Id)
		id = *slc[0].Id
	}
	return fnd, id, nil
}

func findSite(n *nb.ClientWithResponses, slug string) (fnd bool, id types.UUID, err error) {
	rsp, err := n.DcimSitesListWithResponse(
		context.TODO(),
		&nb.DcimSitesListParams{
			SlugIe: &[]string{slug},
		})
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find site %v: %w", &slug, err)
	}
	d := json.NewDecoder(bytes.NewReader(rsp.Body))

	var r nb.PaginatedSiteList
	err = d.Decode(&r)
	// FIX issue with empty emails "failed to pass regex validation"
	if err != nil {
		return fnd, id, fmt.Errorf("failed to decode response finding site %v: %w", &slug, err)
	}

	if *r.Count != 0 {
		fnd = true

		slc := *r.Results
		fmt.Printf("Site-ID: %v\n", *slc[0].Id)
		id = *slc[0].Id
	}
	return fnd, id, nil
}

func main() {
	token, err := NewSecurityProviderNautobotToken("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	check(err)

	c, err := nb.NewClientWithResponses(
		"https://develop.demo.nautobot.com/api/",
		nb.WithRequestEditorFn(token.Intercept),
	)
	check(err)

	////////////////////////////////
	// Read new device parameters
	////////////////////////////////
	dev, err := os.Open("device.json")
	check(err)
	defer dev.Close()

	d := json.NewDecoder(dev)

	var device nb.Device
	err = d.Decode(&device)
	check(err)

	////////////////////////////////
	// Check if device exists already
	////////////////////////////////
	found, devWithIDs, err := getDeviceIDs(c, device)
	check(err)

	if found {
		fmt.Println("Device already present")
		return
	}

	////////////////////////////////
	// Create device
	////////////////////////////////
	created, err := c.DcimDevicesCreateWithResponse(
		context.TODO(),
		nb.DcimDevicesCreateJSONRequestBody(*devWithIDs))
	check(err)

	fmt.Printf("%v", string(created.Body))
}
