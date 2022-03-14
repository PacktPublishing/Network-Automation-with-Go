package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	xr "github.com/nleiva/xrgrpc"
	"github.com/tidwall/gjson"
)

type Router struct {
	Hostname  string `yaml:"hostname"`
	Platform  string `yaml:"platform"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	StrictKey bool   `yaml:"strictkey"`
}

type Inventory struct {
	Routers []Router `yaml:"router"`
}

type Config struct {
	Device    string
	Running   string
	Timestamp time.Time
}

func grpcConfig(r Router, paths string) (c Config) {
	var response string
	port := ":57777"

	// Manually specify target parameters.
	router, err := xr.BuildRouter(
		xr.WithUsername(r.Username),
		xr.WithPassword(r.Password),
		xr.WithHost(r.Hostname + port),
		xr.WithTimeout(5),
	)
	if err != nil {
		fmt.Printf("could not build a router, %v", err)
		return
	}
	if router.Cert != "" {
		fmt.Printf("Not empty: %v", router.Cert)
		return
	}

	conn, ctx, err := xr.Connect(*router)
	if err != nil {
		fmt.Printf("could not setup a client connection to %s, %v", router.Host, err)
		return
	}
	defer conn.Close()

	// Get config for the YANG paths
	if paths != "" {
		file := "yangocpaths.json"
		f, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("could not read file: %v: %v\n", file, err)
			return
		}
		paths = string(f)
	} 
	var id int64 = 1
	response, err = xr.GetConfig(ctx, conn, paths, id)
	if err != nil {
		fmt.Printf("could not get the config from %s, %v", router.Host, err)
		return
	}

	return Config{
		Device:    r.Hostname,
		Running:   response,
		Timestamp: time.Now(),
	}
}

func main() {
	src, err := os.Open("input.yml")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	d := yaml.NewDecoder(src)

	var inv Inventory
	err = d.Decode(&inv)
	if err != nil {
		panic(err)
	}

	var out Config
	// Empty paths forces to read from a local file with default paths.
	paths := ""

    out = grpcConfig(inv.Routers[0], paths)
	//fmt.Printf("Config from %s:\n%s\n", inv.Routers[0].Hostname, out.Running)

	scoped := gjson.Get(out.Running, "data.openconfig-system\\:system.grpc-server")
	fmt.Printf("Scoped config from %s:\n%s\n", inv.Routers[0].Hostname, scoped)
}
