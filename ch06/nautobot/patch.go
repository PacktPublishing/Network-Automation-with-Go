package main

import (
	"fmt"
	"context"
	"net/http"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/types"
	nb "github.com/nautobot/go-nautobot"
)

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

// PaginatedSiteList defines model for PaginatedSiteList.
type PaginatedSiteList struct {
	Count    *int    `json:"count,omitempty"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  *[]Site `json:"results,omitempty"`
}

// Mixin to add `status` choice field to model serializers.
type Site struct {
	// 32-bit autonomous system number
	Asn          *int64                `json:"asn"`
	CircuitCount *int                  `json:"circuit_count,omitempty"`
	Comments     *string               `json:"comments,omitempty"`
	ContactEmail *string               `json:"contact_email,omitempty"`
	ContactName  *string               `json:"contact_name,omitempty"`
	ContactPhone *string               `json:"contact_phone,omitempty"`
	Created      *types.Date           `json:"created,omitempty"`
	CustomFields *nb.Site_CustomFields `json:"custom_fields,omitempty"`
	Description  *string               `json:"description,omitempty"`
	DeviceCount  *int                  `json:"device_count,omitempty"`

	// Human friendly display value
	Display *string `json:"display,omitempty"`

	// Local facility ID or description
	Facility    *string     `json:"facility,omitempty"`
	Id          *types.UUID `json:"id,omitempty"`
	LastUpdated *time.Time  `json:"last_updated,omitempty"`

	// GPS coordinate (latitude)
	Latitude *string `json:"latitude"`

	// GPS coordinate (longitude)
	Longitude       *string `json:"longitude"`
	Name            string  `json:"name"`
	PhysicalAddress *string `json:"physical_address,omitempty"`
	PrefixCount     *int    `json:"prefix_count,omitempty"`
	RackCount       *int    `json:"rack_count,omitempty"`
	Region          *struct {
		// Embedded struct due to allOf(#/components/schemas/NestedRegion)
		nb.NestedRegion `yaml:",inline"`
	} `json:"region"`
	ShippingAddress *string `json:"shipping_address,omitempty"`
	Slug            *string `json:"slug,omitempty"`
	Status          struct {
		Label *nb.SiteStatusLabel `json:"label,omitempty"`
		Value *nb.SiteStatusValue `json:"value,omitempty"`
	} `json:"status"`
	Tags   *[]nb.TagSerializerField `json:"tags,omitempty"`
	Tenant *struct {
		// Embedded struct due to allOf(#/components/schemas/NestedTenant)
		nb.NestedTenant `yaml:",inline"`
	} `json:"tenant"`
	TimeZone            *string `json:"time_zone"`
	Url                 *string `json:"url,omitempty"`
	VirtualmachineCount *int    `json:"virtualmachine_count,omitempty"`
	VlanCount           *int    `json:"vlan_count,omitempty"`
}