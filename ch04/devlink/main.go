package main

import (
	"fmt"

	"github.com/mdlayher/devlink"
)

func main() {
	c, err := devlink.New()
	if err != nil {
		fmt.Println(err)

	}
	defer c.Close()

	devices, err := c.Devices()
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range devices {
		fmt.Printf("%+v", d)
	}

	ports, err := c.Ports()
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range ports {
		fmt.Printf("%+v", p)
	}
}
