package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/types"
	nb "github.com/nautobot/go-nautobot"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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

func createDeviceRole(n *nb.ClientWithResponses, dr nb.DeviceRole) error {
	_, err := n.DcimDeviceRolesCreate(
		context.TODO(),
		nb.DcimDeviceRolesCreateJSONRequestBody(dr))
	if err != nil {
		return fmt.Errorf("failed to create device role %s: %w", *dr.Display, err)
	}
	return nil
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

func createDeviceType(n *nb.ClientWithResponses, dt nb.DeviceType) error {
	man := nb.Manufacturer{
		Display: dt.Manufacturer.Display,
		Name:    dt.Manufacturer.Name,
		Slug:    dt.Manufacturer.Slug,
	}

	found, id, err := findManufacturer(n, *man.Slug)
	if err != nil || !found {
		return fmt.Errorf("error finding manufacturer %s: %w", *man.Display, err)
	}

	ndt := nb.WritableDeviceType{
		Manufacturer: id,
		Display:      dt.Display,
		Model:        dt.Model,
		Slug:         dt.Slug,
		Tags:         &[]nb.TagSerializerField{},
	}

	_, err = n.DcimDeviceTypesCreate(
		context.TODO(),
		nb.DcimDeviceTypesCreateJSONRequestBody(ndt))
	if err != nil {
		return fmt.Errorf("failed to create device type %s: %w", dt.Model, err)
	}
	return nil
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

	var r PaginatedSiteList
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

func createSite(n *nb.ClientWithResponses, s nb.Site) error {
	email := types.Email("contact@example.org")

	ns := nb.WritableSite{
		Name:         s.Name,
		Display:      s.Display,
		Slug:         s.Slug,
		ContactEmail: &email,
	}

	_, err := n.DcimSitesCreate(
		context.TODO(),
		nb.DcimSitesCreateJSONRequestBody(ns))
	if err != nil {
		return fmt.Errorf("failed to create site %s: %w", *ns.Display, err)
	}
	return nil
}

func findManufacturer(n *nb.ClientWithResponses, slug string) (fnd bool, id types.UUID, err error) {
	rsp, err := n.DcimManufacturersListWithResponse(
		context.TODO(),
		&nb.DcimManufacturersListParams{
			SlugIe: &[]string{slug},
		})
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find manufacturer %v: %w", &slug, err)
	}
	d := json.NewDecoder(bytes.NewReader(rsp.Body))
	var r nb.PaginatedManufacturerList
	err = d.Decode(&r)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to decode response finding manufacturer %v: %w", &slug, err)
	}

	if *r.Count != 0 {
		fnd = true

		slc := *r.Results
		fmt.Printf("Manufacturer ID: %v\n", *slc[0].Id)
		id = *slc[0].Id
	}
	return fnd, id, nil
}

func createManufacturer(n *nb.ClientWithResponses, m nb.Manufacturer) error {
	_, err := n.DcimManufacturersCreate(
		context.TODO(),
		nb.DcimManufacturersCreateJSONRequestBody(m))
	if err != nil {
		return fmt.Errorf("failed to create manufacturer %s: %w", m.Name, err)
	}
	return nil
}

