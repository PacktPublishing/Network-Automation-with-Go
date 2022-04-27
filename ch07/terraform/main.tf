terraform {
  required_providers {
    nautobot = {
      version = "0.1"
      source  = "github.com/nleiva/nautobot"
    }
  }
}

variable "manufacturer_name" {
  type    = string
  default = "Juniper"
}

provider "nautobot" {
  url = "https://demo.nautobot.com/api/"
  token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}

data "nautobot_manufacturers" "all" {}

# Only returns Juniper manufacturer
output "juniper" {
  value = {
    for manufacturer in data.nautobot_manufacturers.all.manufacturers :
    manufacturer.id => manufacturer
    if manufacturer.name == var.manufacturer_name
  }
}