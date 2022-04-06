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

resource "netbox_platform" "ceos" {
  name = "Arista EOS"
  slug = "ceos"
}

resource "netbox_platform" "srl" {
  name = "Nokia SR Linux"
  slug = "srl"
}

resource "netbox_platform" "cvx" {
  name = "NVIDIA Cumulus Linux"
  slug = "cvx"
}

resource "netbox_device_role" "container" {
  name      = "Container Router"
  vm_role   = true
  slug      = "container"
  color_hex = "ff0000"
}