func createResources(n *nb.ClientWithResponses) error {
	/////////////////////////////////
	// Manufacturers
	/////////////////////////////////
	man, err := os.Open("manufacturer.json")
	if err != nil {
		return fmt.Errorf("cannot open manufacturers file: %w", err)
	}

	d := json.NewDecoder(man)

	var manInput []nb.Manufacturer
	err = d.Decode(&manInput)
	if err != nil {
		return fmt.Errorf("cannot decode manufacturers data: %w", err)
	}

	for _, vendor := range manInput {
		found, _, err := findManufacturer(n, *vendor.Slug)
		if err != nil {
			return fmt.Errorf("error finding manufacturer %s: %w", *vendor.Display, err)
		}
		if !found {
			err = createManufacturer(n, vendor)
			if err != nil {
				return fmt.Errorf("error creating manufacturer %s: %w", *vendor.Display, err)
			}
		}
	}
	man.Close()

	/////////////////////////////////
	// Device Types
	/////////////////////////////////
	dev, err := os.Open("device-types.json")
	if err != nil {
		return fmt.Errorf("cannot open device types file: %w", err)
	}

	d = json.NewDecoder(dev)

	var devTypes []nb.DeviceType
	err = d.Decode(&devTypes)
	if err != nil {
		return fmt.Errorf("cannot decode device types data: %w", err)
	}

	for _, devType := range devTypes {
		found, _, err := findDeviceType(n, *devType.Slug)
		if err != nil {
			return fmt.Errorf("error finding device type %s: %w", *devType.Display, err)
		}
		if !found {
			err = createDeviceType(n, devType)
			if err != nil {
				return fmt.Errorf("error creating device type %s: %w", *devType.Display, err)
			}
		}
	}
	dev.Close()

	/////////////////////////////////
	// Device Role
	/////////////////////////////////
	rol, err := os.Open("device-roles.json")
	if err != nil {
		return fmt.Errorf("cannot open device roles file: %w", err)
	}

	d = json.NewDecoder(rol)

	var devRoles []nb.DeviceRole
	err = d.Decode(&devRoles)
	if err != nil {
		return fmt.Errorf("cannot decode device roles data: %w", err)
	}

	for _, devRole := range devRoles {
		found, _, err := findDeviceRole(n, *devRole.Slug)
		if err != nil {
			return fmt.Errorf("error finding device role %s: %w", *devRole.Display, err)
		}
		if !found {
			err = createDeviceRole(n, devRole)
			if err != nil {
				return fmt.Errorf("error creating device role %s: %w", *devRole.Display, err)
			}
		}
	}
	dev.Close()

	/////////////////////////////////
	// Sites
	/////////////////////////////////
	sit, err := os.Open("sites.json")
	if err != nil {
		return fmt.Errorf("cannot open sites file: %w", err)
	}
	defer dev.Close()

	d = json.NewDecoder(sit)

	var devSites []nb.Site
	err = d.Decode(&devSites)
	if err != nil {
		return fmt.Errorf("cannot decode sites data: %w", err)
	}

	for _, devSite := range devSites {
		found, _, err := findSite(n, *devSite.Slug)
		if err != nil {
			return fmt.Errorf("error finding site %s: %w", *devSite.Display, err)
		}
		if !found {
			err = createSite(n, devSite)
			if err != nil {
				return fmt.Errorf("error creating site %s: %w", *devSite.Display, err)
			}
		}
	}

	return nil
}

func main() {
	token, err := NewSecurityProviderNautobotToken("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	check(err)

	c, err := nb.NewClientWithResponses(
		"https://demo.nautobot.com/api/",
		nb.WithRequestEditorFn(token.Intercept),
	)
	check(err)

	err = createResources(c)
	check(err)

	/////////////////////////////////
	// Read new device parameters
	/////////////////////////////////
	dev, err := os.Open("device.json")
	check(err)
	defer dev.Close()

	d := json.NewDecoder(dev)

	var device nb.Device
	err = d.Decode(&device)
	check(err)

	/////////////////////////////////
	// Check if device exists already
	/////////////////////////////////
	found, devWithIDs, err := getDeviceIDs(c, device)
	check(err)

	if found {
		res, err := c.DcimDevicesListWithResponse(
			context.TODO(),
			&nb.DcimDevicesListParams{
				NameIe: &[]string{*device.Name},
			})
		check(err)
		fmt.Printf("Device already present: %v\n", string(res.Body))
		return
	}

	/////////////////////////////////
	// Create device
	/////////////////////////////////
	created, err := c.DcimDevicesCreateWithResponse(
		context.TODO(),
		nb.DcimDevicesCreateJSONRequestBody(*devWithIDs))
	check(err)

	fmt.Printf("Device created: %v\n", string(created.Body))
}
