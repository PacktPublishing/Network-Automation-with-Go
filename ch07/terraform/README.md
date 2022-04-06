# Terraform

```
terraform init
terraform plan
terraform apply -auto-approve
```

```bash
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # netbox_device_role.container will be created
  + resource "netbox_device_role" "container" {
      + color_hex = "ff0000"
      + id        = (known after apply)
      + name      = "Container Router"
      + slug      = "container"
      + vm_role   = true
    }

  # netbox_platform.ceos will be created
  + resource "netbox_platform" "ceos" {
      + id   = (known after apply)
      + name = "Arista EOS"
      + slug = "ceos"
    }

  # netbox_platform.cvx will be created
  + resource "netbox_platform" "cvx" {
      + id   = (known after apply)
      + name = "NVIDIA Cumulus Linux"
      + slug = "cvx"
    }

  # netbox_platform.srl will be created
  + resource "netbox_platform" "srl" {
      + id   = (known after apply)
      + name = "Nokia SR Linux"
      + slug = "srl"
    }

Plan: 4 to add, 0 to change, 0 to destroy.
```

## Links

- [API Token](https://demo.netbox.dev/user/api-tokens/)
- [Manufacturer support](https://github.com/e-breuninger/terraform-provider-netbox/pull/142)

### Issues

- [terraform-provider-netbox](https://github.com/e-breuninger/terraform-provider-netbox/issues/145)

### Explore

:-)

- https://networkop.co.uk/post/2019-04-tf-yang/
- https://github.com/networkop/terraform-yang