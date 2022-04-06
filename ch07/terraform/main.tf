terraform {
  required_providers {
    netbox = {
      source = "e-breuninger/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://demo.netbox.dev"
  api_token = "0123456789abcdef0123456789abcdef01234567"
}

resource "netbox_platform" "eos" {
  name = "Arista cEOS"
}

resource "netbox_platform" "srl" {
  name = "Nokia SR Linux" 
}