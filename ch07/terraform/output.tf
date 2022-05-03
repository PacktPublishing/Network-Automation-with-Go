data "nautobot_manufacturers" "all" {}

variable "manufacturer_name" {
  type    = string
  default = "Juniper"
}

# Only returns Juniper manufacturer
output "juniper" {
  value = {
    for manufacturer in data.nautobot_manufacturers.all.manufacturers :
    manufacturer.id => manufacturer
    if manufacturer.name == var.manufacturer_name
  }
}