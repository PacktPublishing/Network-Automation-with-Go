package main

import (
	"context"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	nautobot "github.com/nautobot/go-nautobot"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// bearer, err := securityprovider.NewSecurityProviderBearerToken("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	// check(err)

	// c, err := nautobot.NewClientWithResponses(
	// 	"https://develop.demo.nautobot.com/api/",
	// 	nautobot.WithRequestEditorFn(bearer.Intercept),
	// )
	// check(err)

	basicAuth, err := securityprovider.NewSecurityProviderBasicAuth("demo", "nautobot")
	check(err)

	c, err := nautobot.NewClientWithResponses(
		"https://develop.demo.nautobot.com/api/",
		nautobot.WithRequestEditorFn(basicAuth.Intercept),
	)
	check(err)

	ctx := context.Background()

	resp, err := c.DcimManufacturersListWithResponse(ctx, &nautobot.DcimManufacturersListParams{})
	check(err)

	fmt.Printf("%v", string(resp.Body))
}
