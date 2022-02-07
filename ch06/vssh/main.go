package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yahoo/vssh"
	"golang.org/x/crypto/ssh"
)

type Device struct {
	username string
	password string
	vendor   string
	cmd      string
}

var urlTemplate string = "%s:22"

func buildConfig(username, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(
				func(user, instruction string, questions []string, echos []bool) ([]string, error) {
					answers := make([]string, len(questions))
					for i := range answers {
						answers[i] = password
					}

					return answers, nil
				},
			),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func main() {
	cvxHost := flag.String("cvx-host", "clab-netgo-cvx", "CVX Hostname")
	cvxUser := flag.String("cvx-user", "cumulus", "CVX Username")
	cvxPass := flag.String("cvx-pass", "cumulus", "CVX password")
	srlHost := flag.String("srl-host", "clab-netgo-srl", "SRL Hostname")
	srlUser := flag.String("srl-user", "admin", "SRL Username")
	srlPass := flag.String("srl-pass", "admin", "SRL password")
	ceosHost := flag.String("ceos-host", "clab-netgo-ceos", "CEOS Hostname")
	ceosUser := flag.String("ceos-user", "admin", "CEOS Username")
	ceosPass := flag.String("ceos-pass", "admin", "CEOS password")
	flag.Parse()

	devices := map[string]Device{
		fmt.Sprintf(urlTemplate, *cvxHost): {
			*cvxUser,
			*cvxPass,
			"nvidia",
			"nv show system",
		},
		fmt.Sprintf(urlTemplate, *srlHost): {
			*srlUser,
			*srlPass,
			"nokia",
			"show version",
		},
		fmt.Sprintf(urlTemplate, *ceosHost): {
			username: *ceosUser,
			password: *ceosPass,
			vendor:   "arista",
			cmd:      "show version",
		},
	}

	vs := vssh.New().Start()

	for url, device := range devices {
		vs.AddClient(
			url,
			buildConfig(device.username, device.password),
			vssh.SetLabels(map[string]string{
				"NOS": device.vendor,
			}),
		)
	}
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, device := range devices {
		func(name string, device Device) {
			respChan, err := vs.RunWithLabel(
				ctx,
				device.cmd,
				fmt.Sprintf("NOS == %s", device.vendor),
				10*time.Second,
			)
			if err != nil {
				log.Fatal(err)
			}

			for resp := range respChan {
				if err := resp.Err(); err != nil {
					log.Println("error response: ", err)
					continue
				}

				outTxt, errTxt, _ := resp.GetText(vs)
				fmt.Printf("== Displaying version for device: %s ==\n", name)
				fmt.Println(outTxt, errTxt)
			}
		}(name, device)
	}
}
