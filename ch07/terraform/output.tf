data "nautobot_manufacturers" "all" {
  depends_on = [nautobot_manufacturer.new]
}

variable "manufacturer_name" {
  type    = string
  default = "New Vendor"
}

# Only returns te created manufacturer
output "data_source_example" {
  value = {
    for manufacturer in data.nautobot_manufacturers.all.manufacturers :
    manufacturer.id => manufacturer
    if manufacturer.name == var.manufacturer_name
  }
}
