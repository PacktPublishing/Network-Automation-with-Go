package main

import (
	"context"
	"fmt"
	"os"

	"github.com/karimra/gnmic/api"
	"github.com/karimra/gnmic/target"
	"google.golang.org/protobuf/encoding/prototext"
	"gopkg.in/yaml.v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Router struct {
	Hostname   string `yaml:"hostname"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	SkipVerify bool   `yaml:"skipverify"`
}

type Inventory struct {
	Routers []Router `yaml:"router"`
}

func (r Router) createTarget() (*target.Target, error) {
	return api.NewTarget(
		api.Name("gnmi"),
		api.Address(r.Hostname+":"+r.Port),
		api.Username(r.Username),
		api.Password(r.Password),
		api.SkipVerify(r.SkipVerify),
	)
}

// export ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.18
func main() {
	src, err := os.Open("input.yml")
	check(err)
	defer src.Close()

	d := yaml.NewDecoder(src)

	var inv Inventory
	err = d.Decode(&inv)
	check(err)

	for _, router := range inv.Routers {
		////////////////////////////////
		// Create a target
		////////////////////////////////
		tg, err := router.createTarget()
		check(err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		////////////////////////////////
		// Create a gNMI client
		////////////////////////////////
		err = tg.CreateGNMIClient(ctx)
		check(err)
		defer tg.Close()

		////////////////////////////////
		// Send a gNMI capabilities request to the created target
		////////////////////////////////
		capResp, err := tg.Capabilities(ctx)
		check(err)

		fmt.Println(prototext.Format(capResp))
	}
}
