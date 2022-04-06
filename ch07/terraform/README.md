# Terraform

```
terraform init
terraform plan
terraform apply
```

```bash
⇨  terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions
are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # netbox_platform.eos will be created
  + resource "netbox_platform" "eos" {
      + id   = (known after apply)
      + name = "Arista cEOS"
      + slug = (known after apply)
    }

  # netbox_platform.srl will be created
  + resource "netbox_platform" "srl" {
      + id   = (known after apply)
      + name = "Nokia SR Linux"
      + slug = (known after apply)
    }

Plan: 2 to add, 0 to change, 0 to destroy.
╷
│ Warning: Possibly unsupported Netbox version
│ 
│   with provider["registry.terraform.io/e-breuninger/netbox"],
│   on main.tf line 9, in provider "netbox":
│    9: provider "netbox" {
│ 
│ This provider was tested against Netbox v3.1.3. Your Netbox version is v3.2.0. Unexpected errors
│ may occur.
╵

───────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take
exactly these actions if you run "terraform apply" now.
```

- [API Token](https://demo.netbox.dev/user/api-tokens/)
- [Manufacturer support](https://github.com/e-breuninger/terraform-provider-netbox/pull/142)